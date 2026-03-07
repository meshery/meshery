package models

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// initMeshsyncDataHandler creates a fully isolated in-memory SQLite handler.
//
// SQLite `:memory:` connections within the same process share the same
// database by default, which causes cross-test interference (e.g. one test
// dropping a table breaks every test that runs after it). Using a unique file
// URI per call gives each handler its own private in-memory DB.
func initMeshsyncDataHandler(t testing.TB) (*MeshsyncDataHandler, error) {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=private", uuid.Must(uuid.NewV4()).String())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&component.ComponentDefinition{},
		&model.ModelDefinition{},
		&meshsyncmodel.KubernetesResource{},
		&meshsyncmodel.KubernetesResourceObjectMeta{},
	)
	if err != nil {
		return nil, err
	}

	// Seed a model + component so cache/DB-hit paths have real data.
	m := model.ModelDefinition{Name: "core"}
	db.Create(&m)
	db.Create(&component.ComponentDefinition{
		Component: component.Component{Version: "v1", Kind: "Pod"},
		Model:     &m,
	})

	mockLogger, err := logger.New("test", logger.Options{})
	if err != nil {
		return nil, err
	}

	return NewMeshsyncDataHandler(
		nil, database.Handler{DB: db}, mockLogger, nil,
		uuid.Nil, uuid.Nil, uuid.Nil, "", nil,
	), nil
}

// makeResource creates a KubernetesResource with an explicit string ID.
// Using explicit IDs lets tests look objects up by primary key directly,
// avoiding brittle JOIN-based WHERE clauses whose alias names change with
// GORM versions.
func makeResource(id, name, namespace, kind, apiVersion string) meshsyncmodel.KubernetesResource {
	return meshsyncmodel.KubernetesResource{
		ID: id,
		KubernetesResourceMeta: &meshsyncmodel.KubernetesResourceObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Kind:       kind,
		APIVersion: apiVersion,
	}
}

// seedResources inserts a slice of resources directly into the DB.
func seedResources(t *testing.T, h *MeshsyncDataHandler, objs []meshsyncmodel.KubernetesResource) {
	t.Helper()
	for i := range objs {
		if err := h.dbHandler.Create(&objs[i]).Error; err != nil {
			t.Fatalf("seedResources: %v", err)
		}
	}
}

// drainRegistrationQueue starts a goroutine that discards every entry sent to
// the global MeshSyncRegistrationQueue for the duration of the test.
//
// Why this is necessary:
//   - The global queue channel has a buffer of 10 (see meshsync_register_queue.go).
//   - processAddUpdateBatch calls Send() after every successful write. With no
//     consumer running, the channel fills after 10 items and Send() blocks
//     forever, hanging the test.
//   - The queue is a singleton (sync.Once), so we cannot replace it per test;
//     we must consume from the shared channel instead.
//
// t.Cleanup stops the goroutine after each test so no goroutine leak occurs.
func drainRegistrationQueue(t *testing.T) {
	t.Helper()

	stop := make(chan struct{})
	go func() {
		q := GetMeshSyncRegistrationQueue()
		for {
			select {
			case <-q.RegChan:
				// discard — we only care about DB state in these tests
			case <-stop:
				return
			}
		}
	}()

	t.Cleanup(func() { close(stop) })
}

// ---------------------------------------------------------------------------
// Metadata-cache tests
// ---------------------------------------------------------------------------

func TestGetComponentMetadataWithNilModel(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	t.Run("RecordNotFound_DoesNotPanic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panicked: %v", r)
			}
		}()
		data, modelName := h.getComponentMetadata("does/not/exist", "Ghost")
		if data == nil {
			t.Error("expected non-nil fallback data on miss")
		}
		if modelName != "" {
			t.Errorf("expected empty model name on miss, got %q", modelName)
		}
	})
}

// TestCacheHitReturnsSameData confirms the second call returns identical
// data/model without going back to the database.
func TestCacheHitReturnsSameData(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	data1, model1 := h.getComponentMetadata("v1", "Pod")
	data2, model2 := h.getComponentMetadata("v1", "Pod")

	if model1 != model2 {
		t.Errorf("model mismatch: first=%q second=%q", model1, model2)
	}
	if len(data1) != len(data2) {
		t.Errorf("data length mismatch: first=%d second=%d", len(data1), len(data2))
	}
}

// TestCacheAvoidsDatabaseAfterPrime drops the table after priming the cache
// and confirms the second call never hits the DB.
func TestCacheAvoidsDatabaseAfterPrime(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	_, _ = h.getComponentMetadata("v1", "Pod") // prime

	// Destroy the table — any subsequent DB hit would panic or error.
	h.dbHandler.Migrator().DropTable(&component.ComponentDefinition{})

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panicked after table drop — hit DB instead of cache: %v", r)
		}
	}()

	data, _ := h.getComponentMetadata("v1", "Pod")
	if data == nil {
		t.Error("expected cached data, got nil")
	}
}

