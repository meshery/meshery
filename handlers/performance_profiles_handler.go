package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/models"
)

// SavePerformanceProfileHandler will save performance profile using the current provider's persistence mechanism
func (h *Handler) SavePerformanceProfileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	var parsedBody *models.PerformanceProfile
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}
	fmt.Printf("%+v\n", parsedBody)

	token, err := provider.GetProviderToken(r)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to get user token"), http.StatusInternalServerError)
		return
	}

	resp, err := provider.SavePerformanceProfile(token, parsedBody)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to save the performance profile: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// GetPerformanceProfilesHandler returns the list of all the performance profiles saved by the current user
func (h *Handler) GetPerformanceProfilesHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()

	resp, err := provider.GetPerformanceProfiles(r, q.Get("page"), q.Get("page_size"), q.Get("search"), q.Get("order"))
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to fetch the performance profiles: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// DeletePerformanceProfileHandler deletes a performance profile with the given id
func (h *Handler) DeletePerformanceProfileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	performanceProfileID := mux.Vars(r)["id"]

	resp, err := provider.DeletePerformanceProfile(r, performanceProfileID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to delete the performance profile: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// GetPerformanceProfileHandler fetched the performance profile with the given id
func (h *Handler) GetPerformanceProfileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	performanceProfileID := mux.Vars(r)["id"]

	resp, err := provider.GetPerformanceProfile(r, performanceProfileID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to get the performance profile: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}
