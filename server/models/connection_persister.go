package models

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshery/server/models/environments"
	"github.com/meshery/meshkit/database"
	schemasConnection "github.com/meshery/schemas/models/v1beta1/connection"
	"gorm.io/gorm"
)

// ConnectionPersister is the persister for persisting
// connections on the database
type ConnectionPersister struct {
	DB *database.Handler
}

// GetConnections returns all of the connections
func (cp *ConnectionPersister) GetConnections(search, order string, page, pageSize int, filter string, status []string, kind []string, connType []string, name string) (*connections.ConnectionPage, error) {
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = "updated_at desc"
	}

	query := cp.DB.Model(&connections.Connection{})

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(name) like ?", like)
	}

	if name != "" {
		like := "%" + strings.ToLower(name) + "%"
		query = query.Where("lower(name) like ?", like)
	}

	if len(status) != 0 {
		query = query.Where("status IN (?)", status)
	}

	if len(kind) != 0 {
		query = query.Where("kind IN (?)", kind)
	}

	if len(connType) != 0 {
		query = query.Where("type IN (?)", connType)
	}

	dynamicKeys := []string{"type", "sub_type"}

	// Apply filters using the utility function
	query = utils.ApplyFilters(query, filter, dynamicKeys)

	query = query.Order(order)
	count := int64(0)

	connectionsFetched := []*connections.Connection{}
	query.Table("connections").Count(&count)
	environmentsFetched := []*environments.EnvironmentData{}
	Paginate(uint(page), uint(pageSize))(query).Find(&connectionsFetched)

	for _, connectionFetched := range connectionsFetched {
		cp.DB.Table("environment_connection_mappings").Joins("LEFT JOIN environments ON environments.id = environment_connection_mappings.environment_id").Select("environments.*").
			Where("connection_id = ?", connectionFetched.ID).
			Find(&environmentsFetched)

		connectionFetched.Environments = environmentsFetched
	}
	statusSummary, err := cp.getConnectionsStatusSummary()
	if err != nil {
		return nil, err
	}

	connectionsPage := &connections.ConnectionPage{
		Page:          page,
		PageSize:      pageSize,
		TotalCount:    int(count),
		Connections:   connectionsFetched,
		StatusSummary: statusSummary,
	}

	return connectionsPage, nil
}

// getConnectionsStatusSummary returns a map of connection status to count
func (cp *ConnectionPersister) getConnectionsStatusSummary() (*map[schemasConnection.ConnectionStatus]int, error) {
	var statusCounts []struct {
		Status string `gorm:"column:status"`
		Count  int    `gorm:"column:count"`
	}

	err := cp.DB.Model(&connections.Connection{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching connection status summary: %v", err)
	}

	summary := make(map[schemasConnection.ConnectionStatus]int)
	for _, sc := range statusCounts {
		summary[schemasConnection.ConnectionStatus(sc.Status)] = sc.Count
	}

	return &summary, nil
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

func (cp *ConnectionPersister) DeleteConnection(connection *connections.Connection) (*connections.Connection, error) {
	err := cp.DB.Model(&connection).Find(&connection).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
	}
	err = cp.DB.Delete(connection).Error
	if err != nil {
		return nil, ErrDBDelete(err, cp.fetchUserDetails().UserId)
	}

	return connection, err
}

func (cp *ConnectionPersister) fetchUserDetails() *User {

	return &User{
		UserId:    "meshery",
		FirstName: "Meshery",
		LastName:  "Meshery",
	}
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

func (cp *ConnectionPersister) UpdateConnectionByID(connection *connections.Connection) (*connections.Connection, error) {
	err := cp.DB.Save(connection).Error
	if err != nil {
		return nil, ErrDBPut(err)
	}

	updatedConnection := connections.Connection{}
	err = cp.DB.Model(&updatedConnection).Where("id = ?", connection.ID).First(&updatedConnection).Error
	if err != nil {
		return nil, ErrDBRead(err)
	}
	return connection, nil
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
