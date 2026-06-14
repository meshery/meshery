package models

import (
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
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
	// GetAllEvents(eventFilter *events.EventsFilter, userID core.Uuid, sysID core.Uuid) (*EventsResponse, error)
	// GetEventTypes(userID core.Uuid, sysID core.Uuid) (map[string]interface{}, error)
	DeleteEvent(eventID core.Uuid) error
	UpdateEventStatus(eventID core.Uuid, status string) (*events.Event, error)
	BulkUpdateEventStatus(eventID []*core.Uuid, status string) ([]*events.Event, error)
	BulkDeleteEvent(eventID []*core.Uuid) error
}
