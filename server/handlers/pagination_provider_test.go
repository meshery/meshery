package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/models"
)

type paginationSpyProvider struct {
	models.Provider
	page, size string
}

func (p *paginationSpyProvider) GetWorkspaces(_ string, page string, size string, _, _, _, _ string) ([]byte, error) {
	p.page, p.size = page, size
	return []byte(`{}`), nil
}

func (p *paginationSpyProvider) GetEnvironments(_ string, page string, size string, _, _, _, _ string) ([]byte, error) {
	p.page, p.size = page, size
	return []byte(`{}`), nil
}

func (p *paginationSpyProvider) GetMesheryPatterns(_ string, page string, size string, _, _, _ string, _ []string, _ string, _ []string) ([]byte, error) {
	p.page, p.size = page, size
	return []byte(`{}`), nil
}

func (p *paginationSpyProvider) GetMesheryFilters(_ string, page string, size string, _, _ string, _ []string) ([]byte, error) {
	p.page, p.size = page, size
	return []byte(`{}`), nil
}

func TestProviderPaginationFallback(t *testing.T) {
	cases := []struct{ query, page, size string }{
		{"page=-1&pagesize=abc", "", ""},
		{"page=abc&pagesize=-1", "", ""},
		{"page=4294967296&pagesize=4294967296", "", ""},
		{"page=2&pagesize=15", "2", "15"},
		{"page=0&pagesize=all", "0", "all"},
		{"page=&pagesize=", "", ""},
	}
	endpoints := []struct {
		name    string
		handler func(*Handler, http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)
	}{
		{"workspaces", (*Handler).GetWorkspacesHandler},
		{"environments", (*Handler).GetEnvironments},
		{"patterns", (*Handler).GetMesheryPatternsHandler},
		{"filters", (*Handler).GetMesheryFiltersHandler},
	}
	for _, endpoint := range endpoints {
		for _, tc := range cases {
			t.Run(endpoint.name+"/"+tc.query, func(t *testing.T) {
				provider := &paginationSpyProvider{}
				req := httptest.NewRequest(http.MethodGet, "/?orgId=test&"+tc.query, nil)
				req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "token"))
				rec := httptest.NewRecorder()
				endpoint.handler(&Handler{}, rec, req, nil, nil, provider)
				if rec.Code != http.StatusOK || provider.page != tc.page || provider.size != tc.size {
					t.Fatalf("status=%d page=%q size=%q; want 200, %q, %q", rec.Code, provider.page, provider.size, tc.page, tc.size)
				}
			})
		}
	}
}

func (p *paginationSpyProvider) SMPTestConfigFetch(_ *http.Request, page, size, _, _ string) ([]byte, error) {
	p.page, p.size = page, size
	return []byte(`{}`), nil
}

func TestTestPreferencePaginationFallback(t *testing.T) {
	for _, tc := range []struct{ query, page, size string }{
		{"page=-1&pagesize=abc", "0", "10"},
		{"page=&pagesize=", "0", "10"},
		{"page=2&pagesize=15", "2", "15"},
	} {
		t.Run(tc.query, func(t *testing.T) {
			provider := &paginationSpyProvider{}
			req := httptest.NewRequest(http.MethodGet, "/?"+tc.query, nil)
			rec := httptest.NewRecorder()
			h := &Handler{log: newTestLogger(t)}
			h.UserTestPreferenceGet(rec, req, nil, nil, provider)
			if rec.Code != http.StatusOK || provider.page != tc.page || provider.size != tc.size {
				t.Fatalf("status=%d page=%q size=%q; want 200, %q, %q", rec.Code, provider.page, provider.size, tc.page, tc.size)
			}
		})
	}
}
