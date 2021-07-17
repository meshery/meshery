package handlers

import (
	"net/http"
	"strconv"
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

// swagger:route GET /api/user/prefs UserAPI idGetAnonymousStats
// Handle GET for anonymous stats or Load Test Preferences
//
// Returns anonymous stats for user and Load Test Preferences
// responses:
// 	200: anonymousStatsResponseWrapper

// swagger:route POST /api/user/prefs UserAPI idPostAnonymousStats
// Handle GET for anonymous stats or Load Test Preferences
//
// Updates anonymous stats for user or Update load test preferences
// responses:
// 	200: anonymousStatsResponseWrapper

// AnonymousStatsHandler updates anonymous stats for user or for persisting load test preferences
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

	// if api has been used to update AnonymousStats we return here
	if trackStats {
		if err := json.NewEncoder(w).Encode(prefObj); err != nil {
			obj := "user preferences"
			h.log.Error(ErrEncoding(err, obj))
			http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// or we update Load Test preferences
	qs := req.FormValue("qps")
	qps, err := strconv.Atoi(qs)
	if err != nil {
		err = errors.Wrap(err, "unable to parse qps")
		logrus.Error(err)
		http.Error(w, "please provide a valid value for qps", http.StatusBadRequest)
		return
	}
	if qps < 0 {
		http.Error(w, "please provide a valid value for qps", http.StatusBadRequest)
		return
	}
	dur := req.FormValue("t")
	if _, err = time.ParseDuration(dur); err != nil {
		err = errors.Wrap(err, "unable to parse t as a duration")
		logrus.Error(err)
		http.Error(w, "please provide a valid value for t", http.StatusBadRequest)
		return
	}
	cu := req.FormValue("c")
	c, err := strconv.Atoi(cu)
	if err != nil {
		err = errors.Wrap(err, "unable to parse c")
		logrus.Error(err)
		http.Error(w, "please provide a valid value for c", http.StatusBadRequest)
		return
	}
	if c < 0 {
		http.Error(w, "please provide a valid value for c", http.StatusBadRequest)
		return
	}
	gen := req.FormValue("gen")
	genTrack := false
	// TODO: after we have interfaces for load generators in place, we need to make a generic check, for now using a hard coded one
	for _, lg := range []models.LoadGenerator{models.FortioLG, models.Wrk2LG, models.NighthawkLG} {
		if lg.Name() == gen {
			genTrack = true
		}
	}
	if !genTrack {
		logrus.Error("invalid value for gen")
		http.Error(w, "please provide a valid value for gen (load generator)", http.StatusBadRequest)
		return
	}
	prefObj.LoadTestPreferences = &models.LoadTestPreferences{
		ConcurrentRequests: c,
		Duration:           dur,
		QueriesPerSecond:   qps,
		LoadGenerator:      gen,
	}
	if err = provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
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
