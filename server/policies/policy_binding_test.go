package policies

import (
	"fmt"
	"sort"
	"testing"

	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// marshalProbe counts how many times the component holding it is JSON
// encoded. toGenericMap is a JSON round-trip, so placing a probe in a
// component's configuration turns every conversion of that component into a
// visible increment without touching the code under test.
type marshalProbe struct {
	owner  string
	counts map[string]int
}

func (p marshalProbe) MarshalJSON() ([]byte, error) {
	p.counts[p.owner]++
	return []byte(`"probe"`), nil
}

// buildWildcardBindingFixture returns n Role/RoleBinding/ServiceAccount
// triples wired through a binding relationship whose RoleBinding -> SA side
// uses the "_" array wildcard (subjects/_/name), the shape real Kubernetes
// RBAC relationships use. Every RoleBinding has two subjects; only the last
// one (which "_" resolves to) names the ServiceAccount of the same index.
func buildWildcardBindingFixture(n int, counts map[string]int) (*pattern.PatternFile, *relationship.RelationshipDefinition) {
	probe := func(owner string) marshalProbe { return marshalProbe{owner: owner, counts: counts} }

	comps := make([]*component.ComponentDefinition, 0, 3*n)
	for i := 0; i < n; i++ {
		roleName, rbName, saName := fmt.Sprintf("role-%d", i), fmt.Sprintf("rb-%d", i), fmt.Sprintf("sa-%d", i)

		role := &component.ComponentDefinition{
			Component:     component.Component{Kind: "Role"},
			Configuration: map[string]interface{}{"name": roleName, "probe": probe(roleName)},
		}
		role.ID = staticUUID(roleName)

		rb := &component.ComponentDefinition{
			Component: component.Component{Kind: "RoleBinding"},
			Configuration: map[string]interface{}{
				"roleRef": map[string]interface{}{"name": roleName},
				"subjects": []interface{}{
					map[string]interface{}{"name": "not-the-last-subject"},
					map[string]interface{}{"name": saName},
				},
				"probe": probe(rbName),
			},
		}
		rb.ID = staticUUID(rbName)

		sa := &component.ComponentDefinition{
			Component:     component.Component{Kind: "ServiceAccount"},
			Configuration: map[string]interface{}{"name": saName, "probe": probe(saName)},
		}
		sa.ID = staticUUID(saName)

		comps = append(comps, role, rb, sa)
	}

	roleMutatorRef := relationship.MutatorRef{[]string{"configuration", "name"}}
	rbMutatedRef := relationship.MutatedRef{[]string{"configuration", "roleRef", "name"}}
	rbMutatorRef := relationship.MutatorRef{[]string{"configuration", "subjects", "_", "name"}}
	saMutatedRef := relationship.MutatedRef{[]string{"configuration", "name"}}

	fromMatchFrom := []relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &roleMutatorRef}}
	fromMatchTo := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &rbMutatedRef}}
	toMatchFrom := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatorRef: &rbMutatorRef}}
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
	relDef.ID = staticUUID("rel-binding-wildcard")

	return &pattern.PatternFile{Components: comps}, relDef
}

// identifiedTriples returns "from->binding->to" component names for each
// identified relationship, sorted, so tests can compare results exactly.
func identifiedTriples(t *testing.T, rels []*relationship.RelationshipDefinition, design *pattern.PatternFile) []string {
	t.Helper()
	nameOf := func(id string) string {
		c := componentDeclarationByID(design, id)
		if c == nil {
			t.Fatalf("identified relationship references unknown component %s", id)
		}
		name, _ := c.Configuration["name"].(string)
		if name == "" {
			// RoleBindings carry no top-level name; use the role they reference.
			if ref, ok := c.Configuration["roleRef"].(map[string]interface{}); ok {
				name = "rb(" + ref["name"].(string) + ")"
			}
		}
		return name
	}

	out := make([]string, 0, len(rels))
	for _, rel := range rels {
		ss := (*rel.Selectors)[0]
		from := ss.Allow.From[0]
		to := ss.Allow.To[0]
		bindingID := (*from.Match.To)[0].ID.String()
		out = append(out, fmt.Sprintf("%s->%s->%s", nameOf(from.ID.String()), nameOf(bindingID), nameOf(to.ID.String())))
	}
	sort.Strings(out)
	return out
}

