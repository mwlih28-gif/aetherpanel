package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/aetherpanel/aether-panel/internal/domain/repositories"
	"github.com/google/uuid"
)

var (
	ErrNodeNotFound     = errors.New("node not found")
	ErrNodeOffline      = errors.New("node is offline")
	ErrNodeMaintenance  = errors.New("node is in maintenance mode")
	ErrLocationNotFound = errors.New("location not found")
)

// NodeService handles node operations
type NodeService struct {
	nodeRepo       repositories.NodeRepository
	locationRepo   repositories.LocationRepository
	allocationRepo repositories.AllocationRepository
	serverRepo     repositories.ServerRepository
	auditRepo      repositories.AuditLogRepository
}

// NewNodeService creates a new NodeService
func NewNodeService(
	nodeRepo repositories.NodeRepository,
	locationRepo repositories.LocationRepository,
	allocationRepo repositories.AllocationRepository,
	serverRepo repositories.ServerRepository,
	auditRepo repositories.AuditLogRepository,
) *NodeService {
	return &NodeService{
		nodeRepo:       nodeRepo,
		locationRepo:   locationRepo,
		allocationRepo: allocationRepo,
		serverRepo:     serverRepo,
		auditRepo:      auditRepo,
	}
}

// CreateNodeRequest represents a node creation request
type CreateNodeRequest struct {
	Name           string `json:"name" validate:"required,min=1,max=100"`
	Description    string `json:"description" validate:"max=500"`
	LocationID     uuid.UUID `json:"location_id" validate:"required"`
	FQDN           string `json:"fqdn" validate:"required,fqdn"`
	Scheme         string `json:"scheme" validate:"oneof=http https"`
	DaemonPort     int    `json:"daemon_port" validate:"required,min=1,max=65535"`
	MemoryTotal    int64  `json:"memory_total" validate:"required,min=1024"`
	MemoryOveralloc int   `json:"memory_overalloc" validate:"min=0,max=100"`
	DiskTotal      int64  `json:"disk_total" validate:"required,min=10240"`
	DiskOveralloc  int    `json:"disk_overalloc" validate:"min=0,max=100"`
	CPUTotal       int    `json:"cpu_total" validate:"required,min=100"`
}

// Create creates a new node
func (s *NodeService) Create(ctx context.Context, req *CreateNodeRequest, createdBy uuid.UUID) (*entities.Node, error) {
	// Verify location exists
	if _, err := s.locationRepo.GetByID(ctx, req.LocationID); err != nil {
		return nil, ErrLocationNotFound
	}

	// Generate daemon token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	daemonToken := hex.EncodeToString(tokenBytes)

	node := &entities.Node{
		Name:            req.Name,
		Description:     req.Description,
		LocationID:      req.LocationID,
		FQDN:            req.FQDN,
		Scheme:          req.Scheme,
		DaemonPort:      req.DaemonPort,
		DaemonToken:     daemonToken,
		MemoryTotal:     req.MemoryTotal,
		MemoryOveralloc: req.MemoryOveralloc,
		DiskTotal:       req.DiskTotal,
		DiskOveralloc:   req.DiskOveralloc,
		CPUTotal:        req.CPUTotal,
	}

	if err := s.nodeRepo.Create(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	s.logAudit(ctx, createdBy, entities.AuditActionCreate, "node", &node.ID)
	return node, nil
}

// GetByID retrieves a node by ID
func (s *NodeService) GetByID(ctx context.Context, id uuid.UUID) (*entities.Node, error) {
	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNodeNotFound
	}
	return node, nil
}

// List retrieves all nodes with pagination
func (s *NodeService) List(ctx context.Context, params repositories.ListParams) ([]*entities.Node, int64, error) {
	return s.nodeRepo.List(ctx, params)
}

// GetByLocation retrieves nodes by location
func (s *NodeService) GetByLocation(ctx context.Context, locationID uuid.UUID) ([]*entities.Node, error) {
	return s.nodeRepo.GetByLocationID(ctx, locationID)
}

// Update updates a node
func (s *NodeService) Update(ctx context.Context, id uuid.UUID, req *CreateNodeRequest, updatedBy uuid.UUID) (*entities.Node, error) {
	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNodeNotFound
	}

	node.Name = req.Name
	node.Description = req.Description
	node.LocationID = req.LocationID
	node.FQDN = req.FQDN
	node.Scheme = req.Scheme
	node.DaemonPort = req.DaemonPort
	node.MemoryTotal = req.MemoryTotal
	node.MemoryOveralloc = req.MemoryOveralloc
	node.DiskTotal = req.DiskTotal
	node.DiskOveralloc = req.DiskOveralloc
	node.CPUTotal = req.CPUTotal

	if err := s.nodeRepo.Update(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}

	s.logAudit(ctx, updatedBy, entities.AuditActionUpdate, "node", &id)
	return node, nil
}

// Delete deletes a node
func (s *NodeService) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// Check if node has servers
	count, err := s.serverRepo.CountByNodeID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete node with active servers")
	}

	if err := s.nodeRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.logAudit(ctx, deletedBy, entities.AuditActionDelete, "node", &id)
	return nil
}

// SetMaintenanceMode enables/disables maintenance mode
func (s *NodeService) SetMaintenanceMode(ctx context.Context, id uuid.UUID, maintenance bool, userID uuid.UUID) error {
	if err := s.nodeRepo.SetMaintenanceMode(ctx, id, maintenance); err != nil {
		return err
	}

	s.logAudit(ctx, userID, entities.AuditActionUpdate, "node", &id)
	return nil
}

// UpdateOnlineStatus updates node online status
func (s *NodeService) UpdateOnlineStatus(ctx context.Context, id uuid.UUID, isOnline bool) error {
	return s.nodeRepo.UpdateOnlineStatus(ctx, id, isOnline)
}

