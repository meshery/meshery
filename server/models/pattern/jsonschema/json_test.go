package jsonschema

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

const (
	schemaA = `{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type": "object",
		"properties": {
			"serviceName": { "type": "string" },
			"port": { "type": "integer", "minimum": 1, "maximum": 65535 }
		},
		"required": ["serviceName", "port"]
	}`

	schemaB = `{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type": "object",
		"properties": {
			"replicas": { "type": "integer", "minimum": 0 },
			"strategy": { "type": "string", "enum": ["RollingUpdate", "Recreate"] }
		},
		"required": ["replicas", "strategy"]
	}`
)

func TestValidateBytes_Success(t *testing.T) {
	validData := []byte(`{"serviceName": "web-svc", "port": 8080}`)
	keyErrors, err := ValidateBytes(context.Background(), []byte(schemaA), validData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keyErrors) > 0 {
		t.Fatalf("expected 0 key errors, got: %v", keyErrors)
	}
}

func TestValidateBytes_InvalidData(t *testing.T) {
	invalidData := []byte(`{"serviceName": "web-svc", "port": 99999}`) // port exceeds max
	keyErrors, err := ValidateBytes(context.Background(), []byte(schemaA), invalidData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keyErrors) == 0 {
		t.Fatalf("expected validation errors, got none")
	}
}

func TestValidateBytes_InvalidSchema(t *testing.T) {
	_, err := ValidateBytes(context.Background(), []byte(`{invalid-json`), []byte(`{}`))
	if err == nil {
		t.Fatalf("expected unmarshal error for invalid schema, got nil")
	}
	if !errors.Is(err, ErrSchemaCreate) {
		t.Fatalf("expected ErrSchemaCreate sentinel error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "failed to create schema") {
		t.Fatalf("expected error to contain 'failed to create schema', got: %v", err)
	}
}

func TestValidateBytes_ConcurrentDifferentSchemas(t *testing.T) {
	const iterations = 50
	var wg sync.WaitGroup

	errCh := make(chan error, iterations*2)

	for i := 0; i < iterations; i++ {
		wg.Add(2)

		// Goroutine validating Schema A
		go func(idx int) {
			defer wg.Done()
			data := []byte(fmt.Sprintf(`{"serviceName": "svc-%d", "port": %d}`, idx, 1000+(idx%1000)))
			keyErrs, err := ValidateBytes(context.Background(), []byte(schemaA), data)
			if err != nil {
				errCh <- fmt.Errorf("routine A[%d] unexpected error: %w", idx, err)
				return
			}
			if len(keyErrs) > 0 {
				errCh <- fmt.Errorf("routine A[%d] unexpected validation errors: %v", idx, keyErrs)
				return
			}
		}(i)

		// Goroutine validating Schema B
		go func(idx int) {
			defer wg.Done()
			strategy := "RollingUpdate"
			if idx%2 == 0 {
				strategy = "Recreate"
			}
			data := []byte(fmt.Sprintf(`{"replicas": %d, "strategy": "%s"}`, idx, strategy))
			keyErrs, err := ValidateBytes(context.Background(), []byte(schemaB), data)
			if err != nil {
				errCh <- fmt.Errorf("routine B[%d] unexpected error: %w", idx, err)
				return
			}
			if len(keyErrs) > 0 {
				errCh <- fmt.Errorf("routine B[%d] unexpected validation errors: %v", idx, keyErrs)
				return
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("concurrency validation failure: %v", err)
	}
}
