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
	var req AIRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	resp := AIResponse{
		Response: "Mock response for: " + req.Prompt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}