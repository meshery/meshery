// Package validation defines the trust boundary for designs returned by an AI
// provider. Provider output is untrusted until it passes these checks.
package validation

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	CodeMalformed          = "malformed_output"
	CodeSchemaVersion      = "unsupported_schema_version"
	CodeMetadata           = "invalid_design_metadata"
	CodeSensitiveContent   = "sensitive_content"
	CodeRelationship       = "invalid_relationship"
	CodeDanglingReference  = "dangling_relationship_reference"
	CodeComponent          = "invalid_component"
	CurrentSchemaVersion   = "designs.meshery.io/v1beta3"
)

var acceptedSchemaVersions = map[string]struct{}{
	"designs.meshery.io/v1beta1": {},
	CurrentSchemaVersion:          {},
}

// Issue is a stable, machine-readable validation failure. Path uses JSON-style
// field notation where a location can be identified.
type Issue struct {
	Code    string `json:"code"`
	Path    string `json:"path,omitempty"`
	Message string `json:"message"`
}

func (i Issue) Error() string {
	if i.Path == "" {
		return fmt.Sprintf("%s: %s", i.Code, i.Message)
	}
	return fmt.Sprintf("%s at %s: %s", i.Code, i.Path, i.Message)
}

// ValidationError contains all deterministic issues found in one document.
// The first issue is returned by Error for compatibility with ordinary errors;
// callers that need all failures can inspect Issues.
type ValidationError struct {
	Issues []Issue
}

func (e *ValidationError) Error() string {
	if len(e.Issues) == 0 {
		return "AI design validation failed"
	}
	return e.Issues[0].Error()
}

// Unwrap allows errors.Is/As callers to identify a validation failure.
func (e *ValidationError) Unwrap() error { return errors.New("AI design validation failed") }

// Validate decodes a JSON or YAML Meshery design and applies the AI output
// contract. It never contacts a provider, registry, database, or filesystem.
func Validate(raw []byte) error {
	var err error
	if raw, err = extractDocument(raw); err != nil {
		return &ValidationError{Issues: []Issue{{Code: CodeMalformed, Message: "output must be a single JSON or YAML document, optionally wrapped in one json or yaml Markdown fence"}}}
	}
	var document map[string]interface{}
	if err := decode(raw, &document); err != nil {
		return &ValidationError{Issues: []Issue{{Code: CodeMalformed, Message: "output must be a JSON or YAML object"}}}
	}
	issues := validateDocument(document)
	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}

// extractDocument accepts a bare document or exactly one complete json/yaml
// Markdown fence. Prose, incomplete fences, and multiple fences are rejected
// before YAML parsing so provider commentary cannot be mistaken for content.
func extractDocument(raw []byte) ([]byte, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return nil, errors.New("empty document")
	}
	if !strings.Contains(trimmed, "```") {
		return []byte(trimmed), nil
	}
	lines := strings.Split(trimmed, "\n")
	if len(lines) < 3 {
		return nil, errors.New("incomplete Markdown fence")
	}
	opening := strings.TrimSpace(lines[0])
	if opening != "```json" && opening != "```yaml" {
		return nil, errors.New("Markdown fence must be labelled json or yaml")
	}
	closing := -1
	for index := 1; index < len(lines); index++ {
		if strings.TrimSpace(lines[index]) == "```" {
			closing = index
			break
		}
		if strings.Contains(lines[index], "```") {
			return nil, errors.New("multiple or malformed Markdown fences")
		}
	}
	if closing < 0 {
		return nil, errors.New("incomplete Markdown fence")
	}
	for _, line := range lines[closing+1:] {
		if strings.TrimSpace(line) != "" {
			return nil, errors.New("prose surrounding Markdown fence is not allowed")
		}
	}
	content := strings.TrimSpace(strings.Join(lines[1:closing], "\n"))
	if content == "" {
		return nil, errors.New("empty Markdown fence")
	}
	return []byte(content), nil
}

// decode parses one YAML/JSON document and rejects a second document or a
// malformed document separator. The caller has already removed an optional
// Markdown fence.
func decode(raw []byte, destination *map[string]interface{}) error {
	decoder := yaml.NewDecoder(bytes.NewReader(raw))
	decoder.KnownFields(false)
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if *destination == nil {
		return errors.New("empty document")
	}
	var additional interface{}
	if err := decoder.Decode(&additional); err != io.EOF {
		if err == nil {
			return errors.New("multiple YAML/JSON documents are not allowed")
		}
		return err
	}
	return nil
}

