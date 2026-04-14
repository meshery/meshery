package anthropic_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/internal/ai"
	"github.com/meshery/meshery/server/internal/ai/providers/anthropic"
)

// --- helpers ---

func newTestServer(handler http.HandlerFunc) (*httptest.Server, *anthropic.Provider) {
	srv := httptest.NewServer(handler)
	p := anthropic.New(
		"test-api-key",
		anthropic.WithBaseURL(srv.URL),
		anthropic.WithHTTPClient(srv.Client()),
	)
	return srv, p
}

func simpleRequest() ai.CompleteRequest {
	return ai.CompleteRequest{
		Context:  ai.RequestContext{Provider: "anthropic"},
		Messages: []ai.Message{{Role: ai.RoleUser, Content: "ping"}},
	}
}

func anthropicSuccessResponse(text string) map[string]interface{} {
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": text},
		},
		"stop_reason": "end_turn",
		"usage": map[string]int{
			"input_tokens":  10,
			"output_tokens": 5,
		},
	}
}

// --- Complete ---

func TestComplete_Success(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "test-api-key" {
			t.Errorf("missing/wrong api key header: %s", r.Header.Get("x-api-key"))
		}
		if r.Header.Get("anthropic-version") == "" {
			t.Error("missing anthropic-version header")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(anthropicSuccessResponse("pong"))
	})
	defer srv.Close()

	got, err := p.Complete(context.Background(), simpleRequest())
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if got.Content != "pong" {
		t.Errorf("Content = %q, want %q", got.Content, "pong")
	}
	if got.FinishReason != "end_turn" {
		t.Errorf("FinishReason = %q, want %q", got.FinishReason, "end_turn")
	}
	if got.Usage == nil {
		t.Fatal("Usage is nil")
	}
	if got.Usage.TotalTokens != 15 {
		t.Errorf("TotalTokens = %d, want 15", got.Usage.TotalTokens)
	}
}

func TestComplete_NoContent(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"content":     []interface{}{},
			"stop_reason": "end_turn",
			"usage":       map[string]int{"input_tokens": 0, "output_tokens": 0},
		})
	})
	defer srv.Close()

	_, err := p.Complete(context.Background(), simpleRequest())
	if err == nil {
		t.Fatal("expected error for empty content, got nil")
	}
}

func TestComplete_AuthError(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	defer srv.Close()

	_, err := p.Complete(context.Background(), simpleRequest())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- Stream ---

func TestStream_Success(t *testing.T) {
	// Anthropic SSE format.
	sseLines := []string{
		"event: content_block_start",
		`data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
		"",
		"event: content_block_delta",
		`data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"He"}}`,
		"",
		"event: content_block_delta",
		`data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"llo"}}`,
		"",
		"event: message_delta",
		`data: {"type":"message_delta","delta":{"stop_reason":"end_turn"}}`,
		"",
		"event: message_stop",
		`data: {"type":"message_stop"}`,
		"",
	}

	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		for _, line := range sseLines {
			_, _ = fmt.Fprintln(w, line)
		}
	})
	defer srv.Close()

	reader, err := p.Stream(context.Background(), ai.StreamRequest{
		Context:  ai.RequestContext{Provider: "anthropic"},
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
	if finishReason != "end_turn" {
		t.Errorf("finishReason = %q, want %q", finishReason, "end_turn")
	}
}

func TestStream_HTTPError(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	defer srv.Close()

	_, err := p.Stream(context.Background(), ai.StreamRequest{
		Context:  ai.RequestContext{Provider: "anthropic"},
		Messages: []ai.Message{{Role: ai.RoleUser, Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- Embed ---

func TestEmbed_NotSupported(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		// Should never be called.
		t.Error("unexpected HTTP request to mock server")
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer srv.Close()

	_, err := p.Embed(context.Background(), ai.EmbedRequest{
		Context: ai.RequestContext{Provider: "anthropic"},
		Text:    "test",
	})
	if err == nil {
		t.Fatal("expected error for unsupported Embed, got nil")
	}
}

// --- system prompt ---

func TestComplete_SystemPromptSentAsTopLevelField(t *testing.T) {
	type msgBody struct {
		System   string `json:"system"`
		Messages []struct {
			Role string `json:"role"`
		} `json:"messages"`
	}

	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		var body msgBody
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.System != "You are helpful." {
			t.Errorf("System = %q, want %q", body.System, "You are helpful.")
		}
		// No message should have role "system" inside the messages array.
		for _, m := range body.Messages {
			if m.Role == "system" {
				t.Errorf("found system role inside messages array: %+v", body.Messages)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(anthropicSuccessResponse("ok"))
	})
	defer srv.Close()

	req := simpleRequest()
	req.Context.SystemPrompt = "You are helpful."
	_, err := p.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
}
