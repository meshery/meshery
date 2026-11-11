package models

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/meshery/schemas/models/core"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/v1beta3/environment"
	schemasWorkspace "github.com/meshery/schemas/models/v1beta3/workspace"
	"gorm.io/gorm"
)

// EnvironmentPersister is the persister for persisting
// environments on the database
type EnvironmentPersister struct {
	DB *database.Handler
}

func (ep *EnvironmentPersister) fetchUserDetails() *User {
	return LocalProviderUser()
}

// GetEnvironments returns all of the environments
func (ep *EnvironmentPersister) GetEnvironments(orgID, search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	query := ep.DB.Model(&environment.Environment{})

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

	environmentsFetched := []environment.Environment{}

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
	environmentsPage := &environment.EnvironmentPage{
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

func (ep *EnvironmentPersister) SaveEnvironment(environment *environment.Environment) ([]byte, error) {
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

// DeleteEnvironment deletes an environment and cleans up its connection and workspace mappings.
func (ep *EnvironmentPersister) DeleteEnvironment(env *environment.Environment) ([]byte, error) {
	var existing environment.Environment
	err := ep.DB.Where("id = ?", env.ID).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
		return nil, ErrDBRead(err)
	}
	err = ep.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(env).Error; err != nil {
			return err
		}
		if tx.Migrator().HasTable(&environment.EnvironmentConnectionMapping{}) {
			if err := tx.Where("environment_id = ?", env.ID).Delete(&environment.EnvironmentConnectionMapping{}).Error; err != nil {
				return err
			}
		}
		if tx.Migrator().HasTable(&schemasWorkspace.WorkspacesEnvironmentsMapping{}) {
			if err := tx.Where("environment_id = ?", env.ID).Delete(&schemasWorkspace.WorkspacesEnvironmentsMapping{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, ErrDBDelete(err, ep.fetchUserDetails().ID.String())
	}

	// Marshal the environment to JSON
	envJSON, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}

	return envJSON, nil
}

func (ep *EnvironmentPersister) UpdateEnvironmentByID(env *environment.Environment) (*environment.Environment, error) {
	err := ep.DB.Save(env).Error
	if err != nil {
		return nil, ErrDBPut(err)
	}

	updatedEnvironment := environment.Environment{}
	err = ep.DB.Model(&updatedEnvironment).Where("id = ?", env.ID).First(&updatedEnvironment).Error
	if err != nil {
		return nil, ErrDBRead(err)
	}
	return env, nil
}

// GetEnvironment returns an environment by ID
func (ep *EnvironmentPersister) GetEnvironment(id core.Uuid) (*environment.Environment, error) {
	environment := environment.Environment{}
	query := ep.DB.Where("id = ?", id)
	err := query.First(&environment).Error
	return &environment, err
}

// GetEnvironmentByID returns a single environment by ID
func (ep *EnvironmentPersister) GetEnvironmentByID(environmentID core.Uuid) ([]byte, error) {
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

// UpdateEnvironment updates an environment's fields from the given payload and persists the change
func (ep *EnvironmentPersister) UpdateEnvironment(environmentID core.Uuid, payload *environment.EnvironmentPayload) (*environment.Environment, error) {
	env, err := ep.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}

	env.Name = payload.Name
	env.Description = payload.Description
	env.OrganizationID = core.Uuid(payload.OrgID)

	return ep.UpdateEnvironmentByID(env)
}

// DeleteEnvironmentByID deletes a single environment by ID
func (ep *EnvironmentPersister) DeleteEnvironmentByID(environmentID core.Uuid) ([]byte, error) {
	env, err := ep.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}

	return ep.DeleteEnvironment(env)
}

// AddConnectionToEnvironment adds a connection to an environment. Adding one the
// environment already has returns the existing mapping instead of a second row.
func (ep *EnvironmentPersister) AddConnectionToEnvironment(environmentID, connectionID core.Uuid) ([]byte, error) {
	var existing environment.EnvironmentConnectionMapping
	found := ep.DB.Where("environment_id = ? AND connection_id = ?", environmentID, connectionID).Limit(1).Find(&existing)
	if found.Error != nil {
		return nil, ErrDBRead(found.Error)
	}
	if found.RowsAffected > 0 {
		return json.Marshal(existing)
	}

	envConMapping := environment.EnvironmentConnectionMapping{
		ConnectionID:  &connectionID,
		EnvironmentID: &environmentID,
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
func (ep *EnvironmentPersister) GetEnvironmentConnections(environmentID core.Uuid, search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	isAssigned := true
	deletedAt := false
	// Parse the filter JSON string
	if filter != "" {
		var filterMap map[string]interface{}
		err := json.Unmarshal([]byte(filter), &filterMap)
		if err != nil {
			return nil, err
		}

		if assignedVal, ok := filterMap["assigned"]; ok {
			switch v := assignedVal.(type) {
			case bool:
				isAssigned = v
			case string:
				if b, err := strconv.ParseBool(v); err == nil {
					isAssigned = b
				}
			}
		}

		if delVal, ok := filterMap["deletedAt"]; ok {
			switch v := delVal.(type) {
			case bool:
				deletedAt = v
			case string:
				if b, err := strconv.ParseBool(v); err == nil {
					deletedAt = b
				}
			}
		} else if delVal, ok := filterMap["deleted_at"]; ok {
			switch v := delVal.(type) {
			case bool:
				deletedAt = v
			case string:
				if b, err := strconv.ParseBool(v); err == nil {
					deletedAt = b
				}
			}
		}
	}

	var query *gorm.DB
	if isAssigned {
		// Query for connections that are assigned to the given environment
		query = ep.DB.Table("connections").
			Where("EXISTS (SELECT 1 FROM environment_connection_mappings AS ecm WHERE connections.id = ecm.connection_id AND ecm.environment_id = ? AND ecm.deleted_at IS NULL)", environmentID)
	} else {
		// Query for connections that are not assigned to the given environment
		query = ep.DB.Table("connections").
			Where("NOT EXISTS (SELECT 1 FROM environment_connection_mappings AS ecm WHERE connections.id = ecm.connection_id AND ecm.environment_id = ? AND ecm.deleted_at IS NULL)", environmentID)
	}

	if deletedAt {
		query = query.Where("connections.deleted_at IS NOT NULL")
	} else {
		query = query.Where("connections.deleted_at IS NULL")
	}

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
	if err := query.Count(&count).Error; err != nil {
		return nil, ErrDBRead(err)
	}

	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	var connectionsFetched []*connections.Connection
	pageUint, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, err
	}
	// Fetch all connections if pageSize is "all"
	if pageSize == "all" {
		if err := query.Find(&connectionsFetched).Error; err != nil {
			return nil, ErrDBRead(err)
		}
	} else {
		// Convert page and pageSize from string to uint
		pageSizeUint, err := strconv.ParseUint(pageSize, 10, 32)
		if err != nil {
			return nil, err
		}

		// Fetch connections with pagination
		if err := Paginate(uint(pageUint), uint(pageSizeUint))(query).Find(&connectionsFetched).Error; err != nil {
			return nil, ErrDBRead(err)
		}
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

// DeleteConnectionFromEnvironment deletes a connection from an environment.
// Every mapping of the pair is deleted in a single transaction, including
// duplicates stored before assignments were checked.
func (ep *EnvironmentPersister) DeleteConnectionFromEnvironment(environmentID, connectionID core.Uuid) ([]byte, error) {
	var envConMapping environment.EnvironmentConnectionMapping

	err := ep.DB.Transaction(func(tx *gorm.DB) error {
		// Find the specific connection mapping
		if err := tx.Where("environment_id = ? AND connection_id = ?", environmentID, connectionID).First(&envConMapping).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrResultNotFound(err)
			}
			return ErrDBRead(err)
		}

		// Delete every mapping of the pair, including duplicates stored before
		// assignments were checked.
		if err := tx.Where("environment_id = ? AND connection_id = ?", environmentID, connectionID).Delete(&environment.EnvironmentConnectionMapping{}).Error; err != nil {
			return ErrDBDelete(err, ep.fetchUserDetails().ID.String())
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	envJSON, err := json.Marshal(envConMapping)
	if err != nil {
		return nil, err
	}

	return envJSON, nil
}
