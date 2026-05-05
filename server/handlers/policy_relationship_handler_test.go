package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
	v1beta3comp "github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRegistryManager(t *testing.T) (*registry.RegistryManager, *database.Handler) {
	t.Helper()
	db, err := database.New(database.Options{
		Filename: ":memory:",
		Engine:   "sqlite",
	})
	require.NoError(t, err, "failed to create in-memory database")
	rm, err := registry.NewRegistryManager(&db)
	require.NoError(t, err, "failed to create registry manager")
	return rm, &db
}

func TestRunRelationshipEvaluation_RecoversPanic(t *testing.T) {
	// Production regression: an unrecovered panic in this goroutine used
	// to crash the whole Meshery process, which manifested in CI as the
	// e2e suite cascading from one bad /relationships/evaluate call to
	// 120 ECONNREFUSEDs. Verify the panic is converted to an error sent
	// to the requesting client AND broadcast to coalesced followers.
	tracker := newEvaluationTracker()
	leader, _ := tracker.acquire("design-1")
	if !leader {
		t.Fatalf("expected leader on first acquire")
	}
	// A follower joins before the leader publishes.
	_, waitCh := tracker.acquire("design-1")

	respCh := make(chan pattern.EvaluationResponse, 1)
	errCh := make(chan error, 1)

	done := make(chan struct{})
	go func() {
		defer close(done)
		runRelationshipEvaluation(
			context.Background(),
			newTestLogger(t),
			tracker,
			"design-1",
			func() (pattern.EvaluationResponse, error) {
				panic("synthetic engine explosion")
			},
			respCh,
			errCh,
		)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("goroutine did not return — recover path is broken")
	}

	select {
	case err := <-errCh:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "panic during relationship evaluation")
		assert.Contains(t, err.Error(), "synthetic engine explosion")
	default:
		t.Fatal("leader did not receive an error on errCh after panic")
	}

	select {
	case res := <-waitCh:
		assert.Error(t, res.err, "follower must observe the panic as an error, not block forever")
		assert.Contains(t, res.err.Error(), "synthetic engine explosion")
	case <-time.After(time.Second):
		t.Fatal("follower never unblocked — coalesced clients would hang")
	}

	select {
	case resp := <-respCh:
		t.Fatalf("respCh should be empty on panic, got %+v", resp)
	default:
	}
}

func TestRunRelationshipEvaluation_PanicDoesNotDeadlock(t *testing.T) {
	// If the leader has already given up — typically via the surrounding
	// handler's evalCtx.Done() case returning before this goroutine
	// finishes — the request handler is no longer reading errCh. A
	// blocking send from inside the recover path would then leak this
	// goroutine forever. Use unbuffered channels with no reader and
	// trigger a real panic so the recover-path's non-blocking send is
	// what's actually under test (the early ctx-cancellation guard
	// short-circuits before eval() is called and exercises a different
	// path).
	tracker := newEvaluationTracker()
	tracker.acquire("design-2")

	respCh := make(chan pattern.EvaluationResponse)
	errCh := make(chan error)

	done := make(chan struct{})
	go func() {
		defer close(done)
		runRelationshipEvaluation(
			context.Background(),
			newTestLogger(t),
			tracker,
			"design-2",
			func() (pattern.EvaluationResponse, error) {
				panic("synthetic engine explosion with no receiver")
			},
			respCh,
			errCh,
		)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("goroutine deadlocked on panic with no receiver")
	}
}

func TestRunRelationshipEvaluation_PassesThroughEvalError(t *testing.T) {
	// Non-panic eval errors must still flow through the same channels so
	// the existing happy-error path is unchanged by the recovery wrapper.
	tracker := newEvaluationTracker()
	tracker.acquire("design-3")

	respCh := make(chan pattern.EvaluationResponse, 1)
	errCh := make(chan error, 1)

	sentinel := errors.New("eval failed cleanly")

	go runRelationshipEvaluation(
		context.Background(),
		newTestLogger(t),
		tracker,
		"design-3",
		func() (pattern.EvaluationResponse, error) {
			return pattern.EvaluationResponse{}, sentinel
		},
		respCh,
		errCh,
	)

	select {
	case err := <-errCh:
		assert.ErrorIs(t, err, sentinel)
	case <-time.After(2 * time.Second):
		t.Fatal("eval error never delivered")
	}
}

