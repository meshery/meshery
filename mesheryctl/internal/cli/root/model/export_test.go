package model

import (
	"errors"
	"fmt"
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
	utils.TokenFlag = utils.GetToken(t)
	mesheryctlflags.InitValidators(ModelCmd)

	const modelName = "test-model"
	exportedContent := []byte("exported-model-archive-bytes")

	// RunE builds the query string with url.Values.Encode(), which sorts the
	// keys alphabetically. The defaults are: components/relationships enabled,
	// file_type=oci, output_format=yaml, page=1.
	exportURL := fmt.Sprintf("%s/api/meshmodels/export?components=true&file_type=oci&name=%s&output_format=yaml&page=1&relationships=true", testContext.BaseURL, modelName)
	httpmock.RegisterResponder("GET", exportURL,
		httpmock.NewBytesResponder(200, exportedContent))

	outputDir := t.TempDir()
	buf := utils.SetupMeshkitLoggerTesting(t, false)
	ModelCmd.SetArgs([]string{"export", modelName, "--output-location", outputDir})
	ModelCmd.SetOut(buf)

	err := ModelCmd.Execute()
	if err != nil {
		t.Fatalf("expected export to succeed, got error: %v", err)
	}

	// The default output type "oci" produces a ".tar" file named after the model.
	exportedPath := filepath.Join(outputDir, modelName+".tar")
	got, readErr := os.ReadFile(exportedPath)
	if readErr != nil {
		t.Fatalf("expected exported file at %s: %v", exportedPath, readErr)
	}
	assert.Equal(t, exportedContent, got)
	assert.Contains(t, buf.String(), "Exported model to")
}
