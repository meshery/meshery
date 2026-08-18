package stages

import (
	"testing"

	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

func TestFilterAnnotations(t *testing.T) {
	annotation := &component.ComponentDefinition{}
	annotation.Metadata.IsAnnotation = true

	regular := &component.ComponentDefinition{}
	regular.Metadata.IsAnnotation = false

	data := &Data{
		Pattern: &pattern.PatternFile{
			Components: []*component.ComponentDefinition{
				annotation,
				regular,
			},
		},
	}

	called := false

	FilterAnnotations()(data, nil, func(data *Data, err error) {
		called = true

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(data.Pattern.Components) != 1 {
			t.Fatalf("expected 1 component after filtering, got %d", len(data.Pattern.Components))
		}

		if data.Pattern.Components[0] != regular {
			t.Fatal("expected regular component to remain after filtering")
		}
	})

	if !called {
		t.Fatal("expected next stage to be called")
	}
}
