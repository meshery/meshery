package model

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestExportModel(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	const modelName = "model-test-0"
	const exportUsage = "Usage: mesheryctl model export [model-name]\nRun 'mesheryctl model export --help' to see detailed help message"

	tests := []utils.MesheryCommandTest{
		{
			Name:          "given no argument when model export then throw error",
			Args:          []string{"export"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(errors.New("Please provide a model name. " + exportUsage)),
		},
		{
			Name:          "given an invalid output format when model export then throw error",
			Args:          []string{"export", modelName, "--output-format", "invalid-format"},
			ExpectError:   true,
			ExpectedError: utils.ErrFlagsInvalid(errors.New("Invalid value for --output-format 'invalid-format': valid values are json yaml")),
		},
		{
			Name:          "given an invalid output type when model export then throw error",
			Args:          []string{"export", modelName, "--output-type", "invalid-type"},
			ExpectError:   true,
			ExpectedError: utils.ErrFlagsInvalid(errors.New("Invalid value for --output-type 'invalid-type': valid values are oci tar")),
		},
		{
			Name:          "given an invalid version when model export then throw error",
			Args:          []string{"export", modelName, "--version", "1.0.0"},
			ExpectError:   true,
			ExpectedError: utils.ErrFlagsInvalid(errors.New("Invalid value for --version '1.0.0': version must be in format vX.X.X")),
		},
	}

	mesheryctlflags.InitValidators(ModelCmd)
	utils.InvokeMesheryctlTestCommand(t, update, ModelCmd, tests, currDir, "model")
}

func TestExportModelToFile(t *testing.T) {
	defer utils.ResetCommandFlags(ModelCmd, t)
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	oldToken := utils.TokenFlag
	defer func() { utils.TokenFlag = oldToken }()
	utils.TokenFlag = utils.GetToken(t)
	mesheryctlflags.InitValidators(ModelCmd)

	const modelName = "test-model"
	exportedContent := []byte("exported-model-archive-bytes")

	// Must match the URL export.go requests: httpmock answers only an exact match.
	queryParams := url.Values{}
	queryParams.Set("name", modelName)
	queryParams.Set("outputFormat", "yaml")
	queryParams.Set("fileType", "oci")
	queryParams.Set("components", "true")
	queryParams.Set("relationships", "true")
	queryParams.Set("page", "1")

	exportURL := fmt.Sprintf("%s/api/registry/export?%s", testContext.BaseURL, queryParams.Encode())
	httpmock.RegisterResponder("GET", exportURL,
		httpmock.NewBytesResponder(200, exportedContent))

	outputDir := t.TempDir()
	buf := utils.SetupMeshkitLoggerTesting(t, false)
	ModelCmd.SetArgs([]string{"export", modelName, "--output-location", outputDir})
	ModelCmd.SetOut(buf)

	err := ModelCmd.Execute()
	assert.NoError(t, err)

	// The default output type "oci" produces a ".tar" file named after the model.
	got, err := os.ReadFile(filepath.Join(outputDir, modelName+".tar"))
	assert.NoError(t, err)
	assert.Equal(t, exportedContent, got)
	assert.Contains(t, buf.String(), "Exported model to")
}
