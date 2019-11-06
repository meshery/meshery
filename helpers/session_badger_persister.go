package helpers

import (
	"encoding/json"
	"os"
	"path"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
	"github.com/jinzhu/copier"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// BadgerSessionPersister assists with persisting session in a badger store
type BadgerSessionPersister struct {
	fileName string
	db       *badger.DB
	cache    *sync.Map
	ticker   *time.Ticker
}

// NewBadgerSessionPersister creates a new BadgerSessionPersister instance
func NewBadgerSessionPersister(folderName string) (*BadgerSessionPersister, error) {
	_, err := os.Stat(folderName)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(folderName, os.ModePerm)
			if err != nil {
				logrus.Errorf("Unable to create the directory '%s' due to error: %v.", folderName, err)
				return nil, err
			}
		} else {
			logrus.Errorf("Unable to find/stat the folder '%s': %v,", folderName, err)
			return nil, err
		}
	}

	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	fileName := path.Join(folderName, "db")
	db, err := badger.Open(badger.DefaultOptions(fileName).WithValueLogLoadingMode(options.FileIO).
		WithMaxTableSize(1 << 20).
		WithLevelOneSize(1 << 20).
		WithValueLogFileSize(1 << 20).
		WithValueLogMaxEntries(1 << 10).WithLogger(logger))
	if err != nil {
		logrus.Errorf("Unable to open database: %v.", err)
		return nil, err
	}
	bd := &BadgerSessionPersister{
		fileName: fileName,
		db:       db,
		cache:    &sync.Map{},
		ticker:   time.NewTicker(5 * time.Minute),
	}
	go bd.cleanup()
	return bd, nil
}

// Read reads the session data for the given userID
func (s *BadgerSessionPersister) Read(userID string) (*models.Session, error) {
	data := &models.Session{}
	dataCopyB := []byte{}

	if s.db == nil {
		return nil, errors.New("Connection to DB does not exist.")
	}

	if userID == "" {
		return nil, errors.New("User ID is empty.")
	}

	dataCopyI, ok := s.cache.Load(userID)
	if ok {
		newData, ok1 := dataCopyI.(*models.Session)
		if ok1 {
			return newData, nil
		}
	}

	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userID))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			err = errors.Wrapf(err, "Unable to retrieve data for user: %s.", userID)
			logrus.Error(err)
			return err
		}
		dataCopyB, err = item.ValueCopy(nil)
		if err != nil {
			err = errors.Wrapf(err, "Unable to copy data.")
			logrus.Error(err)
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if len(dataCopyB) > 0 {
		if err := json.Unmarshal(dataCopyB, data); err != nil {
			err = errors.Wrapf(err, "Unable to unmarshal data.")
			logrus.Error(err)
			return nil, err
		}
	}

	_ = s.writeToCache(userID, data)
	return data, nil
}

// Write persists session for the user in the cache
func (s *BadgerSessionPersister) writeToCache(userID string, data *models.Session) error {
	newSess := &models.Session{}
	if err := copier.Copy(newSess, data); err != nil {
		logrus.Errorf("session copy error: %v", err)
		return err
	}
	s.cache.Store(userID, newSess)
	return nil
}

// Write persists session for the user
func (s *BadgerSessionPersister) Write(userID string, data *models.Session) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if userID == "" {
		return errors.New("User ID is empty.")
	}

	if data == nil {
		return errors.New("Given config data is nil.")
	}

	if err := s.writeToCache(userID, data); err != nil {
		return err
	}

	dataB, err := json.Marshal(data)
	if err != nil {
		err = errors.Wrapf(err, "Unable to marshal the user config data.")
		logrus.Error(err)
		return err
	}
	return s.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(userID), dataB); err != nil {
			err = errors.Wrapf(err, "Unable to persist config data.")
			return err
		}
		return nil
	})
}

// Delete removes the session for the user
func (s *BadgerSessionPersister) Delete(userID string) error {
	if s.db == nil {
		return errors.New("Connection to DB does not exist.")
	}

	if userID == "" {
		return errors.New("User ID is empty.")
	}

	s.cache.Delete(userID)
	return s.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete([]byte(userID)); err != nil {
			err = errors.Wrapf(err, "Unable to delete config data for the user: %s.", userID)
			return err
		}
		return nil
	})
}

// Close closes the badger store
func (s *BadgerSessionPersister) Close() {
	if s.db == nil {
		return
	}
	if s.ticker != nil {
		s.ticker.Stop()
	}
	_ = s.db.Close()
	s.cache = nil
}

// Close closes the badger store
func (s *BadgerSessionPersister) cleanup() {
	for range s.ticker.C {
		logrus.Debug("running db gc. . .")
		if err := s.db.RunValueLogGC(0.7); err != nil {
			logrus.Debugf("error while running gc: %v", err)
		}
	}
}
