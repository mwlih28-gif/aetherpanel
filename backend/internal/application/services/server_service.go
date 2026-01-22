package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/aetherpanel/aether-panel/internal/domain/repositories"
	"github.com/google/uuid"
)

var (
	ErrServerNotFound      = errors.New("server not found")
	ErrServerSuspended     = errors.New("server is suspended")
	ErrServerAlreadyRunning = errors.New("server is already running")
	ErrServerNotRunning    = errors.New("server is not running")
	ErrInsufficientResources = errors.New("insufficient resources on node")
	ErrNoAvailableAllocation = errors.New("no available allocation")
	ErrBackupLimitReached  = errors.New("backup limit reached")
)

// ServerService handles server operations
type ServerService struct {
	serverRepo     repositories.ServerRepository
	nodeRepo       repositories.NodeRepository
	allocationRepo repositories.AllocationRepository
	backupRepo     repositories.BackupRepository
	auditRepo      repositories.AuditLogRepository
	nodeClient     NodeClient
}

// NodeClient interface for communicating with node agents
type NodeClient interface {
	StartServer(ctx context.Context, nodeID uuid.UUID, serverID uuid.UUID) error
	StopServer(ctx context.Context, nodeID uuid.UUID, serverID uuid.UUID) error
	RestartServer(ctx context.Context, nodeID uuid.UUID, serverID uuid.UUID) error
	KillServer(ctx context.Context, nodeID uuid.UUID, serverID uuid.UUID) error
	GetServerStatus(ctx context.Context, nodeID uuid.UUID, serverID uuid.UUID) (*ServerStats, error)
	SendCommand(ctx context.Context, nodeID uuid.UUID, serverID uuid.UUID, command string) error
	CreateBackup(ctx context.Context, nodeID uuid.UUID, serverID uuid.UUID, backupID uuid.UUID) error
	RestoreBackup(ctx context.Context, nodeID uuid.UUID, serverID uuid.UUID, backupID uuid.UUID) error
	ReinstallServer(ctx context.Context, nodeID uuid.UUID, serverID uuid.UUID) error
}

// ServerStats represents server resource usage
type ServerStats struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage int64   `json:"memory_usage"`
	MemoryLimit int64   `json:"memory_limit"`
	DiskUsage   int64   `json:"disk_usage"`
	DiskLimit   int64   `json:"disk_limit"`
	NetworkRx   int64   `json:"network_rx"`
	NetworkTx   int64   `json:"network_tx"`
	Uptime      int64   `json:"uptime"`
	Status      string  `json:"status"`
}

// NewServerService creates a new ServerService
func NewServerService(
	serverRepo repositories.ServerRepository,
	nodeRepo repositories.NodeRepository,
	allocationRepo repositories.AllocationRepository,
	backupRepo repositories.BackupRepository,
	auditRepo repositories.AuditLogRepository,
	nodeClient NodeClient,
) *ServerService {
	return &ServerService{
		serverRepo:     serverRepo,
		nodeRepo:       nodeRepo,
		allocationRepo: allocationRepo,
		backupRepo:     backupRepo,
		auditRepo:      auditRepo,
		nodeClient:     nodeClient,
	}
}

// CreateServerRequest represents a server creation request
type CreateServerRequest struct {
	Name          string            `json:"name" validate:"required,min=1,max=100"`
	Description   string            `json:"description" validate:"max=500"`
	OwnerID       uuid.UUID         `json:"owner_id" validate:"required"`
	NodeID        uuid.UUID         `json:"node_id" validate:"required"`
	GameID        uuid.UUID         `json:"game_id" validate:"required"`
	EggID         uuid.UUID         `json:"egg_id" validate:"required"`
	MemoryLimit   int64             `json:"memory_limit" validate:"required,min=128"`
	DiskLimit     int64             `json:"disk_limit" validate:"required,min=1024"`
	CPULimit      int               `json:"cpu_limit" validate:"required,min=1,max=1000"`
	Environment   map[string]string `json:"environment"`
	StartOnCreate bool              `json:"start_on_create"`
}

