package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleWare(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	t.Log("Need to run AuthMiddleware() skipping")
	//_, err := AuthMiddleware(nil, nil)
	//if err != nil {
	//	t.Errorf("AuthMiddleWare() failed with error: %s", err)
	//}
}

// TestWriteMeshkitError_ErrTransientProvider is a focused test for the
// transient-provider response shape emitted by SessionInjectorMiddleware when
// Meshery Cloud is temporarily unreachable. A full middleware-level test would
// require stubbing the 100+ method models.Provider interface, which the
// package doesn't yet have. Instead, we test the migration end-to-end at the
// response-helper level: the same two lines the middleware now executes.
//
// If a provider mocking pattern lands in this package later, promote this to
// a middleware-dispatch test by wiring a stub provider whose GetUserDetails
// returns "Could not reach remote provider".
func TestWriteMeshkitError_ErrTransientProvider(t *testing.T) {
	rec := httptest.NewRecorder()

	writeMeshkitError(
		rec,
		ErrTransientProvider(errors.New("Could not reach remote provider: dial tcp: i/o timeout")),
		http.StatusServiceUnavailable,
	)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type application/json; charset=utf-8, got %q", ct)
	}
	if nosniff := resp.Header.Get("X-Content-Type-Options"); nosniff != "nosniff" {
		t.Errorf("expected X-Content-Type-Options: nosniff, got %q", nosniff)
	}

	var decoded struct {
		Error                string   `json:"error"`
		Code                 string   `json:"code"`
		Severity             string   `json:"severity"`
		ProbableCause        []string `json:"probableCause"`
		SuggestedRemediation []string `json:"suggestedRemediation"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("body did not parse as JSON: %v", err)
	}
	if decoded.Code != ErrTransientProviderCode {
		t.Errorf("expected code %q, got %q", ErrTransientProviderCode, decoded.Code)
	}
	if decoded.Error == "" {
		t.Errorf("expected non-empty error message")
	}
	if len(decoded.SuggestedRemediation) == 0 {
		t.Errorf("expected suggestedRemediation to be populated for transient-provider errors")
	}
}

// TestIsTransientProviderError exercises the detector used by
// SessionInjectorMiddleware to decide between a 401 auth failure and a 503
// transient-provider failure. Regression guard in case the error-message
// substrings the detector relies on drift.
func TestIsTransientProviderError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil is not transient", nil, false},
		{"generic auth failure is not transient", fmt.Errorf("token expired"), false},
		{"Could not reach remote provider is transient", fmt.Errorf("Could not reach remote provider: dial tcp: i/o timeout"), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isTransientProviderError(tc.err); got != tc.want {
				t.Errorf("isTransientProviderError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
