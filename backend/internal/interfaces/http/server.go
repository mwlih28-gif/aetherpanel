package http

import (
	"context"
	"time"

	"github.com/aetherpanel/aether-panel/internal/infrastructure/config"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/redis"
	"github.com/aetherpanel/aether-panel/internal/interfaces/http/handlers"
	"github.com/aetherpanel/aether-panel/internal/interfaces/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NewServer creates and configures a new Fiber server
func NewServer(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               cfg.App.Name,
		ReadTimeout:           cfg.Server.ReadTimeout,
		WriteTimeout:          cfg.Server.WriteTimeout,
		IdleTimeout:           cfg.Server.IdleTimeout,
		BodyLimit:             cfg.Server.BodyLimit * 1024 * 1024,
		DisableStartupMessage: cfg.App.Environment == "production",
		ErrorHandler:          errorHandler,
	})

	// Global middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.App.Debug,
	}))

	app.Use(requestid.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	app.Use(helmet.New())

	// CORS
	if cfg.Server.CORS.Enabled {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     joinStrings(cfg.Server.CORS.AllowOrigins),
			AllowMethods:     joinStrings(cfg.Server.CORS.AllowMethods),
			AllowHeaders:     joinStrings(cfg.Server.CORS.AllowHeaders),
			ExposeHeaders:    joinStrings(cfg.Server.CORS.ExposeHeaders),
			AllowCredentials: cfg.Server.CORS.AllowCredentials,
			MaxAge:           cfg.Server.CORS.MaxAge,
		}))
	}

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.Security.RateLimitRequests,
		Expiration: cfg.Security.RateLimitDuration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	}))

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg, rdb)

	// Initialize handlers
	handler := handlers.NewHandler(cfg, db, rdb)
	authHandler := handlers.NewAuthHandler(cfg, db, rdb)
	serverHandler := handlers.NewServerHandler(cfg, db, rdb)
	nodeHandler := handlers.NewNodeHandler(cfg, db, rdb)
	userHandler := handlers.NewUserHandler(cfg, db, rdb)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"version": cfg.App.Version,
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)

	// Protected routes
	protected := api.Group("", authMiddleware.Authenticate)

	// Auth (protected)
	protected.Post("/auth/logout", authHandler.Logout)
	protected.Get("/auth/me", authHandler.Me)
	protected.Post("/auth/2fa/enable", authHandler.Enable2FA)
	protected.Post("/auth/2fa/verify", authHandler.Verify2FA)
	protected.Post("/auth/2fa/disable", authHandler.Disable2FA)

	// Users
	users := protected.Group("/users")
	users.Get("/", authMiddleware.RequirePermission("users.view"), userHandler.List)
	users.Post("/", authMiddleware.RequirePermission("users.create"), userHandler.Create)
	users.Get("/:id", authMiddleware.RequirePermission("users.view"), userHandler.GetByID)
	users.Put("/:id", authMiddleware.RequirePermission("users.update"), userHandler.Update)
	users.Delete("/:id", authMiddleware.RequirePermission("users.delete"), userHandler.Delete)

	// Servers
	servers := protected.Group("/servers")
	servers.Get("/", serverHandler.List)
	servers.Post("/", authMiddleware.RequirePermission("servers.create"), serverHandler.Create)
	servers.Get("/:id", serverHandler.GetByID)
	servers.Put("/:id", serverHandler.Update)
	servers.Delete("/:id", authMiddleware.RequirePermission("servers.delete"), serverHandler.Delete)

	// Server power actions
	servers.Post("/:id/power/start", serverHandler.Start)
	servers.Post("/:id/power/stop", serverHandler.Stop)
	servers.Post("/:id/power/restart", serverHandler.Restart)
	servers.Post("/:id/power/kill", serverHandler.Kill)

	// Server console
	servers.Post("/:id/command", serverHandler.SendCommand)
	servers.Get("/:id/stats", serverHandler.GetStats)

	// Server backups
	servers.Get("/:id/backups", serverHandler.ListBackups)
	servers.Post("/:id/backups", serverHandler.CreateBackup)
	servers.Delete("/:id/backups/:backupId", serverHandler.DeleteBackup)
	servers.Post("/:id/backups/:backupId/restore", serverHandler.RestoreBackup)

	// Nodes (admin only)
	nodes := protected.Group("/nodes", authMiddleware.RequirePermission("nodes.view"))
	nodes.Get("/", nodeHandler.List)
	nodes.Post("/", authMiddleware.RequirePermission("nodes.create"), nodeHandler.Create)
	nodes.Get("/:id", nodeHandler.GetByID)
	nodes.Put("/:id", authMiddleware.RequirePermission("nodes.update"), nodeHandler.Update)
	nodes.Delete("/:id", authMiddleware.RequirePermission("nodes.delete"), nodeHandler.Delete)
	nodes.Get("/:id/configuration", nodeHandler.GetConfiguration)
	nodes.Post("/:id/allocations", nodeHandler.CreateAllocations)

	// Locations (admin only)
	locations := protected.Group("/locations", authMiddleware.RequirePermission("nodes.view"))
	locations.Get("/", handler.GetLocations)
	locations.Post("/", authMiddleware.RequirePermission("nodes.create"), handler.CreateLocation)
	locations.Get("/:id", handler.GetLocation)
	locations.Put("/:id", authMiddleware.RequirePermission("nodes.update"), handler.UpdateLocation)
	locations.Delete("/:id", authMiddleware.RequirePermission("nodes.delete"), handler.DeleteLocation)

	// Nodes (admin only) - update to use new handler
	nodes.Get("/", handler.GetNodes)
	nodes.Post("/", authMiddleware.RequirePermission("nodes.create"), handler.CreateNode)
	nodes.Get("/:id", handler.GetNode)
	nodes.Put("/:id", authMiddleware.RequirePermission("nodes.update"), handler.UpdateNode)
	nodes.Delete("/:id", authMiddleware.RequirePermission("nodes.delete"), handler.DeleteNode)
	nodes.Get("/:id/configuration", handler.GetNodeConfiguration)

	// Servers - update to use new handler
	servers.Get("/", handler.GetServers)
	servers.Post("/", authMiddleware.RequirePermission("servers.create"), handler.CreateServer)
	servers.Get("/:id", handler.GetServer)
	servers.Put("/:id", handler.UpdateServer)
	servers.Delete("/:id", authMiddleware.RequirePermission("servers.delete"), handler.DeleteServer)

	// Server power actions - update to use new handler
	servers.Post("/:id/start", handler.StartServer)
	servers.Post("/:id/stop", handler.StopServer)
	servers.Post("/:id/restart", handler.RestartServer)

	// WebSocket for real-time console
	app.Get("/ws/console/:serverId", websocket.New(func(c *websocket.Conn) {
		handleConsoleWebSocket(c, cfg, rdb)
	}))

	// WebSocket for real-time stats
	app.Get("/ws/stats/:serverId", websocket.New(func(c *websocket.Conn) {
		handleStatsWebSocket(c, cfg, rdb)
	}))

	return app
}

