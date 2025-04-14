package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
	"gorm.io/gorm"
)

type CredentialPersister struct {
	DB *database.Handler
}

func (crp *CredentialPersister) SaveCredential(credential *Credential) (*Credential, error) {
	err := crp.DB.Transaction(func(tx *gorm.DB) error {
		if result := tx.Table("credentials").Create(&credential); result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error saving user credentials: %v", err)
	}
	return nil, nil
}

func (crp *CredentialPersister) GetCredential(userID string, page, pageSize int, search, order string) (*CredentialsPage, error) {
	result := crp.DB.Select("*").Where("user_id=? and deleted_at is NULL", userID)
	if result.Error != nil {
		return nil, result.Error
	}
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		result = result.Where("(lower(name) like ?)", like)
	}

	result = result.Order(order)

	var count int64
	if err := result.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("error retrieving count of credentials for user id: %s - %v", userID, err)
	}

	var credentialsList []*Credential
	if count > 0 {
		if err := result.Offset(page * pageSize).Limit(pageSize).Find(&credentialsList).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("error retrieving credentials for user id: %s - %v", userID, err)
			}
		}
	}

	credentialsPage := &CredentialsPage{
		Credentials: credentialsList,
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  int(count),
	}

	if result.Error != nil {
		return nil, fmt.Errorf("error getting user credentials: %v", result.Error)
	}

	return credentialsPage, nil
}

func (crp *CredentialPersister) UpdateCredential(credential *Credential) (*Credential, error) {
	updatedCredential := &Credential{}
	if err := crp.DB.Model(*updatedCredential).Where("user_id = ? AND id = ? AND deleted_at is NULL", credential.UserID, credential.ID).Updates(credential); err != nil {
		return nil, fmt.Errorf("error updating user credential: %v", err)
	}

	if err := crp.DB.Where("user_id = ? AND id = ?", credential.UserID, credential.ID).First(updatedCredential).Error; err != nil {
		return nil, fmt.Errorf("error getting updated user credential: %v", err)
	}
	return updatedCredential, nil
}

func (crp *CredentialPersister) DeleteCredential(credentialID uuid.UUID) (*Credential, error) {
	delCredential := &Credential{}
	if err := crp.DB.Model(&Credential{}).Where("id = ?", credentialID).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, err
	}
	if err := crp.DB.Where("id = ?", credentialID).First(delCredential).Error; err != nil {
		return nil, err
	}
	return delCredential, nil
}
