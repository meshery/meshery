package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	comp "github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvokeCompUpdate_InputPaths(t *testing.T) {
	origCompCSVPath, origModelLocation, origModelName := componentCSVFilePath, modelLocation, modelName
	origLogFile, origLog, origRegistryLog := logFile, utils.Log, meshkitRegistryUtils.Log
	origTotalModels, origTotalComponents := totalAggregateModel, totalAggregateComponents
	t.Cleanup(func() {
		componentCSVFilePath, modelLocation, modelName = origCompCSVPath, origModelLocation, origModelName
		logFile, utils.Log, meshkitRegistryUtils.Log = origLogFile, origLog, origRegistryLog
		totalAggregateModel, totalAggregateComponents = origTotalModels, origTotalComponents
	})

	utils.SetupMeshkitLoggerTesting(t, false)
	meshkitRegistryUtils.Log = utils.Log
	logFile = nil
	modelName = "demo"

	root := t.TempDir()
	modelsDir := filepath.Join(root, "models")
	subDir := filepath.Join(root, "mesheryctl")
	modelsWithoutComponent := filepath.Join(root, "models-without-component")
	componentPath := filepath.Join(modelsDir, modelName, "1.0.0", "v1.0.0", "components", "Thing.json")
	require.NoError(t, os.MkdirAll(filepath.Dir(componentPath), 0o755))
	require.NoError(t, os.MkdirAll(subDir, 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(modelsWithoutComponent, modelName, "1.0.0", "v1.0.0"), 0o755))
	componentCSVFilePath = filepath.Join(root, "components.csv")
	csv := "ignored,ignored,ignored,ignored\nregistrant,model,component,description\n" +
		"meshery,demo,Thing,after\n"
	require.NoError(t, os.WriteFile(componentCSVFilePath, []byte(csv), 0o644))

	tests := []struct {
		name            string
		cwd             string
		input           string
		wantErr         bool
		wantDescription string
	}{
		{name: "relative", cwd: root, input: "models", wantDescription: "after"},
		{name: "absolute", cwd: root, input: modelsDir, wantDescription: "after"},
		{name: "relative from subdirectory", cwd: subDir, input: filepath.Join("..", "models"), wantDescription: "after"},
		{name: "valid directory without the model", cwd: root, input: "mesheryctl", wantDescription: "before"},
		{name: "valid model without the component", cwd: root, input: "models-without-component", wantDescription: "before"},
		{name: "missing relative directory", cwd: root, input: "missing-models", wantErr: true, wantDescription: "before"},
		{name: "missing absolute directory", cwd: root, input: filepath.Join(root, "missing-models"), wantErr: true, wantDescription: "before"},
		{name: "regular file", cwd: root, input: componentCSVFilePath, wantErr: true, wantDescription: "before"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chdir(t, tt.cwd)
			original := []byte(`{"description":"before","metadata":{},"component":{"kind":"Thing","schema":"{}"}}`)
			require.NoError(t, os.WriteFile(componentPath, original, 0o644))
			modelLocation = tt.input

			err := InvokeCompUpdate()
			if tt.wantErr {
				require.ErrorContains(t, err, tt.input)
			} else {
				require.NoError(t, err)
			}

			content, err := os.ReadFile(componentPath)
			require.NoError(t, err)
			var updated comp.ComponentDefinition
			require.NoError(t, json.Unmarshal(content, &updated))
			assert.Equal(t, tt.wantDescription, updated.Description)
		})
	}
}

func TestInvokeCompUpdate_ComponentsSheetError(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	// Backup original state
	origSpreadsheetID := spreadsheeetID
	origSheetGID := sheetGID
	origCompCSVPath := componentCSVFilePath
	origModelName := modelName
	origModelLocation := modelLocation
	origLogFile := logFile

	defer func() {
		spreadsheeetID = origSpreadsheetID
		sheetGID = origSheetGID
		componentCSVFilePath = origCompCSVPath
		modelName = origModelName
		modelLocation = origModelLocation
		logFile = origLogFile
	}()

	// Use a missing local file so the components sheet read fails without
	// making a network request to Google Sheets.
	spreadsheeetID = "invalid-spreadsheet-id"
	sheetGID = 0
	componentCSVFilePath = filepath.Join(t.TempDir(), "missing-components.csv")
	modelName = ""
	modelLocation = t.TempDir()
	logFile = nil

	err := InvokeCompUpdate()
	assert.Error(t, err, "InvokeCompUpdate must return an error when the Components sheet cannot be read")
}

func TestUpdateCmd_InvalidCredentials(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	origCred := spreadsheeetCred
	origID := spreadsheeetID
	defer func() {
		spreadsheeetCred = origCred
		spreadsheeetID = origID
	}()

	spreadsheeetCred = "invalid-token"
	spreadsheeetID = "some-id"

	err := updateCmd.RunE(updateCmd, []string{})
	assert.Error(t, err, "updateCmd.RunE must return an error when credentials are invalid")
}

func TestErrUpdateRegistry_WithNon200Status(t *testing.T) {
	err := fmt.Errorf("request to Google Spreadsheet returned HTTP status 404")
	registryErr := ErrUpdateRegistry(err, "./models")
	assert.NotNil(t, registryErr)
	assert.Contains(t, registryErr.Error(), "404")
}
