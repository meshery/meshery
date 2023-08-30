package models

import "github.com/layer5io/meshkit/models/events"


type Event = events.Event
type EventType string
type EventSeverity string

type MesheryEvents interface {
	GetEvent(userID, eventID string) (*Event, error)
	CreateEvent(userID string, data *Event) error
	DeleteEvent(userID, eventID string) error
}