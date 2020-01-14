package handlers

import (
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// UserHandler returns info about the logged in user
func (h *Handler) UserHandler(w http.ResponseWriter, req *http.Request, _ *sessions.Session, _ *models.Preference, user *models.User) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		logrus.Errorf("error getting user data: %v", err)
		http.Error(w, "unable to get session", http.StatusInternalServerError)
		return
	}
}

// AnonymousStatsHandler updates anonymous stats for user
func (h *Handler) AnonymousStatsHandler(w http.ResponseWriter, req *http.Request, _ *sessions.Session, prefObj *models.Preference, user *models.User) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	stats := req.FormValue("anonymousStats")
	aStats, err := strconv.ParseBool(stats)
	if err != nil {
		err = errors.Wrap(err, "unable to parse anonymous_stats")
		logrus.Error(err)
		http.Error(w, "please provide a valid value for anonymous_stats", http.StatusBadRequest)
		return
	}
	prefObj.AnonymousStats = aStats
	if err = h.config.Provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
		logrus.Errorf("unable to save user preferences: %v", err)
		http.Error(w, "unable to save user preferences", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{}"))
}
