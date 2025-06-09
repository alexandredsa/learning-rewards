package main

import (
	"event-processor/internal/logger"
	"event-processor/internal/messaging/kafka"
	"event-processor/internal/service"
	"event-processor/internal/transport"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
)

const (
	defaultPort         = "8081"
	defaultKafkaBrokers = "localhost:29092"
	defaultKafkaTopic   = "learning-events"
)

func main() {
	// Initialize logger
	if err := logger.Init(os.Getenv("DEBUG") != ""); err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	log := logger.Get()
	log.Info("starting event processor service")

	// Get Kafka configuration from environment variables
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = defaultKafkaBrokers
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = defaultKafkaTopic
	}

	// Initialize Kafka producer
	producer, err := kafka.NewProducer(strings.Split(kafkaBrokers, ","), kafkaTopic)
	if err != nil {
		log.Fatal("failed to create Kafka producer", zap.Error(err))
	}
	defer producer.Close()

	svc := service.NewEventService(producer)
	server := transport.NewServer(svc)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Info("event processor service started",
		zap.String("port", port),
		zap.String("kafka_topic", kafkaTopic),
	)

	if err := http.ListenAndServe(":"+port, server.Router()); err != nil {
		log.Fatal("server error", zap.Error(err))
	}
}
