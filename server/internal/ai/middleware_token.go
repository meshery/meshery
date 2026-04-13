package ai

import (
	"context"
	"sync/atomic"
)

// TokenCounter accumulates token usage across provider calls.
type TokenCounter struct {
	InputTokens  atomic.Int64
	OutputTokens atomic.Int64
}

// Add records usage from a single call.
func (c *TokenCounter) Add(usage *TokenUsage) {
	if usage == nil {
		return
	}
	c.InputTokens.Add(int64(usage.InputTokens))
	c.OutputTokens.Add(int64(usage.OutputTokens))
}

// TokenAccountingMiddleware returns a middleware that accumulates token usage
// into the supplied counter. Token accounting happens before the retry
// middleware so that retried requests are not double-counted.
func TokenAccountingMiddleware(counter *TokenCounter) Middleware {
	return func(next LLMProvider) LLMProvider {
		return &tokenProvider{next: next, counter: counter}
	}
}

type tokenProvider struct {
	next    LLMProvider
	counter *TokenCounter
}

func (p *tokenProvider) Complete(ctx context.Context, req CompleteRequest) (CompleteResponse, error) {
	resp, err := p.next.Complete(ctx, req)
	if err == nil {
		p.counter.Add(resp.Usage)
	}
	return resp, err
}

func (p *tokenProvider) Stream(ctx context.Context, req StreamRequest) (StreamReader, error) {
	// Streaming token counts are typically reported in the final chunk.
	// Accounting for streaming is handled at the application layer when
	// the stream is fully consumed and the finish chunk is inspected.
	return p.next.Stream(ctx, req)
}

func (p *tokenProvider) Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error) {
	resp, err := p.next.Embed(ctx, req)
	if err == nil {
		p.counter.Add(resp.Usage)
	}
	return resp, err
}
