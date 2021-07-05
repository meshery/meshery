package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/models"
)

// SaveScheduleHandler will save schedule using the current provider's persistence mechanism
func (h *Handler) SaveScheduleHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	var parsedBody *models.Schedule
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		//failed to read request body
		//fmt.Fprintf(rw, ErrRequestBody(err).Error(), err)
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		h.log.Error(ErrRequestBody(err))
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		//failed to get user token
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		h.log.Error(ErrRetrieveUserToken(err))
		return
	}

	resp, err := provider.SaveSchedule(token, parsedBody)
	if err != nil {
		obj := "schedule"
		//Failed to save the schedule
		http.Error(rw, ErrFailToSave(err, obj).Error(), http.StatusInternalServerError)
		h.log.Error(ErrFailToSave(err, obj))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// GetSchedulesHandler returns the list of all the schedules saved by the current user
func (h *Handler) GetSchedulesHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()

	resp, err := provider.GetSchedules(r, q.Get("page"), q.Get("page_size"), q.Get("order"))
	if err != nil {
		obj := "schedules"
		//unable to get schedules
		http.Error(rw, ErrQueryGet(obj).Error(), http.StatusInternalServerError)
		h.log.Error(ErrQueryGet(obj))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// DeleteScheduleHandler deletes a schedule with the given id
func (h *Handler) DeleteScheduleHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	ScheduleID := mux.Vars(r)["id"]

	resp, err := provider.DeleteSchedule(r, ScheduleID)
	if err != nil {
		obj := "schedule"
		//unable to delete schedules
		http.Error(rw, ErrFailToDelete(err, obj).Error(), http.StatusInternalServerError)
		h.log.Error(ErrFailToDelete(err, obj))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// GetScheduleHandler fetches the schedule with the given id
func (h *Handler) GetScheduleHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	ScheduleID := mux.Vars(r)["id"]

	resp, err := provider.GetSchedule(r, ScheduleID)
	if err != nil {
		obj := "schedule"
		//failed to get schedules
		http.Error(rw, ErrQueryGet(obj).Error(), http.StatusInternalServerError)
		h.log.Error(ErrQueryGet(obj))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}
