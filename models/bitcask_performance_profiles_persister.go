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

// BitCaskPerformanceProfilesPersister assists with locally persisting performance
// profiles
type BitCaskPerformanceProfilesPersister struct {
	fileName string
	db       *bitcask.Bitcask
}

// PerformanceProfilePage represents a page of performance profiles
type PerformanceProfilePage struct {
	Page       uint64                `json:"page"`
	PageSize   uint64                `json:"page_size"`
	TotalCount int                   `json:"total_count"`
	Profiles   []*PerformanceProfile `json:"profiles"`
}

// NewBitCaskPerformanceProfilesPersister returns a pointer to an instance of
// BitCaskPerformanceProfilePersister
func NewBitCaskPerformanceProfilesPersister(folderName string) (*BitCaskPerformanceProfilesPersister, error) {
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

	fileName := path.Join(folderName, "performanceProfileDB")
	db, err := bitcask.Open(fileName, bitcask.WithSync(true))
	if err != nil {
		logrus.Errorf("Unable to open database: %v.", err)
		return nil, err
	}
	bd := &BitCaskPerformanceProfilesPersister{
		fileName: fileName,
		db:       db,
	}

	return bd, nil
}

// GetPerformanceProfiles takes in a page size and page number and returns a performance profiles page accordingly
func (s *BitCaskPerformanceProfilesPersister) GetPerformanceProfiles(page, pageSize uint64) ([]byte, error) {
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
		if err := s.db.Unlock(); err != nil {
			logrus.Error(errors.Wrap(err, "failed to unlock performance profile database"))
		}
	}()

	total := s.db.Len()

	profiles := []*PerformanceProfile{}

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
				profile := &PerformanceProfile{}
				if err := json.Unmarshal(dd, profile); err != nil {
					err = errors.Wrapf(err, "unable to unmarshal data.")
					logrus.Error(err)
					return nil, err
				}

				profiles = append(profiles, profile)
			}
		}

		localIndex++
	}

	bd, err := json.Marshal(&PerformanceProfilePage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		Profiles:   profiles,
	})
	if err != nil {
		err = errors.Wrapf(err, "unable to marshal result data.")
		logrus.Error(err)
		return nil, err
	}

	return bd, nil
}

// GetPerformanceProfile - gets a performanec profile with specific id
func (s *BitCaskPerformanceProfilesPersister) GetPerformanceProfile(id uuid.UUID) (*PerformanceProfile, error) {
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

	keyb := id.Bytes()
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

	profile := &PerformanceProfile{}
	err = json.Unmarshal(data, profile)
	if err != nil {
		err = errors.Wrapf(err, "unable to marshal profile data.")
		logrus.Error(err)
		return nil, err
	}

	return profile, nil
}

// DeletePeformanceProfile - delete a performance profile for a specific id
func (s *BitCaskPerformanceProfilesPersister) DeletePeformanceProfile(id uuid.UUID) ([]byte, error) {
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

	keyb := id.Bytes()
	if !s.db.Has(keyb) {
		err = errors.New("given key not found")
		logrus.Error(err)
		return nil, err
	}

	json, err := s.db.Get(keyb)
	if err != nil {
		err = errors.Wrapf(err, "unable to fetch performance profile data")
		logrus.Error(err)
		return json, err
	}

	err = s.db.Delete(keyb)
	if err != nil {
		err = errors.Wrapf(err, "unable to delete performance profile data")
		logrus.Error(err)
		return json, err
	}

	return json, nil
}

// SavePerformanceProfile persists the performance profile
func (s *BitCaskPerformanceProfilesPersister) SavePerformanceProfile(id uuid.UUID, profile *PerformanceProfile) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
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

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		err = errors.Wrapf(err, "unable to marshal performance profile data to json")
		logrus.Error(err)
		return err
	}

	if err := s.db.Put(id.Bytes(), profileJSON); err != nil {
		err = errors.Wrapf(err, "unable to persist performance profile data.")
		logrus.Error(err)
		return err
	}

	logrus.Info("Successfully persisted performance profile:", *profile)
	return nil
}

// ClosePerformanceProfilePersister closes the badger store
func (s *BitCaskPerformanceProfilesPersister) ClosePerformanceProfilePersister() {
	if s.db == nil {
		return
	}

	logrus.Debug("Attempting to close db connection...")
	if err := s.db.Close(); err != nil {
		logrus.Error("Failed to close the database connection")
	}
}
