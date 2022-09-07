package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/layer5io/meshery/server/models"
	meshkitmodels "github.com/layer5io/meshkit/models"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

type generationPayloadItem struct {
	Name     string `json:"name"`
	Register bool   `json:"register"`
}

type componentGenerationPayload struct {
	Data []generationPayloadItem `json:"data"`
}

type componentGenerationResponseDataItem struct {
	Name       string               `json:"name"`
	Components []v1alpha1.Component `json:"components"`
	Errors     []string             `json:"errors"`
}

// swagger:route POST /api/meshmodel/component/generate MeshmodelComponentGenerate idPostMeshModelComponentGenerate
// Handle POST request for component generation
//
// Generates Meshery Components for the given payload
// responses:
// 	200:

// request body should be json
// request body should be of format - {data: [{name: string, register: boolean}]}
// response format - {data: [{name: string, components: [component], errors: [string] }]}
func (h *Handler) ComponentGenerationHandler(rw http.ResponseWriter, r *http.Request) {
	// Parse the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}
	// Unmarshal request body
	pld := componentGenerationPayload{}
	err = json.Unmarshal(body, &pld)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		return
	}
	// Generate Components
	response := make([]componentGenerationResponseDataItem, 0)
	for _, gpi := range pld.Data {
		responseItem := componentGenerationResponseDataItem{Name: gpi.Name}
		ahpm := models.ArtifactHubPackageManager{
			PackageName: gpi.Name,
		}
		comps, err := generateComponents(ahpm)
		if err != nil {
			h.log.Error(ErrGenerateComponents(err))
			responseItem.Errors = append(responseItem.Errors, err.Error())
			response = append(response, responseItem)
			continue
		}
		responseItem.Components = comps
		h.log.Info("component generation for ", gpi.Name, " successfully completed")
		response = append(response, responseItem)
	}
	// Send response
	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		h.log.Error(ErrGenerateComponents(err))
		http.Error(rw, ErrGenerateComponents(err).Error(), http.StatusInternalServerError)
		return
	}
}

func generateComponents(pm meshkitmodels.PackageManager) ([]v1alpha1.Component, error) {
	components := make([]v1alpha1.Component, 0)
	pkg, err := pm.GetPackage()
	if err != nil {
		return components, ErrGenerateComponents(err)
	}
	components, err = pkg.GenerateComponents()
	if err != nil {
		return components, ErrGenerateComponents(err)
	}
	return components, nil
}
