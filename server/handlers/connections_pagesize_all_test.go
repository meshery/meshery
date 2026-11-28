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

func newPageSizeAllPersister(t *testing.T) (*models.ConnectionPersister, *database.Handler) {
	t.Helper()
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "meshery.db")})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.DBClose() })
	if err := db.AutoMigrate(connections.Connection{}, environment.Environment{}, environment.EnvironmentConnectionMapping{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}
	return &models.ConnectionPersister{DB: &db}, &db
}

func getAllPage(t *testing.T, provider models.Provider, query string) (*httptest.ResponseRecorder, connections.ConnectionPage) {
	t.Helper()
	h := &Handler{config: &models.HandlerConfig{}, log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections?"+query, nil)
	rec := httptest.NewRecorder()
	h.GetConnections(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	var page connections.ConnectionPage
	if rec.Code == http.StatusOK {
		if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}
	return rec, page
}

// statusFlipProvider disconnects one connection right after the first page
// is read, as a status write landing mid-walk would.
type statusFlipProvider struct {
	*models.DefaultLocalProvider
	flip  core.Uuid
	reads int
}

func (p *statusFlipProvider) GetConnections(req *http.Request, userID string, page, pageSize int, search, order, filter string, status, kind, connType []string, name string) (*connections.ConnectionPage, error) {
	res, err := p.DefaultLocalProvider.GetConnections(req, userID, page, pageSize, search, order, filter, status, kind, connType, name)
	p.reads++
	if p.reads == 1 {
		if _, uerr := p.ConnectionPersister.UpdateConnectionStatusByID(p.flip, connections.DISCONNECTED); uerr != nil {
			return nil, uerr
		}
	}
	return res, err
}

// A connection leaving a status filter mid-walk must not shift the offsets of
// the connections after it: every connection that stays connected for the
// whole request is returned.
func TestGetConnectionsPageSizeAllSurvivesStatusChangeDuringWalk(t *testing.T) {
	persister, _ := newPageSizeAllPersister(t)
	const count = 250
	var saved []*connections.Connection
	for i := 0; i < count; i++ {
		c, err := persister.SaveConnection(&connections.Connection{Name: fmt.Sprintf("cluster-%03d", i), Kind: "kubernetes", ConnectionType: "platform", SubType: "orchestration", Status: connections.CONNECTED})
		if err != nil {
			t.Fatalf("save connection: %v", err)
		}
		saved = append(saved, c)
	}
	// The oldest connection sits on the first created_at page.
	stub := &statusFlipProvider{DefaultLocalProvider: &models.DefaultLocalProvider{ConnectionPersister: persister}, flip: saved[0].ID}

	rec, page := getAllPage(t, stub, "pageSize=all&status=connected&kind=kubernetes")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body: %s", rec.Code, rec.Body.String())
	}
	got := make(map[core.Uuid]int, count)
	for _, c := range page.Connections {
		got[c.ID]++
		if c.Status != connections.CONNECTED {
			t.Fatalf("connection %s has status %q, want only connected", c.Name, c.Status)
		}
	}
	for _, c := range saved[1:] {
		if got[c.ID] != 1 {
			t.Fatalf("connection %s returned %d times, want exactly once", c.Name, got[c.ID])
		}
	}
	// The flipped connection was read while still connected, so it is part of
	// the snapshot; what must not happen is a later connection going missing.
	if page.TotalCount != len(page.Connections) {
		t.Fatalf("totalCount=%d, want the %d rows returned", page.TotalCount, len(page.Connections))
	}
}

// Connections created in the same instant must still be walked without
// overlap or gaps between pages.
func TestGetConnectionsPageSizeAllWithTiedCreatedAt(t *testing.T) {
	persister, db := newPageSizeAllPersister(t)
	const count = 250
	ts := time.Now().Add(-time.Hour).UTC()
	for i := 0; i < count; i++ {
		c := &connections.Connection{ID: core.Uuid(uuid.Must(uuid.NewV4())), Name: fmt.Sprintf("cluster-%03d", i), Kind: "kubernetes", Status: connections.CONNECTED, CreatedAt: ts, UpdatedAt: ts}
		if err := db.Create(c).Error; err != nil {
			t.Fatalf("create connection: %v", err)
		}
	}

	rec, page := getAllPage(t, &models.DefaultLocalProvider{ConnectionPersister: persister}, "pageSize=all")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body: %s", rec.Code, rec.Body.String())
	}
	seen := make(map[core.Uuid]bool, count)
	for _, c := range page.Connections {
		seen[c.ID] = true
	}
	if len(page.Connections) != count || len(seen) != count {
		t.Fatalf("returned=%d unique=%d, want %d", len(page.Connections), len(seen), count)
	}
}

// scriptedProvider serves a fixed first page and then either fails or, when
// stuck is set, keeps returning that same page, as a provider that ignores the
// page number would.
type scriptedProvider struct {
	models.Provider
	first []*connections.Connection
	total int
	stuck bool
	reads int
}

func (s *scriptedProvider) GetConnections(_ *http.Request, _ string, page, pageSize int, _, _, _ string, _, _, _ []string, _ string) (*connections.ConnectionPage, error) {
	s.reads++
	if s.reads > 5 {
		return nil, fmt.Errorf("walk did not stop after %d reads", s.reads-1)
	}
	if page > 0 && !s.stuck {
		return nil, fmt.Errorf("provider unavailable")
	}
	return &connections.ConnectionPage{Connections: s.first, TotalCount: s.total, Page: page, PageSize: pageSize}, nil
}

func scriptedConnections(n int) []*connections.Connection {
	out := make([]*connections.Connection, n)
	for i := range out {
		out[i] = &connections.Connection{ID: core.Uuid(uuid.Must(uuid.NewV4())), Name: fmt.Sprintf("cluster-%03d", i), Status: connections.CONNECTED}
	}
	return out
}

// A provider error after the first page must fail the request, not return
// the rows read so far as a successful partial list.
func TestGetConnectionsPageSizeAllProviderErrorMidWalk(t *testing.T) {
	stub := &scriptedProvider{first: scriptedConnections(100), total: 250}
	rec, _ := getAllPage(t, stub, "pageSize=all")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"connections"`) {
		t.Fatalf("error response carries a partial connections list: %s", rec.Body.String())
	}
	if stub.reads != 2 {
		t.Fatalf("provider reads = %d, want 2", stub.reads)
	}
}

// A provider that ignores the page number keeps returning the same full page;
// the walk must stop once a page adds nothing new instead of looping.
func TestGetConnectionsPageSizeAllStopsWhenPageAddsNothing(t *testing.T) {
	stub := &scriptedProvider{first: scriptedConnections(100), total: 250, stuck: true}
	rec, page := getAllPage(t, stub, "pageSize=all")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body: %s", rec.Code, rec.Body.String())
	}
	if stub.reads != 2 || len(page.Connections) != 100 {
		t.Fatalf("reads=%d returned=%d, want reads=2 returned=100", stub.reads, len(page.Connections))
	}
}
