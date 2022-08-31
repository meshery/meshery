// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/layer5io/meshery/server/models"
	SMP "github.com/layer5io/service-mesh-performance/spec"
)

// UserTestPreferenceHandler is used for persisting load test preferences
func (h *Handler) UserTestPreferenceHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost && req.Method != http.MethodGet && req.Method != http.MethodDelete {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if req.Method == http.MethodPost {
		h.UserTestPreferenceStore(w, req, prefObj, user, provider)
		return
	}
	if req.Method == http.MethodDelete {
		h.UserTestPreferenceDelete(w, req, prefObj, user, provider)
		return
	}
	if req.Method == http.MethodGet {
		h.UserTestPreferenceGet(w, req, prefObj, user, provider)
		return
	}
}

// swagger:route POST /api/user/prefs/perf UserAPI idPostLoadPreferences
// Handle POST request for load test preferences
//
// Used for persisting load test preferences
// responses:
// 	200:

// UserTestPreferenceStore is used for persisting load test preferences
func (h *Handler) UserTestPreferenceStore(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		// logrus.Error(err)
		// http.Error(w, msg, http.StatusInternalServerError)
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}
	perfTest := &SMP.PerformanceTestConfig{}
	if err = protojson.Unmarshal(body, perfTest); err != nil {
		obj := "the provided input"
		// err = errors.Wrapf(err, msg)
		// logrus.Error(err)
		// http.Error(w, msg, http.StatusBadRequest)
		h.log.Error(ErrUnmarshal(err, obj))
		http.Error(w, ErrUnmarshal(err, obj).Error(), http.StatusBadRequest)
		return
	}
	if err = models.SMPPerformanceTestConfigValidator(perfTest); err != nil {
		// logrus.Error(err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Error(ErrRecordPreferences(err))
		http.Error(w, ErrRecordPreferences(err).Error(), http.StatusBadRequest)
		return
	}
	tid, err := provider.SMPTestConfigStore(req, perfTest)
	if err != nil {
		obj := "user preference"
		// logrus.Errorf("unable to save user preferences: %v", err)
		// http.Error(w, "unable to save user preferences", http.StatusInternalServerError)
		h.log.Error(ErrFailToSave(err, obj))
		http.Error(w, ErrFailToSave(err, obj).Error(), http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte(tid))
}

// swagger:route GET /api/user/prefs/perf UserAPI idGetLoadPreferences
// Handle GET request for load test preferences
//
// Used for fetching load test preferences
// responses:
// 	200: loadTestPreferencesWrapper

// UserTestPreferenceGet gets the PerformanceTestConfig object
func (h *Handler) UserTestPreferenceGet(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()
	testUUID := q.Get("uuid")
	h.log.Debug(testUUID)
	if testUUID == "" {
		testPage := q.Get("page")
		testPageSize := q.Get("pageSize")
		testSearch := q.Get("search")
		testOrder := q.Get("order")
		h.log.Debug(testPage, testPageSize)
		testObjJSON, err := provider.SMPTestConfigFetch(req, testPage, testPageSize, testSearch, testOrder)
		if err != nil {
			// logrus.Error("error fetching test configs")
			// http.Error(w, "error fetching test configs", http.StatusInternalServerError)
			h.log.Error(ErrTestConfigs)
			http.Error(w, ErrTestConfigs.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(testObjJSON)
	} else {
		testObj, err := provider.SMPTestConfigGet(req, testUUID)
		if err != nil {
			// logrus.Error("error fetching test configs")
			// http.Error(w, "error fetching test configs", http.StatusInternalServerError)
			h.log.Error(ErrTestConfigs)
			http.Error(w, ErrTestConfigs.Error(), http.StatusInternalServerError)
			return
		}
		if testObj == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Printf("%v", testObj)
		data, err := protojson.Marshal(testObj)
		if err != nil {
			// logrus.Errorf("error reading database: %v", err)
			// http.Error(w, "error reading database", http.StatusInternalServerError)
			h.log.Error(ErrReadConfig(err))
			http.Error(w, ErrReadConfig(err).Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(data)
		if err != nil {
			// logrus.Errorf("error writing response: %v", err)
			// http.Error(w, "error writing response", http.StatusInternalServerError)
			h.log.Error(ErrWriteResponse)
			http.Error(w, ErrWriteResponse.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// swagger:route DELETE /api/user/prefs/perf UserAPI idDeleteLoadPreferences
// Handle DELETE request for load test preferences
//
// Used for deleting load test preferences
// responses:
// 	200:

// UserTestPreferenceDelete deletes the PerformanceTestConfig object
func (h *Handler) UserTestPreferenceDelete(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	testUUID := req.URL.Query().Get("uuid")
	if testUUID == "" {
		obj := "field uuid"
		// logrus.Error("field uuid not found")
		// http.Error(w, "field uuid not found", http.StatusBadRequest)
		h.log.Error(ErrQueryGet(obj))
		http.Error(w, ErrQueryGet(obj).Error(), http.StatusBadRequest)
		return
	}
	if err := provider.SMPTestConfigDelete(req, testUUID); err != nil {
		obj := "testConfig"
		// logrus.Errorf("error deleting testConfig: %v", err)
		// http.Error(w, "error deleting testConfig", http.StatusBadRequest)
		h.log.Error(ErrFailToDelete(err, obj))
		http.Error(w, ErrFailToDelete(err, obj).Error(), http.StatusBadRequest)
		return
	}
}
