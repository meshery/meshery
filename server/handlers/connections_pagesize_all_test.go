package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta3/environment"
)

// pageSize=all is what the cluster switcher, the header and the dashboard
// charts send to get every connection in one response. It used to be mapped
// to a single page of 100, so a fleet with more than 100 matching
// connections silently lost rows while totalCount reported the real number.
func TestGetConnectionsPageSizeAllReturnsEveryConnection(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "meshery.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.DBClose() })
	if err := db.AutoMigrate(connections.Connection{}, environment.Environment{}, environment.EnvironmentConnectionMapping{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}
	persister := &models.ConnectionPersister{DB: &db}

	const kubernetesCount = 250
	for i := 0; i < kubernetesCount; i++ {
		if _, err := persister.SaveConnection(&connections.Connection{Name: fmt.Sprintf("cluster-%03d", i), Kind: "kubernetes", ConnectionType: "platform", SubType: "orchestration", Status: connections.CONNECTED}); err != nil {
			t.Fatalf("save connection: %v", err)
		}
	}
	for i := 0; i < 3; i++ {
		if _, err := persister.SaveConnection(&connections.Connection{Name: fmt.Sprintf("prom-%d", i), Kind: "prometheus", ConnectionType: "observability", SubType: "monitoring", Status: connections.CONNECTED}); err != nil {
			t.Fatalf("save connection: %v", err)
		}
	}

	h := &Handler{config: &models.HandlerConfig{}, log: newTestLogger(t)}
	provider := &models.DefaultLocalProvider{ConnectionPersister: persister}

	cases := []struct {
		query          string
		wantReturned   int
		wantTotalCount int
	}{
		{"pageSize=all&kind=kubernetes", kubernetesCount, kubernetesCount},
		{"pagesize=all&kind=kubernetes", kubernetesCount, kubernetesCount},
		// "all" is the whole list, so a page number must not offset into it.
		{"page=1&pageSize=all&kind=kubernetes", kubernetesCount, kubernetesCount},
		{"page=0&pagesize=all", kubernetesCount + 3, kubernetesCount + 3},
		// Numeric page sizes keep the per-request cap.
		{"pageSize=500&kind=kubernetes", 100, kubernetesCount},
		{"pageSize=all&kind=grafana", 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections?"+tc.query, nil)
			rec := httptest.NewRecorder()
			h.GetConnections(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, body: %s", rec.Code, rec.Body.String())
			}

			if strings.Contains(rec.Body.String(), `"connections":null`) {
				t.Fatalf("connections encoded as null, want an empty list")
			}

			var page connections.ConnectionPage
			if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			t.Logf("returned=%d totalCount=%d page=%d pageSize=%d", len(page.Connections), page.TotalCount, page.Page, page.PageSize)

			if len(page.Connections) != tc.wantReturned || page.TotalCount != tc.wantTotalCount {
				t.Fatalf("returned=%d totalCount=%d, want returned=%d totalCount=%d", len(page.Connections), page.TotalCount, tc.wantReturned, tc.wantTotalCount)
			}
			seen := make(map[uuid.UUID]bool, len(page.Connections))
			for _, c := range page.Connections {
				if seen[c.ID] {
					t.Fatalf("connection %s returned twice", c.ID)
				}
				seen[c.ID] = true
			}
		})
	}

	// Rows are saved in ascending name and updated_at order, so a result
	// left in walk order fails both of these.
	t.Run("pageSize=all honours order", func(t *testing.T) {
		get := func(query string) []*connections.Connection {
			req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections?"+query, nil)
			rec := httptest.NewRecorder()
			h.GetConnections(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
			var page connections.ConnectionPage
			if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if len(page.Connections) != kubernetesCount+3 {
				t.Fatalf("returned=%d, want %d", len(page.Connections), kubernetesCount+3)
			}
			return page.Connections
		}

		byName := get("pageSize=all&order=name%20desc")
		if !sort.SliceIsSorted(byName, func(i, j int) bool { return byName[i].Name > byName[j].Name }) || byName[0].Name != "prom-2" {
			t.Fatalf("order=name desc not honoured: first=%q last=%q", byName[0].Name, byName[len(byName)-1].Name)
		}

		byUpdated := get("pageSize=all")
		for i := 1; i < len(byUpdated); i++ {
			if byUpdated[i].UpdatedAt.After(byUpdated[i-1].UpdatedAt) {
				t.Fatalf("default order is not updated_at desc at index %d", i)
			}
		}
		if !byUpdated[0].UpdatedAt.After(byUpdated[len(byUpdated)-1].UpdatedAt) {
			t.Fatalf("seeded updated_at values are not distinct; cannot check ordering")
		}
	})
}

// movingProvider serves connections ordered the way ConnectionPersister
// does and bumps one connection's updated_at right after the first page is
// read, as a status write landing mid-walk would.
type movingProvider struct {
	models.Provider
	conns []*connections.Connection
	moved *connections.Connection
	reads int
}

func (m *movingProvider) GetConnections(_ *http.Request, _ string, page, pageSize int, _, order, _ string, _, _, _ []string, _ string) (*connections.ConnectionPage, error) {
	sorted := append([]*connections.Connection(nil), m.conns...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if order == "created_at asc" {
			return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
		}
		return sorted[i].UpdatedAt.After(sorted[j].UpdatedAt)
	})
	start := min(page*pageSize, len(sorted))
	end := min(start+pageSize, len(sorted))
	m.reads++
	if m.reads == 1 {
		m.moved.UpdatedAt = time.Now().Add(time.Hour)
	}
	return &connections.ConnectionPage{Connections: sorted[start:end], TotalCount: len(sorted), Page: page, PageSize: pageSize}, nil
}

func TestGetConnectionsPageSizeAllSurvivesUpdateDuringWalk(t *testing.T) {
	const count = 250
	base := time.Now().Add(-24 * time.Hour)
	stub := &movingProvider{}
	for i := 0; i < count; i++ {
		ts := base.Add(time.Duration(i) * time.Second)
		stub.conns = append(stub.conns, &connections.Connection{ID: core.Uuid(uuid.Must(uuid.NewV4())), Name: fmt.Sprintf("cluster-%03d", i), CreatedAt: ts, UpdatedAt: ts})
	}
	// Under updated_at desc this connection sits at position 149, on the
	// second page, and jumps to position 0 once page 0 has been read.
	stub.moved = stub.conns[100]

	h := &Handler{config: &models.HandlerConfig{}, log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections?pageSize=all", nil)
	rec := httptest.NewRecorder()
	h.GetConnections(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, stub)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body: %s", rec.Code, rec.Body.String())
	}
	var page connections.ConnectionPage
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	got := make(map[core.Uuid]int, count)
	for _, c := range page.Connections {
		got[c.ID]++
	}
	for _, c := range stub.conns {
		if got[c.ID] != 1 {
			t.Fatalf("connection %s returned %d times, want exactly once", c.Name, got[c.ID])
		}
	}
	if len(page.Connections) != count {
		t.Fatalf("returned=%d, want %d", len(page.Connections), count)
	}
	if page.Connections[0].Name != stub.moved.Name {
		t.Fatalf("first connection = %q, want the most recently updated %q", page.Connections[0].Name, stub.moved.Name)
	}
}