// TestCacheMissAfterEviction evicts a key and confirms the next call
// re-populates the cache.
func TestCacheMissAfterEviction(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	h.getComponentMetadata("v1", "Pod")
	if _, ok := h.metadataCache.Get("v1/Pod"); !ok {
		t.Fatal("expected cache to contain v1/Pod after first call")
	}

	h.metadataCache.Remove("v1/Pod")

	data, _ := h.getComponentMetadata("v1", "Pod")
	if data == nil {
		t.Error("expected non-nil data after re-population")
	}
	if _, ok := h.metadataCache.Get("v1/Pod"); !ok {
		t.Error("expected cache to be repopulated after miss")
	}
}

// TestCacheLRUEvictsOldestEntry fills the cache past its 200-entry capacity
// and verifies the first-inserted key is evicted.
func TestCacheLRUEvictsOldestEntry(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	const cacheCapacity = 200
	firstKey := "v1/First"
	h.metadataCache.Add(firstKey, metadataCacheEntry{data: map[string]any{"x": 1}, model: "first"})

	for i := 1; i <= cacheCapacity; i++ {
		h.metadataCache.Add(fmt.Sprintf("kind-%d/v1", i), metadataCacheEntry{})
	}

	if _, ok := h.metadataCache.Get(firstKey); ok {
		t.Error("expected oldest entry to be evicted after exceeding capacity")
	}
}

// TestCacheConcurrentAccess hammers getComponentMetadata from 50 goroutines.
// Run with -race to catch data races.
func TestCacheConcurrentAccess(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.getComponentMetadata("v1", "Pod")
		}()
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// processAddUpdateBatch tests
// ---------------------------------------------------------------------------

func TestProcessAddBatch_Empty(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panicked on empty batch: %v", r)
		}
	}()
	h.processAddUpdateBatch(nil)
	h.processAddUpdateBatch([]meshsyncmodel.KubernetesResource{})
}

// TestProcessAddBatch_PersistsAllObjects confirms every object in the batch
// ends up in the DB. Uses drainRegistrationQueue to prevent Send() blocking.
func TestProcessAddBatch_PersistsAllObjects(t *testing.T) {
	drainRegistrationQueue(t)

	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	objs := []meshsyncmodel.KubernetesResource{
		makeResource("uid-1", "pod-1", "default", "Pod", "v1"),
		makeResource("uid-2", "pod-2", "default", "Pod", "v1"),
	}

	h.processAddUpdateBatch(objs)

	for _, obj := range objs {
		var result meshsyncmodel.KubernetesResource
		if err := h.dbHandler.First(&result, "id = ?", obj.ID).Error; err != nil {
			t.Errorf("expected object %q to be persisted, got: %v", obj.ID, err)
		}
	}
}

// TestProcessAddBatch_Idempotent confirms that upserting the same objects
// twice does not create duplicate rows.
func TestProcessAddBatch_Idempotent(t *testing.T) {
	drainRegistrationQueue(t)

	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	objs := []meshsyncmodel.KubernetesResource{
		makeResource("uid-1", "pod-1", "default", "Pod", "v1"),
	}

	h.processAddUpdateBatch(objs)
	h.processAddUpdateBatch(objs)

	var count int64
	h.dbHandler.Model(&meshsyncmodel.KubernetesResource{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 row after double upsert, got %d", count)
	}
}

// ---------------------------------------------------------------------------
// processDeleteBatch tests
// ---------------------------------------------------------------------------

func TestProcessDeleteBatch_Empty(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panicked on empty batch: %v", r)
		}
	}()
	h.processDeleteBatch(nil)
	h.processDeleteBatch([]meshsyncmodel.KubernetesResource{})
}

// TestProcessDeleteBatch_RemovesOnlyTargetObjects deletes one specific row
// and confirms the others are untouched.
func TestProcessDeleteBatch_RemovesOnlyTargetObjects(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	toDelete := makeResource("del-id", "to-delete", "default", "Pod", "v1")
	keeper1 := makeResource("keep-1", "keeper-1", "default", "Pod", "v1")
	keeper2 := makeResource("keep-2", "keeper-2", "default", "Pod", "v1")
	seedResources(t, h, []meshsyncmodel.KubernetesResource{toDelete, keeper1, keeper2})

	h.processDeleteBatch([]meshsyncmodel.KubernetesResource{toDelete})

	var remaining int64
	h.dbHandler.Model(&meshsyncmodel.KubernetesResource{}).Count(&remaining)
	if remaining != 2 {
		t.Errorf("expected 2 remaining rows, got %d", remaining)
	}

	var gone meshsyncmodel.KubernetesResource
	if err := h.dbHandler.First(&gone, "id = ?", "del-id").Error; err == nil {
		t.Error("expected deleted object to be absent, but it was still found")
	}
}

// TestProcessDeleteBatch_AllGone confirms a full-batch delete empties the table.
func TestProcessDeleteBatch_AllGone(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	objs := []meshsyncmodel.KubernetesResource{
		makeResource("r1", "res-1", "default", "Pod", "v1"),
		makeResource("r2", "res-2", "default", "Pod", "v1"),
	}
	seedResources(t, h, objs)

	h.processDeleteBatch(objs)

	var count int64
	h.dbHandler.Model(&meshsyncmodel.KubernetesResource{}).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 rows after full delete, got %d", count)
	}
}

