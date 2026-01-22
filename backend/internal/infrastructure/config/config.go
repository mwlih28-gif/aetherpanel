package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Security SecurityConfig `mapstructure:"security"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Mail     MailConfig     `mapstructure:"mail"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"` // development, staging, production
	Debug       bool   `mapstructure:"debug"`
	URL         string `mapstructure:"url"`
	Timezone    string `mapstructure:"timezone"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	BodyLimit       int           `mapstructure:"body_limit"` // MB
	TrustedProxies  []string      `mapstructure:"trusted_proxies"`
	CORS            CORSConfig    `mapstructure:"cors"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	LogLevel        string        `mapstructure:"log_level"` // silent, error, warn, info
}

// DSN returns the database connection string
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// Addr returns the Redis address
func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`
	AccessExpiry     time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry    time.Duration `mapstructure:"refresh_expiry"`
	Issuer           string        `mapstructure:"issuer"`
	Audience         string        `mapstructure:"audience"`
	CookieName       string        `mapstructure:"cookie_name"`
	CookieSecure     bool          `mapstructure:"cookie_secure"`
	CookieHTTPOnly   bool          `mapstructure:"cookie_http_only"`
	CookieSameSite   string        `mapstructure:"cookie_same_site"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	PasswordMinLength    int           `mapstructure:"password_min_length"`
	PasswordRequireUpper bool          `mapstructure:"password_require_upper"`
	PasswordRequireLower bool          `mapstructure:"password_require_lower"`
	PasswordRequireDigit bool          `mapstructure:"password_require_digit"`
	PasswordRequireSpecial bool        `mapstructure:"password_require_special"`
	MaxLoginAttempts     int           `mapstructure:"max_login_attempts"`
	LockoutDuration      time.Duration `mapstructure:"lockout_duration"`
	SessionMaxAge        time.Duration `mapstructure:"session_max_age"`
	TwoFactorIssuer      string        `mapstructure:"two_factor_issuer"`
	RateLimitRequests    int           `mapstructure:"rate_limit_requests"`
	RateLimitDuration    time.Duration `mapstructure:"rate_limit_duration"`
	EncryptionKey        string        `mapstructure:"encryption_key"` // 32 bytes for AES-256
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Driver    string      `mapstructure:"driver"` // local, s3
	LocalPath string      `mapstructure:"local_path"`
	S3        S3Config    `mapstructure:"s3"`
	Backups   BackupStorageConfig `mapstructure:"backups"`
}

// S3Config holds S3 storage configuration
type S3Config struct {
	Endpoint        string `mapstructure:"endpoint"`
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	UseSSL          bool   `mapstructure:"use_ssl"`
}

// BackupStorageConfig holds backup storage configuration
type BackupStorageConfig struct {
	Path           string `mapstructure:"path"`
	MaxSize        int64  `mapstructure:"max_size"` // MB
	RetentionDays  int    `mapstructure:"retention_days"`
}

// MailConfig holds mail configuration
type MailConfig struct {
	Driver     string `mapstructure:"driver"` // smtp, sendgrid, mailgun
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Encryption string `mapstructure:"encryption"` // tls, ssl, none
	FromName   string `mapstructure:"from_name"`
	FromEmail  string `mapstructure:"from_email"`
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Path       string `mapstructure:"path"`
	Port       int    `mapstructure:"port"`
	Prometheus bool   `mapstructure:"prometheus"`
}

// Load loads configuration from file and environment
func Load() (*Config, error) {
	v := viper.New()

	// Set config file
	configPath := os.Getenv("AETHER_CONFIG_PATH")
	if configPath == "" {
		configPath = "."
	}
	v.AddConfigPath(configPath)
	v.AddConfigPath("/etc/aether")
	v.AddConfigPath("$HOME/.aether")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Set defaults
	setDefaults(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Override with environment variables
	v.SetEnvPrefix("AETHER")
	v.AutomaticEnv()

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "Aether Panel")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", true)
	v.SetDefault("app.url", "http://localhost:3000")
	v.SetDefault("app.timezone", "UTC")

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "120s")
	v.SetDefault("server.shutdown_timeout", "30s")
	v.SetDefault("server.body_limit", 50)
	v.SetDefault("server.cors.enabled", true)
	v.SetDefault("server.cors.allow_origins", []string{"*"})
	v.SetDefault("server.cors.allow_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	v.SetDefault("server.cors.allow_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	v.SetDefault("server.cors.allow_credentials", true)
	v.SetDefault("server.cors.max_age", 86400)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "aether")
	v.SetDefault("database.password", "aether")
	v.SetDefault("database.name", "aether_panel")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "1h")
	v.SetDefault("database.conn_max_idle_time", "30m")
	v.SetDefault("database.log_level", "warn")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")

	// JWT defaults
	v.SetDefault("jwt.access_expiry", "15m")
	v.SetDefault("jwt.refresh_expiry", "7d")
	v.SetDefault("jwt.issuer", "aether-panel")
	v.SetDefault("jwt.audience", "aether-panel")
	v.SetDefault("jwt.cookie_name", "aether_token")
	v.SetDefault("jwt.cookie_secure", false)
	v.SetDefault("jwt.cookie_http_only", true)
	v.SetDefault("jwt.cookie_same_site", "lax")

	// Security defaults
	v.SetDefault("security.password_min_length", 8)
	v.SetDefault("security.password_require_upper", true)
	v.SetDefault("security.password_require_lower", true)
	v.SetDefault("security.password_require_digit", true)
	v.SetDefault("security.password_require_special", false)
	v.SetDefault("security.max_login_attempts", 5)
	v.SetDefault("security.lockout_duration", "15m")
	v.SetDefault("security.session_max_age", "24h")
	v.SetDefault("security.two_factor_issuer", "Aether Panel")
	v.SetDefault("security.rate_limit_requests", 100)
	v.SetDefault("security.rate_limit_duration", "1m")

	// Storage defaults
	v.SetDefault("storage.driver", "local")
	v.SetDefault("storage.local_path", "/var/lib/aether")
	v.SetDefault("storage.backups.path", "/var/lib/aether/backups")
	v.SetDefault("storage.backups.max_size", 10240)
	v.SetDefault("storage.backups.retention_days", 30)

	// Mail defaults
	v.SetDefault("mail.driver", "smtp")
	v.SetDefault("mail.port", 587)
	v.SetDefault("mail.encryption", "tls")
	v.SetDefault("mail.from_name", "Aether Panel")

	// Metrics defaults
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.path", "/metrics")
	v.SetDefault("metrics.port", 9090)
	v.SetDefault("metrics.prometheus", true)
}
