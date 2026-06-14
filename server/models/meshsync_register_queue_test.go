package models

import (
	"testing"
	"time"

	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
)

// TestGetMeshSyncRegistrationQueueLazyInit tests that GetMeshSyncRegistrationQueue
// initializes the queue lazily even if InitMeshSyncRegistrationQueue wasn't called
func TestGetMeshSyncRegistrationQueueLazyInit(t *testing.T) {
	// This should not panic even though InitMeshSyncRegistrationQueue may not have been called
	queue := GetMeshSyncRegistrationQueue()
	if queue == nil {
		t.Fatal("GetMeshSyncRegistrationQueue returned nil")
	}
	if queue.RegChan == nil {
		t.Fatal("RegChan is nil")
	}
}

// TestSendWithNilQueue tests that Send doesn't panic when called on a nil queue
func TestSendWithNilQueue(t *testing.T) {
	var nilQueue *MeshSyncRegistrationQueue

	// This should not panic
	nilQueue.Send(MeshSyncRegistrationData{
		Obj: meshsyncmodel.KubernetesResource{},
	})
}

// TestSendWithNilChannel tests that Send doesn't panic when RegChan is nil
func TestSendWithNilChannel(t *testing.T) {
	queue := &MeshSyncRegistrationQueue{
		RegChan: nil,
	}

	// This should not panic
	queue.Send(MeshSyncRegistrationData{
		Obj: meshsyncmodel.KubernetesResource{},
	})
}

// TestSendWithValidQueue tests that Send works correctly with a valid queue
func TestSendWithValidQueue(t *testing.T) {
	queue := GetMeshSyncRegistrationQueue()
	if queue == nil {
		t.Fatal("GetMeshSyncRegistrationQueue returned nil")
	}

	// Send data in a goroutine
	done := make(chan bool)
	go func() {
		queue.Send(MeshSyncRegistrationData{
			Obj: meshsyncmodel.KubernetesResource{
				Kind: "Pod",
			},
		})
		done <- true
	}()

	// Receive data
	select {
	case data := <-queue.RegChan:
		if data.Obj.Kind != "Pod" {
			t.Errorf("Expected Kind 'Pod', got '%s'", data.Obj.Kind)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for data on RegChan")
	}

	// Wait for send to complete
	<-done
}

// TestConcurrentAccess tests that multiple goroutines can safely access the queue
func TestConcurrentAccess(t *testing.T) {
	queue := GetMeshSyncRegistrationQueue()

	// Spawn multiple goroutines trying to send
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			// Multiple goroutines calling GetMeshSyncRegistrationQueue and Send
			q := GetMeshSyncRegistrationQueue()
			q.Send(MeshSyncRegistrationData{
				Obj: meshsyncmodel.KubernetesResource{
					Kind: "Test",
				},
			})
			done <- true
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

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
