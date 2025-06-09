package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"go.uber.org/zap"
)

// Consumer represents a Kafka consumer for user events
type Consumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	logger   *zap.Logger
	handler  func(models.UserEvent) error
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, groupID string, topics []string, logger *zap.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		topics:   topics,
		logger:   logger,
	}, nil
}

// SetHandler sets the event handler function
func (c *Consumer) SetHandler(handler func(models.UserEvent) error) {
	c.handler = handler
}

// Start begins consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	if c.handler == nil {
		return fmt.Errorf("handler not set")
	}

	consumer := &consumerGroupHandler{
		handler: c.handler,
		logger:  c.logger,
	}

	for {
		err := c.consumer.Consume(ctx, c.topics, consumer)
		if err != nil {
			return fmt.Errorf("error from consumer: %w", err)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.consumer.Close()
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler func(models.UserEvent) error
	logger  *zap.Logger
}

// Setup is run at the beginning of a new session
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a claim
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event models.UserEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			h.logger.Error("Failed to unmarshal event", zap.Error(err))
			continue
		}

		if err := h.handler(event); err != nil {
			h.logger.Error("Failed to process event",
				zap.Error(err),
				zap.Any("event", event))
			continue
		}

		session.MarkMessage(message, "")
	}

	return nil
}