// TestIdentifyBindingRelationshipsConvertsEachComponentOnce guards against
// identification re-encoding a component for every candidate pair it is
// part of. Each component may be converted to a generic map at most once per
// identification pass.
func TestIdentifyBindingRelationshipsConvertsEachComponentOnce(t *testing.T) {
	const n = 8
	counts := map[string]int{}
	design, relDef := buildWildcardBindingFixture(n, counts)

	rels := identifyBindingRelationships(relDef, design)
	if len(rels) != n {
		t.Fatalf("expected %d identified relationships, got %d", n, len(rels))
	}

	for owner, c := range counts {
		if c > 1 {
			t.Errorf("component %s was converted %d times during one identification pass, want at most 1", owner, c)
		}
	}
}

// TestIdentifyBindingRelationshipsWildcardResolution checks that "_" paths
// keep resolving against the right component's own array once conversions
// are shared across candidate pairs: each RoleBinding must link only to the
// ServiceAccount named by its last subject.
func TestIdentifyBindingRelationshipsWildcardResolution(t *testing.T) {
	const n = 4
	design, relDef := buildWildcardBindingFixture(n, map[string]int{})

	got := identifiedTriples(t, identifyBindingRelationships(relDef, design), design)

	want := make([]string, 0, n)
	for i := 0; i < n; i++ {
		want = append(want, fmt.Sprintf("role-%d->rb(role-%d)->sa-%d", i, i, i))
	}
	sort.Strings(want)

	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("identified relationships mismatch\n got: %v\nwant: %v", got, want)
	}
}

// TestIdentifyBindingRelationshipsWildcardPerComponentArrayLength makes sure
// a shared conversion is keyed by component: RoleBindings whose subjects
// arrays have different lengths must each resolve "_" to their own last index.
func TestIdentifyBindingRelationshipsWildcardPerComponentArrayLength(t *testing.T) {
	design, relDef := buildWildcardBindingFixture(2, map[string]int{})

	// rb-1 gets a third subject naming sa-0. "_" must now resolve to index 2
	// for rb-1 (-> sa-0) while rb-0 keeps resolving to index 1 (-> sa-0).
	for _, c := range design.Components {
		if c.ID == staticUUID("rb-1") {
			subjects := c.Configuration["subjects"].([]interface{})
			c.Configuration["subjects"] = append(subjects, map[string]interface{}{"name": "sa-0"})
		}
	}

	got := identifiedTriples(t, identifyBindingRelationships(relDef, design), design)
	want := []string{
		"role-0->rb(role-0)->sa-0",
		"role-1->rb(role-1)->sa-0",
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("identified relationships mismatch\n got: %v\nwant: %v", got, want)
	}
}

// TestIsValidBindingTypedSkipsConversionWithoutWildcard checks that a pair
// whose ref paths contain no "_" wildcard is validated without converting
// either component: the generic map is only needed to resolve wildcards.
func TestIsValidBindingTypedSkipsConversionWithoutWildcard(t *testing.T) {
	counts := map[string]int{}
	design, relDef := buildWildcardBindingFixture(1, counts)

	var role, rb *component.ComponentDefinition
	for _, c := range design.Components {
		switch c.Component.Kind {
		case "Role":
			role = c
		case "RoleBinding":
			rb = c
		}
	}

	// Role -> RoleBinding uses configuration/name and configuration/roleRef/name.
	fromSel := (*relDef.Selectors)[0].Allow.From[0]
	if !isValidBindingTyped(role, rb, fromSel, design) {
		t.Fatal("expected role-0 and its RoleBinding to be a valid binding")
	}
	if total := counts["role-0"] + counts["rb-0"]; total != 0 {
		t.Errorf("components were converted %d times for a pair with no wildcard paths, want 0", total)
	}
}

