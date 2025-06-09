package kafka

import (
	"context"
	"encoding/json"
	"event-processor/internal/logger"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// Producer handles publishing events to Kafka
type Producer struct {
	producer sarama.SyncProducer
	topic    string
	logger   *zap.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	log := logger.Get().With(
		zap.String("component", "kafka_producer"),
		zap.Strings("brokers", brokers),
		zap.String("topic", topic),
	)

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	log.Info("connecting to Kafka brokers")
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Error("failed to create Kafka producer", zap.Error(err))
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	log.Info("successfully connected to Kafka")
	return &Producer{
		producer: producer,
		topic:    topic,
		logger:   log,
	}, nil
}

// PublishEvent publishes an event to Kafka
func (p *Producer) PublishEvent(ctx context.Context, event interface{}) error {
	// Convert event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("failed to marshal event to JSON", zap.Error(err))
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create Kafka message
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(eventJSON),
	}

	p.logger.Debug("publishing event to Kafka",
		zap.String("topic", p.topic),
		zap.ByteString("event", eventJSON),
	)

	// Send message to Kafka
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.Error("failed to publish event to Kafka",
			zap.Error(err),
			zap.String("topic", p.topic),
		)
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Info("event published successfully",
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.String("topic", p.topic),
	)
	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	p.logger.Info("closing Kafka producer")
	if err := p.producer.Close(); err != nil {
		p.logger.Error("error closing Kafka producer", zap.Error(err))
		return err
	}
	p.logger.Info("Kafka producer closed successfully")
	return nil
}
