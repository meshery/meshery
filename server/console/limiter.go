package console

import (
	"os"
	"strconv"
	"sync"
)

// defaultMaxConsolesPerUser bounds how many terminal and log consoles one user
// may hold at once. Each open console costs a goroutine set on the Meshery
// server and a streaming connection to the target's API server, so an
// unbounded count lets a single user — or a UI bug that reopens a socket in a
// loop — exhaust both.
const defaultMaxConsolesPerUser = 16

// maxConsolesPerUserEnv overrides [defaultMaxConsolesPerUser]. A value of 0 or
// less disables the limit.
const maxConsolesPerUserEnv = "MESHERY_MAX_SESSIONS_PER_USER"

// Limiter caps concurrent consoles per user.
type Limiter struct {
	mu     sync.Mutex
	counts map[string]int
	max    int
}

// NewLimiter returns a limiter allowing max concurrent consoles per user. A max
// of 0 or less disables limiting.
func NewLimiter(max int) *Limiter {
	return &Limiter{counts: make(map[string]int), max: max}
}

// Acquire reserves a console slot for userID. The returned release function is
// idempotent and must be called when the console ends. On refusal it returns an
// ErrConsoleLimit-coded error and a no-op release, so callers may defer
// unconditionally.
func (l *Limiter) Acquire(userID string) (release func(), err error) {
	if l.max <= 0 {
		return func() {}, nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.counts[userID] >= l.max {
		return func() {}, ErrConsoleLimit(l.max)
	}
	l.counts[userID]++

	var once sync.Once
	return func() {
		once.Do(func() {
			l.mu.Lock()
			defer l.mu.Unlock()
			if l.counts[userID] <= 1 {
				// Drop the key rather than leaving a zero behind, so the map
				// does not grow without bound across a server's lifetime.
				delete(l.counts, userID)
				return
			}
			l.counts[userID]--
		})
	}, nil
}

// Limit is the process-wide console limiter.
var Limit = NewLimiter(maxConsolesPerUserFromEnv())

func maxConsolesPerUserFromEnv() int {
	raw := os.Getenv(maxConsolesPerUserEnv)
	if raw == "" {
		return defaultMaxConsolesPerUser
	}
	max, err := strconv.Atoi(raw)
	if err != nil {
		return defaultMaxConsolesPerUser
	}
	return max
}
