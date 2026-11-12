package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestPublishCmd_PreRunEArgCount(t *testing.T) {
	err := publishCmd.PreRunE(publishCmd, []string{"only-one-arg"})
	assert.Error(t, err, "PreRunE should return error if fewer than 5 arguments are passed")
}

func TestPublishCmd_InvalidCredentials(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	err := publishCmd.RunE(publishCmd, []string{"website", "invalid-token", "some-id", "models", "imgs"})
	assert.Error(t, err, "RunE should return error when credentials are invalid")
}

func TestWebsiteSystem_WriteFailure(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	origOutputFormat := outputFormat
	origModelsOutputPath := modelsOutputPath
	defer func() {
		outputFormat = origOutputFormat
		modelsOutputPath = origModelsOutputPath
	}()

	outputFormat = "js"
	// An invalid path that cannot be written to
	modelsOutputPath = "/nonexistent/invalid/dir/that/cannot/exist"

	err := websiteSystem()
	assert.Error(t, err, "websiteSystem must return an error when data.js cannot be written")
}

func TestSpreadsheetRequestError_Non200Status(t *testing.T) {
	resp := &sheets.Spreadsheet{}
	resp.HTTPStatusCode = 403

	errMsg := spreadsheetRequestError(resp, nil)
	assert.Contains(t, errMsg, "did not succeed")
	assert.Contains(t, errMsg, "403")
}

func TestSpreadsheetRequestError_EmptyResponse(t *testing.T) {
	errMsg := spreadsheetRequestError(nil, nil)
	assert.Contains(t, errMsg, "did not succeed")
	assert.Contains(t, errMsg, "empty response")
}

func TestMesherySystem_Unsupported(t *testing.T) {
	err := mesherySystem()
	assert.Error(t, err, "mesherySystem should return an error as publishing to meshery is unsupported")
	assert.Contains(t, err.Error(), "not supported")
}

func TestErrPublish_MesherySystem(t *testing.T) {
	err := mesherySystem()
	pubErr := ErrPublish(err, "meshery")
	assert.Error(t, pubErr)
	assert.Contains(t, pubErr.Error(), "publishing directly to system 'meshery' is not supported")
}
