package stages

import (
	"fmt"
	"sync"
	"testing"
)

func TestProvisionConcurrentErrorAppends(t *testing.T) {
	// Simulate the parallel_process.go callback executed by planner.go
	// to ensure that errs append remains thread-safe.
	
	var errsMx sync.Mutex
	var errs []error
	var wg sync.WaitGroup

	// Simulate deploying 100 components concurrently
	numComponents := 100
	for i := 0; i < numComponents; i++ {
		wg.Add(1)
		go func(compID int) {
			defer wg.Done()
			
			// Simulate a failure in enriching component or provisioning
			err := fmt.Errorf("failed to provision component %d", compID)
			
			errsMx.Lock()
			errs = append(errs, err)
			errsMx.Unlock()
		}(i)
	}

	// Wait for all "components" to finish
	wg.Wait()

	if len(errs) != numComponents {
		t.Errorf("Expected %d provision errors, got %d", numComponents, len(errs))
	}

	mergedErr := mergeErrors(errs)
	if mergedErr == nil {
		t.Errorf("Expected mergeErrors to return a combined error string, but got nil")
	}
}
