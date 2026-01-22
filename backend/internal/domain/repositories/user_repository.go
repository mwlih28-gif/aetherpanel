package repositories

import (
	"context"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	
	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error
	
	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves users with pagination and filters
	List(ctx context.Context, params ListParams) ([]*entities.User, int64, error)
	
	// GetByResellerID retrieves users by reseller ID
	GetByResellerID(ctx context.Context, resellerID uuid.UUID) ([]*entities.User, error)
	
	// UpdateCredits updates user credits
	UpdateCredits(ctx context.Context, id uuid.UUID, amount float64) error
	
	// IncrementFailedLogin increments failed login count
	IncrementFailedLogin(ctx context.Context, id uuid.UUID) error
	
	// ResetFailedLogin resets failed login count
	ResetFailedLogin(ctx context.Context, id uuid.UUID) error
	
	// UpdateLastLogin updates last login info
	UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string) error
}

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	Create(ctx context.Context, role *entities.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Role, error)
	GetByName(ctx context.Context, name string) (*entities.Role, error)
	GetDefault(ctx context.Context) (*entities.Role, error)
	Update(ctx context.Context, role *entities.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*entities.Role, error)
	AssignPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	RemovePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
}

// PermissionRepository defines the interface for permission data access
type PermissionRepository interface {
	Create(ctx context.Context, permission *entities.Permission) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error)
	GetByName(ctx context.Context, name string) (*entities.Permission, error)
	List(ctx context.Context) ([]*entities.Permission, error)
	GetByCategory(ctx context.Context, category string) ([]*entities.Permission, error)
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]*entities.Permission, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Permission, error)
}

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	Create(ctx context.Context, session *entities.Session) error
	GetByToken(ctx context.Context, token string) (*entities.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*entities.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Session, error)
	Update(ctx context.Context, session *entities.Session) error
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// APIKeyRepository defines the interface for API key data access
type APIKeyRepository interface {
	Create(ctx context.Context, apiKey *entities.APIKey) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.APIKey, error)
	GetByKeyHash(ctx context.Context, keyHash string) (*entities.APIKey, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.APIKey, error)
	Update(ctx context.Context, apiKey *entities.APIKey) error
	Revoke(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ListParams represents common pagination and filter parameters
type ListParams struct {
	Page     int
	PageSize int
	SortBy   string
	SortDir  string // asc, desc
	Search   string
	Filters  map[string]interface{}
}

// DefaultListParams returns default list parameters
func DefaultListParams() ListParams {
	return ListParams{
		Page:     1,
		PageSize: 20,
		SortBy:   "created_at",
		SortDir:  "desc",
		Filters:  make(map[string]interface{}),
	}
}
