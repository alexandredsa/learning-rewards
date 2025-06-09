package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
)

// Producer handles publishing events to Kafka
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

// PublishEvent publishes an event to Kafka
func (p *Producer) PublishEvent(ctx context.Context, event interface{}) error {
	// Convert event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Create Kafka message
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(eventJSON),
	}

	// Send message to Kafka
	_, _, err = p.producer.SendMessage(msg)
	return err
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	return p.producer.Close()
}
