package models

import (
	"strings"
	"testing"

	"github.com/meshery/meshkit/errors"
)

func TestParseWorkspaceFilterDefaults(t *testing.T) {
	for _, filter := range []string{"", "   "} {
		parsed, err := ParseWorkspaceFilter(filter)
		if err != nil {
			t.Fatalf("ParseWorkspaceFilter(%q) returned error: %v", filter, err)
		}
		if !parsed.Assigned {
			t.Errorf("ParseWorkspaceFilter(%q).Assigned = false, want true", filter)
		}
		if parsed.DeletedAt {
			t.Errorf("ParseWorkspaceFilter(%q).DeletedAt = true, want false", filter)
		}
		if parsed.Owner != "" {
			t.Errorf("ParseWorkspaceFilter(%q).Owner = %q, want empty", filter, parsed.Owner)
		}
	}
}

func TestParseWorkspaceFilterValid(t *testing.T) {
	tests := []struct {
		name   string
		filter string
		want   WorkspaceFilter
	}{
		{"assigned false", `{"assigned":false}`, WorkspaceFilter{Assigned: false}},
		{"assigned true", `{"assigned":true}`, WorkspaceFilter{Assigned: true}},
		{"deletedAt true", `{"deletedAt":true}`, WorkspaceFilter{Assigned: true, DeletedAt: true}},
		{"owner only", `{"owner":"abc"}`, WorkspaceFilter{Assigned: true, Owner: "abc"}},
		{"all three", `{"assigned":false,"deletedAt":true,"owner":"abc"}`, WorkspaceFilter{Assigned: false, DeletedAt: true, Owner: "abc"}},
		{"unknown field ignored", `{"nope":1}`, WorkspaceFilter{Assigned: true}},
		{"empty object", `{}`, WorkspaceFilter{Assigned: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseWorkspaceFilter(tt.filter)
			if err != nil {
				t.Fatalf("ParseWorkspaceFilter(%q) returned error: %v", tt.filter, err)
			}
			if got != tt.want {
				t.Errorf("ParseWorkspaceFilter(%q) = %+v, want %+v", tt.filter, got, tt.want)
			}
		})
	}
}

// TestParseWorkspaceFilterRejectsBadInput is the regression for the
// two-grammar defect. Every one of these inputs previously either escaped as a
// raw encoding/json message or was silently accepted as the default filter.
func TestParseWorkspaceFilterRejectsBadInput(t *testing.T) {
	tests := []struct {
		name        string
		filter      string
		wantMessage string
	}{
		{"assigned as string", `{"assigned":"false"}`, "must be a boolean, got a string"},
		{"assigned as number", `{"assigned":0}`, "must be a boolean, got a number"},
		{"assigned as null", `{"assigned":null}`, "must be a boolean, got null"},
		{"deletedAt as string", `{"deletedAt":"true"}`, "must be a boolean, got a string"},
		{"owner as number", `{"owner":7}`, "must be a string, got a number"},
		{"owner as object", `{"owner":{}}`, "must be a string, got an object"},
		// A top-level null decodes into a nil map without error; every lookup
		// on it succeeds with the zero value, so it would otherwise fall
		// through to the default filter.
		{"top-level null", `null`, "must be a JSON object, got null"},
		// The ApplyFilters grammar these endpoints used to also be handed.
		{"applyfilters grammar", "owner abc", ""},
		{"not json at all", "deleted_at Deleted", ""},
		{"json array", `["assigned"]`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseWorkspaceFilter(tt.filter)
			if err == nil {
				t.Fatalf("ParseWorkspaceFilter(%q) = %+v, want error", tt.filter, parsed)
			}

			if code := errors.GetCode(err); code != ErrInvalidWorkspaceFilterCode {
				t.Errorf("error code = %q, want %q", code, ErrInvalidWorkspaceFilterCode)
			}

			if tt.wantMessage != "" && !strings.Contains(errors.GetLDescription(err), tt.wantMessage) {
				t.Errorf("long description = %q, want it to contain %q", errors.GetLDescription(err), tt.wantMessage)
			}
		})
	}
}
