package models

import (
	"strings"
	"testing"

	"github.com/meshery/meshkit/logger"
)

func TestConvertToSpec_MissingRunType(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	m := &MesheryResult{
		Result: map[string]interface{}{
			"RetCodes": map[string]interface{}{},
		},
	}

	_, err = m.ConvertToSpec(log)
	if err == nil {
		t.Fatal("expected error for missing RunType, got nil")
	}
	if !strings.Contains(err.Error(), "RunType") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestConvertToSpec_InvalidRunType(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	m := &MesheryResult{
		Result: map[string]interface{}{
			"RunType":  123,
			"RetCodes": map[string]interface{}{},
		},
	}

	_, err = m.ConvertToSpec(log)
	if err == nil {
		t.Fatal("expected error for non-string RunType, got nil")
	}
	if !strings.Contains(err.Error(), "RunType") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
