package repositories

import (
	"context"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/google/uuid"
)

// PluginRepository defines the interface for plugin data access
type PluginRepository interface {
	Create(ctx context.Context, plugin *entities.Plugin) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Plugin, error)
	GetByExternalID(ctx context.Context, source entities.PluginSource, externalID string) (*entities.Plugin, error)
	Update(ctx context.Context, plugin *entities.Plugin) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params ListParams) ([]*entities.Plugin, int64, error)
	Search(ctx context.Context, query string, filters PluginFilters) ([]*entities.Plugin, int64, error)
	GetFeatured(ctx context.Context, limit int) ([]*entities.Plugin, error)
	GetPopular(ctx context.Context, limit int) ([]*entities.Plugin, error)
	IncrementDownloads(ctx context.Context, id uuid.UUID) error
}

// PluginFilters represents filters for plugin search
type PluginFilters struct {
	Source       entities.PluginSource
	Type         entities.PluginType
	Categories   []string
	GameVersions []string
	Loaders      []string
	IsVerified   *bool
}

// PluginVersionRepository defines the interface for plugin version data access
type PluginVersionRepository interface {
	Create(ctx context.Context, version *entities.PluginVersion) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.PluginVersion, error)
	Update(ctx context.Context, version *entities.PluginVersion) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByPluginID(ctx context.Context, pluginID uuid.UUID) ([]*entities.PluginVersion, error)
	GetLatestByPluginID(ctx context.Context, pluginID uuid.UUID) (*entities.PluginVersion, error)
	GetCompatible(ctx context.Context, pluginID uuid.UUID, gameVersion, loader string) ([]*entities.PluginVersion, error)
}

// InstalledPluginRepository defines the interface for installed plugin data access
type InstalledPluginRepository interface {
	Create(ctx context.Context, installed *entities.InstalledPlugin) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.InstalledPlugin, error)
	Update(ctx context.Context, installed *entities.InstalledPlugin) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByServerID(ctx context.Context, serverID uuid.UUID) ([]*entities.InstalledPlugin, error)
	GetByServerAndPlugin(ctx context.Context, serverID, pluginID uuid.UUID) (*entities.InstalledPlugin, error)
	Enable(ctx context.Context, id uuid.UUID) error
	Disable(ctx context.Context, id uuid.UUID) error
	CountByServerID(ctx context.Context, serverID uuid.UUID) (int64, error)
}

// WorldRepository defines the interface for world data access
type WorldRepository interface {
	Create(ctx context.Context, world *entities.World) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.World, error)
	Update(ctx context.Context, world *entities.World) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByServerID(ctx context.Context, serverID uuid.UUID) ([]*entities.World, error)
	GetByName(ctx context.Context, serverID uuid.UUID, name string) (*entities.World, error)
}

// WorldBackupRepository defines the interface for world backup data access
type WorldBackupRepository interface {
	Create(ctx context.Context, backup *entities.WorldBackup) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.WorldBackup, error)
	Update(ctx context.Context, backup *entities.WorldBackup) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByWorldID(ctx context.Context, worldID uuid.UUID) ([]*entities.WorldBackup, error)
}
