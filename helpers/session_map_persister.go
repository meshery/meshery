package helpers

import (
	"sync"

	"github.com/jinzhu/copier"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// MapSessionPersister assists with persisting session in a badger store
type MapSessionPersister struct {
	db *sync.Map
}

// NewMapSessionPersister creates a new MapSessionPersister instance
func NewMapSessionPersister() (*MapSessionPersister, error) {
	return &MapSessionPersister{
		db: &sync.Map{},
	}, nil
}

// Read reads the session data for the given userID
func (s *MapSessionPersister) Read(userID string) (*models.Session, error) {
	data := &models.Session{}

	if s.db == nil {
		return nil, errors.New("Connection to DB does not exist.")
	}

	if userID == "" {
		return nil, errors.New("User ID is empty.")
	}

	dataCopyB, ok := s.db.Load(userID)
	if ok {
		logrus.Debugf("retrieved session for user with id: %s", userID)
		newData, ok1 := dataCopyB.(*models.Session)
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

// Write persists session for the user
func (s *MapSessionPersister) Write(userID string, data *models.Session) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if userID == "" {
		return errors.New("User ID is empty.")
	}

	if data == nil {
		return errors.New("Given config data is nil.")
	}
	newSess := &models.Session{}
	if err := copier.Copy(newSess, data); err != nil {
		logrus.Errorf("session copy error: %v", err)
		return err
	}

	s.db.Store(userID, newSess)

	return nil
}

// Delete removes the session for the user
func (s *MapSessionPersister) Delete(userID string) error {
	if s.db == nil {
		return errors.New("Connection to DB does not exist.")
	}

	if userID == "" {
		return errors.New("User ID is empty.")
	}
	s.db.Delete(userID)
	return nil
}

// Close closes the DB
func (s *MapSessionPersister) Close() {
	s.db = nil
}
