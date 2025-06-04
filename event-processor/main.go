package main

import (
	"event-processor/internal/db"
	"event-processor/internal/repository"
	"event-processor/internal/service"
	"event-processor/internal/transport"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const defaultPort = "8081"

func main() {
	db, err := db.Connect(os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	repo := repository.NewEventRepository(db)
	svc := service.NewEventService(repo)
	server := transport.NewServer(svc)
	http.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		server.EventHandler(w, r)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	fmt.Printf("\nEvent Processor running at http://localhost:%s/\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
