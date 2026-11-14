package policies

import (
	"encoding/json"
	"strings"

	patching "github.com/meshery/meshkit/utils/patching"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/tidwall/sjson"
)

// Action operation constants.
const (
	UpdateComponentOp              = "update_component"
	UpdateComponentConfigurationOp = "update_component_configuration"
	RemoveComponentConfigurationOp = "remove_component_configuration"
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
	case RemoveComponentConfigurationOp:
		value["id"] = a.ID
		value["path"] = a.UpdatePath
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

func newComponentRemoveAction(id string, path []string) PolicyAction {
	return PolicyAction{
		Op:         RemoveComponentConfigurationOp,
		ID:         id,
		UpdatePath: path,
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

// isConfigurationPath reports whether a mutator/mutated reference addresses the
// component's configuration rather than the component envelope. A reference is a
// literal path into the component object, so only a leading "configuration"
// segment descends into Configuration; every other path (displayName,
// metadata.labels, component.kind) addresses the envelope. This is the rule
// configurationForComponentAtPath reads through, and the one
// applyComponentConfigPatches must write through, so a mutated value lands in
// the object it is read back from.
func isConfigurationPath(path []string) bool {
	return len(path) > 0 && path[0] == "configuration"
}

// getComponentUpdateOp returns the appropriate update op based on whether the
// path addresses the component configuration or the component envelope.
func getComponentUpdateOp(path []string) string {
	if isConfigurationPath(path) {
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

// applyComponentConfigPatches applies UpdateComponentConfiguration, UpdateComponent,
// and RemoveComponentConfiguration actions to component configurations in the design.
func applyComponentConfigPatches(design *pattern.PatternFile, actions []PolicyAction) {
	type updateEntry struct {
		id    string
		path  []string
		value interface{}
	}

	var updates []updateEntry
	var removes []PolicyAction
	for _, a := range actions {
		if a.ID == "" {
			continue
		}
		switch a.Op {
		case UpdateComponentConfigurationOp, UpdateComponentOp:
			updates = append(updates, updateEntry{id: a.ID, path: a.UpdatePath, value: a.UpdateValue})
		case RemoveComponentConfigurationOp:
			removes = append(removes, a)
		}
	}
	if len(updates) == 0 && len(removes) == 0 {
		return
	}

	for _, comp := range design.Components {
		// A reference is a literal path into the component object: a
		// "configuration"-prefixed path patches Configuration, any other path
		// patches the envelope. Splitting here keeps the write target identical
		// to the read target in configurationForComponentAtPath.
		var configPatches, envelopePatches []patching.Patch
		for _, up := range updates {
			if up.id != comp.ID.String() {
				continue
			}
			if isConfigurationPath(up.path) {
				configPatches = append(configPatches, patching.Patch{Path: up.path[1:], Value: up.value})
				continue
			}
			envelopePatches = append(envelopePatches, patching.Patch{Path: up.path, Value: up.value})
		}
		if len(configPatches) > 0 {
			if updated, err := patching.ApplyPatches(comp.Configuration, configPatches); err == nil {
				comp.Configuration = updated
			}
		}

		var configRemoves, envelopeRemoves [][]string
		for _, rm := range removes {
			if rm.ID != comp.ID.String() {
				continue
			}
			if isConfigurationPath(rm.UpdatePath) {
				configRemoves = append(configRemoves, rm.UpdatePath[1:])
				continue
			}
			envelopeRemoves = append(envelopeRemoves, rm.UpdatePath)
		}
		if len(configRemoves) > 0 {
			comp.Configuration = applyRemovesToMap(comp.Configuration, configRemoves)
		}

		if len(envelopePatches) > 0 || len(envelopeRemoves) > 0 {
			patchComponentEnvelope(comp, envelopePatches, envelopeRemoves)
		}
	}
}

// patchComponentEnvelope applies envelope-targeted patches and removals through
// the same JSON projection of the component that
// configurationForComponentAtPath reads through, then folds the result back onto
// the component. Mirrors the Rego engine, which applies update_component by
// json.patch-ing the component object at the raw reference path.
func patchComponentEnvelope(comp *component.ComponentDefinition, patches []patching.Patch, removes [][]string) {
	envelope, err := toGenericMap(comp)
	if err != nil {
		return
	}
	if len(patches) > 0 {
		envelope, err = patching.ApplyPatches(envelope, patches)
		if err != nil {
			return
		}
	}
	if len(removes) > 0 {
		envelope = applyRemovesToMap(envelope, removes)
	}

	raw, err := json.Marshal(envelope)
	if err != nil {
		return
	}
	// Decode into a fresh value so a removed field is reset rather than left at
	// its pre-mutation value, then restore the fields the JSON projection cannot
	// carry: ModelId and Metadata.AdditionalProperties are the only `json:"-"`
	// members of ComponentDefinition, so a new one has to be restored here too.
	modelID := comp.ModelId
	additional := comp.Metadata.AdditionalProperties
	var folded component.ComponentDefinition
	if err := json.Unmarshal(raw, &folded); err != nil {
		return
	}
	folded.ModelId = modelID
	folded.Metadata.AdditionalProperties = additional
	*comp = folded
}

// applyRemovesToMap removes each path from a generic object via sjson.Delete.
// Paths are relative to the object passed in, so callers strip any leading
// "configuration" segment before calling.
func applyRemovesToMap(obj map[string]interface{}, removes [][]string) map[string]interface{} {
	data, err := json.Marshal(obj)
	if err != nil {
		return obj
	}
	jsonStr := string(data)
	for _, path := range removes {
		if len(path) == 0 {
			continue
		}
		jsonStr, err = sjson.Delete(jsonStr, strings.Join(path, "."))
		if err != nil {
			return obj
		}
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return obj
	}
	return result
}
