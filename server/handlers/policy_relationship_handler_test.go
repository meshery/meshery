package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/component"
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
// that ComponentFilter.Get finds it and returns *v1beta3.ComponentDefinition,
// which triggers the *v1beta1.ComponentDefinition type-assertion failure.
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
	}
	id, err := comp.GenerateID()
	require.NoError(t, err)
	comp.ID = id
	_, _, err = rm.RegisterEntity(conn, &comp)
	require.NoError(t, err, "seedTestComponent: RegisterEntity failed")
}

// Regression test for #18915: nil pointer dereference in processEvaluationResponse.
//
// Root cause: ComponentFilter.Get (meshkit v1beta1 filter package) returns
// *v1beta3.ComponentDefinition, but the handler asserts *v1beta1.ComponentDefinition.
// Without the guard the assertion silently returns nil=(*v1beta1.ComponentDefinition)(nil)
// and the next dereference causes a SIGSEGV.
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

	t.Run("type assertion failure routes to unknownComponents", func(t *testing.T) {
		t.Parallel()
		// !ok guard: registry returns *v1beta3.ComponentDefinition (the real production
		// type), which fails the *v1beta1.ComponentDefinition assertion in the handler.
		// Without the guard this is the exact path that caused the production SIGSEGV.
		rm, _ := newTestRegistryManager(t)
		seedTestComponent(t, rm) // seeds via RegisterEntity — no raw SQL
		var got []*component.ComponentDefinition
		require.NotPanics(t, func() {
			got = processEvaluationResponse(rm, pattern.EvaluationRequest{}, makeResp("test-job"))
		}, "must not panic when type assertion on registry entity fails")
		require.Len(t, got, 1)
		assert.Equal(t, "Job", got[0].Component.Kind)
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
