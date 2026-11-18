package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/meshery/meshery/server/helpers"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/generators/artifacthub"
	"github.com/meshery/meshkit/models/meshmodel/registry"

	meshkitmodels "github.com/meshery/meshkit/generators/models"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta3/component"
)

// maxComponentGenerationPayloadSize limits the maximum allowed request body size (4 MB)
// for component generation requests, preventing unbounded memory allocation and denial-of-service
// vulnerabilities via oversized payloads.
const maxComponentGenerationPayloadSize = 4 << 20 // 4MB

type generationPayloadItem struct {
	Name     string `json:"name"`
	Register bool   `json:"register"`
}

type componentGenerationPayload struct {
	Data []generationPayloadItem `json:"data"`
}

type componentGenerationResponseDataItem struct {
	Name       string                          `json:"name"`
	Components []component.ComponentDefinition `json:"components"`
	Errors     []string                        `json:"errors"`
}

// MeshModelGenerationHandler expects the request body to be JSON.
// request body should be of format - {data: [{name: string, register: boolean}]}
// response format - {data: [{name: string, components: [component], errors: [string] }]}
func (h *Handler) MeshModelGenerationHandler(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	// Parse the request body bounded by maxComponentGenerationPayloadSize
	body, err := io.ReadAll(http.MaxBytesReader(rw, r.Body, maxComponentGenerationPayloadSize))
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			if h.log != nil {
				h.log.Error(ErrRequestBody(err))
			}
			writeMeshkitError(rw, ErrRequestBody(err), http.StatusRequestEntityTooLarge)
			return
		}
		if h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	body = bytes.Trim(body, " \t\r\n")
	if len(body) == 0 || bytes.Equal(body, []byte("null")) {
		err := fmt.Errorf("component generation request body is empty or null")
		if h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	// Unmarshal request body
	pld := componentGenerationPayload{}
	err = json.Unmarshal(body, &pld)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	if len(pld.Data) == 0 {
		err := fmt.Errorf("component generation request body contains no data items")
		if h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
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
		if gpi.Register {
			for _, comp := range comps {
				var isModelError bool
				var isRegistranError bool
				utils.WriteSVGsOnFileSystem(&comp)
				host := fmt.Sprintf("%s.artifacthub.meshery", gpi.Name)
				isRegistranError, isModelError, err = h.registryManager.RegisterEntity(registry.RegistrantHostToV1beta3(connection.Connection{
					Kind: artifacthub.ArtifactHub,
					Metadata: map[string]interface{}{
						"name": host,
					},
				}), &comp)
				helpers.HandleError(connection.Connection{
					Kind: artifacthub.ArtifactHub,
				}, &comp, err, isModelError, isRegistranError)
				if err != nil {
					h.log.Error(ErrGenerateComponents(err))
				}

				h.log.Info(comp.DisplayName, " component for ", gpi.Name, " generated")
			}
		}

		responseItem.Components = comps
		response = append(response, responseItem)
	}
	err = helpers.WriteLogsToFiles()
	if err != nil {
		h.log.Error(err)
	}
	// Send response
	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		// Response body has already started streaming via json.Encoder —
		// a partial JSON envelope is on the wire and a fresh error
		// response would corrupt it, so log only.
		h.log.Error(ErrGenerateComponents(err))
		return
	}
}

func generateComponents(pm meshkitmodels.PackageManager) ([]component.ComponentDefinition, error) {
	components := make([]component.ComponentDefinition, 0)
	pkg, err := pm.GetPackage()
	if err != nil {
		return components, ErrGenerateComponents(err)
	}
	components, err = pkg.GenerateComponents("")
	if err != nil {
		return components, ErrGenerateComponents(err)
	}
	return components, nil
}
