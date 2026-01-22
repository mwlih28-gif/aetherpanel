package entities

import (
	"time"

	"github.com/google/uuid"
)

// ServerStatus represents the current status of a game server
type ServerStatus string

const (
	ServerStatusInstalling ServerStatus = "installing"
	ServerStatusStarting   ServerStatus = "starting"
	ServerStatusRunning    ServerStatus = "running"
	ServerStatusStopping   ServerStatus = "stopping"
	ServerStatusStopped    ServerStatus = "stopped"
	ServerStatusRestarting ServerStatus = "restarting"
	ServerStatusError      ServerStatus = "error"
	ServerStatusSuspended  ServerStatus = "suspended"
)

// Server represents a game server instance
type Server struct {
	ID              uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UUID            string       `json:"uuid" gorm:"uniqueIndex;not null;size:36"` // Short UUID for URLs
	Name            string       `json:"name" gorm:"not null;size:100"`
	Description     string       `json:"description" gorm:"size:500"`
	Status          ServerStatus `json:"status" gorm:"type:varchar(20);default:'stopped'"`
	Suspended       bool         `json:"suspended" gorm:"default:false"`
	SuspendedReason string       `json:"suspended_reason" gorm:"size:500"`

	// Ownership
	OwnerID uuid.UUID `json:"owner_id" gorm:"type:uuid;not null;index"`
	Owner   *User     `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`

	// Node & Allocation
	NodeID       uuid.UUID   `json:"node_id" gorm:"type:uuid;not null;index"`
	Node         *Node       `json:"node,omitempty" gorm:"foreignKey:NodeID"`
	AllocationID uuid.UUID   `json:"allocation_id" gorm:"type:uuid;not null"`
	Allocation   *Allocation `json:"allocation,omitempty" gorm:"foreignKey:AllocationID"`

	// Game Configuration
	GameID      uuid.UUID `json:"game_id" gorm:"type:uuid;not null"`
	Game        *Game     `json:"game,omitempty" gorm:"foreignKey:GameID"`
	EggID       uuid.UUID `json:"egg_id" gorm:"type:uuid;not null"`
	Egg         *Egg      `json:"egg,omitempty" gorm:"foreignKey:EggID"`
	DockerImage string    `json:"docker_image" gorm:"size:255"`
	StartupCmd  string    `json:"startup_cmd" gorm:"type:text"`

	// Resource Limits
	MemoryLimit    int64 `json:"memory_limit" gorm:"default:1024"`     // MB
	SwapLimit      int64 `json:"swap_limit" gorm:"default:0"`          // MB
	DiskLimit      int64 `json:"disk_limit" gorm:"default:10240"`      // MB
	CPULimit       int   `json:"cpu_limit" gorm:"default:100"`         // Percentage (100 = 1 core)
	IOWeight       int   `json:"io_weight" gorm:"default:500"`         // 10-1000
	NetworkIn      int64 `json:"network_in" gorm:"default:0"`          // Bytes/s, 0 = unlimited
	NetworkOut     int64 `json:"network_out" gorm:"default:0"`         // Bytes/s, 0 = unlimited
	DatabaseLimit  int   `json:"database_limit" gorm:"default:0"`      // Number of databases
	AllocationLimit int  `json:"allocation_limit" gorm:"default:1"`    // Number of allocations
	BackupLimit    int   `json:"backup_limit" gorm:"default:2"`        // Number of backups

	// Environment Variables (stored as JSON)
	Environment map[string]string `json:"environment" gorm:"type:jsonb;default:'{}'"`

	// Container Info
	ContainerID   string `json:"container_id" gorm:"size:100"`
	InternalID    string `json:"internal_id" gorm:"size:100"` // Docker container name

	// Timestamps
	InstalledAt   *time.Time `json:"installed_at"`
	LastStartedAt *time.Time `json:"last_started_at"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName returns the table name for Server
func (Server) TableName() string {
	return "servers"
}

// IsRunning checks if the server is currently running
func (s *Server) IsRunning() bool {
	return s.Status == ServerStatusRunning
}

// CanStart checks if the server can be started
func (s *Server) CanStart() bool {
	return s.Status == ServerStatusStopped && !s.Suspended
}

// Node represents a physical or virtual machine running the agent
type Node struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string     `json:"name" gorm:"uniqueIndex;not null;size:100"`
	Description string     `json:"description" gorm:"size:500"`
	LocationID  uuid.UUID  `json:"location_id" gorm:"type:uuid;not null;index"`
	Location    *Location  `json:"location,omitempty" gorm:"foreignKey:LocationID"`
	FQDN        string     `json:"fqdn" gorm:"not null;size:255"` // Fully Qualified Domain Name
	Scheme      string     `json:"scheme" gorm:"size:10;default:'https'"`
	DaemonPort  int        `json:"daemon_port" gorm:"default:8443"`
	DaemonToken string     `json:"-" gorm:"size:100"`
	
	// Resource Capacity
	MemoryTotal      int64 `json:"memory_total" gorm:"default:0"`      // MB
	MemoryAllocated  int64 `json:"memory_allocated" gorm:"default:0"`  // MB
	MemoryOveralloc  int   `json:"memory_overalloc" gorm:"default:0"`  // Percentage
	DiskTotal        int64 `json:"disk_total" gorm:"default:0"`        // MB
	DiskAllocated    int64 `json:"disk_allocated" gorm:"default:0"`    // MB
	DiskOveralloc    int   `json:"disk_overalloc" gorm:"default:0"`    // Percentage
	CPUTotal         int   `json:"cpu_total" gorm:"default:0"`         // Percentage (100 = 1 core)
	CPUAllocated     int   `json:"cpu_allocated" gorm:"default:0"`     // Percentage

	// Status
	IsOnline        bool       `json:"is_online" gorm:"default:false"`
	LastCheckedAt   *time.Time `json:"last_checked_at"`
	MaintenanceMode bool       `json:"maintenance_mode" gorm:"default:false"`
	
	// System Info (populated by agent)
	SystemInfo map[string]interface{} `json:"system_info" gorm:"type:jsonb;default:'{}'"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName returns the table name for Node
func (Node) TableName() string {
	return "nodes"
}

// AvailableMemory returns the available memory on the node
func (n *Node) AvailableMemory() int64 {
	maxMemory := n.MemoryTotal + (n.MemoryTotal * int64(n.MemoryOveralloc) / 100)
	return maxMemory - n.MemoryAllocated
}

// AvailableDisk returns the available disk space on the node
func (n *Node) AvailableDisk() int64 {
	maxDisk := n.DiskTotal + (n.DiskTotal * int64(n.DiskOveralloc) / 100)
	return maxDisk - n.DiskAllocated
}

// Location represents a physical location/datacenter
type Location struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ShortCode   string     `json:"short_code" gorm:"uniqueIndex;not null;size:10"`
	Name        string     `json:"name" gorm:"not null;size:100"`
	Description string     `json:"description" gorm:"size:500"`
	Country     string     `json:"country" gorm:"size:2"` // ISO 3166-1 alpha-2
	City        string     `json:"city" gorm:"size:100"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Location
func (Location) TableName() string {
	return "locations"
}

// Allocation represents an IP:Port allocation on a node
type Allocation struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NodeID    uuid.UUID  `json:"node_id" gorm:"type:uuid;not null;index"`
	Node      *Node      `json:"node,omitempty" gorm:"foreignKey:NodeID"`
	IP        string     `json:"ip" gorm:"not null;size:45"`
	Port      int        `json:"port" gorm:"not null"`
	Alias     string     `json:"alias" gorm:"size:255"` // Optional friendly name
	Notes     string     `json:"notes" gorm:"size:500"`
	ServerID  *uuid.UUID `json:"server_id" gorm:"type:uuid;index"`
	Server    *Server    `json:"server,omitempty" gorm:"foreignKey:ServerID"`
	IsPrimary bool       `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Allocation
func (Allocation) TableName() string {
	return "allocations"
}

// Address returns the full address (IP:Port)
func (a *Allocation) Address() string {
	return a.IP + ":" + string(rune(a.Port))
}

// Game represents a supported game type
type Game struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:100"`
	Description string    `json:"description" gorm:"size:500"`
	Icon        string    `json:"icon" gorm:"size:500"`
	Category    string    `json:"category" gorm:"size:50;index"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Game
func (Game) TableName() string {
	return "games"
}

// Egg represents a server configuration template (like Pterodactyl eggs)
type Egg struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GameID          uuid.UUID `json:"game_id" gorm:"type:uuid;not null;index"`
	Game            *Game     `json:"game,omitempty" gorm:"foreignKey:GameID"`
	Name            string    `json:"name" gorm:"not null;size:100"`
	Description     string    `json:"description" gorm:"type:text"`
	Author          string    `json:"author" gorm:"size:100"`
	DockerImages    []string  `json:"docker_images" gorm:"type:jsonb"`
	StartupCommand  string    `json:"startup_command" gorm:"type:text"`
	ConfigFiles     string    `json:"config_files" gorm:"type:jsonb"`
	ConfigStartup   string    `json:"config_startup" gorm:"type:jsonb"`
	ConfigStop      string    `json:"config_stop" gorm:"size:100"`
	ConfigLogs      string    `json:"config_logs" gorm:"type:jsonb"`
	InstallScript   string    `json:"install_script" gorm:"type:text"`
	InstallContainer string   `json:"install_container" gorm:"size:255"`
	InstallEntrypoint string  `json:"install_entrypoint" gorm:"size:255"`
	Variables       []EggVariable `json:"variables,omitempty" gorm:"foreignKey:EggID"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	SortOrder       int       `json:"sort_order" gorm:"default:0"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Egg
func (Egg) TableName() string {
	return "eggs"
}

// EggVariable represents a configurable variable for an egg
type EggVariable struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EggID        uuid.UUID `json:"egg_id" gorm:"type:uuid;not null;index"`
	Name         string    `json:"name" gorm:"not null;size:100"`
	Description  string    `json:"description" gorm:"size:500"`
	EnvVariable  string    `json:"env_variable" gorm:"not null;size:100"`
	DefaultValue string    `json:"default_value" gorm:"size:500"`
	UserViewable bool      `json:"user_viewable" gorm:"default:true"`
	UserEditable bool      `json:"user_editable" gorm:"default:true"`
	Rules        string    `json:"rules" gorm:"size:500"` // Validation rules
	SortOrder    int       `json:"sort_order" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for EggVariable
func (EggVariable) TableName() string {
	return "egg_variables"
}

// ServerVariable represents a server-specific variable value
type ServerVariable struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID      uuid.UUID `json:"server_id" gorm:"type:uuid;not null;index"`
	EggVariableID uuid.UUID `json:"egg_variable_id" gorm:"type:uuid;not null"`
	EggVariable   *EggVariable `json:"egg_variable,omitempty" gorm:"foreignKey:EggVariableID"`
	Value         string    `json:"value" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ServerVariable
func (ServerVariable) TableName() string {
	return "server_variables"
}
