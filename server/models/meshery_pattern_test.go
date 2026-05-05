package models

import (
	"encoding/json"
	"testing"
)

// TestMesheryPattern_UnmarshalAcceptsBothCountSpellings locks in the
// Phase 2.K cascade dual-accept behavior: the five count fields on
// MesheryPattern (viewCount / shareCount / downloadCount / cloneCount /
// deploymentCount) must accept both the canonical camelCase and the
// legacy snake_case spellings on inbound JSON, because remote providers
// (meshery-cloud, Kanvas) may still emit snake_case during the
// deprecation window.
func TestMesheryPattern_UnmarshalAcceptsBothCountSpellings(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "canonical camelCase",
			body: `{"viewCount":1,"shareCount":2,"downloadCount":3,"cloneCount":4,"deploymentCount":5}`,
		},
		{
			name: "legacy snake_case",
			body: `{"view_count":1,"share_count":2,"download_count":3,"clone_count":4,"deployment_count":5}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got MesheryPattern
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if got.ViewCount != 1 || got.ShareCount != 2 || got.DownloadCount != 3 || got.CloneCount != 4 || got.DeploymentCount != 5 {
				t.Errorf("unexpected counts: view=%d share=%d download=%d clone=%d deployment=%d",
					got.ViewCount, got.ShareCount, got.DownloadCount, got.CloneCount, got.DeploymentCount)
			}
		})
	}
}

// TestMesheryPattern_UnmarshalCanonicalWinsOverLegacy verifies the
// documented precedence: when both camelCase and snake_case spellings
// of a count field are present in the same payload, the canonical
// camelCase value wins. This matches the precedence contract adopted
// across the Phase 2.K cascade (see MesheryPatternPOSTRequestBody and
// MesheryPatternUPDATERequestBody for the same pattern).
func TestMesheryPattern_UnmarshalCanonicalWinsOverLegacy(t *testing.T) {
	body := `{"viewCount":10,"view_count":99,"shareCount":20,"share_count":88}`
	var got MesheryPattern
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if got.ViewCount != 10 {
		t.Errorf("ViewCount: canonical should win; want 10, got %d", got.ViewCount)
	}
	if got.ShareCount != 20 {
		t.Errorf("ShareCount: canonical should win; want 20, got %d", got.ShareCount)
	}
}
