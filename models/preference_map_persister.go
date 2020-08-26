package models

import (
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// MapPreferencePersister assists with persisting session in a badger store
type MapPreferencePersister struct {
	db *sync.Map
}

// NewMapPreferencePersister creates a new MapPreferencePersister instance
func NewMapPreferencePersister() (*MapPreferencePersister, error) {
	return &MapPreferencePersister{
		db: &sync.Map{},
	}, nil
}

// ReadFromPersister reads the session data for the given userID
func (s *MapPreferencePersister) ReadFromPersister(userID string) (*Preference, error) {
	data := &Preference{
		AnonymousUsageStats:  true,
		AnonymousPerfResults: true,
	}

	if s.db == nil {
		return nil, errors.New("connection to DB does not exist")
	}

	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	dataCopyB, ok := s.db.Load(userID)
	if ok {
		logrus.Debugf("retrieved session for user with id: %s", userID)
		newData, ok1 := dataCopyB.(*Preference)
		if ok1 {
			logrus.Debugf("session for user with id: %s was read in tact.", userID)
			data = newData
		} else {
			logrus.Warnf("session for user with id: %s was NOT read in tact.", userID)
		}
	} else {
		logrus.Warnf("unable to find session for user with id: %s.", userID)
	}
	return data, nil
}

// WriteToPersister persists session for the user
func (s *MapPreferencePersister) WriteToPersister(userID string, data *Preference) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if userID == "" {
		return errors.New("user ID is empty")
	}

	if data == nil {
		return errors.New("given config data is nil")
	}
	data.UpdatedAt = time.Now()
	newSess := &Preference{
		AnonymousUsageStats:  true,
		AnonymousPerfResults: true,
	}
	if err := copier.Copy(newSess, data); err != nil {
		logrus.Errorf("session copy error: %v", err)
		return err
	}

	s.db.Store(userID, newSess)

	return nil
}

// DeleteFromPersister removes the session for the user
func (s *MapPreferencePersister) DeleteFromPersister(userID string) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if userID == "" {
		return errors.New("user ID is empty")
	}
	s.db.Delete(userID)
	return nil
}

// ClosePersister closes the DB
func (s *MapPreferencePersister) ClosePersister() {
	s.db = nil
}
