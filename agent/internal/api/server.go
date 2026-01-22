package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aetherpanel/aether-panel/agent/internal/config"
	"github.com/aetherpanel/aether-panel/agent/internal/server"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
)

// Server represents the agent API server
type Server struct {
	app     *fiber.App
	config  *config.Config
	manager *server.Manager
	logger  *zap.Logger
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, manager *server.Manager, log *zap.Logger) *Server {
	app := fiber.New(fiber.Config{
		AppName:               "Aether Agent",
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		IdleTimeout:           120 * time.Second,
		DisableStartupMessage: true,
	})

	app.Use(logger.New())

	s := &Server{
		app:     app,
		config:  cfg,
		manager: manager,
		logger:  log,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	// Health check
	s.app.Get("/health", s.healthCheck)

	// API routes with token authentication
	api := s.app.Group("/api", s.authMiddleware)

	// Server management
	api.Post("/servers", s.createServer)
	api.Get("/servers/:id", s.getServer)
	api.Delete("/servers/:id", s.deleteServer)

	// Power actions
	api.Post("/servers/:id/power/start", s.startServer)
	api.Post("/servers/:id/power/stop", s.stopServer)
	api.Post("/servers/:id/power/restart", s.restartServer)
	api.Post("/servers/:id/power/kill", s.killServer)

	// Console
	api.Post("/servers/:id/command", s.sendCommand)
	api.Get("/servers/:id/logs", s.getLogs)

	// Stats
	api.Get("/servers/:id/stats", s.getStats)

	// System info
	api.Get("/system", s.getSystemInfo)

	// WebSocket for console streaming
	s.app.Get("/ws/console/:id", websocket.New(s.consoleWebSocket))
}

// authMiddleware validates the daemon token
func (s *Server) authMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		token = c.Query("token")
	}

	expectedToken := "Bearer " + s.config.Token
	if token != expectedToken && token != s.config.Token {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	return c.Next()
}

// healthCheck returns agent health status
func (s *Server) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"node_id": s.config.NodeID,
		"version": "1.0.0",
	})
}

// createServer creates a new server
func (s *Server) createServer(c *fiber.Ctx) error {
	var cfg server.ServerConfig
	if err := c.BodyParser(&cfg); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := s.manager.CreateServer(c.Context(), &cfg); err != nil {
		s.logger.Error("Failed to create server", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Server created",
	})
}

// getServer returns server information
func (s *Server) getServer(c *fiber.Ctx) error {
	serverID := c.Params("id")

	status, err := s.manager.GetServerStatus(c.Context(), serverID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Server not found",
		})
	}

	stats, _ := s.manager.GetServerStats(c.Context(), serverID)

	return c.JSON(fiber.Map{
		"id":     serverID,
		"status": status,
		"stats":  stats,
	})
}

// deleteServer removes a server
func (s *Server) deleteServer(c *fiber.Ctx) error {
	serverID := c.Params("id")

	if err := s.manager.DeleteServer(c.Context(), serverID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Server deleted",
	})
}

// startServer starts a server
func (s *Server) startServer(c *fiber.Ctx) error {
	serverID := c.Params("id")

	if err := s.manager.StartServer(c.Context(), serverID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"status":  "starting",
	})
}

// stopServer stops a server
func (s *Server) stopServer(c *fiber.Ctx) error {
	serverID := c.Params("id")

	if err := s.manager.StopServer(c.Context(), serverID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"status":  "stopping",
	})
}

// restartServer restarts a server
func (s *Server) restartServer(c *fiber.Ctx) error {
	serverID := c.Params("id")

	if err := s.manager.RestartServer(c.Context(), serverID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"status":  "restarting",
	})
}

// killServer forcefully stops a server
func (s *Server) killServer(c *fiber.Ctx) error {
	serverID := c.Params("id")

	if err := s.manager.KillServer(c.Context(), serverID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"status":  "stopped",
	})
}

// sendCommand sends a command to server console
func (s *Server) sendCommand(c *fiber.Ctx) error {
	serverID := c.Params("id")

	var req struct {
		Command string `json:"command"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := s.manager.SendCommand(c.Context(), serverID, req.Command); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

// getLogs returns server logs
func (s *Server) getLogs(c *fiber.Ctx) error {
	serverID := c.Params("id")
	_ = serverID

	// TODO: Implement log retrieval
	return c.JSON(fiber.Map{
		"logs": []string{},
	})
}

// getStats returns server statistics
func (s *Server) getStats(c *fiber.Ctx) error {
	serverID := c.Params("id")

	stats, err := s.manager.GetServerStats(c.Context(), serverID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Server not found",
		})
	}

	return c.JSON(stats)
}

// getSystemInfo returns node system information
func (s *Server) getSystemInfo(c *fiber.Ctx) error {
	// TODO: Implement system info collection
	return c.JSON(fiber.Map{
		"node_id":    s.config.NodeID,
		"cpu_cores":  0,
		"memory_mb":  0,
		"disk_mb":    0,
		"os":         "linux",
		"docker":     "running",
	})
}

// consoleWebSocket handles WebSocket connections for console streaming
func (s *Server) consoleWebSocket(c *websocket.Conn) {
	serverID := c.Params("id")
	s.logger.Info("Console WebSocket connected", zap.String("server", serverID))

	defer func() {
		s.logger.Info("Console WebSocket disconnected", zap.String("server", serverID))
		c.Close()
	}()

	// TODO: Implement console streaming via Docker attach
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		// Send command to server
		_ = s.manager.SendCommand(context.Background(), serverID, string(msg))
	}
}

// RegisterWithPanel registers this node with the panel
func (s *Server) RegisterWithPanel(ctx context.Context) error {
	client := &http.Client{Timeout: 10 * time.Second}

	payload := map[string]interface{}{
		"node_id": s.config.NodeID,
		"token":   s.config.Token,
		"port":    s.config.API.Port,
	}

	data, _ := json.Marshal(payload)
	_ = data
	_ = client

	// TODO: Implement panel registration
	s.logger.Info("Registered with panel", zap.String("url", s.config.Panel.URL))
	return nil
}

// Start starts the API server
func (s *Server) Start(addr string) error {
	if s.config.API.TLSCert != "" && s.config.API.TLSKey != "" {
		return s.app.ListenTLS(addr, s.config.API.TLSCert, s.config.API.TLSKey)
	}
	return s.app.Listen(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
