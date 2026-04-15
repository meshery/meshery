package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/internal/ai"
)

func TestComplete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Model != "llama3" {
			t.Errorf("model = %q, want %q", req.Model, "llama3")
		}
		if req.Stream == nil || *req.Stream {
			t.Error("expected stream=false for Complete")
		}

		resp := chatResponse{
			Model:           "llama3",
			Message:         ollamaMessage{Role: "assistant", Content: "hello from ollama"},
			Done:            true,
			DoneReason:      "stop",
			PromptEvalCount: 10,
			EvalCount:       5,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	defer srv.Close()

	p := New(Config{BaseURL: srv.URL, Model: "llama3"})
	got, err := p.Complete(context.Background(), ai.CompleteRequest{
		Messages: []ai.Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if got.Content != "hello from ollama" {
		t.Errorf("content = %q, want %q", got.Content, "hello from ollama")
	}
	if got.Usage.PromptTokens != 10 || got.Usage.CompletionTokens != 5 {
		t.Errorf("usage = %+v, want prompt=10 completion=5", got.Usage)
	}
	if got.Usage.TotalTokens != 15 {
		t.Errorf("total_tokens = %d, want 15", got.Usage.TotalTokens)
	}
}

func TestCompleteHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"model not found"}`, http.StatusNotFound)
	}))
	defer srv.Close()

	p := New(Config{BaseURL: srv.URL, Model: "missing"})
	_, err := p.Complete(context.Background(), ai.CompleteRequest{
		Messages: []ai.Message{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

func TestCompleteNoModel(t *testing.T) {
	p := New(Config{BaseURL: "http://localhost:11434"})
	_, err := p.Complete(context.Background(), ai.CompleteRequest{
		Messages: []ai.Message{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error when no model is configured, got nil")
	}
}

func TestEmbed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embeddings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(embedResponse{
			Embedding: []float64{0.1, 0.2, 0.3},
		}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	defer srv.Close()

	p := New(Config{BaseURL: srv.URL, Model: "llama3"})
	got, err := p.Embed(context.Background(), ai.EmbedRequest{Text: "test input"})
	if err != nil {
		t.Fatalf("Embed: %v", err)
	}
	if len(got.Embedding) != 3 {
		t.Errorf("embedding length = %d, want 3", len(got.Embedding))
	}
}

func TestModels(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tagsResponse{
			Models: []struct {
				Name string `json:"name"`
			}{
				{Name: "llama3:latest"},
				{Name: "codellama:7b"},
			},
		}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	defer srv.Close()

	p := New(Config{BaseURL: srv.URL})
	names, err := p.Models(context.Background())
	if err != nil {
		t.Fatalf("Models: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("got %d models, want 2", len(names))
	}
	if names[0] != "llama3:latest" {
		t.Errorf("models[0] = %q, want %q", names[0], "llama3:latest")
	}
}
