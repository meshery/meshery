package models

import (
	"testing"
	"time"

	"github.com/meshery/meshery/server/models/connections"
)

// TestGetConnectionsFiltersByConnectionType checks that the `type` query
// filter matches connections saved through SaveConnection. Those rows only
// carry the type in `connection_type`; the legacy `type` column that meshkit's
// registry migration adds stays empty for them.
func TestGetConnectionsFiltersByConnectionType(t *testing.T) {
	db := newMigratedDB(t)
	// On a real server meshkit's registry manager also migrates the v1beta1
	// Connection, which adds a legacy `type` column next to `connection_type`.
	if err := db.Exec("ALTER TABLE connections ADD COLUMN type text").Error; err != nil {
		t.Fatalf("add legacy type column: %v", err)
	}
	cp := &ConnectionPersister{DB: db}

	for _, c := range []struct{ name, kind, connType, subType string }{
		{"cluster", "kubernetes", "platform", "orchestration"},
		{"prom", "prometheus", "telemetry", "metrics"},
		{"grafana", "grafana", "telemetry", "monitoring"},
	} {
		if _, err := cp.SaveConnection(&connections.Connection{
			ID:             mustUUID(t),
			Name:           c.name,
			Kind:           c.kind,
			ConnectionType: c.connType,
			SubType:        c.subType,
			Status:         connections.CONNECTED,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}); err != nil {
			t.Fatalf("save %s: %v", c.name, err)
		}
	}

	page, err := cp.GetConnections("", "", 0, 10, "", nil, nil, []string{"telemetry"}, "")
	if err != nil {
		t.Fatalf("GetConnections: %v", err)
	}
	if page.TotalCount != 2 || len(page.Connections) != 2 {
		t.Fatalf("type=telemetry returned totalCount=%d, %d rows; want 2", page.TotalCount, len(page.Connections))
	}
	for _, c := range page.Connections {
		if c.ConnectionType != "telemetry" {
			t.Errorf("connection %q has type %q, want telemetry", c.Name, c.ConnectionType)
		}
	}
}

