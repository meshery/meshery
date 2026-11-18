package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	"gorm.io/gorm"
)

// MesheryApplicationPersister is the persister for persisting
// applications on the database
type MesheryApplicationPersister struct {
	DB *database.Handler
}

// MesheryApplicationPage represents a page of applications
type MesheryApplicationPage struct {
	Page         uint64                `json:"page"`
	PageSize     uint64                `json:"pageSize"`
	TotalCount   int                   `json:"totalCount"`
	Applications []*MesheryApplication `json:"applications"`
}

// GetMesheryApplications returns all of the applications
func (maap *MesheryApplicationPersister) GetMesheryApplications(search, order string, page, pageSize uint64, updatedAfter string) ([]byte, error) {
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	count := int64(0)
	applications := []*MesheryApplication{}

	query := maap.DB.Where("updated_at > ?", updatedAfter).Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(meshery_applications.name) like ?)", like)
	}

	if err := query.Table("meshery_applications").Count(&count).Error; err != nil {
		return nil, ErrDBRead(err)
	}

	if err := Paginate(uint(page), uint(pageSize))(query).Find(&applications).Error; err != nil {
		return nil, ErrDBRead(err)
	}

	mesheryApplicationPage := &MesheryApplicationPage{
		Page:         page,
		PageSize:     pageSize,
		TotalCount:   int(count),
		Applications: applications,
	}

	return marshalMesheryApplicationPage(mesheryApplicationPage), nil
}

// DeleteMesheryApplication takes in an application id and delete it if it already exists
func (maap *MesheryApplicationPersister) DeleteMesheryApplication(id core.Uuid) ([]byte, error) {
	application := MesheryApplication{ID: &id}
	res := maap.DB.Delete(&application)
	if res.Error != nil {
		return nil, ErrDBDelete(res.Error, id.String())
	}
	if res.RowsAffected == 0 {
		return nil, ErrResultNotFound(fmt.Errorf("application %s not found", id.String()))
	}

	return marshalMesheryApplication(&application), nil
}

// SaveMesheryApplications batch inserts the given applications
func (maap *MesheryApplicationPersister) SaveMesheryApplications(applications []MesheryApplication) ([]byte, error) {
	finalApplications := []MesheryApplication{}
	for _, application := range applications {
		if application.ID == nil {
			id, err := uuid.NewV4()
			if err != nil {
				return nil, ErrGenerateUUID(err)
			}

			application.ID = &id
		}

		finalApplications = append(finalApplications, application)
	}

	if err := maap.DB.Create(finalApplications).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	return marshalMesheryApplications(finalApplications), nil
}

func (maap *MesheryApplicationPersister) GetMesheryApplication(id core.Uuid) ([]byte, error) {
	var mesheryApplication MesheryApplication
	err := maap.DB.First(&mesheryApplication, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrResultNotFound(fmt.Errorf("application %s not found", id.String()))
		}
		return nil, ErrDBRead(err)
	}
	return marshalMesheryApplication(&mesheryApplication), nil
}

func (maap *MesheryApplicationPersister) GetMesheryApplicationSource(id core.Uuid) ([]byte, error) {
	var mesheryApplication MesheryApplication
	err := maap.DB.First(&mesheryApplication, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrResultNotFound(fmt.Errorf("application %s not found", id.String()))
		}
		return nil, ErrDBRead(err)
	}
	return mesheryApplication.SourceContent, nil
}

func marshalMesheryApplicationPage(maap *MesheryApplicationPage) []byte {
	res, _ := json.Marshal(maap)

	return res
}

func marshalMesheryApplication(ma *MesheryApplication) []byte {
	res, _ := json.Marshal(ma)

	return res
}

func marshalMesheryApplications(mas []MesheryApplication) []byte {
	res, _ := json.Marshal(mas)

	return res
}
