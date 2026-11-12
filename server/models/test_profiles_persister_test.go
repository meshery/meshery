package models

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
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

func TestGetTestConfigNotFound(t *testing.T) {
	db := newMigratedDB(t)
	tpp := &TestProfilesPersister{DB: db}

	_, err := tpp.GetTestConfig(mustUUID(t))
	if errors.GetCode(err) != ErrResultNotFoundCode {
		t.Fatalf("error code = %q, want %q (err: %v)", errors.GetCode(err), ErrResultNotFoundCode, err)
	}
}

func TestDeleteTestConfigNotFoundAndExisting(t *testing.T) {
	db := newMigratedDB(t)
	tpp := &TestProfilesPersister{DB: db}

	err := tpp.DeleteTestConfig(mustUUID(t))
	if errors.GetCode(err) != ErrResultNotFoundCode {
		t.Fatalf("delete missing error code = %q, want %q (err: %v)", errors.GetCode(err), ErrResultNotFoundCode, err)
	}

	id := mustUUID(t)
	payload := []byte(`{"id":"` + id.String() + `","name":"sample-config","duration":"10s"}`)
	if err := tpp.WriteTestConfig(id, payload); err != nil {
		t.Fatalf("WriteTestConfig: %v", err)
	}

	if err := tpp.DeleteTestConfig(id); err != nil {
		t.Fatalf("DeleteTestConfig: %v", err)
	}

	_, err = tpp.GetTestConfig(id)
	if errors.GetCode(err) != ErrResultNotFoundCode {
		t.Fatalf("get after delete error code = %q, want %q (err: %v)", errors.GetCode(err), ErrResultNotFoundCode, err)
	}
}

func TestGetTestConfigsPropagatesDBError(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	tpp := &TestProfilesPersister{DB: &db}

	if body, err := tpp.GetTestConfigs(0, 10); err == nil {
		t.Fatalf("expected DB error on unmigrated table, got body %s", body)
	}
}

func TestDefaultLocalProviderSMPTestConfigDefaultsAndValidation(t *testing.T) {
	db := newMigratedDB(t)
	provider := &DefaultLocalProvider{
		TestProfilesPersister: &TestProfilesPersister{DB: db},
	}

	if _, err := provider.SMPTestConfigGet(nil, "not-a-uuid"); errors.GetCode(err) != ErrInvalidUUIDValueCode {
		t.Fatalf("SMPTestConfigGet malformed error code = %q, want %q", errors.GetCode(err), ErrInvalidUUIDValueCode)
	}
	if _, err := provider.SMPTestConfigGet(nil, core.Uuid(uuid.Nil).String()); errors.GetCode(err) != ErrInvalidUUIDValueCode {
		t.Fatalf("SMPTestConfigGet nil UUID error code = %q, want %q", errors.GetCode(err), ErrInvalidUUIDValueCode)
	}
	if err := provider.SMPTestConfigDelete(nil, "not-a-uuid"); errors.GetCode(err) != ErrInvalidUUIDValueCode {
		t.Fatalf("SMPTestConfigDelete malformed error code = %q, want %q", errors.GetCode(err), ErrInvalidUUIDValueCode)
	}

	cfg := &perfprofile.PerformanceTestConfig{
		Name:     "default-page-test",
		Duration: "10s",
	}
	if _, err := provider.SMPTestConfigStore(nil, cfg); err != nil {
		t.Fatalf("SMPTestConfigStore: %v", err)
	}

	raw, err := provider.SMPTestConfigFetch(nil, "", "", "", "")
	if err != nil {
		t.Fatalf("SMPTestConfigFetch with empty pagination: %v", err)
	}
	var page UserTestProfiles
	if err := json.Unmarshal(raw, &page); err != nil {
		t.Fatalf("unmarshal UserTestProfiles: %v", err)
	}
	if page.TotalCount != 1 || len(page.TestConfigs) != 1 {
		t.Fatalf("got totalCount=%d configs=%d, want 1", page.TotalCount, len(page.TestConfigs))
	}
}