// TestIsValidBindingCachesNonConfigurationLookups checks that a pair whose
// ref paths do not start with "configuration" is validated without
// converting either component more than once per pass: those lookups go
// through the same per-pass cache as wildcard resolution.
func TestIsValidBindingCachesNonConfigurationLookups(t *testing.T) {
	counts := map[string]int{}
	design, _ := buildWildcardBindingFixture(1, counts)

	var role, rb *component.ComponentDefinition
	for _, c := range design.Components {
		switch c.Component.Kind {
		case "Role":
			role = c
		case "RoleBinding":
			rb = c
		}
	}

	// Both components share a display name, so a ["displayName"] ref pair
	// matches. DisplayName lives outside "configuration", so without the
	// cache it would be converted once per candidate check.
	role.DisplayName = "shared-name"
	rb.DisplayName = "shared-name"

	sel := relationship.SelectorItem{
		Kind: strPtr("RoleBinding"),
		Match: &relationship.MatchSelector{
			From: &[]relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &relationship.MutatorRef{[]string{"displayName"}}}},
			To:   &[]relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &relationship.MutatedRef{[]string{"displayName"}}}},
		},
	}

	maps := componentMapCache{}
	for i := 0; i < 4; i++ {
		if !isValidBinding(role, rb, sel, design, maps) {
			t.Fatal("expected role-0 and rb-0 to be a valid binding on displayName")
		}
	}
	for owner, c := range counts {
		if c > 1 {
			t.Errorf("component %s was converted %d times for non-configuration lookups, want at most 1", owner, c)
		}
	}
}

// TestIsValidBindingCachesAliasParentLookups checks that a configuration
// lookup on an aliased component converts the alias parent at most once per
// pass: the parent's map comes from the same per-pass cache, not from a
// fresh conversion on every candidate check.
func TestIsValidBindingCachesAliasParentLookups(t *testing.T) {
	counts := map[string]int{}
	design, relDef := buildWildcardBindingFixture(1, counts)

	var role, rb *component.ComponentDefinition
	for _, c := range design.Components {
		switch c.Component.Kind {
		case "Role":
			role = c
		case "RoleBinding":
			rb = c
		}
	}

	// Alias role-0 to a parent component whose configuration names the role
	// the RoleBinding references, so the ["configuration", "name"] lookup on
	// role-0 reads the parent's configuration through the alias.
	parent := &component.ComponentDefinition{
		Component:     component.Component{Kind: "Role"},
		Configuration: map[string]interface{}{"name": "role-0", "probe": marshalProbe{owner: "parent-0", counts: counts}},
	}
	parent.ID = staticUUID("parent-0")
	design.Components = append(design.Components, parent)

	aliases := map[string]core.ResolvedAlias{
		role.ID.String(): {
			AliasComponentId:     role.ID,
			ResolvedParentId:     parent.ID,
			ResolvedRefFieldPath: []string{"configuration"},
		},
	}
	design.Metadata = &pattern.PatternFile_Metadata{ResolvedAliases: &aliases}

	// Role -> RoleBinding uses configuration/name and configuration/roleRef/name.
	fromSel := (*relDef.Selectors)[0].Allow.From[0]

	maps := componentMapCache{}
	for i := 0; i < 4; i++ {
		if !isValidBinding(role, rb, fromSel, design, maps) {
			t.Fatal("expected role-0 and rb-0 to be a valid binding through the alias")
		}
	}
	for owner, c := range counts {
		if c > 1 {
			t.Errorf("component %s was converted %d times for alias-parent lookups, want at most 1", owner, c)
		}
	}
	if counts["parent-0"] != 1 {
		t.Errorf("alias parent was converted %d times, want exactly 1", counts["parent-0"])
	}
}

// BenchmarkIdentifyBindingWildcard is BenchmarkIdentifyBinding with a "_"
// wildcard in the RoleBinding -> ServiceAccount refs, so the generic-map
// conversion that wildcard resolution needs is part of the measured work.
func BenchmarkIdentifyBindingWildcard(b *testing.B) {
	for _, n := range []int{10, 50, 200} {
		design, relDef := buildWildcardBindingFixture(n, map[string]int{})
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = identifyBindingRelationships(relDef, design)
			}
		})
	}
}
