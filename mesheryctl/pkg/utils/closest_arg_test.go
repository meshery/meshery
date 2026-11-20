// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"testing"
)

// Test cases for closest_arg.go
func TestClosestArgs(t *testing.T) {
	t.Run("Test minIntslice", func(t *testing.T) {
		tests := []struct {
			name          string
			values        []int
			expectedIndex int
			expectedMin   int
		}{
			{
				name:          "ordered ascending",
				values:        []int{1, 2, 3, 4, 5},
				expectedIndex: 0,
				expectedMin:   1,
			},
			{
				name:          "minimum in middle",
				values:        []int{5, 4, 1, 2, 3},
				expectedIndex: 2,
				expectedMin:   1,
			},
			{
				name:          "minimum at end",
				values:        []int{9, 8, 7, 6, 2},
				expectedIndex: 4,
				expectedMin:   2,
			},
			{
				name:          "duplicate minimums preserves first candidate index",
				values:        []int{4, 2, 2, 5, 2},
				expectedIndex: 1,
				expectedMin:   2,
			},
			{
				name:          "single element",
				values:        []int{42},
				expectedIndex: 0,
				expectedMin:   42,
			},
			{
				name:          "negative values",
				values:        []int{-2, 0, -10, 5},
				expectedIndex: 2,
				expectedMin:   -10,
			},
			{
				name:          "empty slice safe boundary",
				values:        []int{},
				expectedIndex: -1,
				expectedMin:   -1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				idx, val := minIntSlice(tt.values)
				if idx != tt.expectedIndex || val != tt.expectedMin {
					t.Errorf("minIntSlice(%v) = (%d, %d), expected (%d, %d)",
						tt.values, idx, val, tt.expectedIndex, tt.expectedMin)
				}
			})
		}
	})

	t.Run("Test minimum", func(t *testing.T) {
		tests := []struct {
			a, b, c  int
			expected int
		}{
			{1, 2, 3, 1},
			{2, 1, 3, 1},
			{3, 2, 1, 1},
			{1, 1, 1, 1},
			{1, 1, 2, 1},
			{2, 1, 1, 1},
			{5, 0, 10, 0},
			{-1, -2, -3, -3},
			{0, 0, 0, 0},
		}

		for _, tt := range tests {
			min := minimum(tt.a, tt.b, tt.c)
			if min != tt.expected {
				t.Errorf("minimum(%d, %d, %d) = %d, expected %d", tt.a, tt.b, tt.c, min, tt.expected)
			}
		}
	})

	t.Run("Test levenshtein", func(t *testing.T) {
		tests := []struct {
			name     string
			str1     string
			str2     string
			expected int
		}{
			{
				name:     "identical strings",
				str1:     "meshery",
				str2:     "meshery",
				expected: 0,
			},
			{
				name:     "both empty strings",
				str1:     "",
				str2:     "",
				expected: 0,
			},
			{
				name:     "first empty second non-empty",
				str1:     "",
				str2:     "meshery",
				expected: 7,
			},
			{
				name:     "first non-empty second empty",
				str1:     "meshery",
				str2:     "",
				expected: 7,
			},
			{
				name:     "single insertion",
				str1:     "mesher",
				str2:     "meshery",
				expected: 1,
			},
			{
				name:     "single deletion",
				str1:     "mesheryctl",
				str2:     "meshery",
				expected: 3,
			},
			{
				name:     "single substitution",
				str1:     "meshary",
				str2:     "meshery",
				expected: 1,
			},
			{
				name:     "two substitutions",
				str1:     "moshory",
				str2:     "meshery",
				expected: 2,
			},
			{
				name:     "adjacent character transposition in standard Levenshtein",
				str1:     "mseh",
				str2:     "mesh",
				expected: 2,
			},
			{
				name:     "completely disjoint strings",
				str1:     "abc",
				str2:     "xyz",
				expected: 3,
			},
			{
				name:     "case sensitivity difference",
				str1:     "Meshery",
				str2:     "meshery",
				expected: 1,
			},
			{
				name:     "rune-aware unicode accents",
				str1:     "mëshëry",
				str2:     "meshery",
				expected: 2,
			},
			{
				name:     "rune-aware CJK characters single deletion",
				str1:     "日本語",
				str2:     "日本",
				expected: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				dist := levenshtein([]rune(tt.str1), []rune(tt.str2))
				if dist != tt.expected {
					t.Errorf("levenshtein(%q, %q) = %d, expected %d", tt.str1, tt.str2, dist, tt.expected)
				}
			})
		}
	})

	t.Run("Test FindClosestArg", func(t *testing.T) {
		tests := []struct {
			name      string
			argPassed string
			validArgs []string
			expected  string
		}{
			{
				name:      "exact match in first position",
				argPassed: "meshery",
				validArgs: []string{"meshery", "mesheryctl", "meshery1", "meshery2"},
				expected:  "meshery",
			},
			{
				name:      "exact match in later position",
				argPassed: "system",
				validArgs: []string{"pattern", "perf", "system", "filter"},
				expected:  "system",
			},
			{
				name:      "misspelling with transposition (sysetm -> system)",
				argPassed: "sysetm",
				validArgs: []string{"pattern", "perf", "system", "filter"},
				expected:  "system",
			},
			{
				name:      "misspelling with transposition (fliter -> filter)",
				argPassed: "fliter",
				validArgs: []string{"filter", "pattern", "system"},
				expected:  "filter",
			},
			{
				name:      "misspelling with dropped character (meshry -> meshery)",
				argPassed: "meshry",
				validArgs: []string{"meshery", "adapter", "pattern"},
				expected:  "meshery",
			},
			{
				name:      "misspelling with substitution (deploi -> deploy)",
				argPassed: "deploi",
				validArgs: []string{"deploy", "delete", "view", "apply"},
				expected:  "deploy",
			},
			{
				name:      "misspelling with truncated suffix (delet -> delete)",
				argPassed: "delet",
				validArgs: []string{"deploy", "delete", "view", "apply"},
				expected:  "delete",
			},
			{
				name:      "misspelling with dropped duplicate letter (aply -> apply)",
				argPassed: "aply",
				validArgs: []string{"apply", "delete", "view", "deploy"},
				expected:  "apply",
			},
			{
				name:      "insertion typo (patterrn -> pattern)",
				argPassed: "patterrn",
				validArgs: []string{"pattern", "filter", "system"},
				expected:  "pattern",
			},
			{
				name:      "deletion typo (per -> perf)",
				argPassed: "per",
				validArgs: []string{"perf", "pattern", "mesh"},
				expected:  "perf",
			},
			{
				name:      "tie-breaker selects first candidate with equal distance",
				argPassed: "cat",
				validArgs: []string{"bat", "hat", "rat"},
				expected:  "bat",
			},
			{
				name:      "empty input string selects shortest validArg",
				argPassed: "",
				validArgs: []string{"longcommand", "mid", "a"},
				expected:  "a",
			},
			{
				name:      "empty validArgs returns empty string without panic",
				argPassed: "test",
				validArgs: []string{},
				expected:  "",
			},
			{
				name:      "unicode accented typo match",
				argPassed: "mëshëry",
				validArgs: []string{"meshery", "system", "pattern"},
				expected:  "meshery",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				closest := FindClosestArg(tt.argPassed, tt.validArgs)
				if closest != tt.expected {
					t.Errorf("FindClosestArg(%q, %v) = %q, expected %q",
						tt.argPassed, tt.validArgs, closest, tt.expected)
				}
			})
		}
	})
}
