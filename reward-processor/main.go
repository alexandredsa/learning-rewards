package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alexandredsa/learning-rewards/reward-processor/internal/database"
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/processor"
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/repository"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/logger"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"go.uber.org/zap"
)

// TODO: Load rules from a DB instance.
var mockedRules = []models.Rule{
	{
		ID:        "rule-001",
		Type:      models.SingleEventRule,
		EventType: "COURSE_COMPLETED",
		Conditions: map[string]string{
			"category": "MATH",
		},
		Reward: models.Reward{
			Type:        models.BadgeReward,
			Description: "Finished a Math course",
		},
		Enabled: true,
	},
	{
		ID:        "rule-002",
		Type:      models.MilestoneRule,
		EventType: "COURSE_COMPLETED",
		Count:     5,
		Conditions: map[string]string{
			"category": "MATH",
		},
		Reward: models.Reward{
			Type:        models.PointsReward,
			Amount:      100,
			Description: "Completed 5 math courses",
		},
		Enabled: true,
	},
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

	// Initialize database
	db, err := database.Connect(getEnv("DATABASE_DSN", ""))
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Create event repository
	eventRepo := repository.NewGormUserEventRepository(db)

	// Get configuration from environment
	cfg := processor.Config{
		KafkaBrokers:   strings.Split(getEnv("KAFKA_BROKERS", "localhost:29092"), ","),
		ConsumerGroup:  getEnv("KAFKA_CONSUMER_GROUP", "reward-processor"),
		ConsumerTopics: strings.Split(getEnv("KAFKA_CONSUMER_TOPICS", "learning-events"), ","),
		ProducerTopic:  getEnv("KAFKA_PRODUCER_TOPIC", "user-rewards"),
		Rules:          mockedRules,
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
