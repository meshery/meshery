package models

import (
	"encoding/json"
	"os"
	"path"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/prologic/bitcask"
	"github.com/sirupsen/logrus"
)

// BitCaskPreferencePersister assists with persisting session in a Bitcask store
type BitCaskPreferencePersister struct {
	fileName string
	db       *bitcask.Bitcask
	cache    *sync.Map
}

// NewBitCaskPreferencePersister creates a new BitCaskPreferencePersister instance
func NewBitCaskPreferencePersister(folderName string) (*BitCaskPreferencePersister, error) {
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
	db, err := bitcask.Open(fileName, bitcask.WithSync(true))
	if err != nil {
		logrus.Errorf("Unable to open database: %v.", err)
		return nil, err
	}
	bd := &BitCaskPreferencePersister{
		fileName: fileName,
		db:       db,
		cache:    &sync.Map{},
	}
	return bd, nil
}

// ReadFromPersister - reads the session data for the given userID
func (s *BitCaskPreferencePersister) ReadFromPersister(userID string) (*Preference, error) {
	if s.db == nil {
		return nil, errors.New("connection to DB does not exist")
	}

	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	data := &Preference{
		AnonymousUsageStats:  true,
		AnonymousPerfResults: true,
	}

	dataCopyI, ok := s.cache.Load(userID)
	if ok {
		newData, ok1 := dataCopyI.(*Preference)
		if ok1 {
			return newData, nil
		}
	}

RETRY:
	locked, err := s.db.TryRLock()
	if err != nil {
		err = errors.Wrapf(err, "Unable to obtain read lock from bitcask store")
		logrus.Error(err)
	}
	if !locked {
		goto RETRY
	}
	defer func() {
		_ = s.db.Unlock()
	}()

	dataCopyB, err := s.db.Get([]byte(userID))
	if err != nil {
		err = errors.Wrapf(err, "Unable to read data from bitcask store")
		logrus.Error(err)
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

// writeToCache persists session for the user in the cache
func (s *BitCaskPreferencePersister) writeToCache(userID string, data *Preference) error {
	newSess := &Preference{
		AnonymousUsageStats:  true,
		AnonymousPerfResults: true,
	}
	if err := copier.Copy(newSess, data); err != nil {
		logrus.Errorf("session copy error: %v", err)
		return err
	}
	s.cache.Store(userID, newSess)
	return nil
}

// WriteToPersister persists session for the user
func (s *BitCaskPreferencePersister) WriteToPersister(userID string, data *Preference) error {
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

RETRY:
	locked, err := s.db.TryLock()
	if err != nil {
		err = errors.Wrapf(err, "Unable to obtain write lock from bitcask store")
		logrus.Error(err)
	}
	if !locked {
		goto RETRY
	}
	defer func() {
		_ = s.db.Unlock()
	}()

	if err := s.writeToCache(userID, data); err != nil {
		return err
	}

	dataB, err := json.Marshal(data)
	if err != nil {
		err = errors.Wrapf(err, "Unable to marshal the user config data.")
		logrus.Error(err)
		return err
	}

	if err := s.db.Put([]byte(userID), dataB); err != nil {
		err = errors.Wrapf(err, "Unable to persist config data.")
		return err
	}
	return nil
}

// DeleteFromPersister removes the session for the user
func (s *BitCaskPreferencePersister) DeleteFromPersister(userID string) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if userID == "" {
		return errors.New("user ID is empty")
	}

RETRY:
	locked, err := s.db.TryLock()
	if err != nil {
		err = errors.Wrapf(err, "Unable to obtain write lock from bitcask store")
		logrus.Error(err)
	}
	if !locked {
		goto RETRY
	}
	defer func() {
		_ = s.db.Unlock()
	}()

	s.cache.Delete(userID)
	if err := s.db.Delete([]byte(userID)); err != nil {
		err = errors.Wrapf(err, "Unable to delete config data for the user: %s.", userID)
		return err
	}
	return nil
}

// ClosePersister closes the badger store
func (s *BitCaskPreferencePersister) ClosePersister() {
	if s.db == nil {
		return
	}
	_ = s.db.Close()
	s.cache = nil
}
