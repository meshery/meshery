package models

import (
	"fmt"

	"github.com/layer5io/meshkit/database"
)

// EventsPersister assists with persisting events in local SQLite DB
type EventsPersister struct{
	DB *database.Handler
}

func (e *EventsPersister) CreateEvent(userID string, event*Event) error {
	fmt.Println("line[8]: ", event)
	return nil
}

func (e *EventsPersister) GetEvent(userID, eventID string) (*Event, error) {
	fmt.Println("line[8]: ", userID, eventID)
	return nil, nil
}

func (e *EventsPersister) DeleteEvent(userID, eventID string) error {
	fmt.Println("line[8]: ", userID, eventID)
	return nil
}
