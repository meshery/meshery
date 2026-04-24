package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/meshery/schemas/models/core"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/v1beta1/environment"
	// NOTE: workspace_persister stays on v1beta1/pattern for the
	// MesheryPattern nested type because schemas' v1beta1/workspace and
	// v1beta3/workspace both type their Designs field against
	// v1beta1/pattern.MesheryPattern (the workspace resource itself is
	// unchanged across versions). Migrating the type graph here would
	// require an upstream schemas change to retype workspace.Designs.
	patternv1beta1 "github.com/meshery/schemas/models/v1beta1/pattern"
	viewv1beta2 "github.com/meshery/schemas/models/v1beta2/view"
	"github.com/meshery/schemas/models/v1beta1/workspace"
	"gorm.io/gorm"
)

// WorkspacePersister is the persister for persisting
// workspaces on the database
type WorkspacePersister struct {
	DB *database.Handler
}

// uuidPtr returns a pointer to the given core.Uuid. Schemas generated
// from meshery/schemas v1.0.9+ use `*core.Uuid` for optional UUID
// fields, so callers building those struct literals need a pointer to
// a local copy. Extracting this out of the hot loop in GetWorkspaces
// keeps the struct literal readable.
func uuidPtr(u core.Uuid) *core.Uuid {
	return &u
}

func (wp *WorkspacePersister) fetchUserDetails() *User {
	return &User{
		UserId:    "meshery",
		FirstName: "Meshery",
		LastName:  "Meshery",
	}
}

// GetWorkspaces returns all of the workspaces
func (wp *WorkspacePersister) GetWorkspaces(orgID, search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	query := wp.DB.Model(&workspace.Workspace{})

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

	workspacesFetched := []workspace.Workspace{}
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

	// Convert fetched workspaces to AvailableWorkspace type
	availableWorkspaces := make([]workspace.AvailableWorkspace, 0, len(workspacesFetched))
	for _, ws := range workspacesFetched {
		aw := workspace.AvailableWorkspace{
			CreatedAt:      ws.CreatedAt,
			DeletedAt:      ws.DeletedAt,
			Description:    ws.Description,
			ID:             ws.ID,
			Name:           ws.Name,
			OrganizationId: uuidPtr(ws.OrganizationID),
			OwnerId:        ws.Owner,
			UpdatedAt:      ws.UpdatedAt,
		}
		availableWorkspaces = append(availableWorkspaces, aw)
	}

	// Prepare the response
	workspacesPage := &workspace.WorkspacePage{
		Page:       int(pageUint),
		PageSize:   len(workspacesFetched),
		TotalCount: int(count),
		Workspaces: availableWorkspaces,
	}

	// Marshal the response to JSON
	wsJSON, err := json.Marshal(workspacesPage)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) SaveWorkspace(workspace *workspace.Workspace) ([]byte, error) {
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

func (wp *WorkspacePersister) DeleteWorkspace(workspace *workspace.Workspace) ([]byte, error) {
	err := wp.DB.Model(&workspace).Find(&workspace).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
	}
	err = wp.DB.Delete(workspace).Error
	if err != nil {
		return nil, ErrDBDelete(err, wp.fetchUserDetails().UserId)
	}

	// Marshal the workspace to JSON
	wsJSON, err := json.Marshal(workspace)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) UpdateWorkspaceByID(selectedWorkspace *workspace.Workspace) (*workspace.Workspace, error) {
	err := wp.DB.Save(selectedWorkspace).Error
	if err != nil {
		return nil, ErrDBPut(err)
	}

	updatedWorkspace := workspace.Workspace{}
	err = wp.DB.Model(&updatedWorkspace).Where("id = ?", selectedWorkspace.ID).First(&updatedWorkspace).Error
	if err != nil {
		return nil, ErrDBRead(err)
	}

	*selectedWorkspace = updatedWorkspace
	return selectedWorkspace, nil
}

