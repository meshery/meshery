package filter

import (
	"fmt"
	"net/url"
	"path/filepath"

	mkerrors "github.com/meshery/meshkit/errors"
	"runtime"
	"testing"

	"github.com/meshery/meshery/server/models"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
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

	utils.InvokeMesheryctlTestListCommand(t, update, FilterCmd, ListTests, currDir, "filter")
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

func TestSelectFilterByName(t *testing.T) {
	original := runFilterSelectPrompt
	t.Cleanup(func() { runFilterSelectPrompt = original })
	first := &models.MesheryFilter{Name: "auth"}
	second := &models.MesheryFilter{Name: "auth"}
	upper := &models.MesheryFilter{Name: "AUTH"}
	partial := &models.MesheryFilter{Name: "oauth2-proxy"}
	tests := []struct {
		name      string
		filters   []*models.MesheryFilter
		query     string
		want      *models.MesheryFilter
		wantError bool
		prompt    bool
	}{
		{"exact before substring", []*models.MesheryFilter{partial, first}, "auth", first, false, false},
		{"exact before case-insensitive", []*models.MesheryFilter{upper, first}, "auth", first, false, false},
		{"case-insensitive", []*models.MesheryFilter{upper}, "auth", upper, false, false},
		{"substring rejected", []*models.MesheryFilter{partial}, "auth", nil, true, false},
		{"duplicate exact prompts", []*models.MesheryFilter{first, second}, "auth", second, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			runFilterSelectPrompt = func(_ string, _ []string) (int, error) { called = true; return 1, nil }
			got, err := selectFilterByName(tt.filters, tt.query)
			if (err != nil) != tt.wantError || got != tt.want || called != tt.prompt {
				t.Fatalf("got %v, err %v, prompted %v", got, err, called)
			}
		})
	}
}

func TestFetchFilterSearchPages(t *testing.T) {
	calls := 0
	got, err := fetchFilterSearch("auth & token", func(raw string) (*models.MesheryFilterPage, error) {
		u, parseErr := url.Parse(raw)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		if u.Query().Get("search") != "auth & token" || u.Query().Get("page") != fmt.Sprint(calls) {
			t.Fatalf("unexpected query %s", raw)
		}
		calls++
		if calls == 1 {
			return &models.MesheryFilterPage{PageSize: 1, TotalCount: 2, Filters: []*models.MesheryFilter{{Name: "partial"}}}, nil
		}
		return &models.MesheryFilterPage{PageSize: 1, TotalCount: 2, Filters: []*models.MesheryFilter{{Name: "auth & token"}}}, nil
	})
	if err != nil || calls != 2 || len(got.Filters) != 2 || got.Filters[1].Name != "auth & token" {
		t.Fatalf("pages=%d result=%v error=%v", calls, got, err)
	}
}

func TestFilterViewNotFoundReturnsTypedError(t *testing.T) {
	mesheryctlflags.InitValidators(FilterCmd)
	run := utils.MesheryCommandRun{Cmd: FilterCmd, Args: []string{"view", "xyz"}, CommandDir: utils.CallerDir(t), Mocks: []utils.MockURL{{URL: "/api/filter", Response: "view.nonexisting.filter.api.response.golden"}}}
	_, err := utils.RunMesheryctlCommand(t, run)
	if err == nil || mkerrors.GetCode(err) != ErrFilterNotFoundCode {
		t.Fatalf("got %v, want %s", err, ErrFilterNotFoundCode)
	}
}
