package main

import (
	"event-processor/internal/messaging/kafka"
	"event-processor/internal/service"
	"event-processor/internal/transport"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	defaultPort         = "8081"
	defaultKafkaBrokers = "localhost:9092"
	defaultKafkaTopic   = "platform-events"
)

func main() {
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
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	svc := service.NewEventService(producer)
	server := transport.NewServer(svc)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	fmt.Printf("\nEvent Processor running at http://localhost:%s/\n", port)
	fmt.Printf("Publishing events to Kafka topic: %s\n", kafkaTopic)

	if err := http.ListenAndServe(":"+port, server.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
