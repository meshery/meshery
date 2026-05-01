package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/stretchr/testify/assert"
)

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
