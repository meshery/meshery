package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/v1beta2/organization"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newLocalProviderWithOrgs returns a Local provider whose organizations table
// holds n organizations, more than the default page size of 10 when n > 10.
func newLocalProviderWithOrgs(t *testing.T, n int) *models.DefaultLocalProvider {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&organization.Organization{}); err != nil {
		t.Fatalf("migrate organizations: %v", err)
	}

	persister := &models.OrganizationPersister{DB: &database.Handler{DB: db}}
	for i := range n {
		org := &organization.Organization{ID: uuid.Must(uuid.NewV4()), Name: fmt.Sprintf("org-%02d", i)}
		if _, err := persister.SaveOrganization(org); err != nil {
			t.Fatalf("seed organization %d: %v", i, err)
		}
	}

	return &models.DefaultLocalProvider{OrganizationPersister: persister}
}

// TestGetOrganizations_LocalProviderPageSize pins how the Local provider reads
// the pagesize query parameter. "all" is what `mesheryctl organization
// --count` sends, and the value the workspace and environment listings already
// accept; it used to fail strconv.ParseUint and surface as a 502.
func TestGetOrganizations_LocalProviderPageSize(t *testing.T) {
	const seeded = 12

	cases := []struct {
		name         string
		query        string
		wantOrgs     int
		wantPageSize uint64
	}{
		{name: "all returns every organization", query: "?pagesize=all", wantOrgs: seeded, wantPageSize: seeded},
		{name: "numeric page size still paginates", query: "?pagesize=5", wantOrgs: 5, wantPageSize: 5},
		{name: "missing page size defaults to 10", query: "", wantOrgs: 10, wantPageSize: 10},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newLocalProviderWithOrgs(t, seeded)

			req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs"+tc.query, nil))
			rec := httptest.NewRecorder()

			h.GetOrganizations(rec, req, nil, nil, provider)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusOK, rec.Body.String())
			}
			var page models.OrganizationsPage
			if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
				t.Fatalf("decode response: %v (body=%q)", err, rec.Body.String())
			}
			if page.TotalCount != seeded {
				t.Errorf("totalCount = %d, want %d", page.TotalCount, seeded)
			}
			if len(page.Organizations) != tc.wantOrgs {
				t.Errorf("got %d organizations, want %d", len(page.Organizations), tc.wantOrgs)
			}
			if page.PageSize != tc.wantPageSize {
				t.Errorf("pageSize = %d, want %d", page.PageSize, tc.wantPageSize)
			}
		})
	}
}
