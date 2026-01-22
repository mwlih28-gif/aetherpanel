package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aetherpanel/aether-panel/internal/infrastructure/config"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/database"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/logger"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/redis"
	"github.com/aetherpanel/aether-panel/internal/interfaces/http"
	"go.uber.org/zap"
)

// @title Aether Panel API
// @version 1.0.0
// @description Next-Generation Game Server Management Platform API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Initialize logger
	log := logger.NewLogger()
	defer log.Sync()

	log.Info("🔥 Starting Aether Panel API Server...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize database connection
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	log.Info("✅ Database connection established")

	// Initialize Redis connection
	rdb, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	log.Info("✅ Redis connection established")

	// Initialize HTTP server
	server := http.NewServer(cfg, db, rdb, log)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Info("🚀 Server starting", zap.String("address", addr))
		if err := server.Listen(addr); err != nil {
			log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.ShutdownWithContext(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	log.Info("👋 Server exited properly")
}
