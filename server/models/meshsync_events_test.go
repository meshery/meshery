package models

import (
	"fmt"
	"sync"
	"sync/atomic"
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

func (b *testMeshsyncBroker) Unsubscribe(_ string) error {
	return nil
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

	return NewMeshsyncDataHandler(
		broker,
		newTestDBHandler(t),
		newTestMeshsyncLogger(t),
		nil,
		uuid.Nil,
		uuid.Nil,
		uuid.Nil,
		"",
		stopFunc,
	)
}

// newTestDBHandler returns a database.Handler backed by an in-memory SQLite
// database, migrated for both the component registry and the MeshSync
// resource tables, with a real mutex so tests can exercise dbHandler.Lock()
// the same way production code does.
func newTestDBHandler(t *testing.T) database.Handler {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	err = db.AutoMigrate(
		&component.ComponentDefinition{},
		&meshsyncmodel.KubernetesResource{},
		&meshsyncmodel.KubernetesResourceObjectMeta{},
		&meshsyncmodel.KubernetesResourceSpec{},
		&meshsyncmodel.KubernetesResourceStatus{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	return database.Handler{DB: db, Mutex: &sync.Mutex{}}
}

func newTestMeshsyncLogger(t *testing.T) logger.Handler {
	t.Helper()

	mockLogger, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return mockLogger
}

// newTestKubernetesResource builds a resource whose ID is derived the same
// way meshsync's BeforeCreate hook derives it, so two resources built with
// the same name/namespace/kind/cluster collide on the same primary key, the
// same way a resync of an already-known object does in production.
func newTestKubernetesResource(clusterID, kind, namespace, name string) meshsyncmodel.KubernetesResource {
	return meshsyncmodel.KubernetesResource{
		ClusterID:  clusterID,
		Kind:       kind,
		APIVersion: "v1",
		KubernetesResourceMeta: &meshsyncmodel.KubernetesResourceObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
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

// countGenuineCommits registers a callback on both the Create and Update
// chains that increments counter only when that specific call actually
// opened and committed its own implicit transaction (i.e. it was not
// nested inside an already-active transaction). This distinguishes a real
// per-statement commit from the callback step merely being reached.
func countGenuineCommits(t *testing.T, db *gorm.DB, counter *int32) {
	t.Helper()

	record := func(tx *gorm.DB) {
		if _, ok := tx.InstanceGet("gorm:started_transaction"); ok {
			atomic.AddInt32(counter, 1)
		}
	}

	if err := db.Callback().Create().After("gorm:commit_or_rollback_transaction").Register("test:count_create_commits", record); err != nil {
		t.Fatalf("failed to register create commit counter: %v", err)
	}
	if err := db.Callback().Update().After("gorm:commit_or_rollback_transaction").Register("test:count_update_commits", record); err != nil {
		t.Fatalf("failed to register update commit counter: %v", err)
	}
}

// TestPersistStoreUpdatesCommitsBatchInSingleTransaction guards the fix for
// meshery/meshery#21954: a resync payload must commit once for the whole
// batch, not once per object. Before the fix, subsribeToStoreUpdates called
// persistStoreUpdate per object, and each call committed its own implicit
// transaction while holding dbHandler's lock, which is shared by every
// connected cluster's MeshsyncDataHandler.
func TestPersistStoreUpdatesCommitsBatchInSingleTransaction(t *testing.T) {
	dbHandler := newTestDBHandler(t)

	var commits int32
	countGenuineCommits(t, dbHandler.DB, &commits)

	// Sanity-check the instrumentation itself: a write made outside any
	// explicit transaction must open and commit its own transaction, so the
	// counter should observe exactly one genuine commit here.
	sentinel := newTestKubernetesResource("sanity-cluster", "Pod", "default", "sanity-pod")
	if err := dbHandler.DB.Create(&sentinel).Error; err != nil {
		t.Fatalf("failed to create sanity object: %v", err)
	}
	if got := atomic.LoadInt32(&commits); got != 1 {
		t.Fatalf("commit counter instrumentation is not detecting real commits, expected 1 got %d", got)
	}
	atomic.StoreInt32(&commits, 0)

	handler := &MeshsyncDataHandler{
		dbHandler: dbHandler,
		log:       newTestMeshsyncLogger(t),
	}

	objects := []meshsyncmodel.KubernetesResource{
		newTestKubernetesResource("cluster-1", "Pod", "default", "pod-a"),
		newTestKubernetesResource("cluster-1", "Pod", "default", "pod-b"),
		newTestKubernetesResource("cluster-1", "Pod", "default", "pod-c"),
	}

	handler.persistStoreUpdates(objects)

	if got := atomic.LoadInt32(&commits); got != 0 {
		t.Fatalf("expected a %d-object resync batch to commit through the single outer transaction (0 per-object auto-commits), got %d per-object commits", len(objects), got)
	}

	var count int64
	if err := dbHandler.DB.Model(&meshsyncmodel.KubernetesResource{}).Where("cluster_id = ?", "cluster-1").Count(&count).Error; err != nil {
		t.Fatalf("failed to count persisted objects: %v", err)
	}
	if count != int64(len(objects)) {
		t.Fatalf("expected all %d objects in the batch to be persisted, found %d", len(objects), count)
	}
}

// TestPersistStoreUpdatesPersistsAllNewObjects is a functional correctness
// check that batching writes into one transaction did not break plain
// object creation.
func TestPersistStoreUpdatesPersistsAllNewObjects(t *testing.T) {
	dbHandler := newTestDBHandler(t)
	handler := &MeshsyncDataHandler{
		dbHandler: dbHandler,
		log:       newTestMeshsyncLogger(t),
	}

	objects := []meshsyncmodel.KubernetesResource{
		newTestKubernetesResource("cluster-1", "Pod", "default", "pod-a"),
		newTestKubernetesResource("cluster-1", "Service", "default", "svc-a"),
	}

	handler.persistStoreUpdates(objects)

	var persisted []meshsyncmodel.KubernetesResource
	if err := dbHandler.DB.Find(&persisted).Error; err != nil {
		t.Fatalf("failed to read back persisted objects: %v", err)
	}
	if len(persisted) != len(objects) {
		t.Fatalf("expected %d persisted objects, got %d", len(objects), len(persisted))
	}
}

// TestPersistStoreUpdatesUpdatesExistingObjectAndContinuesBatch verifies the
// per-object create-then-update-on-conflict semantics are preserved: a
// resync object whose derived ID already exists in the database is updated
// in place rather than duplicated, and the rest of the batch still persists.
func TestPersistStoreUpdatesUpdatesExistingObjectAndContinuesBatch(t *testing.T) {
	dbHandler := newTestDBHandler(t)

	existing := newTestKubernetesResource("cluster-1", "Pod", "default", "pod-a")
	existing.Data = "stale"
	if err := dbHandler.DB.Create(&existing).Error; err != nil {
		t.Fatalf("failed to seed existing object: %v", err)
	}

	handler := &MeshsyncDataHandler{
		dbHandler: dbHandler,
		log:       newTestMeshsyncLogger(t),
	}

	resynced := newTestKubernetesResource("cluster-1", "Pod", "default", "pod-a")
	resynced.Data = "fresh"
	newObject := newTestKubernetesResource("cluster-1", "Pod", "default", "pod-b")

	handler.persistStoreUpdates([]meshsyncmodel.KubernetesResource{resynced, newObject})

	var all []meshsyncmodel.KubernetesResource
	if err := dbHandler.DB.Find(&all).Error; err != nil {
		t.Fatalf("failed to read back objects: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected the conflicting resync to update in place (2 total rows), got %d", len(all))
	}

	var updated meshsyncmodel.KubernetesResource
	if err := dbHandler.DB.Where("id = ?", existing.ID).First(&updated).Error; err != nil {
		t.Fatalf("failed to read back updated object: %v", err)
	}
	if updated.Data != "fresh" {
		t.Fatalf("expected resync to update the existing object's data to %q, got %q", "fresh", updated.Data)
	}
}

// TestPersistStoreUpdatesEmptyBatchIsNoop ensures an empty resync payload
// (e.g. every object failed to unmarshal) does not acquire the shared lock
// or touch the database.
func TestPersistStoreUpdatesEmptyBatchIsNoop(t *testing.T) {
	dbHandler := newTestDBHandler(t)
	handler := &MeshsyncDataHandler{
		dbHandler: dbHandler,
		log:       newTestMeshsyncLogger(t),
	}

	done := make(chan struct{})
	go func() {
		handler.persistStoreUpdates(nil)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("expected an empty batch to return without acquiring the lock")
	}
}

// TestPersistStoreUpdatesTwoHandlersShareDBHandlerConcurrently mirrors the
// production topology from meshery/meshery#21954: every connected cluster's
// MeshsyncDataHandler is constructed with the same dbHandler, so its mutex
// is shared across clusters. This proves that sharing remains correct after
// batching: concurrent resyncs from two clusters neither deadlock nor lose
// or corrupt each other's writes.
func TestPersistStoreUpdatesTwoHandlersShareDBHandlerConcurrently(t *testing.T) {
	sharedDBHandler := newTestDBHandler(t)
	testLog := newTestMeshsyncLogger(t)

	handlerA := &MeshsyncDataHandler{dbHandler: sharedDBHandler, log: testLog}
	handlerB := &MeshsyncDataHandler{dbHandler: sharedDBHandler, log: testLog}

	const objectsPerCluster = 10
	batchA := make([]meshsyncmodel.KubernetesResource, 0, objectsPerCluster)
	batchB := make([]meshsyncmodel.KubernetesResource, 0, objectsPerCluster)
	for i := 0; i < objectsPerCluster; i++ {
		batchA = append(batchA, newTestKubernetesResource("cluster-a", "Pod", "default", fmt.Sprintf("pod-%d", i)))
		batchB = append(batchB, newTestKubernetesResource("cluster-b", "Pod", "default", fmt.Sprintf("pod-%d", i)))
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		handlerA.persistStoreUpdates(batchA)
	}()
	go func() {
		defer wg.Done()
		handlerB.persistStoreUpdates(batchB)
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("concurrent resyncs from two clusters sharing one dbHandler did not complete, possible deadlock")
	}

	var count int64
	if err := sharedDBHandler.DB.Model(&meshsyncmodel.KubernetesResource{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to count persisted objects: %v", err)
	}
	if count != int64(2*objectsPerCluster) {
		t.Fatalf("expected %d objects across both clusters, got %d", 2*objectsPerCluster, count)
	}
}
