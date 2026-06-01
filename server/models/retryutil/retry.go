// Package retryutil provides exponential-backoff retry with a
// context-aware API, wrapping github.com/cenkalti/backoff/v5.
//
// REVIEWERS: This is a TEMPORARY local shim until meshery/meshkit#1007
// merges and github.com/meshery/meshkit/retry becomes available on the
// upstream module. The fork-replace approach (../meshkit or a personal
// fork pseudo-version) would break CI or can't be pushed upstream, so
// we inline the same API here instead.
//
// Once meshkit#1007 merges:
//   1. rm -rf server/models/retryutil/
//   2. In server/models/remote_provider.go, flip the import back to
//      github.com/meshery/meshkit/retry (remove the alias).
//   3. Uncomment the replace in go.mod (or bump the require version).
//
// See CLEANUP.md in this directory for the exact commands.
package retryutil

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/meshery/meshkit/logger"
)

// Operation is a retryable closure.
type Operation func(ctx context.Context) error

// Config controls the retry budget and backoff parameters.
type Config struct {
	MaxAttempts     uint
	InitialInterval time.Duration
	MaxInterval     time.Duration
	MaxElapsedTime  time.Duration
	Multiplier      float64
	Notifier        func(err error, wait time.Duration)
}

func defaultConfig() Config {
	return Config{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		MaxElapsedTime:  2 * time.Minute,
		Multiplier:      1.5,
	}
}

// Option applies a configuration change.
type Option func(*Config)

// WithMaxAttempts sets a hard cap on total calls (includes first attempt).
func WithMaxAttempts(n uint) Option {
	return func(c *Config) { c.MaxAttempts = n }
}

// WithInitialInterval sets the base sleep between retries.
func WithInitialInterval(d time.Duration) Option {
	return func(c *Config) { c.InitialInterval = d }
}

// WithMaxInterval sets the upper bound on each sleep.
func WithMaxInterval(d time.Duration) Option {
	return func(c *Config) { c.MaxInterval = d }
}

// WithMaxElapsedTime sets the wall-clock deadline. Pass 0 to disable.
func WithMaxElapsedTime(d time.Duration) Option {
	return func(c *Config) { c.MaxElapsedTime = d }
}

// WithMultiplier sets the exponential growth factor.
func WithMultiplier(m float64) Option {
	return func(c *Config) {
		if m <= 0 {
			c.Multiplier = 1.5
			return
		}
		c.Multiplier = m
	}
}

// WithLogNotifier emits a warning log entry on each retry.
func WithLogNotifier(log logger.Handler) Option {
	return func(c *Config) {
		c.Notifier = func(err error, wait time.Duration) {
			log.Infof("retry: transient error; retrying in %s", wait.Round(time.Millisecond))
			log.Warn(err)
		}
	}
}

// Do executes op with exponential backoff until success, context
// cancellation, or budget exhaustion.
func Do(ctx context.Context, op Operation, opts ...Option) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	retryOpts := []backoff.RetryOption{
		backoff.WithBackOff(newBackOff(cfg)),
		backoff.WithMaxElapsedTime(cfg.MaxElapsedTime),
	}
	if cfg.Notifier != nil {
		retryOpts = append(retryOpts, backoff.WithNotify(cfg.Notifier))
	}
	if cfg.MaxAttempts > 0 {
		retryOpts = append(retryOpts, backoff.WithMaxTries(cfg.MaxAttempts))
	}

	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		return struct{}{}, op(ctx)
	}, retryOpts...)
	return err
}

func newBackOff(cfg Config) backoff.BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = cfg.InitialInterval
	b.MaxInterval = cfg.MaxInterval
	b.Multiplier = cfg.Multiplier
	return b
}
