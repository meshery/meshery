package models

import (
	"encoding/json"
	"testing"
)

// TestRestrictedAccessJSONTagCanonicalSpelling pins the wire-side spelling
// of RestrictedAccess.IsMesheryUIRestricted to the canonical post-flip
// `isMesheryUIRestricted` (uppercase UI). Meshery Cloud emits this exact
// key on the /{version}/capabilities response; the legacy lowercase-Ui tag
// silently dropped that value and left the playground/restricted-UI gate
// stuck at false, which presented as a cloud.layer5.io ⇄ localhost:9081
// redirect loop on enforced-provider hosts. Guarding the tag here means a
// future tag refactor can't reintroduce that class of bug without flipping
// this assertion as well.
func TestRestrictedAccessJSONTagCanonicalSpelling(t *testing.T) {
	const payload = `{"isMesheryUIRestricted": true, "allowedComponents": {}}`

	var ra RestrictedAccess
	if err := json.Unmarshal([]byte(payload), &ra); err != nil {
		t.Fatalf("unmarshal canonical payload: %v", err)
	}
	if !ra.IsMesheryUIRestricted {
		t.Fatalf("IsMesheryUIRestricted = false; want true — JSON tag must be `isMesheryUIRestricted` (uppercase UI) to match the cloud emitter")
	}

	out, err := json.Marshal(ra)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(out, &raw); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}
	if _, ok := raw["isMesheryUIRestricted"]; !ok {
		t.Fatalf("server output missing canonical key `isMesheryUIRestricted`; got keys: %v", keys(raw))
	}
	if _, ok := raw["isMesheryUiRestricted"]; ok {
		t.Fatalf("server output emitted legacy key `isMesheryUiRestricted`; tag must be canonical only")
	}
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
