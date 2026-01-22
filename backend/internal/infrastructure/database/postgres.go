package database

import (
	"fmt"
	"time"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/aetherpanel/aether-panel/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresConnection creates a new PostgreSQL database connection
func NewPostgresConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// Configure GORM logger
	var logLevel logger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Warn
	}

	gormConfig := &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: false,
		PrepareStmt:                              true,
	}

	// Open connection
	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// AutoMigrate runs database migrations for all entities
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// User & Auth
		&entities.User{},
		&entities.Role{},
		&entities.Permission{},
		&entities.Session{},
		&entities.APIKey{},

		// Server & Node
		&entities.Server{},
		&entities.Node{},
		&entities.Location{},
		&entities.Allocation{},
		&entities.Game{},
		&entities.Egg{},
		&entities.EggVariable{},
		&entities.ServerVariable{},

		// Backup
		&entities.Backup{},
		&entities.BackupSchedule{},
		&entities.Snapshot{},
		&entities.ServerTransfer{},

		// Billing
		&entities.Transaction{},
		&entities.Package{},
		&entities.Subscription{},
		&entities.Invoice{},
		&entities.InvoiceItem{},
		&entities.Coupon{},

		// Plugin & World
		&entities.Plugin{},
		&entities.PluginVersion{},
		&entities.InstalledPlugin{},
		&entities.World{},
		&entities.WorldBackup{},

		// Player
		&entities.Player{},
		&entities.PlayerInventory{},
		&entities.PlayerStats{},
		&entities.ChatLog{},
		&entities.CommandLog{},
		&entities.DeathLog{},

		// Audit & Logs
		&entities.AuditLog{},
		&entities.ActivityLog{},
		&entities.SystemEvent{},
		&entities.Notification{},
	)
}

// SeedDefaultData seeds default data into the database
func SeedDefaultData(db *gorm.DB) error {
	// Seed default permissions
	permissions := getDefaultPermissions()
	for _, p := range permissions {
		if err := db.FirstOrCreate(&p, entities.Permission{Name: p.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed permission %s: %w", p.Name, err)
		}
	}

	// Seed default roles
	roles := getDefaultRoles()
	for _, r := range roles {
		if err := db.FirstOrCreate(&r, entities.Role{Name: r.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed role %s: %w", r.Name, err)
		}
	}

	// Assign permissions to admin role
	var adminRole entities.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return fmt.Errorf("failed to find admin role: %w", err)
	}

	var allPermissions []entities.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return fmt.Errorf("failed to find permissions: %w", err)
	}

	if err := db.Model(&adminRole).Association("Permissions").Replace(allPermissions); err != nil {
		return fmt.Errorf("failed to assign permissions to admin: %w", err)
	}

	// Seed default games
	games := getDefaultGames()
	for _, g := range games {
		if err := db.FirstOrCreate(&g, entities.Game{Name: g.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed game %s: %w", g.Name, err)
		}
	}

	return nil
}

// getDefaultPermissions returns default system permissions
func getDefaultPermissions() []entities.Permission {
	now := time.Now()
	return []entities.Permission{
		// User permissions
		{Name: "users.view", DisplayName: "View Users", Category: "users", CreatedAt: now},
		{Name: "users.create", DisplayName: "Create Users", Category: "users", CreatedAt: now},
		{Name: "users.update", DisplayName: "Update Users", Category: "users", CreatedAt: now},
		{Name: "users.delete", DisplayName: "Delete Users", Category: "users", CreatedAt: now},

		// Server permissions
		{Name: "servers.view", DisplayName: "View Servers", Category: "servers", CreatedAt: now},
		{Name: "servers.create", DisplayName: "Create Servers", Category: "servers", CreatedAt: now},
		{Name: "servers.update", DisplayName: "Update Servers", Category: "servers", CreatedAt: now},
		{Name: "servers.delete", DisplayName: "Delete Servers", Category: "servers", CreatedAt: now},
		{Name: "servers.console", DisplayName: "Access Console", Category: "servers", CreatedAt: now},
		{Name: "servers.files", DisplayName: "Manage Files", Category: "servers", CreatedAt: now},
		{Name: "servers.power", DisplayName: "Power Actions", Category: "servers", CreatedAt: now},
		{Name: "servers.backup", DisplayName: "Manage Backups", Category: "servers", CreatedAt: now},

		// Node permissions
		{Name: "nodes.view", DisplayName: "View Nodes", Category: "nodes", CreatedAt: now},
		{Name: "nodes.create", DisplayName: "Create Nodes", Category: "nodes", CreatedAt: now},
		{Name: "nodes.update", DisplayName: "Update Nodes", Category: "nodes", CreatedAt: now},
		{Name: "nodes.delete", DisplayName: "Delete Nodes", Category: "nodes", CreatedAt: now},

		// Billing permissions
		{Name: "billing.view", DisplayName: "View Billing", Category: "billing", CreatedAt: now},
		{Name: "billing.manage", DisplayName: "Manage Billing", Category: "billing", CreatedAt: now},

		// Admin permissions
		{Name: "admin.settings", DisplayName: "Manage Settings", Category: "admin", CreatedAt: now},
		{Name: "admin.audit", DisplayName: "View Audit Logs", Category: "admin", CreatedAt: now},
	}
}

// getDefaultRoles returns default system roles
func getDefaultRoles() []entities.Role {
	now := time.Now()
	return []entities.Role{
		{
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "Full system access",
			Color:       "#ef4444",
			IsSystem:    true,
			Priority:    100,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "reseller",
			DisplayName: "Reseller",
			Description: "Reseller account with limited admin access",
			Color:       "#8b5cf6",
			IsSystem:    true,
			Priority:    50,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "user",
			DisplayName: "User",
			Description: "Standard user account",
			Color:       "#3b82f6",
			IsSystem:    true,
			IsDefault:   true,
			Priority:    10,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
}

// getDefaultGames returns default supported games
func getDefaultGames() []entities.Game {
	now := time.Now()
	return []entities.Game{
		{Name: "Minecraft Java", Description: "Minecraft Java Edition", Category: "minecraft", SortOrder: 1, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Minecraft Bedrock", Description: "Minecraft Bedrock Edition", Category: "minecraft", SortOrder: 2, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Rust", Description: "Rust Survival Game", Category: "survival", SortOrder: 3, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "ARK: Survival Evolved", Description: "ARK Survival Game", Category: "survival", SortOrder: 4, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Valheim", Description: "Valheim Viking Survival", Category: "survival", SortOrder: 5, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "CS2", Description: "Counter-Strike 2", Category: "fps", SortOrder: 6, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Terraria", Description: "Terraria Sandbox", Category: "sandbox", SortOrder: 7, IsActive: true, CreatedAt: now, UpdatedAt: now},
	}
}
