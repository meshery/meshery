package retryutil

import (
	"errors"

	"github.com/cenkalti/backoff/v5"
	meshkiterrors "github.com/meshery/meshkit/errors"
)

// ErrInvalidConfig is returned when retry configuration validation fails.
var ErrInvalidConfig = errors.New("retry: invalid config")

const (
	ErrRetryCode   = "meshkit-10001"
	ErrContextCode = "meshkit-10002"
	ErrConfigCode  = "meshkit-10003"
)

type retryError struct {
	inner   error
	meshkit *meshkiterrors.Error
}

func (e *retryError) Error() string {
	return e.meshkit.Error()
}

func (e *retryError) Unwrap() []error {
	return []error{e.inner, e.meshkit}
}

func ErrRetry(err error) error {
	return &retryError{
		inner: err,
		meshkit: meshkiterrors.New(ErrRetryCode, meshkiterrors.Alert,
			[]string{"Retry operation failed"},
			[]string{err.Error()},
			[]string{"Operation did not succeed within retry limits"},
			[]string{"Check the underlying operation and retry configuration"}),
	}
}

func ErrContext(err error) error {
	return &retryError{
		inner: err,
		meshkit: meshkiterrors.New(ErrContextCode, meshkiterrors.Alert,
			[]string{"Context canceled or deadline exceeded"},
			[]string{err.Error()},
			[]string{"Operation timed out or context was canceled"},
			[]string{"Check context timeout and ensure the operation completes in time"}),
	}
}

func ErrConfig(err error) error {
	return &retryError{
		inner: err,
		meshkit: meshkiterrors.New(ErrConfigCode, meshkiterrors.Alert,
			[]string{"Invalid retry configuration"},
			[]string{err.Error()},
			[]string{"One or more config values are invalid"},
			[]string{"Ensure all retry configuration values are correct"}),
	}
}

// ErrorDecision controls retry behaviour for a single error.
type ErrorDecision int

const (
	DecisionRetry ErrorDecision = iota
	DecisionStop
)

// ErrorClassifier returns the retry decision for a given error.
// Return DecisionStop for errors that should not be retried (e.g. HTTP 4xx,
// validation failures, auth errors). Return DecisionRetry for transient
// errors (timeouts, 5xx, rate limits).
type ErrorClassifier func(err error) ErrorDecision

// Permanent wraps err to signal no further retries should be attempted.
// Use for non-transient errors (HTTP 4xx, auth failures, validation errors).
func Permanent(err error) error {
	return backoff.Permanent(err)
}

// IsPermanent reports whether err is (or wraps) a PermanentError.
func IsPermanent(err error) bool {
	var pErr *backoff.PermanentError
	return errors.As(err, &pErr)
}
