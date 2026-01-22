package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger instance
func NewLogger() *zap.Logger {
	env := os.Getenv("AETHER_ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	return logger
}

// NewFileLogger creates a logger that writes to a file
func NewFileLogger(filepath string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{filepath}
	config.ErrorOutputPaths = []string{filepath}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return config.Build(zap.AddCaller())
}

// Fields creates a slice of zap fields from key-value pairs
func Fields(keysAndValues ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key, ok := keysAndValues[i].(string)
			if !ok {
				continue
			}
			fields = append(fields, zap.Any(key, keysAndValues[i+1]))
		}
	}
	return fields
}
