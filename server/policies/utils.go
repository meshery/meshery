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

// interfaceToStringSlice converts an interface{} to []string.
func interfaceToStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case []string:
		return val
	case []interface{}:
		return toStringSlice(val)
	default:
		return nil
	}
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

// relationshipPreferenceKey builds the preference key for a relationship.
func relationshipPreferenceKey(rel map[string]interface{}) string {
	kind := strings.ToLower(getMapString(rel, "kind"))
	relType := strings.ToLower(getMapString(rel, "type"))
	subType := strings.ToLower(getMapString(rel, "subType"))
	return kind + "-" + relType + "-" + subType
}
