package models

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	meshkitBroker "github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/database"
	meshkitErrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testMeshsyncBroker struct {
	mu            sync.Mutex
	subscriptions []chan *meshkitBroker.Message
	// subjects counts the active subscriptions per subject, so a test can
	// tell whether a handler is paused (unsubscribed) or receiving.
	subjects map[string]int
	closed   bool
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

func (b *testMeshsyncBroker) SubscribeWithChannel(subject string, _ string, subscription chan *meshkitBroker.Message) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.subjects == nil {
		b.subjects = map[string]int{}
	}
	b.subjects[subject]++
	// A resumed handler re-subscribes the same channel: keep one entry per
	// channel, or CloseConnection would close it twice.
	for _, existing := range b.subscriptions {
		if existing == subscription {
			return nil
		}
	}
	b.subscriptions = append(b.subscriptions, subscription)
	return nil
}

func (b *testMeshsyncBroker) Unsubscribe(subject string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subjects, subject)
	return nil
}

// subscribed reports whether subject currently has a subscription.
func (b *testMeshsyncBroker) subscribed(subject string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.subjects[subject] > 0
}

func (b *testMeshsyncBroker) Info() string {
	return "test-broker"
}

func (b *testMeshsyncBroker) IsConnected() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return !b.closed
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
		database.Handler{DB: db, Mutex: &sync.Mutex{}},
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

func TestGetComponentMetadataWithoutRegistryTablesFallsBackSilently(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	testLogger := &fakeMeshsyncLogger{}
	handler := &MeshsyncDataHandler{
		dbHandler: database.Handler{DB: db},
		log:       testLogger,
	}

	data, model := handler.getComponentMetadata("v1", "Pod")

	if data == nil {
		t.Fatal("expected fallback metadata to be returned")
	}

	if model != "" {
		t.Fatalf("expected empty model name, got %q", model)
	}

	if logged := testLogger.getLogged(); len(logged) != 0 {
		t.Fatalf("expected no error logs when registry tables are absent, got %d", len(logged))
	}
}

