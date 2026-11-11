package models

import (
	"testing"

	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	"github.com/meshery/schemas/models/v1beta1/environment"
)

// newConnectionPersisterFixture returns a persister over an in-memory database
// carrying every table GetConnections reads.
func newConnectionPersisterFixture(t *testing.T) (*ConnectionPersister, *database.Handler) {
	t.Helper()

	db, _ := seedTestRegistry(t)
	if err := db.AutoMigrate(environment.Environment{}, environment.EnvironmentConnectionMapping{}); err != nil {
		t.Fatalf("migrate environment tables: %v", err)
	}
	return &ConnectionPersister{DB: db}, db
}

// TestGetConnectionsExcludesUnseededRegistrants covers meshery/meshery#21074:
// registering Kubernetes components writes a `kubernetes` registrant row into the
// `connections` table, and that row must not be listed beside the user's real
// cluster connection - in the page, the total count or the status summary.
func TestGetConnectionsExcludesUnseededRegistrants(t *testing.T) {
	cp, db := newConnectionPersisterFixture(t)

	cluster, err := cp.SaveConnection(&connections.Connection{
		Name:           "minikube",
		Kind:           "kubernetes",
		ConnectionType: "platform",
		SubType:        "orchestration",
		Status:         connections.CONNECTED,
	})
	if err != nil {
		t.Fatalf("save cluster connection: %v", err)
	}

	registerRegistrant(t, db, "kubernetes", entity.Model)

	// A registrant that seeding has promoted carries a canonical connection type
	// and is a system-owned connection the page must keep showing.
	seededID := registerRegistrant(t, db, "artifacthub", entity.Model)
	if err := db.Model(&connections.Connection{}).Where("id = ?", seededID).Updates(map[string]interface{}{
		"name":            "Artifact Hub",
		"connection_type": "source",
		"sub_type":        "registry",
		"status":          string(connections.CONNECTED),
	}).Error; err != nil {
		t.Fatalf("promote artifacthub registrant: %v", err)
	}

	page, err := cp.GetConnections("", "", 0, 25, "", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("GetConnections: %v", err)
	}

	got := map[string]bool{}
	for _, conn := range page.Connections {
		got[conn.Name] = true
	}
	if !got[cluster.Name] || !got["Artifact Hub"] || len(page.Connections) != 2 {
		t.Errorf("listed connections = %v, want exactly the cluster connection and the seeded Artifact Hub", got)
	}
	if page.TotalCount != 2 {
		t.Errorf("TotalCount = %d, want 2", page.TotalCount)
	}

	total := 0
	for _, n := range *page.StatusSummary {
		total += n
	}
	if total != 2 {
		t.Errorf("status summary counts %d connection(s), want 2: %v", total, *page.StatusSummary)
	}
}
