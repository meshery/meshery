package models

import (
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/meshery/meshery/server/models/httputil"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
)

type publishMetricsWarnSpy struct {
	logger.Handler
	mu       sync.Mutex
	warnings []error
}

func newPublishMetricsWarnSpy(t *testing.T) *publishMetricsWarnSpy {
	t.Helper()
	return &publishMetricsWarnSpy{Handler: newTestLogger(t)}
}

func (s *publishMetricsWarnSpy) Warn(err error) {
	s.mu.Lock()
	s.warnings = append(s.warnings, err)
	s.mu.Unlock()
	s.Handler.Warn(err)
}

func (s *publishMetricsWarnSpy) recordedWarnings() []error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]error(nil), s.warnings...)
}

func testPublishMetricsResult() *MesheryResult {
	return &MesheryResult{
		Name:   "publish-metrics-test",
		Mesh:   "istio",
		TestID: "test-run-1",
		Result: map[string]interface{}{"qps": 1.0},
	}
}

func assertMeshKitErrorStruct(t *testing.T, err error) *meshkiterrors.Error {
	t.Helper()
	var mkErr *meshkiterrors.Error
	if !stderrors.As(err, &mkErr) || mkErr == nil {
		t.Fatalf("expected *meshkit/errors.Error in error chain, got %T (%v)", err, err)
	}
	return mkErr
}

func TestDefaultLocalProviderPublishMetrics_HTTPSuccess(t *testing.T) {
	type capturedRequest struct {
		method string
		path   string
		apiKey string
		body   string
	}
	reqCh := make(chan capturedRequest, 1)
	remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		reqCh <- capturedRequest{
			method: r.Method,
			path:   r.URL.Path,
			apiKey: r.Header.Get("X-API-Key"),
			body:   string(body),
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer remote.Close()

	spy := newPublishMetricsWarnSpy(t)
	provider := &DefaultLocalProvider{
		ProviderBaseURL: remote.URL,
		Log:             spy,
	}

	if err := provider.PublishMetrics("ignored-token", testPublishMetricsResult()); err != nil {
		t.Fatalf("PublishMetrics() unexpected error: %v", err)
	}

	var got capturedRequest
	select {
	case got = <-reqCh:
	default:
		t.Fatal("expected HTTP request to be received by test server")
	}

	if got.method != http.MethodPut {
		t.Errorf("request method = %q, want %q", got.method, http.MethodPut)
	}
	if got.path != "/result/metrics" {
		t.Errorf("request path = %q, want %q", got.path, "/result/metrics")
	}
	if got.apiKey != GlobalTokenForAnonymousResults {
		t.Errorf("X-API-Key = %q, want %q", got.apiKey, GlobalTokenForAnonymousResults)
	}
	if !strings.Contains(got.body, `"testId":"test-run-1"`) {
		t.Errorf("request body %q does not contain serialized testId", got.body)
	}
	if warnings := spy.recordedWarnings(); len(warnings) != 0 {
		t.Fatalf("expected 0 warnings on HTTP 200, got %d: %v", len(warnings), warnings)
	}
}

func TestDefaultLocalProviderPublishMetrics_Non200Status(t *testing.T) {
	cases := []struct {
		name           string
		statusCode     int
		responseBody   string
		wantObj        string
		wantBodySubstr string
	}{
		{
			name:           "HTTP 401 with JSON error body",
			statusCode:     http.StatusUnauthorized,
			responseBody:   `{"error":"invalid api key"}`,
			wantObj:        `metrics: {"error":"invalid api key"}`,
			wantBodySubstr: "invalid api key",
		},
		{
			name:           "HTTP 500 with text body",
			statusCode:     http.StatusInternalServerError,
			responseBody:   "upstream metrics store unavailable\n",
			wantObj:        "metrics: upstream metrics store unavailable",
			wantBodySubstr: "upstream metrics store unavailable",
		},
		{
			name:         "HTTP 502 with empty body falls back without fabricated text",
			statusCode:   http.StatusBadGateway,
			responseBody: "",
			wantObj:      "metrics",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.statusCode)
				if tc.responseBody != "" {
					_, _ = io.WriteString(w, tc.responseBody)
				}
			}))
			defer remote.Close()

			spy := newPublishMetricsWarnSpy(t)
			provider := &DefaultLocalProvider{
				ProviderBaseURL: remote.URL,
				Log:             spy,
			}

			if err := provider.PublishMetrics("ignored-token", testPublishMetricsResult()); err != nil {
				t.Fatalf("PublishMetrics() unexpected error: %v", err)
			}

			warnings := spy.recordedWarnings()
			if len(warnings) != 1 {
				t.Fatalf("expected 1 warning, got %d: %v", len(warnings), warnings)
			}
			warnErr := warnings[0]

			code := meshkiterrors.GetCode(warnErr)
			if code != ErrPostCode {
				t.Fatalf("warning error code = %q, want %q (err: %v)", code, ErrPostCode, warnErr)
			}
			if code == ErrDoRequestCode || code == ErrDataReadCode {
				t.Fatalf("warning misclassified under %q", code)
			}

			gotStatus, ok := httputil.ProviderStatusCode(warnErr)
			if !ok || gotStatus != tc.statusCode {
				t.Errorf("ProviderStatusCode = (%d, %v), want (%d, true)", gotStatus, ok, tc.statusCode)
			}

			mkErr := assertMeshKitErrorStruct(t, warnErr)
			if len(mkErr.ShortDescription) < 2 || mkErr.ShortDescription[1] != tc.wantObj {
				t.Errorf("ShortDescription = %#v, want index 1 = %q", mkErr.ShortDescription, tc.wantObj)
			}
			if tc.wantBodySubstr != "" && !strings.Contains(meshkiterrors.GetSDescription(warnErr), tc.wantBodySubstr) {
				t.Errorf("short description %q does not contain %q", meshkiterrors.GetSDescription(warnErr), tc.wantBodySubstr)
			}

			longDesc := strings.Join(mkErr.LongDescription, " ")
			wantStatusText := fmt.Sprintf("%d", tc.statusCode)
			if !strings.Contains(longDesc, wantStatusText) {
				t.Errorf("LongDescription %q does not include status %s", longDesc, wantStatusText)
			}
			if !strings.Contains(longDesc, http.MethodPut) {
				t.Errorf("LongDescription %q does not include method %s", longDesc, http.MethodPut)
			}
			wantEndpoint := remote.URL + "/result/metrics"
			if !strings.Contains(longDesc, wantEndpoint) {
				t.Errorf("LongDescription %q does not include endpoint %q", longDesc, wantEndpoint)
			}
			if len(mkErr.ProbableCause) != 0 {
				t.Errorf("expected no TCP connectivity probable causes, got %#v", mkErr.ProbableCause)
			}
		})
	}
}

