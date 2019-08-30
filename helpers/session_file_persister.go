package helpers

import (
	"encoding/json"
	"os"
	"path"
	"sync"

	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

// FileSessionPersister assists with persisting session in a file store
type FileSessionPersister struct {
	folder string
	mutex  map[string]*sync.Mutex
}

// NewFileSessionPersister returns a new instance of FileSessionPersister
func NewFileSessionPersister(folder string) *FileSessionPersister {
	return &FileSessionPersister{
		folder: folder,
		mutex:  map[string]*sync.Mutex{},
	}
}

func (s *FileSessionPersister) openFile(userID string) (*os.File, os.FileInfo, error) {
	var fp *os.File
	_, err := os.Stat(s.folder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(s.folder, os.ModePerm)
			if err != nil {
				logrus.Errorf("unable to create the directory '%s' due to error: %v ", s.folder, err)
				return nil, nil, err
			}
		} else {
			logrus.Errorf("unable to find/stat the folder '%s': %v", s.folder, err)
			return nil, nil, err
		}
	}

	fileName := path.Join(s.folder, userID)
	_, err = os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			fp, err = os.Create(fileName)
		} else {
			logrus.Errorf("unable to find/stat the file '%s': %v", fileName, err)
			return nil, nil, err
		}
	} else {
		fp, err = os.OpenFile(fileName, os.O_RDWR, 0644)
	}
	if err != nil {
		logrus.Errorf("error opening file '%s': %v", fileName, err)
		return nil, nil, err
	}

	fs, err := os.Stat(fileName)
	if err != nil {
		logrus.Errorf("unable to stat the file '%s': %v", fileName, err)
		return nil, nil, err
	}
	return fp, fs, nil
}

func (s *FileSessionPersister) lock(userID string) {
	userLock, ok := s.mutex[userID]
	if !ok {
		userLock = &sync.Mutex{}
		s.mutex[userID] = userLock
	}
	userLock.Lock()
}

func (s *FileSessionPersister) unlock(userID string) {
	userLock, ok := s.mutex[userID]
	if !ok {
		userLock = &sync.Mutex{}
		s.mutex[userID] = userLock
	}
	userLock.Unlock()
}

// Read reads the session data for the given userID
func (s *FileSessionPersister) Read(userID string) (*models.Session, error) {
	s.lock(userID)
	defer s.unlock(userID)

	fp, fs, err := s.openFile(userID)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	data := &models.Session{}
	if fs.Size() > 0 {
		err = json.NewDecoder(fp).Decode(data)
		if err != nil {
			logrus.Errorf("error decoding contents from file: %v", err)
			return nil, err
		}
	}
	return data, nil
}

// Write persists session for the user
func (s *FileSessionPersister) Write(userID string, data *models.Session) error {
	s.lock(userID)
	defer s.unlock(userID)

	fp, _, err := s.openFile(userID)
	if err != nil {
		return err
	}
	defer fp.Close()
	err = json.NewEncoder(fp).Encode(data)
	if err != nil {
		logrus.Errorf("error encoding contents to file: %v", err)
		return err
	}
	return nil
}

// Delete removes the session for the user
func (s *FileSessionPersister) Delete(userID string) error {
	s.lock(userID)
	defer s.unlock(userID)

	fp, _, err := s.openFile(userID)
	if err != nil {
		return err
	}
	fp.Close()
	err = os.Remove(path.Join(s.folder, userID))
	if err != nil {
		logrus.Errorf("unable to delete the file: %v", err)
		return err
	}
	return nil
}

// Close closes the persister
func (s *FileSessionPersister) Close() {
}
