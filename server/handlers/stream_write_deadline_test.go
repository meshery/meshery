package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// The server sets no WriteTimeout, so every long-lived SSE stream must bound
// its own writes: without a per-write deadline, a client that stops reading
// blocks the write - and the handler goroutine behind it - indefinitely. The
// tests in this file check that each stream sets a fresh deadline before
// every frame it writes.

// deadlineRecorder is an SSE-capable http.ResponseWriter that also implements
// SetWriteDeadline, so http.ResponseController reaches it, and logs every
// deadline, write and flush in order.
type deadlineRecorder struct {
	mu     sync.Mutex
	header http.Header
	body   bytes.Buffer
	ops    []deadlineRecorderOp
}

type deadlineRecorderOp struct {
	kind     string // "deadline", "write" or "flush"
	at       time.Time
	deadline time.Time
}

func newDeadlineRecorder() *deadlineRecorder {
	return &deadlineRecorder{header: make(http.Header)}
}

func (r *deadlineRecorder) Header() http.Header { return r.header }
func (r *deadlineRecorder) WriteHeader(int)     {}

// Flush is logged so a handler that flushes its headers before setting any
// deadline fails requireDeadlineBeforeEveryWrite's first-operation check.
func (r *deadlineRecorder) Flush() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ops = append(r.ops, deadlineRecorderOp{kind: "flush", at: time.Now()})
}

func (r *deadlineRecorder) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ops = append(r.ops, deadlineRecorderOp{kind: "write", at: time.Now()})
	return r.body.Write(p)
}

func (r *deadlineRecorder) SetWriteDeadline(d time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ops = append(r.ops, deadlineRecorderOp{kind: "deadline", at: time.Now(), deadline: d})
	return nil
}

func (r *deadlineRecorder) snapshot() (string, []deadlineRecorderOp) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.body.String(), append([]deadlineRecorderOp(nil), r.ops...)
}

// waitForBody polls until the recorded body contains want at least n times.
func (r *deadlineRecorder) waitForBody(t *testing.T, want string, n int) {
	t.Helper()
	deadline := time.After(testTimeout)
	for {
		body, _ := r.snapshot()
		if strings.Count(body, want) >= n {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for %d frames containing %q; got: %q", n, want, body)
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// requireDeadlineBeforeEveryWrite fails unless the first operation on the
// connection - including a header flush - sets a write deadline, and every
// write is preceded by a deadline set since the previous write, each in the
// future and at most limit away.
func requireDeadlineBeforeEveryWrite(t *testing.T, ops []deadlineRecorderOp, limit time.Duration) {
	t.Helper()
	if len(ops) == 0 || ops[0].kind != "deadline" {
		t.Fatalf("operations on the connection = %v, want a write deadline first", opKinds(ops))
	}
	fresh := false
	for i, op := range ops {
		switch op.kind {
		case "deadline":
			if !op.deadline.After(op.at) {
				t.Fatalf("op %d set write deadline %v at %v, want a deadline in the future", i, op.deadline, op.at)
			}
			if op.deadline.After(op.at.Add(limit)) {
				t.Fatalf("op %d set write deadline %v, more than %v after %v", i, op.deadline, limit, op.at)
			}
			fresh = true
		case "write":
			if !fresh {
				t.Fatalf("op %d wrote without setting a fresh write deadline first; ops: %v", i, opKinds(ops))
			}
			fresh = false
		}
	}
}

func opKinds(ops []deadlineRecorderOp) []string {
	kinds := make([]string, len(ops))
	for i, op := range ops {
		kinds[i] = op.kind
	}
	return kinds
}

// TestProvidersStreamHandler_BoundsEveryWrite covers /api/providers/stream,
// which is unauthenticated, including the initial header flush.
func TestProvidersStreamHandler_BoundsEveryWrite(t *testing.T) {
	tracker := models.NewProviderTracker(map[string]models.Provider{}, newTestLogger(t))
	h := &Handler{
		config: &models.HandlerConfig{ProviderTracker: tracker},
		log:    newTestLogger(t),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/api/providers/stream", nil).WithContext(ctx)
	rec := newDeadlineRecorder()

	done := make(chan struct{})
	go func() {
		h.ProvidersStreamHandler(rec, req)
		close(done)
	}()

	// Publish records the status even before the handler subscribes, so the
	// subscription snapshot delivers both events whichever side wins the race.
	tracker.Publish(models.ProviderStatusEvent{Key: "first", Status: models.ProviderStatusOnline})
	tracker.Publish(models.ProviderStatusEvent{Key: "second", Status: models.ProviderStatusOnline})
	rec.waitForBody(t, "event: provider", 2)

	cancel()
	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("handler did not return after context cancellation")
	}

	_, ops := rec.snapshot()
	requireDeadlineBeforeEveryWrite(t, ops, providerStreamWriteTimeout)
}

// TestWriteEventStream_BoundsEveryWrite covers /api/system/adapter/events,
// whose frames are written by writeEventStream.
func TestWriteEventStream_BoundsEveryWrite(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	respChan := make(chan []byte, 2)
	rec := newDeadlineRecorder()
	done := make(chan struct{})
	go func() {
		writeEventStream(ctx, rec, respChan, newTestLogger(t), rec)
		close(done)
	}()

	respChan <- []byte(`{"n":1}`)
	respChan <- []byte(`{"n":2}`)
	rec.waitForBody(t, "data: ", 2)

	cancel()
	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("writeEventStream did not return after context cancellation")
	}

	_, ops := rec.snapshot()
	requireDeadlineBeforeEveryWrite(t, ops, eventStreamWriteTimeout)
}

// TestLoadTestHandler_BoundsEveryWrite runs a one-second load test through
// /api/user/performance/profiles/{id}/run and checks every frame of its
// stream, from the "Initiating" message to the final result.
func TestLoadTestHandler_BoundsEveryWrite(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	profile, err := json.Marshal(models.PerformanceProfile{
		Name:     "write-deadline-profile",
		Metadata: map[string]interface{}{},
	})
	if err != nil {
		t.Fatalf("marshal profile: %v", err)
	}
	h := &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}

	query := url.Values{
		"name": {"write-deadline-test"},
		"url":  {target.URL},
		"t":    {"1"},
		"dur":  {"s"},
		"c":    {"1"},
		"qps":  {"5"},
	}
	profileID := uuid.Must(uuid.NewV4()).String()
	req := httptest.NewRequest(http.MethodGet, "/api/user/performance/profiles/"+profileID+"/run?"+query.Encode(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": profileID})
	rec := newDeadlineRecorder()

	// The handler returns once the run completes and its last frame is written.
	h.LoadTestHandler(rec, req, &models.Preference{}, &models.User{}, &savedProfileProvider{profile: profile})

	body, ops := rec.snapshot()
	if n := strings.Count(body, "data: "); n < 2 {
		t.Fatalf("load test stream wrote %d frames, want the initiating message and a final result; body: %q", n, body)
	}
	requireDeadlineBeforeEveryWrite(t, ops, loadTestStreamWriteTimeout)
}
