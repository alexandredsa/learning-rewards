package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/logger"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"go.uber.org/zap"
)

// Producer represents a Kafka producer for reward events
type Producer struct {
	producer sarama.SyncProducer
	topic    string
	log      *zap.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	log := logger.Get()

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	log.Info("Creating Kafka producer",
		zap.Strings("brokers", brokers),
		zap.String("topic", topic))

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Error("Failed to create Kafka producer",
			zap.Strings("brokers", brokers),
			zap.String("topic", topic),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	log.Info("Successfully created Kafka producer")
	return &Producer{
		producer: producer,
		topic:    topic,
		log:      log,
	}, nil
}

// SendReward sends a reward event to Kafka
func (p *Producer) SendReward(reward models.RewardTriggered) error {
	value, err := json.Marshal(reward)
	if err != nil {
		p.log.Error("Failed to marshal reward",
			zap.Error(err),
			zap.Any("reward", reward))
		return fmt.Errorf("failed to marshal reward: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(value),
	}

	p.log.Debug("Sending reward message",
		zap.String("topic", p.topic),
		zap.Any("reward", reward))

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.log.Error("Failed to send reward message",
			zap.String("topic", p.topic),
			zap.Any("reward", reward),
			zap.Error(err))
		return fmt.Errorf("failed to send message: %w", err)
	}

	p.log.Debug("Successfully sent reward message",
		zap.String("topic", p.topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.String("user_id", reward.UserID),
		zap.String("rule_id", reward.RuleID))

	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	p.log.Info("Closing Kafka producer")
	if err := p.producer.Close(); err != nil {
		p.log.Error("Error closing Kafka producer",
			zap.Error(err))
		return fmt.Errorf("error closing producer: %w", err)
	}
	p.log.Info("Successfully closed Kafka producer")
	return nil
}
