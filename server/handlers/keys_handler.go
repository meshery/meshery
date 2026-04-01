package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/meshery/meshery/server/models"
)

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
	if _, err := fmt.Fprint(w, string(resp)); err != nil {
		h.log.Error(err)
	}
}
