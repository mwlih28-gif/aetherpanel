package entities

import (
	"time"

	"github.com/google/uuid"
)

// BackupStatus represents the status of a backup
type BackupStatus string

const (
	BackupStatusPending    BackupStatus = "pending"
	BackupStatusInProgress BackupStatus = "in_progress"
	BackupStatusCompleted  BackupStatus = "completed"
	BackupStatusFailed     BackupStatus = "failed"
	BackupStatusDeleted    BackupStatus = "deleted"
)

// Backup represents a server backup
type Backup struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID    uuid.UUID    `json:"server_id" gorm:"type:uuid;not null;index"`
	Server      *Server      `json:"server,omitempty" gorm:"foreignKey:ServerID"`
	Name        string       `json:"name" gorm:"not null;size:100"`
	Status      BackupStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Checksum    string       `json:"checksum" gorm:"size:64"` // SHA-256
	Size        int64        `json:"size" gorm:"default:0"`   // Bytes
	StoragePath string       `json:"-" gorm:"size:500"`
	IsLocked    bool         `json:"is_locked" gorm:"default:false"`
	IsScheduled bool         `json:"is_scheduled" gorm:"default:false"`
	ScheduleID  *uuid.UUID   `json:"schedule_id" gorm:"type:uuid"`
	ErrorMsg    string       `json:"error_msg" gorm:"size:500"`
	CompletedAt *time.Time   `json:"completed_at"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time   `json:"deleted_at" gorm:"index"`
}

// TableName returns the table name for Backup
func (Backup) TableName() string {
	return "backups"
}

// BackupSchedule represents an automated backup schedule
type BackupSchedule struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID        uuid.UUID  `json:"server_id" gorm:"type:uuid;not null;index"`
	Server          *Server    `json:"server,omitempty" gorm:"foreignKey:ServerID"`
	Name            string     `json:"name" gorm:"not null;size:100"`
	CronExpression  string     `json:"cron_expression" gorm:"not null;size:50"`
	RetentionCount  int        `json:"retention_count" gorm:"default:3"`
	IsActive        bool       `json:"is_active" gorm:"default:true"`
	LastRunAt       *time.Time `json:"last_run_at"`
	NextRunAt       *time.Time `json:"next_run_at"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for BackupSchedule
func (BackupSchedule) TableName() string {
	return "backup_schedules"
}

// Snapshot represents a server snapshot (full state)
type Snapshot struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID    uuid.UUID    `json:"server_id" gorm:"type:uuid;not null;index"`
	Server      *Server      `json:"server,omitempty" gorm:"foreignKey:ServerID"`
	Name        string       `json:"name" gorm:"not null;size:100"`
	Description string       `json:"description" gorm:"size:500"`
	Status      BackupStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Size        int64        `json:"size" gorm:"default:0"`
	StoragePath string       `json:"-" gorm:"size:500"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb;default:'{}'"`
	CompletedAt *time.Time   `json:"completed_at"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time   `json:"deleted_at" gorm:"index"`
}

// TableName returns the table name for Snapshot
func (Snapshot) TableName() string {
	return "snapshots"
}

// ServerTransfer represents a server transfer between nodes
type ServerTransfer struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID     uuid.UUID  `json:"server_id" gorm:"type:uuid;not null;index"`
	Server       *Server    `json:"server,omitempty" gorm:"foreignKey:ServerID"`
	OldNodeID    uuid.UUID  `json:"old_node_id" gorm:"type:uuid;not null"`
	OldNode      *Node      `json:"old_node,omitempty" gorm:"foreignKey:OldNodeID"`
	NewNodeID    uuid.UUID  `json:"new_node_id" gorm:"type:uuid;not null"`
	NewNode      *Node      `json:"new_node,omitempty" gorm:"foreignKey:NewNodeID"`
	OldAllocID   uuid.UUID  `json:"old_allocation_id" gorm:"type:uuid;not null"`
	NewAllocID   uuid.UUID  `json:"new_allocation_id" gorm:"type:uuid;not null"`
	Status       string     `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Progress     int        `json:"progress" gorm:"default:0"` // 0-100
	ErrorMsg     string     `json:"error_msg" gorm:"size:500"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ServerTransfer
func (ServerTransfer) TableName() string {
	return "server_transfers"
}
