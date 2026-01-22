package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Node struct {
	ID                   string    `json:"id" gorm:"primaryKey"`
	UUID                 string    `json:"uuid" gorm:"unique;not null"`
	Name                 string    `json:"name" gorm:"not null"`
	Description          string    `json:"description"`
	LocationID           string    `json:"location_id" gorm:"not null"`
	Location             *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
	FQDN                 string    `json:"fqdn" gorm:"not null"`
	Hostname             string    `json:"hostname" gorm:"not null"`
	IP                   string    `json:"ip" gorm:"not null"`
	Scheme               string    `json:"scheme" gorm:"default:https"`
	BehindProxy          bool      `json:"behind_proxy" gorm:"default:false"`
	MaintenanceMode      bool      `json:"maintenance_mode" gorm:"default:false"`
	Memory               int       `json:"memory" gorm:"not null"`
	MemoryOverallocate   int       `json:"memory_overallocate" gorm:"default:0"`
	Disk                 int       `json:"disk" gorm:"not null"`
	DiskOverallocate     int       `json:"disk_overallocate" gorm:"default:0"`
	UploadSize           int       `json:"upload_size" gorm:"default:100"`
	DaemonListenPort     int       `json:"daemon_listen_port" gorm:"default:8080"`
	DaemonSftpPort       int       `json:"daemon_sftp_port" gorm:"default:2022"`
	DaemonToken          string    `json:"daemon_token" gorm:"not null"`
	PublicKey            bool      `json:"public_key" gorm:"default:true"`
	Status               string    `json:"status" gorm:"default:offline"`
	LastSeen             *time.Time `json:"last_seen"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

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
	var nodes []Node
	
	if err := h.db.Preload("Location").Find(&nodes).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch nodes",
		})
	}

	// Add mock resource data for now
	for i := range nodes {
		nodes[i].Status = "offline" // Will be updated by Wings daemon
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
	var location Location
	if err := h.db.Where("id = ?", req.LocationID).First(&location).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	// Check if FQDN already exists
	var existingNode Node
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

	node := Node{
		ID:                   uuid.New().String(),
		UUID:                 uuid.New().String(),
		Name:                 req.Name,
		Description:          req.Description,
		LocationID:           req.LocationID,
		FQDN:                 req.FQDN,
		Hostname:             req.FQDN,
		IP:                   "", // Will be resolved later
		Scheme:               req.Scheme,
		BehindProxy:          req.BehindProxy,
		MaintenanceMode:      false,
		Memory:               req.Memory,
		MemoryOverallocate:   req.MemoryOverallocate,
		Disk:                 req.Disk,
		DiskOverallocate:     req.DiskOverallocate,
		UploadSize:           req.UploadSize,
		DaemonListenPort:     req.DaemonListenPort,
		DaemonSftpPort:       req.DaemonSftpPort,
		DaemonToken:          token,
		PublicKey:            true,
		Status:               "offline",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
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
	
	var node Node
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

	var node Node
	if err := h.db.Where("id = ?", id).First(&node).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Node not found",
		})
	}

	// Check if location exists
	var location Location
	if err := h.db.Where("id = ?", req.LocationID).First(&location).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	// Check if FQDN already exists (excluding current node)
	var existingNode Node
	if err := h.db.Where("fqdn = ? AND id != ?", req.FQDN, id).First(&existingNode).Error; err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Node with this FQDN already exists",
		})
	}

	// Update node fields
	node.Name = req.Name
	node.Description = req.Description
	node.LocationID = req.LocationID
	node.FQDN = req.FQDN
	node.Hostname = req.FQDN
	node.Scheme = req.Scheme
	node.BehindProxy = req.BehindProxy
	node.MaintenanceMode = req.MaintenanceMode
	node.Memory = req.Memory
	node.MemoryOverallocate = req.MemoryOverallocate
	node.Disk = req.Disk
	node.DiskOverallocate = req.DiskOverallocate
	node.UploadSize = req.UploadSize
	node.DaemonListenPort = req.DaemonListenPort
	node.DaemonSftpPort = req.DaemonSftpPort
	node.UpdatedAt = time.Now()

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
	
	var node Node
	if err := h.db.Where("id = ?", id).First(&node).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Node not found",
		})
	}

	// Check if node has servers
	var serverCount int64
	if err := h.db.Model(&GameServer{}).Where("node_id = ?", id).Count(&serverCount).Error; err != nil {
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
	
	var node Node
	if err := h.db.Where("id = ?", id).First(&node).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Node not found",
		})
	}

	config := fiber.Map{
		"debug": false,
		"uuid":  node.UUID,
		"token_id": node.DaemonToken[:16],
		"token": node.DaemonToken,
		"api": fiber.Map{
			"host": node.FQDN,
			"port": node.DaemonListenPort,
			"ssl": fiber.Map{
				"enabled": node.Scheme == "https",
				"cert":    "/etc/letsencrypt/live/" + node.FQDN + "/fullchain.pem",
				"key":     "/etc/letsencrypt/live/" + node.FQDN + "/privkey.pem",
			},
		},
		"system": fiber.Map{
			"data": "/var/lib/pterodactyl/volumes",
			"sftp": fiber.Map{
				"bind_port": node.DaemonSftpPort,
			},
		},
		"allowed_mounts": []string{},
		"remote":         c.BaseURL(),
	}

	return c.JSON(config)
}
