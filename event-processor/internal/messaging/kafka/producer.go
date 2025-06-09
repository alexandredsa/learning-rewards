package kafka

import (
	"context"
	"encoding/json"
	"event-processor/internal/logger"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	maxRetries = 3
	retryDelay = time.Second * 2
)

// Producer handles publishing events to Kafka
type Producer struct {
	producer sarama.SyncProducer
	topic    string
	logger   *zap.Logger
	config   *sarama.Config
	brokers  []string
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
	config.Producer.Retry.Max = maxRetries
	config.Producer.Retry.Backoff = retryDelay
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	// Try to connect with retries
	var producer sarama.SyncProducer
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Info("attempting to connect to Kafka brokers",
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", maxRetries),
		)

		producer, err = sarama.NewSyncProducer(brokers, config)
		if err == nil {
			break
		}

		log.Warn("failed to connect to Kafka brokers",
			zap.Error(err),
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", maxRetries),
		)

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		log.Error("failed to create Kafka producer after all retries", zap.Error(err))
		return nil, fmt.Errorf("failed to create Kafka producer after %d attempts: %w", maxRetries, err)
	}

	log.Info("successfully connected to Kafka")
	return &Producer{
		producer: producer,
		topic:    topic,
		logger:   log,
		config:   config,
		brokers:  brokers,
	}, nil
}

// ensureConnection checks if the producer is still connected and reconnects if necessary
func (p *Producer) ensureConnection() error {
	// If producer is nil, we need to create a new one
	if p.producer == nil {
		return p.reconnect()
	}
	return nil
}

// reconnect attempts to establish a new connection to Kafka
func (p *Producer) reconnect() error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		p.logger.Info("attempting to connect to Kafka",
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", maxRetries),
		)

		p.producer, err = sarama.NewSyncProducer(p.brokers, p.config)
		if err == nil {
			p.logger.Info("successfully connected to Kafka")
			return nil
		}

		p.logger.Warn("failed to connect to Kafka",
			zap.Error(err),
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", maxRetries),
		)

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return fmt.Errorf("failed to connect to Kafka after %d attempts: %w", maxRetries, err)
}

// PublishEvent publishes an event to Kafka
func (p *Producer) PublishEvent(ctx context.Context, event interface{}) error {
	// Ensure we have a connection before trying to publish
	if err := p.ensureConnection(); err != nil {
		return err
	}

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

	// Send message to Kafka with retries
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		partition, offset, err := p.producer.SendMessage(msg)
		if err == nil {
			p.logger.Info("event published successfully",
				zap.Int32("partition", partition),
				zap.Int64("offset", offset),
				zap.String("topic", p.topic),
			)
			return nil
		}

		lastErr = err
		p.logger.Warn("failed to publish event, will retry",
			zap.Error(err),
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", maxRetries),
		)

		// If we get a connection error, try to reconnect
		if err == sarama.ErrNotConnected || err == sarama.ErrClosedClient {
			if reconnectErr := p.reconnect(); reconnectErr != nil {
				return reconnectErr
			}
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	p.logger.Error("failed to publish event after all retries",
		zap.Error(lastErr),
		zap.String("topic", p.topic),
	)
	return fmt.Errorf("failed to publish event after %d attempts: %w", maxRetries, lastErr)
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
