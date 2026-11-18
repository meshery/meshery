package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/httputil"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
)

type capturedHTTPRequest struct {
	method string
	path   string
}

type remoteProviderMethodCase struct {
	name           string
	wantCode       string
	successStatus  int
	expectedMethod string
	expectedPath   string
	wantStaticObj  string
	call           func(rp *RemoteProvider) error
	successBody    func(t *testing.T) []byte
	verifySuccess  func(t *testing.T, rp *RemoteProvider)
}

func newHTTPTestRemoteProvider(t *testing.T, remoteURL string) *RemoteProvider {
	t.Helper()

	provider := newTestRemoteProvider(t, remoteURL)
	provider.Capabilities = Capabilities{
		{Feature: UsersIdentity, Endpoint: "/api/identity/users"},
		{Feature: PersistPerformanceProfiles, Endpoint: "/api/performance/profiles"},
		{Feature: PersistResults, Endpoint: "/api/performance/results"},
		{Feature: PersistSMIResults, Endpoint: "/api/smi/results"},
		{Feature: PersistResult, Endpoint: "/api/performance/result"},
		{Feature: PersistMetrics, Endpoint: "/api/performance/metrics"},
		{Feature: PersistMesheryPatternResources, Endpoint: "/api/content/patterns/resource"},
	}
	setViperKeyForRemoteAuthTest(t, "opt-token", "test-token")
	return provider
}

func newTokenCookieRequest(method, target string) *http.Request {
	req := httptest.NewRequest(method, target, nil)
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: "test-token"})
	return req
}

func mustMarshalTestJSON(t *testing.T, v any) []byte {
	t.Helper()

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal test fixture: %v", err)
	}
	return data
}

func callWithoutPanic(t *testing.T, fn func() error) error {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("method panicked unexpectedly: %v", r)
		}
	}()
	return fn()
}

func assertCapturedRequest(t *testing.T, reqCh <-chan capturedHTTPRequest, wantMethod, wantPath string) {
	t.Helper()

	select {
	case req := <-reqCh:
		if req.method != wantMethod || req.path != wantPath {
			t.Fatalf("request = %s %s, want %s %s", req.method, req.path, wantMethod, wantPath)
		}
	default:
		t.Fatalf("expected request to mock server (%s %s), but none was recorded", wantMethod, wantPath)
	}
}