// Get workspace by ID
func (wp *WorkspacePersister) GetWorkspace(id core.Uuid) (*workspace.Workspace, error) {
	workspace := workspace.Workspace{}
	query := wp.DB.Where("id = ?", id)
	err := query.First(&workspace).Error
	return &workspace, err
}

// GetWorkspaceByID returns a single workspace by ID
func (wp *WorkspacePersister) GetWorkspaceByID(workspaceID core.Uuid) ([]byte, error) {
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
func (wp *WorkspacePersister) UpdateWorkspace(workspaceID core.Uuid, payload *workspace.WorkspaceUpdatePayload) (*workspace.Workspace, error) {
	ws, err := wp.GetWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}

	if payload.Name != "" {
		ws.Name = payload.Name
	}
	if payload.Description != "" {
		ws.Description = payload.Description
	}
	organizationID := core.Uuid(payload.OrganizationID)
	if organizationID != uuid.Nil {
		ws.OrganizationID = organizationID
	}

	return wp.UpdateWorkspaceByID(ws)
}

// DeleteWorkspaceByID deletes a single workspace by ID
func (wp *WorkspacePersister) DeleteWorkspaceByID(workspaceID core.Uuid) ([]byte, error) {
	ws, err := wp.GetWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}

	return wp.DeleteWorkspace(ws)
}

