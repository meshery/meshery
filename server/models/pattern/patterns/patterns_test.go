package patterns

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentErrorAppends(t *testing.T) {
	// Simulate the exact scenario from patterns.go where multiple goroutines
	// encounter an error and append to an error channel simultaneously.
	
	numRoutines := 100
	errCh := make(chan error, numRoutines)
	var errs []error
	var wg sync.WaitGroup

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := fmt.Errorf("simulated error %d", id)
			
			errCh <- err
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) != numRoutines {
		t.Errorf("Expected %d errors, got %d", numRoutines, len(errs))
	}

	// Also verify that mergeErrors properly handles the collected errors
	mergedErr := mergeErrors(errs)
	if mergedErr == nil {
		t.Errorf("Expected merged error to not be nil")
	}
}
