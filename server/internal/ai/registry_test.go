package ai

import (
	"context"
	"io"
	"testing"
)

// mock implements LLMProvider for registry tests.
type mock struct {
	name string
}

func (m *mock) Name() string { return m.name }
func (m *mock) Complete(context.Context, CompleteRequest) (CompleteResponse, error) {
	return CompleteResponse{}, nil
}
func (m *mock) Stream(context.Context, CompleteRequest) (io.ReadCloser, error) { return nil, nil }
func (m *mock) Embed(context.Context, EmbedRequest) (EmbedResponse, error) {
	return EmbedResponse{}, nil
}

func TestRegistryRegisterAndGet(t *testing.T) {
	r := NewRegistry()
	p := &mock{name: "test"}
	r.Register(p)

	got, err := r.Get("test")
	if err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}
	if got.Name() != "test" {
		t.Fatalf("got provider %q, want %q", got.Name(), "test")
	}
}

func TestRegistryGetMissing(t *testing.T) {
	r := NewRegistry()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing provider, got nil")
	}
}

func TestRegistryHas(t *testing.T) {
	r := NewRegistry()
	r.Register(&mock{name: "ollama"})

	if !r.Has("ollama") {
		t.Fatal("Has(ollama) = false, want true")
	}
	if r.Has("openai") {
		t.Fatal("Has(openai) = true, want false")
	}
}

func TestRegistryNames(t *testing.T) {
	r := NewRegistry()
	r.Register(&mock{name: "ollama"})
	r.Register(&mock{name: "anthropic"})
	r.Register(&mock{name: "openai"})

	names := r.Names()
	want := []string{"anthropic", "ollama", "openai"}
	if len(names) != len(want) {
		t.Fatalf("Names() returned %d items, want %d", len(names), len(want))
	}
	for i, n := range names {
		if n != want[i] {
			t.Errorf("Names()[%d] = %q, want %q", i, n, want[i])
		}
	}
}

func TestRegistryOverwrite(t *testing.T) {
	r := NewRegistry()
	r.Register(&mock{name: "ollama"})
	r.Register(&mock{name: "ollama"})

	names := r.Names()
	if len(names) != 1 {
		t.Fatalf("expected 1 provider after re-register, got %d", len(names))
	}
}

func TestRegistryRegisterNil(t *testing.T) {
	r := NewRegistry()
	r.Register(nil) // should not panic
	if len(r.Names()) != 0 {
		t.Fatal("expected 0 providers after nil register")
	}
}
