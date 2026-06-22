package policies

import (
	"sort"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// TagPolicy handles sibling relationships that group components sharing an
// equal-valued configuration field. Despite the historical "matchLabels" name
// (Kubernetes terminology), the grouping is not restricted to metadata.labels;
// it applies to any equal-valued config field referenced by the selector.
type TagPolicy struct{}

const maxTags = 20

func (p *TagPolicy) Identifier() string {
	// Wire identifier shared with the Rego engine (matchlabels-policy.rego) and
	// used as a staticUUID seed, so it stays stable until renamed in
	// meshery/schemas. See meshery/meshery#19149 (3).
	return "sibling_match_labels_policy"
}

func (p *TagPolicy) IsImplicatedBy(rel *relationship.RelationshipDefinition) bool {
	return strings.EqualFold(rel.RelationshipType, "sibling")
}

func (p *TagPolicy) IsInvalid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return p.IsImplicatedBy(rel)
}

func (p *TagPolicy) AlreadyExists(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return relationshipAlreadyExists(design, rel)
}

func (p *TagPolicy) IdentifyRelationship(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	return identifyTagRelationships(relDef, design)
}

func (p *TagPolicy) SideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	return nil
}

// tagGroup represents a group of components sharing a config field+value.
type tagGroup struct {
	Field      string
	Value      interface{}
	Components []*component.ComponentDefinition
}

// identifyTags finds all groups of components sharing an equal-valued config field.
//
// Single-pass bucketing: each component contributes its (field, value) pairs to
// shared buckets once, instead of re-scanning all (comp, other) pairs. Reduces
// the dominating cost from O(N²·F) to O(N·F) where N = components and F =
// average label fields per component. Output groups are byte-identical to the
// previous double-loop implementation on symmetric fixtures, which is the only
// shape this policy is exercised against.
func identifyTags(design *pattern.PatternFile, relDef *relationship.RelationshipDefinition) []tagGroup {
	type bucketKey struct {
		field, valueKey string
	}
	type bucket struct {
		field      string
		value      interface{}
		components []*component.ComponentDefinition
		seen       map[uuid.UUID]bool
		hasFrom    bool
		hasTo      bool
	}
	buckets := make(map[bucketKey]*bucket)

	// ingest extracts (field, value) pairs from a component's config at the
	// selector's Match.Refs paths and appends the component to each matching
	// bucket. isFrom records the component's role so empty-role buckets can be
	// filtered out at emission time (matches the original pair requirement of
	// at least one feasible-from and one feasible-to component).
	ingest := func(comp *component.ComponentDefinition, sel *relationship.SelectorItem, isFrom bool) {
		if sel == nil || sel.Match == nil || sel.Match.Refs == nil {
			return
		}
		// Dedup buckets touched by this component across multiple paths so a
		// component appears at most once per bucket.
		touched := make(map[bucketKey]bool)
		for _, path := range *sel.Match.Refs {
			cfg := configurationForComponentAtPath(path, comp, design)
			cfgMap, ok := cfg.(map[string]interface{})
			if !ok {
				continue
			}
			for field, value := range cfgMap {
				key := bucketKey{field: field, valueKey: canonicalSeed(value)}
				if touched[key] {
					continue
				}
				touched[key] = true

				b := buckets[key]
				if b == nil {
					b = &bucket{field: field, value: value, seen: make(map[uuid.UUID]bool)}
					buckets[key] = b
				}
				if !b.seen[comp.ID] {
					b.seen[comp.ID] = true
					b.components = append(b.components, comp)
				}
				if isFrom {
					b.hasFrom = true
				} else {
					b.hasTo = true
				}
			}
		}
	}

	for _, comp := range design.Components {
		ingest(comp, isRelationshipFeasibleFrom(comp, relDef), true)
		ingest(comp, isRelationshipFeasibleTo(comp, relDef), false)
	}

	// Iterate buckets in sorted key order so the output (and the truncation
	// when len > maxTags) is stable across runs. Key ordering matches
	// the previous `fmt.Sprintf("%s=%v", field, value)` lexicographic sort.
	sortedKeys := make([]bucketKey, 0, len(buckets))
	for k := range buckets {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		if sortedKeys[i].field != sortedKeys[j].field {
			return sortedKeys[i].field < sortedKeys[j].field
		}
		return sortedKeys[i].valueKey < sortedKeys[j].valueKey
	})

	groups := make([]tagGroup, 0, len(sortedKeys))
	for _, k := range sortedKeys {
		b := buckets[k]
		if !b.hasFrom || !b.hasTo {
			continue
		}
		if len(b.components) < 2 {
			continue
		}
		groups = append(groups, tagGroup{
			Field:      b.field,
			Value:      b.value,
			Components: b.components,
		})
	}
	if len(groups) > maxTags {
		groups = groups[:maxTags]
	}
	return groups
}

// identifyTagRelationships creates relationship definitions from tag groups.
func identifyTagRelationships(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	groups := identifyTags(design, relDef)

	var identified []*relationship.RelationshipDefinition
	seen := make(map[string]bool)

	for _, group := range groups {
		var fromSelectors []relationship.SelectorItem
		for _, comp := range group.Components {
			kind := comp.Component.Kind
			mutatorRef := relationship.MutatorRef{[]string{group.Field}}
			fromSel := relationship.SelectorItem{
				ID:   &comp.ID,
				Kind: &kind,
				RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
					MutatorRef: &mutatorRef,
				},
			}
			if comp.ModelReference.Name != "" {
				fromSel.Model = &comp.ModelReference
			}
			fromSelectors = append(fromSelectors, fromSel)
		}

		ss := relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: fromSelectors,
				To:   fromSelectors,
			},
		}

		compIDs := make([]string, 0, len(group.Components))
		for _, c := range group.Components {
			compIDs = append(compIDs, c.ID.String())
		}
		sort.Strings(compIDs)
		seed := map[string]any{
			"policy": "sibling_match_labels_policy",
			"field":  group.Field,
			"value":  group.Value,
			"comps":  compIDs,
		}
		uuid := staticUUID(seed)
		id := uuid.String()
		if seen[id] {
			continue
		}
		seen[id] = true

		decl := deepCopyRelDef(relDef)
		selectors := relationship.SelectorSet{ss}
		decl.Selectors = &selectors
		decl.ID = uuid
		setRelStatus(decl, StatusIdentified)

		identified = append(identified, decl)
	}
	return identified
}
