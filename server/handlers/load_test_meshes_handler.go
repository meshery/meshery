package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/layer5io/meshery/server/models"
	SMP "github.com/layer5io/service-mesh-performance/spec"
)

// SMPMeshes defines the JSON payload structure for available meshes api
type SMPMeshes struct {
	AvailableMeshes []string `json:"available_meshes,omitempty"`
}

// GetSMPServiceMeshes handles the available meshes request
func (h *Handler) GetSMPServiceMeshes(w http.ResponseWriter, r *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	meshes := SMPMeshes{
		AvailableMeshes: make([]string, 0, len(SMP.ServiceMesh_Type_name)),
	}

	for _, v := range SMP.ServiceMesh_Type_name {
		if v != SMP.ServiceMesh_INVALID_MESH.String() {
			meshes.AvailableMeshes = append(meshes.AvailableMeshes, cases.Title(language.Und).String(strings.ToLower(strings.ReplaceAll(v, "_", " "))))
		}
	}

	if err := json.NewEncoder(w).Encode(meshes); err != nil {
		obj := "meshlist object"
		h.log.Error(ErrEncoding(err, obj))
		http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}
