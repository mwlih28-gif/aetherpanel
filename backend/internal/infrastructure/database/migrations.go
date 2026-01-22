package database

import (
	"github.com/aetherpanel/aether-panel/internal/interfaces/http/handlers"
	"gorm.io/gorm"
)

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&handlers.Location{},
		&handlers.Node{},
		&handlers.GameServer{},
	)
}
