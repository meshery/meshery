package design

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

func TestSelectDesignByName(t *testing.T) {
	origPrompt := runDesignSelectPrompt
	t.Cleanup(func() { runDesignSelectPrompt = origPrompt })

	first := map[string]interface{}{"name": "nginx", "id": "11111111-1111-4111-8111-111111111111"}
	second := map[string]interface{}{"name": "nginx", "id": "22222222-2222-4222-8222-222222222222"}
	upper := map[string]interface{}{"name": "NGINX", "id": "33333333-3333-4333-8333-333333333333"}
	partial := map[string]interface{}{"name": "nginx-staging", "id": "44444444-4444-4444-8444-444444444444"}
	malformed := map[string]interface{}{"name": 12345}

	tests := []struct {
		name      string
		patterns  []interface{}
		query     string
		want      map[string]interface{}
		wantError bool
		prompt    bool
	}{
		{
			name:      "exact before substring",
			patterns:  []interface{}{partial, first},
			query:     "nginx",
			want:      first,
			wantError: false,
			prompt:    false,
		},
		{
			name:      "exact before case-insensitive",
			patterns:  []interface{}{upper, first},
			query:     "nginx",
			want:      first,
			wantError: false,
			prompt:    false,
		},
		{
			name:      "case-insensitive fallback",
			patterns:  []interface{}{upper},
			query:     "nginx",
			want:      upper,
			wantError: false,
			prompt:    false,
		},
		{
			name:      "substring rejected",
			patterns:  []interface{}{partial},
			query:     "nginx",
			want:      nil,
			wantError: true,
			prompt:    false,
		},
		{
			name:      "duplicate exact prompts and selects chosen index",
			patterns:  []interface{}{first, second},
			query:     "nginx",
			want:      second,
			wantError: false,
			prompt:    true,
		},
		{
			name:      "malformed and nil entries are skipped",
			patterns:  []interface{}{nil, malformed, first},
			query:     "nginx",
			want:      first,
			wantError: false,
			prompt:    false,
		},
		{
			name:      "empty patterns returns not found error",
			patterns:  []interface{}{},
			query:     "nginx",
			want:      nil,
			wantError: true,
			prompt:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			runDesignSelectPrompt = func(_ string, _ []string) (int, error) {
				called = true
				return 1, nil
			}

			got, err := selectDesignByName(tt.patterns, tt.query)
			if (err != nil) != tt.wantError {
				t.Fatalf("selectDesignByName() error = %v, wantError = %v", err, tt.wantError)
			}
			if tt.wantError {
				if meshkiterrors.GetCode(err) != ErrDesignsNotFoundCode {
					t.Fatalf("expected error code %s, got %v", ErrDesignsNotFoundCode, err)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("selectDesignByName() = %v, want %v", got, tt.want)
			}
			if called != tt.prompt {
				t.Fatalf("prompt called = %v, want %v", called, tt.prompt)
			}
		})
	}
}

func TestFetchDesignSearchPages(t *testing.T) {
	calls := 0
	got, err := fetchDesignSearch("istio & ingress", func(raw string) (*map[string]interface{}, error) {
		u, parseErr := url.Parse(raw)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		if u.Query().Get("search") != "istio & ingress" || u.Query().Get("page") != fmt.Sprint(calls) {
			t.Fatalf("unexpected query %s", raw)
		}
		if u.Query().Get("pagesize") != "10000" || u.Query().Get("populate") != "pattern_file" {
			t.Fatalf("missing expected query parameters in %s", raw)
		}

		calls++
		if calls == 1 {
			page := map[string]interface{}{
				"pageSize":   float64(1),
				"totalCount": float64(2),
				"patterns":   []interface{}{map[string]interface{}{"name": "partial"}},
			}
			return &page, nil
		}

		page := map[string]interface{}{
			"pageSize":   float64(1),
			"totalCount": float64(2),
			"patterns":   []interface{}{map[string]interface{}{"name": "istio & ingress"}},
		}
		return &page, nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 page calls, got %d", calls)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(got))
	}
	name, _ := got[1].(map[string]interface{})["name"].(string)
	if name != "istio & ingress" {
		t.Fatalf("expected second pattern name 'istio & ingress', got %q", name)
	}
}
