package entities

import (
	"time"

	"github.com/google/uuid"
)

// AuditAction represents the type of audit action
type AuditAction string

const (
	AuditActionCreate  AuditAction = "create"
	AuditActionUpdate  AuditAction = "update"
	AuditActionDelete  AuditAction = "delete"
	AuditActionLogin   AuditAction = "login"
	AuditActionLogout  AuditAction = "logout"
	AuditActionStart   AuditAction = "start"
	AuditActionStop    AuditAction = "stop"
	AuditActionRestart AuditAction = "restart"
	AuditActionBackup  AuditAction = "backup"
	AuditActionRestore AuditAction = "restore"
	AuditActionInstall AuditAction = "install"
	AuditActionCommand AuditAction = "command"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      *uuid.UUID             `json:"user_id" gorm:"type:uuid;index"`
	User        *User                  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Action      AuditAction            `json:"action" gorm:"type:varchar(20);not null;index"`
	Resource    string                 `json:"resource" gorm:"size:50;index"` // user, server, node, etc.
	ResourceID  *uuid.UUID             `json:"resource_id" gorm:"type:uuid;index"`
	Description string                 `json:"description" gorm:"size:500"`
	OldValues   map[string]interface{} `json:"old_values" gorm:"type:jsonb"`
	NewValues   map[string]interface{} `json:"new_values" gorm:"type:jsonb"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	IPAddress   string                 `json:"ip_address" gorm:"size:45"`
	UserAgent   string                 `json:"user_agent" gorm:"size:500"`
	IsSystem    bool                   `json:"is_system" gorm:"default:false"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime;index"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// ActivityLog represents user activity for analytics
type ActivityLog struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ServerID   *uuid.UUID `json:"server_id" gorm:"type:uuid;index"`
	Action     string    `json:"action" gorm:"size:50;not null;index"`
	Details    string    `json:"details" gorm:"size:500"`
	IPAddress  string    `json:"ip_address" gorm:"size:45"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}

// SystemEvent represents system-level events
type SystemEvent struct {
	ID        uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NodeID    *uuid.UUID             `json:"node_id" gorm:"type:uuid;index"`
	Node      *Node                  `json:"node,omitempty" gorm:"foreignKey:NodeID"`
	ServerID  *uuid.UUID             `json:"server_id" gorm:"type:uuid;index"`
	Server    *Server                `json:"server,omitempty" gorm:"foreignKey:ServerID"`
	EventType string                 `json:"event_type" gorm:"size:50;not null;index"`
	Severity  string                 `json:"severity" gorm:"size:20;index"` // info, warning, error, critical
	Message   string                 `json:"message" gorm:"type:text"`
	Data      map[string]interface{} `json:"data" gorm:"type:jsonb"`
	CreatedAt time.Time              `json:"created_at" gorm:"autoCreateTime;index"`
}

func (SystemEvent) TableName() string {
	return "system_events"
}

// Notification represents a user notification
type Notification struct {
	ID        uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID              `json:"user_id" gorm:"type:uuid;not null;index"`
	User      *User                  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Type      string                 `json:"type" gorm:"size:50;not null"`
	Title     string                 `json:"title" gorm:"size:200;not null"`
	Message   string                 `json:"message" gorm:"type:text"`
	Data      map[string]interface{} `json:"data" gorm:"type:jsonb"`
	IsRead    bool                   `json:"is_read" gorm:"default:false;index"`
	ReadAt    *time.Time             `json:"read_at"`
	CreatedAt time.Time              `json:"created_at" gorm:"autoCreateTime"`
}

func (Notification) TableName() string {
	return "notifications"
}
