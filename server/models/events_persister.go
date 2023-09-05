package models

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/models/events"
	"gorm.io/gorm/clause"
)

// EventsPersister assists with persisting events in local SQLite DB
type EventsPersister struct{
	DB *database.Handler
}

func (e *EventsPersister) GetAllEvents(search, status string, eventsFilter *events.EventsFilter, userID uuid.UUID) ([]*events.Event, int64, error) {
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

	if search != "" {
		finder = finder.Where("description LIKE ?", "%"+search+"%")
	}

	if status != "" {
		finder = finder.Where("status = ?", status)
	}

	if eventsFilter.Order == "desc" {
		finder = finder.Order(clause.OrderByColumn{Column: clause.Column{Name: eventsFilter.SortOn}, Desc: true})
	} else {
		finder = finder.Order(eventsFilter.SortOn)
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
		return nil, count, err
	}
	return eventsDB, count, nil
}	

func (e *EventsPersister) UpdateEventStatus(eventID uuid.UUID, status string) (*events.Event, error) {
	err := e.DB.Model(&events.Event{ID: eventID}).Update("status", status).Error
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

func (e *EventsPersister) DeleteEvent(eventID uuid.UUID) error {
	err := e.DB.Delete(&events.Event{ID: eventID}).Error
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