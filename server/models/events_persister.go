package models

import (
	"fmt"

	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/models/events"
)

// EventsPersister assists with persisting events in local SQLite DB
type EventsPersister struct{
	DB *database.Handler
}

func (e *EventsPersister) PersistEvent(event *events.Event) error {
	err := e.DB.Save(event).Error
	if err != nil {
		return ErrPersistEvent(err)
	}
	return nil
}

func (e *EventsPersister) GetEvent(userID, eventID string) (*events.Event, error) {
	fmt.Println("line[8]: ", userID, eventID)
	return nil, nil
}

func (e *EventsPersister) DeleteEvent(userID, eventID string) error {
	fmt.Println("line[8]: ", userID, eventID)
	return nil
}
