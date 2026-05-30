package helpers

import (
	"errors"
	"strings"
	"testing"
)

func TestRedactLogMessageRedactsErrorValues(t *testing.T) {
	got := RedactLogMessage(errors.New("provider failed: Authorization: Bearer raw-token"))

	if strings.Contains(got, "raw-token") {
		t.Fatalf("RedactLogMessage(error) leaked token in %q", got)
	}
	if !strings.Contains(got, LogRedactedValue) {
		t.Fatalf("RedactLogMessage(error) = %q, want redacted placeholder", got)
	}
}

func TestRedactLogMessageRedactsStringValues(t *testing.T) {
	got := RedactLogMessage("api_key=sk-live-secret")

	if strings.Contains(got, "sk-live-secret") {
		t.Fatalf("RedactLogMessage(string) leaked api key in %q", got)
	}
	if !strings.Contains(got, LogRedactedValue) {
		t.Fatalf("RedactLogMessage(string) = %q, want redacted placeholder", got)
	}
}

func TestRedactLogMessageHandlesNil(t *testing.T) {
	if got := RedactLogMessage(nil); got != "" {
		t.Fatalf("RedactLogMessage(nil) = %q, want empty string", got)
	}
}
