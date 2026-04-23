package policies

import (
	"testing"

	"github.com/gofrs/uuid"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

func TestResolveTies_ConflictDropped(t *testing.T) {
	compID := "c1"
	path := []string{"configuration", "spec", "ns"}
	actions := []PolicyAction{
		newComponentUpdateAction(UpdateComponentConfigurationOp, compID, path, "alpha"),
		newComponentUpdateAction(UpdateComponentConfigurationOp, compID, path, "beta"),
	}
	filtered, ties := resolveTies(actions)
	if len(filtered) != 0 {
		t.Errorf("expected all conflicting actions dropped, got %d", len(filtered))
	}
	if len(ties) != 1 {
		t.Fatalf("expected 1 tie recorded, got %d", len(ties))
	}
	if ties[0].ComponentID != compID {
		t.Errorf("tie componentID: got %q want %q", ties[0].ComponentID, compID)
	}
}

func TestResolveTies_AgreementDeduplicated(t *testing.T) {
	compID := "c1"
	path := []string{"configuration", "replicas"}
	actions := []PolicyAction{
		newComponentUpdateAction(UpdateComponentConfigurationOp, compID, path, 3),
		newComponentUpdateAction(UpdateComponentConfigurationOp, compID, path, 3),
	}
	filtered, ties := resolveTies(actions)
	if len(filtered) != 1 {
		t.Errorf("expected 1 action after dedup, got %d", len(filtered))
	}
	if len(ties) != 0 {
		t.Errorf("expected no ties, got %d", len(ties))
	}
}

func TestResolveTies_NonConflictingActionsKept(t *testing.T) {
	actions := []PolicyAction{
		newComponentUpdateAction(UpdateComponentConfigurationOp, "c1", []string{"configuration", "a"}, "x"),
		newComponentUpdateAction(UpdateComponentConfigurationOp, "c2", []string{"configuration", "a"}, "y"),
		newComponentUpdateAction(UpdateComponentConfigurationOp, "c1", []string{"configuration", "b"}, "z"),
	}
	filtered, ties := resolveTies(actions)
	if len(filtered) != 3 {
		t.Errorf("expected all 3 actions kept, got %d", len(filtered))
	}
	if len(ties) != 0 {
		t.Errorf("expected no ties, got %d", len(ties))
	}
}

func TestActionsFingerprint_StableAcrossOrder(t *testing.T) {
	a := []PolicyAction{
		newComponentUpdateAction(UpdateComponentConfigurationOp, "c1", []string{"configuration", "a"}, "x"),
		newComponentUpdateAction(UpdateComponentConfigurationOp, "c2", []string{"configuration", "b"}, "y"),
	}
	b := []PolicyAction{a[1], a[0]}
	if actionsFingerprint(a) != actionsFingerprint(b) {
		t.Error("fingerprint should be order-independent")
	}
}

func TestActionsFingerprint_DifferentForDifferentValues(t *testing.T) {
	a := []PolicyAction{newComponentUpdateAction(UpdateComponentConfigurationOp, "c1", []string{"configuration", "a"}, "x")}
	b := []PolicyAction{newComponentUpdateAction(UpdateComponentConfigurationOp, "c1", []string{"configuration", "a"}, "y")}
	if actionsFingerprint(a) == actionsFingerprint(b) {
		t.Error("different values should produce different fingerprints")
	}
}

func TestRecordAndCheckFlap_TripleInversionTriggers(t *testing.T) {
	hist := map[string][]string{}
	for i, fp := range []string{"A", "B", "A"} {
		if flap := recordAndCheckFlap(hist, "d1", fp); flap {
			t.Fatalf("unexpected flap at step %d", i)
		}
	}
	if !recordAndCheckFlap(hist, "d1", "B") {
		t.Error("expected flap on A-B-A-B pattern")
	}
}

func TestRecordAndCheckFlap_StableDoesNotTrigger(t *testing.T) {
	hist := map[string][]string{}
	for _, fp := range []string{"A", "A", "A", "A"} {
		if flap := recordAndCheckFlap(hist, "d1", fp); flap {
			t.Error("stable sequence should not flap")
		}
	}
}

func TestEngineFlappingDetected_AcrossEvaluations(t *testing.T) {
	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)
	designID, _ := uuid.FromString("00000000-0000-0000-0000-0000000000aa")

	// Seed history so the next evaluation completes an A-B-A-B alternation.
	// The current evaluation produces zero actions, so its fingerprint is
	// fingerprintEmpty. We use a distinct "other" fingerprint to alternate.
	other := "other-fp"
	empty := actionsFingerprint(nil)
	if empty == other {
		t.Fatalf("test setup assumption broken: fingerprints collide")
	}
	engine.flapHist[designID.String()] = []string{other, empty, other}

	design := pattern.PatternFile{ID: designID}
	_, err := engine.EvaluateDesign(design, nil)
	if err == nil {
		t.Fatal("expected ErrFlappingDetected, got nil")
	}
	if code := meshkiterrors.GetCode(err); code != ErrFlappingDetectedCode {
		t.Errorf("expected error code %s, got %s", ErrFlappingDetectedCode, code)
	}
}
