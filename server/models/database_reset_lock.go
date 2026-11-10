package models

import "sync"

// resetMu serializes the whole reset workflow: dropping tables, re-migrating,
// and the background seeding that follows.
//
// The dbHandler mutex cannot cover this. It is released when the handler
// returns, but the seeding goroutine outlives the handler - so a second reset
// could acquire it, drop tables mid-seed, and leave partial state behind. This
// lock is instead released by whichever goroutine finishes the work.
var resetMu sync.Mutex

func TryAcquireResetLock() bool {
	return resetMu.TryLock()
}

func ReleaseResetLock() {
	resetMu.Unlock()
}
