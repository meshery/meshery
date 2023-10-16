package models

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/models/events"
	"gorm.io/gorm/clause"
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

func (e *EventsPersister) GetEventTypes(userID uuid.UUID) (map[string]interface{}, error) {
	eventTypes := make(map[string]interface{}, 2)
	var categories, actions []string
	err := e.DB.Table("events").Distinct("category").Where("user_id = ?", userID).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	eventTypes["category"] = categories
	err = e.DB.Table("events").Distinct("action").Where("user_id = ?", userID).Find(&actions).Error
	if err != nil {
		return nil, err
	}

	eventTypes["action"] = actions
	return eventTypes, err
}

func (e *EventsPersister) GetAllEvents(eventsFilter *events.EventsFilter, userID uuid.UUID) (*EventsResponse, error) {
	eventsDB := []*events.Event{}
	finder := e.DB.Model(&events.Event{}).Where("user_id = ?", userID)

	if len(eventsFilter.Category) != 0 {
		finder = finder.Where("category IN ?", eventsFilter.Category)
	}

	if len(eventsFilter.Action) != 0 {
		finder = finder.Where("action IN ?", eventsFilter.Action)
	}

	if len(eventsFilter.Severity) != 0 {
		finder = finder.Where("severity IN ?", eventsFilter.Severity)
	}

	if eventsFilter.Search != "" {
		finder = finder.Where("description LIKE ?", "%"+eventsFilter.Search+"%")
	}

	if eventsFilter.Status != "" {
		finder = finder.Where("status = ?", eventsFilter.Status)
	}

	if eventsFilter.Order == "asc" {
		finder = finder.Order(eventsFilter.SortOn)
	} else {
		finder = finder.Order(clause.OrderByColumn{Column: clause.Column{Name: eventsFilter.SortOn}, Desc: true})
	}

	var count int64
	finder.Count(&count)

	if eventsFilter.Offset != 0 {
		finder = finder.Offset(eventsFilter.Offset)
	}

	if eventsFilter.Limit != 0 {
		finder = finder.Limit(eventsFilter.Limit)
	}

	err := finder.Scan(&eventsDB).Error
	if err != nil {
		return nil, err
	}

	countBySeverity, err := e.getCountBySeverity(userID, eventsFilter.Status)

	if err != nil {
		return nil, err
	}

	return &EventsResponse{
		Events:               eventsDB,
		PageSize:             eventsFilter.Limit,
		TotalCount:           count,
		CountBySeverityLevel: countBySeverity,
	}, nil
}

func (e *EventsPersister) UpdateEventStatus(eventID uuid.UUID, status string) (*events.Event, error) {
	err := e.DB.Model(&events.Event{ID: eventID, Status: events.EventStatus(status)}).Update("status", status).Error
	if err != nil {
		return nil, err
	}

	updatedEvent := &events.Event{}
	err = e.DB.Find(updatedEvent, "id = ?", eventID).Error
	if err != nil {
		return nil, err
	}
	return updatedEvent, nil
}

func (e *EventsPersister) BulkUpdateEventStatus(eventIDs []*uuid.UUID, status string) ([]*events.Event, error) {

	err := e.DB.Model(&events.Event{Status: events.EventStatus(status)}).Where("id IN ?", eventIDs).Update("status", status).Error
	if err != nil {
		return nil, err
	}
	
	updatedEvent := &[]*events.Event{}
	err = e.DB.Find(updatedEvent, "id IN ?", eventIDs).Error
	if err != nil {
		return nil, err
	}

	return *updatedEvent, nil
}

func (e *EventsPersister) DeleteEvent(eventID uuid.UUID) error {
	err := e.DB.Delete(&events.Event{ID: eventID}).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *EventsPersister) BulkDeleteEvent(eventIDs []*uuid.UUID) error {
	err := e.DB.Where("id IN ?", eventIDs).Delete(&events.Event{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *EventsPersister) PersistEvent(event *events.Event) error {
	err := e.DB.Save(event).Error
	if err != nil {
		return ErrPersistEvent(err)
	}
	return nil
}

func (e *EventsPersister) getCountBySeverity(userID uuid.UUID, eventStatus events.EventStatus) ([]*CountBySeverityLevel, error) {
	if eventStatus == "" {
		eventStatus = events.Unread
	}

	eventsBySeverity := []*CountBySeverityLevel{}
	err := e.DB.Model(&events.Event{}).Select("severity, count(severity) as count").Where("status = ? and user_id = ?", eventStatus, userID).Group("severity").Find(&eventsBySeverity).Error
	if err != nil {
		return nil, err
	}

	return eventsBySeverity, nil
}
