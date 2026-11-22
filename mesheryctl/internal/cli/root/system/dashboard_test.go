package system

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// trackingBody wraps an io.ReadCloser to track when Close is called.
type trackingBody struct {
	io.ReadCloser
	closed *bool
}

// Close records that the body was closed and closes the underlying ReadCloser.
func (b *trackingBody) Close() error {
	*b.closed = true
	return b.ReadCloser.Close()
}

// trackingTransport wraps an http.RoundTripper to intercept responses and wrap their bodies with trackingBody.
type trackingTransport struct {
	base   http.RoundTripper
	closed *bool
}

// RoundTrip executes the HTTP request and wraps the response body with trackingBody.
func (t trackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	resp.Body = &trackingBody{ReadCloser: resp.Body, closed: t.closed}
	return resp, nil
}

func TestKeepConnectionAlive_Success(t *testing.T) {
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	bodyClosed := false
	originalTransport := http.DefaultTransport
	http.DefaultTransport = trackingTransport{
		base:   originalTransport,
		closed: &bodyClosed,
	}
	defer func() {
		http.DefaultTransport = originalTransport
	}()

	var requestReceived atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived.Store(true)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	keepConnectionAlive(server.URL)

	if !requestReceived.Load() {
		t.Fatalf("expected request to be received by server, but it was not")
	}
	if !bodyClosed {
		t.Fatal("expected response body to be closed")
	}
}

func TestKeepConnectionAlive_Failure(t *testing.T) {
	buf := utils.SetupMeshkitLoggerTesting(t, true)

	// Point to a closed server URL that immediately refuses connection
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close() // close immediately so connection fails

	// keepConnectionAlive should log error and return without panicking or succeeding
	keepConnectionAlive(url)

	output := buf.String()
	if !strings.Contains(output, "connection request failed") {
		t.Fatal("expected connection failure to be logged")
	}
	if strings.Contains(output, "connection request success") {
		t.Fatal("did not expect success to be logged after connection failure")
	}
}

// TestOpenMesheryUI verifies that --skip-browser keeps the browser closed, and
// that the user is always pointed to the given URL.
func TestOpenMesheryUI(t *testing.T) {
	const mesheryURL = "http://localhost:9081"

	tests := []struct {
		name        string
		skipBrowser bool
		browserErr  error
		wantOpened  bool
		wantLog     string
	}{
		{name: "opens the browser", wantOpened: true, wantLog: "Opening Meshery UI in browser at " + mesheryURL},
		{name: "skip-browser leaves it closed", skipBrowser: true, wantLog: "Meshery UI available at: " + mesheryURL},
		{name: "browser fails to open", browserErr: errors.New("no browser"), wantOpened: true, wantLog: "please point your browser to " + mesheryURL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, true)

			originalNavigate, originalSkip := navigateToBrowser, dashboardCmdFlags.SkipBrowser
			t.Cleanup(func() { navigateToBrowser, dashboardCmdFlags.SkipBrowser = originalNavigate, originalSkip })

			var opened []string
			navigateToBrowser = func(url string) error {
				opened = append(opened, url)
				return tt.browserErr
			}
			dashboardCmdFlags.SkipBrowser = tt.skipBrowser

			openMesheryUI(mesheryURL, "")

			if tt.wantOpened && (len(opened) != 1 || opened[0] != mesheryURL) {
				t.Errorf("expected the browser to open %s, got %v", mesheryURL, opened)
			}
			if !tt.wantOpened && len(opened) != 0 {
				t.Errorf("expected the browser to stay closed, got %v", opened)
			}
			if !strings.Contains(buf.String(), tt.wantLog) {
				t.Errorf("expected output to contain %q, got %q", tt.wantLog, buf.String())
			}
		})
	}
}
