package models

import (
	"sync"
	"testing"
	"time"

	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
)

// TestGetMeshSyncRegistrationQueueLazyInit tests that GetMeshSyncRegistrationQueue
// initializes the queue lazily even if InitMeshSyncRegistrationQueue wasn't called.
func TestGetMeshSyncRegistrationQueueLazyInit(t *testing.T) {
	queue := GetMeshSyncRegistrationQueue()
	if queue == nil {
		t.Fatal("GetMeshSyncRegistrationQueue returned nil")
	}
	if queue.RegChan == nil {
		t.Fatal("RegChan is nil")
	}
}

// TestSendWithNilQueue tests that Send doesn't panic and returns false when called on a nil queue.
func TestSendWithNilQueue(t *testing.T) {
	var nilQueue *MeshSyncRegistrationQueue

	enqueued := nilQueue.Send(MeshSyncRegistrationData{
		Obj: meshsyncmodel.KubernetesResource{},
	})
	if enqueued {
		t.Fatal("expected Send to return false on nil queue")
	}
}

// TestSendWithNilChannel tests that Send doesn't panic and returns false when RegChan is nil.
func TestSendWithNilChannel(t *testing.T) {
	queue := &MeshSyncRegistrationQueue{
		RegChan: nil,
	}

	enqueued := queue.Send(MeshSyncRegistrationData{
		Obj: meshsyncmodel.KubernetesResource{},
	})
	if enqueued {
		t.Fatal("expected Send to return false on nil RegChan")
	}
}

// TestSendWithValidQueue tests that Send works correctly with an available queue.
func TestSendWithValidQueue(t *testing.T) {
	queue := NewMeshSyncRegistrationQueue(5)

	done := make(chan bool)
	go func() {
		enqueued := queue.Send(MeshSyncRegistrationData{
			Obj: meshsyncmodel.KubernetesResource{
				Kind: "Pod",
			},
		})
		done <- enqueued
	}()

	select {
	case data := <-queue.RegChan:
		if data.Obj.Kind != "Pod" {
			t.Errorf("Expected Kind 'Pod', got '%s'", data.Obj.Kind)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for data on RegChan")
	}

	enqueued := <-done
	if !enqueued {
		t.Fatal("expected Send to return true on non-full queue")
	}
}

// TestSendNonBlockingWhenQueueFull tests that Send does not block when the queue buffer
// is saturated, dropping the overflow element and returning false immediately.
func TestSendNonBlockingWhenQueueFull(t *testing.T) {
	capacity := 3
	q := NewMeshSyncRegistrationQueue(capacity)

	// Fill the queue to capacity
	for i := 0; i < capacity; i++ {
		enqueued := q.Send(MeshSyncRegistrationData{
			Obj: meshsyncmodel.KubernetesResource{Kind: "Pod"},
		})
		if !enqueued {
			t.Fatalf("expected item %d to be enqueued successfully", i)
		}
	}

	// Channel is now full. Next send MUST drop without blocking and return false.
	start := time.Now()
	enqueued := q.Send(MeshSyncRegistrationData{
		Obj: meshsyncmodel.KubernetesResource{Kind: "OverflowPod"},
	})
	duration := time.Since(start)

	if enqueued {
		t.Fatal("expected Send to return false when queue buffer is full")
	}
	if duration > 100*time.Millisecond {
		t.Fatalf("Send took %v on a full queue; expected non-blocking drop", duration)
	}
	if q.Dropped() != 1 {
		t.Fatalf("expected dropped count 1, got %d", q.Dropped())
	}
}

// TestSendDoesNotLeakGoroutines verifies that concurrent calls to Send on a saturated
// queue never block indefinitely, avoiding unbounded goroutine leaks.
func TestSendDoesNotLeakGoroutines(t *testing.T) {
	capacity := 5
	q := NewMeshSyncRegistrationQueue(capacity)
	const numGoroutines = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			_ = q.Send(MeshSyncRegistrationData{
				Obj: meshsyncmodel.KubernetesResource{Kind: "Resource"},
			})
		}(i)
	}

	// If Send were blocking unconditionally, unconsumed goroutines would hang
	// and wg.Wait() would timeout. With non-blocking Send, all goroutines finish promptly.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All goroutines completed promptly
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for goroutines to complete; Send appears to have blocked indefinitely")
	}
}

// TestNewMeshSyncRegistrationQueueCapacity verifies custom and fallback capacity behavior.
func TestNewMeshSyncRegistrationQueueCapacity(t *testing.T) {
	qDefault := NewMeshSyncRegistrationQueue(0)
	if cap(qDefault.RegChan) != DefaultQueueCapacity {
		t.Fatalf("expected default capacity %d, got %d", DefaultQueueCapacity, cap(qDefault.RegChan))
	}

	qNegative := NewMeshSyncRegistrationQueue(-10)
	if cap(qNegative.RegChan) != DefaultQueueCapacity {
		t.Fatalf("expected default capacity %d, got %d", DefaultQueueCapacity, cap(qNegative.RegChan))
	}

	customCap := 42
	qCustom := NewMeshSyncRegistrationQueue(customCap)
	if cap(qCustom.RegChan) != customCap {
		t.Fatalf("expected custom capacity %d, got %d", customCap, cap(qCustom.RegChan))
	}
}

// TestConcurrentAccess tests that multiple goroutines can safely send and receive from the queue.
func TestConcurrentAccess(t *testing.T) {
	const numGoroutines = 20
	queue := NewMeshSyncRegistrationQueue(numGoroutines)
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			enqueued := queue.Send(MeshSyncRegistrationData{
				Obj: meshsyncmodel.KubernetesResource{
					Kind: "Test",
				},
			})
			done <- enqueued
		}(i)
	}

	// Drain the channel
	received := 0
	timeout := time.After(2 * time.Second)
	for received < numGoroutines {
		select {
		case <-queue.RegChan:
			received++
		case <-timeout:
			t.Fatalf("Timeout: only received %d/%d messages", received, numGoroutines)
		}
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

// TestMeshSyncRegistrationQueueObservability verifies SetLogger and Dropped counter behavior.
func TestMeshSyncRegistrationQueueObservability(t *testing.T) {
	var nilQ *MeshSyncRegistrationQueue
	if nilQ.Dropped() != 0 {
		t.Fatalf("expected 0 dropped for nil queue, got %d", nilQ.Dropped())
	}
	nilQ.SetLogger(nil) // should not panic

	q := NewMeshSyncRegistrationQueue(1)
	q.SetLogger(nil)
	if q.Dropped() != 0 {
		t.Fatalf("expected 0 dropped initially, got %d", q.Dropped())
	}

	q.Send(MeshSyncRegistrationData{Obj: meshsyncmodel.KubernetesResource{Kind: "Pod"}})
	q.Send(MeshSyncRegistrationData{Obj: meshsyncmodel.KubernetesResource{Kind: "Overflow"}})
	if q.Dropped() != 1 {
		t.Fatalf("expected 1 dropped, got %d", q.Dropped())
	}
}
