package ai

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"
)

type stubProvider struct{}

func (stubProvider) Complete(context.Context, CompleteRequest) (CompleteResponse, error) {
	return CompleteResponse{}, nil
}

func (stubProvider) Stream(context.Context, StreamRequest) (StreamReader, error) {
	return stubStreamReader{}, nil
}

func (stubProvider) Embed(context.Context, EmbedRequest) (EmbedResponse, error) {
	return EmbedResponse{}, nil
}

type stubStreamReader struct{}

func (stubStreamReader) Recv() (StreamChunk, error) {
	return StreamChunk{}, io.EOF
}

type errorStreamReader struct {
	err error
}

func (r errorStreamReader) Recv() (StreamChunk, error) {
	return StreamChunk{}, r.err
}

func TestRegistryRegisterAndGet(t *testing.T) {
	registry := NewRegistry()
	provider := stubProvider{}

	if err := registry.Register("ollama", provider); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	got, err := registry.Get("ollama")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if !reflect.DeepEqual(got, provider) {
		t.Fatalf("Get() = %#v, want %#v", got, provider)
	}

	if !registry.Has("ollama") {
		t.Fatal("Has() = false, want true")
	}
}

func TestRegistryRegisterValidatesInputs(t *testing.T) {
	registry := NewRegistry()

	testCases := []struct {
		name     string
		key      string
		provider LLMProvider
		wantErr  error
	}{
		{name: "empty name", key: "", provider: stubProvider{}, wantErr: ErrProviderNameRequired},
		{name: "nil provider", key: "ollama", provider: nil, wantErr: ErrProviderNil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := registry.Register(tc.key, tc.provider)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Register() error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestRegistryRejectsDuplicateRegistration(t *testing.T) {
	registry := NewRegistry()

	if err := registry.Register("openai", stubProvider{}); err != nil {
		t.Fatalf("Register() first call error = %v", err)
	}

	err := registry.Register("openai", stubProvider{})
	if !errors.Is(err, ErrProviderAlreadyRegistered) {
		t.Fatalf("Register() duplicate error = %v, want %v", err, ErrProviderAlreadyRegistered)
	}
}

func TestRegistryGetMissingProvider(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Get("missing")
	if !errors.Is(err, ErrProviderNotFound) {
		t.Fatalf("Get() error = %v, want %v", err, ErrProviderNotFound)
	}
}

func TestRegistryNamesAreSorted(t *testing.T) {
	registry := NewRegistry()

	for _, name := range []string{"openai", "anthropic", "ollama"} {
		if err := registry.Register(name, stubProvider{}); err != nil {
			t.Fatalf("Register(%q) error = %v", name, err)
		}
	}

	got := registry.Names()
	want := []string{"anthropic", "ollama", "openai"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Names() = %v, want %v", got, want)
	}
}

func TestDrainStream(t *testing.T) {
	chunks, err := DrainStream(stubStreamReader{})
	if err != nil {
		t.Fatalf("DrainStream() error = %v", err)
	}

	if len(chunks) != 0 {
		t.Fatalf("DrainStream() len = %d, want 0", len(chunks))
	}

	expectedErr := errors.New("stream failed")
	_, err = DrainStream(errorStreamReader{err: expectedErr})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("DrainStream() error = %v, want %v", err, expectedErr)
	}
}
