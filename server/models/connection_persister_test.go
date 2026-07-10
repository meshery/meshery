package models

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newConnectionTestDB creates an in-memory SQLite database for connection tests.
func newConnectionTestDB(t *testing.T) *database.Handler {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	if err := db.AutoMigrate(&connections.Connection{}); err != nil {
		t.Fatalf("failed to auto-migrate connections schema: %v", err)
	}

	return &database.Handler{DB: db}
}

// TestSaveConnection_NewConnection verifies that a connection with a new ID is
// persisted successfully and the returned object carries the same ID.
func TestSaveConnection_NewConnection(t *testing.T) {
	dbHandler := newConnectionTestDB(t)
	cp := &ConnectionPersister{DB: dbHandler}

	conn := &connections.Connection{
		Name:   "test-connection",
		Kind:   "kubernetes",
		Status: connections.CONNECTED,
	}

	saved, err := cp.SaveConnection(conn)
	if err != nil {
		t.Fatalf("expected SaveConnection to succeed for a new connection, got error: %v", err)
	}
	if saved.ID == uuid.Nil {
		t.Fatal("expected SaveConnection to assign a UUID to the connection")
	}

	// Confirm it was actually written to the DB.
	fetched := connections.Connection{}
	if err := dbHandler.First(&fetched, "id = ?", saved.ID).Error; err != nil {
		t.Fatalf("expected connection to exist in DB after save, got: %v", err)
	}
	if fetched.Name != "test-connection" {
		t.Fatalf("unexpected name: got %q, want %q", fetched.Name, "test-connection")
	}
}

// TestSaveConnection_DuplicateIDReturnsError verifies the fix for the inverted
// duplicate-check bug (issue #19151): calling SaveConnection with an ID that
// already exists must return an error instead of silently no-oping.
func TestSaveConnection_DuplicateIDReturnsError(t *testing.T) {
	dbHandler := newConnectionTestDB(t)
	cp := &ConnectionPersister{DB: dbHandler}

	// First save — must succeed.
	conn := &connections.Connection{
		Name:   "original-connection",
		Kind:   "kubernetes",
		Status: connections.CONNECTED,
	}
	saved, err := cp.SaveConnection(conn)
	if err != nil {
		t.Fatalf("first SaveConnection failed unexpectedly: %v", err)
	}

	// Second save with the same ID — must fail with a descriptive error.
	duplicate := &connections.Connection{
		ID:     saved.ID,
		Name:   "duplicate-connection",
		Kind:   "kubernetes",
		Status: connections.CONNECTED,
	}
	_, dupErr := cp.SaveConnection(duplicate)
	if dupErr == nil {
		t.Fatal("expected SaveConnection to return an error for a duplicate ID, but got nil (silent no-op bug is back)")
	}

	// Confirm the error message is descriptive.
	if dupErr.Error() == "" {
		t.Fatal("expected a non-empty error message for the duplicate-ID case")
	}
}

// TestSaveConnection_AutoAssignsUUID verifies that a connection with uuid.Nil
// as its ID gets a new UUID generated before saving.
func TestSaveConnection_AutoAssignsUUID(t *testing.T) {
	dbHandler := newConnectionTestDB(t)
	cp := &ConnectionPersister{DB: dbHandler}

	conn := &connections.Connection{
		// ID intentionally left as uuid.Nil
		Name:   "no-id-connection",
		Kind:   "kubernetes",
		Status: connections.REGISTERED,
	}

	saved, err := cp.SaveConnection(conn)
	if err != nil {
		t.Fatalf("SaveConnection failed: %v", err)
	}
	if saved.ID == uuid.Nil {
		t.Fatal("expected a non-nil UUID to be assigned when connection.ID is uuid.Nil")
	}
}

// TestSaveConnection_MultipleDistinctConnectionsPersist verifies that two
// connections with different IDs are both persisted independently without
// interfering with each other.
func TestSaveConnection_MultipleDistinctConnectionsPersist(t *testing.T) {
	dbHandler := newConnectionTestDB(t)
	cp := &ConnectionPersister{DB: dbHandler}

	for i, name := range []string{"conn-alpha", "conn-beta"} {
		conn := &connections.Connection{
			Name:   name,
			Kind:   "kubernetes",
			Status: connections.DISCOVERED,
		}
		saved, err := cp.SaveConnection(conn)
		if err != nil {
			t.Fatalf("SaveConnection[%d] (%q) failed: %v", i, name, err)
		}
		if saved.Name != name {
			t.Fatalf("unexpected name at index %d: got %q, want %q", i, saved.Name, name)
		}
	}

	// Confirm both rows are in the DB.
	var count int64
	dbHandler.Model(&connections.Connection{}).Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 connections in DB, got %d", count)
	}
}

// TestSaveConnection_ErrRecordNotFoundIsNotPropagated ensures that the normal
// "not found" branch (new record) does not surface gorm.ErrRecordNotFound to
// callers — only the Save result is returned.
func TestSaveConnection_ErrRecordNotFoundIsNotPropagated(t *testing.T) {
	dbHandler := newConnectionTestDB(t)
	cp := &ConnectionPersister{DB: dbHandler}

	conn := &connections.Connection{
		Name:   "brand-new",
		Kind:   "prometheus",
		Status: connections.REGISTERED,
	}

	_, err := cp.SaveConnection(conn)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatal("SaveConnection must not expose gorm.ErrRecordNotFound for a normal insert")
	}
}