// AddEnvironmentToWorkspace adds an environment to a workspace
func (wp *WorkspacePersister) AddEnvironmentToWorkspace(workspaceID, environmentID core.Uuid) ([]byte, error) {
	wsEnvMapping := workspace.WorkspacesEnvironmentsMapping{
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
func (wp *WorkspacePersister) GetWorkspaceEnvironments(workspaceID core.Uuid, search, order, page, pageSize, filter string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = defaultOrderUpdatedAtDesc
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

	if page == "" {
		page = "0"
	}

	if pageSize == "" {
		pageSize = "10"
	}

	environmentsFetched := []environment.Environment{}
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

	environmentsPage := &environment.EnvironmentPage{
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
func (wp *WorkspacePersister) DeleteEnvironmentFromWorkspace(workspaceID, environmentID core.Uuid) ([]byte, error) {
	var wsEnvMapping workspace.WorkspacesEnvironmentsMapping

	// Find the specific environment mapping
	if err := wp.DB.Where("workspace_id = ? AND environment_id = ?", workspaceID, environmentID).First(&wsEnvMapping).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
		return nil, ErrDBRead(err)
	}

	// Delete the environment mapping
	if err := wp.DB.Delete(&wsEnvMapping).Error; err != nil {
		return nil, ErrDBDelete(err, wp.fetchUserDetails().UserId)
	}

	wsJSON, err := json.Marshal(wsEnvMapping)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) AddDesignToWorkspace(workspaceID, designID core.Uuid) ([]byte, error) {

	// delete any existing mapping for the design in the workspace
	_, err := wp.DeleteDesignFromWorkspace(workspaceID, designID)

	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, fmt.Errorf("failed to delete existing design mapping: %w", err)
	}

	wsDesignMapping := workspace.WorkspacesDesignsMapping{
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

func (wp *WorkspacePersister) DeleteDesignFromWorkspace(workspaceID, designID core.Uuid) ([]byte, error) {
	var wsDesignMapping workspace.WorkspacesDesignsMapping

	// Find the specific design mapping
	if err := wp.DB.Where("workspace_id = ? AND design_id = ?", workspaceID, designID).First(&wsDesignMapping).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
		return nil, ErrDBRead(err)
	}

	// Delete the design mapping
	if err := wp.DB.Delete(&wsDesignMapping).Error; err != nil {
		return nil, ErrDBDelete(err, wp.fetchUserDetails().UserId)
	}

	wsJSON, err := json.Marshal(wsDesignMapping)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) GetWorkspaceDesigns(workspaceID core.Uuid, search, order, page, pageSize, filter string, visibility []string) ([]byte, error) {
	// Sanitize the order input
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = defaultOrderUpdatedAtDesc
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

	schemaDesigns, err := schemaMesheryPatterns(designsFetched)
	if err != nil {
		return nil, err
	}

	designsPage := &workspace.MesheryDesignPage{
		Page:       int(pageUint),
		PageSize:   len(designsFetched),
		TotalCount: int(count),
		Designs:    schemaDesigns,
	}

	designsJSON, err := json.Marshal(designsPage)
	if err != nil {
		return nil, err
	}

	return designsJSON, nil
}

func schemaMesheryPatterns(patterns []*MesheryPattern) ([]patternv1beta1.MesheryPattern, error) {
	encoded, err := json.Marshal(patterns)
	if err != nil {
		return nil, err
	}

	decoded := []patternv1beta1.MesheryPattern{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}

func (wp *WorkspacePersister) AddViewToWorkspace(workspaceID, viewID core.Uuid) ([]byte, error) {
	wsViewMapping := workspace.WorkspacesViewsMapping{
		ViewId:      viewID,
		WorkspaceId: workspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, ErrGenerateUUID(err)
	}
	wsViewMapping.ID = id

	if err := wp.DB.Create(wsViewMapping).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	wsJSON, err := json.Marshal(wsViewMapping)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) DeleteViewFromWorkspace(workspaceID, viewID core.Uuid) ([]byte, error) {
	var wsViewMapping workspace.WorkspacesViewsMapping

	if err := wp.DB.Where("workspace_id = ? AND view_id = ?", workspaceID, viewID).First(&wsViewMapping).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
		return nil, ErrDBRead(err)
	}

	if err := wp.DB.Delete(&wsViewMapping).Error; err != nil {
		return nil, ErrDBDelete(err, wp.fetchUserDetails().UserId)
	}

	wsJSON, err := json.Marshal(wsViewMapping)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) GetWorkspaceViews(workspaceID core.Uuid, search, order, page, pageSize, filter string) ([]byte, error) {
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	var workspaceFilter WorkspaceFilter
	workspaceFilter.Assigned = true
	if filter != "" {
		err := json.Unmarshal([]byte(filter), &workspaceFilter)
		if err != nil {
			return nil, err
		}
	}

	query := wp.DB.Table("meshery_views AS v").Select("*")

	if workspaceFilter.Assigned {
		query = query.Where("EXISTS (SELECT 1 FROM workspaces_views_mappings AS wvm WHERE v.id = wvm.view_id AND wvm.workspace_id = ? AND wvm.deleted_at IS NULL)", workspaceID)
	} else {
		query = query.Joins("LEFT JOIN workspaces_views_mappings AS wvm ON v.id = wvm.view_id AND wvm.workspace_id = ?", workspaceID).
			Where("wvm.workspace_id IS NULL")
	}

	if workspaceFilter.DeletedAt {
		query = query.Where("v.deleted_at IS NOT NULL")
	} else {
		query = query.Where("v.deleted_at IS NULL")
	}

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(v.name) LIKE ?", like)
	}

	dynamicKeys := []string{"owner", "organization_id"}
	query = utils.ApplyFilters(query, filter, dynamicKeys)
	query = query.Order(order)

	count := int64(0)
	query.Count(&count)

	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	viewsFetched := []viewv1beta2.MesheryViewWithLocation{}
	pageUint, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, err
	}

	if pageSize == "all" {
		query.Find(&viewsFetched)
	} else {
		pageSizeUint, err := strconv.ParseUint(pageSize, 10, 32)
		if err != nil {
			return nil, err
		}
		Paginate(uint(pageUint), uint(pageSizeUint))(query).Find(&viewsFetched)
	}

	viewsPage := &viewv1beta2.MesheryViewPage{
		Page:       int(pageUint),
		PageSize:   len(viewsFetched),
		TotalCount: int(count),
		Views:      viewsFetched,
	}

	viewsJSON, err := json.Marshal(viewsPage)
	if err != nil {
		return nil, err
	}

	return viewsJSON, nil
}

func (wp *WorkspacePersister) AddTeamToWorkspace(workspaceID, teamID core.Uuid) ([]byte, error) {
	wsTeamMapping := workspace.WorkspacesTeamsMapping{
		TeamId:      teamID,
		WorkspaceId: workspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, ErrGenerateUUID(err)
	}
	wsTeamMapping.ID = id

	if err := wp.DB.Create(wsTeamMapping).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	wsJSON, err := json.Marshal(wsTeamMapping)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) DeleteTeamFromWorkspace(workspaceID, teamID core.Uuid) ([]byte, error) {
	var wsTeamMapping workspace.WorkspacesTeamsMapping

	if err := wp.DB.Where("workspace_id = ? AND team_id = ?", workspaceID, teamID).First(&wsTeamMapping).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrResultNotFound(err)
		}
		return nil, ErrDBRead(err)
	}

	if err := wp.DB.Delete(&wsTeamMapping).Error; err != nil {
		return nil, ErrDBDelete(err, wp.fetchUserDetails().UserId)
	}

	wsJSON, err := json.Marshal(wsTeamMapping)
	if err != nil {
		return nil, err
	}

	return wsJSON, nil
}

