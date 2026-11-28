package stages

import (
	"sync"
	"testing"

	"github.com/meshery/schemas/models/v1beta2/component"
)

// componentWithSchema returns a component definition whose configuration
// schema requires "replicas" to have the given JSON type.
func componentWithSchema(kind, replicasType string) component.ComponentDefinition {
	return component.ComponentDefinition{
		Component: component.Component{
			Kind:    kind,
			Version: "v1",
			Schema: `{"type":"object","properties":{"replicas":{"type":"` +
				replicasType + `"}},"required":["replicas"]}`,
		},
	}
}

func TestValidateWorkload(t *testing.T) {
	integerReplicas := componentWithSchema("IntegerKind", "integer")

	tests := []struct {
		name    string
		config  map[string]interface{}
		def     component.ComponentDefinition
		wantErr bool
	}{
		{
			name:   "configuration matching the schema is accepted",
			config: map[string]interface{}{"replicas": 3},
			def:    integerReplicas,
		},
		{
			name:    "configuration violating the schema is rejected",
			config:  map[string]interface{}{"replicas": "three"},
			def:     integerReplicas,
			wantErr: true,
		},
		{
			name:    "a schema that cannot be parsed is reported",
			config:  map[string]interface{}{"replicas": 3},
			def:     component.ComponentDefinition{Component: component.Component{Schema: "{not json"}},
			wantErr: true,
		},
		{
			name:   "an annotation component without a schema is skipped",
			config: map[string]interface{}{"anything": true},
			def: component.ComponentDefinition{
				Metadata: component.ComponentDefinition_Metadata{IsAnnotation: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateWorkload(tt.config, tt.def)
			if tt.wantErr && err == nil {
				t.Fatal("expected a validation error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

// TestValidateWorkloadConcurrent validates two components with conflicting
// schemas from many goroutines at once. Each configuration is valid for its
// own component and invalid for the other, so a validation that picks up the
// other goroutine's schema fails the test.
func TestValidateWorkloadConcurrent(t *testing.T) {
	integerReplicas := componentWithSchema("IntegerKind", "integer")
	stringReplicas := componentWithSchema("StringKind", "string")
	integerConfig := map[string]interface{}{"replicas": 3}
	stringConfig := map[string]interface{}{"replicas": "three"}

	const workers = 8
	const iterations = 200

	var wg sync.WaitGroup
	errs := make(chan error, workers*iterations)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				var err error
				if (w+i)%2 == 0 {
					err = validateWorkload(integerConfig, integerReplicas)
				} else {
					err = validateWorkload(stringConfig, stringReplicas)
				}
				if err != nil {
					errs <- err
				}
			}
		}(w)
	}
	wg.Wait()
	close(errs)

	failures := 0
	for err := range errs {
		if failures == 0 {
			t.Errorf("a valid configuration was rejected: %v", err)
		}
		failures++
	}
	if failures > 0 {
		t.Errorf("%d of %d concurrent validations used the wrong schema", failures, workers*iterations)
	}
}
