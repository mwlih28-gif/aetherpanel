package repositories

import (
	"context"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/google/uuid"
)

// ServerRepository defines the interface for server data access
type ServerRepository interface {
	Create(ctx context.Context, server *entities.Server) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Server, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.Server, error)
	Update(ctx context.Context, server *entities.Server) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params ListParams) ([]*entities.Server, int64, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entities.Server, error)
	GetByNodeID(ctx context.Context, nodeID uuid.UUID) ([]*entities.Server, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.ServerStatus) error
	UpdateContainerID(ctx context.Context, id uuid.UUID, containerID string) error
	Suspend(ctx context.Context, id uuid.UUID, reason string) error
	Unsuspend(ctx context.Context, id uuid.UUID) error
	CountByNodeID(ctx context.Context, nodeID uuid.UUID) (int64, error)
	CountByOwnerID(ctx context.Context, ownerID uuid.UUID) (int64, error)
}

// NodeRepository defines the interface for node data access
type NodeRepository interface {
	Create(ctx context.Context, node *entities.Node) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Node, error)
	GetByName(ctx context.Context, name string) (*entities.Node, error)
	Update(ctx context.Context, node *entities.Node) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params ListParams) ([]*entities.Node, int64, error)
	GetByLocationID(ctx context.Context, locationID uuid.UUID) ([]*entities.Node, error)
	GetAvailable(ctx context.Context, memoryRequired, diskRequired int64) ([]*entities.Node, error)
	UpdateOnlineStatus(ctx context.Context, id uuid.UUID, isOnline bool) error
	UpdateResources(ctx context.Context, id uuid.UUID, memoryAlloc, diskAlloc int64, cpuAlloc int) error
	SetMaintenanceMode(ctx context.Context, id uuid.UUID, maintenance bool) error
}

// LocationRepository defines the interface for location data access
type LocationRepository interface {
	Create(ctx context.Context, location *entities.Location) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Location, error)
	GetByShortCode(ctx context.Context, shortCode string) (*entities.Location, error)
	Update(ctx context.Context, location *entities.Location) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*entities.Location, error)
}

// AllocationRepository defines the interface for allocation data access
type AllocationRepository interface {
	Create(ctx context.Context, allocation *entities.Allocation) error
	CreateBatch(ctx context.Context, allocations []*entities.Allocation) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Allocation, error)
	Update(ctx context.Context, allocation *entities.Allocation) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByNodeID(ctx context.Context, nodeID uuid.UUID) ([]*entities.Allocation, error)
	GetAvailableByNodeID(ctx context.Context, nodeID uuid.UUID) ([]*entities.Allocation, error)
	GetByServerID(ctx context.Context, serverID uuid.UUID) ([]*entities.Allocation, error)
	AssignToServer(ctx context.Context, id uuid.UUID, serverID uuid.UUID, isPrimary bool) error
	Unassign(ctx context.Context, id uuid.UUID) error
	IsPortAvailable(ctx context.Context, nodeID uuid.UUID, ip string, port int) (bool, error)
}

// GameRepository defines the interface for game data access
type GameRepository interface {
	Create(ctx context.Context, game *entities.Game) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Game, error)
	GetByName(ctx context.Context, name string) (*entities.Game, error)
	Update(ctx context.Context, game *entities.Game) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*entities.Game, error)
	GetActive(ctx context.Context) ([]*entities.Game, error)
}

// EggRepository defines the interface for egg data access
type EggRepository interface {
	Create(ctx context.Context, egg *entities.Egg) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Egg, error)
	Update(ctx context.Context, egg *entities.Egg) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params ListParams) ([]*entities.Egg, int64, error)
	GetByGameID(ctx context.Context, gameID uuid.UUID) ([]*entities.Egg, error)
	GetActive(ctx context.Context) ([]*entities.Egg, error)
}

// EggVariableRepository defines the interface for egg variable data access
type EggVariableRepository interface {
	Create(ctx context.Context, variable *entities.EggVariable) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.EggVariable, error)
	Update(ctx context.Context, variable *entities.EggVariable) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEggID(ctx context.Context, eggID uuid.UUID) ([]*entities.EggVariable, error)
}

// ServerVariableRepository defines the interface for server variable data access
type ServerVariableRepository interface {
	Create(ctx context.Context, variable *entities.ServerVariable) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ServerVariable, error)
	Update(ctx context.Context, variable *entities.ServerVariable) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByServerID(ctx context.Context, serverID uuid.UUID) ([]*entities.ServerVariable, error)
	Upsert(ctx context.Context, serverID, eggVariableID uuid.UUID, value string) error
}
