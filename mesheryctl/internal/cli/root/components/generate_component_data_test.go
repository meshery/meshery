package components

import (
	"reflect"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta3/component"
)

// ComponentDefinition.Model is a pointer, and the API returns "model": null for
// a component that is not tied to a model. Reading the model's name without
// checking crashed mesheryctl component list with a nil pointer dereference;
// it should show N/A, as an unnamed model already does.
func TestGenerateComponentDataToDisplay_ModelColumn(t *testing.T) {
	withModel := component.ComponentDefinition{DisplayName: "Deployment"}
	withModel.Model = &model.ModelDefinition{Name: "kubernetes"}
	withModel.Component.Version = "apps/v1"

	noModel := component.ComponentDefinition{DisplayName: "Orphan"}

	unnamedModel := component.ComponentDefinition{DisplayName: "Unnamed"}
	unnamedModel.Model = &model.ModelDefinition{}

	response := &models.MeshmodelComponentsAPIResponse{
		Components: []component.ComponentDefinition{withModel, noModel, unnamedModel},
		TotalCount: 3,
	}

	rows, total := generateComponentDataToDisplay(response)

	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	zeroID := withModel.ID.String()
	want := [][]string{
		{zeroID, "Deployment", "kubernetes", "apps/v1"},
		{zeroID, "Orphan", "N/A", "N/A"},
		{zeroID, "Unnamed", "N/A", "N/A"},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Errorf("rows = %v\nwant %v", rows, want)
	}
}
