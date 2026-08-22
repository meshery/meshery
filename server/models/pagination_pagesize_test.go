package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
	workspace "github.com/meshery/schemas/models/v1beta3/workspace"
)

// TestGetWorkspacesPageSizeMatchesRequested verifies that GetWorkspaces returns
// the requested pageSize in the response, not len(fetched). Before the fix,
// PageSize was set to len(fetched), which broke the UI hasMore formula on the
// last page and caused infinite empty network requests.
func TestGetWorkspacesPageSizeMatchesRequested(t *testing.T) {
	dbHandler := newWorkspaceTestDB(t)
	if err := dbHandler.AutoMigrate(&workspace.Workspace{}); err != nil {
		t.Fatalf("failed to auto-migrate workspace: %v", err)
	}

	orgID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate org id: %v", err)
	}

	// Insert 13 workspaces so the last page has only 3 items when pageSize=10.
	now := time.Now().UTC()
	for i := 0; i < 13; i++ {
		id, _ := uuid.NewV4()
		ws := workspace.Workspace{
			ID:             id,
			Name:           "ws",
			OrganizationID: orgID,
			Metadata:       core.Map{},
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := dbHandler.Create(&ws).Error; err != nil {
			t.Fatalf("failed to seed workspace: %v", err)
		}
	}

	persister := &WorkspacePersister{DB: dbHandler}

	t.Run("last page - pageSize is requested, not len(fetched)", func(t *testing.T) {
		// Page 1 with pageSize=10 => only 3 rows fetched, but response must report pageSize=10.
		raw, err := persister.GetWorkspaces(orgID.String(), "", "", "1", "10", "")
		if err != nil {
			t.Fatalf("GetWorkspaces failed: %v", err)
		}
		var resp struct {
			PageSize   int `json:"pageSize"`
			TotalCount int `json:"totalCount"`
		}
		if err := json.Unmarshal(raw, &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if resp.PageSize != 10 {
			t.Errorf("pageSize = %d, want 10 (the requested size, not len(fetched))", resp.PageSize)
		}
		if resp.TotalCount != 13 {
			t.Errorf("totalCount = %d, want 13", resp.TotalCount)
		}
	})

	t.Run("pageSize=all sets pageSize to totalCount", func(t *testing.T) {
		raw, err := persister.GetWorkspaces(orgID.String(), "", "", "0", "all", "")
		if err != nil {
			t.Fatalf("GetWorkspaces failed: %v", err)
		}
		var resp struct {
			PageSize   int `json:"pageSize"`
			TotalCount int `json:"totalCount"`
		}
		if err := json.Unmarshal(raw, &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if resp.PageSize != resp.TotalCount {
			t.Errorf("pageSize=all: pageSize=%d, totalCount=%d, want them equal", resp.PageSize, resp.TotalCount)
		}
	})
}
