package policies

import (
	"testing"
	"time"
)

func TestObjectGetNested(t *testing.T) {
	obj := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "value",
			},
		},
		"arr": []interface{}{"x", "y", "z"},
	}

	tests := []struct {
		name     string
		path     []string
		expected interface{}
	}{
		{"simple", []string{"a", "b", "c"}, "value"},
		{"missing", []string{"a", "x"}, "default"},
		{"empty path", nil, obj},
		{"array index", []string{"arr", "1"}, "y"},
		{"array out of bounds", []string{"arr", "5"}, "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := objectGetNested(obj, tt.path, "default")
			if result == nil && tt.expected == nil {
				return
			}
			// For map comparisons, just check non-nil
			if _, ok := result.(map[string]interface{}); ok {
				return
			}
			if result != tt.expected {
				t.Errorf("objectGetNested(%v) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestMatchName(t *testing.T) {
	tests := []struct {
		name, pattern string
		expected      bool
	}{
		{"Namespace", "*", true},
		{"Namespace", "Namespace", true},
		{"Namespace", "Service", false},
		{"Deployment", "Deploy.*", true},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.pattern, func(t *testing.T) {
			result := matchName(tt.name, tt.pattern)
			if result != tt.expected {
				t.Errorf("matchName(%q, %q) = %v, want %v", tt.name, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestMatchValues(t *testing.T) {
	tests := []struct {
		name     string
		from, to interface{}
		strategy string
		expected bool
	}{
		{"equal_match", "test", "test", "equal", true},
		{"equal_mismatch", "test", "other", "equal", false},
		{"equal_as_strings", 42, "42", "equal_as_strings", true},
		{"not_null_both", "a", "b", "not_null", true},
		{"not_null_nil", nil, "b", "not_null", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchValues(tt.from, tt.to, tt.strategy)
			if result != tt.expected {
				t.Errorf("matchValues(%v, %v, %q) = %v, want %v", tt.from, tt.to, tt.strategy, result, tt.expected)
			}
		})
	}
}

func TestPopLast(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"multiple", []string{"a", "b", "c"}, []string{"a", "b"}},
		{"two", []string{"a", "b"}, []string{"a"}},
		{"single", []string{"a"}, nil},
		{"empty", []string{}, nil},
		{"nil", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := popLast(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("popLast(%v) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("popLast(%v)[%d] = %v, want %v", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestIsDirectReference(t *testing.T) {
	tests := []struct {
		name     string
		ref      []string
		expected bool
	}{
		{"all named segments is direct", []string{"configuration", "name"}, true},
		{"wildcard segment makes ref indirect", []string{"configuration", "containers", "_"}, false},
		{"empty ref is direct", []string{}, true},
		{"single named segment is direct", []string{"name"}, true},
		{"single wildcard segment is indirect", []string{"_"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDirectReference(tt.ref)
			if result != tt.expected {
				t.Errorf("isDirectReference(%v) = %v, want %v", tt.ref, result, tt.expected)
			}
		})
	}
}

func TestStringSliceToInterface(t *testing.T) {
	tests := []struct {
		name  string
		input []string
	}{
		{"empty slice returns empty interface slice", []string{}},
		{"single element string converts to interface", []string{"a"}},
		{"multiple elements preserve order and values", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringSliceToInterface(tt.input)
			if len(result) != len(tt.input) {
				t.Fatalf("Expected %d elements, got %d", len(tt.input), len(result))
			}
			for i, v := range result {
				s, ok := v.(string)
				if !ok {
					t.Errorf("Element %d is not a string", i)
				}
				if s != tt.input[i] {
					t.Errorf("Element %d = %q, want %q", i, s, tt.input[i])
				}
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	obj := map[string]interface{}{
		"containers": []interface{}{"nginx", "sidecar", "init"},
	}

	t.Run("no wildcard", func(t *testing.T) {
		path := []string{"name"}
		result := resolvePath(path, obj)
		if len(result) != 1 || result[0] != "name" {
			t.Errorf("Expected [name], got %v", result)
		}
	})

	t.Run("wildcard present", func(t *testing.T) {
		path := []string{"containers", "_", "name"}
		result := resolvePath(path, obj)
		if len(result) != 3 || result[1] != "2" {
			t.Errorf("Expected wildcard resolved to index 2, got %v", result)
		}
	})

	t.Run("wildcard empty array", func(t *testing.T) {
		emptyObj := map[string]interface{}{
			"containers": []interface{}{},
		}
		path := []string{"containers", "_"}
		result := resolvePath(path, emptyObj)
		if len(result) != 2 || result[1] != "0" {
			t.Errorf("Expected wildcard resolved to 0 for empty array, got %v", result)
		}
	})

	t.Run("wildcard non-array target", func(t *testing.T) {
		nonArrayObj := map[string]interface{}{
			"name": "test",
		}
		path := []string{"name", "_"}
		result := resolvePath(path, nonArrayObj)
		if len(result) != 2 || result[1] != "0" {
			t.Errorf("Expected wildcard resolved to 0 for non-array, got %v", result)
		}
	})
}

func TestCanonicalSeedDeterminism(t *testing.T) {
	map1 := map[string]interface{}{"a": "1", "b": "2"}
	map2 := map[string]interface{}{"a": "1", "b": "2"}

	s1 := canonicalSeed(map1)
	s2 := canonicalSeed(map2)
	if s1 != s2 {
		t.Errorf("canonicalSeed not deterministic: %q != %q", s1, s2)
	}

	map3 := map[string]interface{}{"a": "1", "b": "3"}
	s3 := canonicalSeed(map3)
	if s1 == s3 {
		t.Error("Expected different seeds for different maps")
	}

	s4 := canonicalSeed("hello")
	s5 := canonicalSeed("hello")
	if s4 != s5 {
		t.Errorf("canonicalSeed not deterministic for strings: %q != %q", s4, s5)
	}
}

func TestStaticUUIDDeterminism(t *testing.T) {
	seed := map[string]interface{}{"key": "value"}
	u1 := staticUUID(seed)
	u2 := staticUUID(seed)
	if u1 != u2 {
		t.Errorf("staticUUID not deterministic: %v != %v", u1, u2)
	}
}

func TestNewUUIDNonDeterminism(t *testing.T) {
	seed := "test-seed"
	u1 := newUUID(seed)
	time.Sleep(time.Millisecond)
	u2 := newUUID(seed)
	if u1 == u2 {
		t.Error("newUUID should produce different UUIDs on successive calls (time-based)")
	}
}
