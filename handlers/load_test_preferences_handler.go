//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// LoadTestPrefencesHandler is used for persisting load test preferences
func (h *Handler) LoadTestPrefencesHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodGet {
		if err := json.NewEncoder(w).Encode(prefObj); err != nil {
			logrus.Errorf("Error encoding LoadTest preference object: %v", err)
			http.Error(w, "Error encoding LoadTest preference object", http.StatusInternalServerError)
		}
		return
	}
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// qps, _ := strconv.ParseInt(q.Get("qps"), 32)
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
	gen := req.FormValue("loadGenerator")
	genTrack := false
	// TODO: after we have interfaces for load generators in place, we need to make a generic check, for now using a hard coded one
	for _, lg := range []models.LoadGenerator{models.FortioLG, models.Wrk2LG} {
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

	_, _ = w.Write([]byte("{}"))
}