func TestGetComponentMetadataReturnsAssociatedModelName(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	if err := db.AutoMigrate(&modelv1beta1.ModelDefinition{}, &component.ComponentDefinition{}); err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	modelID := uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111")
	componentID := uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222")

	modelDef := modelv1beta1.ModelDefinition{
		ID:   modelID,
		Name: "kubernetes",
	}
	if err := db.Create(&modelDef).Error; err != nil {
		t.Fatalf("Failed to create model definition: %v", err)
	}

	componentDef := component.ComponentDefinition{
		ID:      componentID,
		ModelID: &modelID,
		Component: component.Component{
			Kind:    "Pod",
			Version: "v1",
		},
	}
	if err := db.Create(&componentDef).Error; err != nil {
		t.Fatalf("Failed to create component definition: %v", err)
	}

	mockLogger, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	handler := &MeshsyncDataHandler{
		dbHandler: database.Handler{DB: db},
		log:       mockLogger,
	}

	data, modelName := handler.getComponentMetadata("v1", "Pod")
	if modelName != "kubernetes" {
		t.Fatalf("expected associated model name to be returned, got %q", modelName)
	}

	modelData, ok := data["model"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected returned metadata to include the resolved model, got %#v", data["model"])
	}

	if got := modelData["name"]; got != "kubernetes" {
		t.Fatalf("expected returned metadata to include model name %q, got %#v", "kubernetes", got)
	}
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

// TestMeshsyncDataHandlerPauseAndResumeSubscriptions: pausing unsubscribes both
// MeshSync subjects so a reconnecting connection delivers nothing, and
// resuming re-subscribes the same channels so delivery continues. Pausing or
// resuming twice is a no-op, and resuming a stopped handler does not
// re-subscribe it.
func TestMeshsyncDataHandlerPauseAndResumeSubscriptions(t *testing.T) {
	broker := &testMeshsyncBroker{}
	handler := newTestMeshsyncHandler(t, broker, func() {})

	if err := handler.Run(); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	subscribedBoth := func() bool {
		return broker.subscribed(meshsyncEventsSubject) && broker.subscribed(MeshsyncStoreUpdatesSubject)
	}
	unsubscribedBoth := func() bool {
		return !broker.subscribed(meshsyncEventsSubject) && !broker.subscribed(MeshsyncStoreUpdatesSubject)
	}
	// The events listener subscribes from its own goroutine, so Run may
	// return before it has.
	deadline := time.Now().Add(2 * time.Second)
	for !subscribedBoth() && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if !subscribedBoth() {
		t.Fatal("handler did not subscribe to both MeshSync subjects")
	}

	handler.PauseSubscriptions()
	if !unsubscribedBoth() {
		t.Fatal("pausing did not unsubscribe both MeshSync subjects")
	}
	handler.PauseSubscriptions()
	if !unsubscribedBoth() {
		t.Fatal("pausing twice changed the subscriptions")
	}

	handler.ResumeSubscriptions()
	if !subscribedBoth() {
		t.Fatal("resuming did not re-subscribe both MeshSync subjects")
	}
	// The same channels are re-subscribed, not new ones: a second subscription
	// per channel would deliver every message twice.
	broker.mu.Lock()
	n := len(broker.subscriptions)
	broker.mu.Unlock()
	if n != 2 {
		t.Fatalf("%d subscription channels, want 2 after a resume", n)
	}
	handler.ResumeSubscriptions()
	if !subscribedBoth() {
		t.Fatal("resuming twice changed the subscriptions")
	}

	handler.PauseSubscriptions()
	handler.Stop()
	handler.ResumeSubscriptions()
	if !unsubscribedBoth() {
		t.Fatal("resuming a stopped handler re-subscribed it")
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

func (f *fakeMeshsyncLogger) Info(description ...interface{})          {}
func (f *fakeMeshsyncLogger) Infof(format string, args ...interface{}) {}
func (f *fakeMeshsyncLogger) Debug(description ...interface{})         {}
func (f *fakeMeshsyncLogger) Warn(err error)                           {}

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

func TestReplaceMeshsyncResourceMetadata(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&meshsyncmodel.KubernetesResource{}, &meshsyncmodel.KubernetesResourceObjectMeta{}, &meshsyncmodel.KubernetesKeyValue{}, &meshsyncmodel.KubernetesResourceSpec{}, &meshsyncmodel.KubernetesResourceStatus{}); err != nil {
		t.Fatal(err)
	}
	handler := &MeshsyncDataHandler{dbHandler: database.Handler{DB: db}}
	obj := &meshsyncmodel.KubernetesResource{APIVersion: "v1", Kind: "Pod", ClusterID: "cluster", KubernetesResourceMeta: &meshsyncmodel.KubernetesResourceObjectMeta{Name: "pod", Namespace: "default", Labels: []*meshsyncmodel.KubernetesKeyValue{{Kind: "label", Key: "app", Value: "v1"}, {Kind: "label", Key: "tier", Value: "front"}}, Annotations: []*meshsyncmodel.KubernetesKeyValue{{Kind: "annotation", Key: "old", Value: "value"}}}}
	if err := db.Create(obj).Error; err != nil {
		t.Fatal(err)
	}
	obj.KubernetesResourceMeta.Labels = []*meshsyncmodel.KubernetesKeyValue{{Kind: "label", Key: "app", Value: "v2"}}
	obj.KubernetesResourceMeta.Annotations = nil
	if err := handler.replaceMeshsyncResource(obj); err != nil {
		t.Fatal(err)
	}
	var rows []meshsyncmodel.KubernetesKeyValue
	if err := db.Where("id = ?", obj.ID).Find(&rows).Error; err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Key != "app" || rows[0].Value != "v2" {
		t.Fatalf("unexpected metadata: %+v", rows)
	}
}

func TestMeshsyncEventsAccumulatorTombstonesAndRestoresResource(t *testing.T) {
	handler := newTestMeshsyncHandler(t, &testMeshsyncBroker{}, nil)
	db := handler.dbHandler.DB
	if err := db.AutoMigrate(
		&meshsyncmodel.KubernetesResource{},
		&meshsyncmodel.KubernetesResourceObjectMeta{},
		&meshsyncmodel.KubernetesKeyValue{},
		&meshsyncmodel.KubernetesResourceSpec{},
		&meshsyncmodel.KubernetesResourceStatus{},
	); err != nil {
		t.Fatalf("migrate MeshSync resources: %v", err)
	}

	resource := meshsyncmodel.KubernetesResource{
		APIVersion: "v1",
		Kind:       "Pod",
		ClusterID:  "cluster-1",
		KubernetesResourceMeta: &meshsyncmodel.KubernetesResourceObjectMeta{
			Name:      "recreated-pod",
			Namespace: "default",
			Labels: []*meshsyncmodel.KubernetesKeyValue{
				{Kind: meshsyncmodel.KindLabel, Key: "version", Value: "old"},
			},
		},
	}
	identified := resource
	meshsyncmodel.SetID(&identified)

	if err := handler.meshsyncEventsAccumulator(&meshkitBroker.Message{EventType: meshkitBroker.Add, Object: resource}); err != nil {
		t.Fatalf("persist initial add event: %v", err)
	}
	if err := handler.meshsyncEventsAccumulator(&meshkitBroker.Message{EventType: meshkitBroker.Delete, Object: resource}); err != nil {
		t.Fatalf("persist delete event: %v", err)
	}

	var scoped meshsyncmodel.KubernetesResource
	if err := db.First(&scoped, "id = ?", identified.ID).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("normal query after delete returned %v, want gorm.ErrRecordNotFound", err)
	}

	var tombstoned meshsyncmodel.KubernetesResource
	if err := db.Unscoped().Preload("KubernetesResourceMeta").First(&tombstoned, "id = ?", identified.ID).Error; err != nil {
		t.Fatalf("find tombstoned resource: %v", err)
	}
	if !tombstoned.DeletedAt.Valid {
		t.Fatal("delete event did not set DeletedAt")
	}
	if tombstoned.KubernetesResourceMeta == nil {
		t.Fatal("soft delete removed the resource metadata")
	}

	recreated := resource
	recreated.KubernetesResourceMeta = &meshsyncmodel.KubernetesResourceObjectMeta{
		Name:      "recreated-pod",
		Namespace: "default",
		Labels: []*meshsyncmodel.KubernetesKeyValue{
			{Kind: meshsyncmodel.KindLabel, Key: "version", Value: "new"},
		},
	}
	if err := handler.meshsyncEventsAccumulator(&meshkitBroker.Message{EventType: meshkitBroker.Add, Object: recreated}); err != nil {
		t.Fatalf("persist recreated add event: %v", err)
	}

	var restored meshsyncmodel.KubernetesResource
	if err := db.Preload("KubernetesResourceMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
		First(&restored, "id = ?", identified.ID).Error; err != nil {
		t.Fatalf("find restored resource: %v", err)
	}
	if restored.DeletedAt.Valid {
		t.Fatal("recreated resource is still tombstoned")
	}
	if restored.KubernetesResourceMeta == nil || len(restored.KubernetesResourceMeta.Labels) != 1 {
		t.Fatalf("unexpected restored metadata: %+v", restored.KubernetesResourceMeta)
	}
	label := restored.KubernetesResourceMeta.Labels[0]
	if label.Key != "version" || label.Value != "new" {
		t.Fatalf("restored resource kept stale labels: %+v", label)
	}
}

func TestReplaceMeshsyncResourceRejectsMissingIdentity(t *testing.T) {
	handler := &MeshsyncDataHandler{}
	for _, obj := range []*meshsyncmodel.KubernetesResource{
		{},
		{ID: "existing-id"},
	} {
		if err := handler.replaceMeshsyncResource(obj); err == nil {
			t.Fatalf("expected missing identity to be rejected: %+v", obj)
		}
	}
}