// errorHandler handles errors globally
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   message,
		"code":    code,
		"success": false,
	})
}

// joinStrings joins a slice of strings with comma
func joinStrings(s []string) string {
	if len(s) == 0 {
		return ""
	}
	result := s[0]
	for i := 1; i < len(s); i++ {
		result += "," + s[i]
	}
	return result
}

// handleConsoleWebSocket handles WebSocket connections for server console
func handleConsoleWebSocket(c *websocket.Conn, cfg *config.Config, rdb *redis.Client) {
	serverID := c.Params("serverId")
	
	// Subscribe to console channel
	ctx := context.Background()
	pubsub := rdb.Subscribe(ctx, "console:"+serverID)
	defer pubsub.Close()

	ch := pubsub.Channel()

	// Read messages from Redis and send to client
	go func() {
		for msg := range ch {
			if err := c.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				return
			}
		}
	}()

	// Read messages from client and publish to Redis
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}
		_ = rdb.Publish(ctx, "console:"+serverID+":input", string(msg))
	}
}

// handleStatsWebSocket handles WebSocket connections for server stats
func handleStatsWebSocket(c *websocket.Conn, cfg *config.Config, rdb *redis.Client) {
	serverID := c.Params("serverId")
	
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	ctx := context.Background()
	for range ticker.C {
		// Get stats from Redis cache
		stats, err := rdb.Get(ctx, "stats:"+serverID)
		if err != nil {
			continue
		}
		
		if err := c.WriteMessage(websocket.TextMessage, []byte(stats)); err != nil {
			break
		}
	}
}
