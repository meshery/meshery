package policies

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// objectGetNested traverses a nested map following the given path segments.
// Returns defaultValue if any segment is missing or the path is invalid.
func objectGetNested(obj interface{}, path []string, defaultValue interface{}) interface{} {
	if len(path) == 0 {
		return obj
	}

	current := obj
	for _, key := range path {
		switch m := current.(type) {
		case map[string]interface{}:
			val, ok := m[key]
			if !ok {
				return defaultValue
			}
			current = val
		case []interface{}:
			idx, err := strconv.Atoi(key)
			if err != nil || idx < 0 || idx >= len(m) {
				return defaultValue
			}
			current = m[idx]
		default:
			return defaultValue
		}
	}

	return current
}

// popFirst returns a new slice without the first element.
func popFirst(arr []string) []string {
	if len(arr) <= 1 {
		return nil
	}
	result := make([]string, len(arr)-1)
	copy(result, arr[1:])
	return result
}

// popLast returns a new slice without the last element.
func popLast(arr []string) []string {
	if len(arr) <= 1 {
		return nil
	}
	result := make([]string, len(arr)-1)
	copy(result, arr[:len(arr)-1])
	return result
}

// isDirectReference checks if the reference is direct (does not end with "_").
func isDirectReference(ref []string) bool {
	if len(ref) == 0 {
		return true
	}
	return ref[len(ref)-1] != "_"
}

// canonicalSeed converts a seed to a canonical string representation.
// Uses json.Marshal for maps/slices to ensure deterministic key ordering.
func canonicalSeed(seed interface{}) string {
	switch seed.(type) {
	case map[string]interface{}, []interface{}:
		b, err := json.Marshal(seed)
		if err == nil {
			return string(b)
		}
	}
	return fmt.Sprintf("%v", seed)
}

// newUUID generates a new UUID from a seed and current time.
func newUUID(seed interface{}) uuid.UUID {
	now := fmt.Sprintf("%d", time.Now().UnixNano())
	data := canonicalSeed(seed) + now
	return uuid.NewV5(uuid.NamespaceDNS, data)
}

// staticUUID generates a deterministic UUID from a seed (no time component).
func staticUUID(seed interface{}) uuid.UUID {
	return uuid.NewV5(uuid.NamespaceDNS, canonicalSeed(seed))
}

// matchName checks if a component name matches a selector pattern.
// Supports wildcards ("*"), exact match, and regex match.
func matchName(name, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if name == pattern {
		return true
	}
	matched, err := regexp.MatchString(pattern, name)
	return err == nil && matched
}

// getMapString safely retrieves a string value from a map.
func getMapString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

// resolvePath resolves wildcards ("_") in a path against an object.
// Wildcards represent array indices and are resolved to the last available index.
func resolvePath(path []string, obj map[string]interface{}) []string {
	wildcardIdx := -1
	for i, p := range path {
		if p == "_" {
			wildcardIdx = i
			break
		}
	}

	if wildcardIdx == -1 {
		return path
	}

	prefix := path[:wildcardIdx]
	suffix := path[wildcardIdx+1:]

	val := objectGetNested(obj, prefix, nil)
	arr, ok := val.([]interface{})
	idx := "0"
	if ok && len(arr) > 0 {
		idx = strconv.Itoa(len(arr) - 1)
	}

	result := make([]string, 0, len(path))
	result = append(result, prefix...)
	result = append(result, idx)
	result = append(result, suffix...)
	return result
}


// deepEqual does a deep comparison of two interface{} values.
func deepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// uuidFromString parses a UUID from a string, returning zero UUID on error.
func uuidFromString(s string) (uuid.UUID, error) {
	return uuid.FromString(s)
}

// deepCopyDesign creates a deep copy of a PatternFile via JSON round-trip.
func deepCopyDesign(design *pattern.PatternFile) *pattern.PatternFile {
	data, err := json.Marshal(design)
	if err != nil {
		return design
	}
	var cp pattern.PatternFile
	if err := json.Unmarshal(data, &cp); err != nil {
		return design
	}
	return &cp
}

// deepCopyRelDef creates a deep copy of a RelationshipDefinition via JSON round-trip.
func deepCopyRelDef(rel *relationship.RelationshipDefinition) *relationship.RelationshipDefinition {
	data, err := json.Marshal(rel)
	if err != nil {
		return rel
	}
	var cp relationship.RelationshipDefinition
	if err := json.Unmarshal(data, &cp); err != nil {
		return rel
	}
	return &cp
}

// Relationship status constants. Uses schema-defined values where available.
const (
	StatusApproved   = string(relationship.Approved)
	StatusDeleted    = string(relationship.Deleted)
	StatusPending    = string(relationship.Pending)
	StatusIdentified = "identified"
)

// getRelStatus returns the status string of a relationship, or empty if nil.
func getRelStatus(rel *relationship.RelationshipDefinition) string {
	if rel.Status == nil {
		return ""
	}
	return string(*rel.Status)
}

// setRelStatus sets the status on a relationship.
func setRelStatus(rel *relationship.RelationshipDefinition, status string) {
	s := relationship.RelationshipDefinitionStatus(status)
	rel.Status = &s
}

// getModelsInDesign extracts unique model names from design components.
// Checks both ModelReference.Name and Model.Name for compatibility with
// designs that may only have one or the other populated.
func getModelsInDesign(design *pattern.PatternFile) []string {
	seen := make(map[string]bool)
	var names []string
	for _, comp := range design.Components {
		name := comp.ModelReference.Name
		if name == "" && comp.Model != nil {
			name = comp.Model.Name
		}
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}
	return names
}

// filterRelationshipsInScope filters relationships to those matching models in the design.
func filterRelationshipsInScope(
	allRels []*relationship.RelationshipDefinition,
	modelNames []string,
	design *pattern.PatternFile,
) []*relationship.RelationshipDefinition {
	modelSet := make(map[string]bool, len(modelNames))
	for _, name := range modelNames {
		modelSet[name] = true
	}

	var layerPrefs map[string]interface{}
	if design.Preferences != nil {
		if layers, ok := design.Preferences.Layers["relationships"]; ok {
			if m, ok := layers.(map[string]interface{}); ok {
				layerPrefs = m
			}
		}
	}

	var result []*relationship.RelationshipDefinition
	for _, rel := range allRels {
		if !modelSet[rel.Model.Name] {
			continue
		}
		if layerPrefs != nil {
			key := strings.ToLower(string(rel.Kind)) + "-" + strings.ToLower(rel.RelationshipType) + "-" + strings.ToLower(rel.SubType)
			if v, exists := layerPrefs[key]; exists {
				if b, ok := v.(bool); ok && !b {
					continue
				}
			}
		}
		result = append(result, rel)
	}
	return result
}

// toGenericMap converts any struct to a map[string]interface{} via JSON round-trip.
func toGenericMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}


