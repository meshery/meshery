package models

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/meshery/meshkit/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

// ReadFromPersister reads the session data for the given UserID
func (s *MapPreferencePersister) ReadFromPersister(UserID string) (*Preference, error) {
	logLevel := viper.GetInt("LOG_LEVEL")
	if viper.GetBool("DEBUG") {
		logLevel = int(logrus.DebugLevel)
	}
	log, err := logger.New("meshery", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: logLevel,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	data := NewDefaultPreference()

	if s.db == nil {
		return nil, ErrDBConnection
	}

	if UserID == "" {
		return nil, ErrUserID
	}

	dataCopyB, ok := s.db.Load(UserID)
	if ok {
		log.Debug(fmt.Sprintf("retrieved session for user with id: %s", UserID))
		newData, ok1 := dataCopyB.(*Preference)
		if ok1 {
			log.Debug(fmt.Sprintf("session for user with id: %s was read in tact.", UserID))
			data = newData
		} else {
			log.Warn(ErrSessionNotReadIntact(UserID))
		}
	} else {
		log.Debug(ErrSessionNotFound(UserID))
		// Do not store the freshly created default preference on read-miss.
		// Let callers explicitly persist preferences via WriteToPersister.
	}
	return data, nil
}

// WriteToPersister persists session for the user
func (s *MapPreferencePersister) WriteToPersister(UserID string, data *Preference) error {
	if s.db == nil {
		return ErrDBConnection
	}

	if UserID == "" {
		return ErrUserID
	}

	if data == nil {
		return ErrNilConfigData
	}
	data.UpdatedAt = time.Now()
	newSess := NewDefaultPreference()
	if err := copier.Copy(newSess, data); err != nil {
		return ErrSessionCopy(err)
	}

	s.db.Store(UserID, newSess)

	return nil
}

// DeleteFromPersister removes the session for the user
func (s *MapPreferencePersister) DeleteFromPersister(UserID string) error {
	if s.db == nil {
		return ErrDBConnection
	}

	if UserID == "" {
		return ErrUserID
	}
	s.db.Delete(UserID)
	return nil
}

// ClosePersister closes the DB
func (s *MapPreferencePersister) ClosePersister() {
	s.db = nil
}
