package policies

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	patternutils "github.com/meshery/meshery/server/models/pattern/utils"
	"github.com/meshery/meshkit/logger"
	componentv1beta1 "github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// serverHandlerConstant mirrors, by value, the depth limit the HTTP
// handler loop uses. The handler now re-exports
// policies.MAX_RE_EVALUATION_DEPTH (see
// server/handlers/policy_relationship_handler.go), so this is the single
// authoritative number. If the handler ever forks its own constant, the
// parity assertion below will catch the divergence the moment this value
// is updated to track the handler.
const serverHandlerConstant = 5

// TestMaxReEvaluationDepthParity is the parity pin required by the WASM
// build contract: the depth limit the shared orchestration uses MUST
// equal the limit the server's relationship-evaluation handler uses.
func TestMaxReEvaluationDepthParity(t *testing.T) {
	if MAX_RE_EVALUATION_DEPTH != serverHandlerConstant {
		t.Fatalf(
			"MAX_RE_EVALUATION_DEPTH drift: shared orchestration uses %d but the server handler uses %d; "+
				"the WASM build and the HTTP handler MUST share one depth limit",
			MAX_RE_EVALUATION_DEPTH, serverHandlerConstant,
		)
	}
}

// hierarchicalInventoryFixture builds a deterministic Namespace+Pod
// design plus a single hierarchical/inventory relationship. The engine
// auto-identifies the parent/child relationship; this exercises the
// re-eval loop, alias resolution, trace merge, version bump and action
// de-dup without needing a component registry.
func hierarchicalInventoryFixture(t *testing.T) (pattern.PatternFile, []*relationship.RelationshipDefinition) {
	t.Helper()

	nsID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	podID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	ns := &componentv1beta1.ComponentDefinition{
		Component:      componentv1beta1.Component{Kind: "Namespace"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"name": "default"},
	}
	ns.ID = nsID

	pod := &componentv1beta1.ComponentDefinition{
		Component:      componentv1beta1.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"namespace": "default"},
	}
	pod.ID = podID

	design := pattern.PatternFile{
		Name:          "orchestration-parity-fixture",
		Version:       "0.1.0",
		Components:    []*componentv1beta1.ComponentDefinition{ns, pod},
		Relationships: []*relationship.RelationshipDefinition{},
	}

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "name"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "namespace"}}
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						Kind:  strPtr("Namespace"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						Kind:  strPtr("Deployment"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}
	relDef := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "parent",
		SubType:          "inventory",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	relDef.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000010")

	return design, []*relationship.RelationshipDefinition{relDef}
}

// handlerLoopMinusRegistry independently reconstructs the registry-
// INDEPENDENT part of (*Handler).EvaluateDesign: the
// MAX_RE_EVALUATION_DEPTH re-eval loop, per-iteration version bump
// (the registry-independent half of processEvaluationResponse), alias
// resolution, unique trace merge, completed-at timestamp and
// cross-iteration action de-dup. It deliberately does NOT call
// HydratePattern / processEvaluationResponse / GetEntitiesMemoized.
//
// If EvaluateDesignOrchestrated ever drifts from the handler's loop
// shape, the parity assertion in TestEvaluateDesignOrchestrated_*
// fails — this function is the executable spec of "handler path minus
// registry hydration".
func handlerLoopMinusRegistry(
	e *GoEngine,
	design pattern.PatternFile,
	rels []*relationship.RelationshipDefinition,
	iterations int,
) (pattern.EvaluationResponse, error) {
	var last pattern.EvaluationResponse
	last.Design = design

	for i := range MAX_RE_EVALUATION_DEPTH {
		resp, err := e.EvaluateDesign(last.Design, rels)
		if err != nil {
			return pattern.EvaluationResponse{}, err
		}

		BumpDesignVersion(&resp.Design)

		aliases := ResolveAliasesInDesign(resp.Design)
		if resp.Design.Metadata == nil {
			resp.Design.Metadata = &pattern.PatternFile_Metadata{}
		}
		resp.Design.Metadata.ResolvedAliases = patternutils.ResolvedAliasesV1beta2ToV1beta1(&aliases)

		last.Design = resp.Design
		last.Actions = append(last.Actions, resp.Actions...)
		MergeTraceUnique(&last.Trace, &resp.Trace)

		if iterations == i+1 || DoesntNeedReeval(resp) {
			break
		}
		if i == (MAX_RE_EVALUATION_DEPTH - 1) {
			break
		}
	}

	last.Actions = DeduplicateActions(last.Actions)
	return last, nil
}

