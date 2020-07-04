package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

// Version defines the Json payload structure for version api
type Version struct {
	Build     string `json:"build,omitempty"`
	CommitSHA string `json:"commitsha,omitempty"`
}

// ServerVersionHandler handles the version api request for the server
func (h *Handler) ServerVersionHandler(w http.ResponseWriter, r *http.Request) {

	// Default values incase any errors
	version := &Version{
		Build:     "Unavailable",
		CommitSHA: "Unavailable",
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(os.Getenv("GIT_VERSION")) > 0 {
		version.Build = os.Getenv("GIT_VERSION")
	}

	if len(os.Getenv("GIT_COMMITSHA")) > 0 {
		version.CommitSHA = os.Getenv("GIT_COMMITSHA")
	}

	json.NewEncoder(w).Encode(version)
}
