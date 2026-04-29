package model

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/v1beta1/model"
)

var modelsApiPath = "api/meshmodels/models"

func promptModelSelection(modelSearchTerm string, apiPath string) (*model.ModelDefinition, error) {
	selectedModel := new(model.ModelDefinition)
	err := display.PromptAsyncPagination(
		display.DisplayDataAsync{
			UrlPath:    apiPath,
			SearchTerm: modelSearchTerm,
		},
		formatLabel,
		func(data *models.MeshmodelsAPIResponse) ([]model.ModelDefinition, int64) {
			return data.Models, data.TotalCount
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
