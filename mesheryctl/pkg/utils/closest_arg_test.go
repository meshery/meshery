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
		values := []int{1, 2, 3, 4, 5}
		minIndex, _ := minIntSlice(values)
		if minIndex != 0 {
			t.Errorf("Expected 0, got %d", minIndex)
		}
	})

	t.Run("Test minimum", func(t *testing.T) {
		min := minimum(1, 2, 3)
		if min != 1 {
			t.Errorf("Expected 1, got %d", min)
		}

		min = minimum(2, 1, 3)
		if min != 1 {
			t.Errorf("Expected 1, got %d", min)
		}

		min = minimum(3, 2, 1)
		if min != 1 {
			t.Errorf("Expected 1, got %d", min)
		}
		min = minimum(1, 1, 1)

		if min != 1 {
			t.Errorf("Expected 1, got %d", min)
		}
	})

	t.Run("Test levenshtein", func(t *testing.T) {
		str1 := []rune("meshery")
		str2 := []rune("mesheryctl")

		distance := levenshtein(str1, str2)
		if distance != 3 {
			t.Errorf("Expected 3, got %d", distance)
		}
	})

	t.Run("Test FindClosestArg", func(t *testing.T) {
		argPassed := "meshery"
		validArgs := []string{"meshery", "mesheryctl", "meshery1", "meshery2"}
		closestArg := FindClosestArg(argPassed, validArgs)
		if closestArg != "meshery" {
			t.Errorf("Expected meshery, got %s", closestArg)
		}
	})
}
