package policies

import (
	"encoding/json"

	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils"
	patching "github.com/meshery/meshkit/utils/patching"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
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
		resp.Actions = append(resp.Actions, pattern.Action{
			Op:    action.Op,
			Value: action.Value,
		})
	}

	resp.Trace = buildTrace(allActions, &design, resultDesign)

	return resp, nil
}

type componentUpdatePayload struct {
	ID    string      `json:"id"`
	Value interface{} `json:"value"`
	Path  []string    `json:"path"`
}

// ApplyConfigurationPatches applies component configuration update patches to the response.
func ApplyConfigurationPatches(log logger.Handler, resp *pattern.EvaluationResponse) {
	var updates []componentUpdatePayload
	for _, action := range resp.Actions {
		if action.Op == UpdateComponentConfigurationOp || action.Op == UpdateComponentOp {
			pl, err := utils.MarshalAndUnmarshal[map[string]interface{}, componentUpdatePayload](action.Value)
			if err != nil {
				log.Warn(ErrParsePayload(err))
				continue
			}
			updates = append(updates, pl)
		}
	}

	for _, comp := range resp.Design.Components {
		var patches []patching.Patch
		for _, up := range updates {
			if up.ID == comp.ID.String() {
				path := up.Path
				if len(path) > 0 && path[0] == "configuration" {
					path = path[1:]
				}
				patches = append(patches, patching.Patch{Path: path, Value: up.Value})
			}
		}
		if len(patches) == 0 {
			continue
		}
		updatedConfig, err := patching.ApplyPatches(comp.Configuration, patches)
		if err != nil {
			log.Error(ErrApplyPatch(err))
			continue
		}
		comp.Configuration = updatedConfig
	}
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
			item := getMapMap(action.Value, "item")
			if item == nil {
				continue
			}
			comp, err := mapToComponentDef(item)
			if err == nil {
				trace.ComponentsAdded = append(trace.ComponentsAdded, comp)
			}

		case DeleteComponentOp:
			id := getMapString(action.Value, "id")
			comp := componentDeclarationByID(originalDesign, id)
			if comp != nil {
				trace.ComponentsRemoved = append(trace.ComponentsRemoved, *comp)
			}

		case UpdateComponentOp, UpdateComponentConfigurationOp:
			id := getMapString(action.Value, "id")
			if id == "" || updatedCompIDs[id] {
				continue
			}
			updatedCompIDs[id] = true
			comp := componentDeclarationByID(finalDesign, id)
			if comp != nil {
				trace.ComponentsUpdated = append(trace.ComponentsUpdated, *comp)
			}

		case AddRelationshipOp:
			item := getMapMap(action.Value, "item")
			if item == nil {
				continue
			}
			rel, err := mapToRelationshipDef(item)
			if err == nil {
				trace.RelationshipsAdded = append(trace.RelationshipsAdded, rel)
			}

		case DeleteRelationshipOp:
			relMap := getMapMap(action.Value, "relationship")
			if relMap != nil {
				rel, err := mapToRelationshipDef(relMap)
				if err == nil {
					trace.RelationshipsRemoved = append(trace.RelationshipsRemoved, rel)
				}
			} else {
				id := getMapString(action.Value, "id")
				rel := findRelationshipByID(originalDesign, id)
				if rel != nil {
					trace.RelationshipsRemoved = append(trace.RelationshipsRemoved, *rel)
				}
			}

		case UpdateRelationshipOp:
			id := getMapString(action.Value, "id")
			if id == "" || updatedRelIDs[id] {
				continue
			}
			rel := findRelationshipByID(finalDesign, id)
			if rel == nil {
				continue
			}
			updatedRelIDs[id] = true
			trace.RelationshipsUpdated = append(trace.RelationshipsUpdated, *rel)
		}
	}
	return trace
}

