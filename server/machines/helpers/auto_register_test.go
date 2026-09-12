package helpers

import (
	"testing"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	"github.com/meshery/schemas/models/core"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta3/component"
)

// nonComponentEntity is an entity.Entity that is not a *component.ComponentDefinition.
type nonComponentEntity struct{}

func (nonComponentEntity) Type() entity.EntityType        { return entity.EntityType("non-component") }
func (nonComponentEntity) GetEntityDetail() string        { return "non-component" }
func (nonComponentEntity) GenerateID() (core.Uuid, error) { return core.Uuid{}, nil }
func (nonComponentEntity) GetID() core.Uuid               { return core.Uuid{} }
func (nonComponentEntity) Create(_ *database.Handler, _ core.Uuid) (core.Uuid, error) {
	return core.Uuid{}, nil
}

func connectionDef(name string) *component.ComponentDefinition {
	return &component.ComponentDefinition{
		DisplayName: core.InputString(name),
		Model:       &modelv1beta1.ModelDefinition{},
	}
}

func TestToConnectionDefinitions(t *testing.T) {
	tests := []struct {
		Name     string
		Entities []entity.Entity
		WantLen  int
	}{
		{
			Name:     "given no entities when converting then no definitions are returned",
			Entities: nil,
			WantLen:  0,
		},
		{
			Name: "given only connection definitions when converting then no phantom entries are returned",
			Entities: []entity.Entity{
				connectionDef("grafana"),
				connectionDef("prometheus"),
			},
			WantLen: 2,
		},
		{
			Name: "given a non component entity when converting then it is skipped",
			Entities: []entity.Entity{
				connectionDef("grafana"),
				nonComponentEntity{},
			},
			WantLen: 1,
		},
		{
			Name: "given a definition without a model when converting then it is skipped",
			Entities: []entity.Entity{
				connectionDef("grafana"),
				&component.ComponentDefinition{},
			},
			WantLen: 1,
		},
		{
			Name: "given a nil definition when converting then it is skipped",
			Entities: []entity.Entity{
				connectionDef("grafana"),
				// A typed nil pointer in an entity.Entity passes the type assertion with
				// ok == true, so the asserted pointer itself has to be checked for nil.
				(*component.ComponentDefinition)(nil),
			},
			WantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := toConnectionDefinitions(tt.Entities)

			if len(got) != tt.WantLen {
				t.Fatalf("got %d connection definitions, want %d", len(got), tt.WantLen)
			}

			// getConnectionPayload dereferences Model, so any definition returned here with a
			// nil Model is the server crash reported in #20729.
			for i, def := range got {
				if def.Model == nil {
					t.Errorf("connection definition at index %d has a nil Model, getConnectionPayload would panic on it", i)
				}
			}
		})
	}
}
