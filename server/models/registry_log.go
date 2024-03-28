package models

import "github.com/google/uuid"

type EntityErrorCount struct {
	Attempt int
	Error   error
}

type EntityCountWithErrors struct {
	Model        map[string]EntityErrorCount
	Registry     map[string]EntityErrorCount
	Component    map[string]EntityErrorCount
	Relationship map[uuid.UUID]EntityErrorCount
	Policy       map[uuid.UUID]EntityErrorCount
}
