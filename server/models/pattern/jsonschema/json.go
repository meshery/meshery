package jsonschema

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/qri-io/jsonschema"
)

var (
	// mu synchronizes schema unmarshaling and validation to prevent concurrent
	// map writes in the underlying qri-io/jsonschema library.
	mu sync.Mutex

	// JSONSchema is retained for backwards compatibility.
	JSONSchema = &Schema{}
)

// ValidateBytes parses schemaBytes and validates dataBytes against it in a thread-safe manner.
// A fresh jsonschema.Schema instance is instantiated on each call to prevent concurrent
// schema mutation and cross-component validation pollution.
func ValidateBytes(ctx context.Context, schemaBytes, dataBytes []byte) ([]jsonschema.KeyError, error) {
	mu.Lock()
	defer mu.Unlock()

	rs := &jsonschema.Schema{}
	if err := json.Unmarshal(schemaBytes, rs); err != nil {
		return nil, err
	}
	return rs.ValidateBytes(ctx, dataBytes)
}

// GlobalJSONSchema returns the package's shared Schema instance.
//
// Deprecated: Use ValidateBytes directly for thread-safe component validation.
func GlobalJSONSchema() *Schema {
	return JSONSchema
}

type Schema struct {
	jsonschema.Schema
	Lock sync.Mutex
}

// ValidateBytes validates data against the schema while holding the schema instance's lock.
func (s *Schema) ValidateBytes(ctx context.Context, data []byte) ([]jsonschema.KeyError, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return s.Schema.ValidateBytes(ctx, data)
}
