package registry

import (
	"fmt"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
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

	// Point to an invalid spreadsheet URL so downloading/parsing the components sheet fails
	spreadsheeetID = "invalid-spreadsheet-id"
	sheetGID = 0
	componentCSVFilePath = ""
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
