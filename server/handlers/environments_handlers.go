package handlers

import (
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/server/models"
	"github.com/sirupsen/logrus"
)

// swagger:route GET /api/integrations/environments UserAPI idGetAllEnvironmentsHandler
// Handles GET for all Users
//
// # Users can be further filtered through query parameters
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
// 	200: users

func (h *Handler) GetEnvironments(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	logrus.Infof("Environments handler called")
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()

	resp, err := provider.GetEnvironments(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route GET /api/environment/{id} EnvironmentAPI idGetEnvironmentByIDHandler
// Handle GET for Environment info by ID
//
// Returns Environment info
// responses:
//   200: environmentInfo

// func (h *Handler) GetEnvironmentByIDHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
// 	environmentID := mux.Vars(r)["id"]
// 	resp, err := provider.GetEnvironmentByID(r, environmentID)
// 	if err != nil {
// 		h.log.Error(ErrGetResult(err))
// 		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Fprint(w, string(resp))
// }
