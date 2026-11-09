package handlers

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

func TestEvaluationTracker_SingleLeader(t *testing.T) {
	tr := newEvaluationTracker()

	leader, wait := tr.acquire("d1")
	if !leader {
		t.Fatal("first caller should be the leader")
	}
	if wait != nil {
		t.Fatal("leader should not have a wait channel")
	}

	// publish with no followers should not panic or block
	tr.publish("d1", evalResult{})

	// after publish, the next caller should again be a leader
	leader2, _ := tr.acquire("d1")
	if !leader2 {
		t.Fatal("after publish, next caller should be a new leader")
	}
}

func TestEvaluationTracker_CoalescesConcurrent(t *testing.T) {
	tr := newEvaluationTracker()

	// Leader acquires first and does not publish yet.
	leader, _ := tr.acquire("d1")
	if !leader {
		t.Fatal("first caller should be the leader")
	}

	const followers = 50
	waits := make([]<-chan evalResult, 0, followers)
	for range followers {
		isLeader, w := tr.acquire("d1")
		if isLeader {
			t.Fatal("subsequent callers should be followers")
		}
		waits = append(waits, w)
	}

	// Leader finishes and publishes once. All followers must receive the same result.
	sentinelErr := errors.New("boom")
	tr.publish("d1", evalResult{err: sentinelErr})

	var wg sync.WaitGroup
	var received int32
	for _, w := range waits {
		wg.Add(1)
		go func(ch <-chan evalResult) {
			defer wg.Done()
			select {
			case r := <-ch:
				if !errors.Is(r.err, sentinelErr) {
					t.Errorf("follower got wrong err: %v", r.err)
				}
				atomic.AddInt32(&received, 1)
			case <-time.After(2 * time.Second):
				t.Error("follower timed out waiting for result")
			}
		}(w)
	}
	wg.Wait()

	if got := atomic.LoadInt32(&received); got != followers {
		t.Fatalf("expected %d followers to receive result, got %d", followers, got)
	}
}

func TestEvaluationTracker_PublishIsIdempotent(t *testing.T) {
	tr := newEvaluationTracker()
	_, _ = tr.acquire("d1")

	tr.publish("d1", evalResult{})
	// second publish must be a no-op (in particular, no panic).
	tr.publish("d1", evalResult{})
}

func TestEvaluationTracker_DistinctDesignsAreIndependent(t *testing.T) {
	tr := newEvaluationTracker()

	leader1, _ := tr.acquire("d1")
	leader2, _ := tr.acquire("d2")
	if !leader1 || !leader2 {
		t.Fatal("different designs should each get their own leader")
	}
}

// The tests below pin the composition of the coalescing key itself. Coalescing
// hands the leader's result to every follower verbatim, so the key is the only
// thing standing between "two requests would produce the same response" and
// "one caller is served another caller's evaluated design". Each test pins a
// DENY: an allow-only test passes against a key that ignores identity too.

const evalKeyTestProvider = "Local"

func TestEvaluationKey_DoesNotCoalesceAcrossUsers(t *testing.T) {
	userA := uuid.FromStringOrNil("11111111-1111-4111-8111-111111111111")
	userB := uuid.FromStringOrNil("22222222-2222-4222-8222-222222222222")
	designID := uuid.FromStringOrNil("33333333-3333-4333-8333-333333333333")
	body := []byte(`{"design":{"id":"33333333-3333-4333-8333-333333333333"},"options":{}}`)

	if evaluationKey(userA, evalKeyTestProvider, designID, body) ==
		evaluationKey(userB, evalKeyTestProvider, designID, body) {
		t.Fatal("two different users must not share a coalescing key: the follower " +
			"would be served the leader's evaluated design, and a design id is not a secret")
	}
}

func TestEvaluationKey_DoesNotCoalesceAcrossProviders(t *testing.T) {
	// On a multi-provider deployment one registered provider must not be able to
	// reach another provider's user by minting a colliding user id.
	userID := uuid.FromStringOrNil("11111111-1111-4111-8111-111111111111")
	designID := uuid.FromStringOrNil("33333333-3333-4333-8333-333333333333")
	body := []byte(`{"design":{"id":"33333333-3333-4333-8333-333333333333"}}`)

	if evaluationKey(userID, "Local", designID, body) ==
		evaluationKey(userID, "Meshery", designID, body) {
		t.Fatal("identities from two different providers must not share a coalescing key")
	}
}

