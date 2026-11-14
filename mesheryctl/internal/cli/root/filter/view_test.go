package filter

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gofrs/uuid"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
)

func TestViewCmd(t *testing.T) {
	mesheryctlflags.InitValidators(FilterCmd)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	ListTests := []utils.MesheryListCommandTest{
		{
			Name:             "Fetch Filter View",
			Args:             []string{"view", "KumaTest"},
			ExpectedResponse: "view.filter.output.golden",
			Fixture:          "view.filter.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
		},
		{
			Name:             "Fetch Kuma Filter View with ID",
			Args:             []string{"view", "957fbc9b-a655-4892-823d-375102a9587c"},
			ExpectedResponse: "view.id.filter.output.golden",
			Fixture:          "view.id.filter.api.response.golden",
			URL:              "/api/filter/957fbc9b-a655-4892-823d-375102a9587c",
			ExpectError:      false,
		},
	}

	loggerTests := []utils.MesheryCommandTest{
		{
			Name:           "Fetch Filter View for non existing filter",
			Args:           []string{"view", "xyz"},
			Fixture:        "view.nonexisting.filter.api.response.golden",
			URL:            "/api/filter",
			HttpMethod:     "GET",
			HttpStatusCode: 200,
			ExpectError:    true,
			ExpectedError:  ErrFilterNotFound("xyz", "view"),
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, FilterCmd, ListTests, currDir, "filter")
	utils.InvokeMesheryctlTestCommand(t, update, FilterCmd, loggerTests, currDir, "filter")
}

func TestResolveFilterByName(t *testing.T) {
	u1 := core.Uuid(uuid.Must(uuid.NewV4()))
	u2 := core.Uuid(uuid.Must(uuid.NewV4()))

	tests := []struct {
		name        string
		filters     []*models.MesheryFilter
		searchName  string
		totalCount  int
		wantID      *core.Uuid
		wantErrCode string
	}{
		{
			name:       "one exact match returns that filter",
			filters:    []*models.MesheryFilter{{ID: &u1, Name: "my-filter"}},
			searchName: "my-filter",
			wantID:     &u1,
		},
		{
			name:        "substring match is not the filter",
			filters:    []*models.MesheryFilter{{ID: &u1, Name: "my-filter-v2"}},
			searchName: "my-filter",
			wantErrCode: ErrFilterNotFoundCode,
		},
		{
			name: "exact match wins over substring match",
			filters: []*models.MesheryFilter{
				{ID: &u1, Name: "my-filter-v2"},
				{ID: &u2, Name: "my-filter"},
			},
			searchName: "my-filter",
			wantID:     &u2,
		},
		{
			name:       "case-insensitive match matches when only case differs",
			filters:    []*models.MesheryFilter{{ID: &u1, Name: "My-Filter"}},
			searchName: "my-filter",
			wantID:     &u1,
		},
		{
			name:        "case-insensitive match is not used when results are truncated",
			filters:    []*models.MesheryFilter{{ID: &u1, Name: "My-Filter"}},
			searchName:  "my-filter",
			totalCount:  5,
			wantErrCode: utils.ErrInvalidArgumentCode,
		},
		{
			name: "exact match is used even when results are truncated",
			filters: []*models.MesheryFilter{
				{ID: &u1, Name: "my-filter-v2"},
				{ID: &u2, Name: "my-filter"},
			},
			searchName: "my-filter",
			totalCount: 5,
			wantID:     &u2,
		},
		{
			name:        "no match returns ErrFilterNotFound",
			filters:    []*models.MesheryFilter{},
			searchName: "nonexistent",
			wantErrCode: ErrFilterNotFoundCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalCount := tt.totalCount
			if totalCount == 0 {
				totalCount = len(tt.filters)
			}
			page := &models.MesheryFilterPage{
				TotalCount: totalCount,
				Filters:    tt.filters,
			}

			resolved, err := resolveFilterByName(page, tt.searchName)
			if tt.wantErrCode != "" {
				if err == nil {
					t.Fatalf("expected error code %s, got nil", tt.wantErrCode)
				}
				if code := meshkiterrors.GetCode(err); code != tt.wantErrCode {
					t.Fatalf("expected error code %s, got %s (%v)", tt.wantErrCode, code, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resolved == nil || resolved.ID != tt.wantID {
				t.Fatalf("expected filter ID %v, got %v", tt.wantID, resolved)
			}
		})
	}
}

func TestFilterViewSaveCreatesFileWithExtension(t *testing.T) {
	mesheryctlflags.InitValidators(FilterCmd)

	const filterID = "957fbc9b-a655-4892-823d-375102a9587c"

	mesheryDir := utils.IsolateMesheryHome(t)

	run := utils.MesheryCommandRun{
		Cmd:        FilterCmd,
		Args:       []string{"view", filterID, "--output-format", "json", "--save"},
		CommandDir: utils.CallerDir(t),
		Mocks: []utils.MockURL{
			{
				URL:      "/api/filter",
				Response: "view.filter.api.response.golden",
			},
			{
				URL:      "/api/filter/" + filterID,
				Response: "view.id.filter.api.response.golden",
			},
		},
	}
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	utils.AssertFileExists(t, filepath.Join(mesheryDir, "filter_KumaTest_957fbc9b.json"))
}
