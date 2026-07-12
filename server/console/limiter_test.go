package console

import (
	"sync"
	"testing"

	"github.com/meshery/meshkit/errors"
)

func TestLimiterAllowsUpToMax(t *testing.T) {
	limiter := NewLimiter(2)

	releaseA, err := limiter.Acquire("alice")
	if err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	releaseB, err := limiter.Acquire("alice")
	if err != nil {
		t.Fatalf("second acquire: %v", err)
	}

	if _, err := limiter.Acquire("alice"); err == nil {
		t.Fatal("third acquire succeeded, want a console-limit error")
	} else if code := errors.GetCode(err); code != ErrConsoleLimitCode {
		t.Errorf("error code = %q, want %q", code, ErrConsoleLimitCode)
	}

	// Another user has their own budget.
	if _, err := limiter.Acquire("bob"); err != nil {
		t.Errorf("bob's first acquire: %v", err)
	}

	releaseA()
	if _, err := limiter.Acquire("alice"); err != nil {
		t.Errorf("acquire after release: %v", err)
	}
	releaseB()
}

// TestLimiterRefusalReturnsUsableRelease matters because callers `defer
// release()` before checking the error.
func TestLimiterRefusalReturnsUsableRelease(t *testing.T) {
	limiter := NewLimiter(1)

	release, err := limiter.Acquire("alice")
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	defer release()

	refusedRelease, err := limiter.Acquire("alice")
	if err == nil {
		t.Fatal("second acquire succeeded, want a console-limit error")
	}
	if refusedRelease == nil {
		t.Fatal("refused acquire returned a nil release")
	}
	refusedRelease() // must not free the slot held by the successful acquire

	if _, err := limiter.Acquire("alice"); err == nil {
		t.Error("a refused release freed a slot it never held")
	}
}

func TestLimiterReleaseIsIdempotent(t *testing.T) {
	limiter := NewLimiter(1)

	release, err := limiter.Acquire("alice")
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	release()
	release() // a double release must not create a phantom slot

	if _, err := limiter.Acquire("alice"); err != nil {
		t.Fatalf("acquire after double release: %v", err)
	}
	if _, err := limiter.Acquire("alice"); err == nil {
		t.Error("limit not enforced after a double release")
	}
}

// TestLimiterReleaseDropsKey keeps the counts map from growing without bound
// over a server's lifetime.
func TestLimiterReleaseDropsKey(t *testing.T) {
	limiter := NewLimiter(4)

	release, err := limiter.Acquire("alice")
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	release()

	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	if _, present := limiter.counts["alice"]; present {
		t.Error("counts still holds a zeroed entry for a user with no consoles")
	}
}

func TestLimiterDisabled(t *testing.T) {
	limiter := NewLimiter(0)
	for i := 0; i < 100; i++ {
		if _, err := limiter.Acquire("alice"); err != nil {
			t.Fatalf("acquire %d with limiting disabled: %v", i, err)
		}
	}
}

func TestLimiterIsConcurrencySafe(t *testing.T) {
	const max = 8
	limiter := NewLimiter(max)

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		granted int
	)
	for i := 0; i < max*4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := limiter.Acquire("alice"); err == nil {
				mu.Lock()
				granted++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if granted != max {
		t.Errorf("granted %d consoles concurrently, want exactly %d", granted, max)
	}
}
