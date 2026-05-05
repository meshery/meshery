package policies

import (
	"fmt"
	"strings"

	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// MatchLabelsPolicy handles sibling match-label relationships.
type MatchLabelsPolicy struct{}

const maxMatchLabels = 20

func (p *MatchLabelsPolicy) Identifier() string {
	return "sibling_match_labels_policy"
}

func (p *MatchLabelsPolicy) IsImplicatedBy(rel *relationship.RelationshipDefinition) bool {
	return strings.EqualFold(rel.RelationshipType, "sibling")
}

func (p *MatchLabelsPolicy) IsInvalid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return p.IsImplicatedBy(rel)
}

func (p *MatchLabelsPolicy) AlreadyExists(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return relationshipAlreadyExists(design, rel)
}

func (p *MatchLabelsPolicy) IdentifyRelationship(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	return identifyMatchlabelRelationships(relDef, design)
}

func (p *MatchLabelsPolicy) SideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	return nil
}

// matchLabelGroup represents a group of components sharing a label field+value.
type matchLabelGroup struct {
	Field      string
	Value      interface{}
	Components []*component.ComponentDefinition
}

// identifyMatchlabels finds all matching label groups in the design.
func identifyMatchlabels(design *pattern.PatternFile, relDef *relationship.RelationshipDefinition) []matchLabelGroup {
	type fieldPair struct {
		field      string
		value      interface{}
		components []*component.ComponentDefinition
	}

	var pairs []fieldPair
	for _, comp := range design.Components {
		for _, other := range design.Components {
			if comp.ID == other.ID {
				continue
			}
			from := isRelationshipFeasibleFrom(comp, relDef)
			if from == nil {
				continue
			}
			if isRelationshipFeasibleTo(other, relDef) == nil {
				continue
			}

			if from.Match == nil || from.Match.Refs == nil || len(*from.Match.Refs) == 0 {
				continue
			}
			path := (*from.Match.Refs)[0]

			compConfig := configurationForComponentAtPath(path, comp, design)
			otherConfig := configurationForComponentAtPath(path, other, design)

			compMap, compOk := compConfig.(map[string]interface{})
			otherMap, otherOk := otherConfig.(map[string]interface{})
			if !compOk || !otherOk {
				continue
			}

			for field, value := range compMap {
				if otherValue, exists := otherMap[field]; exists && deepEqual(value, otherValue) {
					pairs = append(pairs, fieldPair{
						field:      field,
						value:      value,
						components: []*component.ComponentDefinition{comp, other},
					})
				}
			}
		}
	}

	// Merge pairs into groups.
	groupMap := make(map[string]*matchLabelGroup)
	for _, pair := range pairs {
		key := fmt.Sprintf("%s=%v", pair.field, pair.value)
		if g, exists := groupMap[key]; exists {
			for _, c := range pair.components {
				found := false
				for _, existing := range g.Components {
					if existing.ID == c.ID {
						found = true
						break
					}
				}
				if !found {
					g.Components = append(g.Components, c)
				}
			}
		} else {
			groupMap[key] = &matchLabelGroup{
				Field:      pair.field,
				Value:      pair.value,
				Components: pair.components,
			}
		}
	}

	groups := make([]matchLabelGroup, 0, len(groupMap))
	for _, g := range groupMap {
		groups = append(groups, *g)
	}
	if len(groups) > maxMatchLabels {
		groups = groups[:maxMatchLabels]
	}
	return groups
}

// identifyMatchlabelRelationships creates relationship definitions from matchlabel groups.
func identifyMatchlabelRelationships(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	groups := identifyMatchlabels(design, relDef)

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

		id := staticUUID(fmt.Sprintf("%v%s", ss, "sibling_match_labels_policy")).String()
		if seen[id] {
			continue
		}
		seen[id] = true

		decl := deepCopyRelDef(relDef)
		selectors := relationship.SelectorSet{ss}
		decl.Selectors = &selectors
		decl.ID = staticUUID(fmt.Sprintf("%v%s", ss, "sibling_match_labels_policy"))
		setRelStatus(decl, StatusIdentified)

		identified = append(identified, decl)
	}
	return identified
}