// validateDocument applies top-level metadata, component, secret, and
// relationship rules in a stable order.
func validateDocument(document map[string]interface{}) []Issue {
	issues := make([]Issue, 0, 4)
	version, ok := stringValue(document["schemaVersion"])
	if !ok || version == "" {
		issues = append(issues, Issue{Code: CodeSchemaVersion, Path: "schemaVersion", Message: "schemaVersion is required"})
	} else if _, accepted := acceptedSchemaVersions[version]; !accepted {
		issues = append(issues, Issue{Code: CodeSchemaVersion, Path: "schemaVersion", Message: "unsupported schemaVersion; accepted versions are designs.meshery.io/v1beta1 and designs.meshery.io/v1beta3"})
	}

	name, ok := stringValue(document["name"])
	if !ok || strings.TrimSpace(name) == "" {
		issues = append(issues, Issue{Code: CodeMetadata, Path: "name", Message: "design name is required"})
	}
	components, ok := sliceValue(document["components"])
	if !ok {
		issues = append(issues, Issue{CodeMetadata, "components", "components must be an array"})
	} else if len(components) == 0 {
		issues = append(issues, Issue{CodeMetadata, "components", "at least one component is required"})
	}
	for index, component := range components {
		issues = append(issues, validateComponent(component, index)...)
	}

	issues = append(issues, sensitiveIssues(document, "")...)
	if ok {
		issues = append(issues, relationshipIssues(document, components)...)
	}
	return issues
}

// Key matching is intentionally conservative. AI output is rejected when a
// camelCase or delimiter-separated secret field is present, even if its value
// is only a placeholder.
var sensitiveKey = regexp.MustCompile(`(?i)(authorization|api[-_.]?(?:key|token)|access[-_.]?(?:key|token)|bearer|token|password|secret|private[-_.]?key|client[-_.]?key[-_.]?data|kubeconfig|credential)`)
var sensitiveValue = regexp.MustCompile(`(?is)(bearer\s+[A-Za-z0-9._~+/=-]{12,}|-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----|client-key-data\s*[:=]|apiVersion:\s*v1\s*\n[\s\S]*kind:\s*Config)`)

func sensitiveIssues(value interface{}, path string) []Issue {
	issues := []Issue{}
	switch typed := value.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			childPath := key
			if path != "" {
				childPath = path + "." + key
			}
			if sensitiveKey.MatchString(key) {
				issues = append(issues, Issue{Code: CodeSensitiveContent, Path: redactedPath(path), Message: "credential or secret-bearing fields are not allowed in AI-generated designs"})
				continue
			}
			issues = append(issues, sensitiveIssues(typed[key], childPath)...)
		}
	case []interface{}:
		for index, item := range typed {
			issues = append(issues, sensitiveIssues(item, fmt.Sprintf("%s[%d]", path, index))...)
		}
	case string:
		if sensitiveValue.MatchString(typed) {
			issues = append(issues, Issue{Code: CodeSensitiveContent, Path: path, Message: "credential, token, private key, or kubeconfig content is not allowed"})
		}
	}
	return issues
}

// validateComponent checks the minimum portable Meshery component envelope.
// Registry-backed schema and relationship evaluation remain downstream.
func validateComponent(value interface{}, index int) []Issue {
	path := fmt.Sprintf("components[%d]", index)
	component, ok := value.(map[string]interface{})
	if !ok {
		return []Issue{{Code: CodeComponent, Path: path, Message: "component entry must be an object"}}
	}
	issues := []Issue{}
	definition, ok := component["component"].(map[string]interface{})
	if !ok {
		issues = append(issues, Issue{Code: CodeComponent, Path: path + ".component", Message: "component definition is required"})
	} else {
		for _, field := range []string{"kind", "version"} {
			if value, exists := definition[field].(string); !exists || strings.TrimSpace(value) == "" {
				issues = append(issues, Issue{Code: CodeComponent, Path: path + ".component." + field, Message: "component kind and version are required"})
			}
		}
	}
	model, ok := component["model"].(map[string]interface{})
	if !ok {
		issues = append(issues, Issue{Code: CodeComponent, Path: path + ".model", Message: "model definition is required"})
	} else if name, exists := model["name"].(string); !exists || strings.TrimSpace(name) == "" {
		issues = append(issues, Issue{Code: CodeComponent, Path: path + ".model.name", Message: "model name is required"})
	}
	if configuration, exists := component["configuration"]; !exists || configuration == nil {
		issues = append(issues, Issue{Code: CodeComponent, Path: path + ".configuration", Message: "configuration is required"})
	}
	return issues
}

