package handlers

import (
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to read request body: %s", err)
		return
	}

	// validate load test data
	qps := prefObj.LoadTestPreferences.QueriesPerSecond
	if qps < 0 {
		http.Error(w, "please provide a valid value for qps", http.StatusBadRequest)
		return
	}

	dur := prefObj.LoadTestPreferences.Duration
	if _, err := time.ParseDuration(dur); err != nil {
		err = errors.Wrap(err, "unable to parse test duration")
		logrus.Error(err)
		http.Error(w, "please provide a valid value for test duration", http.StatusBadRequest)
		return
	}

	c := prefObj.LoadTestPreferences.ConcurrentRequests
	if c < 0 {
		http.Error(w, "please provide a valid value for concurrent requests", http.StatusBadRequest)
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
		logrus.Error("invalid value for load generator")
		http.Error(w, "please specify a valid load generator", http.StatusBadRequest)
		return
	}

	prefObj.AnonymousUsageStats = true

	if err := provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
		logrus.Errorf("unable to save user preferences: %v", err)
		http.Error(w, "unable to save user preferences", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(prefObj); err != nil {
		obj := "user preferences"
		h.log.Error(ErrEncoding(err, obj))
		http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}
