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
