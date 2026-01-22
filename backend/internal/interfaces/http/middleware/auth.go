package middleware

import (
	"context"
	"strings"

	"github.com/aetherpanel/aether-panel/internal/infrastructure/config"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware handles authentication and authorization
type AuthMiddleware struct {
	config *config.Config
	redis  *redis.Client
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(cfg *config.Config, rdb *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
		redis:  rdb,
	}
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

// UserContext keys
const (
	UserIDKey      = "user_id"
	UsernameKey    = "username"
	EmailKey       = "email"
	RoleIDKey      = "role_id"
	RoleNameKey    = "role_name"
	PermissionsKey = "permissions"
)

// Authenticate validates JWT token and sets user context
func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	// Get token from header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authorization header",
		})
	}

	// Extract token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}
	tokenString := parts[1]

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return []byte(m.config.JWT.Secret), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// Check if token is blacklisted
	ctx := context.Background()
	blacklisted, _ := m.redis.Exists(ctx, "blacklist:"+tokenString)
	if blacklisted {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token has been revoked",
		})
	}

	// Get user permissions from cache
	permissions, err := m.getUserPermissions(ctx, claims.UserID)
	if err != nil {
		permissions = []string{}
	}

	// Set user context
	c.Locals(UserIDKey, claims.UserID)
	c.Locals(UsernameKey, claims.Username)
	c.Locals(EmailKey, claims.Email)
	c.Locals(RoleIDKey, claims.RoleID)
	c.Locals(RoleNameKey, claims.RoleName)
	c.Locals(PermissionsKey, permissions)

	return c.Next()
}

// RequirePermission checks if user has required permission
func (m *AuthMiddleware) RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Admin role bypasses all permission checks
		roleName, ok := c.Locals(RoleNameKey).(string)
		if ok && roleName == "admin" {
			return c.Next()
		}

		permissions, ok := c.Locals(PermissionsKey).([]string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		for _, p := range permissions {
			if p == permission || p == "*" {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}

// RequireRole checks if user has required role
func (m *AuthMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals(RoleNameKey).(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient role",
		})
	}
}

// RequireOwnership checks if user owns the resource or is admin
func (m *AuthMiddleware) RequireOwnership(getOwnerID func(*fiber.Ctx) (uuid.UUID, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Admin bypasses ownership check
		roleName, ok := c.Locals(RoleNameKey).(string)
		if ok && roleName == "admin" {
			return c.Next()
		}

		userID, ok := c.Locals(UserIDKey).(uuid.UUID)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		ownerID, err := getOwnerID(c)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Resource not found",
			})
		}

		if userID != ownerID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		return c.Next()
	}
}

// getUserPermissions retrieves user permissions from cache or database
func (m *AuthMiddleware) getUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	cacheKey := redis.BuildKey(redis.PrefixUser, userID.String(), "permissions")
	
	var permissions []string
	err := m.redis.GetJSON(ctx, cacheKey, &permissions)
	if err == nil {
		return permissions, nil
	}

	// If not in cache, would need to fetch from database
	// For now, return empty - actual implementation would query DB
	return []string{}, nil
}

// GetUserID extracts user ID from context
func GetUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	userID, ok := c.Locals(UserIDKey).(uuid.UUID)
	return userID, ok
}

// GetUsername extracts username from context
func GetUsername(c *fiber.Ctx) (string, bool) {
	username, ok := c.Locals(UsernameKey).(string)
	return username, ok
}

// GetRoleName extracts role name from context
func GetRoleName(c *fiber.Ctx) (string, bool) {
	roleName, ok := c.Locals(RoleNameKey).(string)
	return roleName, ok
}

// IsAdmin checks if current user is admin
func IsAdmin(c *fiber.Ctx) bool {
	roleName, ok := GetRoleName(c)
	return ok && roleName == "admin"
}
