package entities

import (
	"time"

	"github.com/google/uuid"
)

// PluginSource represents the source of a plugin
type PluginSource string

const (
	PluginSourceSpigot     PluginSource = "spigot"
	PluginSourceModrinth   PluginSource = "modrinth"
	PluginSourceCurseForge PluginSource = "curseforge"
	PluginSourceCustom     PluginSource = "custom"
	PluginSourceLocal      PluginSource = "local"
)

// PluginType represents the type of plugin/mod
type PluginType string

const (
	PluginTypePlugin  PluginType = "plugin"
	PluginTypeMod     PluginType = "mod"
	PluginTypeModpack PluginType = "modpack"
	PluginTypeWorld   PluginType = "world"
)

// Plugin represents a plugin/mod in the marketplace
type Plugin struct {
	ID              uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ExternalID      string       `json:"external_id" gorm:"size:100;index"` // ID from source
	Source          PluginSource `json:"source" gorm:"type:varchar(20);not null"`
	Type            PluginType   `json:"type" gorm:"type:varchar(20);not null"`
	Name            string       `json:"name" gorm:"not null;size:200"`
	Slug            string       `json:"slug" gorm:"index;size:200"`
	Description     string       `json:"description" gorm:"type:text"`
	Author          string       `json:"author" gorm:"size:100"`
	IconURL         string       `json:"icon_url" gorm:"size:500"`
	WebsiteURL      string       `json:"website_url" gorm:"size:500"`
	SourceURL       string       `json:"source_url" gorm:"size:500"`
	Downloads       int64        `json:"downloads" gorm:"default:0"`
	Rating          float32      `json:"rating" gorm:"type:decimal(3,2);default:0"`
	RatingCount     int          `json:"rating_count" gorm:"default:0"`
	Categories      []string     `json:"categories" gorm:"type:jsonb"`
	Tags            []string     `json:"tags" gorm:"type:jsonb"`
	GameVersions    []string     `json:"game_versions" gorm:"type:jsonb"`
	Loaders         []string     `json:"loaders" gorm:"type:jsonb"` // fabric, forge, spigot, paper, etc.
	IsFeatured      bool         `json:"is_featured" gorm:"default:false"`
	IsVerified      bool         `json:"is_verified" gorm:"default:false"`
	LastSyncedAt    *time.Time   `json:"last_synced_at"`
	CreatedAt       time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Plugin
func (Plugin) TableName() string {
	return "plugins"
}

// PluginVersion represents a specific version of a plugin
type PluginVersion struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PluginID        uuid.UUID  `json:"plugin_id" gorm:"type:uuid;not null;index"`
	Plugin          *Plugin    `json:"plugin,omitempty" gorm:"foreignKey:PluginID"`
	ExternalID      string     `json:"external_id" gorm:"size:100"`
	VersionNumber   string     `json:"version_number" gorm:"not null;size:50"`
	VersionName     string     `json:"version_name" gorm:"size:100"`
	Changelog       string     `json:"changelog" gorm:"type:text"`
	DownloadURL     string     `json:"download_url" gorm:"size:500"`
	FileName        string     `json:"file_name" gorm:"size:255"`
	FileSize        int64      `json:"file_size" gorm:"default:0"`
	FileHash        string     `json:"file_hash" gorm:"size:64"` // SHA-256
	GameVersions    []string   `json:"game_versions" gorm:"type:jsonb"`
	Loaders         []string   `json:"loaders" gorm:"type:jsonb"`
	Dependencies    []PluginDependency `json:"dependencies" gorm:"type:jsonb"`
	IsStable        bool       `json:"is_stable" gorm:"default:true"`
	Downloads       int64      `json:"downloads" gorm:"default:0"`
	ReleasedAt      time.Time  `json:"released_at"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for PluginVersion
func (PluginVersion) TableName() string {
	return "plugin_versions"
}

// PluginDependency represents a plugin dependency
type PluginDependency struct {
	PluginID     string `json:"plugin_id"`
	PluginName   string `json:"plugin_name"`
	VersionRange string `json:"version_range"`
	Required     bool   `json:"required"`
}

// InstalledPlugin represents a plugin installed on a server
type InstalledPlugin struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID        uuid.UUID  `json:"server_id" gorm:"type:uuid;not null;index"`
	Server          *Server    `json:"server,omitempty" gorm:"foreignKey:ServerID"`
	PluginID        *uuid.UUID `json:"plugin_id" gorm:"type:uuid;index"`
	Plugin          *Plugin    `json:"plugin,omitempty" gorm:"foreignKey:PluginID"`
	PluginVersionID *uuid.UUID `json:"plugin_version_id" gorm:"type:uuid"`
	PluginVersion   *PluginVersion `json:"plugin_version,omitempty" gorm:"foreignKey:PluginVersionID"`
	
	// For custom/local plugins
	Name            string     `json:"name" gorm:"size:200"`
	FileName        string     `json:"file_name" gorm:"size:255"`
	FilePath        string     `json:"file_path" gorm:"size:500"`
	FileSize        int64      `json:"file_size" gorm:"default:0"`
	FileHash        string     `json:"file_hash" gorm:"size:64"`
	
	IsEnabled       bool       `json:"is_enabled" gorm:"default:true"`
	AutoUpdate      bool       `json:"auto_update" gorm:"default:false"`
	InstalledBy     uuid.UUID  `json:"installed_by" gorm:"type:uuid"`
	InstalledAt     time.Time  `json:"installed_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for InstalledPlugin
func (InstalledPlugin) TableName() string {
	return "installed_plugins"
}

// World represents a Minecraft world
type World struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID    uuid.UUID  `json:"server_id" gorm:"type:uuid;not null;index"`
	Server      *Server    `json:"server,omitempty" gorm:"foreignKey:ServerID"`
	Name        string     `json:"name" gorm:"not null;size:100"`
	FolderName  string     `json:"folder_name" gorm:"not null;size:100"`
	Dimension   string     `json:"dimension" gorm:"size:50"` // overworld, nether, end
	Seed        string     `json:"seed" gorm:"size:50"`
	GameMode    string     `json:"game_mode" gorm:"size:20"`
	Difficulty  string     `json:"difficulty" gorm:"size:20"`
	Size        int64      `json:"size" gorm:"default:0"` // Bytes
	LastPlayed  *time.Time `json:"last_played"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for World
func (World) TableName() string {
	return "worlds"
}

// WorldBackup represents a world backup
type WorldBackup struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WorldID     uuid.UUID    `json:"world_id" gorm:"type:uuid;not null;index"`
	World       *World       `json:"world,omitempty" gorm:"foreignKey:WorldID"`
	Name        string       `json:"name" gorm:"not null;size:100"`
	Status      BackupStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Size        int64        `json:"size" gorm:"default:0"`
	StoragePath string       `json:"-" gorm:"size:500"`
	Checksum    string       `json:"checksum" gorm:"size:64"`
	CompletedAt *time.Time   `json:"completed_at"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for WorldBackup
func (WorldBackup) TableName() string {
	return "world_backups"
}
