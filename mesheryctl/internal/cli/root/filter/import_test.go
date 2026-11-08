package filter

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

type trackingBody struct {
	io.ReadCloser
	closed bool
}

func (b *trackingBody) Close() error {
	b.closed = true
	return b.ReadCloser.Close()
}

type trackingTransport struct {
	base http.RoundTripper
	body *trackingBody
}

func (t *trackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	resp, err := base.RoundTrip(req)
	if err == nil && resp != nil && resp.Body != nil {
		t.body = &trackingBody{ReadCloser: resp.Body}
		resp.Body = t.body
	}
	return resp, err
}

func setupImportTest(t *testing.T, endpoint string) *trackingTransport {
	t.Helper()
	utils.SetupContextEnv(t)
	viper.Set("contexts.local.endpoint", endpoint)
	origToken := utils.TokenFlag
	t.Cleanup(func() { utils.TokenFlag = origToken })
	utils.TokenFlag = utils.GetToken(t)
	if tokenData, err := os.ReadFile(utils.GetToken(t)); err == nil {
		tempAuth := filepath.Join(t.TempDir(), "auth.json")
		if err := os.WriteFile(tempAuth, tokenData, 0o600); err == nil {
			viper.Set("tokens", []map[string]string{{"name": "default", "location": tempAuth}})
		}
	}
	transport := &trackingTransport{base: http.DefaultTransport}
	previous := http.DefaultTransport
	http.DefaultTransport = transport
	t.Cleanup(func() { http.DefaultTransport = previous })
	return transport
}

func TestImportCmd_BodyClosed(t *testing.T) {
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	t.Run("closes body on 200 OK response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{"status":"imported"}`)
		}))
		defer server.Close()

		transport := setupImportTest(t, server.URL)

		err := importCmd.RunE(importCmd, []string{"http://example.com/filter.wasm"})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if transport.body == nil || !transport.body.closed {
			t.Fatal("expected response body to be closed on 200 OK")
		}
	})

	t.Run("closes body on unexpected non-200 success status (201)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprint(w, `{"error":"bad request"}`)
		}))
		defer server.Close()

		transport := setupImportTest(t, server.URL)

		err := importCmd.RunE(importCmd, []string{"http://example.com/filter.wasm"})
		if err == nil {
			t.Fatal("expected error on 201 response, got nil")
		}
		if transport.body == nil || !transport.body.closed {
			t.Fatal("expected response body to be closed on non-200 response")
		}
	})
}
