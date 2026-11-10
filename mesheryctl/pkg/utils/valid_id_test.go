package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidIDRejectsAmbiguousAndEmptyInputs(t *testing.T) {
	ids := []string{"dc761694-5e61-4dd4-83c9-bca45667ef70", "d7ccca72-adea-42e1-bb19-b28701db050d", "b35c168f-47dc-7e1d-9b60-f4296769f111"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"patterns": []map[string]string{{"id": ids[0]}, {"id": ids[1]}, {"id": ids[2]}}, "totalCount": 3})
	}))
	defer server.Close()
	setupAuthenticatedContext(t, server.URL)

	if _, _, err := ValidId(server.URL, "", "pattern"); err == nil {
		t.Fatal("empty identifier must be rejected")
	}
	if _, _, err := ValidId(server.URL, "d", "pattern"); err == nil || !strings.Contains(err.Error(), ids[0]) || !strings.Contains(err.Error(), ids[1]) {
		t.Fatalf("ambiguous identifier must report both matches, got %v", err)
	}
	if got, isID, err := ValidId(server.URL, ids[0], "pattern"); err != nil || !isID || got != ids[0] {
		t.Fatalf("exact identifier: got %q, %v, %v", got, isID, err)
	}
	if got, isID, err := ValidId(server.URL, ids[2], "pattern"); err != nil || !isID || got != ids[2] {
		t.Fatalf("exact non-v4 identifier: got %q, %v, %v", got, isID, err)
	}
	if got, isID, err := ValidId(server.URL, "dc76", "pattern"); err != nil || !isID || got != ids[0] {
		t.Fatalf("unique prefix: got %q, %v, %v", got, isID, err)
	}
	if got, isID, err := ValidId(server.URL, "my-design", "pattern"); err != nil || isID || got != "my-design" {
		t.Fatalf("name fallback: got %q, %v, %v", got, isID, err)
	}
}

func TestValidIDFindsAmbiguityAcrossPages(t *testing.T) {
	first := "dc761694-5e61-4dd4-83c9-bca45667ef70"
	second := "d7ccca72-adea-42e1-bb19-b28701db050d"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("pagesize") == "" {
			t.Error("pagesize was not provided")
		}
		id := first
		if r.URL.Query().Get("page") == "1" {
			id = second
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"patterns":[{"id":%q}],"pageSize":1,"totalCount":2}`, id)
	}))
	defer server.Close()
	setupAuthenticatedContext(t, server.URL)
	if _, _, err := ValidId(server.URL, "d", "pattern"); err == nil || !strings.Contains(err.Error(), first) || !strings.Contains(err.Error(), second) {
		t.Fatalf("cross-page ambiguity must be rejected, got %v", err)
	}
}

func TestValidIDAllowsShortPageWithoutTotalCount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"patterns":[{"id":"dc761694-5e61-4dd4-83c9-bca45667ef70"}]}`)
	}))
	defer server.Close()
	setupAuthenticatedContext(t, server.URL)
	got, isID, err := ValidId(server.URL, "dc", "pattern")
	if err != nil || !isID || got != "dc761694-5e61-4dd4-83c9-bca45667ef70" {
		t.Fatalf("short page should resolve the ID: got %q, %v, %v", got, isID, err)
	}
}

func TestValidIDRejectsFullPageWithoutTotalCount(t *testing.T) {
	for _, tc := range []struct {
		name            string
		pageSize        int
		includePageSize bool
	}{
		{name: "requested page size", pageSize: 100},
		{name: "server capped page size", pageSize: 10, includePageSize: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				patterns := make([]map[string]string, tc.pageSize)
				for i := range patterns {
					patterns[i] = map[string]string{"id": fmt.Sprintf("d%07x-5e61-4dd4-83c9-bca45667ef70", i)}
				}
				response := map[string]interface{}{"patterns": patterns}
				if tc.includePageSize {
					response["pageSize"] = tc.pageSize
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()
			setupAuthenticatedContext(t, server.URL)
			if _, _, err := ValidId(server.URL, "d", "pattern"); err == nil {
				t.Fatal("must reject a full page without a total count")
			}
		})
	}
}

func TestValidIDLimitsAmbiguityError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		patterns := make([]map[string]string, 12)
		for i := range patterns {
			patterns[i] = map[string]string{"id": fmt.Sprintf("d%07x-5e61-4dd4-83c9-bca45667ef70", i)}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"patterns": patterns, "totalCount": len(patterns)})
	}))
	defer server.Close()
	setupAuthenticatedContext(t, server.URL)
	_, _, err := ValidId(server.URL, "d", "pattern")
	if err == nil || !strings.Contains(err.Error(), "and 2 more") || strings.Count(err.Error(), "-5e61-4dd4-83c9-bca45667ef70") != 10 {
		t.Fatalf("expected 10 matches and the remaining count, got %v", err)
	}
}
