package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alexandredsa/learning-rewards/reward-processor/internal/database"
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/database/seed"
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/processor"
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/repository"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/logger"
	"go.uber.org/zap"
)

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

	// Initialize database
	db, err := database.Connect(getEnv("DATABASE_DSN", ""))
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Create repositories
	eventRepo := repository.NewGormUserEventRepository(db)
	ruleRepo := repository.NewGormRuleRepository(db)

	// Seed rules if needed
	ctx := context.Background()
	if err := seed.SeedRules(ctx, db, log); err != nil {
		log.Fatal("Failed to seed rules", zap.Error(err))
	}

	// Get enabled rules
	rules, err := ruleRepo.GetEnabledRules(ctx)
	if err != nil {
		log.Fatal("Failed to get rules", zap.Error(err))
	}

	// Get configuration from environment
	cfg := processor.Config{
		KafkaBrokers:   strings.Split(getEnv("KAFKA_BROKERS", "localhost:29092"), ","),
		ConsumerGroup:  getEnv("KAFKA_CONSUMER_GROUP", "reward-processor"),
		ConsumerTopics: strings.Split(getEnv("KAFKA_CONSUMER_TOPICS", "learning-events"), ","),
		ProducerTopic:  getEnv("KAFKA_PRODUCER_TOPIC", "user-rewards"),
		Rules:          rules,
	}

	// Create processor
	proc, err := processor.New(cfg, eventRepo, log)
	if err != nil {
		log.Fatal("Failed to create processor", zap.Error(err))
	}
	defer proc.Close()

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Info("Received shutdown signal", zap.Stringer("signal", sig))
		cancel()
	}()

	// Start processing
	if err := proc.Start(ctx); err != nil {
		log.Fatal("Processor error", zap.Error(err))
	}
}
