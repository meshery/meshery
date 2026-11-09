package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

type evalResult struct {
	resp pattern.EvaluationResponse
	err  error
}

// evaluationTracker coalesces concurrent, identical evaluations submitted by the
// same caller: the first caller runs the evaluation, the rest wait for its
// result.
//
// There is one tracker per Meshery Server process - Handler is constructed once,
// in cmd/main.go - so inFlight is a single namespace shared by every session of
// every user. Every key handed to it must therefore already be scoped to a
// caller: build them with evaluationKey and nothing else. A key derived from
// request data alone lets one user join another user's in-flight evaluation as a
// follower and be served that user's evaluated design.
type evaluationTracker struct {
	mu       sync.Mutex
	inFlight map[string][]chan evalResult
}

func newEvaluationTracker() *evaluationTracker {
	return &evaluationTracker{
		inFlight: make(map[string][]chan evalResult),
	}
}

// evaluationKey builds the key concurrent relationship evaluations coalesce on.
//
// A follower is handed the leader's result verbatim, so two requests may share a
// key ONLY when they would produce the same response. That makes the key's
// composition a correctness property rather than a cache-tuning detail, and
// every field below is load-bearing:
//
//   - userID scopes the key to the caller. Keying on the design id alone let any
//     authenticated user become a follower of another user's in-flight
//     evaluation - and be served that user's design, its components and their
//     configuration - just by posting the same design id. A design id is not a
//     secret: the server embeds it in the share URL it builds after a deploy.
//   - providerName keeps two identities minted by different remote providers from
//     colliding on a multi-provider deployment, where one registered provider
//     must not be able to name another provider's user id.
//   - The body digest covers the design and the evaluation options both, so a
//     follower is only ever handed a result computed from the bytes it actually
//     sent. It is also what makes an id-less design safe: PatternFile.ID is a
//     value type, so a body that omits "id" decodes to the nil UUID and every
//     such request in the process would otherwise share one key.
//
// designID is carried for no reason beyond keeping a key greppable against the
// debug log; the digest already implies it.
//
// providerName goes last because it is the only variable-length field. Every
// field before it is fixed width, so no two distinct inputs can render as the
// same key.
//
// A client that serializes the same logical request as different bytes simply
// does not coalesce. That forfeits an optimization, never correctness.
func evaluationKey(userID core.Uuid, providerName string, designID core.Uuid, body []byte) string {
	digest := sha256.Sum256(body)
	return userID.String() + "|" + designID.String() + "|" + hex.EncodeToString(digest[:]) + "|" + providerName
}

// acquire returns (leader=true, nil) for the first caller per key;
// subsequent callers get (false, waitCh) and must read one value from waitCh.
//
// key must come from evaluationKey; see the type comment for why.
func (t *evaluationTracker) acquire(key string) (leader bool, wait <-chan evalResult) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.inFlight[key]; exists {
		ch := make(chan evalResult, 1)
		t.inFlight[key] = append(t.inFlight[key], ch)
		return false, ch
	}
	t.inFlight[key] = nil
	return true, nil
}

// publish broadcasts the result to all waiters. Idempotent: subsequent
// calls for the same key after the entry is cleared are no-ops.
func (t *evaluationTracker) publish(key string, result evalResult) {
	t.mu.Lock()
	waiters, ok := t.inFlight[key]
	if ok {
		delete(t.inFlight, key)
	}
	t.mu.Unlock()

	for _, ch := range waiters {
		ch <- result
	}
}
