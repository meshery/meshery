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
		tokenErr := models.ErrGetToken(errors.New("token missing from request context"))
		if h != nil && h.log != nil {
			h.log.Error(tokenErr)
		}
		writeMeshkitError(w, tokenErr, http.StatusInternalServerError)
		return
	}

	if provider == nil {
		nilProviderErr := ErrNilProvider()
		if h != nil && h.log != nil {
			h.log.Error(nilProviderErr)
		}
		writeMeshkitError(w, nilProviderErr, http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()
	orgID := mux.Vars(req)["orgID"]
	if orgID != "" {
		parsedOrgID, err := uuid.FromString(orgID)
		if err != nil || parsedOrgID == uuid.Nil {
			var invalidOrgIDErr error
			if err != nil {
				invalidOrgIDErr = models.ErrInvalidUUID(fmt.Errorf("invalid organization ID %q: %w", orgID, err))
			} else {
				invalidOrgIDErr = models.ErrInvalidUUID(fmt.Errorf("invalid organization ID %q", orgID))
			}
			if h != nil && h.log != nil {
				h.log.Error(invalidOrgIDErr)
			}
			writeMeshkitError(w, invalidOrgIDErr, http.StatusBadRequest)
			return
		}
	}

	resp, err := provider.GetUsersKeys(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"), orgID)
	if err != nil {
		handlerErr := ErrGetUsersKeys(err)
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
