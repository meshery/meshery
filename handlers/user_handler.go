package handlers

import (
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"github.com/pkg/errors"

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

// swagger:route GET /api/user/prefs UserAPI idGetUserTestPrefs
// Handle GET for User Load Test Preferences
//
// Returns User Load Test Preferences
// responses:
// 	200: userLoadTestPrefsRespWrapper

// swagger:route POST /api/user/prefs UserAPI idPostUserTestPrefs
// Handle GET for User Load Test Preferences
//
// Updates User Load Test Preferences
// responses:
// 	200: userLoadTestPrefsRespWrapper

// UserPrefsHandler updates anonymous stats for user or for persisting load test preferences
func (h *Handler) UserPrefsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodGet {
		if err := json.NewEncoder(w).Encode(prefObj); err != nil {
			obj := "user preference object"
			h.log.Error(ErrEncoding(err, obj))
			http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		}
		return
	}

	defer func() {
		_ = req.Body.Close()
	}()

	// read user preferences from JSON request body
	if err := json.NewDecoder(req.Body).Decode(&prefObj); err != nil {
		h.log.Error(ErrDecoding(err, "user preferences"))
		http.Error(w, ErrDecoding(err, "user preferences").Error(), http.StatusInternalServerError)
		return
	}

	// validate load test data
	qps := prefObj.LoadTestPreferences.QueriesPerSecond
	if qps < 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := fmt.Errorf("QPS value less than 0")
		h.log.Error(ErrSavingUserPreference(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dur := prefObj.LoadTestPreferences.Duration
	if _, err := time.ParseDuration(dur); err != nil {
		err = errors.Wrap(err, "unable to parse test duration")
		h.log.Error(ErrSavingUserPreference(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := prefObj.LoadTestPreferences.ConcurrentRequests
	if c < 0 {
		err := fmt.Errorf("number of concurrent requests less than 0")
		h.log.Error(ErrSavingUserPreference(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loadGen := prefObj.LoadTestPreferences.LoadGenerator
	loadGenSupoorted := false
	for _, lg := range []models.LoadGenerator{models.FortioLG, models.Wrk2LG, models.NighthawkLG} {
		if lg.Name() == loadGen {
			loadGenSupoorted = true
		}
	}
	if !loadGenSupoorted {
		err := fmt.Errorf("invalid load generator: %s", loadGen)
		h.log.Error(ErrSavingUserPreference(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prefObj.AnonymousUsageStats = true

	if err := provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
		err := fmt.Errorf("unable to save user preferences: %v", err)
		h.log.Error(ErrSavingUserPreference(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(prefObj); err != nil {
		obj := "user preferences"
		h.log.Error(ErrEncoding(err, obj))
		http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}
