package repositories

import (
	"context"
	"time"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/google/uuid"
)

// AuditLogRepository defines the interface for audit log data access
type AuditLogRepository interface {
	Create(ctx context.Context, log *entities.AuditLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error)
	List(ctx context.Context, params ListParams) ([]*entities.AuditLog, int64, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, params ListParams) ([]*entities.AuditLog, int64, error)
	GetByResource(ctx context.Context, resource string, resourceID uuid.UUID) ([]*entities.AuditLog, error)
	GetByAction(ctx context.Context, action entities.AuditAction, params ListParams) ([]*entities.AuditLog, int64, error)
	GetByDateRange(ctx context.Context, start, end time.Time, params ListParams) ([]*entities.AuditLog, int64, error)
	DeleteOlderThan(ctx context.Context, before time.Time) error
}

// ActivityLogRepository defines the interface for activity log data access
type ActivityLogRepository interface {
	Create(ctx context.Context, log *entities.ActivityLog) error
	GetByUserID(ctx context.Context, userID uuid.UUID, params ListParams) ([]*entities.ActivityLog, int64, error)
	GetByServerID(ctx context.Context, serverID uuid.UUID, params ListParams) ([]*entities.ActivityLog, int64, error)
	GetRecent(ctx context.Context, limit int) ([]*entities.ActivityLog, error)
	DeleteOlderThan(ctx context.Context, before time.Time) error
}

// SystemEventRepository defines the interface for system event data access
type SystemEventRepository interface {
	Create(ctx context.Context, event *entities.SystemEvent) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.SystemEvent, error)
	List(ctx context.Context, params ListParams) ([]*entities.SystemEvent, int64, error)
	GetByNodeID(ctx context.Context, nodeID uuid.UUID, params ListParams) ([]*entities.SystemEvent, int64, error)
	GetByServerID(ctx context.Context, serverID uuid.UUID, params ListParams) ([]*entities.SystemEvent, int64, error)
	GetBySeverity(ctx context.Context, severity string, params ListParams) ([]*entities.SystemEvent, int64, error)
	GetRecent(ctx context.Context, limit int) ([]*entities.SystemEvent, error)
	DeleteOlderThan(ctx context.Context, before time.Time) error
}

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	Create(ctx context.Context, notification *entities.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Notification, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, params ListParams) ([]*entities.Notification, int64, error)
	GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
}
