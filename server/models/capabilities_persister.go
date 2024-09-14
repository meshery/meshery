package models

import (
	"github.com/layer5io/meshkit/database"
)

// UserCapabilitiesPersister assists with persisting user capabilities in the store
type UserCapabilitiesPersister struct {
	DB *database.Handler
}
type UserCapabilities struct {
	ID           string             `json:"id"`
	Capabilities ProviderProperties `json:"capabilities" gorm:"type:bytes;serializer:json" db:"capabilities"`
}

// ReadFromPersister - reads the capabilities for the given userID
func (u *UserCapabilitiesPersister) ReadCapabilitiesForUser(userID string) (*ProviderProperties, error) {
	if u.DB == nil {
		return nil, ErrDBConnection
	}

	if userID == "" {
		return nil, ErrUserID
	}

	capabilities := &UserCapabilities{}
	err := u.DB.Model(capabilities).Where("id = ?", userID).First(capabilities).Error
	if err != nil {
		return nil, err
	}

	return &capabilities.Capabilities, nil
}

// WriteToPersister persists the capabilities for the user
func (u *UserCapabilitiesPersister) WriteCapabilitiesForUser(userID string, data *ProviderProperties) error {
	if u.DB == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	userCapabilities := &UserCapabilities{
		ID:           userID,
		Capabilities: *data,
	}

	err := u.DB.Model(&UserCapabilities{}).Where("id = ?", userID).Save(userCapabilities).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteFromPersister removes the capabilities for the user
func (u *UserCapabilitiesPersister) DeleteCapabilitiesForUser(userID string) error {
	if u.DB == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	return u.DB.Delete(&UserCapabilities{ID: userID}).Error
}