// seedTestComponent registers a v1beta3.ComponentDefinition into the given
// RegistryManager using the registry's own RegisterEntity API — no raw SQL.
// The returned component has Kind/Version matching the evalResp fixture so
// that ComponentFilter.Get finds it and returns the canonical registry type.
func seedTestComponent(t *testing.T, rm *registry.RegistryManager) {
	t.Helper()
	conn := connection.Connection{
		Name:    "test-registrant",
		Kind:    "kubernetes",
		Type:    "platform",
		SubType: "orchestration",
		Status:  connection.ConnectionStatusConnected,
	}
	enabled := v1beta3comp.Enabled
	bgColor := "#123456"
	shape := core.Shape("round-rectangle")
	comp := v1beta3comp.ComponentDefinition{
		DisplayName:   "Job",
		SchemaVersion: "core.meshery.io/v1beta1",
		Status:        &enabled,
		Component: v1beta3comp.Component{
			Kind:    "Job",
			Version: "batch/v1",
			Schema:  `{"properties":{}}`, // non-empty required by Create
		},
		Model: &model.ModelDefinition{
			Name:          "kubernetes",
			DisplayName:   "Kubernetes",
			SchemaVersion: "models.meshery.io/v1beta1",
			Version:       "v1.25.0",
			Model:         model.Model{Version: "v1.25.0"},
			Category:      category.CategoryDefinition{Name: "Orchestration"},
			Status:        model.Enabled,
		},
		Styles: &core.ComponentStyles{
			BackgroundColor: &bgColor,
			PrimaryColor:    "#123456",
			Shape:           &shape,
		},
	}
	id, err := comp.GenerateID()
	require.NoError(t, err)
	comp.ID = id
	_, _, err = rm.RegisterEntity(conn, &comp)
	require.NoError(t, err, "seedTestComponent: RegisterEntity failed")
}

