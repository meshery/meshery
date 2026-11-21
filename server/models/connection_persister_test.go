package models

import (
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"

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

// TestGetConnectionsSortsBySortableColumns checks every column the connections
// table lets the user sort by. The table sends camelCase column names
// (subType, ConnectionID) and the public `type` alias; each has to reach the
// matching database column instead of falling back to the default order.
func TestGetConnectionsSortsBySortableColumns(t *testing.T) {
	db := newMigratedDB(t)
	cp := &ConnectionPersister{DB: db}

	ids := []string{
		"33333333-3333-4333-8333-333333333333",
		"11111111-1111-4111-8111-111111111111",
		"22222222-2222-4222-8222-222222222222",
	}
	// Rows are saved so that updated_at (the default order) differs from the
	// order of every column under test.
	base := time.Now()
	for i, c := range []struct {
		name, kind, connType, subType string
		status                        connections.ConnectionStatus
	}{
		{"b-conn", "prometheus", "observability", "monitoring", connections.DISCOVERED},
		{"c-conn", "grafana", "platform", "metrics", connections.CONNECTED},
		{"a-conn", "kubernetes", "telemetry", "orchestration", connections.IGNORED},
	} {
		if _, err := cp.SaveConnection(&connections.Connection{
			ID:             core.Uuid(uuid.FromStringOrNil(ids[i])),
			Name:           c.name,
			Kind:           c.kind,
			ConnectionType: c.connType,
			SubType:        c.subType,
			Status:         c.status,
			CreatedAt:      base.Add(time.Duration(i) * time.Minute),
			UpdatedAt:      base.Add(time.Duration(i) * time.Minute),
		}); err != nil {
			t.Fatalf("save %s: %v", c.name, err)
		}
	}

	tests := []struct {
		order string
		want  []string
	}{
		{"kind asc", []string{"c-conn", "a-conn", "b-conn"}},
		{"type asc", []string{"b-conn", "c-conn", "a-conn"}},
		{"subType asc", []string{"c-conn", "b-conn", "a-conn"}},
		{"status asc", []string{"c-conn", "b-conn", "a-conn"}},
		{"ConnectionID asc", []string{"c-conn", "a-conn", "b-conn"}},
		{"kind desc", []string{"b-conn", "a-conn", "c-conn"}},
	}
	for _, tc := range tests {
		t.Run(tc.order, func(t *testing.T) {
			page, err := cp.GetConnections("", tc.order, 0, 10, "", nil, nil, nil, "")
			if err != nil {
				t.Fatalf("GetConnections: %v", err)
			}
			got := make([]string, len(page.Connections))
			for i, c := range page.Connections {
				got[i] = c.Name
			}
			if strings.Join(got, ",") != strings.Join(tc.want, ",") {
				t.Errorf("order %q returned %v, want %v", tc.order, got, tc.want)
			}
		})
	}
}
