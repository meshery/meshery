package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/layer5io/meshery/models"
	SMP "github.com/layer5io/service-mesh-performance/tree/v0.3.1/spec"
	"github.com/sirupsen/logrus"
)

// SMPMeshes defines the JSON payload structure for available meshes api
type SMPMeshes struct {
	AvailableMeshes []string `json:"available_meshes,omitempty"`
}

// GetSMPServiceMeshes handles the available meshes request
func (h *Handler) GetSMPServiceMeshes(w http.ResponseWriter, r *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	meshes := SMPMeshes{
		AvailableMeshes: make([]string, 0, len(SMP.ServiceMeshes_Type_name)),
	}

	for _, v := range SMP.ServiceMeshes_Type_name {
		if v != SMP.ServiceMeshes_INVALID_MESH.String() {
			meshes.AvailableMeshes = append(meshes.AvailableMeshes, strings.Title(strings.ToLower(strings.ReplaceAll(v, "_", " "))))
		}
	}

	if err := json.NewEncoder(w).Encode(meshes); err != nil {
		logrus.Errorf("error encoding meshlist object: %v", err)
		http.Error(w, "Error encoding meshlist object", http.StatusInternalServerError)
		return
	}
}
