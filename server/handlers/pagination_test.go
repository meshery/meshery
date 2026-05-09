package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models"
)

func TestGetPaginationParams(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		wantPage        int
		wantOffset      int
		wantLimit       int
		wantErrContains string
	}{
		{
			name:       "uses defaults when pagination is omitted",
			url:        "/api/system/database",
			wantPage:   0,
			wantOffset: 0,
			wantLimit:  defaultPageSize,
		},
		{
			name:       "accepts explicit page size",
			url:        "/api/system/database?page=2&pagesize=15",
			wantPage:   2,
			wantOffset: 30,
			wantLimit:  15,
		},
		{
			name:       "caps all to max page size",
			url:        "/api/system/database?pagesize=all",
			wantPage:   0,
			wantOffset: 0,
			wantLimit:  maxPageSize,
		},
		{
			name:            "rejects malformed page",
			url:             "/api/system/database?page=abc",
			wantErrContains: "page",
		},
		{
			name:            "rejects negative page",
			url:             "/api/system/database?page=-1",
			wantErrContains: "page",
		},
		{
			name:            "rejects malformed page size",
			url:             "/api/system/database?pagesize=abc",
			wantErrContains: "pagesize",
		},
		{
			name:            "rejects zero page size",
			url:             "/api/system/database?pagesize=0",
			wantErrContains: "pagesize",
		},
		{
			name:            "rejects oversized page size",
			url:             "/api/system/database?pagesize=101",
			wantErrContains: "pagesize",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)

			page, offset, limit, _, _, _, _, err := getPaginationParams(req)
			if tt.wantErrContains != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErrContains)
				}
				if !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Fatalf("expected error to contain %q, got %q", tt.wantErrContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if page != tt.wantPage {
				t.Fatalf("expected page %d, got %d", tt.wantPage, page)
			}
			if offset != tt.wantOffset {
				t.Fatalf("expected offset %d, got %d", tt.wantOffset, offset)
			}
			if limit != tt.wantLimit {
				t.Fatalf("expected limit %d, got %d", tt.wantLimit, limit)
			}
		})
	}
}

func TestGetSystemDatabase_InvalidPaginationReturnsJSON(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/system/database?pagesize=101", nil)
	rec := httptest.NewRecorder()

	h.GetSystemDatabase(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "pagesize") {
		t.Fatalf("expected response body to mention pagesize, got %q", rec.Body.String())
	}

	assertJSONErrorEnvelope(t, rec.Result())
}