func TestDefaultLocalProviderPublishMetrics_BodyLimitAndTruncation(t *testing.T) {
	const limit = 1024
	exactBody := strings.Repeat("a", limit)
	overLimitTail := "SHOULD_NOT_BE_LOGGED"
	overLimitBody := strings.Repeat("x", limit) + overLimitTail
	utf8Prefix := strings.Repeat("u", limit-1)
	utf8SplitRune := "€"
	utf8CrossBody := utf8Prefix + utf8SplitRune + overLimitTail

	cases := []struct {
		name           string
		responseBody   string
		wantObj        string
		wantTruncated  bool
		wantUTF8Prefix string
	}{
		{
			name:          "exactly 1024 bytes is preserved without truncation",
			responseBody:  exactBody,
			wantObj:       "metrics: " + exactBody,
			wantTruncated: false,
		},
		{
			name:          "over 1024 bytes includes only first 1024 bytes with truncation marker",
			responseBody:  overLimitBody,
			wantObj:       "metrics: " + strings.Repeat("x", limit) + "... (truncated)",
			wantTruncated: true,
		},
		{
			name:           "multi-byte UTF-8 rune crossing 1024-byte boundary is trimmed cleanly",
			responseBody:   utf8CrossBody,
			wantObj:        "metrics: " + utf8Prefix + "... (truncated)",
			wantTruncated:  true,
			wantUTF8Prefix: utf8Prefix,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, tc.responseBody)
			}))
			defer remote.Close()

			spy := newPublishMetricsWarnSpy(t)
			provider := &DefaultLocalProvider{
				ProviderBaseURL: remote.URL,
				Log:             spy,
			}

			if err := provider.PublishMetrics("ignored-token", testPublishMetricsResult()); err != nil {
				t.Fatalf("PublishMetrics() unexpected error: %v", err)
			}

			warnings := spy.recordedWarnings()
			if len(warnings) != 1 {
				t.Fatalf("expected 1 warning, got %d: %v", len(warnings), warnings)
			}
			warnErr := warnings[0]
			if got := meshkiterrors.GetCode(warnErr); got != ErrPostCode {
				t.Fatalf("warning error code = %q, want %q", got, ErrPostCode)
			}

			mkErr := assertMeshKitErrorStruct(t, warnErr)
			if len(mkErr.ShortDescription) < 2 {
				t.Fatalf("ShortDescription = %#v, want at least 2 entries", mkErr.ShortDescription)
			}
			if mkErr.ShortDescription[1] != tc.wantObj {
				t.Fatalf("ShortDescription[1] mismatch: got len=%d val=%q, want len=%d", len(mkErr.ShortDescription[1]), mkErr.ShortDescription[1], len(tc.wantObj))
			}
			retainedBody := strings.TrimPrefix(mkErr.ShortDescription[1], "metrics: ")
			retainedBody = strings.TrimSuffix(retainedBody, "... (truncated)")
			if !utf8.ValidString(retainedBody) {
				t.Errorf("retained response body is not valid UTF-8: %q", retainedBody)
			}
			if len(retainedBody) > limit {
				t.Errorf("retained response body length = %d, exceeds %d-byte limit", len(retainedBody), limit)
			}
			if tc.wantTruncated {
				if strings.Contains(warnErr.Error(), overLimitTail) {
					t.Errorf("logged error leaked bytes beyond %d-byte limit", limit)
				}
				if !strings.Contains(mkErr.ShortDescription[1], "... (truncated)") {
					t.Errorf("expected truncation indicator in %q", mkErr.ShortDescription[1])
				}
			} else if strings.Contains(mkErr.ShortDescription[1], "truncated") {
				t.Errorf("unexpected truncation indicator for %d-byte body", limit)
			}
			if tc.wantUTF8Prefix != "" {
				if !strings.HasPrefix(retainedBody, tc.wantUTF8Prefix) {
					t.Errorf("retained body did not preserve complete ASCII prefix of length %d", len(tc.wantUTF8Prefix))
				}
				if strings.ContainsRune(mkErr.ShortDescription[1], utf8.RuneError) || strings.Contains(mkErr.ShortDescription[1], utf8SplitRune) {
					t.Errorf("logged object emitted incomplete or replacement UTF-8 character: %q", mkErr.ShortDescription[1])
				}
			}
		})
	}
}

