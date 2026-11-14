package router

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	"github.com/sirupsen/logrus"
)

type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func capturingLogger(t *testing.T) (logger.Handler, *lockedBuffer) {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	log.SetLevel(logrus.ErrorLevel) // a zero-value logger would drop Error
	out := &lockedBuffer{}
	log.UpdateErrorLogOutput(out)
	return log, out
}

func newTestRouter(log logger.Handler) *mux.Router {
	r := mux.NewRouter()
	r.Use(recoveryMiddleware(log))
	return r
}

func TestRecovery_PassesThroughNormalResponses(t *testing.T) {
	log, out := capturingLogger(t)
	r := newTestRouter(log)
	r.HandleFunc("/ok", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("fine"))
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ok", nil))

	if rec.Code != http.StatusCreated || rec.Body.String() != "fine" {
		t.Errorf("response altered: %d %q", rec.Code, rec.Body.String())
	}
	if out.String() != "" {
		t.Errorf("nothing should be logged without a panic: %q", out.String())
	}
}

// writeKey writes into m. The test calls it with a nil map to trigger the
// runtime panic; a parameter keeps staticcheck (SA5000) from rejecting a write
// it can see is aimed at a nil map.
func writeKey(m map[string]string, k, v string) { m[k] = v }

func TestRecovery_PanicBecomesStructured500(t *testing.T) {
	log, out := capturingLogger(t)
	r := newTestRouter(log)
	r.HandleFunc("/profiles/{id}", func(http.ResponseWriter, *http.Request) {
		writeKey(nil, "boom", "x") // the kind of defect that has crashed the server
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/profiles/abc-123", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500 (body=%q)", rec.Code, rec.Body.String())
	}

	var body struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not a JSON MeshKit error: %v (body=%q)", err, rec.Body.String())
	}
	if want := errors.GetCode(models.ErrHandlerPanic(nil2())); body.Code != want {
		t.Errorf("error code = %q, want %q", body.Code, want)
	}
	if body.Error == "" {
		t.Error("error message is empty")
	}

	// The stack is for the server log. It must not reach the client.
	if strings.Contains(rec.Body.String(), "goroutine") || strings.Contains(rec.Body.String(), ".go:") {
		t.Errorf("response leaks a stack trace: %q", rec.Body.String())
	}
	logged := out.String()
	if !strings.Contains(logged, "recovered panic in http GET /profiles/{id}") || !strings.Contains(logged, "goroutine") {
		t.Errorf("log should carry the route template and the stack: %q", logged)
	}
	// The label is the route template, not the request path, so ids do not
	// make a new counter key per request.
	if strings.Contains(logged, "abc-123") && !strings.Contains(logged, "/profiles/{id}") {
		t.Errorf("label should use the route template: %q", logged)
	}
	if n := utils.RecoveredPanicCounts()["http GET /profiles/{id}"]; n != 1 {
		t.Errorf("count = %d, want 1", n)
	}
}

// nil2 gives ErrHandlerPanic a non-nil cause, only so the test can read its code.
func nil2() error { return http.ErrNotSupported }

func TestRecovery_PanicAfterResponseStartedKeepsTheStatus(t *testing.T) {
	log, _ := capturingLogger(t)
	r := newTestRouter(log)
	r.HandleFunc("/partial", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("partial"))
		panic("late failure")
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/partial", nil))

	if rec.Code != http.StatusAccepted {
		t.Errorf("status = %d, want the already-sent 202", rec.Code)
	}
	if rec.Body.String() != "partial" {
		t.Errorf("an error body was spliced into a started response: %q", rec.Body.String())
	}
}

func TestRecovery_AbortHandlerReachesNetHTTP(t *testing.T) {
	log, _ := capturingLogger(t)
	r := newTestRouter(log)
	r.HandleFunc("/abort", func(http.ResponseWriter, *http.Request) { panic(http.ErrAbortHandler) })

	defer func() {
		if got := recover(); got != http.ErrAbortHandler {
			t.Errorf("recovered %v, want http.ErrAbortHandler to propagate", got)
		}
	}()
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/abort", nil))
	t.Error("http.ErrAbortHandler must not be swallowed")
}

func TestRecovery_NilLoggerStillAnswers500(t *testing.T) {
	r := newTestRouter(nil)
	r.HandleFunc("/nolog", func(http.ResponseWriter, *http.Request) { panic("x") })

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nolog", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

// flushRecorder records whether Flush reached the underlying writer.
type flushRecorder struct {
	*httptest.ResponseRecorder
	flushed bool
}

func (f *flushRecorder) Flush() { f.flushed = true; f.ResponseRecorder.Flush() }

// Streaming endpoints type-assert the writer to http.Flusher. Wrapping must not
// take that away.
func TestRecovery_FlushPassesThrough(t *testing.T) {
	log, _ := capturingLogger(t)
	r := newTestRouter(log)
	r.HandleFunc("/stream", func(w http.ResponseWriter, _ *http.Request) {
		f, ok := w.(http.Flusher)
		if !ok {
			t.Error("the wrapped writer lost http.Flusher")
			return
		}
		_, _ = w.Write([]byte("chunk"))
		f.Flush()
	})

	rec := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/stream", nil))

	if !rec.flushed {
		t.Error("Flush did not reach the underlying writer")
	}
}

// hijackRecorder lets a test observe a WebSocket-style takeover of the connection.
type hijackRecorder struct {
	*httptest.ResponseRecorder
	hijacked bool
	conn     net.Conn
}

func (h *hijackRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.hijacked = true
	server, client := net.Pipe()
	h.conn = client
	return server, bufio.NewReadWriter(bufio.NewReader(server), bufio.NewWriter(server)), nil
}

// GraphQL subscriptions upgrade to WebSocket, which type-asserts http.Hijacker.
// Wrapping must not break that, and a panic after the takeover must not try to
// write an HTTP error onto a connection that is no longer HTTP.
func TestRecovery_HijackPassesThroughAndSuppressesTheErrorBody(t *testing.T) {
	log, _ := capturingLogger(t)
	r := newTestRouter(log)
	r.HandleFunc("/ws", func(w http.ResponseWriter, _ *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Error("the wrapped writer lost http.Hijacker")
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			t.Errorf("Hijack failed: %v", err)
			return
		}
		defer conn.Close()
		panic("failure after the upgrade")
	})

	rec := &hijackRecorder{ResponseRecorder: httptest.NewRecorder()}
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ws", nil))
	if rec.conn != nil {
		defer rec.conn.Close()
	}

	if !rec.hijacked {
		t.Fatal("Hijack did not reach the underlying writer")
	}
	if rec.Body.Len() != 0 || rec.Code != http.StatusOK {
		t.Errorf("an HTTP error was written onto a hijacked connection: %d %q", rec.Code, rec.Body.String())
	}
}

func TestRecovery_HijackWithoutSupportReturnsAnError(t *testing.T) {
	log, _ := capturingLogger(t)
	r := newTestRouter(log)
	r.HandleFunc("/nohijack", func(w http.ResponseWriter, _ *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("wrapper should expose Hijack and report the failure from it")
		}
		if _, _, err := hj.Hijack(); err == nil {
			t.Error("Hijack on a writer that cannot hijack must return an error")
		}
	})

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/nohijack", nil))
}

func TestWithLogger_SetsTheLogger(t *testing.T) {
	log, _ := capturingLogger(t)

	var o routerOptions
	WithLogger(log)(&o)

	if o.log != log {
		t.Error("WithLogger did not set the logger")
	}
}
