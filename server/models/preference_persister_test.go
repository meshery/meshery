package models

import (
	"encoding/json"
	"testing"

	"github.com/meshery/meshkit/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestSessionPreferencePersister(t *testing.T) *SessionPreferencePersister {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	if err := db.AutoMigrate(&UserPreference{}); err != nil {
		t.Fatalf("failed to migrate user preference schema: %v", err)
	}

	return &SessionPreferencePersister{
		DB: &database.Handler{DB: db},
	}
}

func TestPreferenceMarshalJSONEmitsCanonicalAndLegacySelectedOrganizationID(t *testing.T) {
	pref := Preference{
		SelectedOrganizationID: "00000000-0000-0000-0000-000000000001",
		Grafana: &Grafana{
			GrafanaURL:    "http://grafana.example",
			GrafanaAPIKey: "grafana-secret",
		},
		Prometheus: &Prometheus{
			PrometheusURL: "http://prometheus.example",
		},
	}

	data, err := json.Marshal(pref)
	if err != nil {
		t.Fatalf("marshal preference: %v", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal preference: %v", err)
	}

	if got["selectedOrganizationId"] != pref.SelectedOrganizationID {
		t.Fatalf("expected canonical selectedOrganizationId, got %q", got["selectedOrganizationId"])
	}
	if got["selectedOrganizationID"] != pref.SelectedOrganizationID {
		t.Fatalf("expected legacy selectedOrganizationID, got %q", got["selectedOrganizationID"])
	}

	grafana, ok := got["grafana"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected grafana object, got %+v", got["grafana"])
	}
	if grafana["grafanaUrl"] != pref.Grafana.GrafanaURL {
		t.Fatalf("expected canonical grafanaUrl, got %q", grafana["grafanaUrl"])
	}
	if grafana["grafanaApiKey"] != pref.Grafana.GrafanaAPIKey {
		t.Fatalf("expected canonical grafanaApiKey, got %q", grafana["grafanaApiKey"])
	}

	prometheus, ok := got["prometheus"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected prometheus object, got %+v", got["prometheus"])
	}
	if prometheus["prometheusUrl"] != pref.Prometheus.PrometheusURL {
		t.Fatalf("expected canonical prometheusUrl, got %q", prometheus["prometheusUrl"])
	}
}

func TestSessionPreferencePersisterReadFromPersisterReturnsDefaultsWhenMissing(t *testing.T) {
	persister := newTestSessionPreferencePersister(t)

	pref, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("expected missing preference to return defaults, got error: %v", err)
	}

	if pref == nil {
		t.Fatal("expected default preference, got nil")
	}

	if !pref.AnonymousUsageStats {
		t.Error("expected AnonymousUsageStats to default to true")
	}

	if !pref.AnonymousPerfResults {
		t.Error("expected AnonymousPerfResults to default to true")
	}
}

func TestSessionPreferencePersisterReadFromPersisterReturnsPersistedPreference(t *testing.T) {
	persister := newTestSessionPreferencePersister(t)

	want := &Preference{
		AnonymousUsageStats:    false,
		AnonymousPerfResults:   false,
		SelectedOrganizationID: "org-1",
		DashboardPreferences: map[string]interface{}{
			"theme": "dark",
		},
	}

	if err := persister.WriteToPersister("user-1", want); err != nil {
		t.Fatalf("failed to persist preference: %v", err)
	}

	got, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("failed to read persisted preference: %v", err)
	}

	if got == nil {
		t.Fatal("expected persisted preference, got nil")
	}

	if got.AnonymousUsageStats != want.AnonymousUsageStats {
		t.Errorf("got AnonymousUsageStats=%v, want %v", got.AnonymousUsageStats, want.AnonymousUsageStats)
	}

	if got.AnonymousPerfResults != want.AnonymousPerfResults {
		t.Errorf("got AnonymousPerfResults=%v, want %v", got.AnonymousPerfResults, want.AnonymousPerfResults)
	}

	if got.SelectedOrganizationID != want.SelectedOrganizationID {
		t.Errorf("got SelectedOrganizationID=%q, want %q", got.SelectedOrganizationID, want.SelectedOrganizationID)
	}

	if got.DashboardPreferences["theme"] != want.DashboardPreferences["theme"] {
		t.Errorf("got DashboardPreferences[theme]=%v, want %v", got.DashboardPreferences["theme"], want.DashboardPreferences["theme"])
	}
}

func TestSessionPreferencePersisterReadFromPersisterReturnsDefaultsForStoredNull(t *testing.T) {
	persister := newTestSessionPreferencePersister(t)

	if err := persister.DB.Model(&UserPreference{}).Create(&UserPreference{
		ID:              "user-1",
		PreferenceBytes: []byte("null"),
	}).Error; err != nil {
		t.Fatalf("failed to seed null preference: %v", err)
	}

	pref, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("expected stored null preference to return defaults, got error: %v", err)
	}

	if pref == nil {
		t.Fatal("expected default preference for stored null, got nil")
	}

	if !pref.AnonymousUsageStats {
		t.Error("expected AnonymousUsageStats to default to true for stored null preference")
	}

	if !pref.AnonymousPerfResults {
		t.Error("expected AnonymousPerfResults to default to true for stored null preference")
	}
}
