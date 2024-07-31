package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/database"
	"github.com/meshery/schemas/models/v1beta1"
	"gorm.io/gorm"
)

// EnvironmentPersister is the persister for persisting
// environments on the database
type EnvironmentPersister struct {
	DB *database.Handler
}

func (ep *EnvironmentPersister) fetchUserDetails() *User {
	return &User{
		UserID:    "meshery",
		FirstName: "Meshery",
		LastName:  "Meshery",
		AvatarURL: "",
	}
}

// GetEnvironments returns all of the environments
func (ep *EnvironmentPersister) GetEnvironments(orgID, search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = "updated_at desc"
	}

	query := ep.DB.Model(&v1beta1.Environment{})

	// Filter by organization ID
	if orgID != "" {
		query = query.Where("organization_id = ?", orgID)
	}

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(name) like ?", like)
	}

	dynamicKeys := []string{"owner", "organization_id"}
	query = utils.ApplyFilters(query, filter, dynamicKeys)
	query = query.Order(order)

	// Count total records
	count := int64(0)
	query.Table("environments").Count(&count)

	environmentsFetched := []v1beta1.Environment{}

	if page == "" {
		page = "0"
	}

	if pageSize == "" {
		pageSize = "10"
	}

	pageUint, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, err
	}
	// Fetch all environments if pageSize is "all"
	if pageSize == "all" {
		query.Find(&environmentsFetched)
	} else {
		// Convert page and pageSize from string to uint
		pageSizeUint, err := strconv.ParseUint(pageSize, 10, 32)
		if err != nil {
			return nil, err
		}

		// Fetch environments with pagination
		Paginate(uint(pageUint), uint(pageSizeUint))(query).Find(&environmentsFetched)
	}

	// Prepare the response
	environmentsPage := &v1beta1.EnvironmentPage{
		Page:         int(pageUint),
		PageSize:     len(environmentsFetched),
		TotalCount:   int(count),
		Environments: environmentsFetched,
	}

	// Marshal the response to JSON
	envJSON, err := json.Marshal(environmentsPage)
	if err != nil {
		return nil, err
	}

	return envJSON, nil
}

func (ep *EnvironmentPersister) SaveEnvironment(environment *v1beta1.Environment) ([]byte, error) {
	if environment.ID == uuid.Nil {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}
		environment.ID = id
	}

	if err := ep.DB.Create(environment).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	envJSON, err := json.Marshal(environment)
	if err != nil {
		return nil, err
	}

	return envJSON, nil
}

func (ep *EnvironmentPersister) DeleteEnvironment(environment *v1beta1.Environment) ([]byte, error) {
	err := ep.DB.Model(&environment).Find(&environment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
	}
	err = ep.DB.Delete(environment).Error
	if err != nil {
		return nil, ErrDBDelete(err, ep.fetchUserDetails().UserID)
	}

	// Marshal the environment to JSON
	envJSON, err := json.Marshal(environment)
	if err != nil {
		return nil, err
	}

	return envJSON, nil
}

func (ep *EnvironmentPersister) UpdateEnvironmentByID(environment *v1beta1.Environment) (*v1beta1.Environment, error) {
	err := ep.DB.Save(environment).Error
	if err != nil {
		return nil, ErrDBPut(err)
	}

	updatedEnvironment := v1beta1.Environment{}
	err = ep.DB.Model(&updatedEnvironment).Where("id = ?", environment.ID).First(&updatedEnvironment).Error
	if err != nil {
		return nil, ErrDBRead(err)
	}
	return environment, nil
}

// Get environment by ID
func (ep *EnvironmentPersister) GetEnvironment(id uuid.UUID) (*v1beta1.Environment, error) {
	environment := v1beta1.Environment{}
	query := ep.DB.Where("id = ?", id)
	err := query.First(&environment).Error
	return &environment, err
}

