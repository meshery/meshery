package models

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/models/events"
)

type MesheryEvents interface {
	GetAllEvents(eventFilter *events.EventsFilter, userID uuid.UUID) (*EventsResponse, error)
	GetEventTypes(userID uuid.UUID) (map[string]interface{}, error)
	PersistEvent(data *events.Event) error
	DeleteEvent(eventID uuid.UUID) error
	UpdateEventStatus(eventID uuid.UUID, status string) (*events.Event, error)
	BulkUpdateEventStatus(eventID []*uuid.UUID, status string) ([]*events.Event, error)
	BulkDeleteEvent(eventID []*uuid.UUID) error
}
