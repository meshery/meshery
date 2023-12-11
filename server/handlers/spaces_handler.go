package handlers

import (
	"fmt"
	"net/http"
	// "time"

	// "encoding/json"

	// "github.com/gorilla/mux"
	// "github.com/pkg/errors"

	"github.com/layer5io/meshery/server/models"
)


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
