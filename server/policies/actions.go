package policies

import (
	"encoding/json"
	"fmt"
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

// applyActionsToDesign applies all actions to a design, returning a new copy
// and any patch that could not be applied. A failed patch leaves the mutated
// field unchanged, so the next evaluation re-emits the same action; returning
// the error lets the caller report that instead of silently not converging.
func applyActionsToDesign(design *pattern.PatternFile, actions []PolicyAction) (*pattern.PatternFile, []error) {
	if len(actions) == 0 {
		return design, nil
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
	errs := applyComponentConfigPatches(result, actions)

	return result, errs
}

// applyComponentConfigPatches applies UpdateComponentConfiguration, UpdateComponent,
// and RemoveComponentConfiguration actions to component configurations in the design.
func applyComponentConfigPatches(design *pattern.PatternFile, actions []PolicyAction) []error {
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
		return nil
	}

	var errs []error

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
			updated, err := patching.ApplyPatches(comp.Configuration, configPatches)
			if err != nil {
				errs = append(errs, fmt.Errorf("patching configuration of component %s: %w", comp.ID, err))
			} else {
				comp.Configuration = updated
			}
		}

		var configRemoves, envelopeRemoves []PolicyAction
		for _, rm := range removes {
			if rm.ID != comp.ID.String() {
				continue
			}
			if isConfigurationPath(rm.UpdatePath) {
				configRemoves = append(configRemoves, rm)
				continue
			}
			envelopeRemoves = append(envelopeRemoves, rm)
		}
		if len(configRemoves) > 0 {
			comp.Configuration = applyConfigRemoves(comp.Configuration, orderCleanupActions(configRemoves))
		}

		if len(envelopePatches) > 0 || len(envelopeRemoves) > 0 {
			if err := patchComponentEnvelope(comp, envelopePatches, orderCleanupActions(envelopeRemoves)); err != nil {
				errs = append(errs, fmt.Errorf("patching envelope of component %s: %w", comp.ID, err))
			}
		}
	}
	return errs
}

// patchComponentEnvelope applies envelope-targeted patches and removals through
// the same JSON projection of the component that
// configurationForComponentAtPath reads through, then folds the result back onto
// the component. Mirrors the Rego engine, which applies update_component by
// json.patch-ing the component object at the raw reference path.
func patchComponentEnvelope(comp *component.ComponentDefinition, patches []patching.Patch, removes []PolicyAction) error {
	envelope, err := toGenericMap(comp)
	if err != nil {
		return fmt.Errorf("projecting component to JSON: %w", err)
	}
	if len(patches) > 0 {
		envelope, err = patching.ApplyPatches(envelope, patches)
		if err != nil {
			return fmt.Errorf("applying envelope patches: %w", err)
		}
	}
	if len(removes) > 0 {
		// Envelope paths carry no "configuration" prefix, so applyConfigRemoves
		// addresses them against the envelope as-is, with the same per-segment
		// escaping the configuration side uses: a dotted key such as the label
		// "kubernetes.io/service-name" stays one segment.
		envelope = applyConfigRemoves(envelope, removes)
	}

	raw, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("encoding patched envelope: %w", err)
	}
	// Decode into a fresh value so a removed field is reset rather than left at
	// its pre-mutation value. ModelId is the one member the projection drops:
	// it is `json:"-"` and ComponentDefinition has no custom marshaller, unlike
	// ComponentDefinition_Metadata, whose MarshalJSON/UnmarshalJSON inline
	// AdditionalProperties as ordinary keys - so a mutation of an unmodelled
	// metadata key round-trips into the decoded map and must be kept.
	modelID := comp.ModelId
	var folded component.ComponentDefinition
	if err := json.Unmarshal(raw, &folded); err != nil {
		return fmt.Errorf("decoding patched envelope: %w", err)
	}
	folded.ModelId = modelID
	*comp = folded
	return nil
}

// sjsonPathSpecials are the characters gjson/sjson path syntax assigns meaning
// to inside a key. Escaping them makes a path segment address a literal key.
// The colon is included because sjson reads a leading ":" as force-object-key
// syntax, so a literal key starting with ":" would otherwise be rewritten.
//
// This mirrors the escaping meshkit's patching package applies on the set side,
// so a path means the same thing whether a relationship patch writes it or a
// relationship delete removes it.
const sjsonPathSpecials = `.*?|#@:`

// sjsonPath joins a configuration path into an sjson path, escaping specials in
// each segment so every segment addresses one literal key. Segments without
// specials are emitted unchanged, so ordinary nested paths and numeric array
// indices keep their existing meaning.
func sjsonPath(path []string) string {
	escaped := make([]string, len(path))
	for i, segment := range path {
		escaped[i] = escapeSjsonKey(segment)
	}
	return strings.Join(escaped, ".")
}

// escapeSjsonKey backslash-escapes gjson/sjson path specials in a single key.
// The escape character itself is escaped first so backslashes already present
// in a key survive the round trip.
func escapeSjsonKey(key string) string {
	if !strings.ContainsAny(key, sjsonPathSpecials+`\`) {
		return key
	}
	var b strings.Builder
	b.Grow(len(key) + 4)
	for _, r := range key {
		if r == '\\' || strings.ContainsRune(sjsonPathSpecials, r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// applyConfigRemoves removes fields from a component's configuration via sjson.Delete.
// Duplicate removals targeting the same component ID and complete path are deduplicated
// so that actions from multiple selectors do not shift and delete subsequent array elements.
func applyConfigRemoves(config map[string]interface{}, removes []PolicyAction) map[string]interface{} {
	data, err := json.Marshal(config)
	if err != nil {
		return config
	}
	jsonStr := string(data)
	seen := make(map[string]bool)
	for _, rm := range removes {
		path := rm.UpdatePath
		if len(path) > 0 && path[0] == "configuration" {
			path = path[1:]
		}
		if len(path) == 0 {
			continue
		}
		sjPath := sjsonPath(path)
		dedupKey := rm.ID + ":" + sjPath
		if seen[dedupKey] {
			continue
		}
		seen[dedupKey] = true
		jsonStr, err = sjson.Delete(jsonStr, sjPath)
		if err != nil {
			return config
		}
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return config
	}
	return result
}