func remoteProviderMethodCases() []remoteProviderMethodCase {
	resultID := core.Uuid(uuid.Must(uuid.FromString("11111111-1111-4111-8111-111111111111")))

	return []remoteProviderMethodCase{
		{
			name:           "GetUsers",
			wantCode:       ErrFetchCode,
			successStatus:  http.StatusOK,
			expectedMethod: http.MethodGet,
			expectedPath:   "/api/identity/users",
			wantStaticObj:  "Users Data",
			call: func(rp *RemoteProvider) error {
				_, err := rp.GetUsers("test-token", "0", "10", "", "", "")
				return err
			},
			successBody: func(_ *testing.T) []byte {
				return []byte(`{"users":[]}`)
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				body, err := rp.GetUsers("test-token", "0", "10", "", "", "")
				if err != nil {
					t.Fatalf("GetUsers returned unexpected error: %v", err)
				}
				if string(body) != `{"users":[]}` {
					t.Errorf("GetUsers body = %q, want %q", string(body), `{"users":[]}`)
				}
			},
		},
		{
			name:           "FetchResults",
			wantCode:       ErrFetchCode,
			successStatus:  http.StatusOK,
			expectedMethod: http.MethodGet,
			expectedPath:   "/api/performance/profiles/profile-1/results",
			call: func(rp *RemoteProvider) error {
				_, err := rp.FetchResults("test-token", "0", "10", "", "", "profile-1")
				return err
			},
			successBody: func(_ *testing.T) []byte {
				return []byte(`{"results":[]}`)
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				body, err := rp.FetchResults("test-token", "0", "10", "", "", "profile-1")
				if err != nil {
					t.Fatalf("FetchResults returned unexpected error: %v", err)
				}
				if string(body) != `{"results":[]}` {
					t.Errorf("FetchResults body = %q, want %q", string(body), `{"results":[]}`)
				}
			},
		},
		{
			name:           "FetchAllResults",
			wantCode:       ErrFetchCode,
			successStatus:  http.StatusOK,
			expectedMethod: http.MethodGet,
			expectedPath:   "/api/performance/results",
			call: func(rp *RemoteProvider) error {
				_, err := rp.FetchAllResults("test-token", "0", "10", "", "", "", "")
				return err
			},
			successBody: func(_ *testing.T) []byte {
				return []byte(`{"results":[]}`)
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				body, err := rp.FetchAllResults("test-token", "0", "10", "", "", "", "")
				if err != nil {
					t.Fatalf("FetchAllResults returned unexpected error: %v", err)
				}
				if string(body) != `{"results":[]}` {
					t.Errorf("FetchAllResults body = %q, want %q", string(body), `{"results":[]}`)
				}
			},
		},
		{
			name:           "FetchSmiResults",
			wantCode:       ErrFetchCode,
			successStatus:  http.StatusOK,
			expectedMethod: http.MethodGet,
			expectedPath:   "/api/smi/results",
			wantStaticObj:  "SMI Result",
			call: func(rp *RemoteProvider) error {
				_, err := rp.FetchSmiResults(newTokenCookieRequest(http.MethodGet, "/api/smi"), "0", "10", "", "")
				return err
			},
			successBody: func(_ *testing.T) []byte {
				return []byte(`{"results":[]}`)
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				body, err := rp.FetchSmiResults(newTokenCookieRequest(http.MethodGet, "/api/smi"), "0", "10", "", "")
				if err != nil {
					t.Fatalf("FetchSmiResults returned unexpected error: %v", err)
				}
				if string(body) != `{"results":[]}` {
					t.Errorf("FetchSmiResults body = %q, want %q", string(body), `{"results":[]}`)
				}
			},
		},
		{
			name:           "FetchSmiResult",
			wantCode:       ErrFetchCode,
			successStatus:  http.StatusOK,
			expectedMethod: http.MethodGet,
			expectedPath:   "/api/smi/results/" + resultID.String(),
			wantStaticObj:  "SMI Result",
			call: func(rp *RemoteProvider) error {
				_, err := rp.FetchSmiResult(newTokenCookieRequest(http.MethodGet, "/api/smi/"+resultID.String()), "0", "10", "", "", resultID)
				return err
			},
			successBody: func(_ *testing.T) []byte {
				return []byte(`{"id":"11111111-1111-4111-8111-111111111111"}`)
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				body, err := rp.FetchSmiResult(newTokenCookieRequest(http.MethodGet, "/api/smi/"+resultID.String()), "0", "10", "", "", resultID)
				if err != nil {
					t.Fatalf("FetchSmiResult returned unexpected error: %v", err)
				}
				if !strings.Contains(string(body), resultID.String()) {
					t.Errorf("FetchSmiResult body = %q, want result id %q", string(body), resultID.String())
				}
			},
		},
		{
			name:           "GetResult",
			wantCode:       ErrFetchCode,
			successStatus:  http.StatusOK,
			expectedMethod: http.MethodGet,
			expectedPath:   "/api/performance/result/" + resultID.String(),
			call: func(rp *RemoteProvider) error {
				_, err := rp.GetResult("test-token", resultID)
				return err
			},
			successBody: func(t *testing.T) []byte {
				return mustMarshalTestJSON(t, &MesheryResult{ID: resultID, Name: "perf-1"})
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				res, err := rp.GetResult("test-token", resultID)
				if err != nil {
					t.Fatalf("GetResult returned unexpected error: %v", err)
				}
				if res == nil || res.Name != "perf-1" {
					t.Errorf("GetResult = %+v, want Name %q", res, "perf-1")
				}
			},
		},
		{
			name:           "PublishResults",
			wantCode:       ErrPostCode,
			successStatus:  http.StatusCreated,
			expectedMethod: http.MethodPost,
			expectedPath:   "/api/performance/profiles/profile-1/results",
			call: func(rp *RemoteProvider) error {
				_, err := rp.PublishResults(newTokenCookieRequest(http.MethodPost, "/api/performance/results"), &MesheryResult{Name: "perf-1"}, "profile-1")
				return err
			},
			successBody: func(_ *testing.T) []byte {
				return []byte(`{"id":"published-result-id"}`)
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				id, err := rp.PublishResults(newTokenCookieRequest(http.MethodPost, "/api/performance/results"), &MesheryResult{Name: "perf-1"}, "profile-1")
				if err != nil {
					t.Fatalf("PublishResults returned unexpected error: %v", err)
				}
				if id != "published-result-id" {
					t.Errorf("PublishResults id = %q, want %q", id, "published-result-id")
				}
			},
		},
		{
			name:           "PublishSmiResults",
			wantCode:       ErrPostCode,
			successStatus:  http.StatusCreated,
			expectedMethod: http.MethodPost,
			expectedPath:   "/api/smi/results",
			call: func(rp *RemoteProvider) error {
				_, err := rp.PublishSmiResults(&SmiResult{})
				return err
			},
			successBody: func(_ *testing.T) []byte {
				return []byte(`{"id":"published-smi-id"}`)
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				id, err := rp.PublishSmiResults(&SmiResult{})
				if err != nil {
					t.Fatalf("PublishSmiResults returned unexpected error: %v", err)
				}
				if id != "published-smi-id" {
					t.Errorf("PublishSmiResults id = %q, want %q", id, "published-smi-id")
				}
			},
		},
		{
			name:           "PublishMetrics",
			wantCode:       ErrPostCode,
			successStatus:  http.StatusOK,
			expectedMethod: http.MethodPut,
			expectedPath:   "/api/performance/metrics",
			call: func(rp *RemoteProvider) error {
				return rp.PublishMetrics("test-token", &MesheryResult{Name: "perf-1"})
			},
			successBody: func(_ *testing.T) []byte {
				return []byte(`{}`)
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				if err := rp.PublishMetrics("test-token", &MesheryResult{Name: "perf-1"}); err != nil {
					t.Fatalf("PublishMetrics returned unexpected error: %v", err)
				}
			},
		},
		{
			name:           "SaveMesheryPatternResource",
			wantCode:       ErrPostCode,
			successStatus:  http.StatusOK,
			expectedMethod: http.MethodPost,
			expectedPath:   "/api/content/patterns/resource",
			call: func(rp *RemoteProvider) error {
				_, err := rp.SaveMesheryPatternResource("test-token", &PatternResource{Name: "saved-design"})
				return err
			},
			successBody: func(t *testing.T) []byte {
				return mustMarshalTestJSON(t, &PatternResource{Name: "saved-design"})
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				res, err := rp.SaveMesheryPatternResource("test-token", &PatternResource{Name: "saved-design"})
				if err != nil {
					t.Fatalf("SaveMesheryPatternResource returned unexpected error: %v", err)
				}
				if res == nil || res.Name != "saved-design" {
					t.Errorf("SaveMesheryPatternResource = %+v, want Name %q", res, "saved-design")
				}
			},
		},
		{
			name:           "GetMesheryPatternResource",
			wantCode:       ErrFetchCode,
			successStatus:  http.StatusOK,
			expectedMethod: http.MethodGet,
			expectedPath:   "/api/content/patterns/resource/res-1",
			call: func(rp *RemoteProvider) error {
				_, err := rp.GetMesheryPatternResource("test-token", "res-1")
				return err
			},
			successBody: func(t *testing.T) []byte {
				return mustMarshalTestJSON(t, &PatternResource{Name: "fetched-design"})
			},
			verifySuccess: func(t *testing.T, rp *RemoteProvider) {
				res, err := rp.GetMesheryPatternResource("test-token", "res-1")
				if err != nil {
					t.Fatalf("GetMesheryPatternResource returned unexpected error: %v", err)
				}
				if res == nil || res.Name != "fetched-design" {
					t.Errorf("GetMesheryPatternResource = %+v, want Name %q", res, "fetched-design")
				}
			},
		},
	}
}

