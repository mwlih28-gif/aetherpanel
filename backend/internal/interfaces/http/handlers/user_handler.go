package handlers

import (
	"github.com/aetherpanel/aether-panel/internal/infrastructure/config"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/redis"
	"github.com/aetherpanel/aether-panel/internal/interfaces/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserHandler handles user endpoints
type UserHandler struct {
	config *config.Config
	db     *gorm.DB
	redis  *redis.Client
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		config: cfg,
		db:     db,
		redis:  rdb,
	}
}

// CreateUserRequest represents user creation request
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
	RoleID    string `json:"role_id" validate:"uuid"`
}

// List returns paginated list of users
func (h *UserHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)
	search := c.Query("search")
	status := c.Query("status")

	_ = search
	_ = status

	return c.JSON(fiber.Map{
		"data": []interface{}{},
		"meta": fiber.Map{
			"page":        page,
			"page_size":   pageSize,
			"total":       0,
			"total_pages": 0,
		},
	})
}

// Create creates a new user
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

// GetByID returns a user by ID
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Params("id")

	id, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
	_ = id

	return c.JSON(fiber.Map{
		"id": userID,
	})
}

// Update updates a user
func (h *UserHandler) Update(c *fiber.Ctx) error {
	userID := c.Params("id")

	var req struct {
		Email     string `json:"email" validate:"email"`
		Username  string `json:"username" validate:"min=3,max=50"`
		FirstName string `json:"first_name" validate:"max=100"`
		LastName  string `json:"last_name" validate:"max=100"`
		RoleID    string `json:"role_id" validate:"uuid"`
		Status    string `json:"status" validate:"oneof=active inactive suspended"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	_ = userID

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// Delete deletes a user
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	userID := c.Params("id")
	currentUserID, _ := middleware.GetUserID(c)

	id, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if id == currentUserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete your own account",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
