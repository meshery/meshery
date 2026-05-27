// Package utils provides CLI commands and utilities for mesheryctl.
package utils

import (
	"github.com/google/uuid"
)

func IsUUID(value string) bool {
	if _, err := uuid.Parse(value); err != nil {
		return false
	}
	return true
}
