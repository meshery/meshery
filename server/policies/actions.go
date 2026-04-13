package policies

import (
	patching "github.com/meshery/meshkit/utils/patching"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
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
// All fields are typed to avoid map[string]interface{} roundtrips internally.
// Use toPatternAction() to convert to the API response format.
type PolicyAction struct {
	Op           string
	ID           string
	Component    *component.ComponentDefinition
	Relationship *relationship.RelationshipDefinition
	// For UpdateRelationshipOp: path and value are strings.
	Path        string
	StringValue string
	// For component update ops (UpdateComponentOp, UpdateComponentConfigurationOp).
	UpdatePath  []string
	UpdateValue interface{}
}

// toPatternAction converts a typed PolicyAction to the API response format.
func (a PolicyAction) toPatternAction() pattern.Action {
	value := make(map[string]interface{})
	switch a.Op {
	case UpdateRelationshipOp:
		value["id"] = a.ID
		value["path"] = a.Path
		value["value"] = a.StringValue
	case DeleteRelationshipOp:
		value["id"] = a.ID
		if a.Relationship != nil {
			if m, err := toGenericMap(a.Relationship); err == nil {
				value["relationship"] = m
			}
		}
	case UpdateComponentOp, UpdateComponentConfigurationOp:
		value["id"] = a.ID
		value["path"] = a.UpdatePath
		value["value"] = a.UpdateValue
	case AddComponentOp:
		if a.Component != nil {
			if m, err := toGenericMap(a.Component); err == nil {
				value["item"] = m
			}
		}
	case DeleteComponentOp:
		value["id"] = a.ID
		if a.Component != nil {
			if m, err := toGenericMap(a.Component); err == nil {
				value["component"] = m
			}
		}
	case AddRelationshipOp:
		if a.Relationship != nil {
			if m, err := toGenericMap(a.Relationship); err == nil {
				value["item"] = m
			}
		}
	}
	return pattern.Action{Op: a.Op, Value: value}
}

// Action constructors.

func newUpdateRelationshipAction(id, path, value string) PolicyAction {
	return PolicyAction{
		Op:          UpdateRelationshipOp,
		ID:          id,
		Path:        path,
		StringValue: value,
	}
}

func newDeleteRelationshipAction(id string, rel *relationship.RelationshipDefinition) PolicyAction {
	return PolicyAction{
		Op:           DeleteRelationshipOp,
		ID:           id,
		Relationship: rel,
	}
}

func newComponentUpdateAction(op, id string, path []string, value interface{}) PolicyAction {
	return PolicyAction{
		Op:          op,
		ID:          id,
		UpdatePath:  path,
		UpdateValue: value,
	}
}

func newAddComponentAction(comp *component.ComponentDefinition) PolicyAction {
	return PolicyAction{
		Op:        AddComponentOp,
		Component: comp,
	}
}

func newDeleteComponentAction(id string, comp *component.ComponentDefinition) PolicyAction {
	return PolicyAction{
		Op:        DeleteComponentOp,
		ID:        id,
		Component: comp,
	}
}

func newAddRelationshipAction(rel *relationship.RelationshipDefinition) PolicyAction {
	return PolicyAction{
		Op:           AddRelationshipOp,
		Relationship: rel,
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
			if a.ID != "" {
				deleteCompIDs[a.ID] = true
			}
		case DeleteRelationshipOp:
			if a.ID != "" {
				deleteRelIDs[a.ID] = true
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
		if a.Op == AddComponentOp && a.Component != nil {
			id := a.Component.ID.String()
			if !existingCompIDs[id] {
				compCopy := *a.Component
				result.Components = append(result.Components, &compCopy)
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
		if a.Op == AddRelationshipOp && a.Relationship != nil {
			id := a.Relationship.ID.String()
			if !existingRelIDs[id] {
				relCopy := *a.Relationship
				result.Relationships = append(result.Relationships, &relCopy)
				existingRelIDs[id] = true
			}
		}
	}

	// Update relationships (status changes).
	for _, a := range actions {
		if a.Op == UpdateRelationshipOp {
			if a.Path == "/status" && a.ID != "" {
				for _, rel := range result.Relationships {
					if rel.ID.String() == a.ID {
						setRelStatus(rel, a.StringValue)
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
	type updateEntry struct {
		id    string
		path  []string
		value interface{}
	}

	var updates []updateEntry
	for _, a := range actions {
		if a.Op != UpdateComponentConfigurationOp && a.Op != UpdateComponentOp {
			continue
		}
		if a.ID == "" {
			continue
		}
		updates = append(updates, updateEntry{id: a.ID, path: a.UpdatePath, value: a.UpdateValue})
	}
	if len(updates) == 0 {
		return
	}

	for _, comp := range design.Components {
		var patches []patching.Patch
		for _, up := range updates {
			if up.id != comp.ID.String() {
				continue
			}
			path := up.path
			if len(path) > 0 && path[0] == "configuration" {
				path = path[1:]
			}
			patches = append(patches, patching.Patch{Path: path, Value: up.value})
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
