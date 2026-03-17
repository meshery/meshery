package handlers

import (
	"encoding/json"
	"net/http"
)

type AIRequest struct {
	Prompt string `json:"prompt"`
}

type AIResponse struct {
	Response string `json:"response"`
}

// NEW: AI Config struct
type AIConfig struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
}

// NEW: Store config (in-memory)
var aiConfig AIConfig

// Existing test handler
func HandleAITest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req AIRequest

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Create response
	resp := AIResponse{
		Response: "Mock response for: " + req.Prompt,
	}

	// Set header
	w.Header().Set("Content-Type", "application/json")

	// Encode response
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// NEW: AI Config handler
func HandleAIConfig(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodPost:
		var config AIConfig

		// Decode request
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Save config
		aiConfig = config

		// Return saved config
		if err := json.NewEncoder(w).Encode(aiConfig); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodGet:
		// Return current config
		if err := json.NewEncoder(w).Encode(aiConfig); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}