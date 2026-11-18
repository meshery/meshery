package design

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
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

// untotalledPage is a page from a provider that reports no total: the page
// size is then the only end marker.
func untotalledPage(designs ...map[string]interface{}) *map[string]interface{} {
	patterns := make([]interface{}, len(designs))
	for i, d := range designs {
		patterns[i] = d
	}
	return &map[string]interface{}{"patterns": patterns}
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
		if !slices.Equal(offered, want) {
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

	for _, index := range []int{-1, 5} {
		t.Run(fmt.Sprintf("given two designs with the name when the prompt returns index %d then an error is returned", index), func(t *testing.T) {
			runDesignSelectPrompt = func(string, []string) (int, error) { return index, nil }
			_, err := selectDesignByName(designs, "nginx")
			utils.AssertMeshkitErrorsEqual(t, err, ErrDesignInvalidSelection(index, 2))
		})
	}
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

	t.Run("given a short page without a total when searching then it is the last page", func(t *testing.T) {
		var requests []string
		pages := map[string]*map[string]interface{}{
			viewSearchURL("nginx", 0): {"patterns": []interface{}{viewTestDesign("1", "nginx")}},
		}
		got, err := fetchDesignSearch("nginx", fakeDesignSearch(pages, &requests))
		if err != nil || len(got) != 1 || len(requests) != 1 {
			t.Fatalf("got %v, err %v, requests %v", got, err, requests)
		}
	})

	t.Run("given a full page without a total when searching then the next page is requested", func(t *testing.T) {
		var requests []string
		first := make([]map[string]interface{}, designSearchPageSize)
		for i := range first {
			first[i] = viewTestDesign(fmt.Sprintf("p%05d", i), fmt.Sprintf("nginx-%05d", i))
		}
		pages := map[string]*map[string]interface{}{
			viewSearchURL("nginx", 0): untotalledPage(first...),
			viewSearchURL("nginx", 1): untotalledPage(viewTestDesign("exact", "nginx")),
		}
		got, err := fetchDesignSearch("nginx", fakeDesignSearch(pages, &requests))
		if err != nil {
			t.Fatal(err)
		}
		if len(requests) != 2 {
			t.Fatalf("requests %v", requests)
		}
		selected, err := selectDesignByName(got, "nginx")
		if err != nil {
			t.Fatal(err)
		}
		if selected["id"] != "exact" {
			t.Fatalf("selected %v, want id exact", selected["id"])
		}
	})

	t.Run("given a full page without a total that only repeats designs when searching then an error is returned", func(t *testing.T) {
		var requests []string
		first := make([]map[string]interface{}, designSearchPageSize)
		for i := range first {
			first[i] = viewTestDesign(fmt.Sprintf("p%05d", i), fmt.Sprintf("nginx-%05d", i))
		}
		pages := map[string]*map[string]interface{}{
			viewSearchURL("nginx", 0): untotalledPage(first...),
			viewSearchURL("nginx", 1): untotalledPage(first...),
		}
		_, err := fetchDesignSearch("nginx", fakeDesignSearch(pages, &requests))
		utils.AssertMeshkitErrorsEqual(t, err, ErrDesignInvalidApiResponse(`search for "nginx" repeated a full page of designs`))
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

func TestDesignSearchTotal(t *testing.T) {
	tests := []struct {
		name     string
		response map[string]interface{}
		want     int
		wantOK   bool
		wantErr  bool
	}{
		{name: "float64", response: map[string]interface{}{"totalCount": float64(11)}, want: 11, wantOK: true},
		{name: "int", response: map[string]interface{}{"totalCount": 11}, want: 11, wantOK: true},
		{name: "int64", response: map[string]interface{}{"totalCount": int64(11)}, want: 11, wantOK: true},
		{name: "json.Number", response: map[string]interface{}{"totalCount": json.Number("11")}, want: 11, wantOK: true},
		{name: "legacy total_count", response: map[string]interface{}{"total_count": float64(11)}, want: 11, wantOK: true},
		{name: "totalCount preferred over total_count", response: map[string]interface{}{"totalCount": float64(11), "total_count": float64(3)}, want: 11, wantOK: true},
		{name: "missing", response: map[string]interface{}{}},
		{name: "null", response: map[string]interface{}{"totalCount": nil}},
		{name: "string", response: map[string]interface{}{"totalCount": "11"}, wantErr: true},
		{name: "zero", response: map[string]interface{}{"totalCount": float64(0)}, want: 0, wantOK: true},
		{name: "negative float64", response: map[string]interface{}{"totalCount": float64(-1)}, wantErr: true},
		{name: "negative int", response: map[string]interface{}{"totalCount": -1}, wantErr: true},
		{name: "negative int64", response: map[string]interface{}{"totalCount": int64(-1)}, wantErr: true},
		{name: "negative json.Number", response: map[string]interface{}{"totalCount": json.Number("-1")}, wantErr: true},
		{name: "float64 above the int range", response: map[string]interface{}{"totalCount": math.MaxFloat64}, wantErr: true},
		{name: "infinite float64", response: map[string]interface{}{"totalCount": math.Inf(1)}, wantErr: true},
		{name: "json.Number above the int64 range", response: map[string]interface{}{"totalCount": json.Number("9223372036854775808")}, wantErr: true},
		{name: "fractional float64", response: map[string]interface{}{"totalCount": 1.5}, wantErr: true},
		{name: "fractional json.Number", response: map[string]interface{}{"totalCount": json.Number("1.5")}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run("given a "+tt.name+" total when reading it then the count is resolved", func(t *testing.T) {
			got, ok, err := designSearchTotal(tt.response)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %d %v", got, ok)
				}
				return
			}
			if err != nil || got != tt.want || ok != tt.wantOK {
				t.Fatalf("got (%d, %v, %v), want (%d, %v, nil)", got, ok, err, tt.want, tt.wantOK)
			}
		})
	}
}

func TestFetchDesignSearchWithNonFloatTotals(t *testing.T) {
	totals := map[string]interface{}{
		"int":         11,
		"json.Number": json.Number("11"),
	}
	for name, total := range totals {
		t.Run("given a "+name+" total when the exact match is on the second page then it is found", func(t *testing.T) {
			first := make([]interface{}, 10)
			for i := range first {
				first[i] = viewTestDesign(fmt.Sprintf("p%02d", i), fmt.Sprintf("nginx-%02d", i))
			}
			var requests []string
			pages := map[string]*map[string]interface{}{
				viewSearchURL("nginx", 0): {"totalCount": total, "patterns": first},
				viewSearchURL("nginx", 1): {"totalCount": total, "patterns": []interface{}{viewTestDesign("exact", "nginx")}},
			}
			got, err := fetchDesignSearch("nginx", fakeDesignSearch(pages, &requests))
			if err != nil {
				t.Fatal(err)
			}
			selected, err := selectDesignByName(got, "nginx")
			if err != nil || selected["id"] != "exact" || len(requests) != 2 {
				t.Fatalf("selected %v, err %v, requests %v", selected, err, requests)
			}
		})
	}

	t.Run("given a negative total when searching then an invalid response error is returned", func(t *testing.T) {
		_, err := fetchDesignSearch("nginx", func(string) (*map[string]interface{}, error) {
			return &map[string]interface{}{"totalCount": float64(-1), "patterns": []interface{}{}}, nil
		})
		utils.AssertMeshkitErrorsEqual(t, err, ErrDesignInvalidApiResponse("'totalCount' field is not a count: -1"))
	})

	t.Run("given a total that is not a number when searching then an invalid response error is returned", func(t *testing.T) {
		_, err := fetchDesignSearch("nginx", func(string) (*map[string]interface{}, error) {
			return &map[string]interface{}{"totalCount": "eleven", "patterns": []interface{}{}}, nil
		})
		utils.AssertMeshkitErrorsEqual(t, err, ErrDesignInvalidApiResponse("'totalCount' field is not a count: eleven"))
	})
}
