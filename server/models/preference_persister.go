package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/meshery/meshkit/database"
	"gorm.io/gorm"
)

// SessionPreferencePersister assists with persisting session in store.
type SessionPreferencePersister struct {
	DB *database.Handler
}
type UserPreference struct {
	ID              string `json:"userId"`
	PreferenceBytes []byte `json:"preference"`
}

// ReadFromPersister - reads the session data for the given userID
func (s *SessionPreferencePersister) ReadFromPersister(userID string) (*Preference, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}

	if userID == "" {
		return nil, ErrUserID
	}

	data := NewDefaultPreference()
	var u UserPreference
	err := s.DB.Model(&UserPreference{}).Where("id = ?", userID).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, nil
		}
		return nil, ErrDBRead(err)
	}
	err = json.Unmarshal(u.PreferenceBytes, data)
	if err != nil {
		return nil, ErrUnmarshal(err, "User Preference Data")
	}
	return data, nil
}

// WriteToPersister persists session for the user
func (s *SessionPreferencePersister) WriteToPersister(userID string, data *Preference) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	if data == nil {
		return ErrNilConfigData
	}

	var u UserPreference
	err := s.DB.Model(&UserPreference{}).Where("id = ?", userID).First(&u).Error
	if err == nil {
		data.UpdatedAt = time.Now()
		p, err := json.Marshal(data)
		if err != nil {
			return ErrMarshal(err, "User Config Data")
		}
		if err := s.DB.Model(&UserPreference{}).Where("id = ?", userID).Update("preference_bytes", p).Error; err != nil {
			return ErrDBPut(err)
		}
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDBRead(err)
	}

	dataB, err := json.Marshal(data)
	if err != nil {
		return ErrMarshal(err, "User Config Data")
	}

	var p = UserPreference{
		PreferenceBytes: dataB,
		ID:              userID,
	}

	if err := s.DB.Model(&UserPreference{}).Create(&p).Error; err != nil {
		return ErrDBCreate(err)
	}
	return nil
}

// DeleteFromPersister removes the session for the user
func (s *SessionPreferencePersister) DeleteFromPersister(userID string) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	res := s.DB.Model(&UserPreference{}).Where("id = ?", userID).Delete(&UserPreference{})
	if res.Error != nil {
		return ErrDBDelete(res.Error, userID)
	}
	if res.RowsAffected == 0 {
		return ErrResultNotFound(fmt.Errorf("user preference for user %s not found", userID))
	}
	return nil
}

// // ClosePersister closes the badger store
// func (s *SessionPreferencePersister) ClosePersister() {
// 	if s.db == nil {
// 		return
// 	}
// 	_ = s.db.Close()
// 	s.cache = nil
// }
