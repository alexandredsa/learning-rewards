package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"go.uber.org/zap"
)

// Producer represents a Kafka producer for reward events
type Producer struct {
	producer sarama.SyncProducer
	topic    string
	logger   *zap.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string, logger *zap.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}, nil
}

// SendReward sends a reward event to Kafka
func (p *Producer) SendReward(reward models.RewardTriggered) error {
	value, err := json.Marshal(reward)
	if err != nil {
		return fmt.Errorf("failed to marshal reward: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	p.logger.Info("Reward event sent",
		zap.String("topic", p.topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.String("user_id", reward.UserID),
		zap.String("rule_id", reward.RuleID))

	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.producer.Close()
}
