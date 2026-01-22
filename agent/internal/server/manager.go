package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aetherpanel/aether-panel/agent/internal/config"
	"github.com/aetherpanel/aether-panel/agent/internal/docker"
	"go.uber.org/zap"
)

// ServerState represents the state of a managed server
type ServerState struct {
	ID          string
	UUID        string
	ContainerID string
	Status      string
	StartedAt   *time.Time
	Stats       *ServerStats
	mu          sync.RWMutex
}

// ServerStats represents server resource usage
type ServerStats struct {
	CPUPercent    float64   `json:"cpu_percent"`
	MemoryUsage   uint64    `json:"memory_usage"`
	MemoryLimit   uint64    `json:"memory_limit"`
	DiskUsage     uint64    `json:"disk_usage"`
	DiskLimit     uint64    `json:"disk_limit"`
	NetworkRx     uint64    `json:"network_rx"`
	NetworkTx     uint64    `json:"network_tx"`
	Uptime        int64     `json:"uptime"`
	CollectedAt   time.Time `json:"collected_at"`
}

// Manager manages game servers on this node
type Manager struct {
	docker  *docker.Client
	config  *config.Config
	logger  *zap.Logger
	servers map[string]*ServerState
	mu      sync.RWMutex
}

// NewManager creates a new server manager
func NewManager(dockerClient *docker.Client, cfg *config.Config, logger *zap.Logger) *Manager {
	return &Manager{
		docker:  dockerClient,
		config:  cfg,
		logger:  logger,
		servers: make(map[string]*ServerState),
	}
}

// ServerConfig represents server configuration from panel
type ServerConfig struct {
	ID           string            `json:"id"`
	UUID         string            `json:"uuid"`
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	StartupCmd   string            `json:"startup_cmd"`
	Environment  map[string]string `json:"environment"`
	MemoryLimit  int64             `json:"memory_limit"`  // MB
	DiskLimit    int64             `json:"disk_limit"`    // MB
	CPULimit     int               `json:"cpu_limit"`     // percentage
	Allocations  []Allocation      `json:"allocations"`
	Mounts       []Mount           `json:"mounts"`
}

// Allocation represents a port allocation
type Allocation struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	IsPrimary bool   `json:"is_primary"`
}

// Mount represents a volume mount
type Mount struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	ReadOnly bool   `json:"read_only"`
}

// CreateServer creates and starts a new server
func (m *Manager) CreateServer(ctx context.Context, cfg *ServerConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("Creating server", zap.String("id", cfg.ID), zap.String("name", cfg.Name))

	// Create server data directory
	serverPath := filepath.Join(m.config.Storage.ServerDataPath, cfg.UUID)
	if err := os.MkdirAll(serverPath, 0755); err != nil {
		return fmt.Errorf("failed to create server directory: %w", err)
	}

	// Prepare environment variables
	var env []string
	for k, v := range cfg.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	// Prepare mounts
	var mounts []docker.MountConfig
	mounts = append(mounts, docker.MountConfig{
		Source: serverPath,
		Target: "/home/container",
	})
	for _, mount := range cfg.Mounts {
		mounts = append(mounts, docker.MountConfig{
			Source:   mount.Source,
			Target:   mount.Target,
			ReadOnly: mount.ReadOnly,
		})
	}

	// Prepare ports
	var ports []docker.PortConfig
	for _, alloc := range cfg.Allocations {
		ports = append(ports, docker.PortConfig{
			HostIP:   alloc.IP,
			HostPort: fmt.Sprintf("%d", alloc.Port),
			ContPort: fmt.Sprintf("%d", alloc.Port),
			Protocol: "tcp",
		})
		// Also bind UDP for game servers
		ports = append(ports, docker.PortConfig{
			HostIP:   alloc.IP,
			HostPort: fmt.Sprintf("%d", alloc.Port),
			ContPort: fmt.Sprintf("%d", alloc.Port),
			Protocol: "udp",
		})
	}

	// Pull image if needed
	exists, err := m.docker.ImageExists(ctx, cfg.Image)
	if err != nil {
		return fmt.Errorf("failed to check image: %w", err)
	}
	if !exists {
		m.logger.Info("Pulling image", zap.String("image", cfg.Image))
		if err := m.docker.PullImage(ctx, cfg.Image); err != nil {
			return fmt.Errorf("failed to pull image: %w", err)
		}
	}

	// Create container
	containerName := fmt.Sprintf("aether_%s", cfg.UUID)
	containerCfg := &docker.ContainerConfig{
		Name:        containerName,
		Image:       cfg.Image,
		Cmd:         []string{"/bin/bash", "-c", cfg.StartupCmd},
		Env:         env,
		WorkingDir:  "/home/container",
		User:        "container",
		Labels: map[string]string{
			"aether.server.id":   cfg.ID,
			"aether.server.uuid": cfg.UUID,
			"aether.managed":     "true",
		},
		Mounts:      mounts,
		Ports:       ports,
		Memory:      cfg.MemoryLimit * 1024 * 1024,      // Convert MB to bytes
		MemorySwap:  cfg.MemoryLimit * 1024 * 1024 * 2,  // 2x memory for swap
		CPUQuota:    int64(cfg.CPULimit) * 1000,         // CPU quota
		CPUPeriod:   100000,                              // 100ms period
		IOWeight:    500,
		NetworkMode: m.config.Docker.NetworkMode,
		DNS:         m.config.Docker.DNS,
		StopTimeout: m.config.Docker.StopTimeout,
	}

	containerID, err := m.docker.CreateContainer(ctx, containerCfg)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Store server state
	m.servers[cfg.ID] = &ServerState{
		ID:          cfg.ID,
		UUID:        cfg.UUID,
		ContainerID: containerID,
		Status:      "created",
	}

	m.logger.Info("Server created", zap.String("id", cfg.ID), zap.String("container", containerID))
	return nil
}

