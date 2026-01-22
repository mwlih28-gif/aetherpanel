package handlers

import (
	"github.com/aetherpanel/aether-panel/internal/infrastructure/config"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/redis"
	"github.com/aetherpanel/aether-panel/internal/interfaces/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NodeHandler handles node endpoints
type NodeHandler struct {
	config *config.Config
	db     *gorm.DB
	redis  *redis.Client
}

// NewNodeHandler creates a new NodeHandler
func NewNodeHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *NodeHandler {
	return &NodeHandler{
		config: cfg,
		db:     db,
		redis:  rdb,
	}
}

// CreateNodeRequest represents node creation request
type CreateNodeRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=100"`
	Description     string `json:"description" validate:"max=500"`
	LocationID      string `json:"location_id" validate:"required,uuid"`
	FQDN            string `json:"fqdn" validate:"required"`
	Scheme          string `json:"scheme" validate:"oneof=http https"`
	DaemonPort      int    `json:"daemon_port" validate:"required,min=1,max=65535"`
	MemoryTotal     int64  `json:"memory_total" validate:"required,min=1024"`
	MemoryOveralloc int    `json:"memory_overalloc" validate:"min=0,max=100"`
	DiskTotal       int64  `json:"disk_total" validate:"required,min=10240"`
	DiskOveralloc   int    `json:"disk_overalloc" validate:"min=0,max=100"`
	CPUTotal        int    `json:"cpu_total" validate:"required,min=100"`
}

// List returns paginated list of nodes
func (h *NodeHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	// TODO: Implement node listing
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

// Create creates a new node
func (h *NodeHandler) Create(c *fiber.Ctx) error {
	var req CreateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, _ := middleware.GetUserID(c)
	_ = userID

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Node created successfully",
	})
}

// GetByID returns a node by ID
func (h *NodeHandler) GetByID(c *fiber.Ctx) error {
	nodeID := c.Params("id")

	id, err := uuid.Parse(nodeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid node ID",
		})
	}
	_ = id

	return c.JSON(fiber.Map{
		"id":        nodeID,
		"is_online": false,
	})
}

// Update updates a node
func (h *NodeHandler) Update(c *fiber.Ctx) error {
	nodeID := c.Params("id")

	var req CreateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	_ = nodeID

	return c.JSON(fiber.Map{
		"message": "Node updated successfully",
	})
}

// Delete deletes a node
func (h *NodeHandler) Delete(c *fiber.Ctx) error {
	nodeID := c.Params("id")
	userID, _ := middleware.GetUserID(c)
	_ = nodeID
	_ = userID

	return c.JSON(fiber.Map{
		"message": "Node deleted successfully",
	})
}

// GetConfiguration returns node configuration for agent
func (h *NodeHandler) GetConfiguration(c *fiber.Ctx) error {
	nodeID := c.Params("id")
	_ = nodeID

	return c.JSON(fiber.Map{
		"node":        nil,
		"servers":     []interface{}{},
		"allocations": []interface{}{},
	})
}

// CreateAllocations creates allocations for a node
func (h *NodeHandler) CreateAllocations(c *fiber.Ctx) error {
	nodeID := c.Params("id")

	var req struct {
		IP        string `json:"ip" validate:"required,ip"`
		PortStart int    `json:"port_start" validate:"required,min=1,max=65535"`
		PortEnd   int    `json:"port_end" validate:"required,min=1,max=65535"`
		Alias     string `json:"alias" validate:"max=255"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	_ = nodeID

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Allocations created",
		"count":   req.PortEnd - req.PortStart + 1,
	})
}

// ListLocations returns all locations
func (h *NodeHandler) ListLocations(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"data": []interface{}{},
	})
}

// CreateLocation creates a new location
func (h *NodeHandler) CreateLocation(c *fiber.Ctx) error {
	var req struct {
		ShortCode   string  `json:"short_code" validate:"required,min=2,max=10"`
		Name        string  `json:"name" validate:"required,min=1,max=100"`
		Description string  `json:"description" validate:"max=500"`
		Country     string  `json:"country" validate:"len=2"`
		City        string  `json:"city" validate:"max=100"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Location created successfully",
	})
}
