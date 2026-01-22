package handlers

import (
	"net/http"
	"time"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateServerRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
	NodeID      string `json:"node_id" validate:"required,uuid"`
	Game        string `json:"game" validate:"required,oneof=minecraft rust csgo gmod ark valheim"`
	GameVersion string `json:"game_version" validate:"required"`
	Memory      int    `json:"memory" validate:"required,min=128"`
	Disk        int    `json:"disk" validate:"required,min=512"`
	CPU         int    `json:"cpu" validate:"required,min=50,max=400"`
}

type UpdateServerRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
	Memory      int    `json:"memory" validate:"required,min=128"`
	Disk        int    `json:"disk" validate:"required,min=512"`
	CPU         int    `json:"cpu" validate:"required,min=50,max=400"`
}

// GetServers returns all servers
func (h *Handler) GetServers(c *fiber.Ctx) error {
	var servers []entities.Server
	
	if err := h.db.Preload("Node").Preload("Node.Location").Find(&servers).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch servers",
		})
	}

	return c.JSON(fiber.Map{
		"data": servers,
	})
}

// CreateServer creates a new game server
func (h *Handler) CreateServer(c *fiber.Ctx) error {
	var req CreateServerRequest
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

	// Check if node exists
	var node entities.Node
	if err := h.db.Where("id = ?", req.NodeID).First(&node).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Node not found",
		})
	}

	// Check node resources
	var totalMemory, totalDisk int
	h.db.Model(&entities.Server{}).Where("node_id = ?", req.NodeID).Select("COALESCE(SUM(memory), 0), COALESCE(SUM(disk), 0)").Row().Scan(&totalMemory, &totalDisk)

	if totalMemory+req.Memory > node.Memory {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient memory on node",
		})
	}

	if totalDisk+req.Disk > node.Disk {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient disk space on node",
		})
	}

	// Get Docker image and startup command based on game
	dockerImage, startupCmd := getGameConfig(req.Game, req.GameVersion)

	server := entities.Server{
		ID:          uuid.New().String(),
		UUID:        uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		NodeID:      req.NodeID,
		Game:        req.Game,
		GameVersion: req.GameVersion,
		Status:      "stopped",
		Memory:      req.Memory,
		Disk:        req.Disk,
		CPU:         req.CPU,
		DockerImage: dockerImage,
		StartupCmd:  startupCmd,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.db.Create(&server).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create server",
		})
	}

	// Load relationships for response
	h.db.Preload("Node").Preload("Node.Location").First(&server, "id = ?", server.ID)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"data": server,
	})
}

// GetServer returns a specific server
func (h *Handler) GetServer(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var server entities.Server
	if err := h.db.Preload("Node").Preload("Node.Location").Where("id = ?", id).First(&server).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Server not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": server,
	})
}

// UpdateServer updates an existing server
func (h *Handler) UpdateServer(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var req UpdateServerRequest
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

	var server entities.Server
	if err := h.db.Where("id = ?", id).First(&server).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Server not found",
		})
	}

	// Update server fields
	server.Name = req.Name
	server.Description = req.Description
	server.Memory = req.Memory
	server.Disk = req.Disk
	server.CPU = req.CPU
	server.UpdatedAt = time.Now()

	if err := h.db.Save(&server).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update server",
		})
	}

	// Load relationships for response
	h.db.Preload("Node").Preload("Node.Location").First(&server, "id = ?", server.ID)

	return c.JSON(fiber.Map{
		"data": server,
	})
}

// DeleteServer deletes a server
func (h *Handler) DeleteServer(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var server entities.Server
	if err := h.db.Where("id = ?", id).First(&server).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Server not found",
		})
	}

	if server.Status == "running" {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Cannot delete running server. Stop it first.",
		})
	}

	if err := h.db.Delete(&server).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete server",
		})
	}

	return c.Status(http.StatusNoContent).Send(nil)
}

