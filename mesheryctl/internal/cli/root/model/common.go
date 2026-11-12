package model

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/v1beta1/model"
)

var modelsApiPath = "api/registry/models"

func promptModelSelection(modelSearchTerm string, apiPath string) (*model.ModelDefinition, error) {
	selectedModel := new(model.ModelDefinition)
	err := display.PromptAsyncPagination(
		display.DisplayDataAsync{
			UrlPath:    apiPath,
			SearchTerm: modelSearchTerm,
		},
		formatLabel,
		func(data *model.MeshModelModelsPage) ([]model.ModelDefinition, int64) {
			if data == nil || data.Models == nil {
				return nil, 0
			}
			total := int64(0)
			if data.TotalCount != nil {
				total = int64(*data.TotalCount)
			}
			return *data.Models, total
		},
		selectedModel,
	)
	if err != nil {
		// If the error is due to no models found, return nil without error
		// This allows the caller to handle the "not found" case gracefully without treating it as an actual error
		if meshkiterrors.GetCode(err) == utils.ErrNotFoundCode {
			return nil, nil
		}
		return nil, err
	}
	return selectedModel, nil
}

func formatLabel(rows []model.ModelDefinition) []string {
	labels := []string{}

	for _, m := range rows {
		name := fmt.Sprintf("%s, version: %s", m.DisplayName, m.Version)
		labels = append(labels, name)
	}
	return labels
}
