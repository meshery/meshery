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
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// objectGetNestedWithSpecFallback tries objectGetNested with the full path,
// and if nil, retries by skipping "spec" after "configuration".
// This handles relationship definitions that include "spec" in paths
// but components that store fields directly under "configuration".
func objectGetNestedWithSpecFallback(obj interface{}, path []string) interface{} {
	val := objectGetNested(obj, path, nil)
	if val != nil {
		return val
	}
	// Try skipping "spec" if path is like ["configuration", "spec", ...]
	if len(path) > 2 && path[0] == "configuration" && path[1] == "spec" {
		return objectGetNested(obj, append([]string{path[0]}, path[2:]...), nil)
	}
	return nil
}

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

// getMapSlice safely retrieves a []interface{} value from a map.
func getMapSlice(m map[string]interface{}, key string) []interface{} {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	s, ok := v.([]interface{})
	if !ok {
		return nil
	}
	return s
}

// getMapMap safely retrieves a map[string]interface{} value from a map.
func getMapMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	s, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return s
}

// toStringSlice converts []interface{} to []string.
func toStringSlice(arr []interface{}) []string {
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		} else {
			result = append(result, fmt.Sprintf("%v", v))
		}
	}
	return result
}

// deepCopyMap creates a deep copy of a map[string]interface{}.
func deepCopyMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			result[k] = deepCopyMap(val)
		case []interface{}:
			result[k] = deepCopySlice(val)
		default:
			result[k] = v
		}
	}
	return result
}

// deepCopySlice creates a deep copy of a []interface{}.
func deepCopySlice(s []interface{}) []interface{} {
	if s == nil {
		return nil
	}
	result := make([]interface{}, len(s))
	for i, v := range s {
		switch val := v.(type) {
		case map[string]interface{}:
			result[i] = deepCopyMap(val)
		case []interface{}:
			result[i] = deepCopySlice(val)
		default:
			result[i] = v
		}
	}
	return result
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


// stringSliceToInterface converts []string to []interface{}.
func stringSliceToInterface(ss []string) []interface{} {
	result := make([]interface{}, len(ss))
	for i, s := range ss {
		result[i] = s
	}
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

// relationshipPreferenceKey builds the preference key for a relationship map.
func relationshipPreferenceKey(rel map[string]interface{}) string {
	kind := strings.ToLower(getMapString(rel, "kind"))
	relType := strings.ToLower(getMapString(rel, "type"))
	subType := strings.ToLower(getMapString(rel, "subType"))
	return kind + "-" + relType + "-" + subType
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

// mapToComponentDef converts a generic map to a ComponentDefinition.
func mapToComponentDef(m map[string]interface{}) (component.ComponentDefinition, error) {
	var comp component.ComponentDefinition
	data, err := json.Marshal(m)
	if err != nil {
		return comp, err
	}
	err = json.Unmarshal(data, &comp)
	return comp, err
}

// mapToRelationshipDef converts a generic map to a RelationshipDefinition.
func mapToRelationshipDef(m map[string]interface{}) (relationship.RelationshipDefinition, error) {
	var rel relationship.RelationshipDefinition
	data, err := json.Marshal(m)
	if err != nil {
		return rel, err
	}
	err = json.Unmarshal(data, &rel)
	return rel, err
}