// StartServer starts a server
func (h *Handler) StartServer(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var server entities.Server
	if err := h.db.Where("id = ?", id).First(&server).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Server not found",
		})
	}

	if server.Status == "running" {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Server is already running",
		})
	}

	// TODO: Send start command to Wings daemon
	server.Status = "starting"
	server.UpdatedAt = time.Now()
	h.db.Save(&server)

	return c.JSON(fiber.Map{
		"message": "Server start command sent",
		"data":    server,
	})
}

// StopServer stops a server
func (h *Handler) StopServer(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var server entities.Server
	if err := h.db.Where("id = ?", id).First(&server).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Server not found",
		})
	}

	if server.Status == "stopped" {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Server is already stopped",
		})
	}

	// TODO: Send stop command to Wings daemon
	server.Status = "stopping"
	server.UpdatedAt = time.Now()
	h.db.Save(&server)

	return c.JSON(fiber.Map{
		"message": "Server stop command sent",
		"data":    server,
	})
}

// RestartServer restarts a server
func (h *Handler) RestartServer(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var server entities.Server
	if err := h.db.Where("id = ?", id).First(&server).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Server not found",
		})
	}

	// TODO: Send restart command to Wings daemon
	server.Status = "restarting"
	server.UpdatedAt = time.Now()
	h.db.Save(&server)

	return c.JSON(fiber.Map{
		"message": "Server restart command sent",
		"data":    server,
	})
}

// getGameConfig returns Docker image and startup command for a game
func getGameConfig(game, version string) (string, string) {
	switch game {
	case "minecraft":
		return "ghcr.io/pterodactyl/yolks:java_17", "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}"
	case "rust":
		return "ghcr.io/pterodactyl/games:rust", "./RustDedicated -batchmode +server.port {{SERVER_PORT}} +server.identity \"rust\" +rcon.port {{RCON_PORT}} +rcon.web true +server.hostname \"{{SERVER_NAME}}\" +server.level \"{{LEVEL}}\" +server.description \"{{DESCRIPTION}}\" +server.url \"{{SERVER_URL}}\" +server.headerimage \"{{SERVER_IMG}}\" +server.logoimage \"{{SERVER_LOGO}}\" +server.maxplayers {{MAX_PLAYERS}} +rcon.password \"{{RCON_PASS}}\" +server.saveinterval {{SAVE_INTERVAL}} {{ADDITIONAL_ARGS}}"
	case "csgo":
		return "ghcr.io/pterodactyl/games:source", "./srcds_run -game csgo -console -usercon +game_type 0 +game_mode 1 +mapgroup mg_active +map de_dust2 -tickrate {{TICKRATE}} -port {{SERVER_PORT}} +rcon_password \"{{RCON_PASSWORD}}\" +hostname \"{{HOSTNAME}}\" +sv_password \"{{SERVER_PASSWORD}}\""
	case "gmod":
		return "ghcr.io/pterodactyl/games:source", "./srcds_run -game garrysmod -console -usercon +gamemode {{GAMEMODE}} +map {{MAP}} -tickrate {{TICKRATE}} -port {{SERVER_PORT}} +hostname \"{{HOSTNAME}}\" +rcon_password \"{{RCON_PASSWORD}}\" +sv_password \"{{SERVER_PASSWORD}}\""
	case "ark":
		return "ghcr.io/pterodactyl/games:ark", "./Shooterentities.Server {{MAP}}?listen?SessionName=\"{{SESSION_NAME}}\"?ServerPassword={{SERVER_PASSWORD}}?ServerAdminPassword={{ADMIN_PASSWORD}}?Port={{SERVER_PORT}}?QueryPort={{QUERY_PORT}}?MaxPlayers={{MAX_PLAYERS}}"
	case "valheim":
		return "ghcr.io/pterodactyl/games:valheim", "./valheim_server.x86_64 -name \"{{SERVER_NAME}}\" -port {{SERVER_PORT}} -world \"{{WORLD_NAME}}\" -password \"{{PASSWORD}}\" -public {{PUBLIC_SERVER}}"
	default:
		return "ghcr.io/pterodactyl/yolks:debian", "/bin/bash"
	}
}
