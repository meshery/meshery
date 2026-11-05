package policies

import (
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
)

// isRelationshipFeasible checks if a selector matches a component (kind + model name).
func isRelationshipFeasible(sel relationship.SelectorItem, comp *component.ComponentDefinition) bool {
	selKind := selectorItemKind(sel)
	if !matchName(comp.Component.Kind, selKind) {
		return false
	}
	if sel.Model != nil {
		sModelName := sel.Model.Name
		cModelName := comp.ModelReference.Name
		if sModelName != "" && !matchName(cModelName, sModelName) {
			return false
		}
	}
	return true
}

// isRelationshipFeasibleTo checks if a component can be the "to" side of a relationship.
func isRelationshipFeasibleTo(comp *component.ComponentDefinition, rel *relationship.RelationshipDefinition) *relationship.SelectorItem {
	if rel.Selectors == nil {
		return nil
	}
	for _, ss := range *rel.Selectors {
		for i := range ss.Allow.To {
			if isRelationshipFeasible(ss.Allow.To[i], comp) {
				return &ss.Allow.To[i]
			}
		}
	}
	return nil
}

// isRelationshipFeasibleFrom checks if a component can be the "from" side of a relationship.
func isRelationshipFeasibleFrom(comp *component.ComponentDefinition, rel *relationship.RelationshipDefinition) *relationship.SelectorItem {
	if rel.Selectors == nil {
		return nil
	}
	for _, ss := range *rel.Selectors {
		for i := range ss.Allow.From {
			if isRelationshipFeasible(ss.Allow.From[i], comp) {
				return &ss.Allow.From[i]
			}
		}
	}
	return nil
}

// feasibleRelationshipSelectorBetween finds a feasible selector pair between two components.
func feasibleRelationshipSelectorBetween(fromComp, toComp *component.ComponentDefinition, rel *relationship.RelationshipDefinition) *relationship.SelectorSetItem {
	if rel.Selectors == nil {
		return nil
	}
	for i := range *rel.Selectors {
		ss := &(*rel.Selectors)[i]

		// Check deny selectors first.
		if ss.Deny != nil {
			if isSelectorFeasibleBetween(toComp, fromComp, ss.Deny.To, ss.Deny.From) {
				continue
			}
		}

		if isSelectorFeasibleBetween(toComp, fromComp, ss.Allow.To, ss.Allow.From) {
			return ss
		}
	}
	return nil
}

func isSelectorFeasibleBetween(toComp, fromComp *component.ComponentDefinition, toSels, fromSels []relationship.SelectorItem) bool {
	for _, toSel := range toSels {
		if !isRelationshipFeasible(toSel, toComp) {
			continue
		}
		for _, fromSel := range fromSels {
			if isRelationshipFeasible(fromSel, fromComp) {
				return true
			}
		}
	}
	return false
}

// isRelationshipDenied checks if a relationship between two components is denied.
func isRelationshipDenied(fromComp, toComp *component.ComponentDefinition, deny *relationship.Selector) bool {
	if deny == nil || len(deny.From) == 0 || len(deny.To) == 0 {
		return false
	}
	for _, dfSel := range deny.From {
		if !anySelectorMatches(fromComp, dfSel) {
			continue
		}
		for _, dtSel := range deny.To {
			if anySelectorMatches(toComp, dtSel) {
				return true
			}
		}
	}
	return false
}

// anySelectorMatches checks if a selector matches a component's kind and model.
func anySelectorMatches(comp *component.ComponentDefinition, sel relationship.SelectorItem) bool {
	return selectorAndComponentKindMatches(sel, comp) &&
		selectorAndComponentModelMatches(sel, comp)
}

func selectorAndComponentKindMatches(sel relationship.SelectorItem, comp *component.ComponentDefinition) bool {
	return matchName(comp.Component.Kind, selectorItemKind(sel))
}

func selectorAndComponentModelMatches(sel relationship.SelectorItem, comp *component.ComponentDefinition) bool {
	if sel.Model == nil {
		return true
	}
	sName := sel.Model.Name
	cName := comp.ModelReference.Name
	if sName != "*" && sName != cName {
		if !matchName(cName, sName) {
			return false
		}
	}
	return true
}
