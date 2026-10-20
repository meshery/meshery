package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta2/organization"
)

// TestOrganizationsPage_MarshalEmitsBothKeyFlavors locks in the
// deprecation-window contract: MarshalJSON emits both the canonical
// camelCase (`totalCount`, `pageSize`) AND the legacy snake_case
// (`total_count`, `page_size`) spellings so external consumers on
// either vocabulary keep working while they migrate.
func TestOrganizationsPage_MarshalEmitsBothKeyFlavors(t *testing.T) {
	p := OrganizationsPage{
		TotalCount: 7,
		Page:       2,
		PageSize:   25,
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out := string(b)
	for _, want := range []string{`"totalCount":7`, `"total_count":7`, `"pageSize":25`, `"page_size":25`, `"page":2`} {
		if !strings.Contains(out, want) {
			t.Errorf("marshal output missing %q; got %s", want, out)
		}
	}
}

// TestOrganizationsPage_UnmarshalAcceptsEitherKeyFlavor verifies that
// either wire form round-trips back into the struct. Canonical wins if
// both are present.
func TestOrganizationsPage_UnmarshalAcceptsEitherKeyFlavor(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		wantTotal  int
		wantPageSz uint64
	}{
		{"canonical only", `{"totalCount":5,"page":1,"pageSize":10}`, 5, 10},
		{"legacy only", `{"total_count":5,"page":1,"page_size":10}`, 5, 10},
		{"canonical wins when both present", `{"totalCount":9,"total_count":1,"page":1,"pageSize":50,"page_size":5}`, 9, 50},
		{"both flavors absent zeroes the fields", `{"page":1}`, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got OrganizationsPage
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.TotalCount != tc.wantTotal {
				t.Errorf("TotalCount = %d, want %d", got.TotalCount, tc.wantTotal)
			}
			if got.PageSize != tc.wantPageSz {
				t.Errorf("PageSize = %d, want %d", got.PageSize, tc.wantPageSz)
			}
		})
	}
}

// TestOrganizationsPage_UnmarshalResetsFieldsOnReuse locks in std-library
// json.Unmarshal semantics: when the destination struct is reused across
// decodes and the second input omits TotalCount / PageSize entirely, the
// destination's fields must reset to zero rather than carry stale values
// forward from the prior decode.
func TestOrganizationsPage_UnmarshalResetsFieldsOnReuse(t *testing.T) {
	var p OrganizationsPage
	if err := json.Unmarshal([]byte(`{"totalCount":99,"pageSize":77,"page":1}`), &p); err != nil {
		t.Fatalf("first unmarshal: %v", err)
	}
	if p.TotalCount != 99 || p.PageSize != 77 {
		t.Fatalf("prime decode wrong: got TotalCount=%d PageSize=%d", p.TotalCount, p.PageSize)
	}
	if err := json.Unmarshal([]byte(`{"page":2}`), &p); err != nil {
		t.Fatalf("second unmarshal: %v", err)
	}
	if p.TotalCount != 0 {
		t.Errorf("TotalCount leaked stale value %d across reuse", p.TotalCount)
	}
	if p.PageSize != 0 {
		t.Errorf("PageSize leaked stale value %d across reuse", p.PageSize)
	}
	if p.Page != 2 {
		t.Errorf("Page = %d, want 2", p.Page)
	}
}