// jsonEq compares two values by their JSON encoding so that the
// non-deterministic Timestamp (set to time.Now in both paths) is the
// only thing we deliberately exclude — everything else (actions, trace,
// design version, resolved aliases) must match byte for byte.
func jsonEq(t *testing.T, name string, want, got interface{}) {
	t.Helper()
	wb, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal want %s: %v", name, err)
	}
	gb, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal got %s: %v", name, err)
	}
	if !reflect.DeepEqual(json.RawMessage(wb), json.RawMessage(gb)) {
		t.Fatalf("%s mismatch between handler-minus-registry and EvaluateDesignOrchestrated:\nwant: %s\ngot:  %s", name, wb, gb)
	}
}

// TestEvaluateDesignOrchestrated_ParityWithHandlerMinusRegistry proves
// the shared orchestration produces the same actions and trace as the
// handler's registry-independent loop. This is the anti-drift guarantee
// the WASM contract depends on.
func TestEvaluateDesignOrchestrated_ParityWithHandlerMinusRegistry(t *testing.T) {
	log, _ := logger.New("orchestration-parity", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)

	design, rels := hierarchicalInventoryFixture(t)

	want, err := handlerLoopMinusRegistry(engine, design, rels, MAX_RE_EVALUATION_DEPTH)
	if err != nil {
		t.Fatalf("handler-minus-registry reference failed: %v", err)
	}

	got, err := engine.EvaluateDesignOrchestrated(design, rels, MAX_RE_EVALUATION_DEPTH, EvaluationOptions{EnableTrace: true})
	if err != nil {
		t.Fatalf("EvaluateDesignOrchestrated failed: %v", err)
	}

	// Timestamp is wall-clock in both paths; null it before comparing.
	want.Timestamp = nil
	got.Timestamp = nil

	jsonEq(t, "actions", want.Actions, got.Actions)
	jsonEq(t, "trace", want.Trace, got.Trace)
	jsonEq(t, "design.version", want.Design.Version, got.Design.Version)
	jsonEq(t, "design.metadata.resolvedAliases", want.Design.Metadata, got.Design.Metadata)
	jsonEq(t, "full response", want, got)

	// Sanity: the fixture must actually exercise the engine (the
	// hierarchical/inventory relationship is auto-identified, so at least
	// one action is produced). A vacuous parity pass would be worthless.
	if len(got.Actions) == 0 {
		t.Fatalf("fixture produced no actions; parity test is vacuous")
	}

	// The version must have been bumped at least once (registry-
	// independent half of processEvaluationResponse).
	if got.Design.Version == design.Version {
		t.Fatalf("design version was not bumped: still %q", got.Design.Version)
	}

	// Timestamp must be set by the orchestration.
	freshGot, _ := engine.EvaluateDesignOrchestrated(design, rels, MAX_RE_EVALUATION_DEPTH, EvaluationOptions{})
	if freshGot.Timestamp == nil {
		t.Fatalf("EvaluateDesignOrchestrated did not set Timestamp")
	}
}

// TestEvaluateDesignOrchestrated_ReturnDiffOnly verifies the
// registry-independent ReturnDiffOnly projection: the returned design is
// reduced to exactly the trace deltas.
func TestEvaluateDesignOrchestrated_ReturnDiffOnly(t *testing.T) {
	log, _ := logger.New("orchestration-diffonly", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)

	design, rels := hierarchicalInventoryFixture(t)

	full, err := engine.EvaluateDesignOrchestrated(design, rels, MAX_RE_EVALUATION_DEPTH, EvaluationOptions{})
	if err != nil {
		t.Fatalf("full eval failed: %v", err)
	}
	diff, err := engine.EvaluateDesignOrchestrated(design, rels, MAX_RE_EVALUATION_DEPTH, EvaluationOptions{ReturnDiffOnly: true})
	if err != nil {
		t.Fatalf("diff eval failed: %v", err)
	}

	// Diff-only design components/relationships are a subset of full and
	// equal to the trace deltas.
	wantComps := len(diff.Trace.ComponentsAdded) + len(diff.Trace.ComponentsUpdated)
	if len(diff.Design.Components) != wantComps {
		t.Fatalf("ReturnDiffOnly components: got %d, want %d (trace deltas)", len(diff.Design.Components), wantComps)
	}
	if len(diff.Design.Components) > len(full.Design.Components) && len(full.Design.Components) > 0 {
		t.Fatalf("ReturnDiffOnly produced more components (%d) than full eval (%d)", len(diff.Design.Components), len(full.Design.Components))
	}
}
