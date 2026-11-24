package models

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

func TestDecodePerformanceTestConfigSupportsLegacyStoredJSON(t *testing.T) {
	data := []byte(`{
		"smp_version": "v0.7.1",
		"id": "legacy-id",
		"name": "legacy config",
		"labels": {"env": "dev"},
		"clients": [{
			"internal": true,
			"load_generator": "fortio",
			"protocol": 1,
			"connections": 3,
			"rps": 7,
			"headers": {"x-test": "yes"},
			"cookies": {"session": "abc"},
			"body": "{}",
			"content_type": "application/json",
			"endpoint_urls": ["https://meshery.io"],
			"ssl_certificate": "cert",
			"additional_options": "{\"quiet\":true}"
		}],
		"duration": "30s"
	}`)

	config, err := decodePerformanceTestConfig(data)
	if err != nil {
		t.Fatalf("decodePerformanceTestConfig returned error: %v", err)
	}

	if config.SmpVersion != "v0.7.1" || config.ID != "legacy-id" || config.Name != "legacy config" {
		t.Fatalf("decoded top-level fields incorrectly: %+v", config)
	}
	if len(config.Clients) != 1 {
		t.Fatalf("expected 1 client, got %d", len(config.Clients))
	}
	client := config.Clients[0]
	if client.LoadGenerator != "fortio" || client.Protocol != "http" || client.Connections != 3 || client.Rps != 7 {
		t.Fatalf("decoded client scalar fields incorrectly: %+v", client)
	}
	if len(client.EndpointUrls) != 1 || client.EndpointUrls[0] != "https://meshery.io" {
		t.Fatalf("decoded endpoint URLs incorrectly: %+v", client.EndpointUrls)
	}
	if client.ContentType != "application/json" || client.AdditionalOptions == "" {
		t.Fatalf("decoded optional client fields incorrectly: %+v", client)
	}
}

func TestDecodePerformanceTestConfigSupportsSchemaNativeJSON(t *testing.T) {
	original := &perfprofile.PerformanceTestConfig{
		SmpVersion: "v1",
		ID:         "schema-id",
		Name:       "schema config",
		Duration:   "15s",
		Clients: []perfprofile.PerformanceTestClient{
			{
				LoadGenerator: "fortio",
				Protocol:      "tcp",
				Connections:   2,
				Rps:           4,
				EndpointUrls:  []string{"tcp://127.0.0.1:8080"},
			},
		},
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	config, err := decodePerformanceTestConfig(data)
	if err != nil {
		t.Fatalf("decodePerformanceTestConfig returned error: %v", err)
	}

	if config.ID != original.ID || config.SmpVersion != original.SmpVersion || config.Clients[0].Protocol != "tcp" {
		t.Fatalf("schema-native config changed during decode: %+v", config)
	}
}

func TestTestProfilesPersisterWriteTestConfigUpdatesExistingRow(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&PerformanceTestConfig{}); err != nil {
		t.Fatalf("migrate test configs: %v", err)
	}
	persister := &TestProfilesPersister{DB: &db}

	marshalConfig := func(name string) []byte {
		payload, err := json.Marshal(&perfprofile.PerformanceTestConfig{Name: name})
		if err != nil {
			t.Fatalf("marshal test config: %v", err)
		}
		return payload
	}

	existingID := uuid.Must(uuid.NewV4())
	if err := persister.WriteTestConfig(existingID, marshalConfig("old")); err != nil {
		t.Fatalf("write initial config: %v", err)
	}
	if err := persister.WriteTestConfig(existingID, marshalConfig("new")); err != nil {
		t.Fatalf("update config: %v", err)
	}
	updated, err := persister.GetTestConfig(existingID)
	if err != nil {
		t.Fatalf("read updated config: %v", err)
	}
	if updated.Name != "new" {
		t.Fatalf("updated config name = %q, want %q", updated.Name, "new")
	}

	newID := uuid.Must(uuid.NewV4())
	if err := persister.WriteTestConfig(newID, marshalConfig("inserted")); err != nil {
		t.Fatalf("write new config: %v", err)
	}
	inserted, err := persister.GetTestConfig(newID)
	if err != nil {
		t.Fatalf("read new config: %v", err)
	}
	if inserted.Name != "inserted" {
		t.Fatalf("inserted config name = %q, want %q", inserted.Name, "inserted")
	}
}
