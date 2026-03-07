package design

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestOffboardCmd(t *testing.T) {
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
			Name:             "Offboard Application",
			Args:             []string{"offboard", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "offboard.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "offboard.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "Offboard design with invalid file path",
			Args:             []string{"offboard", "-f", invalidFilePath},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFileRead(fmt.Errorf("open %s: no such file or directory", invalidFilePath)),
		},
		{
			Name:             "Offboard design not found",
			Args:             []string{"offboard", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "offboard.empty.response.golden",
					ResponseCode: 200,
				},
			},
			Token:          filepath.Join(fixturesDir, "token.golden"),
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound(filepath.Join(fixturesDir, "sampleDesign.golden")),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}
