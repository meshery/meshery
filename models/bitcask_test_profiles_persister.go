package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/gofrs/uuid"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/pkg/errors"
	"github.com/prologic/bitcask"
	"github.com/sirupsen/logrus"
)

// BitCaskTestProfilesPersister assists with persisting session in a Bitcask store
type BitCaskTestProfilesPersister struct {
	fileName string
	db       *bitcask.Bitcask
}

// UserTestProfiles - represents a page of user test configs
type UserTestProfiles struct {
	Page        uint64                       `json:"page"`
	PageSize    uint64                       `json:"page_size"`
	TotalCount  int                          `json:"total_count"`
	TestConfigs []*SMP.PerformanceTestConfig `json:"test_configs"`
}

// NewBitCaskTestProfilesPersister creates a new BitCaskTestProfilesPersister instance
func NewBitCaskTestProfilesPersister(folderName string) (*BitCaskTestProfilesPersister, error) {
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

	fileName := path.Join(folderName, "testConfigDB")
	db, err := bitcask.Open(fileName, bitcask.WithSync(true))
	if err != nil {
		logrus.Errorf("Unable to open database: %v.", err)
		return nil, err
	}
	bd := &BitCaskTestProfilesPersister{
		fileName: fileName,
		db:       db,
	}
	return bd, nil
}

// GetTestConfigs - gets result for the page and pageSize
func (s *BitCaskTestProfilesPersister) GetTestConfigs(page, pageSize uint64) ([]byte, error) {
	if s.db == nil {
		return nil, errors.New("connection to DB does not exist")
	}

RETRY:
	locked, err := s.db.TryRLock()
	if err != nil {
		err = errors.Wrapf(err, "unable to obtain read lock from bitcask store")
		logrus.Error(err)
	}
	if !locked {
		goto RETRY
	}
	defer func() {
		_ = s.db.Unlock()
	}()

	total := s.db.Len()

	testConfigs := []*SMP.PerformanceTestConfig{}

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
				err = errors.Wrapf(err, "unable to read data from bitcask store")
				logrus.Error(err)
				return nil, err
			}
			if len(dd) > 0 {
				testConfig := &SMP.PerformanceTestConfig{}
				if err := json.Unmarshal(dd, testConfig); err != nil {
					err = errors.Wrapf(err, "unable to unmarshal data.")
					logrus.Error(err)
					return nil, err
				}
				testConfigs = append(testConfigs, testConfig)
			}
		}
		localIndex++
	}

	bd, err := json.Marshal(&UserTestProfiles{
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  total,
		TestConfigs: testConfigs,
	})
	if err != nil {
		err = errors.Wrapf(err, "unable to marshal result data.")
		logrus.Error(err)
		return nil, err
	}

	return bd, nil
}

// GetTestConfig - gets result for a specific key
func (s *BitCaskTestProfilesPersister) GetTestConfig(key uuid.UUID) (*SMP.PerformanceTestConfig, error) {
	if s.db == nil {
		return nil, errors.New("connection to DB does not exist")
	}

RETRY:
	locked, err := s.db.TryRLock()
	if err != nil {
		err = errors.Wrapf(err, "unable to obtain read lock from bitcask store")
		logrus.Error(err)
	}
	if !locked {
		goto RETRY
	}
	defer func() {
		_ = s.db.Unlock()
	}()

	keyb := key.Bytes()
	if !s.db.Has(keyb) {
		err = errors.New("given key not found")
		logrus.Error(err)
		return nil, err
	}

	data, err := s.db.Get(keyb)
	if err != nil {
		err = errors.Wrapf(err, "unable to fetch result data")
		logrus.Error(err)
		return nil, err
	}

	testConfig := &SMP.PerformanceTestConfig{}
	err = json.Unmarshal(data, testConfig)
	if err != nil {
		err = errors.Wrapf(err, "unable to marshal testConfig data.")
		logrus.Error(err)
		return nil, err
	}

	return testConfig, nil
}

// DeleteTestConfig - delete result for a specific key
func (s *BitCaskTestProfilesPersister) DeleteTestConfig(key uuid.UUID) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

RETRY:
	locked, err := s.db.TryRLock()
	if err != nil {
		err = errors.Wrapf(err, "unable to obtain read lock from bitcask store")
		logrus.Error(err)
	}
	if !locked {
		goto RETRY
	}
	defer func() {
		_ = s.db.Unlock()
	}()

	keyb := key.Bytes()
	if !s.db.Has(keyb) {
		err = errors.New("given key not found")
		logrus.Error(err)
		return err
	}

	err = s.db.Delete(keyb)
	if err != nil {
		err = errors.Wrapf(err, "unable to fetch result data")
		logrus.Error(err)
		return err
	}

	return nil
}

// WriteTestConfig persists the result
func (s *BitCaskTestProfilesPersister) WriteTestConfig(key uuid.UUID, result []byte) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if result == nil {
		return errors.New("given result data is nil")
	}

RETRY:
	locked, err := s.db.TryLock()
	if err != nil {
		err = errors.Wrapf(err, "unable to obtain write lock from bitcask store")
		logrus.Error(err)
	}
	if !locked {
		goto RETRY
	}
	defer func() {
		_ = s.db.Unlock()
	}()

	if err := s.db.Put(key.Bytes(), result); err != nil {
		err = errors.Wrapf(err, "unable to persist result data.")
		logrus.Error(err)
		return err
	}
	return nil
}

// CloseTestConfigsPersister closes the badger store
func (s *BitCaskTestProfilesPersister) CloseTestConfigsPersister() {
	if s.db == nil {
		return
	}
	_ = s.db.Close()
}
