package anthropicclaude

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/schemas/models/core"
)


// customTransport redirects all requests to a target URL while preserving original scheme and host conceptually
type customTransport struct {
	Target   *url.URL
	Original http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	newReq := req.Clone(req.Context())
	newReq.URL.Scheme = t.Target.Scheme
	newReq.URL.Host = t.Target.Host
	return t.Original.RoundTrip(newReq)
}

func TestRegisterAction_Execute(t *testing.T) {
	userID, _ := uuid.NewV4()
	sysID := core.Uuid(uuid.Nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: core.Uuid(userID)})
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)

	t.Run("mock server returns 200 -> transitions to connected (Exit)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		targetURL, _ := url.Parse(server.URL)
		origTransport := http.DefaultTransport
		http.DefaultTransport = &customTransport{Target: targetURL, Original: origTransport}
		defer func() { http.DefaultTransport = origTransport }()

		action := &RegisterAction{}
		payload := connections.ConnectionPayload{
			MetaData: map[string]interface{}{
				"baseUrl":      "https://api.anthropic.com",
				"defaultModel": "claude-3-opus",
			},
			CredentialSecret: map[string]interface{}{
				"apiKey": "test-key",
			},
		}

		nextEvent, _, err := action.Execute(ctx, nil, payload)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if nextEvent != machines.Exit {
			t.Fatalf("expected machines.Exit, got %v", nextEvent)
		}
	})

	t.Run("mock server returns 401 -> transitions to failed (NoOp), error returned without response body", func(t *testing.T) {
		secretBody := "secret response body"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(secretBody))
		}))
		defer server.Close()

		targetURL, _ := url.Parse(server.URL)
		origTransport := http.DefaultTransport
		http.DefaultTransport = &customTransport{Target: targetURL, Original: origTransport}
		defer func() { http.DefaultTransport = origTransport }()

		action := &RegisterAction{}
		payload := connections.ConnectionPayload{
			MetaData: map[string]interface{}{
				"baseUrl":      "https://api.anthropic.com",
				"defaultModel": "claude-3-opus",
			},
			CredentialSecret: map[string]interface{}{
				"apiKey": "test-key",
			},
		}

		nextEvent, _, err := action.Execute(ctx, nil, payload)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if nextEvent != machines.NoOp {
			t.Fatalf("expected machines.NoOp, got %v", nextEvent)
		}

		errMsg := err.Error()
		if bytes.Contains([]byte(errMsg), []byte(secretBody)) {
			t.Fatalf("error message should not contain the response body. got: %s", errMsg)
		}
		if !bytes.Contains([]byte(errMsg), []byte("401")) {
			t.Fatalf("error message should contain status 401. got: %s", errMsg)
		}
	})

	t.Run("baseUrl pointing anywhere other than https://api.anthropic.com is rejected", func(t *testing.T) {
		hit := false
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hit = true
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		targetURL, _ := url.Parse(server.URL)
		origTransport := http.DefaultTransport
		http.DefaultTransport = &customTransport{Target: targetURL, Original: origTransport}
		defer func() { http.DefaultTransport = origTransport }()

		action := &RegisterAction{}
		payload := connections.ConnectionPayload{
			MetaData: map[string]interface{}{
				"baseUrl":      "https://other-api.anthropic.com",
				"defaultModel": "claude-3-opus",
			},
			CredentialSecret: map[string]interface{}{
				"apiKey": "test-key",
			},
		}

		nextEvent, _, err := action.Execute(ctx, nil, payload)
		if err == nil {
			t.Fatalf("expected error for invalid base URL, got nil")
		}
		if nextEvent != machines.NoOp {
			t.Fatalf("expected machines.NoOp, got %v", nextEvent)
		}
		if hit {
			t.Fatalf("expected mock server to not be hit")
		}
	})
}
