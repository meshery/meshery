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
