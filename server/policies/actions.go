package policies

import (
	"github.com/meshery/meshkit/utils"
	patching "github.com/meshery/meshkit/utils/patching"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// Action operation constants.
const (
	UpdateComponentOp              = "update_component"
	UpdateComponentConfigurationOp = "update_component_configuration"
	DeleteComponentOp              = "delete_component"
	AddComponentOp                 = "add_component"
	UpdateRelationshipOp           = "update_relationship"
	DeleteRelationshipOp           = "delete_relationship"
	AddRelationshipOp              = "add_relationship"
)

// PolicyAction represents an action produced by the evaluation engine.
type PolicyAction struct {
	Op    string                 `json:"op"`
	Value map[string]interface{} `json:"value"`
}

type componentUpdatePayload struct {
	ID    string      `json:"id"`
	Value interface{} `json:"value"`
	Path  []string    `json:"path"`
}

// Action constructors for type-safe action creation.

func newUpdateRelationshipAction(id, path, value string) PolicyAction {
	return PolicyAction{
		Op: UpdateRelationshipOp,
		Value: map[string]interface{}{
			"id":    id,
			"path":  path,
			"value": value,
		},
	}
}

func newDeleteRelationshipAction(id string, relMap map[string]interface{}) PolicyAction {
	return PolicyAction{
		Op: DeleteRelationshipOp,
		Value: map[string]interface{}{
			"id":           id,
			"relationship": relMap,
		},
	}
}

func newComponentUpdateAction(op, id string, path []string, value interface{}) PolicyAction {
	return PolicyAction{
		Op: op,
		Value: map[string]interface{}{
			"id":    id,
			"path":  path,
			"value": value,
		},
	}
}

func newAddComponentAction(compMap map[string]interface{}) PolicyAction {
	return PolicyAction{
		Op:    AddComponentOp,
		Value: map[string]interface{}{"item": compMap},
	}
}

func newDeleteComponentAction(id string, compMap map[string]interface{}) PolicyAction {
	return PolicyAction{
		Op: DeleteComponentOp,
		Value: map[string]interface{}{
			"id":        id,
			"component": compMap,
		},
	}
}

func newAddRelationshipAction(relMap map[string]interface{}) PolicyAction {
	return PolicyAction{
		Op:    AddRelationshipOp,
		Value: map[string]interface{}{"item": relMap},
	}
}

// getComponentUpdateOp returns the appropriate update op based on whether
// the path starts with "configuration".
func getComponentUpdateOp(path []string) string {
	if len(path) > 0 && path[0] == "configuration" {
		return UpdateComponentConfigurationOp
	}
	return UpdateComponentOp
}

// applyActionsToDesign applies all actions to a design, returning a new copy.
func applyActionsToDesign(design *pattern.PatternFile, actions []PolicyAction) *pattern.PatternFile {
	if len(actions) == 0 {
		return design
	}
	result := deepCopyDesign(design)

	deleteCompIDs := make(map[string]bool)
	deleteRelIDs := make(map[string]bool)
	for _, a := range actions {
		switch a.Op {
		case DeleteComponentOp:
			if id := getMapString(a.Value, "id"); id != "" {
				deleteCompIDs[id] = true
			}
		case DeleteRelationshipOp:
			if id := getMapString(a.Value, "id"); id != "" {
				deleteRelIDs[id] = true
			}
		}
	}

	// Filter deleted components.
	if len(deleteCompIDs) > 0 {
		comps := make([]*component.ComponentDefinition, 0, len(result.Components))
		for _, c := range result.Components {
			if !deleteCompIDs[c.ID.String()] {
				comps = append(comps, c)
			}
		}
		result.Components = comps
	}

	// Filter deleted relationships.
	if len(deleteRelIDs) > 0 {
		rels := make([]*relationship.RelationshipDefinition, 0, len(result.Relationships))
		for _, r := range result.Relationships {
			if !deleteRelIDs[r.ID.String()] {
				rels = append(rels, r)
			}
		}
		result.Relationships = rels
	}

	// Add new components.
	existingCompIDs := make(map[string]bool)
	for _, c := range result.Components {
		existingCompIDs[c.ID.String()] = true
	}
	for _, a := range actions {
		if a.Op == AddComponentOp {
			item := getMapMap(a.Value, "item")
			if item == nil {
				continue
			}
			comp, err := mapToComponentDef(item)
			if err != nil {
				continue
			}
			id := comp.ID.String()
			if !existingCompIDs[id] {
				result.Components = append(result.Components, &comp)
				existingCompIDs[id] = true
			}
		}
	}

	// Add new relationships.
	existingRelIDs := make(map[string]bool)
	for _, r := range result.Relationships {
		existingRelIDs[r.ID.String()] = true
	}
	for _, a := range actions {
		if a.Op == AddRelationshipOp {
			item := getMapMap(a.Value, "item")
			if item == nil {
				continue
			}
			rel, err := mapToRelationshipDef(item)
			if err != nil {
				continue
			}
			id := rel.ID.String()
			if !existingRelIDs[id] {
				result.Relationships = append(result.Relationships, &rel)
				existingRelIDs[id] = true
			}
		}
	}

	// Update relationships (status changes).
	for _, a := range actions {
		if a.Op == UpdateRelationshipOp {
			id := getMapString(a.Value, "id")
			path, _ := a.Value["path"].(string)
			value, _ := a.Value["value"].(string)
			if path == "/status" && id != "" {
				for _, rel := range result.Relationships {
					if rel.ID.String() == id {
						setRelStatus(rel, value)
					}
				}
			}
		}
	}

	// Apply component configuration patches.
	applyComponentConfigPatches(result, actions)

	return result
}

// applyComponentConfigPatches applies UpdateComponentConfiguration and UpdateComponent
// actions to component configurations in the design.
func applyComponentConfigPatches(design *pattern.PatternFile, actions []PolicyAction) {
	var updates []componentUpdatePayload
	for _, a := range actions {
		if a.Op != UpdateComponentConfigurationOp && a.Op != UpdateComponentOp {
			continue
		}
		pl, err := utils.MarshalAndUnmarshal[map[string]interface{}, componentUpdatePayload](a.Value)
		if err != nil {
			continue
		}
		updates = append(updates, pl)
	}
	if len(updates) == 0 {
		return
	}

	for _, comp := range design.Components {
		var patches []patching.Patch
		for _, up := range updates {
			if up.ID != comp.ID.String() {
				continue
			}
			path := up.Path
			if len(path) > 0 && path[0] == "configuration" {
				path = path[1:]
			}
			patches = append(patches, patching.Patch{Path: path, Value: up.Value})
		}
		if len(patches) == 0 {
			continue
		}
		updatedConfig, err := patching.ApplyPatches(comp.Configuration, patches)
		if err != nil {
			continue
		}
		comp.Configuration = updatedConfig
	}
}
