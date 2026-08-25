package jsonschema

import (
	"context"

	"github.com/qri-io/jsonschema"
)

type Schema struct {
	jsonschema.Schema
}

// JsonSchema package creates a local instance of Schema struct which avoids internal concurrent map writes by isolating the schema state.
func (s *Schema) ValidateBytes(ctx context.Context, data []byte) ([]jsonschema.KeyError, error) {
	return s.Schema.ValidateBytes(ctx, data)
}
