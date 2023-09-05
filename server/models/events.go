package models

import "github.com/layer5io/meshkit/models/events"

type MesheryEvents interface {
	GetEvent(userID, eventID string) (*events.Event, error)
	PersistEvent(data *events.Event) error
	DeleteEvent(userID, eventID string) error
}