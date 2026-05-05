// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/meshery/meshery/server/models"
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

// UserTestPreferenceStore is used for persisting load test preferences
func (h *Handler) UserTestPreferenceStore(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	perfTest := &SMP.PerformanceTestConfig{}
	if err = protojson.Unmarshal(body, perfTest); err != nil {
		obj := "the provided input"
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusBadRequest)
		return
	}
	if err = models.SMPPerformanceTestConfigValidator(perfTest); err != nil {
		h.log.Error(ErrRecordPreferences(err))
		writeMeshkitError(w, ErrRecordPreferences(err), http.StatusBadRequest)
		return
	}
	tid, err := provider.SMPTestConfigStore(req, perfTest)
	if err != nil {
		obj := "user preference"
		h.log.Error(ErrFailToSave(err, obj))
		writeMeshkitError(w, ErrFailToSave(err, obj), http.StatusInternalServerError)
		return
	}
	writeJSONMessage(w, map[string]string{"test_uuid": tid}, http.StatusOK)
}

// UserTestPreferenceGet gets the PerformanceTestConfig object
func (h *Handler) UserTestPreferenceGet(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	q := req.URL.Query()
	testUUID := q.Get("uuid")
	h.log.Debug(testUUID)
	if testUUID == "" {
		testPage := q.Get("page")
		testPageSize := q.Get("pagesize")
		testSearch := q.Get("search")
		testOrder := q.Get("order")
		h.log.Debug(testPage, testPageSize)
		testObjJSON, err := provider.SMPTestConfigFetch(req, testPage, testPageSize, testSearch, testOrder)
		if err != nil {
			h.log.Error(ErrTestConfigs)
			writeMeshkitError(w, ErrTestConfigs, http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(testObjJSON)
	} else {
		testObj, err := provider.SMPTestConfigGet(req, testUUID)
		if err != nil {
			h.log.Error(ErrTestConfigs)
			writeMeshkitError(w, ErrTestConfigs, http.StatusInternalServerError)
			return
		}
		if testObj == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Printf("%v", testObj)
		data, err := protojson.Marshal(testObj)
		if err != nil {
			h.log.Error(ErrReadConfig(err))
			writeMeshkitError(w, ErrReadConfig(err), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(data)
		if err != nil {
			// Response body has already been written (protojson.Marshal
			// succeeded and w.Write is partway through) — a fresh error
			// response would corrupt the in-flight body, so log only.
			h.log.Error(ErrWriteResponse(err))
			return
		}
	}
}

// UserTestPreferenceDelete deletes the PerformanceTestConfig object
func (h *Handler) UserTestPreferenceDelete(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	testUUID := req.URL.Query().Get("uuid")
	if testUUID == "" {
		obj := "field uuid"
		h.log.Error(ErrQueryGet(obj))
		writeMeshkitError(w, ErrQueryGet(obj), http.StatusBadRequest)
		return
	}
	if err := provider.SMPTestConfigDelete(req, testUUID); err != nil {
		obj := "testConfig"
		h.log.Error(ErrFailToDelete(err, obj))
		writeMeshkitError(w, ErrFailToDelete(err, obj), http.StatusInternalServerError)
		return
	}
}
