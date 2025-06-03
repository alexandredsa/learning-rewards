package main

import (
	"event-processor/internal/db"
	"event-processor/internal/repository"
	"event-processor/internal/service"
	"event-processor/internal/transport"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db, err := db.Connect(os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	repo := repository.NewEventRepository(db)
	svc := service.NewEventService(repo)
	server := transport.NewServer(svc)

	http.HandleFunc("/events", server.EventHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
