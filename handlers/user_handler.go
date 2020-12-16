package handlers

import (
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// UserHandler returns info about the logged in user
func (h *Handler) UserHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if err := json.NewEncoder(w).Encode(user); err != nil {
		logrus.Errorf("error getting user data: %v", err)
		http.Error(w, "unable to get session", http.StatusInternalServerError)
		return
	}
}

// AnonymousStatsHandler updates anonymous stats for user
func (h *Handler) AnonymousStatsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodGet {
		if err := json.NewEncoder(w).Encode(prefObj); err != nil {
			logrus.Errorf("Error encoding user preference object: %v", err)
			http.Error(w, "Error encoding user preference object", http.StatusInternalServerError)
		}
		return
	}
	// if req.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	var trackStats bool
	usageStats := req.FormValue("anonymousUsageStats")
	if usageStats != "" {
		aUsageStats, err := strconv.ParseBool(usageStats)
		if err != nil {
			err = errors.Wrap(err, "unable to parse anonymousUsageStats")
			logrus.Error(err)
			http.Error(w, "please provide a valid value for anonymousUsageStats", http.StatusBadRequest)
			return
		}
		prefObj.AnonymousUsageStats = aUsageStats
		trackStats = true
	}

	meshMapPreferences := req.FormValue("StartOnZoom")
	if meshMapPreferences != "" {
		aMeshMapPreferences, err := strconv.ParseBool(meshMapPreferences)
		if err != nil {
			err = errors.Wrap(err, "unable to parse meshMapPreferences")
			logrus.Error(err)
			http.Error(w, "please provide a valid value for meshMapPreferences", http.StatusBadRequest)
			return
		}
		prefObj.MeshMapPreferences = &models.MeshMapPreferences{
			StartOnZoom: aMeshMapPreferences,
		}
		trackStats = true
	}

	perfStats := req.FormValue("anonymousPerfResults")
	if perfStats != "" {
		aPerfStats, err := strconv.ParseBool(perfStats)
		if err != nil {
			err = errors.Wrap(err, "unable to parse anonymousPerfResults")
			logrus.Error(err)
			http.Error(w, "please provide a valid value for anonymousPerfResults", http.StatusBadRequest)
			return
		}
		prefObj.AnonymousPerfResults = aPerfStats
		if err = provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
			logrus.Errorf("unable to save user preferences: %v", err)
			http.Error(w, "unable to save user preferences", http.StatusInternalServerError)
			return
		}
		trackStats = true
	}
	if trackStats {
		if err := json.NewEncoder(w).Encode(prefObj); err != nil {
			logrus.Errorf("unable to save user preferences: %v", err)
			http.Error(w, "Unable to decode preferences", http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, "no stats update requested", http.StatusBadRequest)
}
