package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexandredsa/learning-rewards/reward-processor/internal/database"
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/repository"
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/server"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/logger"
	"go.uber.org/zap"
)

const (
	defaultPort     = "8100"
	shutdownTimeout = 10 * time.Second
)

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return defaultPort
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	// Initialize logger
	if err := logger.Initialize(logger.Config{
		Level:      getEnv("LOG_LEVEL", "info"),
		Production: getEnv("ENV", "development") == "production",
	}); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	defer logger.Sync()

	log := logger.Get()

	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connect to database using DSN from environment variable (or default from database package)
	db, err := database.Connect(os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Create repositories
	ruleRepo := repository.NewGormRuleRepository(db)

	// Get port from environment variable or use default
	port := getPort()

	// Create and start server
	srv := server.New(server.Config{
		Port: port,
	})

	// Start server in a goroutine
	go func() {
		if err := srv.Start(ruleRepo); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	log.Info("API server started", zap.String("port", port))

	// Wait for interrupt signal
	<-ctx.Done()
	log.Info("Shutting down API server...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown failed", zap.Error(err))
		os.Exit(1)
	}

	log.Info("Server shutdown complete")
}
