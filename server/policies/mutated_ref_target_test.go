package policies

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// mutationLifecycleCase describes one forward-mutation / re-evaluation / deletion
// lifecycle over a single edge non-binding relationship whose mutatedRef shape is
// the variable under test.
type mutationLifecycleCase struct {
	name       string
	mutatorRef []string
	mutatedRef []string
	// readBack returns the value the engine is expected to have written, read
	// from the object the mutatedRef addresses.
	readBack func(*component.ComponentDefinition) interface{}
	want     interface{}
}

// buildMutationDesign wires a Rule (mutator) and an EventBus (mutated) together
// with one approved edge/non-binding/network relationship.
func buildMutationDesign(mutatorRef, mutatedRef []string, mutatorValue string) (*pattern.PatternFile, uuid.UUID, uuid.UUID, uuid.UUID) {
	ruleID, _ := uuid.FromString("00000000-0000-0000-0000-0000000000a1")
	busID, _ := uuid.FromString("00000000-0000-0000-0000-0000000000a2")
	relID, _ := uuid.FromString("00000000-0000-0000-0000-0000000000a3")

	rule := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Rule"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		DisplayName:    "rule-1",
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"eventBusName": mutatorValue,
				"published":    true,
				"labels":       map[string]interface{}{"app": mutatorValue},
			},
		},
	}
	rule.ID = ruleID

	bus := &component.ComponentDefinition{
		Component:      component.Component{Kind: "EventBus"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		DisplayName:    "event-bus-original",
		Configuration:  map[string]interface{}{},
	}
	bus.ID = busID

	mRef := relationship.MutatorRef{mutatorRef}
	tRef := relationship.MutatedRef{mutatedRef}
	status := relationship.RelationshipDefinitionStatus(StatusApproved)

	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{{
					ID:   &ruleID,
					Kind: strPtr("Rule"),
					RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
						MutatorRef: &mRef,
					},
				}},
				To: []relationship.SelectorItem{{
					ID:   &busID,
					Kind: strPtr("EventBus"),
					RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
						MutatedRef: &tRef,
					},
				}},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Edge,
		RelationshipType: "non-binding",
		SubType:          "network",
		Status:           &status,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID = relID

	design := makePatternFile(
		[]*component.ComponentDefinition{rule, bus},
		[]*relationship.RelationshipDefinition{rel},
	)
	return design, ruleID, busID, relID
}

func componentByID(design *pattern.PatternFile, id uuid.UUID) *component.ComponentDefinition {
	for _, c := range design.Components {
		if c.ID == id {
			return c
		}
	}
	return nil
}

func countMutationActions(actions []pattern.Action, id string) int {
	n := 0
	for _, a := range actions {
		if a.Op != UpdateComponentOp && a.Op != UpdateComponentConfigurationOp {
			continue
		}
		if v, ok := a.Value["id"].(string); ok && v == id {
			n++
		}
	}
	return n
}

func countCleanupActions(actions []pattern.Action, id string) int {
	n := 0
	for _, a := range actions {
		if a.Op != RemoveComponentConfigurationOp && a.Op != UpdateComponentOp && a.Op != UpdateComponentConfigurationOp {
			continue
		}
		if v, ok := a.Value["id"].(string); ok && v == id {
			n++
		}
	}
	return n
}

// roundTrip serializes and reloads a design, proving the observed state survives
// persistence rather than living only in an in-memory struct.
func roundTrip(t *testing.T, design pattern.PatternFile) pattern.PatternFile {
	t.Helper()
	raw, err := json.Marshal(design)
	if err != nil {
		t.Fatalf("marshal design: %v", err)
	}
	var reloaded pattern.PatternFile
	if err := json.Unmarshal(raw, &reloaded); err != nil {
		t.Fatalf("unmarshal design: %v", err)
	}
	return reloaded
}

