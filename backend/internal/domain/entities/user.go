package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusPending   UserStatus = "pending"
)

// User represents a user in the system
type User struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email             string     `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Username          string     `json:"username" gorm:"uniqueIndex;not null;size:50"`
	PasswordHash      string     `json:"-" gorm:"not null"`
	FirstName         string     `json:"first_name" gorm:"size:100"`
	LastName          string     `json:"last_name" gorm:"size:100"`
	Avatar            string     `json:"avatar" gorm:"size:500"`
	Status            UserStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	EmailVerified     bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at"`
	TwoFactorEnabled  bool       `json:"two_factor_enabled" gorm:"default:false"`
	TwoFactorSecret   string     `json:"-" gorm:"size:100"`
	Language          string     `json:"language" gorm:"size:10;default:'en'"`
	Timezone          string     `json:"timezone" gorm:"size:50;default:'UTC'"`
	LastLoginAt       *time.Time `json:"last_login_at"`
	LastLoginIP       string     `json:"last_login_ip" gorm:"size:45"`
	FailedLoginCount  int        `json:"failed_login_count" gorm:"default:0"`
	LockedUntil       *time.Time `json:"locked_until"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`
	Credits           float64    `json:"credits" gorm:"type:decimal(12,2);default:0"`
	RoleID            uuid.UUID  `json:"role_id" gorm:"type:uuid"`
	Role              *Role      `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	ResellerID        *uuid.UUID `json:"reseller_id" gorm:"type:uuid"`
	Reseller          *User      `json:"reseller,omitempty" gorm:"foreignKey:ResellerID"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive && u.DeletedAt == nil
}

// IsLocked checks if the user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return u.FirstName + " " + u.LastName
}

// Role represents a user role with permissions
type Role struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string       `json:"name" gorm:"uniqueIndex;not null;size:50"`
	DisplayName string       `json:"display_name" gorm:"size:100"`
	Description string       `json:"description" gorm:"size:500"`
	Color       string       `json:"color" gorm:"size:7"` // Hex color
	IsSystem    bool         `json:"is_system" gorm:"default:false"`
	IsDefault   bool         `json:"is_default" gorm:"default:false"`
	Priority    int          `json:"priority" gorm:"default:0"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Role
func (Role) TableName() string {
	return "roles"
}

// Permission represents a system permission
type Permission struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:100"`
	DisplayName string    `json:"display_name" gorm:"size:150"`
	Description string    `json:"description" gorm:"size:500"`
	Category    string    `json:"category" gorm:"size:50;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for Permission
func (Permission) TableName() string {
	return "permissions"
}

// Session represents a user session
type Session struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	User         *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Token        string     `json:"-" gorm:"uniqueIndex;not null"`
	RefreshToken string     `json:"-" gorm:"uniqueIndex;not null"`
	IPAddress    string     `json:"ip_address" gorm:"size:45"`
	UserAgent    string     `json:"user_agent" gorm:"size:500"`
	LastActivity time.Time  `json:"last_activity"`
	ExpiresAt    time.Time  `json:"expires_at"`
	RevokedAt    *time.Time `json:"revoked_at"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for Session
func (Session) TableName() string {
	return "sessions"
}

// IsValid checks if the session is still valid
func (s *Session) IsValid() bool {
	return s.RevokedAt == nil && time.Now().Before(s.ExpiresAt)
}

// APIKey represents an API key for programmatic access
type APIKey struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	User        *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name        string     `json:"name" gorm:"not null;size:100"`
	KeyHash     string     `json:"-" gorm:"not null"`
	KeyPrefix   string     `json:"key_prefix" gorm:"size:10"` // First 8 chars for identification
	Permissions []string   `json:"permissions" gorm:"type:jsonb"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for APIKey
func (APIKey) TableName() string {
	return "api_keys"
}
