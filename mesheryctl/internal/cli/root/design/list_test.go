package design

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/core"
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

func TestProcessDesignDataNameTrimming(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		designName   string
		expectedName string
	}{
		{
			name:         "yaml extension is removed",
			designName:   "my-app.yaml",
			expectedName: "my-app",
		},
		{
			name:         "short name with yaml extension",
			designName:   "app.yaml",
			expectedName: "app",
		},
		{
			name:         "json extension is removed",
			designName:   "schema.json",
			expectedName: "schema",
		},
		{
			name:         "wasm extension is removed",
			designName:   "awesome.wasm",
			expectedName: "awesome",
		},
		{
			name:         "name without extension is unchanged",
			designName:   "my-design",
			expectedName: "my-design",
		},
		{
			name:         "name with multiple dots keeps all but last extension",
			designName:   "my.app.config.yaml",
			expectedName: "my.app.config",
		},
		{
			name:         "dotfile with extension keeps dot prefix and removes extension",
			designName:   ".hidden.yaml",
			expectedName: ".hidden",
		},
		{
			name:         "empty name stays empty",
			designName:   "",
			expectedName: "",
		},
	}

	// Use local provider so the output row is [id, name, created, updated]
	// (no owner column), making the name index predictable at [1]. Restore the
	// package-global value so this unit test does not leak state into other tests.
	previousProvider := provider
	provider = "None"
	defer func() { provider = previousProvider }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := core.Uuid(uuid.Must(uuid.NewV4()))
			resp := &models.PatternsAPIResponse{
				TotalCount: 1,
				Patterns: []models.MesheryPattern{
					{
						ID:        &id,
						Name:      tt.designName,
						CreatedAt: &now,
						UpdatedAt: &now,
					},
				},
			}

			rows, count := processDesignData(resp)
			if count != 1 {
				t.Fatalf("expected count 1, got %d", count)
			}
			if len(rows) != 1 {
				t.Fatalf("expected 1 row, got %d", len(rows))
			}
			// Local provider row: [designId, designName, createdAt, updatedAt]
			got := rows[0][1]
			if got != tt.expectedName {
				t.Errorf("processDesignData(%q) name = %q, want %q", tt.designName, got, tt.expectedName)
			}
		})
	}
}
