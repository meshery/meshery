package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestJSONOutputFormatter_Display(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	data := testStruct{Name: "Alice", Age: 30}
	buf := &bytes.Buffer{}
	formatter := NewJSONOutputFormatter(data).WithOutput(buf)

	err := formatter.Display()
	if err != nil {
		t.Fatalf("Display() returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Alice") || !strings.Contains(output, "30") {
		t.Errorf("JSON output missing expected data: %s", output)
	}

	// Test error case: invalid output writer
	badFormatter := NewJSONOutputFormatter(data).WithOutput(badWriter{})
	if err := badFormatter.Display(); err == nil {
		t.Error("Expected error for bad writer, got nil")
	}
}

type badWriter struct{}

func (badWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write error")
}

func TestYAMLOutputFormatter_Display(t *testing.T) {
	type testStruct struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}
	data := testStruct{Name: "Bob", Age: 25}
	buf := &bytes.Buffer{}
	formatter := (&YAMLOutputFormatter[testStruct]{Data: data}).WithOutput(buf)

	err := formatter.Display()
	if err != nil {
		t.Fatalf("Display() returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Bob") || !strings.Contains(output, "25") {
		t.Errorf("YAML output missing expected data: %s", output)
	}

	// Test error case: invalid output writer
	badFormatter := (&YAMLOutputFormatter[testStruct]{Data: data}).WithOutput(badWriter{})
	if err := badFormatter.Display(); err == nil {
		t.Error("Expected error for bad writer, got nil")
	}
}