// TestOrganizationPersister_GetOrganizations_Filter verifies that
// OrganizationPersister.GetOrganizations and DefaultLocalProvider.GetOrganizations
// process `filter` parameters dynamically via utils.ApplyFilters rather than
// incorrectly treating the parameter as an updatedAfter timestamp cutoff.
func TestOrganizationPersister_GetOrganizations_Filter(t *testing.T) {
	db := newMigratedDB(t)
	op := &OrganizationPersister{DB: db}
	lp := &DefaultLocalProvider{OrganizationPersister: op}

	owner1 := uuid.Must(uuid.NewV4())
	owner2 := uuid.Must(uuid.NewV4())

	org1 := &organization.Organization{
		ID:          uuid.Must(uuid.NewV4()),
		Name:        "AlphaOrg",
		Description: "Alpha organization",
		Country:     "US",
		Region:      "West",
		Owner:       owner1,
	}
	org2 := &organization.Organization{
		ID:          uuid.Must(uuid.NewV4()),
		Name:        "BetaOrg",
		Description: "Beta organization",
		Country:     "CA",
		Region:      "East",
		Owner:       owner2,
	}

	if _, err := op.SaveOrganization(org1); err != nil {
		t.Fatalf("SaveOrganization(org1): %v", err)
	}
	if _, err := op.SaveOrganization(org2); err != nil {
		t.Fatalf("SaveOrganization(org2): %v", err)
	}

	t.Run("unfiltered returns all organizations", func(t *testing.T) {
		raw, err := op.GetOrganizations("", "", 0, 10, "")
		if err != nil {
			t.Fatalf("GetOrganizations: %v", err)
		}
		var page OrganizationsPage
		if err := json.Unmarshal(raw, &page); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if page.TotalCount != 2 || len(page.Organizations) != 2 {
			t.Fatalf("got %d orgs, want 2", page.TotalCount)
		}
	})

	t.Run("filter by name returns matching organization", func(t *testing.T) {
		raw, err := op.GetOrganizations("", "", 0, 10, "name AlphaOrg")
		if err != nil {
			t.Fatalf("GetOrganizations: %v", err)
		}
		var page OrganizationsPage
		if err := json.Unmarshal(raw, &page); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if page.TotalCount != 1 || len(page.Organizations) != 1 {
			t.Fatalf("got %d orgs, want 1", page.TotalCount)
		}
		if page.Organizations[0].Name != "AlphaOrg" {
			t.Errorf("got name %q, want AlphaOrg", page.Organizations[0].Name)
		}
	})

	t.Run("filter by owner returns matching organization", func(t *testing.T) {
		raw, err := op.GetOrganizations("", "", 0, 10, fmt.Sprintf("owner %s", owner2))
		if err != nil {
			t.Fatalf("GetOrganizations: %v", err)
		}
		var page OrganizationsPage
		if err := json.Unmarshal(raw, &page); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if page.TotalCount != 1 || len(page.Organizations) != 1 {
			t.Fatalf("got %d orgs, want 1", page.TotalCount)
		}
		if page.Organizations[0].Name != "BetaOrg" {
			t.Errorf("got name %q, want BetaOrg", page.Organizations[0].Name)
		}
	})

	t.Run("arbitrary filter value does not corrupt timestamp query", func(t *testing.T) {
		// Previously, passing "owner" would execute `WHERE updated_at > 'owner'`.
		// With utils.ApplyFilters, a filter with no matching key/value simply returns 0 matching results without SQL timestamp errors.
		raw, err := op.GetOrganizations("", "", 0, 10, "owner non-existent-id")
		if err != nil {
			t.Fatalf("GetOrganizations should not fail on non-timestamp filter: %v", err)
		}
		var page OrganizationsPage
		if err := json.Unmarshal(raw, &page); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if page.TotalCount != 0 {
			t.Fatalf("got %d orgs, want 0", page.TotalCount)
		}
	})

	t.Run("search by substring works", func(t *testing.T) {
		raw, err := op.GetOrganizations("beta", "", 0, 10, "")
		if err != nil {
			t.Fatalf("GetOrganizations: %v", err)
		}
		var page OrganizationsPage
		if err := json.Unmarshal(raw, &page); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if page.TotalCount != 1 || page.Organizations[0].Name != "BetaOrg" {
			t.Fatalf("search failed, got: %+v", page)
		}
	})

	t.Run("ordering works", func(t *testing.T) {
		raw, err := op.GetOrganizations("", "name desc", 0, 10, "")
		if err != nil {
			t.Fatalf("GetOrganizations: %v", err)
		}
		var page OrganizationsPage
		if err := json.Unmarshal(raw, &page); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if len(page.Organizations) != 2 || page.Organizations[0].Name != "BetaOrg" {
			t.Fatalf("ordering failed, got first: %s", page.Organizations[0].Name)
		}
	})

	t.Run("pagination works", func(t *testing.T) {
		raw, err := op.GetOrganizations("", "name asc", 0, 1, "")
		if err != nil {
			t.Fatalf("GetOrganizations: %v", err)
		}
		var page OrganizationsPage
		if err := json.Unmarshal(raw, &page); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if page.TotalCount != 2 || len(page.Organizations) != 1 || page.Organizations[0].Name != "AlphaOrg" {
			t.Fatalf("pagination failed, got: %+v", page)
		}
	})

	t.Run("DefaultLocalProvider delegates filter correctly", func(t *testing.T) {
		raw, err := lp.GetOrganizations("", "0", "10", "", "", "name AlphaOrg")
		if err != nil {
			t.Fatalf("lp.GetOrganizations: %v", err)
		}
		var page OrganizationsPage
		if err := json.Unmarshal(raw, &page); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if page.TotalCount != 1 || page.Organizations[0].Name != "AlphaOrg" {
			t.Fatalf("DefaultLocalProvider filter delegation failed: %+v", page)
		}
	})
}
