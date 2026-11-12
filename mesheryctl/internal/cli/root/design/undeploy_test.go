package design

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestUndeployCmd(t *testing.T) {
	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	_, expectedFileReadErr := os.ReadFile(invalidFilePath)

	// test scenrios for fetching data
	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:             "given valid file when design undeploy then design is undeployed",
			Args:             []string{"undeploy", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "undeploy.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          "/api/pattern/deploy",
					Response:     "undeploy.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given valid relative file when design undeploy then design is undeployed",
			Args:             []string{"undeploy", "-f", filepath.Join("fixtures", "sampleDesign.golden")},
			ExpectedResponse: "undeploy.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          "/api/pattern/deploy",
					Response:     "undeploy.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given invalid file path when design undeploy then error is thrown",
			Args:             []string{"undeploy", "-f", invalidFilePath},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFileRead(expectedFileReadErr),
		},
		{
			Name:             "given URL when design undeploy then error is thrown",
			Args:             []string{"undeploy", "-f", "https://raw.githubusercontent.com/meshery/meshery/master/samples/test.yaml"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("URLs are not currently supported")),
		},
		{
			Name:             "given no args when design undeploy then usage error is returned",
			Args:             []string{"undeploy"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: ErrUndeployDesign(
				fmt.Errorf("provide either a design ID or -f [filepath]"),
			),
		},
		{
			Name:             "given empty file flag when design undeploy then error is thrown",
			Args:             []string{"undeploy", "-f", ""},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: ErrUndeployDesign(
				fmt.Errorf("Usage: mesheryctl design undeploy -f [filepath]"),
			),
		},
		{
			Name:             "given nonexistent design provided when design undeploy then error is thrown",
			Args:             []string{"undeploy", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          "/api/pattern/deploy",
					Response:     "undeploy.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          fmt.Sprintf("%s/api/pattern", testContext.BaseURL),
					Response:     "undeploy.empty.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound(filepath.Join(fixturesDir, "sampleDesign.golden")),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}
