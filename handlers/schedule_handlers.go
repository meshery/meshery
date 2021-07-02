package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/models"
)

// swagger:route POST /api/user/schedules SchedulesAPI idPostSchedules
// Handle POST reqeuest for Schedules
//
// Save schedule using the current provider's persistence mechanism
// responses:
// 	200: singleScheduleResponseWrapper

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
		fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		http.Error(rw, "failed to get user token", http.StatusInternalServerError)
		return
	}

	resp, err := provider.SaveSchedule(token, parsedBody)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to save the schedule: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route GET /api/user/schedules SchedulesAPI idGetSchedules
// Handle GET reqeuest for Schedules
//
// Returns the list of all the schedules saved by the current user
// responses:
// 	200: schedulesResponseWrapper

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
		http.Error(rw, fmt.Sprintf("failed to fetch the schedules: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route DELETE /api/user/schedules/{id} SchedulesAPI idDeleteSchedules
// Handle DELETE reqeuest for Schedules
//
// Deletes a schedule with the given id
// responses:
// 	200: schedulesResponseWrapper

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
		http.Error(rw, fmt.Sprintf("failed to delete the schedule: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route GET /api/user/schedules/{id} SchedulesAPI idGetSingleSchedule
// Handle GET reqeuest for Schedules
//
// Fetches and returns the schedule with the given id
// responses:
// 	200: singleScheduleResponseWrapper

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
		http.Error(rw, fmt.Sprintf("failed to get the schedule: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}
