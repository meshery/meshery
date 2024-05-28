package models

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/database"
	"gorm.io/gorm"
)

// ConnectionPersister is the persister for persisting
// connections on the database
type ConnectionPersister struct {
	DB *database.Handler
}

// GetConnections returns all of the connections
func (cp *ConnectionPersister) GetConnections(search, order string, page, pageSize int, filter string, status []string, kind []string) (*connections.ConnectionPage, error) {
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = "updated_at desc"
	}

	query := cp.DB.Model(&connections.Connection{})

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(name) like ?", like)
	}

	if len(status) != 0 {
		query = query.Where("status IN (?)", status)
	}

	if len(kind) != 0 {
		query = query.Where("kind IN (?)", kind)
	}

	if filter != "" {
		filterArr := strings.Split(filter, " ")
		filterKey := filterArr[0]
		filterVal := strings.Join(filterArr[1:], " ")

		if filterKey == "deleted_at" {
			// Handle deleted_at filter
			if filterVal == "Deleted" {
				query = query.Where("deleted_at IS NOT NULL")
			} else {
				query = query.Where("deleted_at IS NULL")
			}
		} else if filterKey == "type" || filterKey == "sub_type" {
			query = query.Where(fmt.Sprintf("%s = ?", filterKey), filterVal)
		}
	}

	query = query.Order(order)
	count := int64(0)

	connectionsFetched := []*connections.Connection{}

	query.Table("connections").Count(&count)

	Paginate(uint(page), uint(pageSize))(query).Find(&connectionsFetched)
	// query.Table("connections").Count(&count).Find(&connectionsFetched)
	connectionsPage := &connections.ConnectionPage{
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  int(count),
		Connections: connectionsFetched,
	}

	return connectionsPage, nil
}

func (cp *ConnectionPersister) SaveConnection(connection *connections.Connection) (*connections.Connection, error) {
	if connection.ID == uuid.Nil {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}
		connection.ID = id
	}

	err := cp.DB.Transaction(func(tx *gorm.DB) error {
		existingConnection := connections.Connection{}

		// Check if there is already an entry for this context
		if err := tx.First(&existingConnection, "id = ?", connection.ID).Error; err == nil {
			return err
		}

		return tx.Save(&connection).Error
	})

	return connection, err
}

func (cp *ConnectionPersister) UpdateConnectionStatusByID(connectionID uuid.UUID, connectionStatus connections.ConnectionStatus) (*connections.Connection, error) {
	err := cp.DB.Model(&connections.Connection{}).Where("id = ?", connectionID).UpdateColumn("status", connectionStatus).Error
	if err != nil {
		return nil, fmt.Errorf("error updating connection status: %v", err)
	}

	updatedConnection := connections.Connection{}
	err = cp.DB.Model(&updatedConnection).Where("id = ?", connectionID).First(&updatedConnection).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving updated connection: %v", err)
	}

	return &updatedConnection, nil
}

// Get connection by ID
// If kind is provided filter with kind too
func (cp *ConnectionPersister) GetConnection(id uuid.UUID, kind string) (*connections.Connection, error) {
	connection := connections.Connection{}
	query := cp.DB.Where("id = ?", id)
	if kind != "" {
		query = query.Where("kind = ?", kind)
	}
	err := query.First(&connection).Error
	return &connection, err
}
