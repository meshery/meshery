package models

import (
	"encoding/json"
	"strings"

	"github.com/meshery/meshery/server/pkg/encryption"
	"github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"gorm.io/gorm"
)

// MesheryK8sContextPersister is the persister for persisting
// applications on the database
type MesheryK8sContextPersister struct {
	DB     *database.Handler
	EncSvc *encryption.Service
}

// MesheryK8sContextPage represents a page of contexts
type MesheryK8sContextPage struct {
	Page       uint64        `json:"page"`
	PageSize   uint64        `json:"pageSize"`
	TotalCount int           `json:"totalCount"`
	Contexts   []*K8sContext `json:"contexts"`
}

// GetMesheryK8sContexts returns all of the contexts
func (mkcp *MesheryK8sContextPersister) GetMesheryK8sContexts(search, order string, page, pageSize uint64) ([]byte, error) {
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = defaultOrderUpdatedAtDesc
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

	if mkcp.EncSvc != nil {
		for _, ctx := range contexts {
			if ctx == nil {
				continue
			}
			if ctx.Auth != nil {
				decryptedAuth, err := mkcp.EncSvc.DecryptMap(ctx.Auth)
				if err != nil {
					return nil, ErrDecryptK8sContext(err)
				}
				ctx.Auth = sql.Map(decryptedAuth)
			}
			if ctx.Cluster != nil {
				decryptedCluster, err := mkcp.EncSvc.DecryptMap(ctx.Cluster)
				if err != nil {
					return nil, ErrDecryptK8sContext(err)
				}
				ctx.Cluster = sql.Map(decryptedCluster)
			}
		}
	}

	mesheryK8sContextPage := MesheryK8sContextPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Contexts:   contexts,
	}

	resp, _ := json.Marshal(mesheryK8sContextPage)
	return resp, nil
}

func (mkcp *MesheryK8sContextPersister) SaveMesheryK8sContext(mkc K8sContext) (connections.Connection, error) {
	conn := connections.Connection{}
	if mkc.ID == "" {
		id, err := K8sContextGenerateID(mkc)
		if err != nil {
			return conn, ErrContextID
		}

		mkc.ID = id
	}

	// Encrypt only after K8sContextGenerateID has run. The ID is derived from
	// Auth and Cluster, and each Seal uses a fresh nonce, so encrypting first
	// would produce a different ID on every save and break identity deduplication.
	if mkcp.EncSvc != nil {
		if mkc.Auth != nil {
			encryptedAuth, err := mkcp.EncSvc.EncryptMap(mkc.Auth)
			if err != nil {
				return conn, ErrEncryptK8sContext(err)
			}
			mkc.Auth = sql.Map(encryptedAuth)
		}
		if mkc.Cluster != nil {
			encryptedCluster, err := mkcp.EncSvc.EncryptMap(mkc.Cluster)
			if err != nil {
				return conn, ErrEncryptK8sContext(err)
			}
			mkc.Cluster = sql.Map(encryptedCluster)
		}
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

	return conn, err
}

func (mkcp *MesheryK8sContextPersister) GetMesheryK8sContext(id string) (K8sContext, error) {
	var mesheryK8sContext K8sContext

	err := mkcp.DB.First(&mesheryK8sContext, "id = ?", id).Error
	if err != nil {
		return mesheryK8sContext, err
	}

	if mkcp.EncSvc != nil {
		if mesheryK8sContext.Auth != nil {
			decryptedAuth, err := mkcp.EncSvc.DecryptMap(mesheryK8sContext.Auth)
			if err != nil {
				return mesheryK8sContext, ErrDecryptK8sContext(err)
			}
			mesheryK8sContext.Auth = sql.Map(decryptedAuth)
		}
		if mesheryK8sContext.Cluster != nil {
			decryptedCluster, err := mkcp.EncSvc.DecryptMap(mesheryK8sContext.Cluster)
			if err != nil {
				return mesheryK8sContext, ErrDecryptK8sContext(err)
			}
			mesheryK8sContext.Cluster = sql.Map(decryptedCluster)
		}
	}

	return mesheryK8sContext, nil
}

