package model

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

func TestDeleteModel(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	modelId := "d56fb25b-f92c-4cd6-821b-2cfd6bb87259"

	tests := []utils.MesheryCommandTest{
		{
			Name:             "Delete model without arguments",
			Args:             []string{"delete"},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", modelsApiPath),
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("[ model-id ] is required\n\nUsage: mesheryctl model delete [model-id]\nRun 'mesheryctl model delete --help' to see detailed help message")),
		},
		{
			Name:             "Delete model with invalid UUID",
			Args:             []string{"delete", "not-a-uuid"},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", modelsApiPath),
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidUUID(fmt.Errorf("invalid model ID: %q", "not-a-uuid")),
		},
		{
			Name:             "Delete model successfully",
			Args:             []string{"delete", modelId},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s", modelsApiPath, modelId),
			Fixture:          "delete.model.response.golden",
			ExpectedResponse: "delete.model.success.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, ModelCmd, tests, currDir, "model")
}