// TestMutatedRefLifecycle drives the full production evaluation pipeline
// (GoEngine.EvaluateDesign -> applyActionsToDesign) over three mutatedRef shapes
// and asserts the same five-part contract for each:
//
//  1. the first evaluation writes the mutator's value to the object the
//     mutatedRef addresses,
//  2. re-evaluating the resulting design is a fixed point (no second action),
//  3. the value survives serialization and reload,
//  4. deleting the relationship detects the previous mutation and emits cleanup,
//  5. applying cleanup removes the mutated value.
func TestMutatedRefLifecycle(t *testing.T) {
	cases := []mutationLifecycleCase{
		{
			name:       "configuration_prefixed_control",
			mutatorRef: []string{"configuration", "spec", "eventBusName"},
			mutatedRef: []string{"configuration", "spec", "busRef"},
			readBack: func(c *component.ComponentDefinition) interface{} {
				return objectGetNested(c.Configuration, []string{"spec", "busRef"}, nil)
			},
			want: "bus-alpha",
		},
		{
			name:       "unprefixed_envelope_displayName",
			mutatorRef: []string{"configuration", "spec", "eventBusName"},
			mutatedRef: []string{"displayName"},
			readBack: func(c *component.ComponentDefinition) interface{} {
				return string(c.DisplayName)
			},
			want: "bus-alpha",
		},
		{
			name:       "unprefixed_nested_envelope_metadata_published",
			mutatorRef: []string{"configuration", "spec", "published"},
			mutatedRef: []string{"metadata", "published"},
			readBack: func(c *component.ComponentDefinition) interface{} {
				return c.Metadata.Published
			},
			want: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			log, _ := logger.New("mutated-ref", logger.Options{Format: logger.SyslogLogFormat})
			engine := NewGoEngine(log)

			design, _, busID, relID := buildMutationDesign(tc.mutatorRef, tc.mutatedRef, "bus-alpha")

			// Step 1-4: first evaluation and applied mutation.
			first, err := engine.EvaluateDesign(*design, nil)
			if err != nil {
				t.Fatalf("first evaluation failed: %v", err)
			}
			if got := countMutationActions(first.Actions, busID.String()); got != 1 {
				t.Fatalf("first evaluation: want exactly 1 mutation action for the mutated component, got %d (actions=%v)", got, first.Actions)
			}
			bus := componentByID(&first.Design, busID)
			if bus == nil {
				t.Fatal("mutated component missing from first evaluation result")
			}
			if got := tc.readBack(bus); !deepEqual(got, tc.want) {
				t.Errorf("after first evaluation: value at mutatedRef %v = %#v, want %#v", tc.mutatedRef, got, tc.want)
			}

			// Step 5-6: the second evaluation must be a fixed point.
			second, err := engine.EvaluateDesign(first.Design, nil)
			if err != nil {
				t.Fatalf("second evaluation failed: %v", err)
			}
			if got := countMutationActions(second.Actions, busID.String()); got != 0 {
				t.Errorf("second evaluation is not a fixed point: got %d redundant mutation action(s), want 0", got)
			}

			// Step: serialization and reload must not change path interpretation.
			reloaded := roundTrip(t, second.Design)
			third, err := engine.EvaluateDesign(reloaded, nil)
			if err != nil {
				t.Fatalf("post-reload evaluation failed: %v", err)
			}
			if got := countMutationActions(third.Actions, busID.String()); got != 0 {
				t.Errorf("evaluation after serialize/reload is not a fixed point: got %d redundant mutation action(s), want 0", got)
			}

			// Step 7-10: delete the relationship and require cleanup.
			deleted := roundTrip(t, third.Design)
			var target *relationship.RelationshipDefinition
			for _, r := range deleted.Relationships {
				if r.ID == relID {
					target = r
				}
			}
			if target == nil {
				t.Fatal("relationship missing before deletion step")
			}
			setRelStatus(target, StatusDeleted)

			final, err := engine.EvaluateDesign(deleted, nil)
			if err != nil {
				t.Fatalf("deletion evaluation failed: %v", err)
			}
			if got := countCleanupActions(final.Actions, busID.String()); got == 0 {
				t.Fatalf("relationship deletion produced no cleanup action for the mutated component (actions=%v)", final.Actions)
			}
			busAfter := componentByID(&final.Design, busID)
			if busAfter == nil {
				t.Fatal("mutated component missing from deletion result")
			}
			if got := tc.readBack(busAfter); deepEqual(got, tc.want) {
				t.Errorf("after cleanup: value at mutatedRef %v is still %#v, want it removed or reset", tc.mutatedRef, got)
			}
		})
	}
}

