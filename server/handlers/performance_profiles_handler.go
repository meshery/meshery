package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
)

func performanceProfileIDFromRequest(r *http.Request) string {
	return mux.Vars(r)["performanceProfileId"]
}

// SavePerformanceProfileHandler will save performance profile using the current provider's persistence mechanism
func (h *Handler) SavePerformanceProfileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	body = bytes.TrimSpace(body)
	if len(body) == 0 || bytes.Equal(body, []byte("null")) {
		err := fmt.Errorf("performance profile request body is empty or null")
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	parsedBody := &models.PerformanceProfile{}
	err = json.Unmarshal(body, parsedBody)
	if err != nil {
		//failed to read request body
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	if user != nil && user.ID != uuid.Nil {
		parsedBody.Owner = user.ID
	}

	j, _ := json.Marshal(parsedBody)
	h.log.Info("performance profile is ", string(j))

	if provider == nil {
		writeJSONError(rw, "provider not initialized", http.StatusInternalServerError)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		//unable to save user config data
		h.log.Error(ErrRecordPreferences(err))
		writeMeshkitError(rw, ErrRecordPreferences(err), http.StatusInternalServerError)
		return
	}

	resp, err := provider.SavePerformanceProfile(token, parsedBody)
	if err != nil {
		obj := "performance profile"
		//fail to save performance profile
		h.log.Error(ErrFailToSave(err, obj))
		writeMeshkitError(rw, ErrFailToSave(err, obj), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}

// GetPerformanceProfilesHandler returns the list of all the performance profiles saved by the current user
// TODO: make sure cert data is not passed along and used only when test are run add a flag to control this
func (h *Handler) GetPerformanceProfilesHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()

	if provider == nil {
		writeJSONError(rw, "provider not initialized", http.StatusInternalServerError)
		return
	}
	tokenString, _ := r.Context().Value(models.TokenCtxKey).(string)

	resp, err := provider.GetPerformanceProfiles(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"))
	if err != nil {
		fetchErr := ErrFetchProfiles(err)
		h.log.Error(fetchErr)
		writeMeshkitError(rw, fetchErr, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}

// DeletePerformanceProfileHandler deletes a performance profile with the given id
func (h *Handler) DeletePerformanceProfileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	performanceProfileID := performanceProfileIDFromRequest(r)
	if performanceProfileID == "" {
		writeJSONError(rw, "performanceProfileId is required", http.StatusBadRequest)
		return
	}
	if provider == nil {
		writeJSONError(rw, "provider not initialized", http.StatusInternalServerError)
		return
	}

	resp, err := provider.DeletePerformanceProfile(r, performanceProfileID)
	if err != nil {
		obj := "performance profile"
		// a profile that does not exist is a 404, not a server failure
		if errors.GetCode(err) == models.ErrResultNotFoundCode || providerStatus(err) == http.StatusNotFound {
			h.log.Warnf("No %s with ID %q found to delete", obj, performanceProfileID)
			writeMeshkitError(rw, err, http.StatusNotFound)
			return
		}

		//fail to delete performance profile
		h.log.Error(ErrFailToDelete(err, obj))
		// Keep the status a remote provider returned (403, 409, ...) and fall
		// back to 500 only when the failure carries none, as the other delete
		// handlers do.
		writeMeshkitError(rw, ErrFailToDelete(err, obj), providerStatusOrInternal(err))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) GetPerformanceProfileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	performanceProfileID := performanceProfileIDFromRequest(r)
	if performanceProfileID == "" {
		writeJSONError(rw, "performanceProfileId is required", http.StatusBadRequest)
		return
	}
	if provider == nil {
		writeJSONError(rw, "provider not initialized", http.StatusInternalServerError)
		return
	}

	resp, err := provider.GetPerformanceProfile(r, performanceProfileID)
	if err != nil {
		fetchErr := ErrFetchProfile(err)
		h.log.Error(fetchErr)
		writeMeshkitError(rw, fetchErr, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}
