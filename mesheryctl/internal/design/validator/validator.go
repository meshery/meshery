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

package validator

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Pattern represents the minimal fields required to identify a valid Meshery design
type Pattern struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

// ValidateRequiredFields performs lightweight validation on a design file
// and returns a list of validation error messages.
func ValidateRequiredFields(path string) []string {
	var errs []string

	data, err := os.ReadFile(path)
	if err != nil {
		return []string{err.Error()}
	}

	var p Pattern
	if err := yaml.Unmarshal(data, &p); err != nil {
		return []string{"invalid YAML format: " + err.Error()}
	}

	if p.APIVersion == "" {
		errs = append(errs, "apiVersion is missing")
	}

	if p.Kind == "" {
		errs = append(errs, "kind is missing")
	}

	return errs
}
