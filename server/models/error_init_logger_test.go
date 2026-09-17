package models

import (
	"errors"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

func TestErrInitLogger(t *testing.T) {
	underlying := errors.New("failed to open log output file")

	err := ErrInitLogger(underlying)
	if err == nil {
		t.Fatal("ErrInitLogger returned nil, expected a wrapped error")
	}

	if meshkiterrors.GetCode(err) != ErrInitLoggerCode {
		t.Errorf("expected error code %q, got %q", ErrInitLoggerCode, meshkiterrors.GetCode(err))
	}

	if got := err.Error(); got == "" {
		t.Error("expected non-empty error message")
	}
}
