package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Client wraps Docker client with additional functionality
type Client struct {
	cli *client.Client
}

// NewClient creates a new Docker client
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := cli.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping docker: %w", err)
	}

	return &Client{cli: cli}, nil
}

// Close closes the Docker client
func (c *Client) Close() error {
	return c.cli.Close()
}

// ContainerConfig represents container configuration
type ContainerConfig struct {
	Name        string
	Image       string
	Cmd         []string
	Env         []string
	WorkingDir  string
	User        string
	Labels      map[string]string
	Mounts      []MountConfig
	Ports       []PortConfig
	Memory      int64 // bytes
	MemorySwap  int64 // bytes
	CPUQuota    int64
	CPUPeriod   int64
	IOWeight    uint16
	NetworkMode string
	DNS         []string
	StopTimeout int
}

// MountConfig represents a mount configuration
type MountConfig struct {
	Source   string
	Target   string
	ReadOnly bool
}

// PortConfig represents a port configuration
type PortConfig struct {
	HostIP   string
	HostPort string
	ContPort string
	Protocol string // tcp or udp
}

// CreateContainer creates a new container
func (c *Client) CreateContainer(ctx context.Context, cfg *ContainerConfig) (string, error) {
	// Prepare environment
	env := cfg.Env

	// Prepare mounts
	var mounts []mount.Mount
	for _, m := range cfg.Mounts {
		mounts = append(mounts, mount.Mount{
			Type:     mount.TypeBind,
			Source:   m.Source,
			Target:   m.Target,
			ReadOnly: m.ReadOnly,
		})
	}

	// Prepare port bindings
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}
	for _, p := range cfg.Ports {
		port := nat.Port(fmt.Sprintf("%s/%s", p.ContPort, p.Protocol))
		exposedPorts[port] = struct{}{}
		portBindings[port] = []nat.PortBinding{
			{
				HostIP:   p.HostIP,
				HostPort: p.HostPort,
			},
		}
	}

	// Container config
	containerCfg := &container.Config{
		Image:        cfg.Image,
		Cmd:          cfg.Cmd,
		Env:          env,
		WorkingDir:   cfg.WorkingDir,
		User:         cfg.User,
		Labels:       cfg.Labels,
		ExposedPorts: exposedPorts,
		Tty:          true,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}

	// Host config
	stopTimeout := cfg.StopTimeout
	hostCfg := &container.HostConfig{
		Mounts:       mounts,
		PortBindings: portBindings,
		Resources: container.Resources{
			Memory:     cfg.Memory,
			MemorySwap: cfg.MemorySwap,
			CPUQuota:   cfg.CPUQuota,
			CPUPeriod:  cfg.CPUPeriod,
			BlkioWeight: cfg.IOWeight,
		},
		NetworkMode:   container.NetworkMode(cfg.NetworkMode),
		DNS:           cfg.DNS,
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
		StopTimeout:   &stopTimeout,
		LogConfig: container.LogConfig{
			Type: "json-file",
			Config: map[string]string{
				"max-size": "10m",
				"max-file": "3",
			},
		},
	}

	// Network config
	networkCfg := &network.NetworkingConfig{}

	// Create container
	resp, err := c.cli.ContainerCreate(ctx, containerCfg, hostCfg, networkCfg, nil, cfg.Name)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	return resp.ID, nil
}

// StartContainer starts a container
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	return c.cli.ContainerStart(ctx, containerID, container.StartOptions{})
}

// StopContainer stops a container gracefully
func (c *Client) StopContainer(ctx context.Context, containerID string, timeout int) error {
	timeoutDuration := time.Duration(timeout) * time.Second
	return c.cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout, Signal: "SIGTERM"})
	_ = timeoutDuration
}

// KillContainer forcefully stops a container
func (c *Client) KillContainer(ctx context.Context, containerID string) error {
	return c.cli.ContainerKill(ctx, containerID, "SIGKILL")
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(ctx context.Context, containerID string, timeout int) error {
	return c.cli.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout})
}

// RemoveContainer removes a container
func (c *Client) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	return c.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force:         force,
		RemoveVolumes: true,
	})
}

// ContainerStats represents container resource statistics
type ContainerStats struct {
	CPUPercent    float64
	MemoryUsage   uint64
	MemoryLimit   uint64
	MemoryPercent float64
	NetworkRx     uint64
	NetworkTx     uint64
	BlockRead     uint64
	BlockWrite    uint64
	PIDs          uint64
}

// GetContainerStats retrieves container statistics
func (c *Client) GetContainerStats(ctx context.Context, containerID string) (*ContainerStats, error) {
	stats, err := c.cli.ContainerStatsOneShot(ctx, containerID)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	// Parse stats - simplified version
	// In production, would parse the JSON stream properly
	return &ContainerStats{}, nil
}

// GetContainerStatus returns container status
func (c *Client) GetContainerStatus(ctx context.Context, containerID string) (string, error) {
	info, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}
	return info.State.Status, nil
}

// IsContainerRunning checks if container is running
func (c *Client) IsContainerRunning(ctx context.Context, containerID string) (bool, error) {
	status, err := c.GetContainerStatus(ctx, containerID)
	if err != nil {
		return false, err
	}
	return status == "running", nil
}

// ExecCommand executes a command in a container
func (c *Client) ExecCommand(ctx context.Context, containerID string, cmd []string) (string, error) {
	execCfg := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := c.cli.ContainerExecCreate(ctx, containerID, execCfg)
	if err != nil {
		return "", err
	}

	resp, err := c.cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer resp.Close()

	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// AttachContainer attaches to container stdin/stdout
func (c *Client) AttachContainer(ctx context.Context, containerID string) (types.HijackedResponse, error) {
	return c.cli.ContainerAttach(ctx, containerID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
}

// GetContainerLogs retrieves container logs
func (c *Client) GetContainerLogs(ctx context.Context, containerID string, tail string) (io.ReadCloser, error) {
	return c.cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
		Timestamps: true,
	})
}

// PullImage pulls a Docker image
func (c *Client) PullImage(ctx context.Context, image string) error {
	reader, err := c.cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	// Consume the output
	_, err = io.Copy(io.Discard, reader)
	return err
}

// ImageExists checks if an image exists locally
func (c *Client) ImageExists(ctx context.Context, image string) (bool, error) {
	_, _, err := c.cli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ListContainersByLabel lists containers with specific labels
func (c *Client) ListContainersByLabel(ctx context.Context, labels map[string]string) ([]types.Container, error) {
	filterArgs := filters.NewArgs()
	for k, v := range labels {
		filterArgs.Add("label", fmt.Sprintf("%s=%s", k, v))
	}

	return c.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filterArgs,
	})
}

// CreateNetwork creates a Docker network
func (c *Client) CreateNetwork(ctx context.Context, name string) error {
	_, err := c.cli.NetworkCreate(ctx, name, types.NetworkCreate{
		Driver: "bridge",
	})
	return err
}

// GetSystemInfo returns Docker system information
func (c *Client) GetSystemInfo(ctx context.Context) (types.Info, error) {
	return c.cli.Info(ctx)
}
