package jsonschema

import (
	"context"
	"sync"

	"github.com/qri-io/jsonschema"
)

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
