package utils

import "testing"

// TestIsDesignInAlpha2Format guards the design-import rendering regression:
// freshly imported designs (Kubernetes manifests, Helm charts, Compose files)
// are produced by NewPatternFileFromK8sManifest in v1beta3 form. They must be
// recognized as a current design schema - NOT misclassified as legacy
// v1alpha2 - otherwise convertV1alpha2ToV1beta3 reparses them as a
// v1alpha2.PatternFile (which has `services`, not `components`) and strips
// every component, leaving an empty canvas in Kanvas.
func TestIsDesignInAlpha2Format(t *testing.T) {
	tests := []struct {
		name        string
		patternFile string
		wantAlpha2  bool
	}{
		{
			// The regression: import stamps v1beta3; this must be current.
			name:        "v1beta3 imported design is current (not alpha2)",
			patternFile: `{"schemaVersion":"designs.meshery.io/v1beta3","name":"nginx","components":[{"component":{"kind":"Deployment"},"model":{"name":"kubernetes"}}]}`,
			wantAlpha2:  false,
		},
		{
			name:        "v1beta1 design is current (not alpha2)",
			patternFile: `{"schemaVersion":"designs.meshery.io/v1beta1","name":"legacy-beta1","components":[]}`,
			wantAlpha2:  false,
		},
		{
			// Future-proofing: a new design schema version must not regress
			// the way v1beta3 did when only v1beta1 was whitelisted.
			name:        "future designs.meshery.io version is current (not alpha2)",
			patternFile: `{"schemaVersion":"designs.meshery.io/v1beta4","name":"future","components":[]}`,
			wantAlpha2:  false,
		},
		{
			// YAML-encoded pattern file: guards the json+yaml tag pair on the
			// decode struct. With a json-only tag the yaml.v3 fallback in
			// encoding.Unmarshal would miss the camelCase key and misclassify
			// this current design as alpha2.
			name:        "v1beta3 design encoded as YAML is current (not alpha2)",
			patternFile: "schemaVersion: designs.meshery.io/v1beta3\nname: yaml-nginx\ncomponents: []\n",
			wantAlpha2:  false,
		},
		{
			// Genuine legacy format: carries `services`, no versioned
			// schemaVersion. This is the only case that should migrate.
			name:        "v1alpha2 design (services, no schemaVersion) is alpha2",
			patternFile: `{"name":"old","services":{"svc":{"type":"Deployment"}}}`,
			wantAlpha2:  true,
		},
		{
			// Defensive edge cases: a missing or empty schemaVersion has no
			// canonical prefix, so it is classified legacy (migrate).
			name:        "missing schemaVersion key is alpha2",
			patternFile: `{"name":"no-version","components":[]}`,
			wantAlpha2:  true,
		},
		{
			name:        "empty schemaVersion is alpha2",
			patternFile: `{"schemaVersion":"","name":"empty-ver"}`,
			wantAlpha2:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsDesignInAlpha2Format(tt.patternFile)
			if err != nil {
				t.Fatalf("IsDesignInAlpha2Format() unexpected error: %v", err)
			}
			if got != tt.wantAlpha2 {
				t.Errorf("IsDesignInAlpha2Format() = %v, want %v", got, tt.wantAlpha2)
			}
		})
	}
}
