package jsonschema

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/qri-io/jsonschema"
)

var (
	// mu synchronizes schema unmarshaling and validation to prevent concurrent
	// map writes in the underlying qri-io/jsonschema library.
	mu sync.Mutex
)

// ValidateBytes parses schemaBytes and validates dataBytes against it in a thread-safe manner.
// A fresh jsonschema.Schema instance is instantiated on each call to prevent concurrent
// schema mutation and cross-component validation pollution.
func ValidateBytes(ctx context.Context, schemaBytes, dataBytes []byte) ([]jsonschema.KeyError, error) {
	mu.Lock()
	defer mu.Unlock()

	rs := &jsonschema.Schema{}
	if err := json.Unmarshal(schemaBytes, rs); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}
	return rs.ValidateBytes(ctx, dataBytes)
}
