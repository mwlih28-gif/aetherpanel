package repositories

import (
	"context"
	"time"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/google/uuid"
)

// PlayerRepository defines the interface for player data access
type PlayerRepository interface {
	Create(ctx context.Context, player *entities.Player) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Player, error)
	GetByUUID(ctx context.Context, playerUUID string) (*entities.Player, error)
	GetByUsername(ctx context.Context, serverID uuid.UUID, username string) (*entities.Player, error)
	Update(ctx context.Context, player *entities.Player) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByServerID(ctx context.Context, serverID uuid.UUID, params ListParams) ([]*entities.Player, int64, error)
	GetOnlineByServerID(ctx context.Context, serverID uuid.UUID) ([]*entities.Player, error)
	SetOnline(ctx context.Context, id uuid.UUID, ip string) error
	SetOffline(ctx context.Context, id uuid.UUID) error
	Ban(ctx context.Context, id uuid.UUID) error
	Unban(ctx context.Context, id uuid.UUID) error
	UpdatePlayTime(ctx context.Context, id uuid.UUID, seconds int64) error
	CountByServerID(ctx context.Context, serverID uuid.UUID) (int64, error)
	CountOnlineByServerID(ctx context.Context, serverID uuid.UUID) (int64, error)
}

// PlayerInventoryRepository defines the interface for player inventory data access
type PlayerInventoryRepository interface {
	Create(ctx context.Context, inventory *entities.PlayerInventory) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.PlayerInventory, error)
	GetLatestByPlayerID(ctx context.Context, playerID uuid.UUID) (*entities.PlayerInventory, error)
	GetByPlayerID(ctx context.Context, playerID uuid.UUID, limit int) ([]*entities.PlayerInventory, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// PlayerStatsRepository defines the interface for player stats data access
type PlayerStatsRepository interface {
	Upsert(ctx context.Context, playerID uuid.UUID, statType string, value int64) error
	GetByPlayerID(ctx context.Context, playerID uuid.UUID) ([]*entities.PlayerStats, error)
	GetStat(ctx context.Context, playerID uuid.UUID, statType string) (*entities.PlayerStats, error)
	IncrementStat(ctx context.Context, playerID uuid.UUID, statType string, delta int64) error
	GetLeaderboard(ctx context.Context, serverID uuid.UUID, statType string, limit int) ([]*entities.PlayerStats, error)
}

// ChatLogRepository defines the interface for chat log data access
type ChatLogRepository interface {
	Create(ctx context.Context, log *entities.ChatLog) error
	GetByServerID(ctx context.Context, serverID uuid.UUID, params ListParams) ([]*entities.ChatLog, int64, error)
	GetByPlayerID(ctx context.Context, playerID uuid.UUID, params ListParams) ([]*entities.ChatLog, int64, error)
	Search(ctx context.Context, serverID uuid.UUID, query string, start, end time.Time) ([]*entities.ChatLog, error)
	DeleteOlderThan(ctx context.Context, before time.Time) error
}

// CommandLogRepository defines the interface for command log data access
type CommandLogRepository interface {
	Create(ctx context.Context, log *entities.CommandLog) error
	GetByServerID(ctx context.Context, serverID uuid.UUID, params ListParams) ([]*entities.CommandLog, int64, error)
	GetByPlayerID(ctx context.Context, playerID uuid.UUID, params ListParams) ([]*entities.CommandLog, int64, error)
	Search(ctx context.Context, serverID uuid.UUID, query string) ([]*entities.CommandLog, error)
	DeleteOlderThan(ctx context.Context, before time.Time) error
}

// DeathLogRepository defines the interface for death log data access
type DeathLogRepository interface {
	Create(ctx context.Context, log *entities.DeathLog) error
	GetByServerID(ctx context.Context, serverID uuid.UUID, params ListParams) ([]*entities.DeathLog, int64, error)
	GetByPlayerID(ctx context.Context, playerID uuid.UUID, params ListParams) ([]*entities.DeathLog, int64, error)
	GetByKillerID(ctx context.Context, killerID uuid.UUID, params ListParams) ([]*entities.DeathLog, int64, error)
	CountByPlayerID(ctx context.Context, playerID uuid.UUID) (int64, error)
	DeleteOlderThan(ctx context.Context, before time.Time) error
}
