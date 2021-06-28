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
		obj := "user data"
		h.log.Error(ErrEncoding(err, obj))
		http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// AnonymousStatsHandler updates anonymous stats for user
func (h *Handler) AnonymousStatsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodGet {
		if err := json.NewEncoder(w).Encode(prefObj); err != nil {
			obj := "user preference object"
			h.log.Error(ErrEncoding(err, obj))
			http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
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
			obj := "anonymousUsageStats"
			h.log.Error(ErrParseBool(err, obj))

			http.Error(w, ErrParseBool(err, obj).Error(), http.StatusBadRequest)
			return
		}
		prefObj.AnonymousUsageStats = aUsageStats
		trackStats = true
	}

	perfStats := req.FormValue("anonymousPerfResults")
	if perfStats != "" {
		aPerfStats, err := strconv.ParseBool(perfStats)
		if err != nil {
			obj := "anonymousPerfResults"
			h.log.Error(ErrParseBool(err, obj))

			http.Error(w, ErrParseBool(err, obj).Error(), http.StatusBadRequest)
			return
		}
		prefObj.AnonymousPerfResults = aPerfStats
		if err = provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
			h.log.Error(ErrRecordPreferences(err))
			http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
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
