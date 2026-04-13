package ai

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"time"
)

// RetryConfig controls the retry behaviour of the retry middleware.
type RetryConfig struct {
	// MaxAttempts is the total number of attempts (1 means no retries).
	MaxAttempts int
	// InitialDelay is the backoff duration before the first retry.
	InitialDelay time.Duration
	// MaxDelay caps the exponential backoff.
	MaxDelay time.Duration
	// JitterFraction adds randomness (0.0–1.0) to each delay.
	JitterFraction float64
}

// DefaultRetryConfig is a sensible starting point for most providers.
var DefaultRetryConfig = RetryConfig{
	MaxAttempts:    3,
	InitialDelay:   500 * time.Millisecond,
	MaxDelay:       10 * time.Second,
	JitterFraction: 0.2,
}

// RetryableError wraps an underlying error and signals whether the call should
// be retried. Provider implementations may return this type to communicate
// HTTP status codes back to the retry middleware.
type RetryableError struct {
	Err        error
	StatusCode int
}

func (e *RetryableError) Error() string { return e.Err.Error() }
func (e *RetryableError) Unwrap() error { return e.Err }

// isRetryable returns true for transient errors (rate limits, server errors).
// Auth failures and bad requests are never retried.
func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	var re *RetryableError
	if errors.As(err, &re) {
		switch re.StatusCode {
		case http.StatusTooManyRequests,
			http.StatusServiceUnavailable,
			http.StatusBadGateway,
			http.StatusGatewayTimeout,
			http.StatusInternalServerError:
			return true
		}
		return false
	}
	// Network-level errors (context not cancelled) are retryable.
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}
	return true
}

// RetryMiddleware returns a middleware that retries transient errors with
// exponential backoff and jitter.
func RetryMiddleware(cfg RetryConfig) Middleware {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	return func(next LLMProvider) LLMProvider {
		return &retryProvider{next: next, cfg: cfg}
	}
}

type retryProvider struct {
	next LLMProvider
	cfg  RetryConfig
}

func (p *retryProvider) Complete(ctx context.Context, req CompleteRequest) (CompleteResponse, error) {
	var resp CompleteResponse
	err := p.withRetry(ctx, func() error {
		var callErr error
		resp, callErr = p.next.Complete(ctx, req)
		return callErr
	})
	return resp, err
}

func (p *retryProvider) Stream(ctx context.Context, req StreamRequest) (StreamReader, error) {
	// Streaming is not retried automatically because the caller may have already
	// consumed part of the stream. Let the caller decide on retry semantics.
	return p.next.Stream(ctx, req)
}

func (p *retryProvider) Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error) {
	var resp EmbedResponse
	err := p.withRetry(ctx, func() error {
		var callErr error
		resp, callErr = p.next.Embed(ctx, req)
		return callErr
	})
	return resp, err
}

func (p *retryProvider) withRetry(ctx context.Context, fn func() error) error {
	delay := p.cfg.InitialDelay
	var lastErr error

	for attempt := 0; attempt < p.cfg.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if !isRetryable(lastErr) {
			return lastErr
		}
		if attempt == p.cfg.MaxAttempts-1 {
			break
		}

		jitter := time.Duration(float64(delay) * p.cfg.JitterFraction * rand.Float64())
		sleep := delay + jitter
		if sleep > p.cfg.MaxDelay {
			sleep = p.cfg.MaxDelay
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleep):
		}

		delay *= 2
		if delay > p.cfg.MaxDelay {
			delay = p.cfg.MaxDelay
		}
	}

	return lastErr
}
