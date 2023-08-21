package models

import (
	"time"
	"strings"

	"github.com/layer5io/meshkit/database"
	"github.com/gofrs/uuid"
)

// MesheryK8sContextPersister is the persister for persisting
// applications on the database
type CredentialPersister struct {
	DB *database.Handler
}

func (c *CredentialPersister) SaveUserCredentials(credential *Credential) (*Credential, error) {
	if err := c.DB.Create(credential); err != nil {
		return nil, err.Error
	}

	return credential, nil
}

func (c *CredentialPersister) GetUserCredentials(page, pageSize int, search, order string) (*CredentialsPage, error) {
	if order == "" {
		order = "created_at desc"
	}


	qc := c.DB.Where("deleted_at is NULL")

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		qc = qc.Where("(lower(name) like ?)", like)
	}

	qc = qc.Order(order)

	var totalCount int64
	qc.Model(&Credential{}).Count(&totalCount)

	credentials := []*Credential{}
	Paginate(uint(page), uint(pageSize))(qc).Find(&credentials)
	// if count > 0 {
	// 	if err := qc.Paginate(page+1, pageSize).All(&credentials); err != nil {
	// 		if err.Error() != sql.ErrNoRows.Error() {
	// 			logrus.Errorf("error retrieving credentials for user id: %s - %v", userID, err)
	// 			return nil, err
	// 		}
	// 	}
	// }
	// logrus.Debugf("retrieved credentials: %+v", credentials)

	return &CredentialsPage{
		Credentials: credentials,
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  int(totalCount),
	}, nil
}

func (c *CredentialPersister) DeleteUserCredential(credentialID uuid.UUID) (*Credential, error) {

	if err := c.DB.Raw("UPDATE credentials SET deleted_at = ? WHERE id = ?", time.Now(), credentialID); err != nil {
		return nil, err.Error
	}

	delCredential := &Credential{}
	if err := c.DB.Raw("SELECT * FROM credentials WHERE id = ?", credentialID).First(delCredential); err != nil {
		return nil, err.Error
	}

	return delCredential, nil
}


func (c *CredentialPersister) UpdateUserCredential(credential *Credential) (*Credential, error) {
	
	if err := c.DB.Updates(credential); err != nil {
		return nil, err.Error
	}
	return credential, nil
}

func (c *CredentialPersister) GetCredentialByID(credentialID uuid.UUID) (*Credential, error) {
	credential := &Credential{}
	if err := c.DB.Where("id = ?", credentialID).First(credential); err != nil {
		return nil, err.Error
	}
	return credential, nil
}