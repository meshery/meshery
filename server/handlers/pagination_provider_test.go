package handlers

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/mux"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
		{"page=-1&pagesize=abc", "", ""},
		{"page=&pagesize=", "", ""},
		{"page=2&pagesize=15", "2", "15"},
		{"page=2&pagesize=all", "2", "all"},
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

func (p *paginationSpyProvider) GetCatalogMesheryPatterns(_ string, page, size, _, _, _ string, _, _, _, _, _, _, _ []string) ([]byte, error) {
	p.page, p.size = page, size
	return []byte(`{}`), nil
}
func (p *paginationSpyProvider) GetCatalogMesheryFilters(_ string, page, size, _, _ string) ([]byte, error) {
	p.page, p.size = page, size
	return []byte(`{}`), nil
}

func TestCatalogPaginationLeadingPlus(t *testing.T) {
	for _, handler := range []func(*Handler, http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider){(*Handler).GetCatalogMesheryPatternsHandler, (*Handler).GetCatalogMesheryFiltersHandler} {
		provider := &paginationSpyProvider{}
		req := httptest.NewRequest(http.MethodGet, "/?pagesize=%2B15", nil)
		req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "token"))
		rec := httptest.NewRecorder()
		handler(&Handler{}, rec, req, nil, nil, provider)
		if rec.Code != http.StatusOK || provider.size != "15" {
			t.Fatalf("status=%d size=%q", rec.Code, provider.size)
		}
	}
}

func TestLocalProviderPaginationDefaults(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.SmiResultWithID{}, &models.PerformanceTestConfig{}); err != nil {
		t.Fatal(err)
	}
	store := &database.Handler{DB: db, Mutex: &sync.Mutex{}}
	provider := &models.DefaultLocalProvider{SmiResultPersister: &models.SMIResultsPersister{DB: store}, TestProfilesPersister: &models.TestProfilesPersister{DB: store}}
	for _, handler := range []struct {
		name string
		call func(http.ResponseWriter, *http.Request)
	}{
		{"smi", func(w http.ResponseWriter, r *http.Request) {
			(&Handler{}).FetchSmiResultsHandler(w, r, nil, nil, provider)
		}},
		{"profiles", func(w http.ResponseWriter, r *http.Request) {
			(&Handler{log: newTestLogger(t)}).UserTestPreferenceGet(w, r, nil, nil, provider)
		}},
	} {
		for _, tc := range []struct {
			query      string
			page, size uint64
		}{
			{"page=-1&pagesize=15", 0, 15}, {"page=0&pagesize=abc", 0, 10}, {"page=2&pagesize=all", 2, 10},
		} {
			t.Run(handler.name+"/"+tc.query, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/?"+tc.query, nil)
				rec := httptest.NewRecorder()
				handler.call(rec, req)
				var got struct {
					Page     uint64 `json:"page"`
					PageSize uint64 `json:"pageSize"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("decode response %s: %v", rec.Body.String(), err)
				}
				if rec.Code != http.StatusOK || got.Page != tc.page || got.PageSize != tc.size {
					t.Fatalf("status=%d page=%d size=%d; body=%s", rec.Code, got.Page, got.PageSize, rec.Body.String())
				}
			})
		}
	}
}

func (p *paginationSpyProvider) FetchSmiResults(_ *http.Request, page, size, _, _ string) ([]byte, error) {
	p.page, p.size = page, size
	return []byte(`{}`), nil
}
func (p *paginationSpyProvider) FetchSmiResult(_ *http.Request, page, size, _, _ string, _ core.Uuid) ([]byte, error) {
	p.page, p.size = page, size
	return []byte(`{}`), nil
}
func TestRemoteSMIPagination(t *testing.T) {
	for _, single := range []bool{false, true} {
		for _, tc := range []struct{ query, page, size string }{
			{"page=-1&pagesize=abc&pageSize=abc", "", ""},
			{"page=2&pagesize=all&pageSize=all", "2", "all"},
		} {
			p := &paginationSpyProvider{}
			req := httptest.NewRequest(http.MethodGet, "/?"+tc.query, nil)
			req = mux.SetURLVars(req, map[string]string{"id": "11111111-1111-1111-1111-111111111111"})
			rec := httptest.NewRecorder()
			if single {
				(&Handler{}).FetchSingleSmiResultHandler(rec, req, nil, nil, p)
			} else {
				(&Handler{}).FetchSmiResultsHandler(rec, req, nil, nil, p)
			}
			if rec.Code != http.StatusOK || p.page != tc.page || p.size != tc.size {
				t.Fatalf("single=%t status=%d page=%q size=%q", single, rec.Code, p.page, p.size)
			}
		}
	}
}