// StartServer starts a server
func (m *Manager) StartServer(ctx context.Context, serverID string) error {
	m.mu.RLock()
	server, exists := m.servers[serverID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server not found: %s", serverID)
	}

	server.mu.Lock()
	defer server.mu.Unlock()

	if err := m.docker.StartContainer(ctx, server.ContainerID); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	now := time.Now()
	server.Status = "running"
	server.StartedAt = &now

	m.logger.Info("Server started", zap.String("id", serverID))
	return nil
}

// StopServer stops a server gracefully
func (m *Manager) StopServer(ctx context.Context, serverID string) error {
	m.mu.RLock()
	server, exists := m.servers[serverID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server not found: %s", serverID)
	}

	server.mu.Lock()
	defer server.mu.Unlock()

	if err := m.docker.StopContainer(ctx, server.ContainerID, m.config.Docker.StopTimeout); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	server.Status = "stopped"
	server.StartedAt = nil

	m.logger.Info("Server stopped", zap.String("id", serverID))
	return nil
}

// KillServer forcefully stops a server
func (m *Manager) KillServer(ctx context.Context, serverID string) error {
	m.mu.RLock()
	server, exists := m.servers[serverID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server not found: %s", serverID)
	}

	server.mu.Lock()
	defer server.mu.Unlock()

	if err := m.docker.KillContainer(ctx, server.ContainerID); err != nil {
		return fmt.Errorf("failed to kill container: %w", err)
	}

	server.Status = "stopped"
	server.StartedAt = nil

	m.logger.Info("Server killed", zap.String("id", serverID))
	return nil
}

// RestartServer restarts a server
func (m *Manager) RestartServer(ctx context.Context, serverID string) error {
	m.mu.RLock()
	server, exists := m.servers[serverID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server not found: %s", serverID)
	}

	server.mu.Lock()
	defer server.mu.Unlock()

	if err := m.docker.RestartContainer(ctx, server.ContainerID, m.config.Docker.StopTimeout); err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}

	now := time.Now()
	server.Status = "running"
	server.StartedAt = &now

	m.logger.Info("Server restarted", zap.String("id", serverID))
	return nil
}

// SendCommand sends a command to server console
func (m *Manager) SendCommand(ctx context.Context, serverID, command string) error {
	m.mu.RLock()
	server, exists := m.servers[serverID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server not found: %s", serverID)
	}

	// Execute command in container
	_, err := m.docker.ExecCommand(ctx, server.ContainerID, []string{"/bin/bash", "-c", command})
	return err
}