// RegenerateToken regenerates the daemon token
func (s *NodeService) RegenerateToken(ctx context.Context, id uuid.UUID, userID uuid.UUID) (string, error) {
	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return "", ErrNodeNotFound
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	newToken := hex.EncodeToString(tokenBytes)

	node.DaemonToken = newToken
	if err := s.nodeRepo.Update(ctx, node); err != nil {
		return "", err
	}

	s.logAudit(ctx, userID, entities.AuditActionUpdate, "node", &id)
	return newToken, nil
}

// GetConfiguration returns node configuration for the agent
func (s *NodeService) GetConfiguration(ctx context.Context, id uuid.UUID) (*NodeConfiguration, error) {
	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNodeNotFound
	}

	servers, err := s.serverRepo.GetByNodeID(ctx, id)
	if err != nil {
		return nil, err
	}

	allocations, err := s.allocationRepo.GetByNodeID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &NodeConfiguration{
		Node:        node,
		Servers:     servers,
		Allocations: allocations,
	}, nil
}

// NodeConfiguration represents configuration sent to node agent
type NodeConfiguration struct {
	Node        *entities.Node        `json:"node"`
	Servers     []*entities.Server    `json:"servers"`
	Allocations []*entities.Allocation `json:"allocations"`
}

// NodeStats represents node resource statistics
type NodeStats struct {
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryUsed     int64   `json:"memory_used"`
	MemoryTotal    int64   `json:"memory_total"`
	DiskUsed       int64   `json:"disk_used"`
	DiskTotal      int64   `json:"disk_total"`
	NetworkRx      int64   `json:"network_rx"`
	NetworkTx      int64   `json:"network_tx"`
	Uptime         int64   `json:"uptime"`
	ServerCount    int     `json:"server_count"`
	RunningServers int     `json:"running_servers"`
	LastUpdated    time.Time `json:"last_updated"`
}

// CreateAllocationRequest represents an allocation creation request
type CreateAllocationRequest struct {
	NodeID    uuid.UUID `json:"node_id" validate:"required"`
	IP        string    `json:"ip" validate:"required,ip"`
	PortStart int       `json:"port_start" validate:"required,min=1,max=65535"`
	PortEnd   int       `json:"port_end" validate:"required,min=1,max=65535,gtefield=PortStart"`
	Alias     string    `json:"alias" validate:"max=255"`
}

// CreateAllocations creates multiple allocations for a node
func (s *NodeService) CreateAllocations(ctx context.Context, req *CreateAllocationRequest, createdBy uuid.UUID) ([]*entities.Allocation, error) {
	// Verify node exists
	if _, err := s.nodeRepo.GetByID(ctx, req.NodeID); err != nil {
		return nil, ErrNodeNotFound
	}

	var allocations []*entities.Allocation
	for port := req.PortStart; port <= req.PortEnd; port++ {
		// Check if port is available
		available, err := s.allocationRepo.IsPortAvailable(ctx, req.NodeID, req.IP, port)
		if err != nil {
			return nil, err
		}
		if !available {
			continue
		}

		allocation := &entities.Allocation{
			NodeID: req.NodeID,
			IP:     req.IP,
			Port:   port,
			Alias:  req.Alias,
		}
		allocations = append(allocations, allocation)
	}

	if len(allocations) == 0 {
		return nil, errors.New("no available ports in range")
	}

	if err := s.allocationRepo.CreateBatch(ctx, allocations); err != nil {
		return nil, fmt.Errorf("failed to create allocations: %w", err)
	}

	s.logAudit(ctx, createdBy, entities.AuditActionCreate, "allocation", nil)
	return allocations, nil
}

// DeleteAllocation deletes an allocation
func (s *NodeService) DeleteAllocation(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	allocation, err := s.allocationRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("allocation not found")
	}

	if allocation.ServerID != nil {
		return errors.New("cannot delete allocation assigned to a server")
	}

	if err := s.allocationRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.logAudit(ctx, deletedBy, entities.AuditActionDelete, "allocation", &id)
	return nil
}

// Location operations

// CreateLocation creates a new location
func (s *NodeService) CreateLocation(ctx context.Context, location *entities.Location, createdBy uuid.UUID) error {
	if err := s.locationRepo.Create(ctx, location); err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}

	s.logAudit(ctx, createdBy, entities.AuditActionCreate, "location", &location.ID)
	return nil
}

// GetLocations retrieves all locations
func (s *NodeService) GetLocations(ctx context.Context) ([]*entities.Location, error) {
	return s.locationRepo.List(ctx)
}

// UpdateLocation updates a location
func (s *NodeService) UpdateLocation(ctx context.Context, location *entities.Location, updatedBy uuid.UUID) error {
	if err := s.locationRepo.Update(ctx, location); err != nil {
		return err
	}

	s.logAudit(ctx, updatedBy, entities.AuditActionUpdate, "location", &location.ID)
	return nil
}

// DeleteLocation deletes a location
func (s *NodeService) DeleteLocation(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// Check if location has nodes
	nodes, err := s.nodeRepo.GetByLocationID(ctx, id)
	if err != nil {
		return err
	}
	if len(nodes) > 0 {
		return errors.New("cannot delete location with active nodes")
	}

	if err := s.locationRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.logAudit(ctx, deletedBy, entities.AuditActionDelete, "location", &id)
	return nil
}

func (s *NodeService) logAudit(ctx context.Context, userID uuid.UUID, action entities.AuditAction, resource string, resourceID *uuid.UUID) {
	log := &entities.AuditLog{
		UserID:     &userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
	}
	_ = s.auditRepo.Create(ctx, log)
}
