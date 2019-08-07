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

type BadgerSessionPersister struct {
	fileName string
	db       *badger.DB
}

func NewBadgerSessionPersister(folderName string) (*BadgerSessionPersister, error) {
	_, err := os.Stat(folderName)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(folderName, os.ModePerm)
			if err != nil {
				logrus.Errorf("unable to create the directory '%s' due to error: %v ", folderName, err)
				return nil, err
			}
		} else {
			logrus.Errorf("unable to find/stat the folder '%s': %v", err)
			return nil, err
		}
	}

	fileName := path.Join(folderName, "db")
	db, err := badger.Open(badger.DefaultOptions(fileName))
	if err != nil {
		logrus.Errorf("unable to open database: %v", err)
		return nil, err
	}

	return &BadgerSessionPersister{
		fileName: fileName,
		db:       db,
	}, nil
}

func (s *BadgerSessionPersister) Read(userId string) (*models.Session, error) {
	data := &models.Session{}
	dataCopyB := []byte{}

	if s.db == nil {
		return nil, errors.New("connection to DB does not exist")
	}

	if userId == "" {
		return nil, errors.New("user id is empty")
	}

	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userId))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			} else {
				err = errors.Wrapf(err, "unable to retrieve data for user: %s", userId)
				logrus.Error(err)
				return err
			}
		}
		dataCopyB, err = item.ValueCopy(nil)
		if err != nil {
			err = errors.Wrapf(err, "unable to copy data")
			logrus.Error(err)
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if len(dataCopyB) > 0 {
		if err := json.Unmarshal(dataCopyB, data); err != nil {
			err = errors.Wrapf(err, "unable to unmarshal data")
			logrus.Error(err)
			return nil, err
		}
	}
	return data, nil
}

func (s *BadgerSessionPersister) Write(userId string, data *models.Session) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if userId == "" {
		return errors.New("user id is empty")
	}

	if data == nil {
		return errors.New("given config data is nil")
	}

	dataB, err := json.Marshal(data)
	if err != nil {
		err = errors.Wrapf(err, "unable to marshal the user config data")
		logrus.Error(err)
		return err
	}
	return s.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(userId), dataB); err != nil {
			err = errors.Wrapf(err, "unable to persist config data")
			return err
		}
		return nil
	})
}

func (s *BadgerSessionPersister) Delete(userId string) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if userId == "" {
		return errors.New("user id is empty")
	}

	return s.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete([]byte(userId)); err != nil {
			err = errors.Wrapf(err, "unable to delete config data for the user: %s", userId)
			return err
		}
		return nil
	})
}

func (s *BadgerSessionPersister) Close() {
	if s.db == nil {
		return
	}
	s.db.Close()
}
