package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/meshery/meshery/server/models"
)

func (h *Handler) GetUsersKeys(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		writeMeshkitError(w, models.ErrGetToken(errors.New("token missing from request context")), http.StatusInternalServerError)
		return
	}

	if provider == nil {
		nilProviderErr := ErrNilProvider()
		h.log.Error(nilProviderErr)
		writeMeshkitError(w, nilProviderErr, http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()
	orgID := mux.Vars(req)["orgID"]
	if orgID != "" {
		if _, err := uuid.FromString(orgID); err != nil {
			invalidOrgIDErr := models.ErrInvalidUUID(fmt.Errorf("invalid organization ID %q: %w", orgID, err))
			h.log.Error(invalidOrgIDErr)
			writeMeshkitError(w, invalidOrgIDErr, http.StatusBadRequest)
			return
		}
	}

	resp, err := provider.GetUsersKeys(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"), orgID)
	if err != nil {
		handlerErr := ErrGetUsersKeys(err)
		h.log.Error(handlerErr)
		writeMeshkitError(w, handlerErr, providerStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(w, string(resp)); err != nil {
		h.log.Error(err)
	}
}
