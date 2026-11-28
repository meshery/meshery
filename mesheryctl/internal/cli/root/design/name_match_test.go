package design

import (
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
)

// GET /api/pattern?search= is a substring match on the server. When the only
// stored design is "nginx-prod", `design deploy nginx` and `design apply nginx`
// must report that no design is named "nginx" instead of deploying nginx-prod.
func TestDesignByNameRejectsSubstringMatch(t *testing.T) {
	testContext := utils.NewTestHelper(t)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	urls := []utils.MockURL{
		{Method: "GET", URL: testContext.BaseURL + "/api/pattern/types", Response: "view.designTypes.response.golden", ResponseCode: 200},
		{Method: "GET", URL: testContext.BaseURL + "/api/pattern?populate=pattern_file&search=nginx", Response: "deploy.substringMatch.response.golden", ResponseCode: 200},
		{Method: "POST", URL: testContext.BaseURL + "/api/pattern/deploy", Response: "deploy.designdeploy.response.golden", ResponseCode: 200},
	}
	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:          "given only nginx-prod is stored when design deploy nginx then report not found",
			Args:          []string{"deploy", "nginx"},
			URLs:          urls,
			ExpectError:   true,
			ExpectedError: ErrDesignNotFound("nginx"),
		},
		{
			Name:          "given only nginx-prod is stored when design apply nginx then report not found",
			Args:          []string{"apply", "nginx"},
			URLs:          urls,
			ExpectError:   true,
			ExpectedError: ErrDesignNotFound("nginx"),
		},
	}
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}

// The server's `page` parameter is zero-based with 10 designs per page. When 10
// designs containing "nginx" fill page 0 and the design named exactly "nginx" is
// on page 1, `design deploy nginx` and `design apply nginx` must deploy it.
func TestDesignByNameFindsExactMatchOnSecondPage(t *testing.T) {
	testContext := utils.NewTestHelper(t)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	searchURLs := []utils.MockURL{
		{Method: "GET", URL: testContext.BaseURL + "/api/pattern/types", Response: "view.designTypes.response.golden", ResponseCode: 200},
		{Method: "GET", URL: testContext.BaseURL + "/api/pattern?populate=pattern_file&search=nginx", Response: "design.search.page0.response.golden", ResponseCode: 200},
		{Method: "GET", URL: testContext.BaseURL + "/api/pattern?page=1&populate=pattern_file&search=nginx", Response: "design.search.page1.response.golden", ResponseCode: 200},
	}
	withDeploy := func(response string) []utils.MockURL {
		return append(slices.Clone(searchURLs), utils.MockURL{Method: "POST", URL: testContext.BaseURL + "/api/pattern/deploy", Response: response, ResponseCode: 200})
	}
	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:             "given nginx is on the second page when design deploy nginx then design is deployed",
			Args:             []string{"deploy", "nginx"},
			URLs:             withDeploy("deploy.designdeploy.response.golden"),
			ExpectedResponse: "deploy.output.golden",
		},
		{
			Name:             "given nginx is on the second page when design apply nginx then design is applied",
			Args:             []string{"apply", "nginx"},
			URLs:             withDeploy("apply.designDeploy.response.golden"),
			ExpectedResponse: "apply.output.golden",
		},
	}
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}

func TestExactNameMatches(t *testing.T) {
	designs := func(names ...string) []models.MesheryPattern {
		out := make([]models.MesheryPattern, 0, len(names))
		for _, n := range names {
			out = append(out, models.MesheryPattern{Name: n})
		}
		return out
	}
	tests := []struct {
		name     string
		patterns []models.MesheryPattern
		query    string
		want     []string
	}{
		{name: "substring only is not a match", patterns: designs("nginx-prod", "my-nginx"), query: "nginx", want: nil},
		{name: "exact match is kept", patterns: designs("nginx-prod", "nginx"), query: "nginx", want: []string{"nginx"}},
		{name: "exact match wins over case-insensitive match", patterns: designs("NGINX", "nginx"), query: "nginx", want: []string{"nginx"}},
		{name: "case-insensitive match is the fallback", patterns: designs("nginx-prod", "NGINX"), query: "nginx", want: []string{"NGINX"}},
		{name: "every exact match is kept so the user can choose", patterns: designs("nginx", "nginx-prod", "nginx"), query: "nginx", want: []string{"nginx", "nginx"}},
		{name: "name with spaces", patterns: designs("my design", "my design v2"), query: "my design", want: []string{"my design"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, p := range exactNameMatches(tt.patterns, tt.query) {
				got = append(got, p.Name)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("exactNameMatches(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

// The server returns the design search one page at a time (10 per page by
// default). fetchAllPatternPages must gather every page so an exact design past
// the first page is not reported as missing, and it must stop instead of
// looping forever.
func TestFetchAllPatternPages(t *testing.T) {
	design := func(name string) models.MesheryPattern { return models.MesheryPattern{Name: name} }

	t.Run("collects designs across every page until totalCount is reached", func(t *testing.T) {
		pages := []*models.PatternsAPIResponse{
			{Page: 0, PageSize: 2, TotalCount: 3, Patterns: []models.MesheryPattern{design("a"), design("b")}},
			{Page: 1, PageSize: 2, TotalCount: 3, Patterns: []models.MesheryPattern{design("exact")}},
		}
		var requested []uint
		got, err := fetchAllPatternPages(func(page uint) (*models.PatternsAPIResponse, error) {
			requested = append(requested, page)
			return pages[page], nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var names []string
		for _, p := range got {
			names = append(names, p.Name)
		}
		if !slices.Equal(names, []string{"a", "b", "exact"}) {
			t.Errorf("collected %q, want [a b exact]", names)
		}
		if !slices.Equal(requested, []uint{0, 1}) {
			t.Errorf("requested pages %v, want [0 1]", requested)
		}
	})

	t.Run("stops at the first page when it holds every design", func(t *testing.T) {
		calls := 0
		got, err := fetchAllPatternPages(func(_ uint) (*models.PatternsAPIResponse, error) {
			calls++
			return &models.PatternsAPIResponse{Page: 0, PageSize: 10, TotalCount: 1, Patterns: []models.MesheryPattern{design("only")}}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 1 {
			t.Errorf("fetched %d pages, want 1", calls)
		}
		if len(got) != 1 || got[0].Name != "only" {
			t.Errorf("collected %v, want [only]", got)
		}
	})

	t.Run("stops on an empty page instead of looping forever", func(t *testing.T) {
		got, err := fetchAllPatternPages(func(page uint) (*models.PatternsAPIResponse, error) {
			return &models.PatternsAPIResponse{Page: page, PageSize: 10, TotalCount: 999, Patterns: nil}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("collected %v, want none", got)
		}
	})
}
