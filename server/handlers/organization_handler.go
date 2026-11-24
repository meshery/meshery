package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/meshery/meshery/server/models"
)

// organizationsPageSize reads the page size of GET /api/identity/orgs. The
// schemas contract (v1beta2 getOrgs), and so the generated UI client, sends
// camelCase pageSize and all=true; pagesize is the legacy spelling that
// mesheryctl still sends. all=true asks for every organization.
func organizationsPageSize(q url.Values) string {
	if q.Get("all") == queryParamTrue {
		return "all"
	}
	if pageSize := q.Get("pageSize"); pageSize != "" {
		return pageSize
	}
	return q.Get("pagesize")
}

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
	resp, err := provider.GetOrganizations(token, q.Get("page"), organizationsPageSize(q), q.Get("search"), q.Get("order"), q.Get("filter"))
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