// GetEnvironmentByID returns a single environment by ID
func (ep *EnvironmentPersister) GetEnvironmentByID(environmentID uuid.UUID) ([]byte, error) {
	environment, err := ep.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}
	// Marshal the environment to JSON
	envJSON, err := json.Marshal(environment)
	if err != nil {
		return nil, err
	}

	return envJSON, nil
}

// UpdateEnvironmentByID updates a single environment by ID
func (ep *EnvironmentPersister) UpdateEnvironment(environmentID uuid.UUID, payload *v1beta1.EnvironmentPayload) (*v1beta1.Environment, error) {
	env, err := ep.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}

	env.Name = payload.Name
	env.Description = payload.Description
	env.OrganizationId = uuid.FromStringOrNil(payload.OrgId)

	return ep.UpdateEnvironmentByID(env)
}

// DeleteEnvironmentByID deletes a single environment by ID
func (ep *EnvironmentPersister) DeleteEnvironmentByID(environmentID uuid.UUID) ([]byte, error) {
	env, err := ep.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}

	return ep.DeleteEnvironment(env)
}

// AddConnectionToEnvironment adds a connection to an environment
func (ep *EnvironmentPersister) AddConnectionToEnvironment(environmentID, connectionID uuid.UUID) ([]byte, error) {
	envConMapping := v1beta1.EnvironmentConnectionMapping{
		ConnectionId:  connectionID,
		EnvironmentId: environmentID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, ErrGenerateUUID(err)
	}
	envConMapping.ID = id

	// Add connection to environment
	if err := ep.DB.Create(envConMapping).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	envJSON, err := json.Marshal(envConMapping)
	if err != nil {
		return nil, err
	}

	return envJSON, nil
}

// GetEnvironmentConnections returns connections for an environment
func (ep *EnvironmentPersister) GetEnvironmentConnections(environmentID uuid.UUID, search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = "updated_at desc"
	}

	// Build the query to find connections associated with the given environment ID
	query := ep.DB.Table("environment_connection_mappings").
		Select("connections.*").
		Joins("JOIN connections ON environment_connection_mappings.connection_id = connections.id").
		Where("environment_connection_mappings.environment_id = ?", environmentID)

	// Apply search filter
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(connections.name) LIKE ?", like)
	}

	// Apply additional filters
	dynamicKeys := []string{"owner", "organization_id"}
	query = utils.ApplyFilters(query, filter, dynamicKeys)
	query = query.Order(order)

	count := int64(0)
	query.Count(&count)

	var connectionsFetched []*connections.Connection
	pageUint, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, err
	}
	// Fetch all connections if pageSize is "all"
	if pageSize == "all" {
		query.Find(&connectionsFetched)
	} else {
		// Convert page and pageSize from string to uint
		pageSizeUint, err := strconv.ParseUint(pageSize, 10, 32)
		if err != nil {
			return nil, err
		}

		// Fetch connections with pagination
		Paginate(uint(pageUint), uint(pageSizeUint))(query).Find(&connectionsFetched)
	}

	connectionsPage := &connections.ConnectionPage{
		Page:        int(pageUint),
		PageSize:    len(connectionsFetched),
		TotalCount:  int(count),
		Connections: connectionsFetched,
	}

	connsJSON, err := json.Marshal(connectionsPage)
	if err != nil {
		return nil, err
	}

	return connsJSON, nil
}

// DeleteConnectionFromEnvironment deletes a connection from an environment
func (ep *EnvironmentPersister) DeleteConnectionFromEnvironment(environmentID, connectionID uuid.UUID) ([]byte, error) {
	var envConMapping v1beta1.EnvironmentConnectionMapping

	// Find the specific connection mapping
	if err := ep.DB.Where("environment_id = ? AND connection_id = ?", environmentID, connectionID).First(&envConMapping).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
		return nil, ErrDBRead(err)
	}

	// Delete the connection mapping
	if err := ep.DB.Delete(&envConMapping).Error; err != nil {
		return nil, ErrDBDelete(err, ep.fetchUserDetails().UserID)
	}

	envJSON, err := json.Marshal(envConMapping)
	if err != nil {
		return nil, err
	}

	return envJSON, nil
}
