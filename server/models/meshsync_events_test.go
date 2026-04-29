package models

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	meshkitBroker "github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/database"
	meshkitErrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/assert"
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

type fakeMeshsyncLogger struct {
	logger.Handler
	mu     sync.Mutex
	logged []error
}

func (f *fakeMeshsyncLogger) Error(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logged = append(f.logged, err)
}

func (f *fakeMeshsyncLogger) getLogged() []error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]error(nil), f.logged...)
}

func (f *fakeMeshsyncLogger) Info(description ...interface{}) {}
func (f *fakeMeshsyncLogger) Infof(format string, args ...interface{}) {}
func (f *fakeMeshsyncLogger) Debug(description ...interface{}) {}
func (f *fakeMeshsyncLogger) Warn(err error) {}

func TestMeshsyncEventErrorWrappedWithMeshKit(t *testing.T) {
	broker := &testMeshsyncBroker{}
	logger := &fakeMeshsyncLogger{}
	handler := &MeshsyncDataHandler{
		broker: broker,
		log:    logger,
		stopCh: make(chan struct{}),
	}
	handler.listenerWg = &sync.WaitGroup{}
	handler.listenerWg.Add(1)

	go handler.subscribeToMeshsyncEvents()

	// Give the goroutine a moment to subscribe
	time.Sleep(100 * time.Millisecond)

	broker.mu.Lock()
	if len(broker.subscriptions) == 0 {
		broker.mu.Unlock()
		t.Fatalf("expected subscriber to register")
	}
	subChan := broker.subscriptions[0]
	broker.mu.Unlock()

	// Inject error event
	subChan <- &meshkitBroker.Message{
		EventType: meshkitBroker.ErrorEvent,
		Object:    fmt.Errorf("simulated meshsync error"),
	}

	// Wait for processing
	time.Sleep(100 * time.Millisecond)
	handler.StopFunc = func() {}
	close(handler.stopCh)

	logged := logger.getLogged()
	if len(logged) == 0 {
		t.Fatalf("expected error to be logged by meshsync handler")
	}

	// Check if the error is wrapped with the correct MeshKit code
	if mkErr, ok := logged[0].(*meshkitErrors.Error); ok {
		assert.Equal(t, ErrMeshsyncEventCode, mkErr.Code)
	} else {
		t.Fatalf("expected logged error to be a *meshkitErrors.Error, got %T", logged[0])
	}
}

func TestMeshsyncStoreUpdatesErrorWrappedWithMeshKit(t *testing.T) {
	broker := &testMeshsyncBroker{}
	logger := &fakeMeshsyncLogger{}
	handler := &MeshsyncDataHandler{
		broker: broker,
		log:    logger,
		stopCh: make(chan struct{}),
	}
	handler.listenerWg = &sync.WaitGroup{}
	handler.listenerWg.Add(1)

	statusChan := make(chan bool)
	go handler.subsribeToStoreUpdates(statusChan)

	// Wait for subscription to establish
	<-statusChan

	broker.mu.Lock()
	if len(broker.subscriptions) == 0 {
		broker.mu.Unlock()
		t.Fatalf("expected subscriber to register")
	}
	subChan := broker.subscriptions[0]
	broker.mu.Unlock()

	// Inject error event
	subChan <- &meshkitBroker.Message{
		EventType: meshkitBroker.ErrorEvent,
		Object:    fmt.Errorf("simulated store update error"),
	}

	// Wait for processing
	time.Sleep(100 * time.Millisecond)
	close(handler.stopCh)

	logged := logger.getLogged()
	if len(logged) == 0 {
		t.Fatalf("expected error to be logged by meshsync handler")
	}

	// Check if the error is wrapped with the correct MeshKit code
	if mkErr, ok := logged[0].(*meshkitErrors.Error); ok {
		assert.Equal(t, ErrMeshsyncStoreUpdatesCode, mkErr.Code)
	} else {
		t.Fatalf("expected logged error to be a *meshkitErrors.Error, got %T", logged[0])
	}
}

func TestMeshsyncStoreUpdatesNonErrorObjectWrapped(t *testing.T) {
	broker := &testMeshsyncBroker{}
	logger := &fakeMeshsyncLogger{}
	handler := &MeshsyncDataHandler{
		broker: broker,
		log:    logger,
		stopCh: make(chan struct{}),
	}
	handler.listenerWg = &sync.WaitGroup{}
	handler.listenerWg.Add(1)

	statusChan := make(chan bool)
	go handler.subsribeToStoreUpdates(statusChan)

	// Wait for subscription to establish
	<-statusChan

	broker.mu.Lock()
	if len(broker.subscriptions) == 0 {
		broker.mu.Unlock()
		t.Fatalf("expected subscriber to register")
	}
	subChan := broker.subscriptions[0]
	broker.mu.Unlock()

	// Inject error event with a non-error object
	subChan <- &meshkitBroker.Message{
		EventType: meshkitBroker.ErrorEvent,
		Object:    "this is not an error type",
	}

	// Wait for processing
	time.Sleep(100 * time.Millisecond)
	close(handler.stopCh)

	logged := logger.getLogged()
	if len(logged) == 0 {
		t.Fatalf("expected error to be logged by meshsync handler")
	}

	// Check if the error is wrapped with the correct MeshKit code
	if mkErr, ok := logged[0].(*meshkitErrors.Error); ok {
		assert.Equal(t, ErrMeshsyncStoreUpdatesCode, mkErr.Code)
	} else {
		t.Fatalf("expected logged error to be a *meshkitErrors.Error, got %T", logged[0])
	}
}
