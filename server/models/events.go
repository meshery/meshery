package models

import (
	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/events"
)

const (
	Create = "create"
	Update = "update"
)

type EventTypesResponse struct {
	Category []string `json:"category"`
	Action   []string `json:"action"`
}

type MesheryEvents interface {
	// GetAllEvents(eventFilter *events.EventsFilter, userID uuid.UUID, sysID uuid.UUID) (*EventsResponse, error)
	// GetEventTypes(userID uuid.UUID, sysID uuid.UUID) (map[string]interface{}, error)
	DeleteEvent(eventID uuid.UUID) error
	UpdateEventStatus(eventID uuid.UUID, status string) (*events.Event, error)
	BulkUpdateEventStatus(eventID []*uuid.UUID, status string) ([]*events.Event, error)
	BulkDeleteEvent(eventID []*uuid.UUID) error
}
