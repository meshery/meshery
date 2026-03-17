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

	// Encode response with error handling
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}