func TestEvaluationKey_DoesNotCoalesceDifferentBodies(t *testing.T) {
	userID := uuid.FromStringOrNil("11111111-1111-4111-8111-111111111111")
	designID := uuid.FromStringOrNil("33333333-3333-4333-8333-333333333333")

	cases := []struct {
		name string
		a, b []byte
	}{
		{
			// Same caller, same design id, edited between the two requests: the
			// second request must not be answered from the first request's design.
			name: "design edited between requests",
			a:    []byte(`{"design":{"id":"33333333-3333-4333-8333-333333333333","components":[]}}`),
			b:    []byte(`{"design":{"id":"33333333-3333-4333-8333-333333333333","components":[{"id":"x"}]}}`),
		},
		{
			// Options change the shape of the response, so they must change the key.
			name: "different evaluation options",
			a:    []byte(`{"design":{"id":"33333333-3333-4333-8333-333333333333"},"options":{"returnDiffOnly":false}}`),
			b:    []byte(`{"design":{"id":"33333333-3333-4333-8333-333333333333"},"options":{"returnDiffOnly":true}}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if evaluationKey(userID, evalKeyTestProvider, designID, tc.a) ==
				evaluationKey(userID, evalKeyTestProvider, designID, tc.b) {
				t.Fatal("requests that would produce different responses must not share a coalescing key")
			}
		})
	}
}

func TestEvaluationKey_IdlessDesignsDoNotShareOneKey(t *testing.T) {
	// PatternFile.ID is a value type, so a request body that omits "id" decodes
	// to the nil UUID. Keying on the design id alone therefore collapsed every
	// id-less design in the process - an unsaved design in any user's designer -
	// onto a single key.
	if got := uuid.Nil.String(); got != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("assumption behind this test broke: nil uuid renders as %q", got)
	}

	userID := uuid.FromStringOrNil("11111111-1111-4111-8111-111111111111")
	alice := []byte(`{"design":{"name":"alice-draft","components":[]}}`)
	bob := []byte(`{"design":{"name":"bob-draft","components":[]}}`)

	if evaluationKey(userID, evalKeyTestProvider, uuid.Nil, alice) ==
		evaluationKey(userID, evalKeyTestProvider, uuid.Nil, bob) {
		t.Fatal("two different id-less designs must not share a coalescing key")
	}
}

func TestEvaluationKey_CoalescesIdenticalRequestFromSameCaller(t *testing.T) {
	// The rage-click guard this whole mechanism exists for must still work: the
	// same caller re-submitting the same bytes has to hit the same key.
	userID := uuid.FromStringOrNil("11111111-1111-4111-8111-111111111111")
	designID := uuid.FromStringOrNil("33333333-3333-4333-8333-333333333333")
	body := []byte(`{"design":{"id":"33333333-3333-4333-8333-333333333333"},"options":{}}`)

	// A separate backing array holding the same bytes, as two requests would have.
	resubmitted := append([]byte(nil), body...)

	if evaluationKey(userID, evalKeyTestProvider, designID, body) !=
		evaluationKey(userID, evalKeyTestProvider, designID, resubmitted) {
		t.Fatal("an identical request from the same caller must coalesce")
	}
}

func TestEvaluationTracker_DoesNotServeOneUsersResultToAnother(t *testing.T) {
	tr := newEvaluationTracker()

	designID := uuid.FromStringOrNil("33333333-3333-4333-8333-333333333333")
	body := []byte(`{"design":{"id":"33333333-3333-4333-8333-333333333333"},"options":{}}`)
	victimKey := evaluationKey(
		uuid.FromStringOrNil("11111111-1111-4111-8111-111111111111"), evalKeyTestProvider, designID, body)
	otherKey := evaluationKey(
		uuid.FromStringOrNil("22222222-2222-4222-8222-222222222222"), evalKeyTestProvider, designID, body)

	if leader, _ := tr.acquire(victimKey); !leader {
		t.Fatal("first caller should be the leader")
	}

	leader, waitCh := tr.acquire(otherKey)
	if !leader {
		t.Fatal("a second user posting the same design id must run its own evaluation; " +
			"as a follower it would receive the first user's design, components and configuration")
	}
	if waitCh != nil {
		t.Fatal("a leader must not be handed a wait channel")
	}

	// The first caller's result must reach that caller's waiters and nobody else.
	tr.publish(victimKey, evalResult{err: errors.New("first caller's evaluation result")})

	tr.mu.Lock()
	_, otherStillInFlight := tr.inFlight[otherKey]
	tr.mu.Unlock()
	if !otherStillInFlight {
		t.Fatal("publishing one caller's key cleared another caller's in-flight entry")
	}
}
