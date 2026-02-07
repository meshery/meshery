package policies

import "strings"

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

// getComponentUpdateOp returns the appropriate update op based on whether
// the path starts with "configuration".
func getComponentUpdateOp(path []string) string {
	if len(path) > 0 && path[0] == "configuration" {
		return UpdateComponentConfigurationOp
	}
	return UpdateComponentOp
}

// applyUpdatesToItem applies all update actions matching the given op and item ID
// to the original item map, returning the updated version.
func applyUpdatesToItem(original map[string]interface{}, actions []PolicyAction, op string) map[string]interface{} {
	id := getMapString(original, "id")
	if id == "" {
		return original
	}

	applied := false
	result := deepCopyMap(original)

	for _, action := range actions {
		if action.Op != op {
			continue
		}
		actionID := getMapString(action.Value, "id")
		if actionID != id {
			continue
		}

		path := action.Value["path"]
		value := action.Value["value"]

		pathStr, isStr := path.(string)
		if isStr && len(pathStr) > 0 {
			// JSON pointer style path like "/status"
			segments := splitJSONPointer(pathStr)
			setNestedValue(result, segments, value)
			applied = true
		}

		pathSlice, isSlice := path.([]interface{})
		if isSlice {
			segments := toStringSlice(pathSlice)
			setNestedValue(result, segments, value)
			applied = true
		}

		pathStrSlice, isStrSlice := path.([]string)
		if isStrSlice {
			setNestedValue(result, pathStrSlice, value)
			applied = true
		}
	}

	if !applied {
		return original
	}
	return result
}

// splitJSONPointer splits a JSON pointer like "/status" into ["status"].
func splitJSONPointer(p string) []string {
	if p == "" {
		return nil
	}
	if p[0] == '/' {
		p = p[1:]
	}
	if p == "" {
		return nil
	}
	parts := make([]string, 0)
	for _, s := range strings.Split(p, "/") {
		if s != "" {
			parts = append(parts, s)
		}
	}
	return parts
}

// setNestedValue sets a value in a nested map at the given path.
func setNestedValue(m map[string]interface{}, path []string, value interface{}) {
	if len(path) == 0 {
		return
	}
	if len(path) == 1 {
		m[path[0]] = value
		return
	}

	key := path[0]
	child, ok := m[key]
	if !ok {
		child = make(map[string]interface{})
		m[key] = child
	}

	childMap, ok := child.(map[string]interface{})
	if !ok {
		childMap = make(map[string]interface{})
		m[key] = childMap
	}
	setNestedValue(childMap, path[1:], value)
}

// filterOutDeleted removes items whose IDs match any delete action.
func filterOutDeleted(items []map[string]interface{}, actions []PolicyAction, deleteOp string) []map[string]interface{} {
	deletedIDs := make(map[string]bool)
	for _, action := range actions {
		if action.Op == deleteOp {
			id := getMapString(action.Value, "id")
			if id != "" {
				deletedIDs[id] = true
			}
		}
	}

	if len(deletedIDs) == 0 {
		return items
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		id := getMapString(item, "id")
		if !deletedIDs[id] {
			result = append(result, item)
		}
	}
	return result
}

// collectAddedItems collects items from add actions.
func collectAddedItems(actions []PolicyAction, addOp string) []map[string]interface{} {
	var result []map[string]interface{}
	for _, action := range actions {
		if action.Op != addOp {
			continue
		}
		item := getMapMap(action.Value, "item")
		if item != nil {
			result = append(result, item)
		}
	}
	return result
}

// applyAllActionsToDesign applies all actions to a design map.
func applyAllActionsToDesign(design map[string]interface{}, actions []PolicyAction) map[string]interface{} {
	result := deepCopyMap(design)

	// Delete components
	comps := extractMapSlice(result, "components")
	comps = filterOutDeleted(comps, actions, DeleteComponentOp)

	// Add components
	addedComps := collectAddedItems(actions, AddComponentOp)
	comps = appendUniqueByID(comps, addedComps)

	// Update components
	for i, comp := range comps {
		comps[i] = applyUpdatesToItem(comp, actions, UpdateComponentOp)
	}

	result["components"] = mapsToSlice(comps)

	// Delete relationships
	rels := extractMapSlice(result, "relationships")
	rels = filterOutDeleted(rels, actions, DeleteRelationshipOp)

	// Add relationships
	addedRels := collectAddedItems(actions, AddRelationshipOp)
	rels = appendUniqueByID(rels, addedRels)

	// Update relationships
	for i, rel := range rels {
		rels[i] = applyUpdatesToItem(rel, actions, UpdateRelationshipOp)
	}

	result["relationships"] = mapsToSlice(rels)

	return result
}

// extractMapSlice extracts a slice of maps from a design field.
func extractMapSlice(design map[string]interface{}, key string) []map[string]interface{} {
	raw, ok := design[key]
	if !ok {
		return nil
	}
	slice, ok := raw.([]interface{})
	if !ok {
		// might already be []map[string]interface{}
		if ms, ok := raw.([]map[string]interface{}); ok {
			return ms
		}
		return nil
	}
	result := make([]map[string]interface{}, 0, len(slice))
	for _, v := range slice {
		if m, ok := v.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result
}

// mapsToSlice converts []map[string]interface{} to []interface{}.
func mapsToSlice(items []map[string]interface{}) []interface{} {
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = item
	}
	return result
}

// appendUniqueByID appends items to a slice, skipping those whose ID already exists.
func appendUniqueByID(existing, toAdd []map[string]interface{}) []map[string]interface{} {
	ids := make(map[string]bool, len(existing))
	for _, item := range existing {
		id := getMapString(item, "id")
		if id != "" {
			ids[id] = true
		}
	}
	for _, item := range toAdd {
		id := getMapString(item, "id")
		if id == "" || !ids[id] {
			existing = append(existing, item)
			if id != "" {
				ids[id] = true
			}
		}
	}
	return existing
}
