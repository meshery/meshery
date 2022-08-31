package models

import (
	"encoding/json"
	"time"

	"github.com/layer5io/meshkit/database"
)

// PreferencePersister assists with persisting session in store
type SessionPreferencePersister struct {
	DB *database.Handler
}
type UserPreference struct {
	ID              string `json:"user_id"`
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

	data := &Preference{
		AnonymousUsageStats:  true,
		AnonymousPerfResults: true,
	}
	var u UserPreference
	err := s.DB.Model(&UserPreference{}).Where("id = ?", userID).First(&u).Error
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(u.PreferenceBytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteToPersister persists session for the user
func (s *SessionPreferencePersister) WriteToPersister(userID string, data *Preference) error {
	var u UserPreference
	err := s.DB.Model(&UserPreference{}).Where("id = ?", userID).First(&u).Error
	if err == nil {
		data.UpdatedAt = time.Now()
		p, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return s.DB.Model(&UserPreference{}).Where("id = ?", userID).Update("preference_bytes", p).Error
	}
	if s.DB == nil {
		return ErrDBConnection
	}

	if userID == "" {
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
		ID:              userID,
	}

	return s.DB.Model(&UserPreference{}).Create(&p).Error
}

// DeleteFromPersister removes the session for the user
func (s *SessionPreferencePersister) DeleteFromPersister(userID string) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	return s.DB.Model(&UserPreference{}).Where("id = ?", userID).Delete(&UserPreference{}).Error
}

// // ClosePersister closes the badger store
// func (s *SessionPreferencePersister) ClosePersister() {
// 	if s.db == nil {
// 		return
// 	}
// 	_ = s.db.Close()
// 	s.cache = nil
// }
