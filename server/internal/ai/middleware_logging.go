package ai

import (
	"context"
	"fmt"
	"time"
)

// Logger is the minimal logging interface used by the logging middleware.
// Compatible with meshkit/logger.Handler.
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

// LoggingMiddleware returns a middleware that logs request metadata and latency
// for every provider call using the supplied logger.
func LoggingMiddleware(log Logger) Middleware {
	return func(next LLMProvider) LLMProvider {
		return &loggingProvider{next: next, log: log}
	}
}

type loggingProvider struct {
	next LLMProvider
	log  Logger
}

func (p *loggingProvider) Complete(ctx context.Context, req CompleteRequest) (CompleteResponse, error) {
	start := time.Now()
	resp, err := p.next.Complete(ctx, req)
	elapsed := time.Since(start)

	if err != nil {
		p.log.Error(fmt.Sprintf("[ai] Complete provider=%s model=%s latency=%s error=%v",
			req.Context.Provider, req.Context.Model, elapsed, err))
	} else {
		tokens := 0
		if resp.Usage != nil {
			tokens = resp.Usage.TotalTokens
		}
		p.log.Info(fmt.Sprintf("[ai] Complete provider=%s model=%s latency=%s tokens=%d finish=%s",
			req.Context.Provider, req.Context.Model, elapsed, tokens, resp.FinishReason))
	}

	return resp, err
}

func (p *loggingProvider) Stream(ctx context.Context, req StreamRequest) (StreamReader, error) {
	start := time.Now()
	reader, err := p.next.Stream(ctx, req)
	elapsed := time.Since(start)

	if err != nil {
		p.log.Error(fmt.Sprintf("[ai] Stream provider=%s model=%s latency=%s error=%v",
			req.Context.Provider, req.Context.Model, elapsed, err))
	} else {
		p.log.Info(fmt.Sprintf("[ai] Stream provider=%s model=%s ttft=%s",
			req.Context.Provider, req.Context.Model, elapsed))
	}

	return reader, err
}

func (p *loggingProvider) Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error) {
	start := time.Now()
	resp, err := p.next.Embed(ctx, req)
	elapsed := time.Since(start)

	if err != nil {
		p.log.Error(fmt.Sprintf("[ai] Embed provider=%s model=%s latency=%s error=%v",
			req.Context.Provider, req.Context.Model, elapsed, err))
	} else {
		p.log.Info(fmt.Sprintf("[ai] Embed provider=%s model=%s latency=%s dims=%d",
			req.Context.Provider, req.Context.Model, elapsed, len(resp.Vector)))
	}

	return resp, err
}
