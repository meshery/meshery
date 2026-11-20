package ai

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"
)

type recordingProvider struct {
	completeResp CompleteResponse
	streamReader StreamReader
	embedResp    EmbedResponse

	completeReq CompleteRequest
	streamReq   StreamRequest
	embedReq    EmbedRequest
}

func (p *recordingProvider) Complete(_ context.Context, req CompleteRequest) (CompleteResponse, error) {
	p.completeReq = req
	return p.completeResp, nil
}

func (p *recordingProvider) Stream(_ context.Context, req StreamRequest) (StreamReader, error) {
	p.streamReq = req
	return p.streamReader, nil
}

func (p *recordingProvider) Embed(_ context.Context, req EmbedRequest) (EmbedResponse, error) {
	p.embedReq = req
	return p.embedResp, nil
}

type streamSequence struct {
	chunks []StreamChunk
	index  int
}

func (s *streamSequence) Recv() (StreamChunk, error) {
	if s.index >= len(s.chunks) {
		return StreamChunk{}, io.EOF
	}

	chunk := s.chunks[s.index]
	s.index++
	return chunk, nil
}

type middlewareRecorder struct {
	order *[]string
	name  string
}

func (m middlewareRecorder) wrap(next LLMProvider) LLMProvider {
	return middlewareProvider{
		next: next,
		complete: func(ctx context.Context, req CompleteRequest) (CompleteResponse, error) {
			*m.order = append(*m.order, m.name)
			return next.Complete(ctx, req)
		},
	}
}

type middlewareProvider struct {
	next     LLMProvider
	complete func(context.Context, CompleteRequest) (CompleteResponse, error)
}

func (p middlewareProvider) Complete(ctx context.Context, req CompleteRequest) (CompleteResponse, error) {
	if p.complete != nil {
		return p.complete(ctx, req)
	}
	return p.next.Complete(ctx, req)
}

func (p middlewareProvider) Stream(ctx context.Context, req StreamRequest) (StreamReader, error) {
	return p.next.Stream(ctx, req)
}

func (p middlewareProvider) Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error) {
	return p.next.Embed(ctx, req)
}

func TestChainMiddlewarePreservesOrder(t *testing.T) {
	order := make([]string, 0, 2)
	base := &recordingProvider{}

	wrapped := ChainMiddleware(
		base,
		middlewareRecorder{order: &order, name: "token-accounting"}.wrap,
		middlewareRecorder{order: &order, name: "retry"}.wrap,
	)

	_, err := wrapped.Complete(context.Background(), CompleteRequest{})
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	want := []string{"token-accounting", "retry"}
	if !reflect.DeepEqual(order, want) {
		t.Fatalf("middleware order = %v, want %v", order, want)
	}
}

func TestServiceRequiresProviderSelection(t *testing.T) {
	service := NewService(nil)

	_, err := service.Complete(context.Background(), CompleteRequest{})
	if !errors.Is(err, ErrProviderSelectionRequired) {
		t.Fatalf("Complete() error = %v, want %v", err, ErrProviderSelectionRequired)
	}
}

func TestServiceRoutesCompleteRequests(t *testing.T) {
	registry := NewRegistry()
	provider := &recordingProvider{
		completeResp: CompleteResponse{Content: "generated"},
	}
	if err := registry.Register("ollama", provider); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	service := NewService(registry)
	req := CompleteRequest{
		Context: RequestContext{Provider: "ollama"},
		Messages: []Message{
			{Role: RoleUser, Content: "hello"},
		},
	}

	resp, err := service.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if resp.Content != "generated" {
		t.Fatalf("Complete() content = %q, want %q", resp.Content, "generated")
	}
	if !reflect.DeepEqual(provider.completeReq, req) {
		t.Fatalf("provider.completeReq = %#v, want %#v", provider.completeReq, req)
	}
}

func TestServiceRoutesStreamRequests(t *testing.T) {
	registry := NewRegistry()
	provider := &recordingProvider{
		streamReader: &streamSequence{
			chunks: []StreamChunk{{Content: "a"}, {Content: "b"}},
		},
	}
	if err := registry.Register("openai", provider); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	service := NewService(registry)
	req := StreamRequest{
		Context: RequestContext{Provider: "openai"},
	}

	reader, err := service.Stream(context.Background(), req)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	chunks, err := DrainStream(reader)
	if err != nil {
		t.Fatalf("DrainStream() error = %v", err)
	}

	want := []StreamChunk{{Content: "a"}, {Content: "b"}}
	if !reflect.DeepEqual(chunks, want) {
		t.Fatalf("stream chunks = %#v, want %#v", chunks, want)
	}
	if !reflect.DeepEqual(provider.streamReq, req) {
		t.Fatalf("provider.streamReq = %#v, want %#v", provider.streamReq, req)
	}
}

func TestServiceRoutesEmbedRequests(t *testing.T) {
	registry := NewRegistry()
	provider := &recordingProvider{
		embedResp: EmbedResponse{Vector: []float32{1, 2, 3}},
	}
	if err := registry.Register("openai", provider); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	service := NewService(registry)
	req := EmbedRequest{
		Context: RequestContext{Provider: "openai"},
		Text:    "meshery",
	}

	resp, err := service.Embed(context.Background(), req)
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}
	if !reflect.DeepEqual(resp.Vector, []float32{1, 2, 3}) {
		t.Fatalf("Embed() vector = %#v, want %#v", resp.Vector, []float32{1, 2, 3})
	}
	if !reflect.DeepEqual(provider.embedReq, req) {
		t.Fatalf("provider.embedReq = %#v, want %#v", provider.embedReq, req)
	}
}
