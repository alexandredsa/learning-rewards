package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/logger"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"go.uber.org/zap"
)

// Consumer represents a Kafka consumer for user events
type Consumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	log      *zap.Logger
	handler  func(models.UserEvent) error
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, groupID string, topics []string) (*Consumer, error) {
	log := logger.Get()

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	log.Info("Creating Kafka consumer",
		zap.Strings("brokers", brokers),
		zap.String("group_id", groupID),
		zap.Strings("topics", topics))

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Error("Failed to create Kafka consumer",
			zap.Strings("brokers", brokers),
			zap.String("group_id", groupID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	log.Info("Successfully created Kafka consumer")
	return &Consumer{
		consumer: consumer,
		topics:   topics,
		log:      log,
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
		log:     c.log,
	}

	c.log.Info("Starting Kafka consumer",
		zap.Strings("topics", c.topics))

	for {
		select {
		case <-ctx.Done():
			c.log.Info("Context cancelled, stopping consumer")
			return ctx.Err()
		default:
			c.log.Info("Attempting to join consumer group session")
			err := c.consumer.Consume(ctx, c.topics, consumer)
			if err != nil {
				if err == sarama.ErrClosedConsumerGroup {
					c.log.Error("Consumer group was closed",
						zap.Error(err))
					return fmt.Errorf("consumer group was closed: %w", err)
				}
				c.log.Error("Error during consumer group session",
					zap.Error(err),
					zap.String("error_type", fmt.Sprintf("%T", err)))
				// Add a small delay before retrying to avoid tight loop
				time.Sleep(time.Second)
				continue
			}

			if ctx.Err() != nil {
				c.log.Info("Context cancelled during consumer group session")
				return ctx.Err()
			}
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	c.log.Info("Closing Kafka consumer")
	if err := c.consumer.Close(); err != nil {
		c.log.Error("Error closing Kafka consumer",
			zap.Error(err))
		return fmt.Errorf("error closing consumer: %w", err)
	}
	c.log.Info("Successfully closed Kafka consumer")
	return nil
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler func(models.UserEvent) error
	log     *zap.Logger
}

// Setup is run at the beginning of a new session
func (h *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	h.log.Info("Consumer group handler setup started",
		zap.String("member_id", session.MemberID()),
		zap.Int32("generation_id", session.GenerationID()))

	// Log the claims (partitions) assigned to this consumer
	claims := session.Claims()
	for topic, partitions := range claims {
		h.log.Info("Consumer assigned partitions",
			zap.String("topic", topic),
			zap.Int32s("partitions", partitions))
	}

	h.log.Info("Consumer group handler setup completed")
	return nil
}

// Cleanup is run at the end of a session
func (h *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.log.Info("Consumer group handler cleanup started",
		zap.String("member_id", session.MemberID()),
		zap.Int32("generation_id", session.GenerationID()))

	// Log the claims that were being processed
	claims := session.Claims()
	for topic, partitions := range claims {
		h.log.Info("Consumer cleaning up partitions",
			zap.String("topic", topic),
			zap.Int32s("partitions", partitions))
	}

	h.log.Info("Consumer group handler cleanup completed")
	return nil
}

// ConsumeClaim processes messages from a claim
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	h.log.Info("Starting to consume messages",
		zap.String("topic", claim.Topic()),
		zap.Int32("partition", claim.Partition()),
		zap.Int64("initial_offset", claim.InitialOffset()),
		zap.Int64("high_water_mark", claim.HighWaterMarkOffset()))

	messageCount := 0
	lastOffset := claim.InitialOffset()

	for message := range claim.Messages() {
		messageCount++
		lastOffset = message.Offset

		if messageCount%100 == 0 {
			h.log.Info("Message consumption progress",
				zap.String("topic", message.Topic),
				zap.Int32("partition", message.Partition),
				zap.Int64("current_offset", message.Offset),
				zap.Int64("high_water_mark", claim.HighWaterMarkOffset()),
				zap.Int("messages_processed", messageCount))
		}

		h.log.Debug("Processing message",
			zap.String("topic", message.Topic),
			zap.Int32("partition", message.Partition),
			zap.Int64("offset", message.Offset),
			zap.Time("timestamp", message.Timestamp),
			zap.Int("message_size", len(message.Value)))

		var event models.UserEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			h.log.Error("Failed to unmarshal event", zap.Error(err))
			continue
		}

		if err := h.handler(event); err != nil {
			h.log.Error("Failed to process event",
				zap.Error(err),
				zap.Any("event", event))
			continue
		}

		session.MarkMessage(message, "")
		h.log.Debug("Successfully processed and marked message",
			zap.String("topic", message.Topic),
			zap.Int32("partition", message.Partition),
			zap.Int64("offset", message.Offset))
	}

	h.log.Info("Finished consuming messages",
		zap.String("topic", claim.Topic()),
		zap.Int32("partition", claim.Partition()),
		zap.Int64("last_offset", lastOffset),
		zap.Int64("high_water_mark", claim.HighWaterMarkOffset()),
		zap.Int("total_messages_processed", messageCount))
	return nil
}