// Create creates a new server
func (s *ServerService) Create(ctx context.Context, req *CreateServerRequest, createdBy uuid.UUID) (*entities.Server, error) {
	// Verify node exists and has capacity
	node, err := s.nodeRepo.GetByID(ctx, req.NodeID)
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	if node.AvailableMemory() < req.MemoryLimit {
		return nil, ErrInsufficientResources
	}

	if node.AvailableDisk() < req.DiskLimit {
		return nil, ErrInsufficientResources
	}

	// Find available allocation
	allocations, err := s.allocationRepo.GetAvailableByNodeID(ctx, req.NodeID)
	if err != nil || len(allocations) == 0 {
		return nil, ErrNoAvailableAllocation
	}
	allocation := allocations[0]

	// Generate short UUID
	shortUUID := uuid.New().String()[:8]

	// Create server
	server := &entities.Server{
		UUID:         shortUUID,
		Name:         req.Name,
		Description:  req.Description,
		Status:       entities.ServerStatusInstalling,
		OwnerID:      req.OwnerID,
		NodeID:       req.NodeID,
		AllocationID: allocation.ID,
		GameID:       req.GameID,
		EggID:        req.EggID,
		MemoryLimit:  req.MemoryLimit,
		DiskLimit:    req.DiskLimit,
		CPULimit:     req.CPULimit,
		Environment:  req.Environment,
	}

	if err := s.serverRepo.Create(ctx, server); err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	// Assign allocation to server
	if err := s.allocationRepo.AssignToServer(ctx, allocation.ID, server.ID, true); err != nil {
		return nil, fmt.Errorf("failed to assign allocation: %w", err)
	}

	// Update node allocated resources
	if err := s.nodeRepo.UpdateResources(ctx, req.NodeID,
		node.MemoryAllocated+req.MemoryLimit,
		node.DiskAllocated+req.DiskLimit,
		node.CPUAllocated+req.CPULimit,
	); err != nil {
		return nil, fmt.Errorf("failed to update node resources: %w", err)
	}

	// Log audit
	s.logAudit(ctx, createdBy, entities.AuditActionCreate, "server", &server.ID)

	return server, nil
}

// GetByID retrieves a server by ID
func (s *ServerService) GetByID(ctx context.Context, id uuid.UUID) (*entities.Server, error) {
	server, err := s.serverRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrServerNotFound
	}
	return server, nil
}

// GetByUUID retrieves a server by short UUID
func (s *ServerService) GetByUUID(ctx context.Context, uuid string) (*entities.Server, error) {
	server, err := s.serverRepo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, ErrServerNotFound
	}
	return server, nil
}

// List retrieves servers with pagination
func (s *ServerService) List(ctx context.Context, params repositories.ListParams) ([]*entities.Server, int64, error) {
	return s.serverRepo.List(ctx, params)
}

// GetByOwner retrieves servers owned by a user
func (s *ServerService) GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*entities.Server, error) {
	return s.serverRepo.GetByOwnerID(ctx, ownerID)
}

// Start starts a server
func (s *ServerService) Start(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) error {
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return ErrServerNotFound
	}

	if server.Suspended {
		return ErrServerSuspended
	}

	if server.IsRunning() {
		return ErrServerAlreadyRunning
	}

	// Update status
	if err := s.serverRepo.UpdateStatus(ctx, serverID, entities.ServerStatusStarting); err != nil {
		return err
	}

	// Send start command to node
	if err := s.nodeClient.StartServer(ctx, server.NodeID, serverID); err != nil {
		_ = s.serverRepo.UpdateStatus(ctx, serverID, entities.ServerStatusError)
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Update last started
	now := time.Now()
	server.LastStartedAt = &now
	_ = s.serverRepo.Update(ctx, server)

	s.logAudit(ctx, userID, entities.AuditActionStart, "server", &serverID)
	return nil
}

// Stop stops a server gracefully
func (s *ServerService) Stop(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) error {
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return ErrServerNotFound
	}

	if !server.IsRunning() {
		return ErrServerNotRunning
	}

	if err := s.serverRepo.UpdateStatus(ctx, serverID, entities.ServerStatusStopping); err != nil {
		return err
	}

	if err := s.nodeClient.StopServer(ctx, server.NodeID, serverID); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	s.logAudit(ctx, userID, entities.AuditActionStop, "server", &serverID)
	return nil
}

// Restart restarts a server
func (s *ServerService) Restart(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) error {
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return ErrServerNotFound
	}

	if server.Suspended {
		return ErrServerSuspended
	}

	if err := s.serverRepo.UpdateStatus(ctx, serverID, entities.ServerStatusRestarting); err != nil {
		return err
	}

	if err := s.nodeClient.RestartServer(ctx, server.NodeID, serverID); err != nil {
		return fmt.Errorf("failed to restart server: %w", err)
	}

	s.logAudit(ctx, userID, entities.AuditActionRestart, "server", &serverID)
	return nil
}

// Kill forcefully stops a server
func (s *ServerService) Kill(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) error {
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return ErrServerNotFound
	}

	if err := s.nodeClient.KillServer(ctx, server.NodeID, serverID); err != nil {
		return fmt.Errorf("failed to kill server: %w", err)
	}

	_ = s.serverRepo.UpdateStatus(ctx, serverID, entities.ServerStatusStopped)
	s.logAudit(ctx, userID, entities.AuditActionStop, "server", &serverID)
	return nil
}

