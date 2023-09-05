package models

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/models/events"
)

type MesheryEvents interface {
	GetAllEvents(search, status string, eventFilter *events.EventsFilter, userID uuid.UUID) ([]*events.Event, int64, error)
	PersistEvent(data *events.Event) error
	DeleteEvent(eventID uuid.UUID) error
	UpdateEventStatus(eventID uuid.UUID, status string) (*events.Event, error)
}