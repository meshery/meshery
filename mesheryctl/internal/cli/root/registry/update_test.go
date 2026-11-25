package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	modelLocation = "./models"
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

func TestInvokeCompUpdate_Summary(t *testing.T) {
	origCompCSVPath := componentCSVFilePath
	origModelName := modelName
	origModelLocation := modelLocation
	origLogFile := logFile
	origRegistryLog := meshkitRegistryUtils.Log
	t.Cleanup(func() {
		componentCSVFilePath = origCompCSVPath
		modelName = origModelName
		modelLocation = origModelLocation
		logFile = origLogFile
		meshkitRegistryUtils.Log = origRegistryLog
	})

	out := utils.SetupMeshkitLoggerTesting(t, false)
	meshkitRegistryUtils.Log = utils.Log
	logFile = nil
	modelName = ""

	root := t.TempDir()
	csv := "ignored,ignored,ignored,ignored\nregistrant,model,component,description\n"
	for _, c := range []string{"BrowserSet", "Config", "DeviceSet", "License", "Quota"} {
		csv += "meshery,aerokube," + c + ",repro\n"
		for _, v := range []string{"2.8.1", "2.8.2", "2.8.3"} {
			dir := filepath.Join(root, "models", "aerokube", v, "v1.0.0", "components")
			require.NoError(t, os.MkdirAll(dir, 0o755))
			require.NoError(t, os.WriteFile(filepath.Join(dir, c+".json"), []byte("{}"), 0o644))
		}
	}

	// "ghost" is listed in the CSV, but its only component has no JSON file.
	ghostDir := filepath.Join(root, "models", "ghost", "1.0.0", "v1.0.0", "components")
	require.NoError(t, os.MkdirAll(ghostDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(ghostDir, "Other.json"), []byte("{}"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(ghostDir, "README.md"), []byte("notes"), 0o644))
	csv += "meshery,ghost,Missing,repro\n"

	// MeshKit 1.0.29 can drop the last row, so end with a blank one.
	csv += ",,,\n"
	componentCSVFilePath = filepath.Join(root, "components.csv")
	require.NoError(t, os.WriteFile(componentCSVFilePath, []byte(csv), 0o644))

	t.Chdir(root)
	modelLocation = "models"
	require.NoError(t, InvokeCompUpdate())

	assert.Contains(t, out.String(), "For 1 models updated 15 components")
	assert.Contains(t, out.String(), "Updated 1 models and 15 components")
	assert.Contains(t, out.String(), "ghost-1.0.0, updated 0 out of 1 components")
	for _, v := range []string{"2.8.1", "2.8.2", "2.8.3"} {
		assert.Contains(t, out.String(), "aerokube-"+v+", updated 5 out of 5 components")
	}
}
