package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/meshery/meshery/server/models"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	_model "github.com/meshery/schemas/models/v1beta1/model"
)

// SMPMeshes defines the JSON payload structure for available meshes api
type SMPMeshes struct {
	AvailableMeshes []string `json:"availableMeshes,omitempty"`
}

// GetSMPServiceMeshes handles the available meshes request. The selectable
// technologies are sourced from the Meshery Registry's model names rather than
// the legacy, fixed service-mesh enum, so users can target any registered
// model and are not constrained to service-mesh technologies.
func (h *Handler) GetSMPServiceMeshes(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	if req != nil && req.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeMeshkitError(w, ErrMethodNotAllowed(req.Method), http.StatusMethodNotAllowed)
		return
	}

	if h.registryManager == nil {
		err := fmt.Errorf("registry manager not initialized")
		if h.log != nil {
			h.log.Error(ErrRetrieveData(err))
		}
		writeMeshkitError(w, ErrRetrieveData(err), http.StatusInternalServerError)
		return
	}

	entities, _, _, err := h.registryManager.GetEntities(&regv1beta1.ModelFilter{})
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrRetrieveData(err))
		}
		writeMeshkitError(w, ErrRetrieveData(err), http.StatusInternalServerError)
		return
	}

	seen := make(map[string]struct{})
	meshes := SMPMeshes{
		AvailableMeshes: make([]string, 0, len(entities)),
	}
	for _, entity := range entities {
		model, ok := entity.(*_model.ModelDefinition)
		if !ok || model == nil {
			continue
		}
		name := model.DisplayName
		if name == "" {
			name = model.Name
		}
		if name == "" {
			continue
		}
		if _, dup := seen[name]; dup {
			continue
		}
		seen[name] = struct{}{}
		meshes.AvailableMeshes = append(meshes.AvailableMeshes, name)
	}
	sort.Strings(meshes.AvailableMeshes)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(meshes); err != nil {
		// Response body has already started streaming via json.Encoder —
		// a partial JSON envelope is on the wire and a fresh error
		// response would corrupt it, so log only.
		obj := "meshlist object"
		if h.log != nil {
			h.log.Error(models.ErrEncoding(err, obj))
		}
		return
	}
}
