package model

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestHasCSVs(t *testing.T) {
	tests := []struct {
		name           string
		dirPath        string
		expectedResult bool
	}{
		{
			name:           "directory with CSVs",
			dirPath:        "./fixtures/with_csvs",
			expectedResult: true,
		},
		{
			name:           "directory without CSVs",
			dirPath:        "./fixtures/without_csvs",
			expectedResult: false,
		},
		{
			name:           "non-existent directory",
			dirPath:        "./fixtures/invalid_path",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := hasCSVs(tc.dirPath)

			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestImportModelReturnsErrorForURLImportFailure(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}

	currDir := filepath.Dir(filename)

	tests := []utils.MesheryCommandTest{
		{
			Name:           "given a URL import failure when importing model then return the error",
			Args:           []string{"import", "https://example.com/model"},
			URL:            "/api/meshmodels/register",
			HttpMethod:     "POST",
			HttpStatusCode: 500,
			Fixture:        "model.import.error.api.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrMesheryServerInternalError(errors.New("internal server error\n")),
			IsOutputGolden: false,
		},
	}

	mesheryctlflags.InitValidators(ModelCmd)
	utils.InvokeMesheryctlTestCommand(t, update, ModelCmd, tests, currDir, "model")
}
