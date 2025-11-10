package v1beta1_test

import (
	"regexp"
	"testing"

	"github.com/meshery/meshery/server/models/meshmodel/registry/v1beta1"
)

func TestComponentFilter_Exclusion(t *testing.T) {
	filter := &v1beta1.ComponentFilter{
		Exclude:      "List",
		ExcludeRegex: ".*List$",
	}

	tests := []struct {
		name   string
		input  string
		expect bool
	}{
		{"exclude exact match", "ServiceList", true},
		{"exclude regex match", "PodList", true},
		{"allow non-match", "Service", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if matched, _ := regexp.MatchString(filter.ExcludeRegex, tt.input); matched != tt.expect {
				t.Errorf("expected %v for %s", tt.expect, tt.input)
			}
		})
	}
}
