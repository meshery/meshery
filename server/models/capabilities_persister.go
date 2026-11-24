package models

import (
	"errors"
	"fmt"

	"github.com/meshery/meshkit/database"
	"gorm.io/gorm"
)

// UserCapabilitiesPersister assists with persisting user capabilities in the store
type UserCapabilitiesPersister struct {
	DB *database.Handler
}
type UserCapabilities struct {
	ID           string             `json:"id"`
	Capabilities ProviderProperties `json:"capabilities" gorm:"type:bytes;serializer:json" db:"capabilities"`
}

// ReadCapabilitiesForUser - reads the capabilities for the given userID
func (u *UserCapabilitiesPersister) ReadCapabilitiesForUser(userID string) (*ProviderProperties, error) {
	if u == nil || u.DB == nil {
		return nil, ErrDBConnection
	}

	if userID == "" {
		return nil, ErrUserID
	}

	capabilities := &UserCapabilities{}
	err := u.DB.Model(capabilities).Where("id = ?", userID).First(capabilities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrResultNotFound(err)
		}
		return nil, err
	}

	return &capabilities.Capabilities, nil
}

// WriteCapabilitiesForUser persists the capabilities for the user
func (u *UserCapabilitiesPersister) WriteCapabilitiesForUser(userID string, data *ProviderProperties) error {
	if u == nil || u.DB == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	if data == nil {
		return ErrNilConfigData
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

// DeleteCapabilitiesForUser removes the capabilities for the user
func (u *UserCapabilitiesPersister) DeleteCapabilitiesForUser(userID string) error {
	if u == nil || u.DB == nil {
		return ErrDBConnection
	}

	if userID == "" {
		return ErrUserID
	}

	res := u.DB.Delete(&UserCapabilities{ID: userID})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrResultNotFound(fmt.Errorf("capabilities for user %s not found", userID))
	}
	return nil
}
