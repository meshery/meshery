package models

import (
	"encoding/json"
	"strings"

	"github.com/layer5io/meshkit/database"
	"gorm.io/gorm"
)

// MesheryK8sContextPersister is the persister for persisting
// applications on the database
type MesheryK8sContextPersister struct {
	DB *database.Handler
}

// MesheryK8sContextPage represents a page of contexts
type MesheryK8sContextPage struct {
	Page       uint64        `json:"page"`
	PageSize   uint64        `json:"page_size"`
	TotalCount int           `json:"total_count"`
	Contexts   []*K8sContext `json:"contexts"`
}

// GetMesheryK8sContexts returns all of the contexts
func (mkcp *MesheryK8sContextPersister) GetMesheryK8sContexts(search, order string, page, pageSize uint64) ([]byte, error) {
	order = sanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = "updated_at desc"
	}

	count := int64(0)
	contexts := []*K8sContext{}

	query := mkcp.DB.Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(name) like ?)", like)
	}

	query.Model(K8sContext{}).Count(&count)

	Paginate(uint(page), uint(pageSize))(query).Find(&contexts)

	mesheryK8sContextPage := MesheryK8sContextPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Contexts:   contexts,
	}

	resp, _ := json.Marshal(mesheryK8sContextPage)
	return resp, nil
}

// DeleteMesheryK8sContext takes in an application id and delete it if it already exists
func (mkcp *MesheryK8sContextPersister) DeleteMesheryK8sContext(id string) (K8sContext, error) {
	context := K8sContext{ID: id}
	mkcp.DB.Delete(&context)

	return context, nil
}

func (mkcp *MesheryK8sContextPersister) SaveMesheryK8sContext(mkc K8sContext) (K8sContext, error) {
	if mkc.ID == "" {
		id, err := K8sContextGenerateID(mkc)
		if err != nil {
			return mkc, ErrContextID
		}

		mkc.ID = id
	}

	// Perform the operation in a transaction
	err := mkcp.DB.Transaction(func(tx *gorm.DB) error {
		var mesheryK8sContext K8sContext

		// Check if there is already an entry for this context
		if err := tx.First(&mesheryK8sContext, "id = ?", mkc.ID).Error; err == nil {
			return ErrContextAlreadyPersisted
		}

		return tx.Save(&mkc).Error
	})

	return mkc, err
}

func (mkcp *MesheryK8sContextPersister) GetMesheryK8sContext(id string) (K8sContext, error) {
	var mesheryK8sContext K8sContext

	err := mkcp.DB.First(&mesheryK8sContext, "id = ?", id).Error
	return mesheryK8sContext, err
}

// func (mkcp *MesheryK8sContextPersister) SetMesheryK8sCurrentContext(id string) error {
// 	// Perform the operation in a transaction
// 	return mkcp.DB.Transaction(func(tx *gorm.DB) error {
// 		var mesheryK8sContext K8sContext

// 		// Get context which is currently in use
// 		if err := tx.First(&mesheryK8sContext, "is_current_context = true").Error; err != nil {
// 			return err
// 		}

// 		// If the context id matches with the provided id then skip the next steps
// 		if mesheryK8sContext.ID == id {
// 			return nil
// 		}

// 		if err := tx.Save(&mesheryK8sContext).Error; err != nil {
// 			return err
// 		}

// 		// Set the specified context as active
// 		return tx.Model(K8sContext{}).Where("id = ?", id).Update("is_current_context", true).Error
// 	})
// }

// func (mkcp *MesheryK8sContextPersister) GetMesheryK8sCurrentContext() (K8sContext, error) {
// 	var mesheryK8sContext K8sContext

// 	err := mkcp.DB.First(&mesheryK8sContext, "is_current_context = true").Error

// 	return mesheryK8sContext, err
// }