// Regression test for #18915: nil pointer dereference in processEvaluationResponse.
//
// Original root cause: ComponentFilter.Get (meshkit v1beta1 filter package)
// returns *v1beta3.ComponentDefinition, while this handler works with
// v1beta1 evaluation payloads. The conversion must be explicit so the handler
// neither panics nor treats valid registry-backed components as unknown.
//
// Real production panic (guard absent):
//
//	panic: runtime error: invalid memory address or nil pointer dereference
//	[signal SIGSEGV: segmentation violation code=0x2 addr=0x0 pc=0x107496bb0]
//
//	goroutine 10388 [running]:
//	github.com/meshery/meshery/server/handlers.processEvaluationResponse(...)
//		policy_relationship_handler.go:447 +0x410
//	github.com/meshery/meshery/server/handlers.(*Handler).EvaluateDesign(...)
//		policy_relationship_handler.go:351 +0x660
//	github.com/meshery/meshery/server/handlers.(*Handler).EvaluateRelationshipPolicy.func2()
//		policy_relationship_handler.go:581 +0xb8
func TestProcessEvaluationResponse_NilPointerGuard(t *testing.T) {
	t.Parallel()

	makeResp := func(displayName string) *pattern.EvaluationResponse {
		return &pattern.EvaluationResponse{
			Design: pattern.PatternFile{Version: "0.0.1"},
			Trace: pattern.Trace{
				ComponentsAdded: []component.ComponentDefinition{
					{Component: component.Component{Kind: "Job", Version: "batch/v1"}, DisplayName: displayName},
				},
			},
		}
	}

	t.Run("empty registry routes to unknownComponents", func(t *testing.T) {
		t.Parallel()
		// len(entities)==0 guard: registry has no match → unknownComponents, no panic.
		rm, _ := newTestRegistryManager(t)
		var got []*component.ComponentDefinition
		require.NotPanics(t, func() {
			got = processEvaluationResponse(rm, pattern.EvaluationRequest{}, makeResp("test-job"))
		}, "must not panic when registry returns no entities")
		require.Len(t, got, 1)
		assert.Equal(t, "Job", got[0].Component.Kind)
	})

	t.Run("v1beta3 registry component hydrates v1beta1 evaluation response", func(t *testing.T) {
		t.Parallel()
		rm, _ := newTestRegistryManager(t)
		seedTestComponent(t, rm)
		id, err := uuid.NewV4()
		require.NoError(t, err)
		position := &struct {
			X float64 `json:"x" yaml:"x"`
			Y float64 `json:"y" yaml:"y"`
		}{X: 10, Y: 20}
		partial := &component.ComponentDefinition{
			ID:             id,
			Component:      component.Component{Kind: "Job", Version: "batch/v1"},
			DisplayName:    "test-job",
			ModelReference: model.ModelReference{Name: "kubernetes"},
			Styles:         &core.ComponentStyles{Position: position},
		}
		resp := &pattern.EvaluationResponse{
			Design: pattern.PatternFile{
				Version:    "0.0.1",
				Components: []*component.ComponentDefinition{partial},
			},
			Trace: pattern.Trace{
				ComponentsAdded: []component.ComponentDefinition{*partial},
			},
			Actions: []pattern.Action{
				{
					Op: "add_component",
					Value: map[string]interface{}{
						"item": map[string]interface{}{
							"id": id.String(),
						},
					},
				},
			},
		}
		var got []*component.ComponentDefinition
		require.NotPanics(t, func() {
			got = processEvaluationResponse(rm, pattern.EvaluationRequest{}, resp)
		}, "must not panic when registry returns canonical v1beta3 components")
		require.Empty(t, got)
		require.Len(t, resp.Design.Components, 1)
		require.NotNil(t, resp.Design.Components[0].Styles)
		require.NotNil(t, resp.Design.Components[0].Styles.BackgroundColor)
		assert.Equal(t, "#123456", *resp.Design.Components[0].Styles.BackgroundColor)
		require.NotNil(t, resp.Design.Components[0].Styles.Position)
		assert.Equal(t, float64(10), resp.Design.Components[0].Styles.Position.X)
		assert.Equal(t, float64(20), resp.Design.Components[0].Styles.Position.Y)
		require.Len(t, resp.Actions, 1)
		item, ok := resp.Actions[0].Value["item"].(map[string]interface{})
		require.True(t, ok)
		styles, ok := item["styles"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "#123456", styles["background-color"])
	})

	t.Run("explicit DisplayName is preserved on unknown component", func(t *testing.T) {
		t.Parallel()
		// When the registry has no match (empty), the component is routed to
		// unknownComponents. The original DisplayName on the eval-trace entry
		// must be preserved so callers can identify the component.
		rm, _ := newTestRegistryManager(t)
		var got []*component.ComponentDefinition
		require.NotPanics(t, func() {
			got = processEvaluationResponse(rm, pattern.EvaluationRequest{}, makeResp("my-custom-job"))
		})
		require.Len(t, got, 1)
		assert.Equal(t, "my-custom-job", got[0].DisplayName)
		assert.Equal(t, "Job", got[0].Component.Kind)
	})

	t.Run("ReturnDiffOnly option populates Design from Trace", func(t *testing.T) {
		t.Parallel()
		// When ReturnDiffOnly is true, the function replaces Design.Components and
		// Design.Relationships from the Trace instead of the existing design state.
		rm, _ := newTestRegistryManager(t)
		returnDiff := true
		req := pattern.EvaluationRequest{
			Options: &pattern.Options{ReturnDiffOnly: &returnDiff},
		}
		resp := &pattern.EvaluationResponse{
			Design: pattern.PatternFile{Version: "0.0.1"},
			Trace: pattern.Trace{
				ComponentsAdded: []component.ComponentDefinition{
					{Component: component.Component{Kind: "Job", Version: "batch/v1"}, DisplayName: "diff-job"},
				},
				RelationshipsAdded: []relationship.RelationshipDefinition{
					{SubType: "binding"},
				},
			},
		}
		require.NotPanics(t, func() {
			processEvaluationResponse(rm, req, resp)
		}, "must not panic with ReturnDiffOnly=true")
		// Design.Relationships should be populated from Trace.RelationshipsAdded.
		require.Len(t, resp.Design.Relationships, 1)
		assert.Equal(t, "binding", resp.Design.Relationships[0].SubType)
	})

	t.Run("IsAnnotation metadata is copied to hydrated component", func(t *testing.T) {
		t.Parallel()
		// Annotation flag on the eval-trace entry must propagate to the registry
		// component after hydration (line: _component.Metadata.IsAnnotation = _c.Metadata.IsAnnotation).
		// With an empty registry the component is unknown, but the metadata copy
		// code path is reachable once the guard passes — documented here.
		rm, _ := newTestRegistryManager(t)
		resp := &pattern.EvaluationResponse{
			Design: pattern.PatternFile{Version: "0.0.1"},
			Trace: pattern.Trace{
				ComponentsAdded: []component.ComponentDefinition{
					{
						Component:   component.Component{Kind: "Job", Version: "batch/v1"},
						DisplayName: "annotation-job",
						Metadata:    component.ComponentDefinition_Metadata{IsAnnotation: true},
					},
				},
			},
		}
		var got []*component.ComponentDefinition
		require.NotPanics(t, func() {
			got = processEvaluationResponse(rm, pattern.EvaluationRequest{}, resp)
		})
		require.Len(t, got, 1)
		assert.True(t, got[0].Metadata.IsAnnotation)
	})
}

