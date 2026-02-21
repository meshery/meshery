package policies

import (
	"encoding/json"

	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils"
	patching "github.com/meshery/meshkit/utils/patching"
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
		},
	}
}

// EvaluateDesign evaluates a design file using the Go policy engine.
func (e *GoEngine) EvaluateDesign(
	design pattern.PatternFile,
	registeredRelationships []interface{},
) (pattern.EvaluationResponse, error) {
	var resp pattern.EvaluationResponse

	designMap, err := toGenericMap(design)
	if err != nil {
		return resp, ErrConvertDesign(err)
	}

	var relMaps []map[string]interface{}
	for _, r := range registeredRelationships {
		rm, err := toGenericMapFromInterface(r)
		if err != nil {
			continue
		}
		relMaps = append(relMaps, rm)
	}

	modelsInDesign := getModelsInDesign(designMap)
	relsInScope := filterRelationshipsInScope(relMaps, modelsInDesign, designMap)

	e.log.Info("models in design: ", len(modelsInDesign), ", registered rels: ", len(relMaps), ", rels in scope: ", len(relsInScope))

	resultDesign, allActions := e.evaluate(designMap, relsInScope)

	resp, err = fromGenericMap(resultDesign)
	if err != nil {
		return resp, ErrConvertResult(err)
	}

	for _, action := range allActions {
		resp.Actions = append(resp.Actions, pattern.Action{
			Op:    action.Op,
			Value: action.Value,
		})
	}

	applyConfigurationPatches(e.log, &resp)

	return resp, nil
}

type componentUpdatePayload struct {
	ID    string      `json:"id"`
	Value interface{} `json:"value"`
	Path  []string    `json:"path"`
}

// applyConfigurationPatches applies component configuration update patches to the response.
func applyConfigurationPatches(log logger.Handler, resp *pattern.EvaluationResponse) {
	var updates []componentUpdatePayload
	for _, action := range resp.Actions {
		if action.Op == UpdateComponentConfigurationOp {
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
			if up.ID == comp.Id.String() {
				patches = append(patches, patching.Patch{Path: up.Path[1:], Value: up.Value})
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

// evaluate runs the full evaluation pipeline on a design map.
func (e *GoEngine) evaluate(designMap map[string]interface{}, relsInScope []map[string]interface{}) (map[string]interface{}, []PolicyAction) {
	var allActions []PolicyAction

	// Phase 1: Validate existing relationships
	for _, policy := range e.policies {
		validationActions := validateRelationshipsInDesign(designMap, policy)
		allActions = append(allActions, validationActions...)
	}

	designWithValidatedRels := applyAllActionsToDesign(designMap, allActions)

	// Phase 1.5: Identify missing components (additions)
	var additionActions []PolicyAction
	for _, relDef := range relsInScope {
		additionActions = append(additionActions, identifyAdditions(relDef, designWithValidatedRels)...)
	}
	if len(additionActions) > 0 {
		allActions = append(allActions, additionActions...)
		designWithValidatedRels = applyAllActionsToDesign(designWithValidatedRels, additionActions)
	}

	// Phase 2: Identify new relationships
	var identifyActions []PolicyAction
	for _, policy := range e.policies {
		actions := identifyRelationshipsInDesign(designWithValidatedRels, relsInScope, policy)
		identifyActions = append(identifyActions, actions...)
	}

	combinedActions := append(allActions, identifyActions...)
	designWithIdentified := applyAllActionsToDesign(designMap, combinedActions)

	// Phase 3: Generate and apply actions
	designForActions := designWithIdentified

	var phaseActions []PolicyAction
	for _, policy := range e.policies {
		actions := generateActionsToApplyOnDesign(designForActions, policy)
		phaseActions = append(phaseActions, actions...)
	}

	allActions = append(allActions, identifyActions...)
	allActions = append(allActions, phaseActions...)

	finalDesign := applyAllActionsToDesign(designForActions, phaseActions)

	e.log.Info("evaluation complete: ", countItems(finalDesign, "relationships"), " relationships, ", len(allActions), " actions")

	return finalDesign, allActions
}

// getModelsInDesign extracts unique model references from design components.
func getModelsInDesign(design map[string]interface{}) []map[string]interface{} {
	comps := getMapSlice(design, "components")
	seen := make(map[string]bool)
	var models []map[string]interface{}

	for _, c := range comps {
		comp, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		model := getMapMap(comp, "model")
		if model == nil {
			continue
		}
		name := getMapString(model, "name")
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		models = append(models, model)
	}
	return models
}

// filterRelationshipsInScope filters relationships to those matching models in the design
// and respecting user preferences.
func filterRelationshipsInScope(
	allRels []map[string]interface{},
	modelsInDesign []map[string]interface{},
	design map[string]interface{},
) []map[string]interface{} {
	preferences := getMapMap(design, "preferences")
	var layerPrefs map[string]interface{}
	if preferences != nil {
		layers := getMapMap(preferences, "layers")
		if layers != nil {
			layerPrefs = getMapMap(layers, "relationships")
		}
	}

	var result []map[string]interface{}
	for _, rel := range allRels {
		relModel := getMapMap(rel, "model")
		relModelName := getMapString(relModel, "name")

		matched := false
		for _, model := range modelsInDesign {
			if getMapString(model, "name") == relModelName {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}

		if layerPrefs != nil {
			key := relationshipPreferenceKey(rel)
			if v, exists := layerPrefs[key]; exists {
				if b, ok := v.(bool); ok && !b {
					continue
				}
			}
		}

		result = append(result, rel)
	}
	return result
}

func countItems(design map[string]interface{}, key string) int {
	return len(getMapSlice(design, key))
}

func toGenericMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func toGenericMapFromInterface(v interface{}) (map[string]interface{}, error) {
	if m, ok := v.(map[string]interface{}); ok {
		return m, nil
	}
	return toGenericMap(v)
}

func fromGenericMap(designMap map[string]interface{}) (pattern.EvaluationResponse, error) {
	var resp pattern.EvaluationResponse

	data, err := json.Marshal(designMap)
	if err != nil {
		return resp, err
	}

	var designFile pattern.PatternFile
	if err := json.Unmarshal(data, &designFile); err != nil {
		return resp, err
	}

	resp.Design = designFile
	resp.Trace = pattern.Trace{}

	return resp, nil
}
