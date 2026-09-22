package models

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newConnectionTestDB(t *testing.T) *database.Handler {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	return &database.Handler{DB: db}
}

// TestSaveConnectionPreservesDisconnectedStatus pins the invariant shouldDriveDiscovery depends on:
// if re-import ever resurrects a DISCONNECTED connection as DISCOVERED, the gate silently stops working.
func TestSaveConnectionPreservesDisconnectedStatus(t *testing.T) {
	dbHandler := newConnectionTestDB(t)
	if err := dbHandler.AutoMigrate(&connections.Connection{}); err != nil {
		t.Fatalf("failed to auto-migrate connections: %v", err)
	}

	connectionID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate connection id: %v", err)
	}

	now := time.Now().UTC().Round(time.Second)
	initialConn := connections.Connection{
		ID:        core.Uuid(connectionID),
		CreatedAt: now,
		UpdatedAt: now,
		Kind:      "kubernetes",
		Status:    connections.DISCONNECTED,
	}

	if err := dbHandler.Create(&initialConn).Error; err != nil {
		t.Fatalf("failed to persist initial connection: %v", err)
	}

	persister := &ConnectionPersister{DB: dbHandler}

	payloadConn := connections.Connection{
		ID:     core.Uuid(connectionID),
		Kind:   "kubernetes",
		Status: connections.DISCOVERED,
	}

	returnedConn, err := persister.SaveConnection(&payloadConn)
	if err != nil {
		t.Fatalf("expected SaveConnection to succeed, got %v", err)
	}

	if returnedConn.Status != connections.DISCONNECTED {
		t.Fatalf("expected returned connection status to remain DISCONNECTED, got %s", returnedConn.Status)
	}

	stored := connections.Connection{}
	if err := dbHandler.First(&stored, "id = ?", connectionID).Error; err != nil {
		t.Fatalf("failed to read persisted connection: %v", err)
	}

	if stored.Status != connections.DISCONNECTED {
		t.Fatalf("expected stored connection status to remain DISCONNECTED, got %s", stored.Status)
	}
}
