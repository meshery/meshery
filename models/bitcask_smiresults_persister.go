package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/prologic/bitcask"
	"github.com/sirupsen/logrus"
)

// BitCaskSmiResultsPersister assists with persisting session in a Bitcask store
type BitCaskSmiResultsPersister struct {
	fileName string
	db       *bitcask.Bitcask
}

// SmiResultPage - represents a page of meshery results
type SmiResultPage struct {
	Page       uint64       `json:"page"`
	PageSize   uint64       `json:"page_size"`
	TotalCount int          `json:"total_count"`
	Results    []*SmiResult `json:"results"`
}

// NewBitCaskSmiResultsPersister creates a new BitCaskSmiResultsPersister instance
func NewBitCaskSmiResultsPersister(folderName string) (*BitCaskSmiResultsPersister, error) {
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

	fileName := path.Join(folderName, "smiresultDB")
	db, err := bitcask.Open(fileName, bitcask.WithSync(true))
	if err != nil {
		logrus.Errorf("Unable to open database: %v.", err)
		return nil, err
	}
	bd := &BitCaskSmiResultsPersister{
		fileName: fileName,
		db:       db,
	}
	return bd, nil
}

// GetSmiResults - gets result for the page and pageSize
func (s *BitCaskSmiResultsPersister) GetResults(page, pageSize uint64) ([]byte, error) {
	if s.db == nil {
		return nil, errors.New("connection to DB does not exist")
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

	total := s.db.Len()

	results := []*SmiResult{}

	start := page * pageSize
	end := (page+1)*pageSize - 1
	logrus.Debugf("received page: %d, page size: %d, total: %d", page, pageSize, total)
	logrus.Debugf("computed start index: %d, end index: %d", start, end)

	if start > uint64(total) {
		return nil, fmt.Errorf("index out of range")
	}
	var localIndex uint64

	for k := range s.db.Keys() {
		if localIndex >= start && localIndex <= end {
			dd, err := s.db.Get(k)
			if err != nil {
				err = errors.Wrapf(err, "Unable to read data from bitcask store")
				logrus.Error(err)
				return nil, err
			}
			if len(dd) > 0 {
				result := &SmiResult{}
				if err := json.Unmarshal(dd, result); err != nil {
					err = errors.Wrapf(err, "Unable to unmarshal data.")
					logrus.Error(err)
					return nil, err
				}
				results = append(results, result)
			}
		}
		localIndex++
	}

	bd, err := json.Marshal(&SmiResultPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		Results:    results,
	})
	if err != nil {
		err = errors.Wrapf(err, "Unable to marshal result data.")
		logrus.Error(err)
		return nil, err
	}

	return bd, nil
}

// WriteSmiResult persists the result
func (s *BitCaskSmiResultsPersister) WriteResult(key uuid.UUID, result []byte) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if result == nil {
		return errors.New("given result data is nil")
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

	if err := s.db.Put(key.Bytes(), result); err != nil {
		err = errors.Wrapf(err, "Unable to persist result data.")
		logrus.Error(err)
		return err
	}
	return nil
}

// CloseSmiResultPersister closes the badger store
func (s *BitCaskSmiResultsPersister) CloseResultPersister() {
	if s.db == nil {
		return
	}
	_ = s.db.Close()
}
