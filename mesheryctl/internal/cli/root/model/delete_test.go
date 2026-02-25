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
			Name:             "given no model name or ID provided when running model delete then ErrInvalidArgument is returned",
			Args:             []string{"delete"},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", modelsApiPath),
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(errors.New(errDeleteInvalidArg)),
		},
		{
			Name:             "given a valid model-id when running model delete then model is deleted successfully",
			Args:             []string{"delete", modelId},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s", modelsApiPath, modelId),
			Fixture:          "delete.model.response.golden",
			ExpectedResponse: "delete.model.success.golden",
			ExpectError:      false,
		},
		{
			Name:             "given a non-existent model name when running model delete then ErrModelNotFound is returned",
			Args:             []string{"delete", "nonexistent-model"},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s?page=0&pagesize=10&search=nonexistent-model", modelsApiPath),
			Fixture:          "delete.model.empty.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrNotFound(fmt.Errorf("no results for %s\n", "nonexistent-model")),
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, ModelCmd, tests, currDir, "model")
}
