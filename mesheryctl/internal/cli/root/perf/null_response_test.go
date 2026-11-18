package perf

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// A server, or a proxy or provider in front of it, can answer 200 with the JSON
// value null; json.Marshal writes exactly that for a nil slice or pointer. It has
// to read as "nothing there". Decoding into a pointer left the page nil, and the
// next field access crashed the command.
func TestFetchPerformanceProfiles_NullResponseIsEmpty(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	httpmock.RegisterResponder(http.MethodGet, testContext.BaseURL+"/api/performance/profiles",
		httpmock.NewStringResponder(http.StatusOK, "null"))

	profiles, _, err := fetchPerformanceProfiles(testContext.BaseURL, "", 10, 0)
	if err != nil {
		t.Fatalf("a null response should read as an empty page, got error: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("profiles = %v, want none", profiles)
	}
}

func TestFetchPerformanceProfileResults_NullResponseIsEmpty(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	profileID := "a947dff3-1415-4ca6-8922-e35b3232bd78"
	httpmock.RegisterResponder(http.MethodGet, testContext.BaseURL+"/api/user/performance/profiles/"+profileID+"/results",
		httpmock.NewStringResponder(http.StatusOK, "null"))

	results, _, err := fetchPerformanceProfileResults(testContext.BaseURL, profileID, 10, 0)
	if err != nil {
		t.Fatalf("a null response should read as an empty page, got error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("results = %v, want none", results)
	}
}
