package policies

import (
	"fmt"
	"testing"

	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// buildMutatorBenchFixture returns n Deployments and n ConfigMaps, so only half
// the components match either side of the selector pair. That is the shape the
// pre-fix code paid most for: the kind of every component was re-tested inside
// the inner loop, making the kind check itself O(C²) even though the eventual
// pair set is much smaller.
//
// Each Deployment's mutator path holds a ConfigMap name, and one ConfigMap
// carries the matching value, so matchingMutators does real work rather than
// short-circuiting on the first ref.
//
// Pre-fix work: O(C²) kind tests, each one a regexp compile on the mismatch path.
// Post-fix work: O(C) kind tests, and the pair enumeration is unchanged.
func buildMutatorBenchFixture(n int) (*pattern.PatternFile, *relationship.RelationshipDefinition, relationship.SelectorItem, relationship.SelectorItem) {
	comps := make([]*component.ComponentDefinition, 0, 2*n)
	for i := 0; i < n; i++ {
		d := &component.ComponentDefinition{
			Component: component.Component{Kind: "Deployment"},
			Configuration: map[string]interface{}{
				"spec": map[string]interface{}{
					"configMapName": fmt.Sprintf("cm-%d", i%(n/4+1)),
				},
			},
		}
		d.ID = staticUUID(fmt.Sprintf("deploy-%d", i))
		comps = append(comps, d)

		c := &component.ComponentDefinition{
			Component: component.Component{Kind: "ConfigMap"},
			Configuration: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name": fmt.Sprintf("cm-%d", i%(n/4+1)),
				},
			},
		}
		c.ID = staticUUID(fmt.Sprintf("cm-%d", i))
		comps = append(comps, c)
	}

	fromKind := "Deployment"
	toKind := "ConfigMap"
	mutatorRefs := [][]string{{"configuration", "spec", "configMapName"}}
	mutatedRefs := [][]string{{"configuration", "metadata", "name"}}

	fromSel := relationship.SelectorItem{
		Kind: &fromKind,
		RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
			MutatorRef: &mutatorRefs,
		},
	}
	toSel := relationship.SelectorItem{
		Kind: &toKind,
		RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
			MutatedRef: &mutatedRefs,
		},
	}

	relDef := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "parent",
		Selectors: &relationship.SelectorSet{
			relationship.SelectorSetItem{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{fromSel},
					To:   []relationship.SelectorItem{toSel},
				},
			},
		},
	}
	relDef.ID = staticUUID("rel-mutator-bench")

	return &pattern.PatternFile{Components: comps}, relDef, fromSel, toSel
}

// BenchmarkIdentifyMutatorMutated measures the third identify path as the
// component count grows. The other two paths already have this coverage:
// BenchmarkIdentifyMatchlabels for the O(N²·F) → O(N·F) rewrite and
// BenchmarkIdentifyBinding for the O(n³) → O(n²) one.
func BenchmarkIdentifyMutatorMutated(b *testing.B) {
	for _, n := range []int{10, 50, 200} {
		design, relDef, _, _ := buildMutatorBenchFixture(n)
		b.Run(fmt.Sprintf("N=%d", 2*n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields(relDef, design)
			}
		})
	}
}

// BenchmarkMatchComponentPairs isolates the pair loop from the selector
// plumbing, so the kind-bucketing delta is visible on its own.
func BenchmarkMatchComponentPairs(b *testing.B) {
	for _, n := range []int{10, 50, 200} {
		design, relDef, fromSel, toSel := buildMutatorBenchFixture(n)
		b.Run(fmt.Sprintf("N=%d", 2*n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = matchComponentPairs(fromSel, toSel, design.Components, relDef, design)
			}
		})
	}
}

// BenchmarkMatchName covers the mismatch path, which is the one the pair loop
// hits most and the one that used to compile a regex every call.
func BenchmarkMatchName(b *testing.B) {
	b.Run("literal_mismatch", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = matchName("ConfigMap", "Deployment")
		}
	})
	b.Run("regex_mismatch", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = matchName("ConfigMap", "Deploy.*")
		}
	})
}
