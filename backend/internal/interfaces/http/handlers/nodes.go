package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateNodeRequest struct {
	Name                 string `json:"name" validate:"required,min=1,max=100"`
	Description          string `json:"description" validate:"max=500"`
	LocationID           string `json:"location_id" validate:"required,uuid"`
	FQDN                 string `json:"fqdn" validate:"required,fqdn"`
	Scheme               string `json:"scheme" validate:"required,oneof=http https"`
	BehindProxy          bool   `json:"behind_proxy"`
	Memory               int    `json:"memory" validate:"required,min=128"`
	MemoryOverallocate   int    `json:"memory_overallocate" validate:"min=0,max=500"`
	Disk                 int    `json:"disk" validate:"required,min=1024"`
	DiskOverallocate     int    `json:"disk_overallocate" validate:"min=0,max=500"`
	UploadSize           int    `json:"upload_size" validate:"min=1,max=1000"`
	DaemonListenPort     int    `json:"daemon_listen_port" validate:"required,min=1024,max=65535"`
	DaemonSftpPort       int    `json:"daemon_sftp_port" validate:"required,min=1024,max=65535"`
}

type UpdateNodeRequest struct {
	Name                 string `json:"name" validate:"required,min=1,max=100"`
	Description          string `json:"description" validate:"max=500"`
	LocationID           string `json:"location_id" validate:"required,uuid"`
	FQDN                 string `json:"fqdn" validate:"required,fqdn"`
	Scheme               string `json:"scheme" validate:"required,oneof=http https"`
	BehindProxy          bool   `json:"behind_proxy"`
	MaintenanceMode      bool   `json:"maintenance_mode"`
	Memory               int    `json:"memory" validate:"required,min=128"`
	MemoryOverallocate   int    `json:"memory_overallocate" validate:"min=0,max=500"`
	Disk                 int    `json:"disk" validate:"required,min=1024"`
	DiskOverallocate     int    `json:"disk_overallocate" validate:"min=0,max=500"`
	UploadSize           int    `json:"upload_size" validate:"min=1,max=1000"`
	DaemonListenPort     int    `json:"daemon_listen_port" validate:"required,min=1024,max=65535"`
	DaemonSftpPort       int    `json:"daemon_sftp_port" validate:"required,min=1024,max=65535"`
}

// generateToken generates a secure random token
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetNodes returns all nodes
func (h *Handler) GetNodes(c *fiber.Ctx) error {
	var nodes []entities.Node
	
	if err := h.db.Preload("Location").Find(&nodes).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch nodes",
		})
	}

	// Set status based on IsOnline field
	for i := range nodes {
		if nodes[i].IsOnline {
			// Status will be determined by Wings daemon, but for now set based on online status
		}
	}

	return c.JSON(fiber.Map{
		"data": nodes,
	})
}

// CreateNode creates a new node
func (h *Handler) CreateNode(c *fiber.Ctx) error {
	var req CreateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	// Check if location exists
	var location entities.Location
	if err := h.db.Where("id = ?", req.LocationID).First(&location).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	// Check if FQDN already exists
	var existingNode entities.Node
	if err := h.db.Where("fqdn = ?", req.FQDN).First(&existingNode).Error; err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Node with this FQDN already exists",
		})
	}

	// Generate daemon token
	token, err := generateToken()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate daemon token",
		})
	}

	locationUUID, _ := uuid.Parse(req.LocationID)
	
	node := entities.Node{
		Name:             req.Name,
		Description:      req.Description,
		LocationID:       locationUUID,
		FQDN:             req.FQDN,
		Scheme:           req.Scheme,
		DaemonPort:       req.DaemonListenPort,
		DaemonToken:      token,
		MemoryTotal:      int64(req.Memory),
		MemoryOveralloc:  req.MemoryOverallocate,
		DiskTotal:        int64(req.Disk),
		DiskOveralloc:    req.DiskOverallocate,
		CPUTotal:         100, // Default 1 core
		IsOnline:         false,
		MaintenanceMode:  req.BehindProxy, // Use BehindProxy as maintenance mode for now
	}

	if err := h.db.Create(&node).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create node",
		})
	}

	// Load location for response
	h.db.Preload("Location").First(&node, "id = ?", node.ID)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"data": node,
	})
}

// GetNode returns a specific node
func (h *Handler) GetNode(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var node entities.Node
	if err := h.db.Preload("Location").Where("id = ?", id).First(&node).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Node not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": node,
	})
}

// UpdateNode updates an existing node
func (h *Handler) UpdateNode(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var req UpdateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	var node entities.Node
	if err := h.db.Where("id = ?", id).First(&node).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Node not found",
		})
	}

	// Check if location exists
	var location entities.Location
	if err := h.db.Where("id = ?", req.LocationID).First(&location).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	// Check if FQDN already exists (excluding current node)
	var existingNode entities.Node
	if err := h.db.Where("fqdn = ? AND id != ?", req.FQDN, id).First(&existingNode).Error; err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Node with this FQDN already exists",
		})
	}

	// Update node fields
	locationUUID, _ := uuid.Parse(req.LocationID)
	
	node.Name = req.Name
	node.Description = req.Description
	node.LocationID = locationUUID
	node.FQDN = req.FQDN
	node.Scheme = req.Scheme
	node.MaintenanceMode = req.MaintenanceMode
	node.MemoryTotal = int64(req.Memory)
	node.MemoryOveralloc = req.MemoryOverallocate
	node.DiskTotal = int64(req.Disk)
	node.DiskOveralloc = req.DiskOverallocate
	node.DaemonPort = req.DaemonListenPort

	if err := h.db.Save(&node).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update node",
		})
	}

	// Load location for response
	h.db.Preload("Location").First(&node, "id = ?", node.ID)

	return c.JSON(fiber.Map{
		"data": node,
	})
}

// DeleteNode deletes a node
func (h *Handler) DeleteNode(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var node entities.Node
	if err := h.db.Where("id = ?", id).First(&node).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Node not found",
		})
	}

	// Check if node has servers
	var serverCount int64
	if err := h.db.Model(&entities.Server{}).Where("node_id = ?", id).Count(&serverCount).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check node usage",
		})
	}

	if serverCount > 0 {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Cannot delete node that has servers assigned to it",
		})
	}

	if err := h.db.Delete(&node).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete node",
		})
	}

	return c.Status(http.StatusNoContent).Send(nil)
}

// GetNodeConfiguration returns the Wings configuration for a node
func (h *Handler) GetNodeConfiguration(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var node entities.Node
	if err := h.db.Where("id = ?", id).First(&node).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Node not found",
		})
	}

	config := fiber.Map{
		"debug": false,
		"uuid":  node.ID.String(),
		"token_id": node.DaemonToken[:16],
		"token": node.DaemonToken,
		"api": fiber.Map{
			"host": node.FQDN,
			"port": node.DaemonPort,
			"ssl": fiber.Map{
				"enabled": node.Scheme == "https",
				"cert":    "/etc/letsencrypt/live/" + node.FQDN + "/fullchain.pem",
				"key":     "/etc/letsencrypt/live/" + node.FQDN + "/privkey.pem",
			},
		},
		"system": fiber.Map{
			"data": "/var/lib/pterodactyl/volumes",
			"sftp": fiber.Map{
				"bind_port": 2022, // Default SFTP port
			},
		},
		"allowed_mounts": []string{},
		"remote":         c.BaseURL(),
	}

	return c.JSON(config)
}
