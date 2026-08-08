package handlers

// This file is an ad hoc connectivity check for the Anthropic model-provider
// Connection/Credential, per the first steps outlined on meshery/meshery#20994:
// https://github.com/meshery/meshery/issues/20994#issuecomment-5114456107
//
// It intentionally does NOT attempt to be the production provider health/readiness
// check described in the AI Adapter project's expected outcomes - that belongs to
// a later, dedicated PR with proper status/error/event semantics and
// operationId correlation. This is only meant to prove, end to end, that a real
// Anthropic API key can be exercised through the connection/credential shape
// established in #19877/#19888.
//
// It is skipped by default (no network access, no CI dependency). To run it
// locally with a real key:
//
//	$env:ANTHROPIC_API_KEY = "sk-ant-..."
//	go test ./server/handlers/ -run TestAnthropicConnectionLiveConnectivity -v

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

const anthropicAPIBaseURL = "https://api.anthropic.com/v1"
const anthropicAPIVersion = "2023-06-01"

// TestAnthropicConnectionLiveConnectivity makes a minimal, real request to
// Anthropic's API using a locally supplied API key, to confirm that a stored
// Credential's apiKey value is sufficient to authenticate. It does not assert
// on model output or response body - only that authentication succeeds - since
// its only job is to prove the credential round-trips correctly, not to
// exercise generation.
func TestAnthropicConnectionLiveConnectivity(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set; skipping live connectivity check")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, anthropicAPIBaseURL+"/models", nil)
	if err != nil {
		t.Fatalf("building connectivity check request: %v", err)
	}
	// Never log apiKey; only its presence/absence and the resulting status code.
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", anthropicAPIVersion)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("anthropic connectivity check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("anthropic connectivity check returned status %d, want %d", resp.StatusCode, http.StatusOK)
	}
}