package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/meshery/meshery/server/models"
)

func (h *Handler) GetOrganizations(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		writeMeshkitError(w, models.ErrGetToken(errors.New("token missing from request context")), http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()
	resp, err := provider.GetOrganizations(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(w, string(resp)); err != nil {
		h.log.Error(err)
	}
}
