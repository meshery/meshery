package handlers

import (
	"net/http"
	"sync"
	"testing"
)

// mockFlusher implements http.Flusher for testing
type mockFlusher struct {
	flushCount int
	mu         sync.Mutex
}

func (m *mockFlusher) Flush() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.flushCount++
}

func (m *mockFlusher) Header() http.Header {
	return http.Header{}
}

func (m *mockFlusher) Write(b []byte) (int, error) {
	return len(b), nil
}

func (m *mockFlusher) WriteHeader(statusCode int) {}

// resetFlusherMap is a helper to reset the global flusherMap for testing
func resetFlusherMap(t *testing.T) {
	t.Cleanup(func() {
		flusherMapMu.Lock()
		flusherMap = make(map[string]http.Flusher)
		flusherMapMu.Unlock()
	})
	// Also reset at start of test
	flusherMapMu.Lock()
	flusherMap = make(map[string]http.Flusher)
	flusherMapMu.Unlock()
}

func TestSetFlusher(t *testing.T) {
	resetFlusherMap(t)

	flusher := &mockFlusher{}
	setFlusher("test-client", flusher)

	got, ok := getFlusher("test-client")
	if !ok {
		t.Error("Expected flusher to be found, but it wasn't")
	}
	if got != flusher {
		t.Error("Retrieved flusher doesn't match stored flusher")
	}
}

func TestGetFlusher_NotFound(t *testing.T) {
	resetFlusherMap(t)

	_, ok := getFlusher("nonexistent-client")
	if ok {
		t.Error("Expected flusher not to be found, but it was")
	}
}

func TestGetFlusher_EmptyMap(t *testing.T) {
	resetFlusherMap(t)

	flusher, ok := getFlusher("any-client")
	if ok {
		t.Error("Expected flusher not to be found in empty map")
	}
	if flusher != nil {
		t.Error("Expected nil flusher from empty map")
	}
}

func TestDeleteFlusher(t *testing.T) {
	resetFlusherMap(t)

	flusher := &mockFlusher{}
	setFlusher("test-client", flusher)

	// Verify it exists
	_, ok := getFlusher("test-client")
	if !ok {
		t.Fatal("Expected flusher to exist before deletion")
	}

	// Delete it
	deleteFlusher("test-client")

	// Verify it's gone
	_, ok = getFlusher("test-client")
	if ok {
		t.Error("Expected flusher to be deleted, but it still exists")
	}
}

func TestDeleteFlusher_NonexistentKey(t *testing.T) {
	resetFlusherMap(t)

	// Should not panic when deleting nonexistent key
	deleteFlusher("nonexistent-client")
}

// TestConcurrentFlusherAccess verifies thread-safety of flusher operations.
// This test should be run with -race flag: go test -race ./...
func TestConcurrentFlusherAccess(t *testing.T) {
	resetFlusherMap(t)

	const numGoroutines = 100
	const numOperations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // 3 types of operations

	// Concurrent writers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				flusher := &mockFlusher{}
				setFlusher("client-"+string(rune('A'+id%26)), flusher)
			}
		}(i)
	}

	// Concurrent readers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				getFlusher("client-" + string(rune('A'+id%26)))
			}
		}(i)
	}

	// Concurrent deleters
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				deleteFlusher("client-" + string(rune('A'+id%26)))
			}
		}(i)
	}

	wg.Wait()
}

// TestConcurrentReadWrite specifically tests read-write races
func TestConcurrentReadWrite(t *testing.T) {
	resetFlusherMap(t)

	const iterations = 1000
	var wg sync.WaitGroup

	// Start writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			setFlusher("shared-client", &mockFlusher{})
		}
	}()

	// Start multiple readers
	for r := 0; r < 10; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				getFlusher("shared-client")
			}
		}()
	}

	// Start deleter
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			deleteFlusher("shared-client")
		}
	}()

	wg.Wait()
}

// TestMultipleClients tests operations with multiple different clients
func TestMultipleClients(t *testing.T) {
	resetFlusherMap(t)

	clients := []string{"ui", "cli", "api", "mobile", "web"}
	flushers := make(map[string]*mockFlusher)

	// Add all clients
	for _, client := range clients {
		f := &mockFlusher{}
		flushers[client] = f
		setFlusher(client, f)
	}

	// Verify all clients exist
	for _, client := range clients {
		got, ok := getFlusher(client)
		if !ok {
			t.Errorf("Expected client %s to exist", client)
		}
		if got != flushers[client] {
			t.Errorf("Flusher mismatch for client %s", client)
		}
	}

	// Delete some clients
	deleteFlusher("cli")
	deleteFlusher("mobile")

	// Verify deleted clients are gone
	if _, ok := getFlusher("cli"); ok {
		t.Error("Expected cli to be deleted")
	}
	if _, ok := getFlusher("mobile"); ok {
		t.Error("Expected mobile to be deleted")
	}

	// Verify remaining clients still exist
	if _, ok := getFlusher("ui"); !ok {
		t.Error("Expected ui to still exist")
	}
	if _, ok := getFlusher("api"); !ok {
		t.Error("Expected api to still exist")
	}
	if _, ok := getFlusher("web"); !ok {
		t.Error("Expected web to still exist")
	}
}

// BenchmarkSetFlusher benchmarks the set operation
func BenchmarkSetFlusher(b *testing.B) {
	flusher := &mockFlusher{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setFlusher("benchmark-client", flusher)
	}
}

// BenchmarkGetFlusher benchmarks the get operation
func BenchmarkGetFlusher(b *testing.B) {
	flusher := &mockFlusher{}
	setFlusher("benchmark-client", flusher)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getFlusher("benchmark-client")
	}
}

// BenchmarkConcurrentAccess benchmarks concurrent operations
func BenchmarkConcurrentAccess(b *testing.B) {
	flusher := &mockFlusher{}
	setFlusher("benchmark-client", flusher)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			setFlusher("benchmark-client", flusher)
			getFlusher("benchmark-client")
		}
	})
}
