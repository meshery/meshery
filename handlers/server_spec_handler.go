package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Version defines the Json payload structure for version api\
type Version struct {
	Build     string `json:"build,omitempty"`
	CommitSHA string `json:"commitsha,omitempty"`
}

// ServerVersionHandler handles the version api request for the server
func (h *Handler) ServerVersionHandler(w http.ResponseWriter, r *http.Request) {
	// Default values incase any errors
	version := &Version{
		Build:     viper.GetString("BUILD"),
		CommitSHA: viper.GetString("COMMITSHA"),
	}

	// if r.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(version)
	if err != nil {
		logrus.Errorf("unable to send data: %v", err)
	}
}
