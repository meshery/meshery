package models

import (
	"github.com/layer5io/meshkit/database"
)

// PreferencePersister assists with persisting session in store
type CapabilitiesPersister struct {
	DB *database.Handler
}
type UserCapabilities struct {
	ID                string       `json:"id"`
	CapabilitiesBytes Capabilities `json:"capabilities" gorm:"type:bytes;serializer:json" db:"capabilities"`
}

// ReadFromPersister - reads the session data for the given userID
func (s *CapabilitiesPersister) ReadFromPersister(userID string) (*Capabilities, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}

	if userID == "" {
		return nil, ErrUserID
	}

	capabilities := &Capabilities{}
	err := s.DB.Model(&UserCapabilities{}).Where("id = ?", userID).First(capabilities).Error
	if err != nil {
		return nil, err
	}

	return capabilities, nil
}

// WriteToPersister persists session for the user
func (s *CapabilitiesPersister) WriteToPersister(userID string, data *Capabilities) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	err := s.DB.Model(&UserCapabilities{}).Save(data).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteFromPersister removes the session for the user
func (s *CapabilitiesPersister) DeleteFromPersister(userID string) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	return s.DB.Delete(&UserCapabilities{ID: userID}).Error
}