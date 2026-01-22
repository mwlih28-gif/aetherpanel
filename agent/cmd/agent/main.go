package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aetherpanel/aether-panel/agent/internal/config"
	"github.com/aetherpanel/aether-panel/agent/internal/docker"
	"github.com/aetherpanel/aether-panel/agent/internal/server"
	"github.com/aetherpanel/aether-panel/agent/internal/api"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("🔥 Starting Aether Node Agent...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		logger.Fatal("Failed to connect to Docker", zap.Error(err))
	}
	defer dockerClient.Close()
	logger.Info("✅ Docker connection established")

	// Initialize server manager
	serverManager := server.NewManager(dockerClient, cfg, logger)

	// Initialize API server for panel communication
	apiServer := api.NewServer(cfg, serverManager, logger)

	// Start background services
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start health check routine
	go serverManager.StartHealthCheck(ctx)

	// Start metrics collection
	go serverManager.StartMetricsCollection(ctx)

	// Start console streaming
	go serverManager.StartConsoleStreaming(ctx)

	// Register with panel
	if err := apiServer.RegisterWithPanel(ctx); err != nil {
		logger.Warn("Failed to register with panel, will retry", zap.Error(err))
	}

	// Start API server
	go func() {
		addr := fmt.Sprintf(":%d", cfg.API.Port)
		logger.Info("🚀 Agent API starting", zap.String("address", addr))
		if err := apiServer.Start(addr); err != nil {
			logger.Fatal("API server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Shutting down agent...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("API server shutdown error", zap.Error(err))
	}

	logger.Info("👋 Agent stopped")
}
