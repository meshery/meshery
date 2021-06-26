package handlers

import (
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/layer5io/meshery/models"
)

// UserHandler returns info about the logged in user
func (h *Handler) UserHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.log.Error(ErrUserData(err))
		http.Error(w, ErrUserData(err).Error(), http.StatusInternalServerError)
		return
	}
}

// AnonymousStatsHandler updates anonymous stats for user
func (h *Handler) AnonymousStatsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodGet {
		if err := json.NewEncoder(w).Encode(prefObj); err != nil {
			h.log.Error(ErrUserPreferenceObject(err))
			http.Error(w, ErrUserPreferenceObject(err).Error(), http.StatusInternalServerError)
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
			h.log.Error(err)

			http.Error(w, "please provide a valid value for anonymousUsageStats", http.StatusBadRequest)
			return
		}
		prefObj.AnonymousUsageStats = aUsageStats
		trackStats = true
	}

	perfStats := req.FormValue("anonymousPerfResults")
	if perfStats != "" {
		aPerfStats, err := strconv.ParseBool(perfStats)
		if err != nil {
			h.log.Error(err)

			http.Error(w, "please provide a valid value for anonymousPerfResults", http.StatusBadRequest)
			return
		}
		prefObj.AnonymousPerfResults = aPerfStats
		if err = provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
			h.log.Error(ErrSaveUserPreference(err))
			http.Error(w, ErrSaveUserPreference(err).Error(), http.StatusInternalServerError)
			return
		}
		trackStats = true
	}
	if trackStats {
		if err := json.NewEncoder(w).Encode(prefObj); err != nil {
			h.log.Error(ErrEncodeUserPreference(err))
			http.Error(w, ErrEncodeUserPreference(err).Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, "no stats update requested", http.StatusBadRequest)
}
