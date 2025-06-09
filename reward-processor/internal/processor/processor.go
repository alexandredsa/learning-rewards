package processor

import (
	"context"

	"github.com/alexandredsa/learning-rewards/reward-processor/internal/kafka"
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/repository"
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/rules"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"go.uber.org/zap"
)

// Config holds the processor configuration
type Config struct {
	KafkaBrokers   []string
	ConsumerGroup  string
	ConsumerTopics []string
	ProducerTopic  string
	Rules          []models.Rule
}

// Processor handles the reward processing logic
type Processor struct {
	consumer *kafka.Consumer
	producer *kafka.Producer
	engine   *rules.Engine
	logger   *zap.Logger
}

// New creates a new reward processor
func New(cfg Config, eventRepo repository.UserEventRepository, logger *zap.Logger) (*Processor, error) {
	// Create rules engine with repository
	engine := rules.NewEngine(cfg.Rules, eventRepo, logger)

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.ConsumerGroup,
		cfg.ConsumerTopics,
	)
	if err != nil {
		return nil, err
	}

	// Create Kafka producer
	producer, err := kafka.NewProducer(
		cfg.KafkaBrokers,
		cfg.ProducerTopic,
	)
	if err != nil {
		consumer.Close()
		return nil, err
	}

	p := &Processor{
		consumer: consumer,
		producer: producer,
		engine:   engine,
		logger:   logger,
	}

	// Set up event handler
	consumer.SetHandler(p.handleEvent)

	return p, nil
}

// handleEvent processes a single user event
func (p *Processor) handleEvent(event models.UserEvent) error {
	// Process event through rules engine
	triggered, err := p.engine.EvaluateEvent(context.Background(), event)
	if err != nil {
		p.logger.Error("Failed to evaluate event",
			zap.Error(err),
			zap.Any("event", event))
		return err
	}

	// Send triggered rewards
	for _, reward := range triggered {
		if err := p.producer.SendReward(reward); err != nil {
			p.logger.Error("Failed to send reward",
				zap.Error(err),
				zap.Any("reward", reward))
			return err
		}
	}

	return nil
}

// Start begins processing events
func (p *Processor) Start(ctx context.Context) error {
	p.logger.Info("Starting reward processor")
	return p.consumer.Start(ctx)
}

// Close closes the processor and its resources
func (p *Processor) Close() error {
	if err := p.consumer.Close(); err != nil {
		p.logger.Error("Error closing consumer", zap.Error(err))
	}
	if err := p.producer.Close(); err != nil {
		p.logger.Error("Error closing producer", zap.Error(err))
	}
	return nil
}
