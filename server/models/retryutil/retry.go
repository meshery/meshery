// Package retryutil provides exponential-backoff retry with a
// context-aware API, wrapping github.com/cenkalti/backoff/v5.
//
// This is a TEMPORARY local shim until meshery/meshkit#1007 merges and
// github.com/meshery/meshkit/retry becomes available on the upstream
// module. Once meshkit#1007 merges:
//  1. rm -rf server/models/retryutil/
//  2. In server/models/remote_provider.go, flip the import to
//     github.com/meshery/meshkit/retry (remove the alias).
//  3. Bump github.com/meshery/meshkit in go.mod.
package retryutil

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/cenkalti/backoff/v5"
)

// Operation is a retryable closure.
type Operation func(ctx context.Context) error

// Do executes op with exponential backoff until success, permanent error,
// context cancellation, or budget exhaustion. Config via opts (default:
// 500ms initial, 1.5x growth, 30% jitter, 2min cap).
//
// When a ErrorClassifier is configured via WithErrorClassifier, every non-nil
// error from op (except those explicitly wrapped with Permanent) is passed to
// the classifier before the retry decision is made.
func Do(ctx context.Context, op Operation, opts ...Option) error {
	if err := ctx.Err(); err != nil {
		return ErrContext(err)
	}
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	if err := validateConfig(cfg); err != nil {
		return ErrConfig(err)
	}

	apply := op
	if cfg.ErrorClassifier != nil {
		apply = func(ctx context.Context) error {
			err := op(ctx)
			if err == nil {
				return nil
			}
			var pErr *backoff.PermanentError
			if errors.As(err, &pErr) {
				return err
			}
			if cfg.ErrorClassifier(err) == DecisionStop {
				return backoff.Permanent(err)
			}
			return err
		}
	}

	retryOpts := []backoff.RetryOption{
		backoff.WithBackOff(buildBackOff(cfg)),
		backoff.WithMaxElapsedTime(cfg.MaxElapsedTime),
		backoff.WithNotify(cfg.Notifier),
	}
	if cfg.MaxAttempts > 0 {
		retryOpts = append(retryOpts, backoff.WithMaxTries(cfg.MaxAttempts))
	}

	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		return struct{}{}, apply(ctx)
	}, retryOpts...)
	if err != nil {
		return ErrRetry(err)
	}
	return nil
}

func validateConfig(cfg Config) error {
	if cfg.InitialInterval <= 0 {
		return fmt.Errorf("%w: InitialInterval must be > 0, got %v", ErrInvalidConfig, cfg.InitialInterval)
	}
	if cfg.MaxInterval <= 0 {
		return fmt.Errorf("%w: MaxInterval must be > 0, got %v", ErrInvalidConfig, cfg.MaxInterval)
	}
	if cfg.MaxInterval < cfg.InitialInterval {
		return fmt.Errorf("%w: MaxInterval (%v) must be >= InitialInterval (%v)", ErrInvalidConfig, cfg.MaxInterval, cfg.InitialInterval)
	}
	if cfg.MaxElapsedTime < 0 {
		return fmt.Errorf("%w: MaxElapsedTime must be >= 0, got %v", ErrInvalidConfig, cfg.MaxElapsedTime)
	}
	if math.IsNaN(cfg.Multiplier) || math.IsInf(cfg.Multiplier, 0) || cfg.Multiplier < 1 {
		return fmt.Errorf("%w: Multiplier must be finite and >= 1, got %v", ErrInvalidConfig, cfg.Multiplier)
	}
	if math.IsNaN(cfg.RandomizationFactor) || cfg.RandomizationFactor < 0 || cfg.RandomizationFactor > 1 {
		return fmt.Errorf("%w: RandomizationFactor must be finite and in [0,1], got %v", ErrInvalidConfig, cfg.RandomizationFactor)
	}
	return nil
}

func buildBackOff(cfg Config) backoff.BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = cfg.InitialInterval
	b.MaxInterval = cfg.MaxInterval
	b.Multiplier = cfg.Multiplier
	b.RandomizationFactor = cfg.RandomizationFactor
	return b
}
