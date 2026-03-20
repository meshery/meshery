package models

import "github.com/gofrs/uuid"

func validateUUID(id, field string) (uuid.UUID, error) {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return uuid.Nil, ErrProviderInvalidUUID(err, field)
	}

	return parsedID, nil
}
