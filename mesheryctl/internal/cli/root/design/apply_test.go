package design

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
	"github.com/pkg/errors"
)

// TestApplyCmd verifies the output and errors of design apply for local files.
func TestApplyCmd(t *testing.T) {
	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	// test scenrios for fetching data
	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:             "given valid file when design apply then design is applied",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "apply.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/import",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "apply.designDeploy.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "given valid file when design apply with --skip-save then design is applied without saving",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "--skip-save"},
			ExpectedResponse: "apply.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "apply.designDeploy.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "given invalid file path when design apply then error is thrown",
			Args:             []string{"apply", "-f", invalidFilePath},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFileRead(errors.Errorf(errInvalidPathMsg, invalidFilePath)),
		},
		{
			Name:             "given invalid server response when design apply then error is thrown",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/import",
					Response:     "apply.invalidJSON.response.golden",
					ResponseCode: 200,
				},
			},
			Token:          filepath.Join(fixturesDir, "token.golden"),
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError: func() error {
				// Replicate the exact JSON unmarshal error
				var response []*pattern.MesheryPattern
				innerErr := json.Unmarshal([]byte(`{ "patterns": [ { "id": "123", "name": "incomplete-json"`), &response)

				return utils.ErrUnmarshal(innerErr)
			}(),
		},
		{
			Name:             "given no args and no file flag when design apply then validation error is returned",
			Args:             []string{"apply"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("please provide a design name or specify a file using --file")),
		},
		{
			Name:             "given empty file flag when design apply then error is returned",
			Args:             []string{"apply", "-f", ""},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("file path cannot be empty")),
		},
		{
			Name:             "given both design name and file flag when design apply then error is returned",
			Args:             []string{"apply", "my-design", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("cannot specify both a design name and the --file flag")),
		},
		{
			Name:             "given empty design name when design apply then validation error is returned",
			Args:             []string{"apply", ""},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("please provide a design name or specify a file using --file")),
		},
		{
			Name:             "given whitespace design name when design apply then validation error is returned",
			Args:             []string{"apply", "   "},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("please provide a design name or specify a file using --file")),
		},
		{
			Name:             "given multi-word design name as separate arguments when design apply then design is applied",
			Args:             []string{"apply", "Istio", "Service", "Mesh"},
			ExpectedResponse: "apply.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&search=Istio+Service+Mesh",
					Response:     "list.design.local.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "apply.designDeploy.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "given quoted multi-word design name when design apply then design is applied",
			Args:             []string{"apply", "Istio Service Mesh"},
			ExpectedResponse: "apply.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&search=Istio+Service+Mesh",
					Response:     "list.design.local.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "apply.designDeploy.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "given non-existent design name when design apply then throw error",
			Args:             []string{"apply", "nonexistent-design"},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&search=nonexistent-design",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
			},
			Token:          filepath.Join(fixturesDir, "token.golden"),
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound("nonexistent-design"),
		},
		{
			Name:             "given multi-word design name with spaces when design apply and not found then throw error without raw percent encoding",
			Args:             []string{"apply", "my", "sample", "design"},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&search=my+sample+design",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
			},
			Token:          filepath.Join(fixturesDir, "token.golden"),
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound("my sample design"),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}