// TestRemoteProviderHTTPErrors_AllAffectedMethods verifies all 11 RemoteProvider
// methods handle non-success HTTP status codes (including non-refreshable 401 and
// 4xx/5xx responses) without panicking and preserve MeshKit codes and status tags.
func TestRemoteProviderHTTPErrors_AllAffectedMethods(t *testing.T) {
	statuses := []int{
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	}
	const errBody = "upstream provider error body"

	for _, tc := range remoteProviderMethodCases() {
		for _, status := range statuses {
			subName := tc.name + "/status_" + strconv.Itoa(status)
			t.Run(subName, func(t *testing.T) {
				reqCh := make(chan capturedHTTPRequest, 2)
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					select {
					case reqCh <- capturedHTTPRequest{method: r.Method, path: r.URL.Path}:
					default:
					}
					w.WriteHeader(status)
					_, _ = w.Write([]byte(errBody))
				}))
				defer server.Close()

				rp := newHTTPTestRemoteProvider(t, server.URL)
				err := callWithoutPanic(t, func() error {
					return tc.call(rp)
				})
				if err == nil {
					t.Fatalf("%s with status %d returned nil error", tc.name, status)
				}

				assertCapturedRequest(t, reqCh, tc.expectedMethod, tc.expectedPath)

				if gotCode := meshkiterrors.GetCode(err); gotCode != tc.wantCode {
					t.Errorf("error code = %q, want %q (err: %v)", gotCode, tc.wantCode, err)
				}
				if gotStatus, ok := httputil.ProviderStatusCode(err); !ok || gotStatus != status {
					t.Errorf("ProviderStatusCode = (%d, %v), want (%d, true)", gotStatus, ok, status)
				}

				sdesc := meshkiterrors.GetSDescription(err)
				if tc.wantStaticObj != "" {
					if !strings.Contains(sdesc, tc.wantStaticObj) {
						t.Errorf("short description %q does not contain %q", sdesc, tc.wantStaticObj)
					}
				} else if !strings.Contains(sdesc, errBody) {
					t.Errorf("short description %q does not contain response body %q", sdesc, errBody)
				}

				mkErr := assertMeshKitErrorStruct(t, err)
				if len(mkErr.LongDescription) < 2 {
					t.Fatalf("LongDescription = %#v, want status prefix and cause element", mkErr.LongDescription)
				}
				causeDesc := mkErr.LongDescription[1]
				if !strings.Contains(causeDesc, "unexpected HTTP status") || !strings.Contains(causeDesc, strconv.Itoa(status)) {
					t.Errorf("error cause %q should contain %q and status %d", causeDesc, "unexpected HTTP status", status)
				}
			})
		}
	}
}

