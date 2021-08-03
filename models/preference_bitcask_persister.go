package models

import (
	"encoding/json"
	"os"
	"path"
	"sync"
	"time"

	"git.mills.io/prologic/bitcask"
	"github.com/jinzhu/copier"
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
				return nil, ErrMakeDir(err, folderName)
			}
		} else {
			return nil, ErrFolderStat(err, folderName)
		}
	}

	fileName := path.Join(folderName, "db")
	db, err := bitcask.Open(fileName, bitcask.WithSync(true), bitcask.WithMaxValueSize(uint64(1<<32)))
	if err != nil {
		return nil, ErrDBOpen(err)
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
		return nil, ErrDBConnection
	}

	if userID == "" {
		return nil, ErrUserID
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

	dataCopyB, err := s.db.Get([]byte(userID))
	if err != nil {
		return nil, ErrDBRead(err)
	}
	if len(dataCopyB) > 0 {
		if err := json.Unmarshal(dataCopyB, data); err != nil {
			return nil, ErrUnmarshal(err, "data")
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
		return ErrCopy(err, "session")
	}
	s.cache.Store(userID, newSess)
	return nil
}

// WriteToPersister persists session for the user
func (s *BitCaskPreferencePersister) WriteToPersister(userID string, data *Preference) error {
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

	if err := s.writeToCache(userID, data); err != nil {
		return err
	}

	dataB, err := json.Marshal(data)
	if err != nil {
		return ErrMarshal(err, "User Config Data")
	}

	if err := s.db.Put([]byte(userID), dataB); err != nil {
		return ErrDBPut(err)
	}
	return nil
}

// DeleteFromPersister removes the session for the user
func (s *BitCaskPreferencePersister) DeleteFromPersister(userID string) error {
	if s.db == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	s.cache.Delete(userID)
	if err := s.db.Delete([]byte(userID)); err != nil {
		return ErrDBDelete(err, userID)
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
