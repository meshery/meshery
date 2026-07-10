package sessions

import (
	"os"
	"strconv"
	"sync"
)

// defaultMaxSessionsPerUser bounds how many terminal and log sessions one user
// may hold at once. Each open session costs a goroutine set on the Meshery
// server and a streaming connection to the target's API server, so an
// unbounded count lets a single user — or a UI bug that reopens a socket in a
// loop — exhaust both.
const defaultMaxSessionsPerUser = 16

// maxSessionsPerUserEnv overrides [defaultMaxSessionsPerUser]. A value of 0 or
// less disables the limit.
const maxSessionsPerUserEnv = "MESHERY_MAX_SESSIONS_PER_USER"

// Limiter caps concurrent sessions per user.
type Limiter struct {
	mu     sync.Mutex
	counts map[string]int
	max    int
}

// NewLimiter returns a limiter allowing max concurrent sessions per user. A max
// of 0 or less disables limiting.
func NewLimiter(max int) *Limiter {
	return &Limiter{counts: make(map[string]int), max: max}
}

// Acquire reserves a session slot for userID. The returned release function is
// idempotent and must be called when the session ends. On refusal it returns an
// ErrSessionLimit-coded error and a no-op release, so callers may defer
// unconditionally.
func (l *Limiter) Acquire(userID string) (release func(), err error) {
	if l.max <= 0 {
		return func() {}, nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.counts[userID] >= l.max {
		return func() {}, ErrSessionLimit(l.max)
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

// Limit is the process-wide session limiter.
var Limit = NewLimiter(maxSessionsPerUserFromEnv())

func maxSessionsPerUserFromEnv() int {
	raw := os.Getenv(maxSessionsPerUserEnv)
	if raw == "" {
		return defaultMaxSessionsPerUser
	}
	max, err := strconv.Atoi(raw)
	if err != nil {
		return defaultMaxSessionsPerUser
	}
	return max
}