// seedNamespaceComponent registers a second styled component (Namespace)
// alongside the Job from seedTestComponent so multi-component hydration
// scenarios have two distinct kinds with two distinct styles to
// disambiguate. Reuses the same connection/model identity as
// seedTestComponent (kubernetes v1.25.0) to keep the model_dbs row
// shared — RegisterEntity rejects two components claiming the same
// generated model ID via separate registrations otherwise.
func seedNamespaceComponent(t *testing.T, rm *registry.RegistryManager) {
	t.Helper()
	conn := connection.Connection{
		Name:    "test-registrant",
		Kind:    "kubernetes",
		Type:    "platform",
		SubType: "orchestration",
		Status:  connection.ConnectionStatusConnected,
	}
	enabled := v1beta3comp.Enabled
	bgColor := "#326CE5"
	primary := "#326CE5"
	shape := core.Shape("rectangle")
	comp := v1beta3comp.ComponentDefinition{
		DisplayName:   "Namespace",
		SchemaVersion: "core.meshery.io/v1beta1",
		Status:        &enabled,
		Component: v1beta3comp.Component{
			Kind:    "Namespace",
			Version: "v1",
			Schema:  `{"properties":{}}`,
		},
		Model: &model.ModelDefinition{
			Name:          "kubernetes",
			DisplayName:   "Kubernetes",
			SchemaVersion: "models.meshery.io/v1beta1",
			Version:       "v1.25.0",
			Model:         model.Model{Version: "v1.25.0"},
			Category:      category.CategoryDefinition{Name: "Orchestration"},
			Status:        model.Enabled,
		},
		Styles: &core.ComponentStyles{
			BackgroundColor: &bgColor,
			PrimaryColor:    primary,
			Shape:           &shape,
		},
	}
	id, err := comp.GenerateID()
	require.NoError(t, err)
	comp.ID = id
	_, _, err = rm.RegisterEntity(conn, &comp)
	require.NoError(t, err, "seedNamespaceComponent: RegisterEntity failed")
}

