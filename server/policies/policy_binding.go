package policies

import (
	"strings"

	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// EdgeBindingPolicy handles edge binding relationships (3-party: from, binding, to).
type EdgeBindingPolicy struct{}

func (p *EdgeBindingPolicy) Identifier() string {
	return "edge_binding"
}

func (p *EdgeBindingPolicy) IsImplicatedBy(rel *relationship.RelationshipDefinition) bool {
	return strings.EqualFold(string(rel.Kind), "edge") &&
		strings.EqualFold(rel.RelationshipType, "binding")
}

func (p *EdgeBindingPolicy) IsInvalid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return fromOrToComponentsDontExist(rel, design)
}

func (p *EdgeBindingPolicy) AlreadyExists(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return relationshipAlreadyExists(design, rel)
}

func (p *EdgeBindingPolicy) IdentifyRelationship(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	return identifyBindingRelationships(relDef, design)
}

func (p *EdgeBindingPolicy) SideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	if getRelStatus(rel) == StatusDeleted {
		return nil
	}
	actions := patchMutatorsAction(rel, design)
	if len(actions) > 0 {
		return actions
	}
	return patchBindingMatchFields(rel, design)
}

// patchBindingMatchFields patches the binding component's configuration
// using values from the match fields of the binding relationship selectors.
func patchBindingMatchFields(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	if rel.Selectors == nil {
		return nil
	}
	var actions []PolicyAction

	for _, ss := range *rel.Selectors {
		for _, selArr := range [][]relationship.SelectorItem{ss.Allow.From, ss.Allow.To} {
			if len(selArr) == 0 {
				continue
			}
			sel := selArr[0]
			if sel.Match == nil {
				continue
			}
			compID := selectorItemID(sel)
			comp := componentDeclarationByID(design, compID)
			if comp == nil {
				continue
			}

			actions = append(actions, patchBindingMatchFieldsForSelector(sel.Match, comp, design)...)
		}
	}
	return actions
}

// patchBindingMatchFieldsForSelector processes match fields for a single selector.
func patchBindingMatchFieldsForSelector(match *relationship.MatchSelector, _ *component.ComponentDefinition, design *pattern.PatternFile) []PolicyAction {
	var actions []PolicyAction

	type matchSide struct {
		self  *[]relationship.MatchSelectorItem
		other *[]relationship.MatchSelectorItem
	}
	sides := []matchSide{
		{match.From, match.To},
		{match.To, match.From},
	}
	for _, ms := range sides {
		if ms.self == nil || len(*ms.self) == 0 {
			continue
		}
		mSel := (*ms.self)[0]
		if mSel.MutatorRef == nil {
			continue
		}
		matchComp := findMatchSelectorComponent(mSel, design)
		if matchComp == nil {
			continue
		}

		if ms.other == nil || len(*ms.other) == 0 {
			continue
		}
		otherSel := (*ms.other)[0]
		if otherSel.MutatedRef == nil {
			continue
		}
		otherComp := findMatchSelectorComponent(otherSel, design)
		if otherComp == nil {
			continue
		}

		count := min(len(*mSel.MutatorRef), len(*otherSel.MutatedRef))

		for i := range count {
			mutatorPath := (*mSel.MutatorRef)[i]
			mutatedPath := (*otherSel.MutatedRef)[i]

			mutatorValue := configurationForComponentAtPath(mutatorPath, matchComp, design)
			oldValue := configurationForComponentAtPath(mutatedPath, otherComp, design)

			if deepEqual(mutatorValue, oldValue) || mutatorValue == nil {
				continue
			}
			actions = append(actions, newComponentUpdateAction(getComponentUpdateOp(mutatedPath), otherComp.ID.String(), mutatedPath, mutatorValue))
		}
	}
	return actions
}

// findMatchSelectorComponent finds the component referenced by a MatchSelectorItem.
func findMatchSelectorComponent(mSel relationship.MatchSelectorItem, design *pattern.PatternFile) *component.ComponentDefinition {
	if mSel.ID == nil {
		return nil
	}
	return componentDeclarationByID(design, mSel.ID.String())
}

