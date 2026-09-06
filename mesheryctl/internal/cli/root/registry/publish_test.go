package registry

import (
	"os"
	"path/filepath"
	"testing"

	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestMesherySystem(t *testing.T) {
	modelName := "ExampleModel"
	modelRegistrant := "ExampleRegistrant"
	modelNameFormatted := "examplemodel" // FormatName("ExampleModel")

	tests := []struct {
		name        string
		models      []meshkitRegistryUtils.ModelCSV
		components  map[string]map[string][]meshkitRegistryUtils.ComponentCSV
		preSetup    func(t *testing.T, rootDir, subDir string)
		wantErr     bool
		checkOutput func(t *testing.T, rootDir, subDir string)
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
			name: "GenerateIcons error when model path is a file",
			models: []meshkitRegistryUtils.ModelCSV{
				{Model: modelName, Registrant: modelRegistrant},
			},
			components: map[string]map[string][]meshkitRegistryUtils.ComponentCSV{
				modelRegistrant: {},
			},
			preSetup: func(t *testing.T, rootDir, subDir string) {
				blockPath := filepath.Join(rootDir, modelNameFormatted)
				err := os.WriteFile(blockPath, []byte("block"), 0644)
				require.NoError(t, err)
			},
			wantErr: true,
		},
		{
			name: "WriteModelDefToFileSystem error when model def path is a file",
			models: []meshkitRegistryUtils.ModelCSV{
				{Model: modelName, Registrant: modelRegistrant, SVGColor: "<svg></svg>", SVGWhite: "<svg></svg>"},
			},
			components: map[string]map[string][]meshkitRegistryUtils.ComponentCSV{
				modelRegistrant: {},
			},
			preSetup: func(t *testing.T, rootDir, subDir string) {
				blockPath := filepath.Join(subDir, modelName) // modelDef.Name == model.Model
				err := os.WriteFile(blockPath, []byte("block"), 0644)
				require.NoError(t, err)
			},
			wantErr: true,
		},
		{
			name: "valid model with output verification",
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
							SVGColor:  "<svg></svg>",
							SVGWhite:  "<svg></svg>",
						},
					},
				},
			},
			wantErr: false,
			checkOutput: func(t *testing.T, rootDir, subDir string) {
				// Model SVG icons: rootDir/<formattedModel>/icons/{color,white}/<formattedModel>-{color,white}.svg
				modelColorIcon := filepath.Join(rootDir, modelNameFormatted, "icons", "color", modelNameFormatted+"-color.svg")
				modelWhiteIcon := filepath.Join(rootDir, modelNameFormatted, "icons", "white", modelNameFormatted+"-white.svg")
				assert.FileExists(t, modelColorIcon, "model color icon should exist")
				assert.FileExists(t, modelWhiteIcon, "model white icon should exist")

				// Component SVG icons: rootDir/<formattedModel>/components/<compName>/icons/{color,white}/<compName>-{color,white}.svg
				compName := "test-component" // FormatName(FormatToReadableString("TestComponent"))
				compColorIcon := filepath.Join(rootDir, modelNameFormatted, "components", compName, "icons", "color", compName+"-color.svg")
				compWhiteIcon := filepath.Join(rootDir, modelNameFormatted, "components", compName, "icons", "white", compName+"-white.svg")
				assert.FileExists(t, compColorIcon, "component color icon should exist")
				assert.FileExists(t, compWhiteIcon, "component white icon should exist")

				// Model definition JSON: subDir/<ModelName>/model.json
				modelJSON := filepath.Join(subDir, modelName, "model.json")
				assert.FileExists(t, modelJSON, "model definition JSON should exist")

				data, err := os.ReadFile(modelJSON)
				require.NoError(t, err)
				assert.Contains(t, string(data), `"name": "`+modelName+`"`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.SetupMeshkitLoggerTesting(t, false)

			rootDir := t.TempDir()
			subDir := filepath.Join(rootDir, "sub")
			err := os.MkdirAll(subDir, 0755)
			require.NoError(t, err)

			oldCWD, err := os.Getwd()
			require.NoError(t, err)
			err = os.Chdir(subDir)
			require.NoError(t, err)

			if tt.preSetup != nil {
				tt.preSetup(t, rootDir, subDir)
			}

			oldModels := models
			oldComponents := components
			oldModelsOutputPath := modelsOutputPath
			oldImgsOutputPath := imgsOutputPath

			models = tt.models
			components = tt.components
			modelsOutputPath = "."
			imgsOutputPath = "."

			t.Cleanup(func() {
				os.Chdir(oldCWD)
				models = oldModels
				components = oldComponents
				modelsOutputPath = oldModelsOutputPath
				imgsOutputPath = oldImgsOutputPath
			})

			err = mesherySystem()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkOutput != nil {
				tt.checkOutput(t, rootDir, subDir)
			}
		})
	}
}
