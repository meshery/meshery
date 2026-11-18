package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/meshery/meshery/server/models"
)

func (h *Handler) GetOrganizations(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	if !h.requireProvider(w, provider) {
		return
	}

	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		writeMeshkitError(w, models.ErrGetToken(errors.New("token missing from request context")), http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()
	resp, err := provider.GetOrganizations(token, q.Get("page"), getPageSizeParam(q), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		handlerErr := ErrGetOrganizations(err)
		if h != nil && h.log != nil {
			h.log.Error(handlerErr)
		}
		writeMeshkitError(w, handlerErr, providerStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(w, string(resp)); err != nil {
		if h != nil && h.log != nil {
			h.log.Error(err)
		}
	}
}
