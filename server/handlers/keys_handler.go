package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/layer5io/meshery/server/models"
)

// swagger:route GET /api/identity/orgs/{orgID}/users/keys UserKeysAPI idGetAllUsersKeysHandler;
// Handles GET for all Keys for users
//
// ```?order={field}``` orders on the passed field
//
// ```?search={}``` If search is non empty then a greedy search is performed
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 10
//
// ```?filter={filter}``` Filter keys
// responses:
// 200: keys

func (h *Handler) GetUsersKeys(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()
	orgID := mux.Vars(req)["orgID"]
	resp, err := provider.GetUsersKeys(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"), orgID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}
