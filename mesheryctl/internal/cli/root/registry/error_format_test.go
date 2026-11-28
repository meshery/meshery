package registry

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

// TestErrGenerateComponentFormat is a regression test for issue #22681: the
// short description carried bare "%s" format verbs that were never fed through
// Sprintf, so compName and modelName were appended as stray slice elements and
// the literal "%s" surfaced to the user. The short description must now contain
// both interpolated values and no literal format verb.
func TestErrGenerateComponentFormat(t *testing.T) {
	const (
		compName  = "my-component"
		modelName = "my-model"
	)
	underlying := errors.New("schema validation failed")

	err := ErrGenerateComponent(underlying, modelName, compName)
	if err == nil {
		t.Fatal("ErrGenerateComponent returned nil, expected a wrapped error")
	}

	sdesc := meshkiterrors.GetSDescription(err)
	if !strings.Contains(sdesc, compName) {
		t.Errorf("short description %q does not contain interpolated compName %q", sdesc, compName)
	}
	if !strings.Contains(sdesc, modelName) {
		t.Errorf("short description %q does not contain interpolated modelName %q", sdesc, modelName)
	}
	if strings.Contains(sdesc, "%s") || strings.Contains(sdesc, "%v") {
		t.Errorf("short description %q still contains a literal format verb", sdesc)
	}
}

func TestRegistryErrors_NilErrorPanicGuards(t *testing.T) {
	tests := []struct {
		name string
		fn   func() error
	}{
		{"ErrUpdateRegistry with nil error", func() error { return ErrUpdateRegistry(nil, "./models") }},
		{"ErrGenerateModel with nil error", func() error { return ErrGenerateModel(nil, "my-model") }},
		{"ErrGenerateComponent with nil error", func() error { return ErrGenerateComponent(nil, "my-model", "my-comp") }},
		{"ErrUpdateModel with nil error", func() error { return ErrUpdateModel(nil, "my-model") }},
		{"ErrUpdateComponent with nil error", func() error { return ErrUpdateComponent(nil, "my-model", "my-comp") }},
		{"ErrParsingSheet with nil error", func() error { return ErrParsingSheet(nil, "models") }},
		{"ErrParsingRelationship with nil error", func() error { return ErrParsingRelationship(nil, "my-model") }},
		{"ErrPublish with nil error", func() error { return ErrPublish(nil, "website") }},
		{"ErrModelGenerationFailed with nil error", func() error { return ErrModelGenerationFailed("my-model", nil) }},
		{"ErrPurgeReadModelsDir with nil error", func() error { return ErrPurgeReadModelsDir(nil, "./models") }},
		{"ErrPurgeRemove with nil error", func() error { return ErrPurgeRemove(nil, "./models/v1.0.0") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				err := tt.fn()
				assert.NotNil(t, err)
			})
		})
	}
}

func TestErrUpdateRegistry_NilAndCustomErrors(t *testing.T) {
	// Assert nil error does not panic and returns non-nil error
	errNil := ErrUpdateRegistry(nil, "./models")
	assert.NotNil(t, errNil)

	// Assert custom error is captured
	customErr := fmt.Errorf("request to Google Spreadsheet returned HTTP status 500")
	errCustom := ErrUpdateRegistry(customErr, "./models")
	assert.NotNil(t, errCustom)
	assert.Contains(t, errCustom.Error(), "500")
}
