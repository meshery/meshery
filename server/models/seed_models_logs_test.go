package models

import (
	"reflect"
	"testing"

	"github.com/meshery/meshkit/models/meshmodel/core/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/connection"
)

func host(kind string, models, components, relationships, policies int64) v1beta1.MeshModelHostsWithEntitySummary {
	return v1beta1.MeshModelHostsWithEntitySummary{
		Connection: connection.Connection{Kind: kind},
		Summary: v1beta1.EntitySummary{
			Models:        models,
			Components:    components,
			Relationships: relationships,
			Policies:      policies,
		},
	}
}

// TestAggregateSummariesByKind guards the de-duplication of registrant startup
// summaries: each registrant Kind must be reported exactly once with its counts
// folded together, regardless of how many connections back that Kind.
func TestAggregateSummariesByKind(t *testing.T) {
	tests := []struct {
		name          string
		hosts         []v1beta1.MeshModelHostsWithEntitySummary
		wantOrder     []string
		wantSummaries map[string]v1beta1.EntitySummary
	}{
		{
			name:          "empty input",
			hosts:         nil,
			wantOrder:     []string{},
			wantSummaries: map[string]v1beta1.EntitySummary{},
		},
		{
			name: "distinct kinds are preserved in first-seen order",
			hosts: []v1beta1.MeshModelHostsWithEntitySummary{
				host("artifacthub", 150, 1326, 26, 0),
				host("github", 129, 1441, 559, 0),
				host("meshery", 9, 1158, 4, 0),
			},
			wantOrder: []string{"artifacthub", "github", "meshery"},
			wantSummaries: map[string]v1beta1.EntitySummary{
				"artifacthub": {Models: 150, Components: 1326, Relationships: 26},
				"github":      {Models: 129, Components: 1441, Relationships: 559},
				"meshery":     {Models: 9, Components: 1158, Relationships: 4},
			},
		},
		{
			// Reproduces issue #19588: two artifacthub connections must collapse
			// into a single summary line (152 = 150+2, 1336 = 1326+10).
			name: "duplicate kinds are folded into one summary",
			hosts: []v1beta1.MeshModelHostsWithEntitySummary{
				host("artifacthub", 150, 1326, 26, 0),
				host("github", 129, 1441, 559, 0),
				host("artifacthub", 2, 10, 0, 0),
				host("meshery", 9, 1158, 4, 1),
			},
			wantOrder: []string{"artifacthub", "github", "meshery"},
			wantSummaries: map[string]v1beta1.EntitySummary{
				"artifacthub": {Models: 152, Components: 1336, Relationships: 26},
				"github":      {Models: 129, Components: 1441, Relationships: 559},
				"meshery":     {Models: 9, Components: 1158, Relationships: 4, Policies: 1},
			},
		},
		{
			// An empty Kind cannot be named in the summary, so it must be skipped
			// rather than produce a "For registrant  imported" line.
			name: "empty kind is skipped",
			hosts: []v1beta1.MeshModelHostsWithEntitySummary{
				host("", 5, 5, 0, 0),
				host("github", 129, 1441, 559, 0),
			},
			wantOrder: []string{"github"},
			wantSummaries: map[string]v1beta1.EntitySummary{
				"github": {Models: 129, Components: 1441, Relationships: 559},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOrder, gotSummaries := aggregateSummariesByKind(tt.hosts)

			if !reflect.DeepEqual(gotOrder, tt.wantOrder) {
				t.Errorf("ordered kinds = %v, want %v", gotOrder, tt.wantOrder)
			}
			// Cardinality must equal the number of distinct kinds, not the number
			// of input hosts - this is the regression the de-duplication fixes.
			if len(gotOrder) != len(tt.wantOrder) {
				t.Errorf("summary line count = %d, want %d", len(gotOrder), len(tt.wantOrder))
			}
			if !reflect.DeepEqual(gotSummaries, tt.wantSummaries) {
				t.Errorf("summaries = %#v, want %#v", gotSummaries, tt.wantSummaries)
			}
		})
	}
}

// TestFormatRegistrantSummary covers the one-line summary rendering, including
// the defensive zero-entities fallback so the message is never left dangling.
func TestFormatRegistrantSummary(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		summary v1beta1.EntitySummary
		want    string
	}{
		{
			name:    "all entity types present",
			kind:    "artifacthub",
			summary: v1beta1.EntitySummary{Models: 150, Components: 1326, Relationships: 26, Policies: 3},
			want:    "For registrant artifacthub imported 150 models, 1326 components, 26 relationships, 3 policies.",
		},
		{
			name:    "only some entity types present",
			kind:    "meshery",
			summary: v1beta1.EntitySummary{Models: 9, Components: 1158, Relationships: 4},
			want:    "For registrant meshery imported 9 models, 1158 components, 4 relationships.",
		},
		{
			name:    "no entities falls back to 0 entities",
			kind:    "github",
			summary: v1beta1.EntitySummary{},
			want:    "For registrant github imported 0 entities.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatRegistrantSummary(tt.kind, tt.summary); got != tt.want {
				t.Errorf("formatRegistrantSummary() = %q, want %q", got, tt.want)
			}
		})
	}
}
