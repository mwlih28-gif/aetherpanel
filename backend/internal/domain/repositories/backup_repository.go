package repositories

import (
	"context"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/google/uuid"
)

// BackupRepository defines the interface for backup data access
type BackupRepository interface {
	Create(ctx context.Context, backup *entities.Backup) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Backup, error)
	Update(ctx context.Context, backup *entities.Backup) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByServerID(ctx context.Context, serverID uuid.UUID) ([]*entities.Backup, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.BackupStatus) error
	CountByServerID(ctx context.Context, serverID uuid.UUID) (int64, error)
	GetOldestByServerID(ctx context.Context, serverID uuid.UUID, limit int) ([]*entities.Backup, error)
	Lock(ctx context.Context, id uuid.UUID) error
	Unlock(ctx context.Context, id uuid.UUID) error
}

// BackupScheduleRepository defines the interface for backup schedule data access
type BackupScheduleRepository interface {
	Create(ctx context.Context, schedule *entities.BackupSchedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.BackupSchedule, error)
	Update(ctx context.Context, schedule *entities.BackupSchedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByServerID(ctx context.Context, serverID uuid.UUID) ([]*entities.BackupSchedule, error)
	GetActive(ctx context.Context) ([]*entities.BackupSchedule, error)
	GetDue(ctx context.Context) ([]*entities.BackupSchedule, error)
	UpdateLastRun(ctx context.Context, id uuid.UUID) error
}

// SnapshotRepository defines the interface for snapshot data access
type SnapshotRepository interface {
	Create(ctx context.Context, snapshot *entities.Snapshot) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Snapshot, error)
	Update(ctx context.Context, snapshot *entities.Snapshot) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByServerID(ctx context.Context, serverID uuid.UUID) ([]*entities.Snapshot, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.BackupStatus) error
}

// ServerTransferRepository defines the interface for server transfer data access
type ServerTransferRepository interface {
	Create(ctx context.Context, transfer *entities.ServerTransfer) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ServerTransfer, error)
	Update(ctx context.Context, transfer *entities.ServerTransfer) error
	GetByServerID(ctx context.Context, serverID uuid.UUID) ([]*entities.ServerTransfer, error)
	GetPending(ctx context.Context) ([]*entities.ServerTransfer, error)
	UpdateProgress(ctx context.Context, id uuid.UUID, progress int) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}
