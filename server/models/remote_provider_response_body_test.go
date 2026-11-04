package models

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/schemas/models/core"
)

// responseBodyTestHost is the remote provider host in these tests. Only
// requests to it reach the tracking transport; it is never dialed.
const responseBodyTestHost = "remote-provider.invalid"

// trackedBody is a response body that records whether it was closed.
type trackedBody struct {
	io.Reader
	closed atomic.Bool
}

func (b *trackedBody) Close() error {
	b.closed.Store(true)
	return nil
}

// bodyTrackingTransport answers requests for responseBodyTestHost with a fixed
// status and a trackedBody, and passes any other request on to next.
type bodyTrackingTransport struct {
	status int
	next   http.RoundTripper

	mu     sync.Mutex
	bodies []*trackedBody
}

func (rt *bodyTrackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host != responseBodyTestHost {
		return rt.next.RoundTrip(req)
	}
	if req.Body != nil {
		_ = req.Body.Close()
	}
	body := &trackedBody{Reader: strings.NewReader("{}")}
	rt.mu.Lock()
	rt.bodies = append(rt.bodies, body)
	rt.mu.Unlock()
	return &http.Response{
		StatusCode: rt.status,
		Header:     make(http.Header),
		Body:       body,
		Request:    req,
	}, nil
}

// expectResponseBodiesClosed runs call once for each response status the
// remote provider might send, with http.DefaultTransport (which the remote
// provider's HTTP client uses) replaced by a bodyTrackingTransport. It fails if
// a response body is left open. An unclosed body keeps its connection out of
// the idle pool, so every call that leaks one leaks a connection.
func expectResponseBodiesClosed(t *testing.T, call func()) {
	t.Helper()
	for _, status := range []int{http.StatusOK, http.StatusCreated, http.StatusInternalServerError} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			rt := &bodyTrackingTransport{status: status, next: http.DefaultTransport}
			http.DefaultTransport = rt
			t.Cleanup(func() { http.DefaultTransport = rt.next })

			call()

			rt.mu.Lock()
			defer rt.mu.Unlock()
			if len(rt.bodies) == 0 {
				t.Fatal("no request reached the remote provider")
			}
			for i, body := range rt.bodies {
				if !body.closed.Load() {
					t.Errorf("response body %d of %d was not closed", i+1, len(rt.bodies))
				}
			}
		})
	}
}

func newResponseBodyTestProvider(t *testing.T, feature Feature, endpoint string) *RemoteProvider {
	t.Helper()
	provider := newTestRemoteProvider(t, "http://"+responseBodyTestHost)
	provider.Capabilities = Capabilities{{Feature: feature, Endpoint: endpoint}}
	return provider
}

// newShareRequest returns an inbound share request that carries a token cookie.
func newShareRequest() *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/share", strings.NewReader(`{}`))
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: "test-token"})
	return req
}

func TestExecutePrefSync_ClosesResponseBody(t *testing.T) {
	provider := newResponseBodyTestProvider(t, SyncPrefs, "/user/preferences")

	expectResponseBodiesClosed(t, func() {
		provider.executePrefSync("test-token", &Preference{})
	})
}

func TestUpdateEventStatus_ClosesResponseBody(t *testing.T) {
	provider := newResponseBodyTestProvider(t, PersistEvents, "/api/events")

	expectResponseBodiesClosed(t, func() {
		_ = provider.UpdateEventStatus("test-token", core.Uuid{}, "read")
	})
}

func TestBulkUpdateEventStatus_ClosesResponseBody(t *testing.T) {
	provider := newResponseBodyTestProvider(t, PersistEvents, "/api/events")

	expectResponseBodiesClosed(t, func() {
		_ = provider.BulkUpdateEventStatus("test-token", []*core.Uuid{{}}, "read")
	})
}

func TestDeleteEvent_ClosesResponseBody(t *testing.T) {
	provider := newResponseBodyTestProvider(t, PersistEvents, "/api/events")

	expectResponseBodiesClosed(t, func() {
		_ = provider.DeleteEvent("test-token", core.Uuid{})
	})
}

func TestBulkDeleteEvent_ClosesResponseBody(t *testing.T) {
	provider := newResponseBodyTestProvider(t, PersistEvents, "/api/events")

	expectResponseBodiesClosed(t, func() {
		_ = provider.BulkDeleteEvent("test-token", []*core.Uuid{{}})
	})
}

func TestUpdateConnectionStatusByID_ClosesResponseBody(t *testing.T) {
	provider := newResponseBodyTestProvider(t, PersistConnection, "/api/integrations/connections")

	expectResponseBodiesClosed(t, func() {
		_, _, _ = provider.UpdateConnectionStatusByID("test-token", core.Uuid{}, connections.CONNECTED)
	})
}

func TestShareDesign_ClosesResponseBody(t *testing.T) {
	provider := newResponseBodyTestProvider(t, ShareDesigns, "/api/content/design/share")

	expectResponseBodiesClosed(t, func() {
		_, _ = provider.ShareDesign(newShareRequest())
	})
}

func TestShareFilter_ClosesResponseBody(t *testing.T) {
	provider := newResponseBodyTestProvider(t, ShareFilters, "/api/content/filter/share")

	expectResponseBodiesClosed(t, func() {
		_, _ = provider.ShareFilter(newShareRequest())
	})
}
