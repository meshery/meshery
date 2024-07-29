package models

import (
	"encoding/json"
	"strconv"
	"strings"

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
func (ep *EnvironmentPersister) GetEnvironments(search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = "updated_at desc"
	}

	// Convert page and pageSize from string to uint
	pageUint, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, err
	}
	pageSizeUint, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		return nil, err
	}

	// Prepare the query
	query := ep.DB.Model(&v1beta1.Environment{})
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

	// Fetch environments with pagination
	environmentsFetched := []v1beta1.Environment{}
	Paginate(uint(pageUint), uint(pageSizeUint))(query).Find(&environmentsFetched)

	// Prepare the response
	environmentsPage := &v1beta1.EnvironmentPage{
		Page:         int(pageUint),
		PageSize:     int(pageSizeUint),
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
func (ep *EnvironmentPersister) AddConnectionToEnvironment(environmentID, connectionID uuid.UUID) (*v1beta1.Environment, error) {
	env, err := ep.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}

	// Add connection to environment

	return env, nil
}

// GetEnvironmentConnections returns connections for an environment
func (ep *EnvironmentPersister) GetEnvironmentConnections(environmentID uuid.UUID) ([]*connections.Connection, error) {
	// env, err := ep.GetEnvironment(environmentID)
	// if err != nil {
	// 	return nil, err
	// }

	// // Get connections of an environment
	// var conns []*connections.Connection

	// return conns, nil
	return nil, nil
}

// DeleteConnectionFromEnvironment deletes a connection from an environment
func (ep *EnvironmentPersister) DeleteConnectionFromEnvironment(environmentID, connectionID uuid.UUID) (*v1beta1.Environment, error) {
	env, err := ep.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}

	// Delete connection from environment

	return env, nil
}
