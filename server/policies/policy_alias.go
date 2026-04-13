package policies

import (
	"fmt"
	"strings"

	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// AliasPolicy handles alias relationships (hierarchical parent alias).
type AliasPolicy struct{}

func (p *AliasPolicy) Identifier() string {
	return "alias_relationships_policy"
}

func (p *AliasPolicy) IsImplicatedBy(rel *relationship.RelationshipDefinition) bool {
	return isAliasRelationship(rel)
}

func (p *AliasPolicy) IsInvalid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return !isAliasRelationshipValid(rel, design)
}

func (p *AliasPolicy) AlreadyExists(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return aliasRelationshipAlreadyExists(design, rel)
}

func (p *AliasPolicy) IdentifyRelationship(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	var allIdentified []*relationship.RelationshipDefinition
	for _, comp := range design.Components {
		if isRelationshipFeasibleTo(comp, relDef) == nil {
			continue
		}
		identified := identifyAliasRelationships(comp, relDef, design)
		allIdentified = append(allIdentified, identified...)
	}
	return allIdentified
}

func (p *AliasPolicy) SideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	status := getRelStatus(rel)
	if status == StatusIdentified || status == StatusPending {
		return aliasAddComponentSideEffects(rel, design)
	}
	if status == StatusDeleted {
		return aliasDeleteComponentSideEffects(rel, design)
	}
	return nil
}

// isAliasRelationship checks if a relationship is an alias type.
func isAliasRelationship(rel *relationship.RelationshipDefinition) bool {
	return strings.EqualFold(string(rel.Kind), "hierarchical") &&
		strings.EqualFold(rel.RelationshipType, "parent") &&
		strings.EqualFold(rel.SubType, "alias")
}

// isAliasRelationshipValid checks if an alias relationship is still valid.
func isAliasRelationshipValid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	status := getRelStatus(rel)
	if status == StatusPending {
		return true
	}
	if status != StatusApproved {
		return false
	}

	fromComp := componentDeclarationByID(design, fromComponentID(rel))
	if fromComp == nil {
		return false
	}
	toComp := componentDeclarationByID(design, toComponentID(rel))
	if toComp == nil {
		return false
	}

	ref := aliasRefFromRelationship(rel)
	if ref == nil {
		return false
	}
	config := getComponentConfiguration(toComp, design)
	return objectGetNested(config, popFirst(ref), nil) != nil
}

// aliasRefFromRelationship extracts the mutatorRef path from an alias relationship.
func aliasRefFromRelationship(rel *relationship.RelationshipDefinition) []string {
	if rel.Selectors == nil {
		return nil
	}
	for _, ss := range *rel.Selectors {
		for _, fromSel := range ss.Allow.From {
			if fromSel.RelationshipDefinitionSelectorsPatch != nil && fromSel.RelationshipDefinitionSelectorsPatch.MutatorRef != nil && len(*fromSel.RelationshipDefinitionSelectorsPatch.MutatorRef) > 0 {
				return (*fromSel.RelationshipDefinitionSelectorsPatch.MutatorRef)[0]
			}
		}
	}
	return nil
}

// identifyAliasRelationships identifies alias relationships for a component.
func identifyAliasRelationships(comp *component.ComponentDefinition, relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	if relDef.Selectors == nil {
		return nil
	}
	var identified []*relationship.RelationshipDefinition

	for _, ss := range *relDef.Selectors {
		for _, fromSel := range ss.Allow.From {
			for _, toSel := range ss.Allow.To {
				if fromSel.RelationshipDefinitionSelectorsPatch == nil || fromSel.RelationshipDefinitionSelectorsPatch.MutatorRef == nil || len(*fromSel.RelationshipDefinitionSelectorsPatch.MutatorRef) == 0 {
					continue
				}
				ref := (*fromSel.RelationshipDefinitionSelectorsPatch.MutatorRef)[0]
				paths := getArrayAwareConfigPaths(ref, comp, design)

				for _, path := range paths {
					mutatorRef := relationship.MutatorRef{path}
					mutatedRef := relationship.MutatedRef{path}
					selectorPatch := relationship.RelationshipDefinitionSelectorsPatch{
						MutatorRef: &mutatorRef,
						MutatedRef: &mutatedRef,
					}

					compID := staticUUID(map[string]interface{}{"c": comp.ID.String(), "s": path})

					newFromSel := fromSel
					newFromSel.ID = &compID
					newFromSel.RelationshipDefinitionSelectorsPatch = &selectorPatch

					newToSel := toSel
					newToSel.ID = &comp.ID
					newToSel.RelationshipDefinitionSelectorsPatch = &selectorPatch

					selectorSetItem := relationship.SelectorSetItem{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{newFromSel},
							To:   []relationship.SelectorItem{newToSel},
						},
					}

					decl := deepCopyRelDef(relDef)
					selectors := relationship.SelectorSet{selectorSetItem}
					decl.Selectors = &selectors
					decl.ID = staticUUID(selectorSetItem)
					setRelStatus(decl, StatusIdentified)

					identified = append(identified, decl)
				}
			}
		}
	}
	return identified
}