// SendCommand sends a command to the server console
func (s *ServerService) SendCommand(ctx context.Context, serverID uuid.UUID, command string, userID uuid.UUID) error {
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return ErrServerNotFound
	}

	if !server.IsRunning() {
		return ErrServerNotRunning
	}

	if err := s.nodeClient.SendCommand(ctx, server.NodeID, serverID, command); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	s.logAudit(ctx, userID, entities.AuditActionCommand, "server", &serverID)
	return nil
}

// GetStats retrieves server resource statistics
func (s *ServerService) GetStats(ctx context.Context, serverID uuid.UUID) (*ServerStats, error) {
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return nil, ErrServerNotFound
	}

	return s.nodeClient.GetServerStatus(ctx, server.NodeID, serverID)
}

// Suspend suspends a server
func (s *ServerService) Suspend(ctx context.Context, serverID uuid.UUID, reason string, userID uuid.UUID) error {
	// Stop server first if running
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return ErrServerNotFound
	}

	if server.IsRunning() {
		_ = s.nodeClient.StopServer(ctx, server.NodeID, serverID)
	}

	if err := s.serverRepo.Suspend(ctx, serverID, reason); err != nil {
		return err
	}

	s.logAudit(ctx, userID, entities.AuditActionUpdate, "server", &serverID)
	return nil
}

// Unsuspend unsuspends a server
func (s *ServerService) Unsuspend(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) error {
	if err := s.serverRepo.Unsuspend(ctx, serverID); err != nil {
		return err
	}

	s.logAudit(ctx, userID, entities.AuditActionUpdate, "server", &serverID)
	return nil
}

// CreateBackup creates a server backup
func (s *ServerService) CreateBackup(ctx context.Context, serverID uuid.UUID, name string, userID uuid.UUID) (*entities.Backup, error) {
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return nil, ErrServerNotFound
	}

	// Check backup limit
	count, err := s.backupRepo.CountByServerID(ctx, serverID)
	if err != nil {
		return nil, err
	}

	if int(count) >= server.BackupLimit {
		return nil, ErrBackupLimitReached
	}

	// Create backup record
	backup := &entities.Backup{
		ServerID: serverID,
		Name:     name,
		Status:   entities.BackupStatusPending,
	}

	if err := s.backupRepo.Create(ctx, backup); err != nil {
		return nil, err
	}

	// Trigger backup on node
	if err := s.nodeClient.CreateBackup(ctx, server.NodeID, serverID, backup.ID); err != nil {
		_ = s.backupRepo.UpdateStatus(ctx, backup.ID, entities.BackupStatusFailed)
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	s.logAudit(ctx, userID, entities.AuditActionBackup, "server", &serverID)
	return backup, nil
}

// Reinstall reinstalls a server
func (s *ServerService) Reinstall(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) error {
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return ErrServerNotFound
	}

	// Stop if running
	if server.IsRunning() {
		_ = s.nodeClient.StopServer(ctx, server.NodeID, serverID)
	}

	if err := s.serverRepo.UpdateStatus(ctx, serverID, entities.ServerStatusInstalling); err != nil {
		return err
	}

	if err := s.nodeClient.ReinstallServer(ctx, server.NodeID, serverID); err != nil {
		return fmt.Errorf("failed to reinstall server: %w", err)
	}

	s.logAudit(ctx, userID, entities.AuditActionInstall, "server", &serverID)
	return nil
}

// Delete deletes a server
func (s *ServerService) Delete(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) error {
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return ErrServerNotFound
	}

	// Stop if running
	if server.IsRunning() {
		_ = s.nodeClient.KillServer(ctx, server.NodeID, serverID)
	}

	// Free allocation
	_ = s.allocationRepo.Unassign(ctx, server.AllocationID)

	// Update node resources
	node, _ := s.nodeRepo.GetByID(ctx, server.NodeID)
	if node != nil {
		_ = s.nodeRepo.UpdateResources(ctx, server.NodeID,
			node.MemoryAllocated-server.MemoryLimit,
			node.DiskAllocated-server.DiskLimit,
			node.CPUAllocated-server.CPULimit,
		)
	}

	// Delete server
	if err := s.serverRepo.Delete(ctx, serverID); err != nil {
		return err
	}

	s.logAudit(ctx, userID, entities.AuditActionDelete, "server", &serverID)
	return nil
}

func (s *ServerService) logAudit(ctx context.Context, userID uuid.UUID, action entities.AuditAction, resource string, resourceID *uuid.UUID) {
	log := &entities.AuditLog{
		UserID:     &userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
	}
	_ = s.auditRepo.Create(ctx, log)
}