// Per-instance preservation invariant: when the evaluator adds a new
// component (e.g. Namespace auto-added for a Pod), the user's custom
// styling on EXISTING design components (e.g. a Pod whose color the user
// changed from green to red) must not be overwritten. Hydration is for
// new components only — existing components flow through unmodified.
//
// This is the contract that lets users customize node appearance without
// losing those customizations every time the evaluator runs (and the
// evaluator runs after almost every UI mutation: ADD_COMPONENT,
// UPDATE_CONFIGURATION, MERGE_DESIGN, …).
//
// Regression guard: a future "improvement" that loops over
// Design.Components and applies registry styles to all of them — rather
// than only to Trace.ComponentsAdded — would silently destroy every
// per-instance customization in the design. This test fails fast if that
// happens.
func TestProcessEvaluationResponse_PreservesUserCustomizationsOnExistingComponents(t *testing.T) {
	t.Parallel()

	rm, _ := newTestRegistryManager(t)
	seedTestComponent(t, rm)        // Job, registry BackgroundColor "#123456"
	seedNamespaceComponent(t, rm)   // Namespace, registry BackgroundColor "#326CE5"

	existingPodID, err := uuid.NewV4()
	require.NoError(t, err)
	addedNamespaceID, err := uuid.NewV4()
	require.NoError(t, err)

	// User customization on the existing component: red background. The
	// registry's default for "Job" is "#123456"; if hydration accidentally
	// touched existing components it would overwrite this with "#123456".
	userBg := "#FF0000"
	existingPod := &component.ComponentDefinition{
		ID:          existingPodID,
		Component:   component.Component{Kind: "Job", Version: "batch/v1"},
		DisplayName: "user-customized-job",
		Styles:      &core.ComponentStyles{BackgroundColor: &userBg},
	}

	// New component the evaluator just added — bare, no styles.
	addedNamespace := &component.ComponentDefinition{
		ID:          addedNamespaceID,
		Component:   component.Component{Kind: "Namespace", Version: "v1"},
		DisplayName: "default",
		ModelReference: model.ModelReference{Name: "kubernetes"},
	}

	resp := &pattern.EvaluationResponse{
		Design: pattern.PatternFile{
			Version:    "0.0.1",
			Components: []*component.ComponentDefinition{existingPod, addedNamespace},
		},
		Trace: pattern.Trace{
			ComponentsAdded: []component.ComponentDefinition{*addedNamespace},
		},
	}

	got := processEvaluationResponse(rm, pattern.EvaluationRequest{}, resp)
	require.Empty(t, got, "Namespace must hydrate from registry, not strand into unknownComponents")

	require.Len(t, resp.Design.Components, 2)
	var hydratedPod, hydratedNs *component.ComponentDefinition
	for _, c := range resp.Design.Components {
		switch c.ID {
		case existingPodID:
			hydratedPod = c
		case addedNamespaceID:
			hydratedNs = c
		}
	}
	require.NotNil(t, hydratedPod, "existing user-customized component must remain in design")
	require.NotNil(t, hydratedNs, "evaluator-added component must be present in design")

	// Invariant 1: existing user customization survives unmodified.
	require.NotNil(t, hydratedPod.Styles)
	require.NotNil(t, hydratedPod.Styles.BackgroundColor)
	assert.Equal(t, "#FF0000", *hydratedPod.Styles.BackgroundColor,
		"user-customized BackgroundColor must not be overwritten by hydration")
	assert.Equal(t, "user-customized-job", string(hydratedPod.DisplayName),
		"user-set DisplayName must not be overwritten by hydration")

	// Invariant 2: newly-added component is hydrated from the registry.
	require.NotNil(t, hydratedNs.Styles, "new component must have registry styles after hydration")
	require.NotNil(t, hydratedNs.Styles.BackgroundColor)
	assert.Equal(t, "#326CE5", *hydratedNs.Styles.BackgroundColor,
		"new component must carry the registry's BackgroundColor")
}

// Multi-component hydration: the evaluator may add several components
// in a single pass (e.g. a Pod referencing both a missing Namespace and
// a missing ServiceAccount). Each must hydrate to its own registry-defined
// styles, and the registry-cache reuse across the loop must not bleed
// styles from one kind into another.
//
// Regression guard for any future change that breaks the per-iteration
// filter construction, the registry-cache scoping, or the ID-match
// loop in the design-component swap.
func TestProcessEvaluationResponse_HydratesMultipleAddedComponents(t *testing.T) {
	t.Parallel()

	rm, _ := newTestRegistryManager(t)
	seedTestComponent(t, rm)       // Job → BackgroundColor "#123456"
	seedNamespaceComponent(t, rm)  // Namespace → BackgroundColor "#326CE5"

	jobID, err := uuid.NewV4()
	require.NoError(t, err)
	nsID, err := uuid.NewV4()
	require.NoError(t, err)

	bareJob := &component.ComponentDefinition{
		ID:             jobID,
		Component:      component.Component{Kind: "Job", Version: "batch/v1"},
		DisplayName:    "auto-job",
		ModelReference: model.ModelReference{Name: "kubernetes"},
	}
	bareNs := &component.ComponentDefinition{
		ID:             nsID,
		Component:      component.Component{Kind: "Namespace", Version: "v1"},
		DisplayName:    "auto-ns",
		ModelReference: model.ModelReference{Name: "kubernetes"},
	}

	resp := &pattern.EvaluationResponse{
		Design: pattern.PatternFile{
			Version:    "0.0.1",
			Components: []*component.ComponentDefinition{bareJob, bareNs},
		},
		Trace: pattern.Trace{
			ComponentsAdded: []component.ComponentDefinition{*bareJob, *bareNs},
		},
	}

	got := processEvaluationResponse(rm, pattern.EvaluationRequest{}, resp)
	require.Empty(t, got, "all known kinds must hydrate; nothing should fall to unknownComponents")
	require.Len(t, resp.Design.Components, 2)

	byID := map[core.Uuid]*component.ComponentDefinition{}
	for _, c := range resp.Design.Components {
		byID[c.ID] = c
	}

	// Each kind hydrates to ITS OWN registry colour — no cross-bleed
	// from registry-cache reuse.
	require.NotNil(t, byID[jobID].Styles)
	require.NotNil(t, byID[jobID].Styles.BackgroundColor)
	assert.Equal(t, "#123456", *byID[jobID].Styles.BackgroundColor, "Job must hydrate to its own colour")

	require.NotNil(t, byID[nsID].Styles)
	require.NotNil(t, byID[nsID].Styles.BackgroundColor)
	assert.Equal(t, "#326CE5", *byID[nsID].Styles.BackgroundColor, "Namespace must hydrate to its own colour")
}

