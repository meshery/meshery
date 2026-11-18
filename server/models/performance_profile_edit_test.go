package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/meshery/schemas/models/core"
)

// TestSavePerformanceProfileEditKeepsCreatedAt stores a profile, then saves it
// again with the body the UI's performance form sends when a profile is edited:
// the ID and the form fields, with no createdAt. The creation time must survive
// the edit, in the stored row and in the save response.
func TestSavePerformanceProfileEditKeepsCreatedAt(t *testing.T) {
	db := newMigratedDB(t)
	provider := &DefaultLocalProvider{PerformanceProfilesPersister: &PerformanceProfilePersister{DB: db}}

	id := mustUUID(t)
	created := time.Date(2026, 11, 1, 9, 30, 0, 0, time.UTC)
	stored := PerformanceProfile{
		ID:                id,
		Name:              "original",
		LoadGenerators:    []string{"fortio"},
		Endpoints:         []string{"https://meshery.io"},
		ConcurrentRequest: 1,
		QPS:               5,
		Duration:          "30s",
		Metadata:          core.Map{},
		CreatedAt:         created,
		UpdatedAt:         created,
	}
	if err := db.Create(&stored).Error; err != nil {
		t.Fatalf("seed profile: %v", err)
	}

	// The body built by ui/components/performance/performance-helpers.tsx.
	edit := `{"id":"` + id.String() + `","name":"renamed","loadGenerators":["fortio"],` +
		`"endpoints":["https://meshery.io"],"serviceMesh":"","concurrentRequest":1,"qps":5,"duration":"30s",` +
		`"requestHeaders":"","requestBody":"","requestCookies":"","contentType":"",` +
		`"metadata":{"additional_options":[""],"ca_certificate":{}}}`
	var profile PerformanceProfile
	if err := json.Unmarshal([]byte(edit), &profile); err != nil {
		t.Fatalf("decode edit body: %v", err)
	}

	resp, err := provider.SavePerformanceProfile("", &profile)
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := provider.PerformanceProfilesPersister.GetPerformanceProfile(id)
	if err != nil {
		t.Fatalf("read profile back: %v", err)
	}
	if got.Name != "renamed" {
		t.Errorf("name = %q, want %q", got.Name, "renamed")
	}
	if !got.CreatedAt.Equal(created) {
		t.Errorf("stored created_at = %v, want %v to survive the edit", got.CreatedAt, created)
	}

	var saved PerformanceProfile
	if err := json.Unmarshal(resp, &saved); err != nil {
		t.Fatalf("decode response %s: %v", resp, err)
	}
	if !saved.CreatedAt.Equal(created) {
		t.Errorf("response createdAt = %v, want %v", saved.CreatedAt, created)
	}
}
