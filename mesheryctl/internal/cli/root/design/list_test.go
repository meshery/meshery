package design

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
)

func TestDesignListCmd(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenrios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "given no arguments provided when design list then design list is displayed",
			Args:             []string{"list", "--page", "1"},
			ExpectedResponse: "list.design.output.golden",
			Fixture:          "list.design.api.response.golden",
			URL:              "/api/pattern?page=0&pagesize=10",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, DesignCmd, tests, currDir, "design")
}

// Test_processDesignData_DropsOnlyTheExtensionFromDesignNames pins the NAME
// column against `strings.Trim`, whose second argument is a cutset and not a
// suffix: every character of the extension was also stripped off the front of
// the name, so `mesh.yaml` listed as `esh` and `my-app.yaml` as `-app`.
func Test_processDesignData_DropsOnlyTheExtensionFromDesignNames(t *testing.T) {
	previousProvider := provider
	provider = models.LocalProviderName
	t.Cleanup(func() { provider = previousProvider })

	designID := uuid.FromStringOrNil("76430bfa-e725-4444-a79f-76112bef1a78")
	timestamp := time.Date(2021, 8, 1, 8, 28, 50, 0, time.UTC)

	tests := []struct {
		name       string
		designName string
		expected   string
	}{
		{
			name:       "given a name whose first character recurs in its extension when design list then only the extension is dropped",
			designName: "mesh.yaml",
			expected:   "mesh",
		},
		{
			name:       "given a name starting with extension characters when design list then the leading characters are kept",
			designName: "my-app.yaml",
			expected:   "my-app",
		},
		{
			name:       "given a name without an extension when design list then the name is displayed as is",
			designName: "Untitled Design",
			expected:   "Untitled Design",
		},
		{
			name:       "given an archived design name when design list then the double extension is dropped",
			designName: "bookinfo.tar.gz",
			expected:   "bookinfo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := &models.PatternsAPIResponse{
				TotalCount: 1,
				Patterns: []models.MesheryPattern{
					{
						ID:        &designID,
						Name:      tt.designName,
						CreatedAt: &timestamp,
						UpdatedAt: &timestamp,
					},
				},
			}

			displayData, totalCount := processDesignData(response)

			assert.Equal(t, int64(1), totalCount)
			assert.Equal(t, tt.expected, displayData[0][1])
		})
	}
}
