package system

import (
	"fmt"
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

func TestDashboardCmdArgs(t *testing.T) {
	expectedErr := utils.ErrInvalidArgument(
		fmt.Errorf("unexpected argument(s)\n\nUsage: mesheryctl system dashboard [flags]\nRun 'mesheryctl system dashboard --help' to see detailed help message"),
	)

	tests := []struct {
		name          string
		args          []string
		wantErr       bool
		expectedError error
	}{
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: false,
		},
		{
			name:          "single unexpected argument",
			args:          []string{"foo"},
			wantErr:       true,
			expectedError: expectedErr,
		},
		{
			name:          "multiple unexpected arguments",
			args:          []string{"foo", "bar"},
			wantErr:       true,
			expectedError: expectedErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dashboardCmd.Args(dashboardCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("dashboardCmd.Args() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil {
				utils.AssertMeshkitErrorsEqual(t, err, tt.expectedError)
			}
		})
	}
}