// redactedPath removes a sensitive field name from diagnostics while retaining
// enough location context for callers to identify the rejected object.
func redactedPath(parent string) string {
	if parent == "" {
		return "<redacted>"
	}
	return parent + ".<redacted>"
}

// relationshipIssues validates instance-edge endpoints. Selector-only
// relationship definitions are intentionally deferred to registry validation.
func relationshipIssues(document map[string]interface{}, components []interface{}) []Issue {
	relationships, ok := sliceValue(document["relationships"])
	if !ok {
		return nil
	}
	ids := map[string]struct{}{}
	for index, item := range components {
		component, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		id, _ := stringValue(component["id"])
		if id == "" {
			if nested, nestedOK := component["component"].(map[string]interface{}); nestedOK {
				id, _ = stringValue(nested["id"])
			}
		}
		if id != "" {
			ids[id] = struct{}{}
		} else {
			_ = index // component identity is optional in legacy declarations
		}
	}
	issues := []Issue{}
	for index, item := range relationships {
		relationship, ok := item.(map[string]interface{})
		if !ok {
			issues = append(issues, Issue{Code: CodeRelationship, Path: fmt.Sprintf("relationships[%d]", index), Message: "relationship must be an object"})
			continue
		}
		from, hasFrom, validFrom := endpoint(relationship, "source", "from", "origin")
		to, hasTo, validTo := endpoint(relationship, "target", "to", "destination")
		if !hasFrom && !hasTo {
			continue // schema relationship definitions may contain selectors, not instance edges
		}
		path := fmt.Sprintf("relationships[%d]", index)
		if !hasFrom || !hasTo || !validFrom || !validTo {
			issues = append(issues, Issue{Code: CodeRelationship, Path: path, Message: "relationship edges must define both source and target"})
			continue
		}
		if _, exists := ids[from]; !exists {
			issues = append(issues, Issue{Code: CodeDanglingReference, Path: path + ".source", Message: "source component is not declared in components"})
		}
		if _, exists := ids[to]; !exists {
			issues = append(issues, Issue{Code: CodeDanglingReference, Path: path + ".target", Message: "target component is not declared in components"})
		}
	}
	return issues
}

// endpoint distinguishes an absent endpoint alias from an explicitly empty or
// incorrectly typed alias, which must be rejected as a malformed edge.
func endpoint(relationship map[string]interface{}, names ...string) (string, bool, bool) {
	present := false
	valid := true
	resolved := ""
	for _, name := range names {
		value, exists := relationship[name]
		if !exists {
			continue
		}
		present = true
		if text, ok := stringValue(value); ok && text != "" {
			if resolved == "" {
				resolved = text
			}
			continue
		}
		if object, ok := value.(map[string]interface{}); ok {
			found := false
			for _, key := range []string{"id", "componentId", "componentID"} {
				if text, ok := stringValue(object[key]); ok && text != "" {
					if resolved == "" {
						resolved = text
					}
					found = true
					break
				}
			}
			if !found {
				valid = false
			}
			continue
		}
		valid = false
	}
	return resolved, present, valid && resolved != ""
}

// stringValue returns a string without coercing arbitrary YAML values.
func stringValue(value interface{}) (string, bool) {
	text, ok := value.(string)
	return text, ok
}

// sliceValue returns an array without coercing arbitrary YAML values.
func sliceValue(value interface{}) ([]interface{}, bool) {
	slice, ok := value.([]interface{})
	return slice, ok
}

// IsValidationError reports whether err was produced by Validate.
func IsValidationError(err error) bool {
	var target *ValidationError
	return errors.As(err, &target)
}
