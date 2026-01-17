package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// handleAPI handles both GET and POST requests to /api
func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Enable CORS for development
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		s.handleGetState(w, r)
	case "POST":
		s.handlePostState(w, r)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// handleGetState returns the current state as JSON
func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := BuildState(s.ledStrip, s.ledBar, s.videoLight1, s.videoLight2)
	if err != nil {
		log.Printf("Error building state: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Failed to read state: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(state); err != nil {
		log.Printf("Error encoding state: %v", err)
		http.Error(w, `{"error":"Failed to encode state"}`, http.StatusInternalServerError)
		return
	}
}

// handlePostState applies the provided state
func (s *Server) handlePostState(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, `{"error":"Failed to read request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var state State
	if err := json.Unmarshal(body, &state); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Invalid JSON: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Validate state
	if err := state.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Validation failed: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Apply state to drivers
	if err := ApplyState(&state, s.ledStrip, s.ledBar, s.videoLight1, s.videoLight2); err != nil {
		log.Printf("Error applying state: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Failed to apply state: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the updated state
	updatedState, err := BuildState(s.ledStrip, s.ledBar, s.videoLight1, s.videoLight2)
	if err != nil {
		log.Printf("Error building updated state: %v", err)
		http.Error(w, `{"error":"State applied but failed to read back"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(updatedState); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	log.Println("Web: State updated successfully")
}

// handleHealth returns a simple health check response
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
