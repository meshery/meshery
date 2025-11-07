package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// HealthResponse represents the response structure for health endpoints
type HealthResponse struct {
	Status string `json:"status"`
}

// respondHealthOK is a helper function to send a standardized health check response
func respondHealthOK(w http.ResponseWriter) {
	response := HealthResponse{
		Status: "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error but don't change response status since headers are already written
		logrus.Error("Failed to encode health check response: ", err)
	}
}

// HealthHandler handles the general health check endpoint
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	respondHealthOK(w)
}

// HealthLiveHandler handles the liveness probe endpoint
func (h *Handler) HealthLiveHandler(w http.ResponseWriter, r *http.Request) {
	respondHealthOK(w)
}

// HealthReadyHandler handles the readiness probe endpoint
func (h *Handler) HealthReadyHandler(w http.ResponseWriter, r *http.Request) {
	respondHealthOK(w)
}
