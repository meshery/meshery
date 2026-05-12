// Package policies implements Meshery's native Go relationship evaluation
// engine. It replaces the legacy Open Policy Agent (OPA) Rego rules with
// typed Go policies that operate directly on the v1beta1 pattern model and
// the v1beta2 relationship model.
//
// # Evaluation pipeline
//
// [GoEngine.EvaluateDesign] is the entry point. It filters the registered
// relationships down to those whose model appears in the design, then runs
// the pipeline in [GoEngine.evaluate]:
//
//  1. Validate. For every implicated existing relationship whose
//     components no longer exist in the design
//     (fromOrToComponentsDontExist), call
//     [RelationshipPolicy.IsInvalid]. When both conditions hold, the
//     relationship receives a status update to "deleted".
//  2. Identify. For every relationship in scope, call
//     [RelationshipPolicy.IdentifyRelationship]. Each returned relationship
//     is deduplicated against the existing design via
//     [RelationshipPolicy.AlreadyExists] and emitted as an add-relationship
//     action.
//  2a. Inventory additions. [identifyInventoryAdditions] auto-creates any
//      parent declarations a child references but the design does not
//      contain (Namespace for a namespaced resource, etc.). This mirrors
//      the Rego engine's identify_additions step.
//  3. Side effects. For every implicated relationship, approve identified
//     ones, clean up deleted ones, and apply policy-specific side effects
//     (configuration patches, alias node creation, ...).
//
// Each phase produces a flat slice of [PolicyAction] values. The engine
// applies them with [applyActionsToDesign] between phases so later phases
// see the mutated design. The final action list is converted to
// [pattern.Action] for the HTTP response and a [pattern.Trace] is built
// describing what changed.
//
// # The RelationshipPolicy interface
//
// Each built-in policy implements [RelationshipPolicy]:
//
//   - Identifier returns a stable policy name used for logging.
//   - IsImplicatedBy gates the policy on relationship kind/type/subType.
//   - IsInvalid decides whether an existing relationship should be deleted.
//   - AlreadyExists prevents duplicate identification.
//   - IdentifyRelationship proposes new relationships by walking the design.
//   - SideEffects produces configuration mutations or component
//     add/delete actions tied to the relationship's lifecycle.
//
// # Built-in policies
//
// Registered in [NewGoEngine] in this order:
//
//   - MatchLabelsPolicy           (policy_matchlabels.go)
//     Sibling relationships. Groups components that share equal values at
//     the configured Match.Refs path(s).
//   - AliasPolicy                 (policy_alias.go)
//     Hierarchical parent/alias. Creates an alias node that resolves to a
//     reference path inside another component (a container port, for
//     example).
//   - EdgeNonBindingPolicy        (policy_edge_network.go)
//     Logical edge with no configuration mutation.
//   - EdgeBindingPolicy           (policy_binding.go)
//     3-party binding (from, binding, to). Patches the bound component's
//     configuration with the mutator's value and reverses the patch on
//     deletion via [cleanupActionForPath].
//   - HierarchicalParentChildPolicy (policy_hierarchical.go)
//     Inventory parent/child. Connects a child whose configuration field
//     value equals its parent's mutator value (Namespace/Deployment, etc.).
//   - HierarchicalWalletPolicy    (policy_wallet.go)
//     Custodial parent/wallet relationships. Structurally identical to
//     HierarchicalParentChildPolicy (same matching helpers and mutator
//     patching), but gated to subType="wallet" and not subject to phase
//     2.5 inventory auto-additions. Used for ownership/custody semantics
//     where the parent must already exist in the design.
//
// # Determinism contract
//
// Two designs with the same components and relationships must produce
// identical evaluation output across runs and across processes. Anything
// that affects the result must be deterministic:
//
//   - Map iteration order must not leak into output. Where a policy ranges
//     over a map (binding kinds, matchLabel groups), the keys are sorted
//     before iteration.
//   - Relationship UUIDs are generated with [staticUUID], a UUIDv5 over a
//     canonical seed (json.Marshal of a typed struct or map). Never seed
//     with %v formatting of a struct that contains pointers; pointer
//     addresses are not stable.
//   - Selector / match-ref slices are ranged in full. Reading only [0]
//     drops multi-container Pods, multi-port containers, and multi-ref
//     bindings.
//
// # Adding a new policy
//
//  1. Add a policy_<name>.go file with a struct that implements
//     [RelationshipPolicy].
//  2. Register the struct in [NewGoEngine]. Order does not affect
//     correctness because phases run per-policy, but matching the existing
//     order keeps trace output stable.
//  3. Reuse the helpers in eval_rules.go and feasibility.go
//     (identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields,
//     feasibleRelationshipSelectorBetween, matchingMutators, ...) instead
//     of rebuilding selector traversal.
//  4. Seed relationship IDs with [staticUUID] over a canonical map. Do not
//     include pointers or unsorted slices in the seed.
//  5. Add coverage in engine_test.go and, for hot policies, in
//     bench_test.go.
package policies
