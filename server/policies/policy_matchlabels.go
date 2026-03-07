package policies

import (
	"fmt"
	"strings"
)

// MatchLabelsPolicy handles sibling match-label relationships.
type MatchLabelsPolicy struct{}

const maxMatchLabels = 20

func (p *MatchLabelsPolicy) Identifier() string {
	return "sibling_match_labels_policy"
}

func (p *MatchLabelsPolicy) IsImplicatedBy(rel map[string]interface{}) bool {
	return strings.EqualFold(getMapString(rel, "type"), "sibling")
}

func (p *MatchLabelsPolicy) IsInvalid(rel, designFile map[string]interface{}) bool {
	// Matchlabel rels are stateless - invalidate all then re-identify
	return p.IsImplicatedBy(rel)
}

func (p *MatchLabelsPolicy) AlreadyExists(rel, designFile map[string]interface{}) bool {
	return false // matchlabels always re-identify
}

func (p *MatchLabelsPolicy) IdentifyRelationship(relDef, designFile map[string]interface{}) []map[string]interface{} {
	return identifyMatchlabelRelationships(relDef, designFile)
}

func (p *MatchLabelsPolicy) SideEffects(rel, designFile map[string]interface{}) []PolicyAction {
	return nil // no side effects
}

// matchLabelGroup represents a group of components sharing a label field+value.
type matchLabelGroup struct {
	Field      string
	Value      interface{}
	Components []map[string]interface{}
}

// identifyMatchlabels finds all matching label groups in the design.
func identifyMatchlabels(designFile, relationship map[string]interface{}) []matchLabelGroup {
	comps := getMapSlice(designFile, "components")

	type fieldPair struct {
		field      string
		value      interface{}
		components []map[string]interface{}
	}

	// Collect field pairs
	var pairs []fieldPair
	for _, ci := range comps {
		comp, ok := ci.(map[string]interface{})
		if !ok {
			continue
		}
		for _, oi := range comps {
			other, ok := oi.(map[string]interface{})
			if !ok {
				continue
			}
			if getMapString(comp, "id") == getMapString(other, "id") {
				continue
			}

			from := isRelationshipFeasibleFrom(comp, relationship)
			if from == nil {
				continue
			}
			if isRelationshipFeasibleTo(other, relationship) == nil {
				continue
			}

			// Get match refs
			match := getMapMap(from, "match")
			if match == nil {
				continue
			}
			refs := getMapSlice(match, "refs")
			if len(refs) == 0 {
				continue
			}

			path := interfaceToStringSlice(refs[0])

			compConfig := configurationForComponentAtPath(path, comp, designFile)
			otherConfig := configurationForComponentAtPath(path, other, designFile)

			compMap, compOk := compConfig.(map[string]interface{})
			otherMap, otherOk := otherConfig.(map[string]interface{})
			if !compOk || !otherOk {
				continue
			}

			for field, value := range compMap {
				if otherValue, exists := otherMap[field]; exists {
					if deepEqual(value, otherValue) {
						pairs = append(pairs, fieldPair{
							field:      field,
							value:      value,
							components: []map[string]interface{}{comp, other},
						})
					}
				}
			}
		}
	}

	// Merge pairs into groups
	groupMap := make(map[string]*matchLabelGroup)
	for _, pair := range pairs {
		key := fmt.Sprintf("%s=%v", pair.field, pair.value)
		if g, exists := groupMap[key]; exists {
			for _, comp := range pair.components {
				found := false
				compID := getMapString(comp, "id")
				for _, existing := range g.Components {
					if getMapString(existing, "id") == compID {
						found = true
						break
					}
				}
				if !found {
					g.Components = append(g.Components, comp)
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

	// Truncate to max
	if len(groups) > maxMatchLabels {
		groups = groups[:maxMatchLabels]
	}

	return groups
}

// identifyMatchlabelRelationships creates relationship declarations from matchlabel groups.
func identifyMatchlabelRelationships(relDef, designFile map[string]interface{}) []map[string]interface{} {
	groups := identifyMatchlabels(designFile, relDef)

	var identified []map[string]interface{}
	seen := make(map[string]bool)

	for _, group := range groups {
		var fromSelectors []interface{}
		for _, comp := range group.Components {
			fromSel := map[string]interface{}{
				"id":    getMapString(comp, "id"),
				"kind":  getMapString(getMapMap(comp, "component"), "kind"),
				"model": comp["model"],
				"patch": map[string]interface{}{
					"patchStrategy": "replace",
					"mutatorRef":    []interface{}{[]interface{}{group.Field}},
				},
			}
			fromSelectors = append(fromSelectors, fromSel)
		}

		selectorDecl := map[string]interface{}{
			"allow": map[string]interface{}{
				"from": fromSelectors,
				"to":   fromSelectors,
			},
			"deny": map[string]interface{}{
				"from": []interface{}{},
				"to":   []interface{}{},
			},
		}

		id := staticUUID(fmt.Sprintf("%v%s", selectorDecl, "sibling_match_labels_policy")).String()
		if seen[id] {
			continue
		}
		seen[id] = true

		rel := deepCopyMap(relDef)
		rel["selectors"] = []interface{}{selectorDecl}
		rel["id"] = id
		rel["status"] = "identified"

		identified = append(identified, rel)
	}

	return identified
}
