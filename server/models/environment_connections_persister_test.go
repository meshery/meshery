package models

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	connectionv1beta1 "github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta3/environment"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newEnvironmentConnectionsTestDB builds an in-memory database carrying just
// the two tables GetEnvironmentConnections queries.
func newEnvironmentConnectionsTestDB(t *testing.T) *database.Handler {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	handler := &database.Handler{DB: db}
	// The real `connections` table carries two column sets for the same
	// fields: meshkit's registry manager AutoMigrates the v1beta1 Connection
	// (`type`, `user_id`) while the server AutoMigrates the canonical v1beta3
	// Connection (`connection_type`, `owner`) - see legacyRegistrantColumnsByID
	// in seed_connections.go. GetEnvironmentConnections' SELECT list names the
	// v1beta1 spellings, so a harness that migrates only v1beta3 fails with
	// "no such column: connections.type" and tests nothing.
	// Migrated one model per call on purpose: both Connection structs map to
	// the `connections` table, and gorm collapses same-table models passed in a
	// single AutoMigrate call, so only one column set would survive.
	for _, model := range []interface{}{
		&connectionv1beta1.Connection{},
		&connections.Connection{},
		&environment.EnvironmentConnectionMapping{},
	} {
		if err := handler.AutoMigrate(model); err != nil {
			t.Fatalf("failed to auto-migrate %T: %v", model, err)
		}
	}

	return handler
}

func newTestUUID(t *testing.T) uuid.UUID {
	t.Helper()

	id, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate uuid: %v", err)
	}

	return id
}

// seedConnection persists one connection and assigns it to environmentID.
func seedConnection(t *testing.T, handler *database.Handler, environmentID core.Uuid, name string, owner core.Uuid) {
	t.Helper()

	connectionID := newTestUUID(t)
	conn := connections.Connection{
		ID:             connectionID,
		Name:           name,
		Kind:           "kubernetes",
		ConnectionType: "platform",
		Owner:          &owner,
	}
	if err := handler.Create(&conn).Error; err != nil {
		t.Fatalf("failed to persist connection %q: %v", name, err)
	}

	envID := environmentID
	mapping := environment.EnvironmentConnectionMapping{
		ID:            newTestUUID(t),
		EnvironmentID: &envID,
		ConnectionID:  &connectionID,
	}
	if err := handler.Create(&mapping).Error; err != nil {
		t.Fatalf("failed to map connection %q to environment: %v", name, err)
	}
}

func fetchConnections(t *testing.T, ep *EnvironmentPersister, environmentID core.Uuid, filter string) *connections.ConnectionPage {
	t.Helper()

	raw, err := ep.GetEnvironmentConnections(environmentID, "", "", "0", "10", filter)
	if err != nil {
		t.Fatalf("GetEnvironmentConnections(filter=%q) returned error: %v", filter, err)
	}

	page := &connections.ConnectionPage{}
	if err := json.Unmarshal(raw, page); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return page
}

// TestGetEnvironmentConnectionsOwnerFilterIsApplied is the regression for the
// substance of issue #21826: the `owner` filter was declared but unreachable,
// because the JSON grammar this endpoint accepts could never satisfy
// utils.ApplyFilters' "key value" grammar. Before the fix this returned both
// connections.
func TestGetEnvironmentConnectionsOwnerFilterIsApplied(t *testing.T) {
	handler := newEnvironmentConnectionsTestDB(t)
	ep := &EnvironmentPersister{DB: handler}

	environmentID := newTestUUID(t)
	wantedOwner := newTestUUID(t)
	otherOwner := newTestUUID(t)

	seedConnection(t, handler, environmentID, "owned-by-wanted", wantedOwner)
	seedConnection(t, handler, environmentID, "owned-by-other", otherOwner)

	all := fetchConnections(t, ep, environmentID, "")
	if all.TotalCount != 2 {
		t.Fatalf("unfiltered TotalCount = %d, want 2", all.TotalCount)
	}

	filtered := fetchConnections(t, ep, environmentID, `{"owner":"`+wantedOwner.String()+`"}`)
	if filtered.TotalCount != 1 {
		t.Fatalf("owner-filtered TotalCount = %d, want 1", filtered.TotalCount)
	}
	if len(filtered.Connections) != 1 {
		t.Fatalf("owner-filtered returned %d connections, want 1", len(filtered.Connections))
	}
	if got := filtered.Connections[0].Name; got != "owned-by-wanted" {
		t.Errorf("owner-filtered connection name = %q, want %q", got, "owned-by-wanted")
	}
}

// TestGetEnvironmentConnectionsAssignedFilter pins the value Meshery UI sends.
func TestGetEnvironmentConnectionsAssignedFilter(t *testing.T) {
	handler := newEnvironmentConnectionsTestDB(t)
	ep := &EnvironmentPersister{DB: handler}

	environmentID := newTestUUID(t)
	seedConnection(t, handler, environmentID, "assigned-one", newTestUUID(t))

	assigned := fetchConnections(t, ep, environmentID, `{"assigned":true}`)
	if assigned.TotalCount != 1 {
		t.Errorf(`filter={"assigned":true} TotalCount = %d, want 1`, assigned.TotalCount)
	}

	// The unassigned complement: the one seeded connection is assigned, so
	// nothing is left over.
	unassigned := fetchConnections(t, ep, environmentID, `{"assigned":false}`)
	if unassigned.TotalCount != 0 {
		t.Errorf(`filter={"assigned":false} TotalCount = %d, want 0`, unassigned.TotalCount)
	}
}

// TestGetEnvironmentConnectionsRejectsMalformedFilter is the panic regression.
// `{"assigned":"false"}` previously reached an unchecked type assertion and
// took the handler goroutine down; it must now come back as an error.
func TestGetEnvironmentConnectionsRejectsMalformedFilter(t *testing.T) {
	handler := newEnvironmentConnectionsTestDB(t)
	ep := &EnvironmentPersister{DB: handler}

	environmentID := newTestUUID(t)

	for _, filter := range []string{`{"assigned":"false"}`, `{"assigned":0}`, `{"assigned":null}`, "owner abc"} {
		t.Run(filter, func(t *testing.T) {
			// Must return an error rather than panic.
			if _, err := ep.GetEnvironmentConnections(environmentID, "", "", "0", "10", filter); err == nil {
				t.Fatalf("GetEnvironmentConnections(filter=%q) returned no error, want one", filter)
			}
		})
	}
}
