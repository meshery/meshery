package models

import (
	"encoding/json"
	"fmt"
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

	if page == "" {
		page = "0"
	}

	if pageSize == "" {
		pageSize = "10"
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
		Page:       int(pageUint),
		PageSize:   len(workspacesFetched),
		TotalCount: int(count),
		Workspaces: workspacesFetched,
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
	ws.OrganizationId = uuid.FromStringOrNil(payload.OrganizationID)

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
	wsEnvMapping := v1beta1.WorkspacesEnvironmentsMapping{
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

type WorkspaceFilter struct {
	Assigned  bool `json:"assigned"`
	DeletedAt bool `json:"deletedAt"`
}

// GetWorkspaceEnvironments returns environments for a workspace
func (wp *WorkspacePersister) GetWorkspaceEnvironments(workspaceID uuid.UUID, search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = "updated_at desc"
	}
	fmt.Println("got here")
	// Parse the filter parameter
	var workspaceFilter WorkspaceFilter
	workspaceFilter.Assigned = true // Default to true
	if filter != "" {
		err := json.Unmarshal([]byte(filter), &workspaceFilter)
		if err != nil {
			return nil, err
		}
	}

	query := wp.DB.Table("environments AS e").Select("*")

	// Handle filter for whether environments are assigned or not
	if workspaceFilter.Assigned {
		// Environments assigned to the workspace
		query = query.Where("EXISTS (SELECT 1 FROM workspaces_environments_mappings AS wem WHERE e.id = wem.environment_id AND wem.workspace_id = ? AND wem.deleted_at IS NULL)", workspaceID)
	} else {
		// Environments not assigned to the workspace
		query = query.Joins("LEFT JOIN workspaces_environments_mappings AS wem ON e.id = wem.environment_id AND wem.workspace_id = ?", workspaceID).
			Where("wem.workspace_id IS NULL")
	}

	if workspaceFilter.DeletedAt {
		query = query.Where("e.deleted_at IS NOT NULL")
	} else {
		query = query.Where("e.deleted_at IS NULL")
	}

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(e.name) LIKE ? OR lower(e.description) LIKE ?", like, like)
	}

	// Apply additional filters
	dynamicKeys := []string{"owner", "organization_id"}
	query = utils.ApplyFilters(query, filter, dynamicKeys)
	query = query.Order(order)

	count := int64(0)
	query.Count(&count)
	fmt.Println("count done")

	if page == "" {
		page = "0"
	}

	if pageSize == "" {
		pageSize = "10"
	}

	environmentsFetched := []v1beta1.Environment{}
	pageUint, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, err
	}

	if pageSize == "all" {
		query.Find(&environmentsFetched)
	} else {
		pageSizeUint, err := strconv.ParseUint(pageSize, 10, 32)
		if err != nil {
			return nil, err
		}

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
	var wsEnvMapping v1beta1.WorkspacesEnvironmentsMapping

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

func (wp *WorkspacePersister) AddDesignToWorkspace(workspaceID, designID uuid.UUID) ([]byte, error) {
	wsDesignMapping := v1beta1.WorkspacesDesignsMapping{
		DesignId:    designID,
		WorkspaceId: workspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, ErrGenerateUUID(err)
	}
	wsDesignMapping.ID = id

	// Add design to workspace
	if err := wp.DB.Create(wsDesignMapping).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	wsJSON, err := json.Marshal(wsDesignMapping)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) DeleteDesignFromWorkspace(workspaceID, designID uuid.UUID) ([]byte, error) {
	var wsDesignMapping v1beta1.WorkspacesDesignsMapping

	// Find the specific design mapping
	if err := wp.DB.Where("workspace_id = ? AND design_id = ?", workspaceID, designID).First(&wsDesignMapping).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
		return nil, ErrDBRead(err)
	}

	// Delete the design mapping
	if err := wp.DB.Delete(&wsDesignMapping).Error; err != nil {
		return nil, ErrDBDelete(err, wp.fetchUserDetails().UserID)
	}

	wsJSON, err := json.Marshal(wsDesignMapping)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) GetWorkspaceDesigns(workspaceID uuid.UUID, search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = "updated_at desc"
	}

	// Parse the filter parameter
	var workspaceFilter WorkspaceFilter
	workspaceFilter.Assigned = true // Default to true
	if filter != "" {
		err := json.Unmarshal([]byte(filter), &workspaceFilter)
		if err != nil {
			return nil, err
		}
	}

	// Build the query to find designs associated with the given workspace ID
	query := wp.DB.Table("meshery_patterns AS d").Select("*")

	// Build the query to find designs associated with the given workspace ID
	if workspaceFilter.Assigned {
		// Designs assigned to the workspace
		query = query.Where("EXISTS (SELECT 1 FROM workspaces_designs_mappings AS wdm WHERE d.id = wdm.design_id AND wdm.workspace_id = ? AND wdm.deleted_at IS NULL)", workspaceID)
	} else {
		// Designs not assigned to the workspace
		query = query.Joins("LEFT JOIN workspaces_designs_mappings AS wdm ON d.id = wdm.design_id AND wdm.workspace_id = ?", workspaceID).
			Where("wdm.workspace_id IS NULL")
	}

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(d.name) LIKE ?", like)
	}

	dynamicKeys := []string{"owner", "organization_id"}
	query = utils.ApplyFilters(query, filter, dynamicKeys)
	query = query.Order(order)

	count := int64(0)
	query.Count(&count)

	// Handle pagination
	if page == "" {
		page = "0"
	}

	if pageSize == "" {
		pageSize = "10"
	}

	designsFetched := []*MesheryPattern{}
	pageUint, err := strconv.ParseUint(page, 10, 64)
	if err != nil {
		return nil, err
	}

	if pageSize == "all" {
		query.Find(&designsFetched)
	} else {
		pageSizeUint, err := strconv.ParseUint(pageSize, 10, 64)
		if err != nil {
			return nil, err
		}

		Paginate(uint(pageUint), uint(pageSizeUint))(query).Find(&designsFetched)
	}

	designsPage := &MesheryDesignPage{
		Page:       int(pageUint),
		PageSize:   len(designsFetched),
		TotalCount: int(count),
		Designs:    designsFetched,
	}

	designsJSON, err := json.Marshal(designsPage)
	if err != nil {
		return nil, err
	}

	return designsJSON, nil
}
