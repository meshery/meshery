package models

import (
	"database/sql"
	"strings"
	"reflect"
	"errors"
	"fmt"
	"time"

	"github.com/layer5io/meshkit/database"
	"gorm.io/gorm"
	"github.com/gofrs/uuid"
)

// MesheryK8sContextPersister is the persister for persisting
// applications on the database
type ConnectionPersister struct {
	DB *database.Handler
}

type MesheryK8sContextPage struct {
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalCount int           `json:"total_count"`
	Contexts   []*K8sContext `json:"contexts"`
}


var (
  queryConnections = `
	SELECT
	id,
	name,
	COALESCE(credential_id, uuid_nil()) AS credential_id,
	type,
	sub_type,
	kind,
	metadata,
	status,
	user_id,
	created_at,
	updated_at,
	deleted_at
	FROM connections
	`

	queryK8sContexts = `with k8scontexts as (
Select 
conn.metadata->>'id' as id,
conn.id as connection_id,
conn.metadata->>'name' as name,
cred.secret->'auth' as auth,
cred.secret->'cluster' as cluster,
conn.metadata->>'server' as server,
COALESCE(conn.metadata->>'version', 'NA') as version,
uuid(conn.metadata->>'meshery_instance_id') as meshery_instance_id,
uuid(conn.metadata->>'kubernetes_server_id') as kubernetes_server_id,
conn.user_id as owner,
conn.created_at as created_at,
conn.updated_at as updated_at
from connections conn
join credentials cred on cred.id = conn.credential_id
where conn.kind = 'kubernetes' and conn.deleted_at is null
)`
)

func (c *ConnectionPersister) UpsertConnection(connection *Connection) (*Connection, error) {

	// TODO: Fix this hack and make things generic
	// if credential id is missing then skip that column, else to add it
	// for the case of GitHub connection, credential is not stored in credential table but for all other connections it is stored in credential table
	// if connection.CredentialID == uuid.Nil {
	// 	if err := c.DB.Create(connection, "credential_id"); err != nil {
	// 		return nil, err
	// 	}
	// } else {
	// 	if err := c.cxn.Create(connection); err != nil {
	// 		return nil, err
	// 	}
	// }

	if err := c.DB.Create(connection); err != nil {
			return nil, err.Error
		}

	createdConnection := &Connection{}
	if err := c.DB.Raw(fmt.Sprintf("%s WHERE name = ? and deleted_at is null order by updated_at desc", queryConnections), connection.Name).First(createdConnection); err != nil {
		return nil, err.Error
	}

	return createdConnection, nil
}


func (c *ConnectionPersister) CheckK8sExistence(connection *Connection) (*Connection, error) {
	var k8sID string
	var mesheryInstanceID string

	v := reflect.ValueOf(connection.Metadata)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			if key.Interface().(string) == "id" {
				k8sID = strct.Interface().(string)
			}
			if key.Interface().(string) == "meshery_instance_id" {
				mesheryInstanceID = strct.Interface().(string)
			}
		}
	}

	if k8sID == "" {
		return nil, errors.New("empty k8s id; cannot check for existence")
	}

	existingConnection := Connection{}
	if err := c.DB.Raw(fmt.Sprintf("%s where metadata->>'id' = ? and kind = ? and metadata->>'meshery_instance_id' = ? and deleted_at is null", queryConnections), k8sID, connection.Kind, mesheryInstanceID).First(&existingConnection); err != nil {
		if err.Error.Error() != sql.ErrNoRows.Error() {	
			return nil, err.Error
		}
		if err.Error.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		}
	}

	return &existingConnection, nil
}


