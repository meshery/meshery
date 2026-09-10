package validation

import (
	"errors"
	"fmt"
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
	validYAML := []byte("schemaVersion: designs.meshery.io/v1beta3\nname: yaml-design\ncomponents:\n  - id: one\n    component:\n      kind: Deployment\n      version: apps/v1\n    model:\n      name: kubernetes\n    configuration: {}\n")
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

func TestValidateMarkdownFencesAndDocuments(t *testing.T) {
	validJSON := "```json\n{" +
		"\"schemaVersion\":\"designs.meshery.io/v1beta3\",\"name\":\"fenced\",\"components\":[{" +
		"\"id\":\"one\",\"component\":{\"kind\":\"Deployment\",\"version\":\"apps/v1\"}," +
		"\"model\":{\"name\":\"kubernetes\"},\"configuration\":{}}]}\n```"
	for name, input := range map[string]string{
		"json fence": validJSON,
		"yaml fence": "```yaml\nschemaVersion: designs.meshery.io/v1beta3\nname: fenced\ncomponents:\n  - component:\n      kind: Deployment\n      version: apps/v1\n    model:\n      name: kubernetes\n    configuration: {}\n```",
	} {
		t.Run(name, func(t *testing.T) {
			if err := Validate([]byte(input)); err != nil {
				t.Fatalf("accepted fence rejected: %v", err)
			}
		})
	}
	for name, input := range map[string]string{
		"surrounding prose": "Here is the design:\n" + validJSON,
		"incomplete fence":  "```json\n{}",
		"multiple fences":   validJSON + "\n\n```json\n{}\n```",
		"second document":   validJSON + "\n---\napiToken: do-not-return-this\n",
	} {
		t.Run(name, func(t *testing.T) {
			if err := Validate([]byte(input)); err == nil {
				t.Fatal("invalid fenced or multi-document input was accepted")
			}
		})
	}
}

func TestValidateFixtureEdgeCases(t *testing.T) {
	tests := []struct {
		file string
		ok   bool
	}{
		{file: "fenced_valid_design.md", ok: true},
		{file: "fenced_valid_design.yaml.md", ok: true},
		{file: "fenced_surrounding_prose.md"},
		{file: "fenced_incomplete.md"},
		{file: "fenced_multiple.md"},
		{file: "second_document.yaml"},
		{file: "non_object_component.json"},
		{file: "incomplete_component.json"},
		{file: "empty_relationship_endpoint.json"},
		{file: "wrong_type_relationship_endpoint.json"},
	}
	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("fixtures", test.file))
			if err != nil {
				t.Fatal(err)
			}
			err = Validate(raw)
			if test.ok && err != nil {
				t.Fatalf("valid fixture rejected: %v", err)
			}
			if !test.ok && err == nil {
				t.Fatal("invalid fixture was accepted")
			}
		})
	}
}

func TestValidateComponentEnvelope(t *testing.T) {
	base := `{"schemaVersion":"designs.meshery.io/v1beta3","name":"x","components":%s}`
	for name, components := range map[string]string{
		"non-object": `["not-a-component"]`,
		"missing component": `[{"model":{"name":"kubernetes"},"configuration":{}}]`,
		"missing model": `[{"component":{"kind":"Deployment","version":"apps/v1"},"configuration":{}}]`,
		"missing kind": `[{"component":{"version":"apps/v1"},"model":{"name":"kubernetes"},"configuration":{}}]`,
		"missing version": `[{"component":{"kind":"Deployment"},"model":{"name":"kubernetes"},"configuration":{}}]`,
		"missing configuration": `[{"component":{"kind":"Deployment","version":"apps/v1"},"model":{"name":"kubernetes"}}]`,
	} {
		t.Run(name, func(t *testing.T) {
			err := Validate([]byte(fmt.Sprintf(base, components)))
			var validationErr *ValidationError
			if !errors.As(err, &validationErr) || !hasIssue(validationErr, CodeComponent) {
				t.Fatalf("expected component issue, got %v", err)
			}
		})
	}
}

func TestValidateDoesNotExposeSecretValues(t *testing.T) {
	err := Validate([]byte(`{"schemaVersion":"designs.meshery.io/v1beta3","name":"x","components":[{"id":"x","component":{"kind":"Deployment","version":"apps/v1"},"model":{"name":"kubernetes"},"configuration":{"apiToken":"do-not-return-this"}}]}`))
	if err == nil || strings.Contains(err.Error(), "do-not-return-this") || strings.Contains(err.Error(), "apiToken") {
		t.Fatalf("secret was accepted or exposed: %v", err)
	}
}

func TestValidateRelationshipEndpointTypes(t *testing.T) {
	base := `{"schemaVersion":"designs.meshery.io/v1beta3","name":"x","components":[{"id":"one","component":{"kind":"Deployment","version":"apps/v1"},"model":{"name":"kubernetes"},"configuration":{}}],"relationships":[%s]}`
	for name, relationship := range map[string]string{
		"empty source": `{"source":"","target":"one"}`,
		"empty target": `{"source":"one","target":""}`,
		"wrong source type": `{"source":false,"target":"one"}`,
		"wrong target type": `{"source":"one","target":[]}`,
	} {
		t.Run(name, func(t *testing.T) {
			err := Validate([]byte(fmt.Sprintf(base, relationship)))
			var validationErr *ValidationError
			if !errors.As(err, &validationErr) || !hasIssue(validationErr, CodeRelationship) {
				t.Fatalf("expected relationship issue, got %v", err)
			}
		})
	}
}

func hasIssue(validationErr *ValidationError, code string) bool {
	for _, issue := range validationErr.Issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}
