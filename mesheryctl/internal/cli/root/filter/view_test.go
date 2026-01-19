package filter

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestViewCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	testContext := utils.NewTestHelper(t)
	base := testContext.BaseURL

	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "Fetch Filter View",
			Args: []string{"view", "KumaTest"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: fmt.Sprintf("%s/api/filter?page_size=10000", base), Response: "view.filter.api.response.golden", ResponseCode: 200},
				{Method: "GET", URL: fmt.Sprintf("%s/api/filter?search=KumaTest", base), Response: "view.filter.api.response.golden", ResponseCode: 200},
			},
			ExpectedContains: []string{"name: KumaTest", "id: 957fbc9b-a655-4892-823d-375102a9587c"},
			ExpectError:      false,
		},
		{
			Name: "Fetch Kuma Filter View with ID",
			Args: []string{"view", "957fbc9b-a655-4892-823d-375102a9587c"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: fmt.Sprintf("%s/api/filter?page_size=10000", base), Response: "view.filter.api.response.golden", ResponseCode: 200},
				{Method: "GET", URL: fmt.Sprintf("%s/api/filter/957fbc9b-a655-4892-823d-375102a9587c", base), Response: "view.id.filter.api.response.golden", ResponseCode: 200},
			},
			ExpectedContains: []string{"name: KumaTest", "id: 957fbc9b-a655-4892-823d-375102a9587c"},
			ExpectError:      false,
		},
		{
			Name: "Fetch Filter View for non existing filter",
			Args: []string{"view", "xyz"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: fmt.Sprintf("%s/api/filter?page_size=10000", base), Response: "view.filter.api.response.golden", ResponseCode: 200},
				{Method: "GET", URL: fmt.Sprintf("%s/api/filter?search=xyz", base), Response: "view.nonexisting.filter.api.response.golden", ResponseCode: 200},
			},
			ExpectedContains: []string{"filter with name: xyz not found"},
			ExpectError:      false,
		},
	}

	utils.RunMesheryctlMultiURLTests(t, update, FilterCmd, tests, currDir, "filter", func() {})
}
