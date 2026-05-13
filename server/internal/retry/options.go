package retry

import (
	"time"

	"github.com/meshery/meshkit/logger"
)

const (
	// DefaultInitialInterval is the wait before the second attempt.
	DefaultInitialInterval = 500 * time.Millisecond

	// DefaultMaxInterval caps the per-attempt wait ceiling.
	DefaultMaxInterval = 30 * time.Second

	// DefaultMaxElapsedTime caps the total wall-clock time across all attempts.
	// Zero disables this limit; pair with WithMaxAttempts instead.
	DefaultMaxElapsedTime = 2 * time.Minute

	// DefaultMultiplier controls exponential growth rate of the wait interval.
	DefaultMultiplier = 1.5

	// DefaultRandomizationFactor adds ±jitter to each wait interval.
	// Never set to 0 in production — it prevents thundering-herd stampedes.
	DefaultRandomizationFactor = 0.3
)

// Config holds all tuneable knobs for a retry operation.
// The zero value is valid; Do applies production-safe defaults before the
// first attempt. Override individual fields with the With* option functions.
type Config struct {
	// MaxAttempts caps the total number of calls (first attempt + retries).
	// 0 means unlimited — rely on MaxElapsedTime to bound the loop instead.
	MaxAttempts uint64

	// InitialInterval is the backoff wait before the second attempt.
	InitialInterval time.Duration

	// MaxInterval is the upper ceiling for a single per-attempt wait.
	MaxInterval time.Duration

	// MaxElapsedTime is the upper bound on total time spent across all
	// attempts (including wait intervals). 0 disables the wall-clock cap.
	MaxElapsedTime time.Duration

	// Multiplier is the factor by which the interval grows each attempt.
	Multiplier float64

	// RandomizationFactor adds jitter in the range
	//   [interval * (1 - f), interval * (1 + f)].
	RandomizationFactor float64

	// Notifier is called after each transient failure, before the next
	// wait interval begins. It receives the error and the computed wait
	// duration. nil means no-op.
	Notifier func(err error, wait time.Duration)
}

// defaultConfig returns a Config populated with production-safe defaults.
func defaultConfig() Config {
	return Config{
		InitialInterval:     DefaultInitialInterval,
		MaxInterval:         DefaultMaxInterval,
		MaxElapsedTime:      DefaultMaxElapsedTime,
		Multiplier:          DefaultMultiplier,
		RandomizationFactor: DefaultRandomizationFactor,
	}
}

// Option is a functional option that mutates a Config.
type Option func(*Config)

// WithMaxAttempts sets a hard cap on the total number of calls.
// The count includes the first (non-retry) attempt, so
// WithMaxAttempts(1) means "try once, no retries".
func WithMaxAttempts(n uint64) Option {
	return func(c *Config) { c.MaxAttempts = n }
}

// WithInitialInterval overrides the wait before the second attempt.
func WithInitialInterval(d time.Duration) Option {
	return func(c *Config) { c.InitialInterval = d }
}

// WithMaxInterval overrides the per-attempt wait ceiling.
func WithMaxInterval(d time.Duration) Option {
	return func(c *Config) { c.MaxInterval = d }
}

// WithMaxElapsedTime sets the wall-clock deadline across all attempts.
// Pass 0 to disable the elapsed-time cap entirely (use WithMaxAttempts
// to bound the loop in that case).
func WithMaxElapsedTime(d time.Duration) Option {
	return func(c *Config) { c.MaxElapsedTime = d }
}

// WithMultiplier overrides the exponential growth rate of wait intervals.
func WithMultiplier(m float64) Option {
	return func(c *Config) { c.Multiplier = m }
}

// WithJitter overrides the randomization factor (range: 0.0–1.0).
// Do not set to 0.0 in production — this disables jitter entirely, which
// can cause synchronized retry storms under load.
func WithJitter(f float64) Option {
	return func(c *Config) { c.RandomizationFactor = f }
}

// WithNotifier wires in an arbitrary per-attempt callback called on each
// transient failure, before the next wait interval begins.
func WithNotifier(n func(err error, wait time.Duration)) Option {
	return func(c *Config) { c.Notifier = n }
}

// WithLogNotifier builds a Notifier that emits a structured Warn log entry
// via MeshKit's logger.Handler each time a transient error triggers a retry.
//
// Example:
//
//	err := retry.Do(ctx, op, retry.WithLogNotifier(l.Log))
func WithLogNotifier(log logger.Handler) Option {
	return WithNotifier(func(err error, wait time.Duration) {
		log.Warnf("retry: transient error — retrying in %s: %v", wait.Round(time.Millisecond), err)
	})
}
