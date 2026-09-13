package models

import (
	"strings"
	"testing"

	"github.com/meshery/meshkit/errors"
)

// TestParseEnvironmentConnectionsFilterDefaults pins the no-filter behaviour the
// endpoint has always had: assigned connections, unnarrowed.
func TestParseEnvironmentConnectionsFilterDefaults(t *testing.T) {
	for _, filter := range []string{"", "   "} {
		parsed, err := ParseEnvironmentConnectionsFilter(filter)
		if err != nil {
			t.Fatalf("ParseEnvironmentConnectionsFilter(%q) returned error: %v", filter, err)
		}
		if !parsed.Assigned {
			t.Errorf("ParseEnvironmentConnectionsFilter(%q) Assigned = false, want true", filter)
		}
		if parsed.Owner != "" {
			t.Errorf("ParseEnvironmentConnectionsFilter(%q) Owner = %q, want empty", filter, parsed.Owner)
		}
	}
}

func TestParseEnvironmentConnectionsFilterValid(t *testing.T) {
	tests := []struct {
		name         string
		filter       string
		wantAssigned bool
		wantOwner    string
	}{
		// The value Meshery UI sends; it must keep working verbatim.
		{"ui assigned false", `{"assigned":false}`, false, ""},
		{"assigned true", `{"assigned":true}`, true, ""},
		{"owner only", `{"owner":"1d3b1e0a"}`, true, "1d3b1e0a"},
		{"owner and assigned", `{"assigned":false,"owner":"1d3b1e0a"}`, false, "1d3b1e0a"},
		// An unrecognised field is ignored rather than rejected.
		{"unknown field ignored", `{"organization_id":"x"}`, true, ""},
		{"empty object", `{}`, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseEnvironmentConnectionsFilter(tt.filter)
			if err != nil {
				t.Fatalf("ParseEnvironmentConnectionsFilter(%q) returned error: %v", tt.filter, err)
			}
			if parsed.Assigned != tt.wantAssigned {
				t.Errorf("Assigned = %v, want %v", parsed.Assigned, tt.wantAssigned)
			}
			if parsed.Owner != tt.wantOwner {
				t.Errorf("Owner = %q, want %q", parsed.Owner, tt.wantOwner)
			}
		})
	}
}

// TestParseEnvironmentConnectionsFilterRejectsBadInput is the regression for
// issue #21826. Every one of these inputs previously either panicked the
// handler goroutine (the unchecked `assignedVal.(bool)`) or escaped as a raw
// encoding/json message. They must now come back as the structured
// caller-error the handler renders as 400.
func TestParseEnvironmentConnectionsFilterRejectsBadInput(t *testing.T) {
	tests := []struct {
		name        string
		filter      string
		wantMessage string
	}{
		{"assigned as string", `{"assigned":"false"}`, "must be a boolean, got a string"},
		{"assigned as number", `{"assigned":0}`, "must be a boolean, got a number"},
		{"assigned as null", `{"assigned":null}`, "must be a boolean, got null"},
		{"assigned as object", `{"assigned":{}}`, "must be a boolean, got an object"},
		{"owner as number", `{"owner":7}`, "must be a string, got a number"},
		// The ApplyFilters grammar this endpoint used to also be handed.
		{"applyfilters grammar", "owner abc", ""},
		{"not json at all", "deleted_at Deleted", ""},
		{"json array", `["assigned"]`, ""},
		// json.Unmarshal decodes a top-level "null" into a nil map without
		// error, and every subsequent map lookup on a nil map succeeds by
		// returning the zero value - so this silently fell through to the
		// default filter instead of being rejected.
		{"top-level null", `null`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The call must return, not panic.
			parsed, err := ParseEnvironmentConnectionsFilter(tt.filter)
			if err == nil {
				t.Fatalf("ParseEnvironmentConnectionsFilter(%q) = %+v, want error", tt.filter, parsed)
			}

			if code := errors.GetCode(err); code != ErrInvalidEnvironmentConnectionsFilterCode {
				t.Errorf("error code = %q, want %q", code, ErrInvalidEnvironmentConnectionsFilterCode)
			}

			// The specific detail lands in the long description, the same slot
			// ErrUnmarshal puts its underlying error in.
			if tt.wantMessage != "" && !strings.Contains(errors.GetLDescription(err), tt.wantMessage) {
				t.Errorf("long description = %q, want it to contain %q", errors.GetLDescription(err), tt.wantMessage)
			}
		})
	}
}
