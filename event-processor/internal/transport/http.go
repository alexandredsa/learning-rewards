package transport

import (
	"encoding/json"
	"event-processor/internal/service"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	svc    service.EventService
	router *mux.Router
}

func NewServer(svc service.EventService) *Server {
	s := &Server{
		svc:    svc,
		router: mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/events", s.handleEvent).Methods(http.MethodPost)
	s.router.HandleFunc("/health", s.handleHealth).Methods(http.MethodGet)
}

func (s *Server) Router() *mux.Router {
	return s.router
}

type EventRequest struct {
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	CourseID  string    `json:"course_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *Server) handleEvent(w http.ResponseWriter, r *http.Request) {
	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := s.svc.ProcessEvent(ctx, req.UserID, req.EventType, req.CourseID, req.Timestamp)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"event accepted"}`))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
