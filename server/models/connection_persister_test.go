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

	fetched := connections.Connection{}
	if err := dbHandler.First(&fetched, "id = ?", saved.ID).Error; err != nil {
		t.Fatalf("expected connection to exist in DB after save, got: %v", err)
	}
	if fetched.Name != "test-connection" {
		t.Fatalf("unexpected name: got %q, want %q", fetched.Name, "test-connection")
	}
}

func TestSaveConnection_DuplicateIDReturnsError(t *testing.T) {
	dbHandler := newConnectionTestDB(t)
	cp := &ConnectionPersister{DB: dbHandler}

	conn := &connections.Connection{
		Name:   "original-connection",
		Kind:   "kubernetes",
		Status: connections.CONNECTED,
	}
	saved, err := cp.SaveConnection(conn)
	if err != nil {
		t.Fatalf("first SaveConnection failed unexpectedly: %v", err)
	}

	duplicate := &connections.Connection{
		ID:     saved.ID,
		Name:   "duplicate-connection",
		Kind:   "kubernetes",
		Status: connections.CONNECTED,
	}
	_, dupErr := cp.SaveConnection(duplicate)
	if dupErr == nil {
		t.Fatal("expected SaveConnection to return an error for a duplicate ID, but got nil")
	}
	if dupErr.Error() == "" {
		t.Fatal("expected a non-empty error message for the duplicate-ID case")
	}
}

func TestSaveConnection_AutoAssignsUUID(t *testing.T) {
	dbHandler := newConnectionTestDB(t)
	cp := &ConnectionPersister{DB: dbHandler}

	conn := &connections.Connection{
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

	var count int64
	dbHandler.Model(&connections.Connection{}).Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 connections in DB, got %d", count)
	}
}

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
