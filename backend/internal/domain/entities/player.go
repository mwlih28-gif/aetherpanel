package entities

import (
	"time"

	"github.com/google/uuid"
)

// Player represents a game player
type Player struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID    uuid.UUID  `json:"server_id" gorm:"type:uuid;not null;index"`
	Server      *Server    `json:"server,omitempty" gorm:"foreignKey:ServerID"`
	UUID        string     `json:"uuid" gorm:"size:36;index"` // Minecraft UUID
	Username    string     `json:"username" gorm:"not null;size:50;index"`
	DisplayName string     `json:"display_name" gorm:"size:100"`
	SkinURL     string     `json:"skin_url" gorm:"size:500"`
	IsOnline    bool       `json:"is_online" gorm:"default:false"`
	IsBanned    bool       `json:"is_banned" gorm:"default:false"`
	IsOp        bool       `json:"is_op" gorm:"default:false"`
	IsWhitelisted bool     `json:"is_whitelisted" gorm:"default:false"`
	FirstJoinAt *time.Time `json:"first_join_at"`
	LastJoinAt  *time.Time `json:"last_join_at"`
	LastLeaveAt *time.Time `json:"last_leave_at"`
	PlayTime    int64      `json:"play_time" gorm:"default:0"` // Seconds
	JoinCount   int        `json:"join_count" gorm:"default:0"`
	IPAddress   string     `json:"ip_address" gorm:"size:45"`
	Country     string     `json:"country" gorm:"size:2"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Player) TableName() string {
	return "players"
}

// PlayerInventory represents a player's inventory snapshot
type PlayerInventory struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PlayerID    uuid.UUID `json:"player_id" gorm:"type:uuid;not null;index"`
	Player      *Player   `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	WorldID     *uuid.UUID `json:"world_id" gorm:"type:uuid"`
	World       *World    `json:"world,omitempty" gorm:"foreignKey:WorldID"`
	InventoryType string  `json:"inventory_type" gorm:"size:50"` // main, ender_chest, armor
	Slots       map[string]interface{} `json:"slots" gorm:"type:jsonb"`
	SnapshotAt  time.Time `json:"snapshot_at"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (PlayerInventory) TableName() string {
	return "player_inventories"
}

// PlayerStats represents player statistics
type PlayerStats struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PlayerID    uuid.UUID `json:"player_id" gorm:"type:uuid;not null;index"`
	Player      *Player   `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	StatType    string    `json:"stat_type" gorm:"size:50;index"` // kills, deaths, blocks_mined, etc.
	StatValue   int64     `json:"stat_value" gorm:"default:0"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PlayerStats) TableName() string {
	return "player_stats"
}

// ChatLog represents a chat message log
type ChatLog struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID  uuid.UUID `json:"server_id" gorm:"type:uuid;not null;index"`
	PlayerID  *uuid.UUID `json:"player_id" gorm:"type:uuid;index"`
	Player    *Player   `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	Username  string    `json:"username" gorm:"size:50"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	Channel   string    `json:"channel" gorm:"size:50"` // global, local, private, etc.
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}

func (ChatLog) TableName() string {
	return "chat_logs"
}

// CommandLog represents a command execution log
type CommandLog struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID  uuid.UUID `json:"server_id" gorm:"type:uuid;not null;index"`
	PlayerID  *uuid.UUID `json:"player_id" gorm:"type:uuid;index"`
	Player    *Player   `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	Username  string    `json:"username" gorm:"size:50"`
	Command   string    `json:"command" gorm:"type:text;not null"`
	IsConsole bool      `json:"is_console" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}

func (CommandLog) TableName() string {
	return "command_logs"
}

// DeathLog represents a player death event
type DeathLog struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServerID    uuid.UUID `json:"server_id" gorm:"type:uuid;not null;index"`
	PlayerID    uuid.UUID `json:"player_id" gorm:"type:uuid;not null;index"`
	Player      *Player   `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	KillerID    *uuid.UUID `json:"killer_id" gorm:"type:uuid"`
	Killer      *Player   `json:"killer,omitempty" gorm:"foreignKey:KillerID"`
	DeathCause  string    `json:"death_cause" gorm:"size:100"`
	DeathMessage string   `json:"death_message" gorm:"size:500"`
	WorldName   string    `json:"world_name" gorm:"size:100"`
	X           float64   `json:"x"`
	Y           float64   `json:"y"`
	Z           float64   `json:"z"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}

func (DeathLog) TableName() string {
	return "death_logs"
}
