// Package validation provides output validation and re-prompt logic for
// AI-generated Meshery designs.
package validation

import (
	"encoding/json"
	"fmt"
	"strings"

	patternV1beta1 "github.com/meshery/schemas/models/v1beta1/pattern"
)

// ValidationResult holds the outcome of validating an AI-generated output.
type ValidationResult struct {
	Valid   bool
	Design  *patternV1beta1.PatternFile
	Errors  []string
	RawJSON string
}

// Validator validates AI-generated JSON against the Meshery PatternFile schema.
type Validator struct{}

// NewValidator creates a new Validator.
func NewValidator() *Validator {
	return &Validator{}
}

// Validate attempts to parse and validate the given string as a Meshery
// PatternFile (v1beta1). It returns a ValidationResult describing whether
// the output is usable and any errors found.
func (v *Validator) Validate(output string) ValidationResult {
	result := ValidationResult{RawJSON: output}

	// Strip any markdown fencing that the LLM may have added despite the prompt.
	cleaned := stripMarkdownFences(output)

	var pf patternV1beta1.PatternFile
	if err := json.Unmarshal([]byte(cleaned), &pf); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("invalid JSON: %v", err))
		return result
	}

	result.Design = &pf
	result.Errors = validatePatternFile(&pf)
	result.Valid = len(result.Errors) == 0

	return result
}

// validatePatternFile performs semantic checks on the parsed PatternFile.
func validatePatternFile(pf *patternV1beta1.PatternFile) []string {
	var errs []string

	if pf.Name == "" {
		errs = append(errs, "design 'name' is required")
	}

	if pf.SchemaVersion == "" {
		errs = append(errs, "design 'schemaVersion' is required (expected 'designs.meshery.io/v1beta1')")
	}

	if len(pf.Components) == 0 {
		errs = append(errs, "design must contain at least one component")
	}

	for i, comp := range pf.Components {
		prefix := fmt.Sprintf("component[%d]", i)
		if comp == nil {
			errs = append(errs, fmt.Sprintf("%s: component is nil", prefix))
			continue
		}

		if comp.Component.Kind == "" {
			errs = append(errs, fmt.Sprintf("%s: 'component.kind' is required", prefix))
		}

		if comp.DisplayName == "" {
			errs = append(errs, fmt.Sprintf("%s: 'displayName' is required", prefix))
		}
	}

	return errs
}

// stripMarkdownFences removes ```json ... ``` or ``` ... ``` wrappers that
// LLMs sometimes add around their JSON output.
func stripMarkdownFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		// Remove opening fence line (```json or ```)
		idx := strings.Index(s, "\n")
		if idx != -1 {
			s = s[idx+1:]
		} else {
			// Single-line fence: strip fence and common language tags
			s = strings.TrimPrefix(s, "```")
			s = strings.TrimPrefix(s, "json")
			s = strings.TrimPrefix(s, "yaml")
		}
	}
	if strings.HasSuffix(s, "```") {
		s = s[:strings.LastIndex(s, "```")]
	}
	return strings.TrimSpace(s)
}
