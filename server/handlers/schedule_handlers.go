package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
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
		writeMeshkitError(rw, ErrFailToSave(err, obj), http.StatusInternalServerError)
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
		obj := "schedules"
		//unable to get schedules
		if h != nil && h.log != nil {
			h.log.Error(ErrQueryGet(obj))
		}
		writeMeshkitError(rw, ErrQueryGet(obj), http.StatusInternalServerError)
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

	resp, err := provider.DeleteSchedule(r, scheduleID)
	if err != nil {
		obj := "schedule"
		//unable to delete schedules
		if h != nil && h.log != nil {
			h.log.Error(ErrFailToDelete(err, obj))
		}
		writeMeshkitError(rw, ErrFailToDelete(err, obj), http.StatusInternalServerError)
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

	resp, err := provider.GetSchedule(r, scheduleID)
	if err != nil {
		obj := "schedule"
		//failed to get schedules
		if h != nil && h.log != nil {
			h.log.Error(ErrQueryGet(obj))
		}
		writeMeshkitError(rw, ErrQueryGet(obj), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		if h != nil && h.log != nil {
			h.log.Error(err)
		}
	}
}
