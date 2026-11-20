package handlers

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/sirupsen/logrus"
)

// The swagger annotations on LogLevelResponse and LogLevelRequest, and the
// LogLevel enum in meshery/schemas v1beta3/event, all name `warning`. That is
// what logrus serializes WarnLevel to, and what this endpoint therefore reports
// and stores. Both documents said `warn` for a while, which is a value the
// endpoint never returns. Pin the serialization so the documents cannot drift
// back, and pin the alias so nobody "fixes" the parser into rejecting it.
func TestAvailableLogLevelsUseCanonicalLogrusNames(t *testing.T) {
	levels := getAvailableLogLevels()

	want := []string{"panic", "fatal", "error", "warning", "info", "debug", "trace"}
	if !slices.Equal(levels, want) {
		t.Fatalf("expected available levels %v, got %v", want, levels)
	}
	if slices.Contains(levels, "warn") {
		t.Errorf("available levels must report the canonical %q, never %q: %v", "warning", "warn", levels)
	}
}

func TestParseLevelAcceptsWarnAsAnAliasOfWarning(t *testing.T) {
	for _, requested := range []string{"warning", "warn"} {
		level, err := logrus.ParseLevel(requested)
		if err != nil {
			t.Fatalf("expected %q to be accepted, got %v", requested, err)
		}
		if level != logrus.WarnLevel {
			t.Errorf("expected %q to resolve to WarnLevel, got %v", requested, level)
		}
		if level.String() != "warning" {
			t.Errorf("expected %q to serialize back as %q, got %q", requested, "warning", level.String())
		}
	}
}

// The request body carries the canonical `eventLogLevel` key and the legacy
// `event_log_level` one; canonical wins when both are present.
func TestLogLevelRequestAcceptsCanonicalAndLegacyKeys(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{"canonical key", `{"eventLogLevel":"warning"}`, "warning"},
		{"legacy key", `{"event_log_level":"warning"}`, "warning"},
		{"canonical wins over legacy", `{"eventLogLevel":"debug","event_log_level":"warning"}`, "debug"},
		{"neither key", `{}`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req LogLevelRequest
			if err := json.Unmarshal([]byte(tt.body), &req); err != nil {
				t.Fatalf("decoding %s: %v", tt.body, err)
			}
			if req.LogLevel != tt.want {
				t.Errorf("expected log level %q from %s, got %q", tt.want, tt.body, req.LogLevel)
			}
		})
	}
}
