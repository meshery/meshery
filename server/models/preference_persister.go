package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/meshery/meshkit/database"
	"gorm.io/gorm"
)

// PreferencePersister assists with persisting session in store
type SessionPreferencePersister struct {
	DB *database.Handler
}
type UserPreference struct {
	ID              string `json:"UserID"`
	PreferenceBytes []byte `json:"preference"`
}

// ReadFromPersister - reads the session data for the given UserID
func (s *SessionPreferencePersister) ReadFromPersister(UserID string) (*Preference, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}

	if UserID == "" {
		return nil, ErrUserID
	}

	data := NewDefaultPreference()
	var u UserPreference
	err := s.DB.Model(&UserPreference{}).Where("id = ?", UserID).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, nil
		}
		return nil, err
	}
	err = json.Unmarshal(u.PreferenceBytes, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteToPersister persists session for the user
func (s *SessionPreferencePersister) WriteToPersister(UserID string, data *Preference) error {
	var u UserPreference
	err := s.DB.Model(&UserPreference{}).Where("id = ?", UserID).First(&u).Error
	if err == nil {
		data.UpdatedAt = time.Now()
		p, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return s.DB.Model(&UserPreference{}).Where("id = ?", UserID).Update("preference_bytes", p).Error
	}
	if s.DB == nil {
		return ErrDBConnection
	}

	if UserID == "" {
		return ErrUserID
	}

	if data == nil {
		return ErrNilConfigData
	}

	dataB, err := json.Marshal(data)
	if err != nil {
		return ErrMarshal(err, "User Config Data")
	}

	var p = UserPreference{
		PreferenceBytes: dataB,
		ID:              UserID,
	}

	return s.DB.Model(&UserPreference{}).Create(&p).Error
}

// DeleteFromPersister removes the session for the user
func (s *SessionPreferencePersister) DeleteFromPersister(UserID string) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if UserID == "" {
		return ErrUserID
	}

	return s.DB.Model(&UserPreference{}).Where("id = ?", UserID).Delete(&UserPreference{}).Error
}

// // ClosePersister closes the badger store
// func (s *SessionPreferencePersister) ClosePersister() {
// 	if s.db == nil {
// 		return
// 	}
// 	_ = s.db.Close()
// 	s.cache = nil
// }