// identifyBindingRelationships identifies 3-party binding relationships.
func identifyBindingRelationships(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	if relDef.Selectors == nil {
		return nil
	}
	var identified []*relationship.RelationshipDefinition

	for _, ss := range *relDef.Selectors {
		fromSelByKind := make(map[string]relationship.SelectorItem)
		for _, fromSel := range ss.Allow.From {
			fromSelByKind[selectorItemKind(fromSel)] = fromSel
		}

		bindingKinds := collectBindingKindsTyped(ss.Allow.From)
		toSelByPair := indexToSelectorsByKindPairTyped(ss.Allow.To)

		fromComps := filterComponentsByKindSetTyped(design.Components, fromSelByKind)
		toComps := filterComponentsByToSelectorsTyped(design.Components, ss.Allow.To)

		for bindingKind := range bindingKinds {
			bindingDecls := filterComponentsByKindTyped(design.Components, bindingKind)
			if len(bindingDecls) == 0 {
				continue
			}

			for _, fromDecl := range fromComps {
				fromKind := fromDecl.Component.Kind
				fromSel, exists := fromSelByKind[fromKind]
				if !exists {
					continue
				}

				for _, bindingDecl := range bindingDecls {
					if fromDecl.ID == bindingDecl.ID {
						continue
					}
					if !isValidBindingTyped(fromDecl, bindingDecl, fromSel, design) {
						continue
					}

					for _, toDecl := range toComps {
						if toDecl.ID == bindingDecl.ID {
							continue
						}
						toKind := toDecl.Component.Kind
						toSel, exists := toSelByPair[toKind+"#"+bindingKind]
						if !exists {
							continue
						}

						if ss.Deny != nil && isRelationshipDenied(fromDecl, toDecl, ss.Deny) {
							continue
						}
						if !isValidBindingTyped(bindingDecl, toDecl, toSel, design) {
							continue
						}

						newFromSel := fromSel
						newFromSel.ID = &fromDecl.ID
						patchBindingSelectorIDsTyped(&newFromSel, fromDecl.ID.String(), bindingDecl.ID.String())

						newToSel := toSel
						newToSel.ID = &toDecl.ID
						patchBindingSelectorIDsTyped(&newToSel, bindingDecl.ID.String(), toDecl.ID.String())

						selectorSetItem := relationship.SelectorSetItem{
							Allow: relationship.Selector{
								From: []relationship.SelectorItem{newFromSel},
								To:   []relationship.SelectorItem{newToSel},
							},
						}

						seed := map[string]any{
							"from":    fromDecl.ID.String(),
							"binding": bindingDecl.ID.String(),
							"to":      toDecl.ID.String(),
							"relId":   relDef.ID.String(),
						}

						decl := deepCopyRelDef(relDef)
						decl.ID = staticUUID(seed)
						selectors := relationship.SelectorSet{selectorSetItem}
						decl.Selectors = &selectors
						setRelStatus(decl, StatusApproved)

						identified = append(identified, decl)
					}
				}
			}
		}
	}
	return identified
}

// collectBindingKindsTyped extracts binding component kinds from from selectors' match fields.
func collectBindingKindsTyped(fromSels []relationship.SelectorItem) map[string]bool {
	kinds := make(map[string]bool)
	for _, fromSel := range fromSels {
		if fromSel.Match == nil {
			continue
		}
		for _, matchSels := range []*[]relationship.MatchSelectorItem{fromSel.Match.From, fromSel.Match.To} {
			if matchSels == nil {
				continue
			}
			for _, mSel := range *matchSels {
				if mSel.Kind != "" && !strings.EqualFold(mSel.Kind, "self") {
					kinds[mSel.Kind] = true
				}
			}
		}
	}
	return kinds
}

// indexToSelectorsByKindPairTyped indexes to selectors by "toKind#matchKind".
func indexToSelectorsByKindPairTyped(toSels []relationship.SelectorItem) map[string]relationship.SelectorItem {
	result := make(map[string]relationship.SelectorItem)
	for _, toSel := range toSels {
		toKind := selectorItemKind(toSel)
		if toSel.Match == nil {
			continue
		}
		for _, matchSels := range []*[]relationship.MatchSelectorItem{toSel.Match.From, toSel.Match.To} {
			if matchSels == nil {
				continue
			}
			for _, mSel := range *matchSels {
				if mSel.Kind != "" && !strings.EqualFold(mSel.Kind, "self") {
					result[toKind+"#"+mSel.Kind] = toSel
				}
			}
		}
	}
	return result
}

