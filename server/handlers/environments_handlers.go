package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta1/environment"
)

// environmentPayloadWire is a handler-local dual-accept wrapper around
// environment.EnvironmentPayload. The schemas-generated struct tags
// OrgId as json:"organization_id" (required by the current published
// v1beta1 contract), but the canonical wire contract and every in-repo
// consumer now emit the camelCase `organizationId`. Because Go's
// encoding/json case-insensitive tag fallback will not match across an
// underscore boundary, a struct tagged `organization_id` would silently
// drop a JSON key of `organizationId`. This wrapper intercepts both
// spellings during the Phase 2 deprecation window. Canonical wins when
// both are present. Retire once schemas v1beta2 flips the tag and this
// handler consumes the new version.
type environmentPayloadWire struct {
	environment.EnvironmentPayload
}

func (p *environmentPayloadWire) UnmarshalJSON(data []byte) error {
	type alias environment.EnvironmentPayload
	aux := struct {
		*alias
		OrgIdCamel  *openapi_types.UUID `json:"organizationId,omitempty"`
		OrgIdSnake  *openapi_types.UUID `json:"organization_id,omitempty"`
	}{alias: (*alias)(&p.EnvironmentPayload)}

	// Zero OrgId so a reused receiver does not carry stale data when the
	// next payload omits both spellings.
	p.OrgId = openapi_types.UUID{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Canonical wins when both are supplied.
	switch {
	case aux.OrgIdCamel != nil:
		p.OrgId = *aux.OrgIdCamel
	case aux.OrgIdSnake != nil:
		p.OrgId = *aux.OrgIdSnake
	}
	return nil
}

func (h *Handler) GetEnvironments(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		writeMeshkitError(w, ErrFetchToken(fmt.Errorf("token not found in request context")), http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()

	orgID := q.Get("orgId")
	if orgID == "" {
		h.log.Error(errors.New("orgId is required"))
		writeJSONError(w, "orgId is required", http.StatusBadRequest)
		return
	}
	resp, err := provider.GetEnvironments(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"), orgID)
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

func (h *Handler) GetEnvironmentByIDHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	environmentID := mux.Vars(r)["id"]
	q := r.URL.Query()
	orgID := q.Get("orgId")
	if orgID == "" {
		h.log.Error(errors.New("orgId is required"))
		writeJSONError(w, "orgId is required", http.StatusBadRequest)
		return
	}
	resp, err := provider.GetEnvironmentByID(r, environmentID, orgID)
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

func (h *Handler) SaveEnvironment(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	wire := environmentPayloadWire{}
	err = json.Unmarshal(bd, &wire)
	obj := "environment"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusInternalServerError)
		return
	}

	environment := wire.EnvironmentPayload
	resp, err := provider.SaveEnvironment(req, &environment, "", false)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	description := fmt.Sprintf("Environment %s created.", environment.Name)

	h.log.Info(description)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := fmt.Fprint(w, string(resp)); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) DeleteEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	environmentID := mux.Vars(r)["id"]
	resp, err := provider.DeleteEnvironment(r, environmentID)
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

func (h *Handler) UpdateEnvironmentHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	environmentID := mux.Vars(req)["id"]
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	wire := environmentPayloadWire{}
	err = json.Unmarshal(bd, &wire)
	obj := "environment"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusInternalServerError)
		return
	}

	environment := wire.EnvironmentPayload
	resp, err := provider.UpdateEnvironment(req, &environment, environmentID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, models.ErrMarshal(err, "environment update response"), http.StatusInternalServerError)
		return
	}
	description := fmt.Sprintf("Environment %s updated.", environment.Name)
	h.log.Info(description)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(respJSON)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		// Headers already committed; log only. Writing another body would corrupt the stream.
		return
	}
}

func (h *Handler) AddConnectionToEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	environmentID := mux.Vars(r)["environmentID"]
	connectionID := mux.Vars(r)["connectionID"]
	resp, err := provider.AddConnectionToEnvironment(r, environmentID, connectionID)
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

func (h *Handler) RemoveConnectionFromEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	environmentID := mux.Vars(r)["environmentID"]
	connectionID := mux.Vars(r)["connectionID"]
	resp, err := provider.RemoveConnectionFromEnvironment(r, environmentID, connectionID)
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

func (h *Handler) GetConnectionsOfEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	environmentID := mux.Vars(r)["environmentID"]
	q := r.URL.Query()
	resp, err := provider.GetConnectionsOfEnvironment(r, environmentID, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(w, string(resp)); err != nil {
		h.log.Error(err)
	}
}