// GetServerStatus returns server status
func (m *Manager) GetServerStatus(ctx context.Context, serverID string) (string, error) {
	m.mu.RLock()
	server, exists := m.servers[serverID]
	m.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("server not found: %s", serverID)
	}

	return m.docker.GetContainerStatus(ctx, server.ContainerID)
}

// GetServerStats returns server resource statistics
func (m *Manager) GetServerStats(ctx context.Context, serverID string) (*ServerStats, error) {
	m.mu.RLock()
	server, exists := m.servers[serverID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("server not found: %s", serverID)
	}

	server.mu.RLock()
	defer server.mu.RUnlock()

	return server.Stats, nil
}

// StartHealthCheck starts the health check routine
func (m *Manager) StartHealthCheck(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.checkAllServers(ctx)
		}
	}
}

// checkAllServers checks health of all servers
func (m *Manager) checkAllServers(ctx context.Context) {
	m.mu.RLock()
	servers := make([]*ServerState, 0, len(m.servers))
	for _, s := range m.servers {
		servers = append(servers, s)
	}
	m.mu.RUnlock()

	for _, server := range servers {
		status, err := m.docker.GetContainerStatus(ctx, server.ContainerID)
		if err != nil {
			m.logger.Warn("Failed to get container status",
				zap.String("id", server.ID),
				zap.Error(err))
			continue
		}

		server.mu.Lock()
		server.Status = status
		server.mu.Unlock()
	}
}

// StartMetricsCollection starts collecting metrics for all servers
func (m *Manager) StartMetricsCollection(ctx context.Context) {
	interval := time.Duration(m.config.Metrics.CollectInterval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.collectAllMetrics(ctx)
		}
	}
}

// collectAllMetrics collects metrics for all servers
func (m *Manager) collectAllMetrics(ctx context.Context) {
	m.mu.RLock()
	servers := make([]*ServerState, 0, len(m.servers))
	for _, s := range m.servers {
		servers = append(servers, s)
	}
	m.mu.RUnlock()

	for _, server := range servers {
		stats, err := m.docker.GetContainerStats(ctx, server.ContainerID)
		if err != nil {
			continue
		}

		server.mu.Lock()
		server.Stats = &ServerStats{
			CPUPercent:  stats.CPUPercent,
			MemoryUsage: stats.MemoryUsage,
			MemoryLimit: stats.MemoryLimit,
			NetworkRx:   stats.NetworkRx,
			NetworkTx:   stats.NetworkTx,
			CollectedAt: time.Now(),
		}
		if server.StartedAt != nil {
			server.Stats.Uptime = int64(time.Since(*server.StartedAt).Seconds())
		}
		server.mu.Unlock()
	}
}

// StartConsoleStreaming starts console streaming for all servers
func (m *Manager) StartConsoleStreaming(ctx context.Context) {
	// Console streaming would be implemented here
	// Using Docker attach and Redis pub/sub for real-time console
	<-ctx.Done()
}

// DeleteServer removes a server
func (m *Manager) DeleteServer(ctx context.Context, serverID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	server, exists := m.servers[serverID]
	if !exists {
		return fmt.Errorf("server not found: %s", serverID)
	}

	// Remove container
	if err := m.docker.RemoveContainer(ctx, server.ContainerID, true); err != nil {
		m.logger.Warn("Failed to remove container", zap.Error(err))
	}

	// Remove from map
	delete(m.servers, serverID)

	m.logger.Info("Server deleted", zap.String("id", serverID))
	return nil
}

// LoadServers loads existing servers from panel configuration
func (m *Manager) LoadServers(ctx context.Context, configs []ServerConfig) error {
	for _, cfg := range configs {
		// Check if container already exists
		containers, err := m.docker.ListContainersByLabel(ctx, map[string]string{
			"aether.server.id": cfg.ID,
		})
		if err != nil {
			m.logger.Warn("Failed to list containers", zap.Error(err))
			continue
		}

		if len(containers) > 0 {
			// Container exists, just track it
			container := containers[0]
			m.servers[cfg.ID] = &ServerState{
				ID:          cfg.ID,
				UUID:        cfg.UUID,
				ContainerID: container.ID,
				Status:      container.State,
			}
		}
	}

	return nil
}
