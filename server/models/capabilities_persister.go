package models

import (
	"github.com/meshery/meshkit/database"
)

// UserCapabilitiesPersister assists with persisting user capabilities in the store
type UserCapabilitiesPersister struct {
	DB *database.Handler
}
type UserCapabilities struct {
	ID           string             `json:"id"`
	Capabilities ProviderProperties `json:"capabilities" gorm:"type:bytes;serializer:json" db:"capabilities"`
}

// ReadFromPersister - reads the capabilities for the given UserID
func (u *UserCapabilitiesPersister) ReadCapabilitiesForUser(UserID string) (*ProviderProperties, error) {
	if u.DB == nil {
		return nil, ErrDBConnection
	}

	if UserID == "" {
		return nil, ErrUserID
	}

	capabilities := &UserCapabilities{}
	err := u.DB.Model(capabilities).Where("id = ?", UserID).First(capabilities).Error
	if err != nil {
		return nil, err
	}

	return &capabilities.Capabilities, nil
}

// WriteToPersister persists the capabilities for the user
func (u *UserCapabilitiesPersister) WriteCapabilitiesForUser(UserID string, data *ProviderProperties) error {
	if u.DB == nil {
		return ErrDBConnection
	}

	if UserID == "" {
		return ErrUserID
	}

	userCapabilities := &UserCapabilities{
		ID:           UserID,
		Capabilities: *data,
	}

	err := u.DB.Model(&UserCapabilities{}).Where("id = ?", UserID).Save(userCapabilities).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteFromPersister removes the capabilities for the user
func (u *UserCapabilitiesPersister) DeleteCapabilitiesForUser(UserID string) error {
	if u.DB == nil {
		return ErrDBConnection
	}

	if UserID == "" {
		return ErrUserID
	}

	return u.DB.Delete(&UserCapabilities{ID: UserID}).Error
}
