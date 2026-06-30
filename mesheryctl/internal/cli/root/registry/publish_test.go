package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestMesherySystem(t *testing.T) {
	modelName := "ExampleModel"
	modelRegistrant := "ExampleRegistrant"

	tests := []struct {
		name       string
		models     []meshkitRegistryUtils.ModelCSV
		components map[string]map[string][]meshkitRegistryUtils.ComponentCSV
		wantErr    bool
	}{
		{
			name:       "empty models list",
			models:     []meshkitRegistryUtils.ModelCSV{},
			components: map[string]map[string][]meshkitRegistryUtils.ComponentCSV{},
			wantErr:    false,
		},
		{
			name: "model registrant not in components map",
			models: []meshkitRegistryUtils.ModelCSV{
				{Model: modelName, Registrant: modelRegistrant},
			},
			components: map[string]map[string][]meshkitRegistryUtils.ComponentCSV{},
			wantErr:    false,
		},
		{
			name: "registrant exists but model key missing",
			models: []meshkitRegistryUtils.ModelCSV{
				{Model: modelName, Registrant: modelRegistrant},
			},
			components: map[string]map[string][]meshkitRegistryUtils.ComponentCSV{
				modelRegistrant: {},
			},
			wantErr: false,
		},
		{
			name: "valid model with components",
			models: []meshkitRegistryUtils.ModelCSV{
				{
					Model:             modelName,
					Registrant:        modelRegistrant,
					ModelDisplayName:  "Example Model",
					Description:       "An example model for testing",
					Category:          "Test",
					SVGColor:          "<svg></svg>",
					SVGWhite:          "<svg></svg>",
					PublishToRegistry: "true",
				},
			},
			components: map[string]map[string][]meshkitRegistryUtils.ComponentCSV{
				modelRegistrant: {
					modelName: {
						{
							Component: "TestComponent",
							Model:     modelName,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple models with one missing registrant",
			models: []meshkitRegistryUtils.ModelCSV{
				{Model: "ModelA", Registrant: "RegA"},
				{Model: "ModelB", Registrant: "RegB"},
			},
			components: map[string]map[string][]meshkitRegistryUtils.ComponentCSV{
				"RegA": {
					"ModelA": {
						{Component: "CompA", Model: "ModelA"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.SetupMeshkitLoggerTesting(t, false)

			modelsDir := t.TempDir()
			imgsDir := t.TempDir()

			oldModels := models
			oldComponents := components
			oldModelsOutputPath := modelsOutputPath
			oldImgsOutputPath := imgsOutputPath

			models = tt.models
			components = tt.components
			modelsOutputPath = modelsDir
			imgsOutputPath = imgsDir

			t.Cleanup(func() {
				models = oldModels
				components = oldComponents
				modelsOutputPath = oldModelsOutputPath
				imgsOutputPath = oldImgsOutputPath
			})

			err := mesherySystem()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
