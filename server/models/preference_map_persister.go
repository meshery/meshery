package models

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/layer5io/meshkit/logger"
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

// ReadFromPersister reads the session data for the given userID
func (s *MapPreferencePersister) ReadFromPersister(userID string) (*Preference, error) {
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
	data := &Preference{
		AnonymousUsageStats:  true,
		AnonymousPerfResults: true,
	}

	if s.db == nil {
		return nil, ErrDBConnection
	}

	if userID == "" {
		return nil, ErrUserID
	}

	dataCopyB, ok := s.db.Load(userID)
	if ok {
		log.Debug(fmt.Sprintf("retrieved session for user with id: %s", userID))
		newData, ok1 := dataCopyB.(*Preference)
		if ok1 {
			log.Debug(fmt.Sprintf("session for user with id: %s was read in tact.", userID))
			data = newData
		} else {
			log.Warn(ErrSessionNotReadIntact(userID))
		}
	} else {
		log.Warn(ErrSessionNotFound(userID))
	}
	return data, nil
}

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
	data.UpdatedAt = time.Now()
	newSess := &Preference{
		AnonymousUsageStats:  true,
		AnonymousPerfResults: true,
	}
	if err := copier.Copy(newSess, data); err != nil {
		return ErrSessionCopy(err)
	}

	s.db.Store(userID, newSess)

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
