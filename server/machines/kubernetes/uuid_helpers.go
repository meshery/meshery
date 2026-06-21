package kubernetes

import "github.com/google/uuid"

// parseUUIDOrNil parses s as a UUID and returns uuid.Nil if parsing fails.
// This is a compatibility shim replacing gofrs/uuid's FromStringOrNil.
func parseUUIDOrNil(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}
