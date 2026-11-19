package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"gorm.io/gorm"
)

// SaveScheduleHandler will save schedule using the current provider's persistence mechanism
func (h *Handler) SaveScheduleHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	if !h.requireProvider(rw, provider) {
		return
	}

	var parsedBody *models.Schedule
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		//failed to read request body
		if h != nil && h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	if parsedBody == nil {
		err := fmt.Errorf("request body cannot be empty")
		if h != nil && h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		//failed to get user token
		if h != nil && h.log != nil {
			h.log.Error(ErrRetrieveUserToken(err))
		}
		writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	resp, err := provider.SaveSchedule(token, parsedBody)
	if err != nil {
		obj := "schedule"
		//Failed to save the schedule
		if h != nil && h.log != nil {
			h.log.Error(ErrFailToSave(err, obj))
		}
		writeMeshkitError(rw, ErrFailToSave(err, obj), providerStatusOrInternal(err))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		if h != nil && h.log != nil {
			h.log.Error(err)
		}
	}
}

func (h *Handler) GetSchedulesHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	if !h.requireProvider(rw, provider) {
		return
	}

	q := r.URL.Query()

	resp, err := provider.GetSchedules(r, q.Get("page"), q.Get("pagesize"), q.Get("order"))
	if err != nil {
		if isPaginationError(err) {
			writeMeshkitError(rw, err, http.StatusBadRequest)
			return
		}
		obj := "schedules"
		//unable to get schedules
		if h != nil && h.log != nil {
			h.log.Error(ErrQueryGet(obj))
		}
		writeMeshkitError(rw, ErrQueryGet(obj), providerStatusOrInternal(err))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		if h != nil && h.log != nil {
			h.log.Error(err)
		}
	}
}

// DeleteScheduleHandler deletes a schedule with the given id
func (h *Handler) DeleteScheduleHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	if !h.requireProvider(rw, provider) {
		return
	}

	scheduleID := mux.Vars(r)["id"]
	if scheduleID == "" {
		missingIDErr := ErrMissingRouteVariable("id")
		if h != nil && h.log != nil {
			h.log.Error(missingIDErr)
		}
		writeMeshkitError(rw, missingIDErr, http.StatusBadRequest)
		return
	}

	// Validate the UUID format using the shared helper.
	if _, ok := h.requirePathUUID(rw, r, "id"); !ok {
		return
	}
	// Additionally reject the zero UUID which the helper treats as valid.
	if scheduleID == uuid.Nil.String() {
		invalidUUIDErr := models.ErrInvalidUUID(fmt.Errorf("invalid schedule id: %q", scheduleID))
		if h != nil && h.log != nil {
			h.log.Error(invalidUUIDErr)
		}
		writeMeshkitError(rw, invalidUUIDErr, http.StatusBadRequest)
		return
	}

	resp, err := provider.DeleteSchedule(r, scheduleID)
	if err != nil {
		obj := "schedule"
		// a schedule that does not exist is a 404, not a server failure
		if errors.Is(err, gorm.ErrRecordNotFound) || meshkiterrors.GetCode(err) == models.ErrResultNotFoundCode || providerStatus(err) == http.StatusNotFound {
			if h != nil && h.log != nil {
				h.log.Warnf("No %s with ID %q found to delete", obj, scheduleID)
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = models.ErrResultNotFound(err)
			}
			writeMeshkitError(rw, err, http.StatusNotFound)
			return
		}

		//unable to delete schedules
		if h != nil && h.log != nil {
			h.log.Error(ErrFailToDelete(err, obj))
		}
		writeMeshkitError(rw, ErrFailToDelete(err, obj), providerStatusOrInternal(err))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		if h != nil && h.log != nil {
			h.log.Error(err)
		}
	}
}

// GetScheduleHandler fetches the schedule with the given id
func (h *Handler) GetScheduleHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	if !h.requireProvider(rw, provider) {
		return
	}

	scheduleID := mux.Vars(r)["id"]
	if scheduleID == "" {
		missingIDErr := ErrMissingRouteVariable("id")
		if h != nil && h.log != nil {
			h.log.Error(missingIDErr)
		}
		writeMeshkitError(rw, missingIDErr, http.StatusBadRequest)
		return
	}

	// Validate the UUID format using the shared helper.
	if _, ok := h.requirePathUUID(rw, r, "id"); !ok {
		return
	}
	// Additionally reject the zero UUID which the helper treats as valid.
	if scheduleID == uuid.Nil.String() {
		invalidUUIDErr := models.ErrInvalidUUID(fmt.Errorf("invalid schedule id: %q", scheduleID))
		if h != nil && h.log != nil {
			h.log.Error(invalidUUIDErr)
		}
		writeMeshkitError(rw, invalidUUIDErr, http.StatusBadRequest)
		return
	}

	resp, err := provider.GetSchedule(r, scheduleID)
	if err != nil {
		obj := "schedule"
		if errors.Is(err, gorm.ErrRecordNotFound) || meshkiterrors.GetCode(err) == models.ErrResultNotFoundCode || providerStatus(err) == http.StatusNotFound {
			if h != nil && h.log != nil {
				h.log.Warnf("No %s with ID %q found", obj, scheduleID)
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = models.ErrResultNotFound(err)
			}
			writeMeshkitError(rw, err, http.StatusNotFound)
			return
		}

		//failed to get schedules
		if h != nil && h.log != nil {
			h.log.Error(ErrQueryGet(obj))
		}
		writeMeshkitError(rw, ErrQueryGet(obj), providerStatusOrInternal(err))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		if h != nil && h.log != nil {
			h.log.Error(err)
		}
	}
}
