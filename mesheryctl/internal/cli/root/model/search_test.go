package model

import (
	goerrors "errors"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestSearchModel(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	querySearch := "model-test"
	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "given no argument provided when running mesheryctl model search then an error message is displayed",
			Args:             []string{"search"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(goerrors.New(errSearchModelName)),
		},
		{
			Name:             "given a valid argument provided when running mesheryctl model search valid-name then the output is displayed",
			Args:             []string{"search", "model-test", "--page", "1", "--pagesize", "10"},
			URL:              fmt.Sprintf("/%s?search=%s&page=0&pagesize=10", modelsApiPath, url.QueryEscape(querySearch)),
			Fixture:          "search.model.api.response.golden",
			ExpectedResponse: "search.model.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given a quoted multi-word argument provided when running mesheryctl model search then the whole phrase is searched",
			Args:             []string{"search", "model test", "--page", "1", "--pagesize", "10"},
			URL:              fmt.Sprintf("/%s?search=%s&page=0&pagesize=10", modelsApiPath, url.QueryEscape("model test")),
			Fixture:          "search.model.api.response.golden",
			ExpectedResponse: "search.model.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given two arguments provided when running mesheryctl model search then an error message is displayed",
			Args:             []string{"search", "aws", "ec2"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(goerrors.New(errInvalidArg)),
		},
		{
			Name:             "given three arguments provided when running mesheryctl model search then an error message is displayed",
			Args:             []string{"search", "aws", "ec2", "lambda"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(goerrors.New(errInvalidArg)),
		},
		{
			Name:             "given an extra argument after a flag when running mesheryctl model search then an error message is displayed",
			Args:             []string{"search", "aws", "--pagesize", "5", "ec2"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(goerrors.New(errInvalidArg)),
		},
	}

	mesheryctlflags.InitValidators(ModelCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, ModelCmd, tests, currDir, "model")
}