func TestDefaultLocalProviderPublishMetrics_BodyReadFailure(t *testing.T) {
	hijackErrCh := make(chan error, 1)
	remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			hijackErrCh <- fmt.Errorf("ResponseWriter %T does not implement http.Hijacker", w)
			return
		}
		conn, bufrw, err := hj.Hijack()
		if err != nil {
			hijackErrCh <- fmt.Errorf("Hijack failed: %w", err)
			return
		}
		if _, err := bufrw.WriteString("HTTP/1.1 500 Internal Server Error\r\nContent-Length: 100\r\nConnection: close\r\n\r\npartial"); err != nil {
			_ = conn.Close()
			hijackErrCh <- fmt.Errorf("WriteString failed: %w", err)
			return
		}
		if err := bufrw.Flush(); err != nil {
			_ = conn.Close()
			hijackErrCh <- fmt.Errorf("Flush failed: %w", err)
			return
		}
		// Signal handler setup success before closing the socket so PublishMetrics
		// cannot observe EOF and return before hijackErrCh is populated.
		hijackErrCh <- nil
		_ = conn.Close()
	}))
	defer remote.Close()

	spy := newPublishMetricsWarnSpy(t)
	provider := &DefaultLocalProvider{
		ProviderBaseURL: remote.URL,
		Log:             spy,
	}

	if err := provider.PublishMetrics("ignored-token", testPublishMetricsResult()); err != nil {
		t.Fatalf("PublishMetrics() unexpected error: %v", err)
	}

	select {
	case herr := <-hijackErrCh:
		if herr != nil {
			t.Fatalf("test server setup failed: %v", herr)
		}
	default:
		t.Fatal("expected hijacked test server handler to run")
	}

	warnings := spy.recordedWarnings()
	if len(warnings) == 0 {
		t.Fatal("expected at least 1 warning on response-body read failure, got 0")
	}

	var dataReadWarnings int
	for _, warnErr := range warnings {
		code := meshkiterrors.GetCode(warnErr)
		switch code {
		case ErrDataReadCode:
			dataReadWarnings++
			mkErr := assertMeshKitErrorStruct(t, warnErr)
			longDesc := strings.Join(mkErr.LongDescription, " ")
			if !strings.Contains(longDesc, "EOF") {
				t.Errorf("ErrDataRead LongDescription %q does not capture premature EOF cause", longDesc)
			}
		case ErrPostCode, ErrDoRequestCode:
			t.Fatalf("body read failure was misclassified under %q (err: %v)", code, warnErr)
		default:
			if !strings.Contains(warnErr.Error(), "error closing response body") {
				t.Fatalf("unexpected untyped warning on body read failure: %v", warnErr)
			}
		}
	}
	if dataReadWarnings != 1 {
		t.Errorf("expected 1 ErrDataRead warning, got %d (warnings: %v)", dataReadWarnings, warnings)
	}
}

func TestDefaultLocalProviderPublishMetrics_TransportFailure(t *testing.T) {
	remote := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	closedURL := remote.URL
	remote.Close()

	spy := newPublishMetricsWarnSpy(t)
	provider := &DefaultLocalProvider{
		ProviderBaseURL: closedURL,
		Log:             spy,
	}

	if err := provider.PublishMetrics("ignored-token", testPublishMetricsResult()); err != nil {
		t.Fatalf("PublishMetrics() unexpected error: %v", err)
	}

	warnings := spy.recordedWarnings()
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning on transport failure, got %d: %v", len(warnings), warnings)
	}
	warnErr := warnings[0]

	code := meshkiterrors.GetCode(warnErr)
	if code != ErrDoRequestCode {
		t.Fatalf("warning error code = %q, want %q (err: %v)", code, ErrDoRequestCode, warnErr)
	}
	if code == ErrPostCode || code == ErrDataReadCode {
		t.Fatalf("transport failure misclassified under %q", code)
	}
	if status, ok := httputil.ProviderStatusCode(warnErr); ok {
		t.Errorf("expected no provider HTTP status for transport failure, got %d", status)
	}
}
