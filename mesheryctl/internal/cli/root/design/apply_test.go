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
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy?contexts=all",
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
					URL:          testContext.BaseURL + "/api/pattern/deploy?contexts=all",
					Response:     "apply.designDeploy.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "given valid file and single --contexts when design apply then target context is passed",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "--contexts", "k8s-cluster-1"},
			ExpectedResponse: "apply.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy?contexts=k8s-cluster-1",
					Response:     "apply.designDeploy.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "given valid file and multiple --contexts when design apply then repeated context query params are passed",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "-c", "k8s-cluster-1,k8s-cluster-2"},
			ExpectedResponse: "apply.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy?contexts=k8s-cluster-1&contexts=k8s-cluster-2",
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
					URL:          testContext.BaseURL + "/api/pattern",
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
			Name:             "given invalid combination of all and specific context ID when design apply then return error",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "--contexts", "all,k8s-cluster-1"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    fmt.Errorf("cannot combine 'all' with specific context IDs"),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}
