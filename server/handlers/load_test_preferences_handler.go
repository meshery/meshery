// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
	"gorm.io/gorm"
)

func parseTestPreferenceUUID(raw string) (string, error) {
	uid, err := uuid.FromString(raw)
	if err != nil || uid == uuid.Nil {
		return "", models.ErrInvalidUUID(fmt.Errorf("invalid test config uuid: %q", raw))
	}
	return uid.String(), nil
}

func testPreferenceErrorStatus(err error) int {
	switch errors.GetCode(err) {
	case models.ErrInvalidUUIDValueCode, models.ErrPageNumberCode, models.ErrPageSizeCode:
		return http.StatusBadRequest
	case models.ErrResultNotFoundCode:
		return http.StatusNotFound
	}
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound
	}
	return providerStatusOrInternal(err)
}

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
		writeMeshkitError(w, ErrMethodNotAllowed(req.Method), http.StatusMethodNotAllowed)
	}
}

// UserTestPreferenceStore is used for persisting load test preferences
func (h *Handler) UserTestPreferenceStore(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	defer func() {
		_ = req.Body.Close()
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	perfTest := &perfprofile.PerformanceTestConfig{}
	if err = json.Unmarshal(body, perfTest); err != nil {
		obj := "the provided input"
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusBadRequest)
		return
	}
	if err = models.PerformanceTestConfigValidator(perfTest); err != nil {
		h.log.Error(ErrRecordPreferences(err))
		writeMeshkitError(w, ErrRecordPreferences(err), http.StatusBadRequest)
		return
	}
	tid, err := provider.SMPTestConfigStore(req, perfTest)
	if err != nil {
		obj := "user preference"
		h.log.Error(ErrFailToSave(err, obj))
		writeMeshkitError(w, ErrFailToSave(err, obj), testPreferenceErrorStatus(err))
		return
	}
	writeJSONMessage(w, map[string]string{"test_uuid": tid}, http.StatusOK)
}

// UserTestPreferenceGet gets the PerformanceTestConfig object
func (h *Handler) UserTestPreferenceGet(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	q := req.URL.Query()
	testUUID := q.Get("uuid")
	h.log.Debug(testUUID)
	if !q.Has("uuid") {
		testPage := q.Get("page")
		testPageSize := q.Get("pagesize")
		testSearch := q.Get("search")
		testOrder := q.Get("order")
		h.log.Debug(testPage, testPageSize)
		testObjJSON, err := provider.SMPTestConfigFetch(req, testPage, testPageSize, testSearch, testOrder)
		if err != nil {
			status := testPreferenceErrorStatus(err)
			if status == http.StatusBadRequest {
				h.log.Warn(err)
				writeMeshkitError(w, err, status)
				return
			}
			fetchErr := ErrTestConfigs(err)
			h.log.Error(fetchErr)
			writeMeshkitError(w, fetchErr, status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(testObjJSON); err != nil {
			h.log.Error(ErrWriteResponse(err))
		}
		return
	}

	validUUID, err := parseTestPreferenceUUID(testUUID)
	if err != nil {
		h.log.Warn(err)
		writeMeshkitError(w, err, http.StatusBadRequest)
		return
	}

	testObj, err := provider.SMPTestConfigGet(req, validUUID)
	if err != nil {
		status := testPreferenceErrorStatus(err)
		if status == http.StatusNotFound {
			h.log.Warnf("Test config with UUID %q not found", validUUID)
			writeMeshkitError(w, err, http.StatusNotFound)
			return
		}
		fetchErr := ErrTestConfigs(err)
		h.log.Error(fetchErr)
		writeMeshkitError(w, fetchErr, status)
		return
	}
	if testObj == nil {
		writeMeshkitError(w, models.ErrResultNotFound(fmt.Errorf("test config %s not found", validUUID)), http.StatusNotFound)
		return
	}
	data, err := json.Marshal(testObj)
	if err != nil {
		h.log.Error(ErrReadConfig(err))
		writeMeshkitError(w, ErrReadConfig(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(data); err != nil {
		// Response body has already been written (json.Marshal
		// succeeded and w.Write is partway through) — a fresh error
		// response would corrupt the in-flight body, so log only.
		h.log.Error(ErrWriteResponse(err))
		return
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
	validUUID, err := parseTestPreferenceUUID(testUUID)
	if err != nil {
		h.log.Warn(err)
		writeMeshkitError(w, err, http.StatusBadRequest)
		return
	}
	if err := provider.SMPTestConfigDelete(req, validUUID); err != nil {
		status := testPreferenceErrorStatus(err)
		if status == http.StatusNotFound {
			h.log.Warnf("No test config with UUID %q found to delete", validUUID)
			writeMeshkitError(w, err, http.StatusNotFound)
			return
		}
		obj := "testConfig"
		h.log.Error(ErrFailToDelete(err, obj))
		writeMeshkitError(w, ErrFailToDelete(err, obj), status)
		return
	}
	writeJSONEmptyObject(w, http.StatusOK)
}
