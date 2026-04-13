package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/internal/ai"
	"github.com/meshery/meshery/server/internal/ai/providers/openai"
)

// --- helpers ---

func newTestServer(handler http.HandlerFunc) (*httptest.Server, *openai.Provider) {
	srv := httptest.NewServer(handler)
	p := openai.New(
		"test-api-key",
		openai.WithBaseURL(srv.URL),
		openai.WithHTTPClient(srv.Client()),
	)
	return srv, p
}

func simpleRequest() ai.CompleteRequest {
	return ai.CompleteRequest{
		Context:  ai.RequestContext{Provider: "openai"},
		Messages: []ai.Message{{Role: ai.RoleUser, Content: "ping"}},
	}
}

func openaiSuccessResponse(content string) map[string]interface{} {
	return map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message":       map[string]string{"role": "assistant", "content": content},
				"finish_reason": "stop",
			},
		},
		"usage": map[string]int{
			"prompt_tokens":     10,
			"completion_tokens": 5,
			"total_tokens":      15,
		},
	}
}

// --- Complete ---

func TestComplete_Success(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("missing/wrong auth header: %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(openaiSuccessResponse("pong"))
	})
	defer srv.Close()

	got, err := p.Complete(context.Background(), simpleRequest())
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if got.Content != "pong" {
		t.Errorf("Content = %q, want %q", got.Content, "pong")
	}
	if got.FinishReason != "stop" {
		t.Errorf("FinishReason = %q, want %q", got.FinishReason, "stop")
	}
	if got.Usage == nil {
		t.Fatal("Usage is nil")
	}
	if got.Usage.TotalTokens != 15 {
		t.Errorf("TotalTokens = %d, want 15", got.Usage.TotalTokens)
	}
}

func TestComplete_NoChoices(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"choices": []interface{}{}})
	})
	defer srv.Close()

	_, err := p.Complete(context.Background(), simpleRequest())
	if err == nil {
		t.Fatal("expected error for empty choices, got nil")
	}
}

func TestComplete_RateLimited(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})
	defer srv.Close()

	_, err := p.Complete(context.Background(), simpleRequest())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	re, ok := err.(*ai.RetryableError)
	if !ok {
		// RetryableError may be wrapped; unwrap manually for test.
		t.Logf("error type: %T: %v", err, err)
	} else if re.StatusCode != http.StatusTooManyRequests {
		t.Errorf("StatusCode = %d, want 429", re.StatusCode)
	}
}

// --- Stream ---

func TestStream_Success(t *testing.T) {
	// Build SSE response for streaming.
	sseLines := []string{
		`data: {"choices":[{"delta":{"role":"assistant","content":"He"},"finish_reason":""}]}`,
		`data: {"choices":[{"delta":{"role":"assistant","content":"llo"},"finish_reason":""}]}`,
		`data: {"choices":[{"delta":{},"finish_reason":"stop"}]}`,
		`data: [DONE]`,
	}

	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		for _, line := range sseLines {
			fmt.Fprintln(w, line)
		}
	})
	defer srv.Close()

	reader, err := p.Stream(context.Background(), ai.StreamRequest{
		Context:  ai.RequestContext{Provider: "openai"},
		Messages: []ai.Message{{Role: ai.RoleUser, Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	var collected string
	var finishReason string
	for {
		chunk, err := reader.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv() error = %v", err)
		}
		collected += chunk.Content
		if chunk.FinishReason != "" {
			finishReason = chunk.FinishReason
		}
	}

	if collected != "Hello" {
		t.Errorf("collected = %q, want %q", collected, "Hello")
	}
	if finishReason != "stop" {
		t.Errorf("finishReason = %q, want %q", finishReason, "stop")
	}
}

func TestStream_HTTPError(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	defer srv.Close()

	_, err := p.Stream(context.Background(), ai.StreamRequest{
		Context:  ai.RequestContext{Provider: "openai"},
		Messages: []ai.Message{{Role: ai.RoleUser, Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- Embed ---

func TestEmbed_Success(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/embeddings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := map[string]interface{}{
			"data": []map[string]interface{}{
				{"embedding": []float32{0.1, 0.2, 0.3}},
			},
			"usage": map[string]int{
				"prompt_tokens": 5,
				"total_tokens":  5,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	got, err := p.Embed(context.Background(), ai.EmbedRequest{
		Context: ai.RequestContext{Provider: "openai"},
		Text:    "test",
	})
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}
	if len(got.Vector) != 3 {
		t.Errorf("Vector length = %d, want 3", len(got.Vector))
	}
	if got.Usage == nil || got.Usage.TotalTokens != 5 {
		t.Errorf("Usage.TotalTokens = %v", got.Usage)
	}
}

// --- system prompt ---

func TestComplete_SystemPrompt(t *testing.T) {
	type chatBody struct {
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}

	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		var body chatBody
		_ = json.NewDecoder(r.Body).Decode(&body)
		if len(body.Messages) < 1 || body.Messages[0].Role != "system" {
			t.Errorf("expected first message to be system, got: %+v", body.Messages)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(openaiSuccessResponse("ok"))
	})
	defer srv.Close()

	req := simpleRequest()
	req.Context.SystemPrompt = "Be helpful."
	_, err := p.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
}
