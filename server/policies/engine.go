package policies

import (
	"encoding/json"

	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// GoEngine is the native Go policy evaluation engine.
type GoEngine struct {
	policies []RelationshipPolicy
	log      logger.Handler
}

// NewGoEngine creates a new Go policy engine with all built-in policies registered.
func NewGoEngine(log logger.Handler) *GoEngine {
	return &GoEngine{
		log: log,
		policies: []RelationshipPolicy{
			&MatchLabelsPolicy{},
			&AliasPolicy{},
			&EdgeNonBindingPolicy{},
			&EdgeBindingPolicy{},
			&HierarchicalParentChildPolicy{},
			&HierarchicalWalletPolicy{},
		},
	}
}

// ConvertRelationships converts entity interfaces to typed relationship definitions.
func ConvertRelationships(registeredRelationships []interface{}) []*relationship.RelationshipDefinition {
	var rels []*relationship.RelationshipDefinition
	for _, r := range registeredRelationships {
		switch v := r.(type) {
		case *relationship.RelationshipDefinition:
			rels = append(rels, v)
		case relationship.RelationshipDefinition:
			rels = append(rels, &v)
		default:
			data, err := json.Marshal(r)
			if err != nil {
				continue
			}
			var rel relationship.RelationshipDefinition
			if err := json.Unmarshal(data, &rel); err != nil {
				continue
			}
			rels = append(rels, &rel)
		}
	}
	return rels
}

// EvaluateDesign evaluates a design file using the Go policy engine.
func (e *GoEngine) EvaluateDesign(
	design pattern.PatternFile,
	registeredRelationships []*relationship.RelationshipDefinition,
) (pattern.EvaluationResponse, error) {
	var resp pattern.EvaluationResponse

	// Ensure ModelReference.Name is populated for all components.
	// Some designs from cloud only have Model.Name set, not ModelReference.Name.
	for _, comp := range design.Components {
		if comp.ModelReference.Name == "" && comp.Model != nil {
			comp.ModelReference.Name = comp.Model.Name
		}
	}

	modelsInDesign := getModelsInDesign(&design)
	relsInScope := filterRelationshipsInScope(registeredRelationships, modelsInDesign, &design)

	e.log.Info("models in design: ", len(modelsInDesign), ", registered rels: ", len(registeredRelationships), ", rels in scope: ", len(relsInScope))

	resultDesign, allActions := e.evaluate(&design, relsInScope)

	resp.Design = *resultDesign

	for _, action := range allActions {
		resp.Actions = append(resp.Actions, action.toPatternAction())
	}

	resp.Trace = buildTrace(allActions, &design, resultDesign)

	return resp, nil
}

// evaluate runs the full evaluation pipeline on a design.
func (e *GoEngine) evaluate(design *pattern.PatternFile, relsInScope []*relationship.RelationshipDefinition) (*pattern.PatternFile, []PolicyAction) {
	var allActions []PolicyAction

	// Phase 1: Validate existing relationships.
	for _, policy := range e.policies {
		validationActions := validateRelationshipsInDesign(design, policy)
		allActions = append(allActions, validationActions...)
	}

	designWithValidatedRels := applyActionsToDesign(design, allActions)

	// Phase 2: Identify new relationships.
	var identifyActions []PolicyAction
	for _, policy := range e.policies {
		actions := identifyRelationshipsInDesign(designWithValidatedRels, relsInScope, policy)
		identifyActions = append(identifyActions, actions...)
	}

	combinedActions := append(allActions, identifyActions...)
	designWithIdentified := applyActionsToDesign(design, combinedActions)

	// Phase 3: Generate and apply actions.
	var phaseActions []PolicyAction
	for _, policy := range e.policies {
		actions := generateActionsToApplyOnDesign(designWithIdentified, policy)
		phaseActions = append(phaseActions, actions...)
	}

	allActions = append(allActions, identifyActions...)
	allActions = append(allActions, phaseActions...)

	finalDesign := applyActionsToDesign(designWithIdentified, phaseActions)

	e.log.Info("evaluation complete: ", len(finalDesign.Relationships), " relationships, ", len(allActions), " actions")

	return finalDesign, allActions
}

// buildTrace constructs a Trace from the actions produced during evaluation.
func buildTrace(actions []PolicyAction, originalDesign, finalDesign *pattern.PatternFile) pattern.Trace {
	var trace pattern.Trace
	updatedCompIDs := make(map[string]bool)
	updatedRelIDs := make(map[string]bool)

	for _, action := range actions {
		switch action.Op {
		case AddComponentOp:
			if action.Component != nil {
				trace.ComponentsAdded = append(trace.ComponentsAdded, *action.Component)
			}

		case DeleteComponentOp:
			comp := componentDeclarationByID(originalDesign, action.ID)
			if comp != nil {
				trace.ComponentsRemoved = append(trace.ComponentsRemoved, *comp)
			}

		case UpdateComponentOp, UpdateComponentConfigurationOp:
			if action.ID == "" || updatedCompIDs[action.ID] {
				continue
			}
			updatedCompIDs[action.ID] = true
			comp := componentDeclarationByID(finalDesign, action.ID)
			if comp != nil {
				trace.ComponentsUpdated = append(trace.ComponentsUpdated, *comp)
			}

		case AddRelationshipOp:
			if action.Relationship != nil {
				trace.RelationshipsAdded = append(trace.RelationshipsAdded, *action.Relationship)
			}

		case DeleteRelationshipOp:
			if action.Relationship != nil {
				trace.RelationshipsRemoved = append(trace.RelationshipsRemoved, *action.Relationship)
			} else {
				rel := findRelationshipByID(originalDesign, action.ID)
				if rel != nil {
					trace.RelationshipsRemoved = append(trace.RelationshipsRemoved, *rel)
				}
			}

		case UpdateRelationshipOp:
			if action.ID == "" || updatedRelIDs[action.ID] {
				continue
			}
			rel := findRelationshipByID(finalDesign, action.ID)
			if rel == nil {
				continue
			}
			updatedRelIDs[action.ID] = true
			trace.RelationshipsUpdated = append(trace.RelationshipsUpdated, *rel)
		}
	}
	return trace
}

