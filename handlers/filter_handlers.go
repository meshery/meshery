package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/models"
)

// FetchAllFiltersHandler will handle requests of both type GET and POST
// on the route /api/experimental/filter
func (h *Handler) FetchAllFiltersHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	if r.Method == http.MethodGet {
		h.GetMesheryFiltersHandler(rw, r, prefObj, user, provider)
		return
	}

	if r.Method == http.MethodPost {
		h.SaveFilterFile(rw, r, prefObj, user, provider)
		return
	}
}

// NewFilterHandler will handle requests type GET
// on the route /api/experimental/filter/{id}
func (h *Handler) NewFilterHandler(
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

func (h *Handler) SaveFilterFile(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	var parsedBody *models.MesheryFilter

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

	resp, err := provider.SaveMesheryFilter(token, parsedBody)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to save the filter: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// GetMesheryFiltersHandler returns the list of all the patterns saved by the current user
func (h *Handler) GetMesheryFiltersHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()

	resp, err := provider.GetMesheryFilters(r, q.Get("page"), q.Get("page_size"), q.Get("search"), q.Get("order"))
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to fetch the filters: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}
