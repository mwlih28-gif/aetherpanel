module github.com/aetherpanel/aether-panel

go 1.22

require (
	github.com/gofiber/fiber/v2 v2.52.0
	github.com/gofiber/websocket/v2 v2.2.1
	github.com/gofiber/contrib/jwt v1.0.8
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.5.2
	github.com/redis/go-redis/v9 v9.4.0
	github.com/docker/docker v25.0.2+incompatible
	github.com/docker/go-connections v0.5.0
	github.com/spf13/viper v1.18.2
	github.com/go-playground/validator/v10 v10.17.0
	github.com/pquerna/otp v1.4.0
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.18.0
	gorm.io/gorm v1.25.6
	gorm.io/driver/postgres v1.5.4
	github.com/prometheus/client_golang v1.18.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/shirou/gopsutil/v3 v3.24.1
)