func (d *ConnectionPersister) GetConnections(page, pageSize int, search, order string) (*ConnectionPage, error) {
	query := fmt.Sprintf(`
	%s
	`, queryConnections)

	if order == "" {
		order = "updated_at desc"
	}


	var like string
	if search != "" {
		like = "%" + strings.ToLower(search) + "%"
		query += ` where (
			lower(name) like ? or
			lower(type) like ? or
			lower(sub_type) like ? or
			lower(kind) like ? or
			lower(status) like ?
		)`
	}

	query += " order by " + order

	var qc *gorm.DB
	if like == "" {
		qc = d.DB.Raw(query)
	} else {
		qc = d.DB.Raw(query, like, like, like, like, like)
	}

	var totalCount int64
	qc.Model(&Connection{}).Count(&totalCount)
	

	connections := []*Connection{}

	Paginate(uint(page), uint(pageSize))(qc).Find(&connections)
	// if totalCount > 0 {
	// 	if err := qc.Paginate(page+1, pageSize).All(&connections); err != nil {
	// 		if err.Error() != sql.ErrNoRows.Error() {
	// 			logrus.Errorf("error retrieving connections: %v", err)
	// 			return nil, err
	// 		}
	// 	}
	// }

	return &ConnectionPage{
		Connections: connections,
		TotalCount:  int(totalCount),
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

func (c *ConnectionPersister) UpdateConnectionById(conn *Connection) (*Connection, error) {
	if err := c.DB.Model(Connection{}).Updates(conn); err != nil {
		return nil, err.Error
	}

	updatedConnection := &Connection{}
	if err := c.DB.Raw(fmt.Sprintf("%s WHERE id = ? and deleted_at is null order by updated_at desc", queryConnections), conn.ID).First(updatedConnection); err != nil {
		return nil, err.Error
	}

	return updatedConnection, nil
}


func (c *ConnectionPersister) GetConnectionById(connID string) (*Connection, error) {
	conn := &Connection{}
	if err := c.DB.Raw(fmt.Sprintf("%s WHERE id = ? and deleted_at is null", queryConnections), connID).First(conn); err != nil {
		return nil, err.Error
	}

	return conn, nil
}

func (c *ConnectionPersister) DeleteConnection(connID *uuid.UUID) (*Connection, error) {
	conn := &Connection{}
	if err := c.DB.Raw(fmt.Sprintf("%s where id = ? and deleted_at is NULL", queryConnections), connID).First(conn); err != nil {
		return nil, err.Error
	}

	if err := c.DB.Raw("UPDATE connections SET deleted_at = ?, status = ? WHERE id = ?", time.Now(), DELETED, connID); err != nil {
		return nil, err.Error
	}

	return conn, nil
}

func (c *ConnectionPersister) GetK8sContexts(userID, instanceID string, page, pageSize int, search, order string) (*MesheryK8sContextPage, error) {
	order = sanitizeOrderInput(order, []string{"updated_at", "created_at", "name"})
	if order == "" {
		order = "updated_at desc"
	}

	query := fmt.Sprintf(`%s
select k8scontexts.* from k8scontexts
where k8scontexts.meshery_instance_id = ? and k8scontexts.owner = ?`, queryK8sContexts)

	var like string
	if search != "" {
		like = "%" + strings.ToLower(search) + "%"
		query += ` and ( lower(k8scontexts.name) like ? )`
	}

	query += " order by " + order

	var qc *gorm.DB
	if like == "" {
		qc = c.DB.Raw(query, instanceID, userID)
	} else {
		qc = c.DB.Raw(query, instanceID, userID, like)
	}

	var totalCount int64
	qc.Model(&Connection{}).Count(&totalCount)
	
	contexts := []*K8sContext{}

	Paginate(uint(page), uint(pageSize))(qc).Find(&contexts)


	return &MesheryK8sContextPage{
		Contexts:   contexts,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(totalCount),
	}, nil
}


func (c *ConnectionPersister) GetK8sContext(userID, connectionID string) (*K8sContext, error) {
	context := &K8sContext{}

	query := fmt.Sprintf(`%s 
select k8scontexts.* from k8scontexts
where k8scontexts.connection_id = ? and k8scontexts.owner = ?`, queryK8sContexts)

	if err := c.DB.Raw(query, connectionID, userID).First(context); err != nil {
		if err.Error.Error() != sql.ErrNoRows.Error() {
			return nil, err.Error
		}
	}

	return context, nil
}



// // GetMesheryK8sContexts returns all of the contexts
// func (mkcp *MesheryK8sContextPersister) GetMesheryK8sContexts(search, order string, page, pageSize uint64) ([]byte, error) {
// 	order = sanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

// 	if order == "" {
// 		order = "updated_at desc"
// 	}

// 	count := int64(0)
// 	contexts := []*K8sContext{}

// 	query := mkcp.DB.Order(order)

// 	if search != "" {
// 		like := "%" + strings.ToLower(search) + "%"
// 		query = query.Where("(lower(name) like ?)", like)
// 	}

// 	query.Model(K8sContext{}).Count(&count)

// 	Paginate(uint(page), uint(pageSize))(query).Find(&contexts)

// 	mesheryK8sContextPage := MesheryK8sContextPage{
// 		Page:       page,
// 		PageSize:   pageSize,
// 		TotalCount: int(count),
// 		Contexts:   contexts,
// 	}

// 	resp, _ := json.Marshal(mesheryK8sContextPage)
// 	return resp, nil
// }

// // DeleteMesheryK8sContext takes in an application id and delete it if it already exists
// func (mkcp *MesheryK8sContextPersister) DeleteMesheryK8sContext(id string) (K8sContext, error) {
// 	context := K8sContext{ID: id}
// 	mkcp.DB.Delete(&context)

// 	return context, nil
// }

// func (mkcp *MesheryK8sContextPersister) SaveMesheryK8sContext(mkc K8sContext) (K8sContext, error) {
// 	if mkc.ID == "" {
// 		id, err := K8sContextGenerateID(mkc)
// 		if err != nil {
// 			return mkc, ErrContextID
// 		}

// 		mkc.ID = id
// 	}

// 	// Perform the operation in a transaction
// 	err := mkcp.DB.Transaction(func(tx *gorm.DB) error {
// 		var mesheryK8sContext K8sContext

// 		// Check if there is already an entry for this context
// 		if err := tx.First(&mesheryK8sContext, "id = ?", mkc.ID).Error; err == nil {
// 			return ErrContextAlreadyPersisted
// 		}

// 		return tx.Save(&mkc).Error
// 	})

// 	return mkc, err
// }

// func (mkcp *MesheryK8sContextPersister) GetMesheryK8sContext(id string) (K8sContext, error) {
// 	var mesheryK8sContext K8sContext

// 	err := mkcp.DB.First(&mesheryK8sContext, "id = ?", id).Error
// 	return mesheryK8sContext, err
// }

// // func (mkcp *MesheryK8sContextPersister) SetMesheryK8sCurrentContext(id string) error {
// // 	// Perform the operation in a transaction
// // 	return mkcp.DB.Transaction(func(tx *gorm.DB) error {
// // 		var mesheryK8sContext K8sContext

// // 		// Get context which is currently in use
// // 		if err := tx.First(&mesheryK8sContext, "is_current_context = true").Error; err != nil {
// // 			return err
// // 		}

// // 		// If the context id matches with the provided id then skip the next steps
// // 		if mesheryK8sContext.ID == id {
// // 			return nil
// // 		}

// // 		if err := tx.Save(&mesheryK8sContext).Error; err != nil {
// // 			return err
// // 		}

// // 		// Set the specified context as active
// // 		return tx.Model(K8sContext{}).Where("id = ?", id).Update("is_current_context", true).Error
// // 	})
// // }

// // func (mkcp *MesheryK8sContextPersister) GetMesheryK8sCurrentContext() (K8sContext, error) {
// // 	var mesheryK8sContext K8sContext

// // 	err := mkcp.DB.First(&mesheryK8sContext, "is_current_context = true").Error

// // 	return mesheryK8sContext, err
// // }