// patchBindingSelectorIDsTyped sets the IDs in a binding selector's match field.
func patchBindingSelectorIDsTyped(sel *relationship.SelectorItem, matchFromID, matchToID string) {
	if sel.Match == nil {
		return
	}
	fromUUID, _ := uuidFromString(matchFromID)
	toUUID, _ := uuidFromString(matchToID)
	if sel.Match.From != nil && len(*sel.Match.From) > 0 {
		(*sel.Match.From)[0].ID = &fromUUID
	}
	if sel.Match.To != nil && len(*sel.Match.To) > 0 {
		(*sel.Match.To)[0].ID = &toUUID
	}
}

// isValidBindingTyped checks if two components match via a binding selector's match field.
func isValidBindingTyped(comp1, comp2 *component.ComponentDefinition, sel relationship.SelectorItem, design *pattern.PatternFile) bool {
	if sel.Match == nil {
		return false
	}

	mutatorComp, mutatorPaths := extractMutatorFromMatch(sel.Match, comp1, comp2)
	mutatedComp, mutatedPaths := extractMutatedFromMatch(sel.Match, comp1, comp2)
	if mutatorComp == nil || mutatedComp == nil || len(mutatorPaths) == 0 || len(mutatedPaths) == 0 {
		return false
	}

	mutatorMap, _ := toGenericMap(mutatorComp)
	mutatedMap, _ := toGenericMap(mutatedComp)

	count := min(len(mutatorPaths), len(mutatedPaths))

	for i := range count {
		resolvedFrom := resolvePath(mutatorPaths[i], mutatorMap)
		resolvedTo := resolvePath(mutatedPaths[i], mutatedMap)
		val1 := configurationForComponentAtPath(resolvedFrom, mutatorComp, design)
		val2 := configurationForComponentAtPath(resolvedTo, mutatedComp, design)
		if val1 != nil && val2 != nil && !deepEqual(val1, val2) {
			return false
		}
	}
	return true
}

// extractMutatorFromMatch finds the mutator component and ref paths from a match field.
func extractMutatorFromMatch(match *relationship.MatchSelector, comp1, comp2 *component.ComponentDefinition) (*component.ComponentDefinition, [][]string) {
	if match.From != nil && len(*match.From) > 0 {
		if (*match.From)[0].MutatorRef != nil {
			return comp1, *(*match.From)[0].MutatorRef
		}
	}
	if match.To != nil && len(*match.To) > 0 {
		if (*match.To)[0].MutatorRef != nil {
			return comp2, *(*match.To)[0].MutatorRef
		}
	}
	return nil, nil
}

// extractMutatedFromMatch finds the mutated component and ref paths from a match field.
func extractMutatedFromMatch(match *relationship.MatchSelector, comp1, comp2 *component.ComponentDefinition) (*component.ComponentDefinition, [][]string) {
	if match.From != nil && len(*match.From) > 0 {
		if (*match.From)[0].MutatedRef != nil {
			return comp1, *(*match.From)[0].MutatedRef
		}
	}
	if match.To != nil && len(*match.To) > 0 {
		if (*match.To)[0].MutatedRef != nil {
			return comp2, *(*match.To)[0].MutatedRef
		}
	}
	return nil, nil
}

// filterComponentsByKindTyped returns components matching the given kind.
func filterComponentsByKindTyped(comps []*component.ComponentDefinition, kind string) []*component.ComponentDefinition {
	var result []*component.ComponentDefinition
	for _, comp := range comps {
		if matchName(comp.Component.Kind, kind) {
			result = append(result, comp)
		}
	}
	return result
}

// filterComponentsByKindSetTyped returns components whose kind is in the given set.
func filterComponentsByKindSetTyped(comps []*component.ComponentDefinition, kindSet map[string]relationship.SelectorItem) []*component.ComponentDefinition {
	var result []*component.ComponentDefinition
	for _, comp := range comps {
		if _, exists := kindSet[comp.Component.Kind]; exists {
			result = append(result, comp)
		}
	}
	return result
}

// filterComponentsByToSelectorsTyped returns components matching any to selector kind.
func filterComponentsByToSelectorsTyped(comps []*component.ComponentDefinition, toSels []relationship.SelectorItem) []*component.ComponentDefinition {
	toKinds := make(map[string]bool)
	for _, toSel := range toSels {
		toKinds[selectorItemKind(toSel)] = true
	}
	var result []*component.ComponentDefinition
	for _, comp := range comps {
		if toKinds[comp.Component.Kind] || toKinds["*"] {
			result = append(result, comp)
		}
	}
	return result
}
