package retry

import (
	"errors"

	"github.com/cenkalti/backoff/v4"
)

// PermanentError is the type returned by Permanent.
// Callers should not create PermanentError values directly; use Permanent(err).
type PermanentError = backoff.PermanentError

// Permanent wraps err to signal Do that no further retries should be
// attempted. The unwrapped error is returned to the caller of Do.
//
// Use this for errors that cannot be resolved by waiting:
//   - HTTP 4xx (except 429 Too Many Requests)
//   - Authentication / authorisation failures
//   - JSON decode / schema validation failures
//   - Business-logic invariant violations
//
// Do NOT use Permanent for context-cancellation errors; simply return
// ctx.Err() and the retry loop will stop on its own.
func Permanent(err error) error {
	return backoff.Permanent(err)
}

// IsPermanent reports whether err is (or wraps) a PermanentError.
// Useful in callers that need to distinguish permanent failures from
// transient ones after Do returns.
func IsPermanent(err error) bool {
	var pErr *backoff.PermanentError
	return errors.As(err, &pErr)
}
