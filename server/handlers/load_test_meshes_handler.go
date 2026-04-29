package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/meshery/meshery/server/models"
)

// SMPMeshes defines the JSON payload structure for available meshes api
type SMPMeshes struct {
	AvailableMeshes []string `json:"availableMeshes,omitempty"`
}

// GetSMPServiceMeshes handles the available meshes request
func (h *Handler) GetSMPServiceMeshes(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	meshes := SMPMeshes{
		AvailableMeshes: make([]string, 0, len(SMP.ServiceMesh_Type_name)),
	}

	for _, v := range SMP.ServiceMesh_Type_name {
		if v != SMP.ServiceMesh_INVALID_MESH.String() {
			meshes.AvailableMeshes = append(meshes.AvailableMeshes, cases.Title(language.Und).String(strings.ToLower(strings.ReplaceAll(v, "_", " "))))
		}
	}

	if err := json.NewEncoder(w).Encode(meshes); err != nil {
		// Response body has already started streaming via json.Encoder —
		// a partial JSON envelope is on the wire and a fresh error
		// response would corrupt it, so log only.
		obj := "meshlist object"
		h.log.Error(models.ErrEncoding(err, obj))
		return
	}
}
