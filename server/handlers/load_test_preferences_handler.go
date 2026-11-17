// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/meshery/meshery/server/models"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

// maxPerformancePayloadSize defines the maximum permissible payload size for performance preference requests (4MB).
const maxPerformancePayloadSize = 4 * 1024 * 1024

// UserTestPreferenceHandler is used for persisting load test preferences
func (h *Handler) UserTestPreferenceHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	switch req.Method {
	case http.MethodPost:
		h.UserTestPreferenceStore(w, req, prefObj, user, provider)
	case http.MethodDelete:
		h.UserTestPreferenceDelete(w, req, prefObj, user, provider)
	case http.MethodGet:
		h.UserTestPreferenceGet(w, req, prefObj, user, provider)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// UserTestPreferenceStore is used for persisting load test preferences
func (h *Handler) UserTestPreferenceStore(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	if !h.requireProvider(w, provider) {
		return
	}

	defer func() {
		_ = req.Body.Close()
	}()

	body, err := io.ReadAll(io.LimitReader(req.Body, maxPerformancePayloadSize+1))
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(w, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	if len(body) > maxPerformancePayloadSize {
		err := fmt.Errorf("request body exceeds maximum allowed size of %dMB", maxPerformancePayloadSize/(1024*1024))
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, err, http.StatusRequestEntityTooLarge)
		return
	}
	if len(body) == 0 {
		if h.log != nil {
			h.log.Error(ErrRequestBody(fmt.Errorf("empty request body")))
		}
		writeMeshkitError(w, ErrRequestBody(fmt.Errorf("empty request body")), http.StatusBadRequest)
		return
	}

	perfTest := &perfprofile.PerformanceTestConfig{}
	if err = json.Unmarshal(body, perfTest); err != nil {
		obj := "the provided input"
		if h.log != nil {
			h.log.Error(models.ErrUnmarshal(err, obj))
		}
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusBadRequest)
		return
	}
	if err = models.PerformanceTestConfigValidator(perfTest); err != nil {
		if h.log != nil {
			h.log.Error(ErrRecordPreferences(err))
		}
		writeMeshkitError(w, ErrRecordPreferences(err), http.StatusBadRequest)
		return
	}
	tid, err := provider.SMPTestConfigStore(req, perfTest)
	if err != nil {
		obj := "user preference"
		if h.log != nil {
			h.log.Error(ErrFailToSave(err, obj))
		}
		writeMeshkitError(w, ErrFailToSave(err, obj), http.StatusInternalServerError)
		return
	}
	writeJSONMessage(w, map[string]string{"test_uuid": tid}, http.StatusOK)
}

// UserTestPreferenceGet gets the PerformanceTestConfig object
func (h *Handler) UserTestPreferenceGet(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	if !h.requireProvider(w, provider) {
		return
	}

	q := req.URL.Query()
	testUUID := q.Get("uuid")
	if h.log != nil {
		h.log.Debug(testUUID)
	}
	if testUUID == "" {
		testPage := q.Get("page")
		testPageSize := q.Get("pagesize")
		testSearch := q.Get("search")
		testOrder := q.Get("order")
		if h.log != nil {
			h.log.Debug(testPage, testPageSize)
		}
		testObjJSON, err := provider.SMPTestConfigFetch(req, testPage, testPageSize, testSearch, testOrder)
		if err != nil {
			if h.log != nil {
				h.log.Error(ErrTestConfigs)
			}
			writeMeshkitError(w, ErrTestConfigs, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(testObjJSON)
	} else {
		testObj, err := provider.SMPTestConfigGet(req, testUUID)
		if err != nil {
			if h.log != nil {
				h.log.Error(ErrTestConfigs)
			}
			writeMeshkitError(w, ErrTestConfigs, http.StatusInternalServerError)
			return
		}
		if testObj == nil {
			writeMeshkitError(w, ErrTestConfigs, http.StatusNotFound)
			return
		}
		writeJSONMessage(w, testObj, http.StatusOK)
	}
}

// UserTestPreferenceDelete deletes the PerformanceTestConfig object
func (h *Handler) UserTestPreferenceDelete(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	if !h.requireProvider(w, provider) {
		return
	}

	testUUID := req.URL.Query().Get("uuid")
	if testUUID == "" {
		obj := "field uuid"
		if h.log != nil {
			h.log.Error(ErrQueryGet(obj))
		}
		writeMeshkitError(w, ErrQueryGet(obj), http.StatusBadRequest)
		return
	}
	if err := provider.SMPTestConfigDelete(req, testUUID); err != nil {
		obj := "testConfig"
		if h.log != nil {
			h.log.Error(ErrFailToDelete(err, obj))
		}
		writeMeshkitError(w, ErrFailToDelete(err, obj), http.StatusInternalServerError)
		return
	}

	writeJSONEmptyObject(w, http.StatusOK)
}
