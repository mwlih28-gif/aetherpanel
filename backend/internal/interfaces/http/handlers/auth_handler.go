package handlers

import (
	"github.com/aetherpanel/aether-panel/internal/infrastructure/config"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/redis"
	"github.com/aetherpanel/aether-panel/internal/interfaces/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	config *config.Config
	db     *gorm.DB
	redis  *redis.Client
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{
		config: cfg,
		db:     db,
		redis:  rdb,
	}
}

// LoginRequest represents login request body
type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	TwoFACode string `json:"two_fa_code"`
}

// RegisterRequest represents registration request body
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO: Implement actual login logic using AuthService
	// This is a placeholder response
	return c.JSON(fiber.Map{
		"message": "Login endpoint",
		"email":   req.Email,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO: Implement actual registration logic
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registration successful",
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO: Implement token refresh logic
	return c.JSON(fiber.Map{
		"message": "Token refreshed",
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// TODO: Invalidate session
	_ = userID

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// Me returns current user info
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	username, _ := middleware.GetUsername(c)
	roleName, _ := middleware.GetRoleName(c)

	return c.JSON(fiber.Map{
		"user_id":  userID,
		"username": username,
		"role":     roleName,
	})
}

// Enable2FA initiates 2FA setup
func (h *AuthHandler) Enable2FA(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// TODO: Generate TOTP secret and QR code
	_ = userID

	return c.JSON(fiber.Map{
		"secret":  "PLACEHOLDER_SECRET",
		"qr_code": "data:image/png;base64,...",
	})
}

// Verify2FA verifies and enables 2FA
func (h *AuthHandler) Verify2FA(c *fiber.Ctx) error {
	var req struct {
		Code string `json:"code" validate:"required,len=6"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// TODO: Verify TOTP code and enable 2FA
	_ = userID

	return c.JSON(fiber.Map{
		"message": "2FA enabled successfully",
	})
}

// Disable2FA disables 2FA
func (h *AuthHandler) Disable2FA(c *fiber.Ctx) error {
	var req struct {
		Password string `json:"password" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// TODO: Verify password and disable 2FA
	_ = userID

	return c.JSON(fiber.Map{
		"message": "2FA disabled successfully",
	})
}

// ForgotPassword initiates password reset
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO: Send password reset email
	return c.JSON(fiber.Map{
		"message": "If the email exists, a reset link has been sent",
	})
}

// ResetPassword resets user password
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token    string `json:"token" validate:"required"`
		Password string `json:"password" validate:"required,min=8"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO: Verify token and reset password
	return c.JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}

// parseUUID parses a UUID from string
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