// getArrayAwareConfigPaths resolves paths, handling array wildcards.
func getArrayAwareConfigPaths(ref []string, comp *component.ComponentDefinition, design *pattern.PatternFile) [][]string {
	if isDirectReference(ref) {
		config := getComponentConfiguration(comp, design)
		value := objectGetNested(config, popFirst(ref), nil)
		if value == nil {
			return nil
		}
		return [][]string{ref}
	}

	directRef := popLast(ref)
	config := getComponentConfiguration(comp, design)
	items := objectGetNested(config, popFirst(directRef), nil)

	arr, ok := items.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}

	var paths [][]string
	for i := range arr {
		if arr[i] == nil {
			continue
		}
		path := make([]string, len(directRef))
		copy(path, directRef)
		path = append(path, fmt.Sprintf("%d", i))
		paths = append(paths, path)
	}
	return paths
}

// aliasRelationshipAlreadyExists checks for duplicate alias relationships.
func aliasRelationshipAlreadyExists(design *pattern.PatternFile, rel *relationship.RelationshipDefinition) bool {
	for _, existing := range design.Relationships {
		if getRelStatus(existing) == StatusDeleted || !isAliasRelationship(existing) {
			continue
		}
		existingTo := getFirstToSelectorItem(existing)
		newTo := getFirstToSelectorItem(rel)
		if existingTo != nil && newTo != nil {
			if sameRelationshipSelectorClause(*existingTo, *newTo) {
				return true
			}
		}
	}
	return false
}

// getFirstToSelectorItem returns the first "to" selector item from a relationship.
func getFirstToSelectorItem(rel *relationship.RelationshipDefinition) *relationship.SelectorItem {
	if rel.Selectors == nil {
		return nil
	}
	for _, ss := range *rel.Selectors {
		if len(ss.Allow.To) > 0 {
			return &ss.Allow.To[0]
		}
	}
	return nil
}

// aliasAddComponentSideEffects creates actions to add alias components.
func aliasAddComponentSideEffects(rel *relationship.RelationshipDefinition, _ *pattern.PatternFile) []PolicyAction {
	if rel.Selectors == nil {
		return nil
	}
	var actions []PolicyAction
	for _, ss := range *rel.Selectors {
		for _, fromSel := range ss.Allow.From {
			if fromSel.RelationshipDefinitionSelectorsPatch == nil || fromSel.RelationshipDefinitionSelectorsPatch.MutatorRef == nil || len(*fromSel.RelationshipDefinitionSelectorsPatch.MutatorRef) == 0 {
				continue
			}
			path := (*fromSel.RelationshipDefinitionSelectorsPatch.MutatorRef)[0]
			length := len(path)
			displayName := ""
			if length >= 2 {
				displayName = fmt.Sprintf("%s.%s", path[length-2], path[length-1])
			}

			newComp := component.ComponentDefinition{
				ID:          *fromSel.ID,
				DisplayName: displayName,
				Component:   component.Component{Kind: selectorItemKind(fromSel)},
				Metadata:    component.ComponentDefinition_Metadata{IsAnnotation: true},
			}
			if fromSel.Model != nil {
				newComp.ModelReference = *fromSel.Model
			}
			actions = append(actions, newAddComponentAction(&newComp))
		}
	}
	return actions
}

// aliasDeleteComponentSideEffects creates actions to delete alias components.
func aliasDeleteComponentSideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	if rel.Selectors == nil {
		return nil
	}
	var actions []PolicyAction
	for _, ss := range *rel.Selectors {
		for _, fromSel := range ss.Allow.From {
			fID := selectorItemID(fromSel)
			comp := componentDeclarationByID(design, fID)
			actions = append(actions, newDeleteComponentAction(fID, comp))
		}
	}
	return actions
}
