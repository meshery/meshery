package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/layer5io/meshery/server/models"
	"github.com/sirupsen/logrus"
)

// LogLevelResponse represents the response structure for log level operations
// swagger:response LogLevelResponse
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

// LogLevelRequest represents the request body for setting log level
// swagger:parameters LogLevelRequest
type LogLevelRequest struct {
	// Desired log level to set
	// required: true
	// enum: panic,fatal,error,warn,info,debug,trace
	// example: debug
	LogLevel string `json:"event_log_level"`
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
	case http.MethodPost:
		h.ServerEventConfigurationSet(w, req, prefObj, user, provider)
	case http.MethodGet:
		h.ServerEventConfigurationGet(w, req, prefObj, user, provider)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// swagger:route POST /api/system/events system setLogLevel
//
// # Set server event logging level
//
// Updates the server's logging level. The new level must be one of the available logging levels.
// Invalid log levels will result in a 400 Bad Request response.
//
// Responses:
//
//	200: LogLevelResponse       Successfully updated log level
//	400: LogLevelResponse       Invalid log level provided
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

// swagger:route GET /api/system/events system getLogLevel
//
// # Get current logging configuration for the server events
//
// Retrieves the current server log level and lists all available log levels.
//
// Produces:
//   - application/json
//
// Responses:
//
//	200: LogLevelResponse       Successfully retrieved log level
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
