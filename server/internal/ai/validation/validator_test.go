package validation

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateFixtures(t *testing.T) {
	tests := []struct {
		name string
		file string
		code string
	}{
		{name: "valid", file: "valid_ai_design.json"},
		{name: "missing schema", file: "missing_schema_design.json", code: CodeSchemaVersion},
		{name: "sensitive token", file: "sensitive_token_leak_design.json", code: CodeSensitiveContent},
		{name: "dangling relationship", file: "dangling_relationship_design.json", code: CodeDanglingReference},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("fixtures", test.file))
			if err != nil {
				t.Fatal(err)
			}
			err = Validate(raw)
			if test.code == "" {
				if err != nil {
					t.Fatalf("valid fixture rejected: %v", err)
				}
				return
			}
			var validationErr *ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected ValidationError, got %v", err)
			}
			found := false
			for _, issue := range validationErr.Issues {
				if issue.Code == test.code {
					found = true
				}
				if strings.Contains(issue.Message, "Bearer") || strings.Contains(issue.Message, "abcdefghijklmnop") {
					t.Fatal("validation message leaked sensitive content")
				}
			}
			if !found {
				t.Fatalf("expected issue code %q, got %#v", test.code, validationErr.Issues)
			}
		})
	}
}

func TestValidateYAMLAndMalformedOutput(t *testing.T) {
	validYAML := []byte("schemaVersion: designs.meshery.io/v1beta3\nname: yaml-design\ncomponents:\n  - id: one\n    component:\n      kind: Deployment\n")
	if err := Validate(validYAML); err != nil {
		t.Fatalf("valid YAML rejected: %v", err)
	}
	if err := Validate([]byte("[not: a design")); err == nil {
		t.Fatal("malformed output was accepted")
	} else {
		var validationErr *ValidationError
		if !errors.As(err, &validationErr) || validationErr.Issues[0].Code != CodeMalformed {
			t.Fatalf("unexpected malformed error: %v", err)
		}
	}
}

func TestValidateDoesNotExposeSecretValues(t *testing.T) {
	err := Validate([]byte(`{"schemaVersion":"designs.meshery.io/v1beta3","name":"x","components":[{"id":"x","configuration":{"apiToken":"do-not-return-this"}}]}`))
	if err == nil || strings.Contains(err.Error(), "do-not-return-this") {
		t.Fatalf("secret was accepted or exposed: %v", err)
	}
}
