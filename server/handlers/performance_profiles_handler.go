package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"gorm.io/gorm"
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

	const maxBodyLimit = 4 << 20 // 4MB
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyLimit))
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	body = bytes.TrimSpace(body)
	if len(body) == 0 || bytes.Equal(body, []byte("null")) {
		err := fmt.Errorf("performance profile request body is empty or null")
		if h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	parsedBody := &models.PerformanceProfile{}
	err = json.Unmarshal(body, parsedBody)
	if err != nil {
		// failed to read request body
		if h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	if user != nil && user.ID != uuid.Nil {
		parsedBody.Owner = user.ID
	}

	if h.log != nil {
		j, _ := json.Marshal(parsedBody)
		h.log.Info("performance profile is ", string(j))
	}

	if provider == nil {
		writeJSONError(rw, "provider not initialized", http.StatusInternalServerError)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		// unable to save user config data
		if h.log != nil {
			h.log.Error(ErrRecordPreferences(err))
		}
		writeMeshkitError(rw, ErrRecordPreferences(err), http.StatusInternalServerError)
		return
	}

	resp, err := provider.SavePerformanceProfile(token, parsedBody)
	if err != nil {
		obj := "performance profile"
		// fail to save performance profile
		if h.log != nil {
			h.log.Error(ErrFailToSave(err, obj))
		}
		writeMeshkitError(rw, ErrFailToSave(err, obj), providerStatusOrInternal(err))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		if h.log != nil {
			h.log.Error(err)
		}
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
	tokenString, ok := r.Context().Value(models.TokenCtxKey).(string)
	if !ok || tokenString == "" {
		fetchTokenErr := ErrFetchToken(fmt.Errorf("token not found in request context"))
		if h.log != nil {
			h.log.Error(fetchTokenErr)
		}
		writeMeshkitError(rw, fetchTokenErr, http.StatusInternalServerError)
		return
	}

	resp, err := provider.GetPerformanceProfiles(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"))
	if err != nil {
		fetchErr := ErrFetchProfiles(err)
		if h.log != nil {
			h.log.Error(fetchErr)
		}
		writeMeshkitError(rw, fetchErr, providerStatusOrInternal(err))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		if h.log != nil {
			h.log.Error(err)
		}
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

	key, err := uuid.FromString(performanceProfileID)
	if err != nil || key == uuid.Nil {
		invalidUUIDErr := models.ErrInvalidUUID(fmt.Errorf("invalid performance profile id: %q", performanceProfileID))
		if h.log != nil {
			h.log.Error(invalidUUIDErr)
		}
		writeMeshkitError(rw, invalidUUIDErr, http.StatusBadRequest)
		return
	}

	resp, err := provider.DeletePerformanceProfile(r, performanceProfileID)
	if err != nil {
		obj := "performance profile"
		// a profile that does not exist is a 404, not a server failure
		if errors.Is(err, gorm.ErrRecordNotFound) || meshkiterrors.GetCode(err) == models.ErrResultNotFoundCode || providerStatus(err) == http.StatusNotFound {
			if h.log != nil {
				h.log.Warnf("No %s with ID %q found to delete", obj, performanceProfileID)
			}
			writeMeshkitError(rw, err, http.StatusNotFound)
			return
		}

		// fail to delete performance profile
		if h.log != nil {
			h.log.Error(ErrFailToDelete(err, obj))
		}
		// Keep the status a remote provider returned (403, 409, ...) and fall
		// back to 500 only when the failure carries none, as the other delete
		// handlers do.
		writeMeshkitError(rw, ErrFailToDelete(err, obj), providerStatusOrInternal(err))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		if h.log != nil {
			h.log.Error(err)
		}
	}
}

// GetPerformanceProfileHandler fetches a performance profile with the given id from the current provider
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

	key, err := uuid.FromString(performanceProfileID)
	if err != nil || key == uuid.Nil {
		invalidUUIDErr := models.ErrInvalidUUID(fmt.Errorf("invalid performance profile id: %q", performanceProfileID))
		if h.log != nil {
			h.log.Error(invalidUUIDErr)
		}
		writeMeshkitError(rw, invalidUUIDErr, http.StatusBadRequest)
		return
	}

	resp, err := provider.GetPerformanceProfile(r, performanceProfileID)
	if err != nil {
		obj := "performance profile"
		// a profile that does not exist is a 404, not a server failure
		if errors.Is(err, gorm.ErrRecordNotFound) || meshkiterrors.GetCode(err) == models.ErrResultNotFoundCode || providerStatus(err) == http.StatusNotFound {
			if h.log != nil {
				h.log.Warnf("No %s with ID %q found", obj, performanceProfileID)
			}
			writeMeshkitError(rw, err, http.StatusNotFound)
			return
		}

		fetchErr := ErrFetchProfile(err)
		if h.log != nil {
			h.log.Error(fetchErr)
		}
		writeMeshkitError(rw, fetchErr, providerStatusOrInternal(err))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		if h.log != nil {
			h.log.Error(err)
		}
	}
}
