package models

import (
	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/spf13/viper"
)

// EventsPersister assists with persisting events in local SQLite DB
type EventsPersister struct {
	DB *database.Handler
}

// swagger:response EventsResponse
type EventsResponse struct {
	Events               []*events.Event         `json:"events"`
	Page                 int                     `json:"page"`
	PageSize             int                     `json:"page_size"`
	CountBySeverityLevel []*CountBySeverityLevel `json:"count_by_severity_level"`
	TotalCount           int64                   `json:"total_count"`
}

type CountBySeverityLevel struct {
	Severity string `json:"severity"`
	Count    int    `json:"count"`
}

func (e *EventsPersister) getCountBySeverity(userID uuid.UUID, eventStatus events.EventStatus) ([]*CountBySeverityLevel, error) {
	if eventStatus == "" {
		eventStatus = events.Unread
	}
	// Get the system ID from the config for the current instance. This is used to filter events that are not associated with the user but are associated with the system
	systemID := viper.GetString("INSTANCE_ID")
	sysID := uuid.FromStringOrNil(systemID)

	eventsBySeverity := []*CountBySeverityLevel{}
	err := e.DB.Model(&events.Event{}).
		Select("severity, count(severity) as count").
		Where("status = ? AND (user_id = ? OR user_id = ?)", eventStatus, userID, sysID).
		Group("severity").
		Find(&eventsBySeverity).Error

	if err != nil {
		return nil, err
	}

	return eventsBySeverity, nil
}
