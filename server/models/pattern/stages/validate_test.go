package stages

import "testing"

// TestHydrateComponentWithOriginalType covers the spec shapes that reach
// hydrateComponentWithOriginalType for EnvoyFilter/WasmPlugin components. A
// non-object spec (string, number, slice) used to hit an unchecked type
// assertion spec.(map[string]interface{}) and panic, which surfaced as an
// HTTP 500 with a stack trace during design deploy instead of a clean
// validation error. It must now return an error without panicking.
func TestHydrateComponentWithOriginalType(t *testing.T) {
	tests := []struct {
		name     string
		compType string
		spec     interface{}
		wantErr  bool
	}{
		{name: "nil spec", compType: "EnvoyFilter", spec: nil, wantErr: true},
		{name: "string spec", compType: "EnvoyFilter", spec: "not-a-map", wantErr: true},
		{name: "number spec", compType: "WasmPlugin", spec: 42, wantErr: true},
		{name: "slice spec", compType: "EnvoyFilter", spec: []interface{}{"a", "b"}, wantErr: true},
		{name: "valid object spec", compType: "WasmPlugin", spec: map[string]interface{}{}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hydrateComponentWithOriginalType(tt.compType, tt.spec)
			if (err != nil) != tt.wantErr {
				t.Fatalf("hydrateComponentWithOriginalType(%s, %T) err = %v, wantErr = %v", tt.compType, tt.spec, err, tt.wantErr)
			}
		})
	}
}
