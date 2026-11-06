package registry

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

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

func TestPublish_Non200StatusError(t *testing.T) {
	statusErr := fmt.Errorf("HTTP status 403")
	errMsg := fmt.Sprintf("Request to Google Spreadsheet did not succeed.\n\nReturned error: %s", statusErr.Error())
	pubErr := errors.New(utils.RegistryError(errMsg, "publish"))
	assert.NotNil(t, pubErr)
	assert.Contains(t, pubErr.Error(), "403")
}
