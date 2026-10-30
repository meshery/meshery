package design

import (
	"fmt"
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
					URL:          "/api/pattern/deploy?contexts=all",
					Response:     "undeploy.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given valid file and single --contexts when design undeploy then target context is passed",
			Args:             []string{"undeploy", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "--contexts", "k8s-cluster-1"},
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
					URL:          "/api/pattern/deploy?contexts=k8s-cluster-1",
					Response:     "undeploy.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given valid file and multiple --contexts when design undeploy then repeated context query params are passed",
			Args:             []string{"undeploy", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "-c", "k8s-cluster-1,k8s-cluster-2"},
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
					URL:          "/api/pattern/deploy?contexts=k8s-cluster-1&contexts=k8s-cluster-2",
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
			ExpectedError:    utils.ErrFileRead(fmt.Errorf("open %s: no such file or directory", invalidFilePath)),
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
		{
			Name:             "given invalid combination of all and specific context ID when design undeploy then return error",
			Args:             []string{"undeploy", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "--contexts", "all,k8s-cluster-1"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    fmt.Errorf("cannot combine 'all' with specific context IDs"),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}
