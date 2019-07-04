package helpers

import (
	"encoding/json"
	"os"
	"path"
	"sync"

	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

type FileSessionPersister struct {
	folder string
	lock   map[string]*sync.Mutex
}

func NewFileSessionPersister(folder string) *FileSessionPersister {
	return &FileSessionPersister{
		folder: folder,
		lock:   map[string]*sync.Mutex{},
	}
}

func (s *FileSessionPersister) openFile(userId string) (*os.File, os.FileInfo, error) {
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
			logrus.Errorf("unable to find/stat the folder '%s': %v", err)
			return nil, nil, err
		}
	}

	fileName := path.Join(s.folder, userId)
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

func (s *FileSessionPersister) Lock(userId string) {
	userLock, ok := s.lock[userId]
	if !ok {
		userLock = &sync.Mutex{}
		s.lock[userId] = userLock
	}
	userLock.Lock()
}

func (s *FileSessionPersister) Unlock(userId string) {
	userLock, ok := s.lock[userId]
	if !ok {
		userLock = &sync.Mutex{}
		s.lock[userId] = userLock
	}
	userLock.Unlock()
}

func (s *FileSessionPersister) Read(userId string) (*models.Session, error) {
	fp, fs, err := s.openFile(userId)
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

func (s *FileSessionPersister) Write(userId string, data *models.Session) error {
	fp, _, err := s.openFile(userId)
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

func (s *FileSessionPersister) Delete(userId string) error {
	fp, _, err := s.openFile(userId)
	if err != nil {
		return err
	}
	fp.Close()
	err = os.Remove(path.Join(s.folder, userId))
	if err != nil {
		logrus.Errorf("unable to delete the file: %v", err)
		return err
	}
	return nil
}
