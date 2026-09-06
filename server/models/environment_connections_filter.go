package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

// EnvironmentConnectionsFilter is the decoded `filter` query parameter of
// GET /api/environments/{environmentID}/connections.
//
// The wire form is a JSON object - the grammar Meshery UI has always sent
// (`{"assigned":false}`). It is deliberately NOT the space-separated
// "key value" form that helpers/utils.ApplyFilters reads: this endpoint used to
// hand the same string to both parsers, so whichever grammar a caller picked,
// the other one rejected or silently ignored it. Decoding once, here, is what
// keeps the two from fighting. See issue #21826.
//
// Only the fields this endpoint can actually serve are accepted. `owner` maps
// to the `connections.owner` column. There is deliberately no
// `organization_id`: the previous ApplyFilters call listed one, but the
// `connections` table has no such column, so it could only ever have produced a
// SQL error had it been reachable.
type EnvironmentConnectionsFilter struct {
	// Assigned selects the connections assigned to the environment (true, the
	// default when the caller says nothing) or the unassigned complement.
	Assigned bool

	// Owner, when non-empty, narrows the result to connections with that owner.
	Owner string
}

// ParseEnvironmentConnectionsFilter decodes the `filter` query parameter.
//
// An empty filter yields the default: assigned connections, unnarrowed. Any
// other input must be a JSON object whose recognised fields carry the right
// type; anything else is a caller error and is reported as one rather than
// panicking (`assigned` was previously read with an unchecked type assertion,
// so a non-boolean took down the request goroutine) or surfacing a bare
// encoding/json message.
//
// Unrecognised fields are ignored, so a caller sending a field this endpoint
// does not serve gets the default behaviour rather than a rejection.
func ParseEnvironmentConnectionsFilter(filter string) (EnvironmentConnectionsFilter, error) {
	parsed := EnvironmentConnectionsFilter{Assigned: true}

	if strings.TrimSpace(filter) == "" {
		return parsed, nil
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(filter), &raw); err != nil {
		return parsed, ErrInvalidEnvironmentConnectionsFilter(err)
	}

	if value, ok := raw["assigned"]; ok {
		assigned, ok := value.(bool)
		if !ok {
			return parsed, ErrInvalidEnvironmentConnectionsFilter(
				fmt.Errorf(`"assigned" must be a boolean, got %s`, jsonTypeName(value)))
		}
		parsed.Assigned = assigned
	}

	if value, ok := raw["owner"]; ok {
		owner, ok := value.(string)
		if !ok {
			return parsed, ErrInvalidEnvironmentConnectionsFilter(
				fmt.Errorf(`"owner" must be a string, got %s`, jsonTypeName(value)))
		}
		parsed.Owner = owner
	}

	return parsed, nil
}

// jsonTypeName names the JSON type encoding/json decoded a value into, so the
// error tells the caller what they actually sent rather than a Go type name.
func jsonTypeName(value interface{}) string {
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
