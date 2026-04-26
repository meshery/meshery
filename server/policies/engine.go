package policies

import (
	"sync"

	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// Logger is the minimal subset of meshkit/logger.Handler that the engine uses.
// Declared here so js/wasm builds don't transitively pull in gorm/logrus via
// the full Handler interface.
type Logger interface {
	Info(args ...any)
	Warnf(format string, args ...any)
}

// GoEngine is the native Go policy evaluation engine.
type GoEngine struct {
	policies []RelationshipPolicy
	log      Logger

	mu       sync.Mutex
	flapHist map[string][]string // designID -> recent action fingerprints
	lastTies []TieInfo           // ties detected in the most recent evaluation
}

// NewGoEngine creates a new Go policy engine with all built-in policies registered.
func NewGoEngine(log Logger) *GoEngine {
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
		flapHist: make(map[string][]string),
	}
}

func cloneTieInfo(t TieInfo) TieInfo {
	cloned := t
	if t.Path != nil {
		cloned.Path = make([]string, len(t.Path))
		copy(cloned.Path, t.Path)
	}
	return cloned
}

// LastTies returns the ties detected during the most recent evaluation.
func (e *GoEngine) LastTies() []TieInfo {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]TieInfo, len(e.lastTies))
	for i, tie := range e.lastTies {
		out[i] = cloneTieInfo(tie)
	}
	return out
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
			if bridged := bridgeRelationship(r); bridged != nil {
				rels = append(rels, bridged)
			}
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

	if err := e.checkFlapping(design.ID.String(), allActions); err != nil {
		return resp, err
	}
	return resp, nil
}

// hasTrackableDesignID reports whether designID is safe to use as a flapping
// history key. Empty and zero UUID values are treated as "no design ID" and
// should not share history across unrelated designs.
func hasTrackableDesignID(designID string) bool {
	return designID != "" && designID != "00000000-0000-0000-0000-000000000000"
}

// checkFlapping records the fingerprint of this evaluation's actions and
// returns ErrFlappingDetected when the last four evaluations for designID
// alternate between two results.
func (e *GoEngine) checkFlapping(designID string, actions []PolicyAction) error {
	if !hasTrackableDesignID(designID) {
		return nil
	}
	fp := actionsFingerprint(actions)
	e.mu.Lock()
	flap := recordAndCheckFlap(e.flapHist, designID, fp)
	hist := append([]string(nil), e.flapHist[designID]...)
	e.mu.Unlock()
	if flap {
		e.log.Warnf("flapping detected for design %s", designID)
		return ErrFlappingDetected(designID)
	}
	if len(hist) >= 3 && hist[len(hist)-1] == hist[len(hist)-3] && hist[len(hist)-1] != hist[len(hist)-2] {
		e.log.Warnf("flapping candidate for design %s (1 inversion)", designID)
	}
	return nil
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

	// Phase 2.5: Identify inventory parents the design is missing. Mirrors
	// identify_additions.rego — without this the Go engine silently skips
	// auto-creating Namespaces for namespaced resources, diverging from the
	// Rego engine's behavior.
	inventoryActions := identifyInventoryAdditions(designWithIdentified, relsInScope)
	if len(inventoryActions) > 0 {
		combinedActions = append(combinedActions, inventoryActions...)
		designWithIdentified = applyActionsToDesign(design, combinedActions)
	}

	// Phase 3: Generate and apply actions.
	var phaseActions []PolicyAction
	for _, policy := range e.policies {
		actions := generateActionsToApplyOnDesign(designWithIdentified, policy)
		phaseActions = append(phaseActions, actions...)
	}

	phaseActions, ties := resolveTies(phaseActions)
	e.mu.Lock()
	e.lastTies = ties
	e.mu.Unlock()
	for _, t := range ties {
		e.log.Warnf("tie detected on component %s path %v; dropping conflicting actions", t.ComponentID, t.Path)
	}

	allActions = append(allActions, identifyActions...)
	allActions = append(allActions, inventoryActions...)
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
