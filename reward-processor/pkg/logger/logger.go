package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// log is the singleton instance of the logger
	log *zap.Logger
)

// Config holds the logger configuration
type Config struct {
	Level      string
	Production bool
}

// Initialize sets up the logger with the given configuration
func Initialize(cfg Config) error {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	var config zap.Config
	if cfg.Production {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// Get returns the singleton logger instance
func Get() *zap.Logger {
	if log == nil {
		// Initialize with default config if not initialized
		_ = Initialize(Config{
			Level:      "info",
			Production: false,
		})
	}
	return log
}

// Sync flushes any buffered log entries
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}

// With creates a child logger with the given fields
func With(fields ...zap.Field) *zap.Logger {
	return Get().With(fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Fatal logs a fatal message and then calls os.Exit(1)
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}
