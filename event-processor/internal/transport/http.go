package transport

import (
	"encoding/json"
	"event-processor/internal/service"
	"net/http"
	"time"
)

type Server struct {
	svc service.EventService
}

func NewServer(svc service.EventService) *Server {
	return &Server{svc: svc}
}

type EventRequest struct {
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	CourseID  string    `json:"course_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *Server) EventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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
