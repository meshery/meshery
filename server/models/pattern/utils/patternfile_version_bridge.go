// This file bridges the v1beta1/pattern.PatternFile and
// v1beta3/design.PatternFile representations by round-tripping through
// JSON. It exists because the evaluation-engine (rego policies,
// EvaluationRequest / EvaluationResponse) is a documented carve-out that
// still consumes v1beta1/pattern, while meshkit's HydratePattern /
// DehydratePattern helpers operate exclusively on v1beta3/design.
//
// The nested component representation differs between the two versions
// (v1beta1/component vs v1beta2/component) but their JSON wire forms
// overlap on the fields the policy engine and dehydrate helper read, so
// a Marshal + Unmarshal is both straightforward and forward-compatible
// with future field additions.
package utils

import (
	"encoding/json"

	patternv1beta1 "github.com/meshery/schemas/models/v1beta1/pattern"
	designv1beta3 "github.com/meshery/schemas/models/v1beta3/design"
)

// PatternV1beta3ToV1beta1 converts a v1beta3/design.PatternFile into the
// v1beta1/pattern.PatternFile shape that the evaluation engine still
// consumes. Errors surface only if JSON encoding / decoding fails, which
// would indicate a malformed PatternFile to begin with.
func PatternV1beta3ToV1beta1(src *designv1beta3.PatternFile) (*patternv1beta1.PatternFile, error) {
	if src == nil {
		return nil, nil
	}
	raw, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	dst := &patternv1beta1.PatternFile{}
	if err := json.Unmarshal(raw, dst); err != nil {
		return nil, err
	}
	return dst, nil
}

// PatternV1beta1ToV1beta3 converts a v1beta1/pattern.PatternFile back to
// the canonical-casing v1beta3/design.PatternFile so meshkit helpers
// (Hydrate/Dehydrate) can operate on it. Inverse of PatternV1beta3ToV1beta1.
func PatternV1beta1ToV1beta3(src *patternv1beta1.PatternFile) (*designv1beta3.PatternFile, error) {
	if src == nil {
		return nil, nil
	}
	raw, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	dst := &designv1beta3.PatternFile{}
	if err := json.Unmarshal(raw, dst); err != nil {
		return nil, err
	}
	return dst, nil
}
