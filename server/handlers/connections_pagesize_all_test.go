package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
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
}