func (wp *WorkspacePersister) GetWorkspaceTeams(workspaceID core.Uuid, search, order, page, pageSize, filter string) ([]byte, error) {
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	var workspaceFilter WorkspaceFilter
	workspaceFilter.Assigned = true
	if filter != "" {
		err := json.Unmarshal([]byte(filter), &workspaceFilter)
		if err != nil {
			return nil, err
		}
	}

	query := wp.DB.Table("teams AS t").Select("*")

	if workspaceFilter.Assigned {
		query = query.Where("EXISTS (SELECT 1 FROM workspaces_teams_mappings AS wtm WHERE t.id = wtm.team_id AND wtm.workspace_id = ? AND wtm.deleted_at IS NULL)", workspaceID)
	} else {
		query = query.Joins("LEFT JOIN workspaces_teams_mappings AS wtm ON t.id = wtm.team_id AND wtm.workspace_id = ?", workspaceID).
			Where("wtm.workspace_id IS NULL")
	}

	if workspaceFilter.DeletedAt {
		query = query.Where("t.deleted_at IS NOT NULL")
	} else {
		query = query.Where("t.deleted_at IS NULL")
	}

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(t.name) LIKE ?", like)
	}

	dynamicKeys := []string{"owner", "organization_id"}
	query = utils.ApplyFilters(query, filter, dynamicKeys)
	query = query.Order(order)

	count := int64(0)
	query.Count(&count)

	if page == "" {
		page = "0"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	type Team struct {
		ID        core.Uuid  `json:"id" db:"id"`
		Name      string     `json:"name" db:"name"`
		CreatedAt time.Time  `json:"created_at" db:"created_at"`
		UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	}

	teamsFetched := []Team{}
	pageUint, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return nil, err
	}

	if pageSize == "all" {
		query.Find(&teamsFetched)
	} else {
		pageSizeUint, err := strconv.ParseUint(pageSize, 10, 32)
		if err != nil {
			return nil, err
		}
		Paginate(uint(pageUint), uint(pageSizeUint))(query).Find(&teamsFetched)
	}

	type TeamPage struct {
		Page       int    `json:"page"`
		PageSize   int    `json:"page_size"`
		TotalCount int    `json:"total_count"`
		Teams      []Team `json:"teams"`
	}

	teamsPage := &TeamPage{
		Page:       int(pageUint),
		PageSize:   len(teamsFetched),
		TotalCount: int(count),
		Teams:      teamsFetched,
	}

	teamsJSON, err := json.Marshal(teamsPage)
	if err != nil {
		return nil, err
	}

	return teamsJSON, nil
}
