package ollama_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/internal/ai"
	"github.com/meshery/meshery/server/internal/ai/providers/ollama"
)

// --- helpers ---

func newTestServer(handler http.HandlerFunc) (*httptest.Server, *ollama.Provider) {
	srv := httptest.NewServer(handler)
	p := ollama.New(
		ollama.WithEndpoint(srv.URL),
		ollama.WithHTTPClient(srv.Client()),
	)
	return srv, p
}

func simpleRequest() ai.CompleteRequest {
	return ai.CompleteRequest{
		Context:  ai.RequestContext{Provider: "ollama"},
		Messages: []ai.Message{{Role: ai.RoleUser, Content: "ping"}},
	}
}

// --- Complete ---

func TestComplete_Success(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := map[string]interface{}{
			"model":   "llama3",
			"message": map[string]string{"role": "assistant", "content": "pong"},
			"done":    true,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
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
}

func TestComplete_ServerError(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer srv.Close()

	_, err := p.Complete(context.Background(), simpleRequest())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var re *ai.RetryableError
	if !isRetryableError(err, &re) {
		t.Fatalf("expected RetryableError, got %T: %v", err, err)
	}
	if re.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", re.StatusCode, http.StatusInternalServerError)
	}
}

func TestComplete_ContextCancelled(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		// Never responds — the client context will cancel first.
		<-r.Context().Done()
	})
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := p.Complete(ctx, simpleRequest())
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

// --- Stream ---

func TestStream_Success(t *testing.T) {
	chunks := []map[string]interface{}{
		{"message": map[string]string{"role": "assistant", "content": "He"}, "done": false},
		{"message": map[string]string{"role": "assistant", "content": "llo"}, "done": false},
		{"message": map[string]string{"role": "assistant", "content": ""}, "done": true},
	}

	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		for _, c := range chunks {
			_ = json.NewEncoder(w).Encode(c)
		}
	})
	defer srv.Close()

	reader, err := p.Stream(context.Background(), ai.StreamRequest{
		Context:  ai.RequestContext{Provider: "ollama"},
		Messages: []ai.Message{{Role: ai.RoleUser, Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	var collected string
	for {
		chunk, err := reader.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv() error = %v", err)
		}
		collected += chunk.Content
	}

	if collected != "Hello" {
		t.Errorf("collected = %q, want %q", collected, "Hello")
	}
}

// --- Embed ---

func TestEmbed_Success(t *testing.T) {
	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embeddings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := map[string]interface{}{
			"embedding": []float32{0.1, 0.2, 0.3},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	got, err := p.Embed(context.Background(), ai.EmbedRequest{
		Context: ai.RequestContext{Provider: "ollama"},
		Text:    "test text",
	})
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}
	if len(got.Vector) != 3 {
		t.Errorf("Vector length = %d, want 3", len(got.Vector))
	}
}

// --- request body validation ---

func TestComplete_SendsSystemPrompt(t *testing.T) {
	type chatReq struct {
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}

	srv, p := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		var req chatReq
		_ = json.NewDecoder(r.Body).Decode(&req)
		if len(req.Messages) < 2 || req.Messages[0].Role != "system" {
			t.Errorf("expected system message as first message, got: %+v", req.Messages)
		}
		resp := map[string]interface{}{
			"message": map[string]string{"role": "assistant", "content": "ok"},
			"done":    true,
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	req := simpleRequest()
	req.Context.SystemPrompt = "You are a helpful assistant."
	_, err := p.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
}

// isRetryableError uses errors.As to unwrap the error chain and find a RetryableError.
func isRetryableError(err error, target **ai.RetryableError) bool {
	return errors.As(err, target)
}
