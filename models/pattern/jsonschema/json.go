package jsonschema

import (
	"context"
	"sync"

	"github.com/qri-io/jsonschema"
)

var JSONSchema = &Schema{}

// This approach is bad, and we are not the ones implementing this. The qri-io/jsonschema internally is creating a global instance of Schema.
// Hence for concurrent operations, we have to make sure that the package qri-io/jsonschema is accessed in a thread-safe way.
//So use this Global instance for now.
func GlobalJSONSchema() *Schema {
	return JSONSchema
}

type Schema struct {
	jsonschema.Schema
	Lock sync.Mutex
}

// JsonSchema package creates a global instance(without any locks) of Schema struct which panics when concurrent routines try to call ValidateBytes.
// So this package creates a thin shim to avoid internal concurrent map writes
func (s *Schema) ValidateBytes(ctx context.Context, data []byte) ([]jsonschema.KeyError, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return s.Schema.ValidateBytes(ctx, data)
}
