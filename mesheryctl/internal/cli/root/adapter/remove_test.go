package adapter

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestRemoveMesh(t *testing.T) {
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "Test Remove Linkerd ",
			Args: []string{"remove", "linkerd"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/system/sync",
					Response:     "sync.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/adapter/operation",
					Response:     "deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectedContains: []string{"Verifying prerequisites"},
			ExpectError:      false,
		},
		{
			Name: "Test Remove Linkerd with Namespace",
			Args: []string{"remove", "linkerd", "--namespace", "linkerd-ns"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/system/sync",
					Response:     "sync.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/adapter/operation",
					Response:     "deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectedContains: []string{"Verifying prerequisites"},
			ExpectError:      false,
		},
		{
			Name: "Test Remove Cilium ",
			Args: []string{"remove", "cilium", "service", "mesh"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/system/sync",
					Response:     "sync.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/adapter/operation",
					Response:     "deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectedContains: []string{"Verifying prerequisites"},
			ExpectError:      false,
		},
	}

	utils.RunMesheryctlMultiURLTests(t, update, AdapterCmd, tests, currDir, "adapter", func() {})
}
