package handlers

import (
	"github.com/aetherpanel/aether-panel/internal/infrastructure/config"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/redis"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Handler contains all dependencies for HTTP handlers
type Handler struct {
	cfg       *config.Config
	db        *gorm.DB
	redis     *redis.Client
	validator *validator.Validate
}

// NewHandler creates a new handler instance
func NewHandler(cfg *config.Config, db *gorm.DB, redis *redis.Client) *Handler {
	return &Handler{
		cfg:       cfg,
		db:        db,
		redis:     redis,
		validator: validator.New(),
	}
}
