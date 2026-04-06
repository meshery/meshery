package patterns

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentErrorAppends(t *testing.T) {
	// Simulate the exact scenario from patterns.go where multiple goroutines
	// could encounter an error and append to the errs slice simultaneously.
	
	var errsMx sync.Mutex
	var errs []error
	var wg sync.WaitGroup

	numRoutines := 100
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := fmt.Errorf("simulated error %d", id)
			
			// If errsMx is removed, running 'go test -race' will detect a data race here.
			errsMx.Lock()
			errs = append(errs, err)
			errsMx.Unlock()
		}(i)
	}

	wg.Wait()

	if len(errs) != numRoutines {
		t.Errorf("Expected %d errors, got %d", numRoutines, len(errs))
	}

	// Also verify that mergeErrors properly handles the collected errors
	mergedErr := mergeErrors(errs)
	if mergedErr == nil {
		t.Errorf("Expected merged error to not be nil")
	}
}
