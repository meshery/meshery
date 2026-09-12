package utils

import (
	"github.com/gofrs/uuid"
)

func IsUUID(value string) bool {
	if _, err := uuid.FromString(value); err != nil {
		return false
	}
	return true
}
