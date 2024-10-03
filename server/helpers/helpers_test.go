package helpers

import (
	"os"
	"testing"
)

func TestMergeStringMaps(t *testing.T) {
	map1 := map[string]string{"key1": "value1", "key2": "value2"}
	map2 := map[string]string{"key2": "override", "key3": "value3"}

	expected := map[string]string{"key1": "value1", "key2": "override", "key3": "value3"}
	result := MergeStringMaps(map1, map2)

	if len(result) != len(expected) {
		t.Errorf("Expected map length %d, got %d", len(expected), len(result))
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Expected key %s to have value %s, but got %s", k, v, result[k])
		}
	}
}

func TestResolveFSRef(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "test content"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	result, err := ResolveFSRef(tmpfile.Name())
	if err != nil {
		t.Fatalf("ResolveFSRef failed: %v", err)
	}
	if result != content {
		t.Errorf("Expected file content '%s', got '%s'", content, result)
	}
}
