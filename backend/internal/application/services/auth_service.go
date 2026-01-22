package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/aetherpanel/aether-panel/internal/domain/repositories"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account is locked")
	ErrAccountInactive    = errors.New("account is not active")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrInvalid2FACode     = errors.New("invalid 2FA code")
	ErrEmailNotVerified   = errors.New("email not verified")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
	auditRepo   repositories.AuditLogRepository
	config      *config.Config
}

// NewAuthService creates a new AuthService
func NewAuthService(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	auditRepo repositories.AuditLogRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		auditRepo:   auditRepo,
		config:      cfg,
	}
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// Claims represents JWT claims
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	RoleID   uuid.UUID `json:"role_id"`
	RoleName string    `json:"role_name"`
	jwt.RegisteredClaims
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	TwoFACode string `json:"two_fa_code"`
	IPAddress string `json:"-"`
	UserAgent string `json:"-"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	User         *entities.User `json:"user"`
	Tokens       *TokenPair     `json:"tokens"`
	Requires2FA  bool           `json:"requires_2fa"`
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if account is locked
	if user.IsLocked() {
		return nil, ErrAccountLocked
	}

	// Check if account is active
	if !user.IsActive() {
		return nil, ErrAccountInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Increment failed login count
		_ = s.userRepo.IncrementFailedLogin(ctx, user.ID)
		
		// Check if should lock account
		if user.FailedLoginCount+1 >= s.config.Security.MaxLoginAttempts {
			lockUntil := time.Now().Add(s.config.Security.LockoutDuration)
			user.LockedUntil = &lockUntil
			_ = s.userRepo.Update(ctx, user)
		}
		
		return nil, ErrInvalidCredentials
	}

	// Check 2FA if enabled
	if user.TwoFactorEnabled {
		if req.TwoFACode == "" {
			return &LoginResponse{Requires2FA: true}, nil
		}
		
		if !totp.Validate(req.TwoFACode, user.TwoFactorSecret) {
			return nil, ErrInvalid2FACode
		}
	}

	// Reset failed login count
	_ = s.userRepo.ResetFailedLogin(ctx, user.ID)

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID, req.IPAddress)

	// Generate tokens
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	session := &entities.Session{
		UserID:       user.ID,
		Token:        tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		LastActivity: time.Now(),
		ExpiresAt:    tokens.ExpiresAt,
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Log audit
	s.logAudit(ctx, user.ID, entities.AuditActionLogin, "user", &user.ID, req.IPAddress, req.UserAgent)

	// Clear sensitive data
	user.PasswordHash = ""
	user.TwoFactorSecret = ""

	return &LoginResponse{
		User:   user,
		Tokens: tokens,
	}, nil
}

// Logout invalidates a user session
func (s *AuthService) Logout(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, ip, ua string) error {
	if err := s.sessionRepo.Revoke(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	s.logAudit(ctx, userID, entities.AuditActionLogout, "user", &userID, ip, ua)
	return nil
}

// LogoutAll invalidates all user sessions
func (s *AuthService) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return s.sessionRepo.RevokeAllByUserID(ctx, userID)
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Find session by refresh token
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !session.IsValid() {
		return nil, ErrTokenExpired
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !user.IsActive() {
		return nil, ErrAccountInactive
	}

	// Generate new tokens
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update session
	session.Token = tokens.AccessToken
	session.RefreshToken = tokens.RefreshToken
	session.ExpiresAt = tokens.ExpiresAt
	session.LastActivity = time.Now()
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return tokens, nil
}

// ValidateToken validates an access token and returns claims
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// Enable2FA enables 2FA for a user
func (s *AuthService) Enable2FA(ctx context.Context, userID uuid.UUID) (string, string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	// Generate TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.config.Security.TwoFactorIssuer,
		AccountName: user.Email,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Store secret (not enabled yet until verified)
	user.TwoFactorSecret = key.Secret()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

// Verify2FA verifies and enables 2FA
func (s *AuthService) Verify2FA(ctx context.Context, userID uuid.UUID, code string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !totp.Validate(code, user.TwoFactorSecret) {
		return ErrInvalid2FACode
	}

	user.TwoFactorEnabled = true
	return s.userRepo.Update(ctx, user)
}

// Disable2FA disables 2FA for a user
func (s *AuthService) Disable2FA(ctx context.Context, userID uuid.UUID, password string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}

	user.TwoFactorEnabled = false
	user.TwoFactorSecret = ""
	return s.userRepo.Update(ctx, user)
}

// generateTokenPair generates access and refresh tokens
func (s *AuthService) generateTokenPair(user *entities.User) (*TokenPair, error) {
	now := time.Now()
	accessExpiry := now.Add(s.config.JWT.AccessExpiry)

	// Access token claims
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RoleID:   user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.JWT.Issuer,
			Subject:   user.ID.String(),
			Audience:  jwt.ClaimStrings{s.config.JWT.Audience},
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	if user.Role != nil {
		claims.RoleName = user.Role.Name
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, err
	}
	refreshToken := base64.URLEncoding.EncodeToString(refreshTokenBytes)

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry,
		TokenType:    "Bearer",
	}, nil
}

// logAudit logs an audit event
func (s *AuthService) logAudit(ctx context.Context, userID uuid.UUID, action entities.AuditAction, resource string, resourceID *uuid.UUID, ip, ua string) {
	log := &entities.AuditLog{
		UserID:     &userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IPAddress:  ip,
		UserAgent:  ua,
	}
	_ = s.auditRepo.Create(ctx, log)
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
