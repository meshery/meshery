package design

import (
	"errors"
	"fmt"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func viewTestDesign(id, name string) map[string]interface{} {
	return map[string]interface{}{"id": id, "name": name}
}

func viewTestPage(total int, designs ...map[string]interface{}) *map[string]interface{} {
	patterns := make([]interface{}, len(designs))
	for i, d := range designs {
		patterns[i] = d
	}
	return &map[string]interface{}{"totalCount": float64(total), "patterns": patterns}
}

func viewSearchURL(name string, page int) string {
	return fmt.Sprintf("api/pattern?page=%d&pagesize=10000&populate=pattern_file&search=%s", page, name)
}

// fakeDesignSearch serves pages by request URL and records each request.
func fakeDesignSearch(pages map[string]*map[string]interface{}, requests *[]string) func(string) (*map[string]interface{}, error) {
	return func(u string) (*map[string]interface{}, error) {
		*requests = append(*requests, u)
		page, ok := pages[u]
		if !ok {
			return nil, fmt.Errorf("unexpected request %s", u)
		}
		return page, nil
	}
}

func designNames(designs []map[string]interface{}) []string {
	names := make([]string, len(designs))
	for i, d := range designs {
		names[i], _ = d["name"].(string)
	}
	return names
}

func TestSelectDesignByName(t *testing.T) {
	tests := []struct {
		name    string
		designs []map[string]interface{}
		arg     string
		wantID  string
		wantErr error
	}{
		{
			name:    "given the exact match first when selecting then it is chosen",
			designs: []map[string]interface{}{viewTestDesign("1", "nginx"), viewTestDesign("2", "nginx-staging")},
			arg:     "nginx",
			wantID:  "1",
		},
		{
			name:    "given the exact match after a substring match when selecting then the exact match is chosen",
			designs: []map[string]interface{}{viewTestDesign("2", "nginx-staging"), viewTestDesign("1", "nginx")},
			arg:     "nginx",
			wantID:  "1",
		},
		{
			name:    "given the exact match after several substring matches when selecting then the exact match is chosen",
			designs: []map[string]interface{}{viewTestDesign("2", "nginx-a"), viewTestDesign("3", "nginx-b"), viewTestDesign("1", "nginx")},
			arg:     "nginx",
			wantID:  "1",
		},
		{
			name:    "given only a match that differs in case when selecting then it is chosen",
			designs: []map[string]interface{}{viewTestDesign("2", "nginx-staging"), viewTestDesign("1", "NGINX")},
			arg:     "nginx",
			wantID:  "1",
		},
		{
			name:    "given an exact match and a match that differs in case when selecting then the exact match is chosen",
			designs: []map[string]interface{}{viewTestDesign("2", "NGINX"), viewTestDesign("1", "nginx")},
			arg:     "nginx",
			wantID:  "1",
		},
		{
			name:    "given a name with spaces when selecting then the whole name is matched",
			designs: []map[string]interface{}{viewTestDesign("2", "my nginx design v2"), viewTestDesign("1", "my nginx design")},
			arg:     "my nginx design",
			wantID:  "1",
		},
		{
			name:    "given only a substring match when selecting then the design is not found",
			designs: []map[string]interface{}{viewTestDesign("2", "nginx-staging")},
			arg:     "nginx",
			wantErr: ErrDesignNotFound("nginx"),
		},
		{
			name:    "given several substring matches when selecting then the design is not found",
			designs: []map[string]interface{}{viewTestDesign("2", "nginx-staging"), viewTestDesign("3", "nginx-prod")},
			arg:     "nginx",
			wantErr: ErrDesignNotFound("nginx"),
		},
		{
			name:    "given no designs when selecting then the design is not found",
			arg:     "nginx",
			wantErr: ErrDesignNotFound("nginx"),
		},
		{
			name:    "given a design without a name when selecting then it is not matched",
			designs: []map[string]interface{}{{"id": "2"}, viewTestDesign("1", "nginx")},
			arg:     "nginx",
			wantID:  "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectDesignByName(tt.designs, tt.arg)
			if tt.wantErr != nil {
				utils.AssertMeshkitErrorsEqual(t, err, tt.wantErr)
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got["id"] != tt.wantID {
				t.Fatalf("selected %v, want id %s", got, tt.wantID)
			}
		})
	}
}

func TestSelectDesignByNameIsIndependentOfOrder(t *testing.T) {
	a, b := viewTestDesign("1", "nginx"), viewTestDesign("2", "nginx-staging")
	first, err := selectDesignByName([]map[string]interface{}{a, b}, "nginx")
	if err != nil {
		t.Fatal(err)
	}
	second, err := selectDesignByName([]map[string]interface{}{b, a}, "nginx")
	if err != nil {
		t.Fatal(err)
	}
	if first["id"] != "1" || second["id"] != "1" {
		t.Fatalf("selection depends on order: %v, %v", first["id"], second["id"])
	}
}

func TestSelectDesignByNameWithDuplicateNames(t *testing.T) {
	original := runDesignSelectPrompt
	t.Cleanup(func() { runDesignSelectPrompt = original })

	designs := []map[string]interface{}{viewTestDesign("2", "nginx-staging"), viewTestDesign("1", "nginx"), viewTestDesign("3", "nginx")}

	t.Run("given two designs with the name when selecting then the user picks one of them", func(t *testing.T) {
		var offered []string
		runDesignSelectPrompt = func(_ string, items []string) (int, error) {
			offered = items
			return 1, nil
		}
		got, err := selectDesignByName(designs, "nginx")
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"Name: nginx, ID: 1", "Name: nginx, ID: 3"}
		if fmt.Sprint(offered) != fmt.Sprint(want) {
			t.Fatalf("offered %v, want %v", offered, want)
		}
		if got["id"] != "3" {
			t.Fatalf("selected %v, want id 3", got["id"])
		}
	})

	t.Run("given two designs with the name when the prompt fails then the error is returned", func(t *testing.T) {
		runDesignSelectPrompt = func(string, []string) (int, error) { return 0, utils.ErrPromptCancelled() }
		_, err := selectDesignByName(designs, "nginx")
		utils.AssertMeshkitErrorsEqual(t, err, utils.ErrPromptCancelled())
	})

	t.Run("given two designs with the name when the prompt returns an invalid index then an error is returned", func(t *testing.T) {
		runDesignSelectPrompt = func(string, []string) (int, error) { return 5, nil }
		if _, err := selectDesignByName(designs, "nginx"); err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestFetchDesignSearch(t *testing.T) {
	pageOf := func(prefix string, n int) []map[string]interface{} {
		designs := make([]map[string]interface{}, n)
		for i := range designs {
			designs[i] = viewTestDesign(fmt.Sprintf("%s%02d", prefix, i), fmt.Sprintf("nginx-%02d", i))
		}
		return designs
	}

	t.Run("given a complete first page when searching then one request is made", func(t *testing.T) {
		var requests []string
		pages := map[string]*map[string]interface{}{
			viewSearchURL("nginx", 0): viewTestPage(2, viewTestDesign("2", "nginx-staging"), viewTestDesign("1", "nginx")),
		}
		got, err := fetchDesignSearch("nginx", fakeDesignSearch(pages, &requests))
		if err != nil {
			t.Fatal(err)
		}
		if len(requests) != 1 || fmt.Sprint(designNames(got)) != "[nginx-staging nginx]" {
			t.Fatalf("requests %v, designs %v", requests, designNames(got))
		}
	})

	t.Run("given the exact match on the second page when viewing then it is found", func(t *testing.T) {
		var requests []string
		pages := map[string]*map[string]interface{}{
			viewSearchURL("nginx", 0): viewTestPage(11, pageOf("p", 10)...),
			viewSearchURL("nginx", 1): viewTestPage(11, viewTestDesign("exact", "nginx")),
		}
		got, err := fetchDesignSearch("nginx", fakeDesignSearch(pages, &requests))
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 11 || len(requests) != 2 {
			t.Fatalf("got %d designs in %d requests", len(got), len(requests))
		}
		selected, err := selectDesignByName(got, "nginx")
		if err != nil {
			t.Fatal(err)
		}
		if selected["id"] != "exact" {
			t.Fatalf("selected %v", selected["id"])
		}
	})

	t.Run("given a design repeated on the next page when searching then it is kept once", func(t *testing.T) {
		var requests []string
		repeated := viewTestDesign("1", "nginx")
		pages := map[string]*map[string]interface{}{
			viewSearchURL("nginx", 0): viewTestPage(2, repeated),
			viewSearchURL("nginx", 1): viewTestPage(2, repeated, viewTestDesign("2", "nginx-staging")),
		}
		got, err := fetchDesignSearch("nginx", fakeDesignSearch(pages, &requests))
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprint(designNames(got)) != "[nginx nginx-staging]" {
			t.Fatalf("designs %v", designNames(got))
		}
	})

	t.Run("given a page that adds no designs before the total when searching then an error is returned", func(t *testing.T) {
		var requests []string
		pages := map[string]*map[string]interface{}{
			viewSearchURL("nginx", 0): viewTestPage(5, viewTestDesign("1", "nginx")),
			viewSearchURL("nginx", 1): viewTestPage(5),
		}
		_, err := fetchDesignSearch("nginx", fakeDesignSearch(pages, &requests))
		utils.AssertMeshkitErrorsEqual(t, err, ErrDesignInvalidApiResponse(`search for "nginx" returned 1 of 5 designs`))
		if len(requests) != 2 {
			t.Fatalf("requests %v", requests)
		}
	})

	t.Run("given a page that only repeats designs when searching then an error is returned", func(t *testing.T) {
		var requests []string
		repeated := viewTestDesign("1", "nginx")
		pages := map[string]*map[string]interface{}{
			viewSearchURL("nginx", 0): viewTestPage(3, repeated),
			viewSearchURL("nginx", 1): viewTestPage(3, repeated),
		}
		if _, err := fetchDesignSearch("nginx", fakeDesignSearch(pages, &requests)); err == nil {
			t.Fatal("expected an error")
		}
		if len(requests) != 2 {
			t.Fatalf("requests %v", requests)
		}
	})

	t.Run("given a response without a total when searching then the first page is used", func(t *testing.T) {
		var requests []string
		pages := map[string]*map[string]interface{}{
			viewSearchURL("nginx", 0): {"patterns": []interface{}{viewTestDesign("1", "nginx")}},
		}
		got, err := fetchDesignSearch("nginx", fakeDesignSearch(pages, &requests))
		if err != nil || len(got) != 1 || len(requests) != 1 {
			t.Fatalf("got %v, err %v, requests %v", got, err, requests)
		}
	})

	t.Run("given a name with spaces when searching then it is query-escaped", func(t *testing.T) {
		var requests []string
		pages := map[string]*map[string]interface{}{
			viewSearchURL("my+design", 0): viewTestPage(0),
		}
		if _, err := fetchDesignSearch("my design", fakeDesignSearch(pages, &requests)); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("given the request fails when searching then its error is returned", func(t *testing.T) {
		apiErr := errors.New("connection refused")
		_, err := fetchDesignSearch("nginx", func(string) (*map[string]interface{}, error) { return nil, apiErr })
		if !errors.Is(err, apiErr) {
			t.Fatalf("got %v, want %v", err, apiErr)
		}
	})

	malformed := []struct {
		name     string
		response *map[string]interface{}
		want     error
	}{
		{"an empty response", nil, ErrDesignInvalidApiResponse("empty response")},
		{"no patterns field", &map[string]interface{}{"totalCount": float64(0)}, ErrDesignInvalidApiResponse("'patterns' field missing")},
		{"a null patterns field", &map[string]interface{}{"patterns": nil}, ErrDesignInvalidApiResponse("'patterns' field is not of expected type")},
		{"a patterns field that is not a list", &map[string]interface{}{"patterns": "x"}, ErrDesignInvalidApiResponse("'patterns' field is not of expected type")},
		{"a design that is not an object", &map[string]interface{}{"patterns": []interface{}{"x"}}, ErrDesignInvalidApiResponse("design entry is not an object")},
	}
	for _, tt := range malformed {
		t.Run("given "+tt.name+" when searching then an invalid response error is returned", func(t *testing.T) {
			_, err := fetchDesignSearch("nginx", func(string) (*map[string]interface{}, error) { return tt.response, nil })
			utils.AssertMeshkitErrorsEqual(t, err, tt.want)
		})
	}
}
