package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the response structure for health endpoints
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler handles the general health check endpoint
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "ok",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HealthLiveHandler handles the liveness probe endpoint
func (h *Handler) HealthLiveHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "ok",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HealthReadyHandler handles the readiness probe endpoint
func (h *Handler) HealthReadyHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "ok",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
