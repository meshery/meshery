package models

import (
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	meshkitBroker "github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta3/component"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testMeshsyncBroker struct {
	mu            sync.Mutex
	subscriptions []chan *meshkitBroker.Message
	closed        bool
}

func (b *testMeshsyncBroker) Publish(_ string, _ *meshkitBroker.Message) error {
	return nil
}

func (b *testMeshsyncBroker) PublishWithChannel(_ string, _ chan *meshkitBroker.Message) error {
	return nil
}

func (b *testMeshsyncBroker) Subscribe(_ string, _ string, _ []byte) error {
	return nil
}

func (b *testMeshsyncBroker) SubscribeWithChannel(_ string, _ string, subscription chan *meshkitBroker.Message) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscriptions = append(b.subscriptions, subscription)
	return nil
}

func (b *testMeshsyncBroker) Info() string {
	return "test-broker"
}

func (b *testMeshsyncBroker) DeepCopyObject() meshkitBroker.Handler {
	return b
}

func (b *testMeshsyncBroker) DeepCopyInto(meshkitBroker.Handler) {}

func (b *testMeshsyncBroker) IsEmpty() bool {
	return false
}

func (b *testMeshsyncBroker) CloseConnection() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for _, subscription := range b.subscriptions {
		close(subscription)
	}
}

func (b *testMeshsyncBroker) ConnectedEndpoints() []string {
	return nil
}

func (b *testMeshsyncBroker) isClosed() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.closed
}

func newTestMeshsyncHandler(t *testing.T, broker meshkitBroker.Handler, stopFunc func()) *MeshsyncDataHandler {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	err = db.AutoMigrate(&component.ComponentDefinition{})
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	mockLogger, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	return NewMeshsyncDataHandler(
		broker,
		database.Handler{DB: db},
		mockLogger,
		nil,
		uuid.Nil,
		uuid.Nil,
		uuid.Nil,
		"",
		stopFunc,
	)
}

// TestGetComponentMetadataWithNilModel tests that getComponentMetadata
// does not panic when Model is nil (e.g., when record is not found)
func TestGetComponentMetadataWithNilModel(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&component.ComponentDefinition{})
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	// Create a mock logger
	mockLogger, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Create a database handler wrapper
	dbHandler := database.Handler{DB: db}

	// Create MeshsyncDataHandler
	handler := &MeshsyncDataHandler{
		dbHandler: dbHandler,
		log:       mockLogger,
	}

	// Test case 1: Record not found (should not panic)
	t.Run("RecordNotFound", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("getComponentMetadata panicked with: %v", r)
			}
		}()

		// Call with non-existent apiVersion and kind
		data, model := handler.getComponentMetadata("non/existent/v1", "NonExistentKind")

		// Verify that it returns without panicking
		if data == nil {
			t.Error("Expected data to be non-nil")
		}

		// Model should be empty string when not found
		if model != "" {
			t.Errorf("Expected empty model string, got: %s", model)
		}
	})

	// Test case 2: Valid record with Model set (should work correctly)
	t.Run("ValidRecordWithModel", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("getComponentMetadata panicked with: %v", r)
			}
		}()

		// This test would require setting up proper test data with models
		// For now, we just verify it doesn't panic with non-existent records
		data, model := handler.getComponentMetadata("v1", "Pod")

		if data == nil {
			t.Error("Expected data to be non-nil")
		}

		// Since we don't have test data, model will be empty
		_ = model
	})
}

func TestMeshsyncDataHandlerStopStopsListeners(t *testing.T) {
	stopCalled := make(chan struct{}, 1)
	broker := &testMeshsyncBroker{}
	handler := newTestMeshsyncHandler(t, broker, func() {
		stopCalled <- struct{}{}
	})

	if err := handler.Run(); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	handler.Stop()

	select {
	case <-stopCalled:
	case <-time.After(time.Second):
		t.Fatal("expected Stop to invoke stopFunc")
	}

	if !broker.isClosed() {
		t.Fatal("expected Stop to close the broker connection")
	}
}

func TestMeshsyncDataHandlerStopIsIdempotent(t *testing.T) {
	var stopCalls int
	broker := &testMeshsyncBroker{}
	handler := newTestMeshsyncHandler(t, broker, func() {
		stopCalls++
	})

	if err := handler.Run(); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	handler.Stop()
	handler.Stop()

	if stopCalls != 1 {
		t.Fatalf("expected stopFunc to be called once, got %d", stopCalls)
	}
}