func TestMapConnectionOrderAlias(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"type asc maps to connection_type asc", "type asc", "connection_type asc"},
		{"type desc maps to connection_type desc", "type desc", "connection_type desc"},
		{"Type capitalized maps to connection_type", "Type asc", "connection_type asc"},
		{"ConnectionID asc maps to id asc", "ConnectionID asc", "id asc"},
		{"connectionId desc maps to id desc", "connectionId desc", "id desc"},
		{"connection_id asc maps to id asc", "connection_id asc", "id asc"},
		{"kind asc remains unchanged", "kind asc", "kind asc"},
		{"kind desc remains unchanged", "kind desc", "kind desc"},
		{"subType asc remains unchanged for snake_case folding", "subType asc", "subType asc"},
		{"status desc remains unchanged", "status desc", "status desc"},
		{"id asc remains unchanged", "id asc", "id asc"},
		{"name asc remains unchanged", "name asc", "name asc"},
		{"empty input returns empty", "", ""},
		{"single word returns input", "kind", "kind"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapConnectionOrderAlias(tt.input)
			if got != tt.expected {
				t.Errorf("mapConnectionOrderAlias(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGetConnectionsSorting(t *testing.T) {
	db := newMigratedDB(t)
	cp := &ConnectionPersister{DB: db}

	c1 := &connections.Connection{
		ID:             mustUUID(t),
		Name:           "alpha-conn",
		Kind:           "kubernetes",
		ConnectionType: "platform",
		SubType:        "orchestration",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now().Add(-3 * time.Hour),
		UpdatedAt:      time.Now().Add(-3 * time.Hour),
	}
	c2 := &connections.Connection{
		ID:             mustUUID(t),
		Name:           "beta-conn",
		Kind:           "grafana",
		ConnectionType: "telemetry",
		SubType:        "monitoring",
		Status:         connections.DISCOVERED,
		CreatedAt:      time.Now().Add(-2 * time.Hour),
		UpdatedAt:      time.Now().Add(-2 * time.Hour),
	}
	c3 := &connections.Connection{
		ID:             mustUUID(t),
		Name:           "gamma-conn",
		Kind:           "prometheus",
		ConnectionType: "telemetry",
		SubType:        "metrics",
		Status:         connections.REGISTERED,
		CreatedAt:      time.Now().Add(-1 * time.Hour),
		UpdatedAt:      time.Now().Add(-1 * time.Hour),
	}

	for _, c := range []*connections.Connection{c1, c2, c3} {
		if _, err := cp.SaveConnection(c); err != nil {
			t.Fatalf("save connection %s: %v", c.Name, err)
		}
	}

	// Explicitly space updated_at timestamps so updated_at desc ordering is deterministic
	now := time.Now()
	db.Model(&connections.Connection{}).Where("name = ?", "alpha-conn").UpdateColumn("updated_at", now.Add(-30*time.Minute))
	db.Model(&connections.Connection{}).Where("name = ?", "beta-conn").UpdateColumn("updated_at", now.Add(-20*time.Minute))
	db.Model(&connections.Connection{}).Where("name = ?", "gamma-conn").UpdateColumn("updated_at", now.Add(-10*time.Minute))

	// 1. Sort by kind asc: grafana (c2), kubernetes (c1), prometheus (c3)
	page, err := cp.GetConnections("", "kind asc", 0, 10, "", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("GetConnections kind asc: %v", err)
	}
	if len(page.Connections) != 3 {
		t.Fatalf("expected 3 connections, got %d", len(page.Connections))
	}
	if page.Connections[0].Kind != "grafana" || page.Connections[1].Kind != "kubernetes" || page.Connections[2].Kind != "prometheus" {
		t.Errorf("unexpected order for kind asc: %v, %v, %v",
			page.Connections[0].Kind, page.Connections[1].Kind, page.Connections[2].Kind)
	}

	// 2. Sort by kind desc: prometheus (c3), kubernetes (c1), grafana (c2)
	page, err = cp.GetConnections("", "kind desc", 0, 10, "", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("GetConnections kind desc: %v", err)
	}
	if page.Connections[0].Kind != "prometheus" || page.Connections[1].Kind != "kubernetes" || page.Connections[2].Kind != "grafana" {
		t.Errorf("unexpected order for kind desc: %v, %v, %v",
			page.Connections[0].Kind, page.Connections[1].Kind, page.Connections[2].Kind)
	}

	// 3. Sort by wire alias "type asc": platform (c1), telemetry (c2/c3)
	page, err = cp.GetConnections("", "type asc", 0, 10, "", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("GetConnections type asc: %v", err)
	}
	if page.Connections[0].ConnectionType != "platform" {
		t.Errorf("expected first element to have type 'platform', got %q", page.Connections[0].ConnectionType)
	}

	// 4. Sort by wire alias "type desc": telemetry (c2/c3), platform (c1)
	page, err = cp.GetConnections("", "type desc", 0, 10, "", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("GetConnections type desc: %v", err)
	}
	if page.Connections[2].ConnectionType != "platform" {
		t.Errorf("expected last element to have type 'platform', got %q", page.Connections[2].ConnectionType)
	}

	// 5. Sort by camelCase "subType asc": metrics (c3), monitoring (c2), orchestration (c1)
	page, err = cp.GetConnections("", "subType asc", 0, 10, "", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("GetConnections subType asc: %v", err)
	}
	if page.Connections[0].SubType != "metrics" || page.Connections[1].SubType != "monitoring" || page.Connections[2].SubType != "orchestration" {
		t.Errorf("unexpected order for subType asc: %v, %v, %v",
			page.Connections[0].SubType, page.Connections[1].SubType, page.Connections[2].SubType)
	}

	// 6. Sort by name asc: alpha-conn (c1), beta-conn (c2), gamma-conn (c3)
	page, err = cp.GetConnections("", "name asc", 0, 10, "", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("GetConnections name asc: %v", err)
	}
	if page.Connections[0].Name != "alpha-conn" || page.Connections[2].Name != "gamma-conn" {
		t.Errorf("unexpected order for name asc: first=%s, last=%s", page.Connections[0].Name, page.Connections[2].Name)
	}

	// 7. Sort by ConnectionID / id
	page, err = cp.GetConnections("", "ConnectionID asc", 0, 10, "", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("GetConnections ConnectionID asc: %v", err)
	}
	if len(page.Connections) != 3 {
		t.Fatalf("expected 3 connections, got %d", len(page.Connections))
	}

	// 8. Invalid column falls back safely to defaultOrderUpdatedAtDesc (gamma, beta, alpha)
	page, err = cp.GetConnections("", "invalidColumn asc", 0, 10, "", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("GetConnections invalidColumn asc: %v", err)
	}
	if page.Connections[0].Name != "gamma-conn" || page.Connections[2].Name != "alpha-conn" {
		t.Errorf("expected fallback to updated_at desc, got first=%s, last=%s", page.Connections[0].Name, page.Connections[2].Name)
	}
}
