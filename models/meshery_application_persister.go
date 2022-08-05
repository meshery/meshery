package models

import (
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
)

// MesheryApplicationPersister is the persister for persisting
// applications on the database
type MesheryApplicationPersister struct {
	DB *database.Handler
}

// MesheryApplicationPage represents a page of applications
type MesheryApplicationPage struct {
	Page         uint64                `json:"page"`
	PageSize     uint64                `json:"page_size"`
	TotalCount   int                   `json:"total_count"`
	Applications []*MesheryApplication `json:"applications"`
}

// GetMesheryApplications returns all of the applications
func (maap *MesheryApplicationPersister) GetMesheryApplications(search, order string, page, pageSize uint64) ([]byte, error) {
	order = sanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = "updated_at desc"
	}

	count := int64(0)
	applications := []*MesheryApplication{}

	query := maap.DB.Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(meshery_applications.name) like ?)", like)
	}

	query.Table("meshery_applications").Count(&count)

	Paginate(uint(page), uint(pageSize))(query).Find(&applications)

	mesheryApplicationPage := &MesheryApplicationPage{
		Page:         page,
		PageSize:     pageSize,
		TotalCount:   int(count),
		Applications: applications,
	}

	return marshalMesheryApplicationPage(mesheryApplicationPage), nil
}

// DeleteMesheryApplication takes in an application id and delete it if it already exists
func (maap *MesheryApplicationPersister) DeleteMesheryApplication(id uuid.UUID) ([]byte, error) {
	application := MesheryApplication{ID: &id}
	err := maap.DB.Delete(&application).Error

	return marshalMesheryApplication(&application), err
}

func (maap *MesheryApplicationPersister) SaveMesheryApplication(application *MesheryApplication) ([]byte, error) {
	if application.ID == nil {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}

		application.ID = &id
	}

	return marshalMesheryApplications([]MesheryApplication{*application}), maap.DB.Save(application).Error
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

	return marshalMesheryApplications(finalApplications), maap.DB.Create(finalApplications).Error
}

func (maap *MesheryApplicationPersister) GetMesheryApplication(id uuid.UUID) ([]byte, error) {
	var mesheryApplication MesheryApplication
	err := maap.DB.First(&mesheryApplication, id).Error
	return marshalMesheryApplication(&mesheryApplication), err
}

func (maap *MesheryApplicationPersister) GetMesheryApplicationSource(id uuid.UUID) ([]byte, error) {
	var mesheryApplication MesheryApplication
	err := maap.DB.First(&mesheryApplication, id).Error
	return mesheryApplication.SourceContent, err
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
