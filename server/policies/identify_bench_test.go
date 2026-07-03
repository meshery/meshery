package policies

import (
	"fmt"
	"testing"

	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// buildMatchlabelsBenchFixture returns n Pods sharing label values across a
// small number of groups so identifyMatchlabels has nontrivial bucketing work.
// Each pod has two labels (app, version) whose values come from a small pool,
// producing several non-trivial groups regardless of n.
func buildMatchlabelsBenchFixture(n int) (*pattern.PatternFile, *relationship.RelationshipDefinition) {
	comps := make([]*component.ComponentDefinition, 0, n)
	groups := n/5 + 1
	for i := 0; i < n; i++ {
		c := &component.ComponentDefinition{
			Component: component.Component{Kind: "Pod"},
			Configuration: map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app":     fmt.Sprintf("g-%d", i%groups),
						"version": fmt.Sprintf("v-%d", i%5),
					},
				},
			},
		}
		c.ID = staticUUID(fmt.Sprintf("pod-%d", i))
		comps = append(comps, c)
	}

	kind := "Pod"
	refs := [][]string{{"configuration", "metadata", "labels"}}
	relDef := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "sibling",
		Selectors: &relationship.SelectorSet{
			relationship.SelectorSetItem{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{{Kind: &kind, Match: &relationship.MatchSelector{Refs: &refs}}},
					To:   []relationship.SelectorItem{{Kind: &kind, Match: &relationship.MatchSelector{Refs: &refs}}},
				},
			},
		},
	}
	relDef.ID = staticUUID("rel-matchlabels-bench")

	return &pattern.PatternFile{Components: comps}, relDef
}

// BenchmarkIdentifyMatchlabels measures identifyMatchlabels as the component
// count grows, so the perf delta from the O(N²·F) → O(N·F) rewrite is visible.
func BenchmarkIdentifyMatchlabels(b *testing.B) {
	for _, n := range []int{10, 50, 200} {
		design, relDef := buildMatchlabelsBenchFixture(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = identifyMatchlabels(design, relDef)
			}
		})
	}
}

// buildBindingBenchFixture returns 3n components arranged so each RoleBinding
// fans out to multiple matching Roles, which is the case where the cubic
// (fromDecl × bindingDecl × toDecl) re-evaluation dominates pre-fix.
//
// Roles cycle through bindingFanout distinct names so each name covers
// n/bindingFanout Roles. RoleBindings cycle through the same names, so every
// RoleBinding matches n/bindingFanout Roles. Each RoleBinding's subjects[0].name
// targets exactly one matching ServiceAccount.
//
// Pre-fix work: O(n²) outer (Role, RB) valid pairs × n inner toDecl scans = O(n³).
// Post-fix work: O(n²) for (RB, SA) precompute + O(n²) outer = O(n²) total.
func buildBindingBenchFixture(n int) (*pattern.PatternFile, *relationship.RelationshipDefinition) {
	const bindingFanout = 5
	comps := make([]*component.ComponentDefinition, 0, 3*n)
	for i := 0; i < n; i++ {
		role := &component.ComponentDefinition{
			Component:     component.Component{Kind: "Role"},
			Configuration: map[string]interface{}{"name": fmt.Sprintf("role-%d", i%bindingFanout)},
		}
		role.ID = staticUUID(fmt.Sprintf("role-%d", i))

		rb := &component.ComponentDefinition{
			Component: component.Component{Kind: "RoleBinding"},
			Configuration: map[string]interface{}{
				"roleRef":  map[string]interface{}{"name": fmt.Sprintf("role-%d", i%bindingFanout)},
				"subjects": []interface{}{map[string]interface{}{"name": fmt.Sprintf("sa-%d", i)}},
			},
		}
		rb.ID = staticUUID(fmt.Sprintf("rb-%d", i))

		sa := &component.ComponentDefinition{
			Component:     component.Component{Kind: "ServiceAccount"},
			Configuration: map[string]interface{}{"name": fmt.Sprintf("sa-%d", i)},
		}
		sa.ID = staticUUID(fmt.Sprintf("sa-%d", i))

		comps = append(comps, role, rb, sa)
	}

	roleMutatorRef := relationship.MutatorRef{[]string{"configuration", "name"}}
	rbMutatedRef := relationship.MutatedRef{[]string{"configuration", "roleRef", "name"}}
	saMutatorRef := relationship.MutatorRef{[]string{"configuration", "subjects", "0", "name"}}
	saMutatedRef := relationship.MutatedRef{[]string{"configuration", "name"}}

	fromMatchFrom := []relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &roleMutatorRef}}
	fromMatchTo := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &rbMutatedRef}}
	toMatchFrom := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatorRef: &saMutatorRef}}
	toMatchTo := []relationship.MatchSelectorItem{{Kind: "self", MutatedRef: &saMutatedRef}}

	relDef := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "binding",
		Selectors: &relationship.SelectorSet{
			relationship.SelectorSetItem{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{Kind: strPtr("Role"), Match: &relationship.MatchSelector{From: &fromMatchFrom, To: &fromMatchTo}},
					},
					To: []relationship.SelectorItem{
						{Kind: strPtr("ServiceAccount"), Match: &relationship.MatchSelector{From: &toMatchFrom, To: &toMatchTo}},
					},
				},
			},
		},
	}
	relDef.ID = staticUUID("rel-binding-bench")

	return &pattern.PatternFile{Components: comps}, relDef
}

// BenchmarkIdentifyBinding measures identifyBindingRelationships as the
// component count grows, so the perf delta from hoisting (bindingDecl, toDecl)
// validity out of the fromDecl loop is visible.
func BenchmarkIdentifyBinding(b *testing.B) {
	for _, n := range []int{10, 50, 200} {
		design, relDef := buildBindingBenchFixture(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = identifyBindingRelationships(relDef, design)
			}
		})
	}
}
