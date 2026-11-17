package components

import (
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
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
	withModel.ID = uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))
	withModel.Model = &model.ModelDefinition{Name: "kubernetes"}
	withModel.Component.Version = "apps/v1"

	noModel := component.ComponentDefinition{DisplayName: "Orphan"}
	noModel.ID = uuid.Must(uuid.FromString("22222222-2222-2222-2222-222222222222"))

	unnamedModel := component.ComponentDefinition{DisplayName: "Unnamed"}
	unnamedModel.ID = uuid.Must(uuid.FromString("33333333-3333-3333-3333-333333333333"))
	unnamedModel.Model = &model.ModelDefinition{}

	response := &models.MeshmodelComponentsAPIResponse{
		Components: []component.ComponentDefinition{withModel, noModel, unnamedModel},
		TotalCount: 3,
	}

	rows, total := generateComponentDataToDisplay(response)

	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	want := [][]string{
		{"11111111-1111-1111-1111-111111111111", "Deployment", "kubernetes", "apps/v1"},
		{"22222222-2222-2222-2222-222222222222", "Orphan", "N/A", "N/A"},
		{"33333333-3333-3333-3333-333333333333", "Unnamed", "N/A", "N/A"},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Errorf("rows = %v\nwant %v", rows, want)
	}
}