// TestRemoteProviderHTTPErrors_201MethodsReject200OK ensures PublishResults and
// PublishSmiResults treat a 200 OK response (when 201 Created is required) as an
// ErrPost failure without panicking or attaching an out-of-range status tag.
func TestRemoteProviderHTTPErrors_201MethodsReject200OK(t *testing.T) {
	cases := []struct {
		name string
		call func(rp *RemoteProvider) (string, error)
	}{
		{
			name: "PublishResults",
			call: func(rp *RemoteProvider) (string, error) {
				return rp.PublishResults(newTokenCookieRequest(http.MethodPost, "/api/performance/results"), &MesheryResult{Name: "perf-1"}, "profile-1")
			},
		},
		{
			name: "PublishSmiResults",
			call: func(rp *RemoteProvider) (string, error) {
				return rp.PublishSmiResults(&SmiResult{})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"id":"ignored"}`))
			}))
			defer server.Close()

			rp := newHTTPTestRemoteProvider(t, server.URL)
			var id string
			err := callWithoutPanic(t, func() error {
				var callErr error
				id, callErr = tc.call(rp)
				return callErr
			})
			if err == nil {
				t.Fatalf("%s unexpectedly succeeded on HTTP 200 (expected 201)", tc.name)
			}
			if id != "" {
				t.Errorf("%s id = %q, want empty string on error", tc.name, id)
			}
			if gotCode := meshkiterrors.GetCode(err); gotCode != ErrPostCode {
				t.Errorf("error code = %q, want %q", gotCode, ErrPostCode)
			}
			if gotStatus, ok := httputil.ProviderStatusCode(err); ok {
				t.Errorf("ProviderStatusCode = (%d, true), want (_, false) for HTTP 200", gotStatus)
			}
		})
	}
}

// TestSaveMesheryPatternResource_ErrorPaths verifies SaveMesheryPatternResource's
// fixed control flow across unexpected HTTP statuses (body preservation) and
// unreachable provider failures (resp == nil).
func TestSaveMesheryPatternResource_ErrorPaths(t *testing.T) {
	t.Run("non-200 status preserves response body instead of struct formatting", func(t *testing.T) {
		const errBody = "validation failed: invalid design payload"
		for _, tc := range []struct {
			name       string
			status     int
			wantTagged bool
		}{
			{name: "422 Unprocessable Entity", status: http.StatusUnprocessableEntity, wantTagged: true},
			{name: "201 Created", status: http.StatusCreated, wantTagged: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(tc.status)
					_, _ = w.Write([]byte(errBody))
				}))
				defer server.Close()

				rp := newHTTPTestRemoteProvider(t, server.URL)
				var res *PatternResource
				err := callWithoutPanic(t, func() error {
					var callErr error
					res, callErr = rp.SaveMesheryPatternResource("test-token", &PatternResource{Name: "design"})
					return callErr
				})
				if err == nil {
					t.Fatalf("expected error on HTTP %d, got nil", tc.status)
				}
				if res != nil {
					t.Fatalf("expected nil PatternResource on error, got %+v", res)
				}
				if gotCode := meshkiterrors.GetCode(err); gotCode != ErrPostCode {
					t.Errorf("error code = %q, want %q", gotCode, ErrPostCode)
				}

				gotStatus, ok := httputil.ProviderStatusCode(err)
				if tc.wantTagged && (!ok || gotStatus != tc.status) {
					t.Errorf("ProviderStatusCode = (%d, %v), want (%d, true)", gotStatus, ok, tc.status)
				}
				if !tc.wantTagged && ok {
					t.Errorf("ProviderStatusCode = (%d, true), want (_, false)", gotStatus)
				}

				sdesc := meshkiterrors.GetSDescription(err)
				if !strings.Contains(sdesc, errBody) {
					t.Errorf("short description %q does not contain body %q", sdesc, errBody)
				}
				if strings.Contains(sdesc, "&{") {
					t.Errorf("short description %q still formats raw http.Response Body struct", sdesc)
				}
			})
		}
	})

	t.Run("unreachable provider with nil response returns ErrUnreachableRemoteProvider", func(t *testing.T) {
		closedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		closedURL := closedServer.URL
		closedServer.Close()

		rp := newHTTPTestRemoteProvider(t, closedURL)
		var res *PatternResource
		err := callWithoutPanic(t, func() error {
			var callErr error
			res, callErr = rp.SaveMesheryPatternResource("test-token", &PatternResource{Name: "design"})
			return callErr
		})
		if err == nil {
			t.Fatal("expected error for unreachable provider, got nil")
		}
		if res != nil {
			t.Fatalf("expected nil PatternResource on error, got %+v", res)
		}
		if gotCode := meshkiterrors.GetCode(err); gotCode != ErrUnreachableRemoteProviderCode {
			t.Errorf("error code = %q, want %q", gotCode, ErrUnreachableRemoteProviderCode)
		}
		if gotStatus, ok := httputil.ProviderStatusCode(err); ok {
			t.Errorf("ProviderStatusCode = (%d, true), want (_, false)", gotStatus)
		}
	})

	t.Run("body read failure returns ErrDataRead", func(t *testing.T) {
		hijackErrCh := make(chan error, 1)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				_, _ = io.Copy(io.Discard, r.Body)
				_ = r.Body.Close()
			}
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
			// Signal setup success before closing the socket so SaveMesheryPatternResource
			// cannot observe EOF and return before hijackErrCh is populated.
			hijackErrCh <- nil
			_ = conn.Close()
		}))
		defer server.Close()

		rp := newHTTPTestRemoteProvider(t, server.URL)
		var res *PatternResource
		err := callWithoutPanic(t, func() error {
			var callErr error
			res, callErr = rp.SaveMesheryPatternResource("test-token", &PatternResource{Name: "design"})
			return callErr
		})

		select {
		case herr := <-hijackErrCh:
			if herr != nil {
				t.Fatalf("test server setup failed: %v", herr)
			}
		default:
			t.Fatal("expected hijacked test server handler to run")
		}

		if err == nil {
			t.Fatal("expected error on truncated response body, got nil")
		}
		if res != nil {
			t.Fatalf("expected nil PatternResource on body read failure, got %+v", res)
		}
		if gotCode := meshkiterrors.GetCode(err); gotCode != ErrDataReadCode {
			t.Errorf("error code = %q, want %q (err: %v)", gotCode, ErrDataReadCode, err)
		}
		if gotStatus, ok := httputil.ProviderStatusCode(err); ok {
			t.Errorf("ProviderStatusCode = (%d, true), want (_, false) on body read failure", gotStatus)
		}
		if sdesc := meshkiterrors.GetSDescription(err); !strings.Contains(sdesc, "Design resource") {
			t.Errorf("short description %q does not contain %q", sdesc, "Design resource")
		}

		mkErr := assertMeshKitErrorStruct(t, err)
		longDesc := strings.Join(mkErr.LongDescription, " ")
		hasReadFailureCause := strings.Contains(longDesc, "unexpected EOF") ||
			strings.Contains(longDesc, "connection reset by peer") ||
			strings.Contains(longDesc, "forcibly closed")
		if !hasReadFailureCause {
			t.Errorf("ErrDataRead LongDescription %q does not capture premature EOF or connection reset cause", longDesc)
		}
	})
}

// TestRemoteProviderHTTPErrors_SuccessPaths ensures none of the HTTP error-handling
// fixes alter existing 200 OK / 201 Created success paths across all 11 methods.
func TestRemoteProviderHTTPErrors_SuccessPaths(t *testing.T) {
	for _, tc := range remoteProviderMethodCases() {
		t.Run(tc.name, func(t *testing.T) {
			reqCh := make(chan capturedHTTPRequest, 2)
			payload := tc.successBody(t)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				select {
				case reqCh <- capturedHTTPRequest{method: r.Method, path: r.URL.Path}:
				default:
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.successStatus)
				_, _ = w.Write(payload)
			}))
			defer server.Close()

			rp := newHTTPTestRemoteProvider(t, server.URL)
			tc.verifySuccess(t, rp)
			assertCapturedRequest(t, reqCh, tc.expectedMethod, tc.expectedPath)
		})
	}
}
