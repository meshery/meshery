package models

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/layer5io/meshkit/logger"
)

// MapPreferencePersister assists with persisting session in a badger store
type MapPreferencePersister struct {
	db     *sync.Map
	logger logger.Handler
}

// NewMapPreferencePersister creates a new MapPreferencePersister instance
func NewMapPreferencePersister(log logger.Handler) (*MapPreferencePersister, error) {
	return &MapPreferencePersister{
		db:     &sync.Map{},
		logger: log,
	}, nil
}

// @TODO revise to use logging
// ReadFromPersister reads the session data for the given userID
func (s *MapPreferencePersister) ReadFromPersister(userID string) (*Preference, error) {
	data := &Preference{
		AnonymousUsageStats:  true,
		AnonymousPerfResults: true,
	}

	if s.db == nil {
		return data, nil
	}

	if userID == "" {
		return nil, fmt.Errorf("user ID is empty")
	}

	dataCopyB, ok := s.db.Load(userID)
	if !ok {
		return nil, fmt.Errorf("unable to find session for user with id: %s", userID)
	}

	s.logger.Debugf("retrieved session for user with id: %s", userID)

	newData, ok1 := dataCopyB.(*Preference)
	if !ok1 {
		s.logger.Warnf("session for user with id: %s was NOT intact", userID)
		return data, fmt.Errorf("unable to cast session for user with id: %s", userID)
	}

	s.logger.Debugf("session for user with id: %s was read intact", userID)
	return newData, nil
}

// @TODO revise to add more error handling
// WriteToPersister persists session for the user
func (s *MapPreferencePersister) WriteToPersister(userID string, data *Preference) error {
	if s.db == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	if data == nil {
		return ErrNilConfigData
	}

	dataCopy := &Preference{}
	if err := copier.Copy(dataCopy, data); err != nil {
		return ErrSessionCopy(err)
	}
	dataCopy.UpdatedAt = time.Now()

	if _, ok := s.db.LoadOrStore(userID, dataCopy); ok {
		s.logger.Debugf("session for user with id: %s was updated", userID)
	} else {
		s.logger.Debugf("session for user with id: %s was created", userID)
	}

	return nil
}

// DeleteFromPersister removes the session for the user
func (s *MapPreferencePersister) DeleteFromPersister(userID string) error {
	if s.db == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}
	s.db.Delete(userID)
	return nil
}

// ClosePersister closes the DB
func (s *MapPreferencePersister) ClosePersister() {
	s.db = nil
}
