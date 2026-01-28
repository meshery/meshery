package utils

import "github.com/google/uuid"

func IsUUID(value string) bool {
	if _, err := uuid.Parse(value); err != nil {
		return false
	}
	return true
}
