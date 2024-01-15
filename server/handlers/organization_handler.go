package handlers

import (
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/server/models"
)

// swagger:route GET /api/identity/orgs OrganizationsAPI idGetOrganizations
// Handles GET for all Organizations
//
//
// ```?order={field}``` orders on the passed field
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 20
//
// ```?search={environments_name}``` If search is non empty then a greedy search is performed
//
// ```?filter={condition}```
// responses:
// 	200: orgsResponseWrapper

func (h *Handler) GetOrganizations(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()
	resp, err := provider.GetOrganizations(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}