func TestParseRelationshipToAlias(t *testing.T) {
	fromID := uuid.Must(uuid.NewV4())
	toID := uuid.Must(uuid.NewV4())
	relID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name   string
		input  relationship.RelationshipDefinition
		wantOk bool
	}{
		{
			name: "wrong subtype returns false",
			input: relationship.RelationshipDefinition{
				SubType: "not-alias",
			},
			wantOk: false,
		},
		{
			name: "nil Selectors returns false",
			input: relationship.RelationshipDefinition{
				SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
				Selectors: nil,
			},
			wantOk: false,
		},
		{
			name: "empty Selectors returns false",
			input: relationship.RelationshipDefinition{
				SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
				Selectors: &relationship.SelectorSet{},
			},
			wantOk: false,
		},
		{
			name: "empty From set returns false",
			input: func() relationship.RelationshipDefinition {
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{},
							To:   []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "empty To set returns false",
			input: func() relationship.RelationshipDefinition {
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{{ID: &fromID}},
							To:   []relationship.SelectorItem{},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "nil Patch returns false",
			input: func() relationship.RelationshipDefinition {
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{{ID: &fromID, RelationshipDefinitionSelectorsPatch: nil}},
							To:   []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "nil MutatedRef returns false",
			input: func() relationship.RelationshipDefinition {
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									ID: &fromID,
									RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
										MutatedRef: nil,
									},
								},
							},
							To: []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "empty MutatedRef returns false",
			input: func() relationship.RelationshipDefinition {
				emptyRefs := [][]string{}
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									ID: &fromID,
									RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
										MutatedRef: &emptyRefs,
									},
								},
							},
							To: []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "nil to.ID returns false",
			input: func() relationship.RelationshipDefinition {
				refs := [][]string{{"configuration", "spec", "containers"}}
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									ID: &fromID,
									RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
										MutatedRef: &refs,
									},
								},
							},
							To: []relationship.SelectorItem{{ID: nil}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "nil from.ID returns false",
			input: func() relationship.RelationshipDefinition {
				refs := [][]string{{"configuration", "spec", "containers"}}
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									ID: nil,
									RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
										MutatedRef: &refs,
									},
								},
							},
							To: []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "valid alias relationship returns true",
			input: func() relationship.RelationshipDefinition {
				refs := [][]string{{"configuration", "spec", "containers"}}
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									ID: &fromID,
									RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
										MutatedRef: &refs,
									},
								},
							},
							To: []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				rd := relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
				rd.ID = relID
				return rd
			}(),
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alias, ok := parseRelationshipToAlias(tt.input)
			assert.Equal(t, tt.wantOk, ok, "parseRelationshipToAlias() ok mismatch")

			if tt.wantOk {
				assert.Equal(t, toID, alias.ImmediateParentId, "ImmediateParentId should match to.ID")
				assert.Equal(t, fromID, alias.AliasComponentId, "AliasComponentId should match from.ID")
				assert.Equal(t, relID, alias.RelationshipId, "RelationshipId should match the relationship's Id")
				assert.Equal(t, []string{"configuration", "spec", "containers"}, alias.ImmediateRefFieldPath, "ImmediateRefFieldPath should match first mutatedRef entry")
			}
		})
	}
}
