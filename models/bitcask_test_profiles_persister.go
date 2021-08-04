package models

import (
	"encoding/json"
	"os"
	"path"

	"git.mills.io/prologic/bitcask"
	"github.com/gofrs/uuid"
	SMP "github.com/layer5io/service-mesh-performance/spec"
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
				return nil, ErrMakeDir(err, folderName)
			}
		} else {
			return nil, ErrFolderStat(err, folderName)
		}
	}

	fileName := path.Join(folderName, "testConfigDB")
	db, err := bitcask.Open(fileName, bitcask.WithSync(true))
	if err != nil {
		return nil, ErrDBOpen(err)
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
		return nil, ErrDBConnection
	}

	total := s.db.Len()

	testConfigs := []*SMP.PerformanceTestConfig{}

	start := page * pageSize
	end := (page+1)*pageSize - 1
	logrus.Debugf("received page: %d, page size: %d, total: %d", page, pageSize, total)
	logrus.Debugf("computed start index: %d, end index: %d", start, end)

	if start > uint64(total) {
		return nil, ErrIndexOutOfRange
	}
	var localIndex uint64

	for k := range s.db.Keys() {
		if localIndex >= start && localIndex <= end {
			dd, err := s.db.Get(k)
			if err != nil {
				return nil, ErrDBRead(err)
			}
			if len(dd) > 0 {
				testConfig := &SMP.PerformanceTestConfig{}
				if err := json.Unmarshal(dd, testConfig); err != nil {
					return nil, ErrUnmarshal(err, "data")
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
		return nil, ErrMarshal(err, "result data")
	}

	return bd, nil
}

// GetTestConfig - gets result for a specific key
func (s *BitCaskTestProfilesPersister) GetTestConfig(key uuid.UUID) (*SMP.PerformanceTestConfig, error) {
	if s.db == nil {
		return nil, ErrDBConnection
	}

	keyb := key.Bytes()
	if !s.db.Has(keyb) {
		return nil, ErrNilKeys
	}

	data, err := s.db.Get(keyb)
	if err != nil {
		return nil, ErrFetchData(err)
	}

	testConfig := &SMP.PerformanceTestConfig{}
	err = json.Unmarshal(data, testConfig)
	if err != nil {
		return nil, ErrMarshal(err, "testConfig data")
	}

	return testConfig, nil
}

// DeleteTestConfig - delete result for a specific key
func (s *BitCaskTestProfilesPersister) DeleteTestConfig(key uuid.UUID) error {
	if s.db == nil {
		return ErrDBConnection
	}

	keyb := key.Bytes()
	if !s.db.Has(keyb) {
		return ErrNilKeys
	}

	err := s.db.Delete(keyb)
	if err != nil {
		return ErrFetchData(err)
	}

	return nil
}

// WriteTestConfig persists the result
func (s *BitCaskTestProfilesPersister) WriteTestConfig(key uuid.UUID, result []byte) error {
	if s.db == nil {
		return ErrDBConnection
	}

	if result == nil {
		return ErrResultData()
	}

	if err := s.db.Put(key.Bytes(), result); err != nil {
		return ErrUnableToPersistsResult(err)
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
