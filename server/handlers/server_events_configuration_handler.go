package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/meshery/meshery/server/models"
	"github.com/sirupsen/logrus"
)

type LogLevelResponse struct {
	// Status of the operation (success/error)
	// example: success/error
	Status string `json:"status,omitempty"`

	// Current log level of the server
	// example: info
	// required: true
	LogLevel string `json:"event_log_level"`

	// List of available logging levels
	// example: ["panic","fatal","error","warn","info","debug","trace"]
	Available []string `json:"available_levels,omitempty"`
}

type LogLevelRequest struct {
	// Desired log level to set
	// required: true
	// enum: panic,fatal,error,warn,info,debug,trace
	// example: debug
	LogLevel string `json:"eventLogLevel"`
}

// UnmarshalJSON dual-accepts the canonical `eventLogLevel` wire key and the
// legacy `event_log_level` spelling during the Phase 2 deprecation window.
// Go's encoding/json case-insensitive tag fallback does NOT cross an
// underscore boundary, so a struct tagged `eventLogLevel` would silently
// drop a payload keyed `event_log_level`. Canonical wins when both are
// present. Retire the fallback once Phase 3 consumer migration completes.
func (r *LogLevelRequest) UnmarshalJSON(data []byte) error {
	aux := struct {
		Canonical *string `json:"eventLogLevel,omitempty"`
		Legacy    *string `json:"event_log_level,omitempty"`
	}{}
	r.LogLevel = ""
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	switch {
	case aux.Canonical != nil:
		r.LogLevel = *aux.Canonical
	case aux.Legacy != nil:
		r.LogLevel = *aux.Legacy
	}
	return nil
}

// getAvailableLogLevels returns all valid logging levels supported by the system
func getAvailableLogLevels() []string {
	levels := make([]string, len(logrus.AllLevels))
	for i, level := range logrus.AllLevels {
		levels[i] = level.String()
	}
	return levels
}

func (h *Handler) ServerEventConfigurationHandler(w http.ResponseWriter, req *http.Request,
	prefObj *models.Preference, user *models.User, provider models.Provider) {

	switch req.Method {
	case http.MethodPut:
		h.ServerEventConfigurationSet(w, req, prefObj, user, provider)
	case http.MethodGet:
		h.ServerEventConfigurationGet(w, req, prefObj, user, provider)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) ServerEventConfigurationSet(w http.ResponseWriter, req *http.Request,
	prefObj *models.Preference, user *models.User, provider models.Provider) {

	var logLevelReq LogLevelRequest
	if err := json.NewDecoder(req.Body).Decode(&logLevelReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(LogLevelResponse{
			Status:    "error",
			LogLevel:  h.log.GetLevel().String(),
			Available: getAvailableLogLevels(),
		}); err != nil {
			h.log.Error(err)
		}
		return
	}

	requestedLevel := strings.ToLower(strings.TrimSpace(logLevelReq.LogLevel))
	level, err := logrus.ParseLevel(requestedLevel)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(LogLevelResponse{
			Status:    "error",
			LogLevel:  h.log.GetLevel().String(),
			Available: getAvailableLogLevels(),
		}); err != nil {
			h.log.Error(err)
		}
		return
	}

	h.log.SetLevel(level)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(LogLevelResponse{
		Status:    "success",
		LogLevel:  level.String(),
		Available: getAvailableLogLevels(),
	}); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) ServerEventConfigurationGet(w http.ResponseWriter, req *http.Request,
	prefObj *models.Preference, user *models.User, provider models.Provider) {

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(LogLevelResponse{
		LogLevel:  h.log.GetLevel().String(),
		Available: getAvailableLogLevels(),
	}); err != nil {
		h.log.Error(err)
	}
}
