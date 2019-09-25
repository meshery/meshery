package helpers

import (
	"encoding/json"
	"os"
	"path"

	"github.com/dgraph-io/badger"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// BadgerSessionPersister assists with persisting session in a badger store
type BadgerSessionPersister struct {
	fileName string
	db       *badger.DB
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

	fileName := path.Join(folderName, "db")
	db, err := badger.Open(badger.DefaultOptions(fileName))
	if err != nil {
		logrus.Errorf("Unable to open database: %v.", err)
		return nil, err
	}

	return &BadgerSessionPersister{
		fileName: fileName,
		db:       db,
	}, nil
}

// Read reads the session data for the given userID
func (s *BadgerSessionPersister) Read(userID string) (*models.Session, error) {
	data := &models.Session{}
	dataCopyB := []byte{}

	if s.db == nil {
		return nil, errors.New("Connection to DB does not exist") //error strings should not be capitalized or end with punctuation or a newline
	}

	if userID == "" {
		return nil, errors.New("User ID is empty") //error strings should not be capitalized or end with punctuation or a newline
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
	return data, nil
}

// Write persists session for the user
func (s *BadgerSessionPersister) Write(userID string, data *models.Session) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if userID == "" {
		return errors.New("User ID is empty") //error strings should not be capitalized or end with punctuation or a newline
	}

	if data == nil {
		return errors.New("Given config data is nil") //error strings should not be capitalized or end with punctuation or a newline
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
		return errors.New("Connection to DB does not exist") //error strings should not be capitalized or end with punctuation or a newline
	}

	if userID == "" {
		return errors.New("User ID is empty") //error strings should not be capitalized or end with punctuation or a newline
	}

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
	s.db.Close()
}
