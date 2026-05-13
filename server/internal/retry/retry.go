// Package retry provides a centralized, context-aware exponential-backoff
// retry abstraction for the Meshery server.
//
// # Overview
//
// The single entry-point is [Do]. It wraps [github.com/cenkalti/backoff/v4]
// and adds:
//   - Production-safe defaults (initial 500 ms, ×1.5 growth, ±30 % jitter,
//     2 min total wall-clock cap).
//   - Context-first design — the loop terminates immediately when ctx is
//     cancelled or times out.
//   - Permanent-error escape hatch — wrap non-retryable errors with
//     [Permanent] to stop the loop without exhausting the retry budget.
//   - Composable functional options (see the With* helpers in options.go).
//   - Optional per-attempt logging via [WithLogNotifier].
//
// # Basic usage
//
//	err := retry.Do(ctx, func() error {
//	    resp, err := http.Get(url)
//	    if err != nil {
//	        return err // transient — will be retried
//	    }
//	    if resp.StatusCode == http.StatusBadRequest {
//	        return retry.Permanent(fmt.Errorf("bad request: %d", resp.StatusCode))
//	    }
//	    return nil
//	}, retry.WithMaxAttempts(5), retry.WithLogNotifier(log))
//
// # Error classification
//
//   - Network errors, HTTP 5xx, HTTP 429  → return err (transient)
//   - HTTP 4xx (except 429), auth, decode → return retry.Permanent(err)
//   - Context cancelled / deadline exceeded → return ctx.Err() (loop exits
//     automatically; do not wrap in Permanent)
package retry

import (
	"context"

	"github.com/cenkalti/backoff/v4"
)

// Operation is the function signature for retryable work.
//
// Return nil on success, [Permanent](err) to stop without further retries,
// or any plain error to trigger the next backoff wait and retry.
type Operation func() error

// Do executes op with exponential backoff until one of the following occurs:
//   - op returns nil (success)
//   - op returns a [Permanent] error (stop immediately; unwrapped error returned)
//   - ctx is cancelled or its deadline is exceeded
//   - the configured MaxAttempts or MaxElapsedTime budget is exhausted
//
// All retry configuration is supplied via the variadic opts. When no opts are
// given, [defaultConfig] values are used (500 ms initial, ×1.5 growth, ±30 %
// jitter, 2-min elapsed cap). Individual knobs are overridden with the With*
// helpers defined in options.go.
//
// Do is safe for concurrent use; each call creates its own backoff state.
func Do(ctx context.Context, op Operation, opts ...Option) error {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	b := buildBackOff(cfg)
	bCtx := backoff.WithContext(b, ctx)

	if cfg.Notifier != nil {
		return backoff.RetryNotify(backoff.Operation(op), bCtx, cfg.Notifier)
	}
	return backoff.Retry(backoff.Operation(op), bCtx)
}

// buildBackOff constructs a cenkalti/backoff policy from the supplied Config.
// When MaxAttempts is set, the exponential backoff is wrapped with a
// WithMaxRetries limiter; otherwise the raw exponential policy is returned and
// MaxElapsedTime acts as the sole termination condition.
func buildBackOff(cfg Config) backoff.BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = cfg.InitialInterval
	b.MaxInterval = cfg.MaxInterval
	b.MaxElapsedTime = cfg.MaxElapsedTime
	b.Multiplier = cfg.Multiplier
	b.RandomizationFactor = cfg.RandomizationFactor

	// MaxAttempts is the total number of calls (1st attempt + retries).
	// cenkalti/backoff's WithMaxRetries counts extra retries after the first,
	// so we subtract 1.
	if cfg.MaxAttempts > 0 {
		return backoff.WithMaxRetries(b, cfg.MaxAttempts-1)
	}
	return b
}