// TestUnprefixedMutatedRefReadWriteTargets isolates the read/write asymmetry
// itself: the path that patchMutatorsAction reads through and the object that
// applyActionsToDesign writes into must be the same object.
func TestUnprefixedMutatedRefReadWriteTargets(t *testing.T) {
	design, _, busID, _ := buildMutationDesign(
		[]string{"configuration", "spec", "eventBusName"},
		[]string{"displayName"},
		"bus-alpha",
	)

	actions := (&EdgeNonBindingPolicy{}).SideEffects(design.Relationships[0], design)
	if len(actions) != 1 {
		t.Fatalf("want 1 mutation action, got %d", len(actions))
	}
	if actions[0].Op != UpdateComponentOp {
		t.Fatalf("an unprefixed mutatedRef must produce %q, got %q", UpdateComponentOp, actions[0].Op)
	}

	applied := applyActionsToDesign(design, actions)
	bus := componentByID(applied, busID)

	if got := string(bus.DisplayName); got != "bus-alpha" {
		t.Errorf("write target: envelope displayName = %q, want %q", got, "bus-alpha")
	}
	if stray := objectGetNested(bus.Configuration, []string{"displayName"}, nil); stray != nil {
		t.Errorf("write target: configuration.displayName = %#v, want nil - an envelope path must not be written into configuration", stray)
	}

	// The read side must observe what the write side produced.
	readBack := configurationForComponentAtPath([]string{"displayName"}, bus, applied)
	if readBack != "bus-alpha" {
		t.Errorf("read target: configurationForComponentAtPath(displayName) = %#v, want %q", readBack, "bus-alpha")
	}
}

// TestConfigurationMutatedRefStaysInConfiguration is the opposite guard: a
// configuration-prefixed reference must not start writing to the envelope now
// that envelope writes are supported.
func TestConfigurationMutatedRefStaysInConfiguration(t *testing.T) {
	design, _, busID, _ := buildMutationDesign(
		[]string{"configuration", "spec", "eventBusName"},
		[]string{"configuration", "spec", "busRef"},
		"bus-alpha",
	)

	actions := (&EdgeNonBindingPolicy{}).SideEffects(design.Relationships[0], design)
	if len(actions) != 1 || actions[0].Op != UpdateComponentConfigurationOp {
		t.Fatalf("want 1 %q action, got %+v", UpdateComponentConfigurationOp, actions)
	}

	applied := applyActionsToDesign(design, actions)
	bus := componentByID(applied, busID)

	if got := objectGetNested(bus.Configuration, []string{"spec", "busRef"}, nil); got != "bus-alpha" {
		t.Errorf("configuration.spec.busRef = %#v, want %q", got, "bus-alpha")
	}
	if got := string(bus.DisplayName); got != "event-bus-original" {
		t.Errorf("a configuration-prefixed reference must not touch the envelope: displayName = %q, want %q", got, "event-bus-original")
	}
}

// TestEnvelopeMutatedRefDoesNotLeakIntoConfiguration pins the behaviour for an
// unprefixed reference naming a field the component envelope cannot hold -
// ComponentDefinition_Metadata is a closed struct with no "labels" member, so
// "metadata.labels" has no sink. The engine must not silently divert it into
// Configuration, which is what made such a reference appear to work while being
// unreadable on the next evaluation. Relationship definitions using this shape
// are themselves defective and belong to the definition-repair class tracked by
// #21490.
func TestEnvelopeMutatedRefDoesNotLeakIntoConfiguration(t *testing.T) {
	design, _, busID, _ := buildMutationDesign(
		[]string{"configuration", "spec", "labels"},
		[]string{"metadata", "labels"},
		"bus-alpha",
	)

	actions := (&EdgeNonBindingPolicy{}).SideEffects(design.Relationships[0], design)
	if len(actions) != 1 || actions[0].Op != UpdateComponentOp {
		t.Fatalf("want 1 %q action, got %+v", UpdateComponentOp, actions)
	}

	applied := applyActionsToDesign(design, actions)
	bus := componentByID(applied, busID)

	if stray := objectGetNested(bus.Configuration, []string{"metadata"}, nil); stray != nil {
		t.Errorf("configuration.metadata = %#v, want nil - an envelope reference must never be written into configuration", stray)
	}
}
