package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

// WorkspaceFilter is the decoded `filter` query parameter of the
// workspace-scoped list endpoints:
//
//	GET /api/workspaces/{id}/{environments,designs,views,teams}
//
// The wire form is a JSON object - the grammar Meshery UI has always sent
// (`{"assigned":false}`). It is deliberately NOT the space-separated
// "key value" form that helpers/utils.ApplyFilters reads: these endpoints used
// to hand the same string to both parsers, so whichever grammar a caller
// picked, the other one rejected or silently ignored it. Decoding once, here,
// is what keeps the two from fighting. See issue #21826 for the same defect in
// GetEnvironmentConnections.
type WorkspaceFilter struct {
	// Assigned selects the resources assigned to the workspace (true, the
	// default when the caller says nothing) or the unassigned complement.
	Assigned bool

	// DeletedAt selects soft-deleted rows instead of live ones. It is not
	// honoured by the designs listing, whose `meshery_patterns` table has no
	// `deleted_at` column.
	DeletedAt bool

	// Owner, when non-empty, narrows the result to rows with that owner. Only
	// the environments and teams listings can serve it; see
	// ParseWorkspaceFilter for why the other two cannot.
	Owner string
}

// ParseWorkspaceFilter decodes the `filter` query parameter.
//
// An empty filter yields the default: assigned, live resources, unnarrowed.
// Any other input must be a JSON object whose recognised fields carry the right
// type; anything else is a caller error and is reported as one rather than
// surfacing a bare encoding/json message.
//
// Unrecognised fields are ignored, so a caller sending a field this endpoint
// does not serve gets the default behaviour rather than a rejection.
//
// `owner` is decoded here for every caller but is only applied by the
// environments and teams listings, because only those two tables carry an
// `owner` column. `meshery_patterns` has no owner column at all (its
// `MesheryPattern.UserID` is `gorm:"-"`, so it is never persisted), and on
// `meshery_views` the field is named UserID, from which gorm derives
// `user_id` - the `db:"owner"` tag never reaches gorm. The previous
// ApplyFilters call listed `owner` and `organization_id` for all four, which
// could only have produced a SQL error on designs and views had it ever been
// reachable.
//
// `organization_id` is dropped everywhere rather than carried forward.
// `meshery_patterns` has no such column; where the column does exist - on
// `environments` and `meshery_views` - the enclosing workspace already scopes
// the listing, so the key adds nothing a caller could use.
func ParseWorkspaceFilter(filter string) (WorkspaceFilter, error) {
	parsed := WorkspaceFilter{Assigned: true}

	if strings.TrimSpace(filter) == "" {
		return parsed, nil
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(filter), &raw); err != nil {
		return parsed, ErrInvalidWorkspaceFilter(err)
	}
	// A top-level `null` decodes into a nil map without error, and every
	// lookup on a nil map succeeds by returning the zero value, so without
	// this the filter would silently fall through to the default.
	if raw == nil {
		return parsed, ErrInvalidWorkspaceFilter(
			fmt.Errorf("filter must be a JSON object, got null"))
	}

	if value, ok := raw["assigned"]; ok {
		assigned, ok := value.(bool)
		if !ok {
			return parsed, ErrInvalidWorkspaceFilter(
				fmt.Errorf(`"assigned" must be a boolean, got %s`, workspaceFilterJSONType(value)))
		}
		parsed.Assigned = assigned
	}

	if value, ok := raw["deletedAt"]; ok {
		deletedAt, ok := value.(bool)
		if !ok {
			return parsed, ErrInvalidWorkspaceFilter(
				fmt.Errorf(`"deletedAt" must be a boolean, got %s`, workspaceFilterJSONType(value)))
		}
		parsed.DeletedAt = deletedAt
	}

	if value, ok := raw["owner"]; ok {
		owner, ok := value.(string)
		if !ok {
			return parsed, ErrInvalidWorkspaceFilter(
				fmt.Errorf(`"owner" must be a string, got %s`, workspaceFilterJSONType(value)))
		}
		parsed.Owner = owner
	}

	return parsed, nil
}

// workspaceFilterJSONType names the JSON type encoding/json decoded a value
// into, so the error tells the caller what they actually sent rather than a Go
// type name.
func workspaceFilterJSONType(value interface{}) string {
	switch value.(type) {
	case nil:
		return "null"
	case bool:
		return "a boolean"
	case float64:
		return "a number"
	case string:
		return "a string"
	case []interface{}:
		return "an array"
	case map[string]interface{}:
		return "an object"
	default:
		return fmt.Sprintf("%T", value)
	}
}