// ---------------------------------------------------------------------------
// eventProcessor tests
// ---------------------------------------------------------------------------

// TestEventProcessor_TickerFlushesSmallBatch sends 3 events (well below
// maxBatchSize=100) and waits for the 500ms ticker to flush them.
func TestEventProcessor_TickerFlushesSmallBatch(t *testing.T) {
	drainRegistrationQueue(t)

	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	go h.eventProcessor()

	for i := 0; i < 3; i++ {
		h.eventBuffer <- meshsyncInternalEvent{
			EventType: broker.Add,
			Object:    makeResource(fmt.Sprintf("ticker-uid-%d", i), fmt.Sprintf("ticker-pod-%d", i), "default", "Pod", "v1"),
		}
	}

	time.Sleep(700 * time.Millisecond)

	var count int64
	h.dbHandler.Model(&meshsyncmodel.KubernetesResource{}).Count(&count)
	if count != 3 {
		t.Errorf("expected 3 rows after ticker flush, got %d", count)
	}
}

// TestEventProcessor_MaxBatchSizeFlush sends exactly maxBatchSize events and
// expects them to flush before the 500ms ticker fires.
func TestEventProcessor_MaxBatchSizeFlush(t *testing.T) {
	drainRegistrationQueue(t)

	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	go h.eventProcessor()

	const maxBatchSize = 100
	for i := 0; i < maxBatchSize; i++ {
		h.eventBuffer <- meshsyncInternalEvent{
			EventType: broker.Add,
			Object:    makeResource(fmt.Sprintf("batch-uid-%d", i), fmt.Sprintf("batch-pod-%d", i), "default", "Pod", "v1"),
		}
	}

	// Count-based flush should fire before the 500ms ticker.
	time.Sleep(100 * time.Millisecond)

	var count int64
	h.dbHandler.Model(&meshsyncmodel.KubernetesResource{}).Count(&count)
	if count != maxBatchSize {
		t.Errorf("expected %d rows after max-batch flush, got %d", maxBatchSize, count)
	}
}

// TestEventProcessor_DeleteEventRemovesObjects confirms Delete events routed
// through the buffer actually remove rows.
// Delete events do not touch the registration queue, so no drain is needed.
func TestEventProcessor_DeleteEventRemovesObjects(t *testing.T) {
	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	seed := makeResource("del-uid", "to-remove", "default", "Pod", "v1")
	seedResources(t, h, []meshsyncmodel.KubernetesResource{seed})

	go h.eventProcessor()

	h.eventBuffer <- meshsyncInternalEvent{
		EventType: broker.Delete,
		Object:    seed,
	}

	time.Sleep(700 * time.Millisecond)

	var count int64
	h.dbHandler.Model(&meshsyncmodel.KubernetesResource{}).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 rows after delete event, got %d", count)
	}
}

// TestEventProcessor_MixedEvents sends one Add and one Delete in the same
// batch window and expects a net row count of 1.
func TestEventProcessor_MixedEvents(t *testing.T) {
	drainRegistrationQueue(t)

	h, err := initMeshsyncDataHandler(t)
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	existing := makeResource("exist-uid", "existing", "default", "Pod", "v1")
	seedResources(t, h, []meshsyncmodel.KubernetesResource{existing})

	go h.eventProcessor()

	h.eventBuffer <- meshsyncInternalEvent{
		EventType: broker.Add,
		Object:    makeResource("new-uid", "new-pod", "default", "Pod", "v1"),
	}
	h.eventBuffer <- meshsyncInternalEvent{
		EventType: broker.Delete,
		Object:    existing,
	}

	time.Sleep(700 * time.Millisecond)

	var count int64
	h.dbHandler.Model(&meshsyncmodel.KubernetesResource{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 row after mixed add+delete, got %d", count)
	}
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func BenchmarkMetadataCache(b *testing.B) {
	h, err := initMeshsyncDataHandler(b)
	if err != nil {
		b.Fatalf("init: %v", err)
	}

	b.Run("Uncached", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h.metadataCache.Remove("v1/Pod")
			h.getComponentMetadata("v1", "Pod")
		}
	})

	b.Run("Cached", func(b *testing.B) {
		h.getComponentMetadata("v1", "Pod") // prime
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.getComponentMetadata("v1", "Pod")
		}
	})
}

func BenchmarkProcessAddBatch(b *testing.B) {
	// Benchmarks don't use t.Cleanup, so drain manually.
	go func() {
		q := GetMeshSyncRegistrationQueue()
		for range q.RegChan {
		}
	}()

	h, err := initMeshsyncDataHandler(b)
	if err != nil {
		b.Fatalf("init: %v", err)
	}

	const batchSize = 100
	objs := make([]meshsyncmodel.KubernetesResource, batchSize)
	for i := range objs {
		objs[i] = makeResource(
			fmt.Sprintf("bench-uid-%d", i),
			fmt.Sprintf("bench-pod-%d", i),
			"default", "Pod", "v1",
		)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.processAddUpdateBatch(objs)
	}
}
