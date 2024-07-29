package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshkit/database"
	"github.com/meshery/schemas/models/v1beta1"
	"gorm.io/gorm"
)

// WorkspacePersister is the persister for persisting
// workspaces on the database
type WorkspacePersister struct {
	DB *database.Handler
}

func (wp *WorkspacePersister) fetchUserDetails() *User {
	return &User{
		UserID:    "meshery",
		FirstName: "Meshery",
		LastName:  "Meshery",
		AvatarURL: "",
	}
}

// GetWorkspaces returns all of the workspaces
func (wp *WorkspacePersister) GetWorkspaces(orgID, search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = "updated_at desc"
	}

	query := wp.DB.Model(&v1beta1.Workspace{})
	
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
	query.Table("workspaces").Count(&count)

	workspacesFetched := []v1beta1.Workspace{}
	pageUint, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, err
	}
	// Fetch all workspaces if pageSize is "all"
	if pageSize == "all" {
		query.Find(&workspacesFetched)
	} else {
		// Convert page and pageSize from string to uint
		pageSizeUint, err := strconv.ParseUint(pageSize, 10, 32)
		if err != nil {
			return nil, err
		}
		
		// Fetch workspaces with pagination
		Paginate(uint(pageUint), uint(pageSizeUint))(query).Find(&workspacesFetched)
	}

	// Prepare the response
	workspacesPage := &v1beta1.WorkspacePage{
		Page:         int(pageUint),
		PageSize:     len(workspacesFetched),
		TotalCount:   int(count),
		Workspaces:   workspacesFetched,
	}

	// Marshal the response to JSON
	wsJSON, err := json.Marshal(workspacesPage)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) SaveWorkspace(workspace *v1beta1.Workspace) ([]byte, error) {
	if workspace.ID == uuid.Nil {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}
		workspace.ID = id
	}

	if err := wp.DB.Create(workspace).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	wsJSON, err := json.Marshal(workspace)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) DeleteWorkspace(workspace *v1beta1.Workspace) ([]byte, error) {
	err := wp.DB.Model(&workspace).Find(&workspace).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
	}
	err = wp.DB.Delete(workspace).Error
	if err != nil {
		return nil, ErrDBDelete(err, wp.fetchUserDetails().UserID)
	}

	// Marshal the workspace to JSON
	wsJSON, err := json.Marshal(workspace)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) UpdateWorkspaceByID(workspace *v1beta1.Workspace) (*v1beta1.Workspace, error) {
	err := wp.DB.Save(workspace).Error
	if err != nil {
		return nil, ErrDBPut(err)
	}

	updatedWorkspace := v1beta1.Workspace{}
	err = wp.DB.Model(&updatedWorkspace).Where("id = ?", workspace.ID).First(&updatedWorkspace).Error
	if err != nil {
		return nil, ErrDBRead(err)
	}
	return workspace, nil
}

// Get workspace by ID
func (wp *WorkspacePersister) GetWorkspace(id uuid.UUID) (*v1beta1.Workspace, error) {
	workspace := v1beta1.Workspace{}
	query := wp.DB.Where("id = ?", id)
	err := query.First(&workspace).Error
	return &workspace, err
}

// GetWorkspaceByID returns a single workspace by ID
func (wp *WorkspacePersister) GetWorkspaceByID(workspaceID uuid.UUID) ([]byte, error) {
	workspace, err := wp.GetWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}
	// Marshal the workspace to JSON
	wsJSON, err := json.Marshal(workspace)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

// UpdateWorkspaceByID updates a single workspace by ID
func (wp *WorkspacePersister) UpdateWorkspace(workspaceID uuid.UUID, payload *v1beta1.WorkspacePayload) (*v1beta1.Workspace, error) {
	ws, err := wp.GetWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}

	ws.Name = payload.Name
	ws.Description = payload.Description
	ws.OrganizationId = uuid.FromStringOrNil(payload.OrgId)

	return wp.UpdateWorkspaceByID(ws)
}

// DeleteWorkspaceByID deletes a single workspace by ID
func (wp *WorkspacePersister) DeleteWorkspaceByID(workspaceID uuid.UUID) ([]byte, error) {
	ws, err := wp.GetWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}

	return wp.DeleteWorkspace(ws)
}

// AddEnvironmentToWorkspace adds an environment to a workspace
func (wp *WorkspacePersister) AddEnvironmentToWorkspace(workspaceID, environmentID uuid.UUID) ([]byte, error) {
	wsEnvMapping := v1beta1.WorkspaceEnvironmentMapping{
		EnvironmentId: environmentID,
		WorkspaceId:   workspaceID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, ErrGenerateUUID(err)
	}
	wsEnvMapping.ID = id

	// Add environment to workspace
	if err := wp.DB.Create(wsEnvMapping).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	wsJSON, err := json.Marshal(wsEnvMapping)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

// GetWorkspaceEnvironments returns environments for a workspace
func (wp *WorkspacePersister) GetWorkspaceEnvironments(workspaceID uuid.UUID, search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = "updated_at desc"
	}

	// Build the query to find environments associated with the given workspace ID
	query := wp.DB.Table("workspace_environment_mappings").
		Select("environments.*").
		Joins("JOIN environments ON workspace_environment_mappings.environment_id = environments.id").
		Where("workspace_environment_mappings.workspace_id = ?", workspaceID)

	// Apply search filter
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(environments.name) LIKE ?", like)
	}

	// Apply additional filters
	dynamicKeys := []string{"owner", "organization_id"}
	query = utils.ApplyFilters(query, filter, dynamicKeys)
	query = query.Order(order)

	count := int64(0)
	query.Count(&count)

	var environmentsFetched []*v1beta1.Environment
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

	environmentsPage := &v1beta1.EnvironmentPage{
		Page:         int(pageUint),
		PageSize:     len(environmentsFetched),
		TotalCount:   int(count),
		Environments: environmentsFetched,
	}

	envsJSON, err := json.Marshal(environmentsPage)
	if err != nil {
		return nil, err
	}

	return envsJSON, nil
}

// DeleteEnvironmentFromWorkspace deletes an environment from a workspace
func (wp *WorkspacePersister) DeleteEnvironmentFromWorkspace(workspaceID, environmentID uuid.UUID) ([]byte, error) {
	var wsEnvMapping v1beta1.WorkspaceEnvironmentMapping

	// Find the specific environment mapping
	if err := wp.DB.Where("workspace_id = ? AND environment_id = ?", workspaceID, environmentID).First(&wsEnvMapping).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
		return nil, ErrDBRead(err)
	}

	// Delete the environment mapping
	if err := wp.DB.Delete(&wsEnvMapping).Error; err != nil {
		return nil, ErrDBDelete(err, wp.fetchUserDetails().UserID)
	}

	wsJSON, err := json.Marshal(wsEnvMapping)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}
