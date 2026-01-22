package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds agent configuration
type Config struct {
	NodeID      string      `mapstructure:"node_id"`
	Token       string      `mapstructure:"token"`
	Panel       PanelConfig `mapstructure:"panel"`
	API         APIConfig   `mapstructure:"api"`
	Docker      DockerConfig `mapstructure:"docker"`
	Storage     StorageConfig `mapstructure:"storage"`
	Metrics     MetricsConfig `mapstructure:"metrics"`
}

// PanelConfig holds panel connection settings
type PanelConfig struct {
	URL           string `mapstructure:"url"`
	Insecure      bool   `mapstructure:"insecure"`
	RetryInterval int    `mapstructure:"retry_interval"`
	MaxRetries    int    `mapstructure:"max_retries"`
}

// APIConfig holds agent API settings
type APIConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	TLSCert   string `mapstructure:"tls_cert"`
	TLSKey    string `mapstructure:"tls_key"`
}

// DockerConfig holds Docker settings
type DockerConfig struct {
	Socket          string            `mapstructure:"socket"`
	Network         string            `mapstructure:"network"`
	NetworkMode     string            `mapstructure:"network_mode"`
	DNS             []string          `mapstructure:"dns"`
	LogDriver       string            `mapstructure:"log_driver"`
	LogOpts         map[string]string `mapstructure:"log_opts"`
	StopTimeout     int               `mapstructure:"stop_timeout"`
	PullPolicy      string            `mapstructure:"pull_policy"`
}

// StorageConfig holds storage settings
type StorageConfig struct {
	ServerDataPath string `mapstructure:"server_data_path"`
	BackupPath     string `mapstructure:"backup_path"`
	TmpPath        string `mapstructure:"tmp_path"`
}

// MetricsConfig holds metrics settings
type MetricsConfig struct {
	Enabled         bool `mapstructure:"enabled"`
	CollectInterval int  `mapstructure:"collect_interval"`
	RetentionHours  int  `mapstructure:"retention_hours"`
}

// Load loads configuration from file and environment
func Load() (*Config, error) {
	v := viper.New()

	configPath := os.Getenv("AETHER_AGENT_CONFIG")
	if configPath == "" {
		configPath = "/etc/aether"
	}

	v.AddConfigPath(configPath)
	v.AddConfigPath(".")
	v.SetConfigName("agent")
	v.SetConfigType("yaml")

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	v.SetEnvPrefix("AETHER_AGENT")
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	// Panel defaults
	v.SetDefault("panel.url", "http://localhost:8080")
	v.SetDefault("panel.insecure", false)
	v.SetDefault("panel.retry_interval", 30)
	v.SetDefault("panel.max_retries", 10)

	// API defaults
	v.SetDefault("api.host", "0.0.0.0")
	v.SetDefault("api.port", 8443)

	// Docker defaults
	v.SetDefault("docker.socket", "unix:///var/run/docker.sock")
	v.SetDefault("docker.network", "aether_network")
	v.SetDefault("docker.network_mode", "bridge")
	v.SetDefault("docker.dns", []string{"1.1.1.1", "8.8.8.8"})
	v.SetDefault("docker.log_driver", "json-file")
	v.SetDefault("docker.log_opts", map[string]string{
		"max-size": "10m",
		"max-file": "3",
	})
	v.SetDefault("docker.stop_timeout", 30)
	v.SetDefault("docker.pull_policy", "if-not-present")

	// Storage defaults
	v.SetDefault("storage.server_data_path", "/var/lib/aether/servers")
	v.SetDefault("storage.backup_path", "/var/lib/aether/backups")
	v.SetDefault("storage.tmp_path", "/tmp/aether")

	// Metrics defaults
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.collect_interval", 5)
	v.SetDefault("metrics.retention_hours", 24)
}
