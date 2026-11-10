package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/meshery/schemas/models/core"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"gorm.io/gorm/clause"
)

// MesheryFilterPersister is the persister for persisting
// wasm filters on the database
type MesheryFilterPersister struct {
	DB *database.Handler
}

// stampLocalProviderFilterOwner sets the built-in provider's single user as the
// filter owner. MesheryFilter.Owner is gorm:"-" (never a DB column) and the
// local provider is single-user, so every filter it persists is owned by the
// local "meshery" user. Previously the persister wrote an empty-string owner on
// save, which is unresolvable: the UI had no id to look a profile up with and
// the owner never rendered.
//
// Published filters are left untouched for the same reason published designs
// are - see stampLocalProviderOwner.
func stampLocalProviderFilterOwner(f *MesheryFilter) {
	if f == nil || !LocalProviderOwnsContent(f.Visibility) {
		return
	}
	owner := LocalProviderUserID.String()
	f.Owner = &owner
}

// MesheryFilterPage represents a page of filters
type MesheryFilterPage struct {
	Page       uint64           `json:"page"`
	PageSize   uint64           `json:"pageSize"`
	TotalCount int              `json:"totalCount"`
	Filters    []*MesheryFilter `json:"filters"`
}

// GetMesheryFilters returns all of the 'private' filters. Though private has no meaning here since there is only
// one local user. We make this distinction to be consistent with the remote provider
func (mfp *MesheryFilterPersister) GetMesheryFilters(search, order string, page, pageSize uint64, visibility []string) ([]byte, error) {
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	count := int64(0)
	filters := []*MesheryFilter{}

	query := mfp.DB.Table("meshery_filters")

	if len(visibility) > 0 {
		query = query.Where("visibility in (?)", visibility)
	}

	query = query.Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(meshery_filters.name) like ?)", like)
	}

	query.Count(&count)
	Paginate(uint(page), uint(pageSize))(query).Find(&filters)

	for _, f := range filters {
		stampLocalProviderFilterOwner(f)
	}

	mesheryFilterPage := &MesheryFilterPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Filters:    filters,
	}

	return marshalMesheryFilterPage(mesheryFilterPage), nil
}

// GetMesheryCatalogFilters returns all of the published filters
func (mfp *MesheryFilterPersister) GetMesheryCatalogFilters(page, pageSize, search, order string) ([]byte, error) {
	var err error
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	var pg int
	if page != "" {
		pg, err = strconv.Atoi(page)

		if err != nil || pg < 0 {
			pg = 0
		}
	} else {
		pg = 0
	}

	// 0 page size is for all records
	var pgSize int
	if pageSize != "" {
		pgSize, err = strconv.Atoi(pageSize)

		if err != nil || pgSize < 0 {
			pgSize = 0
		}
	} else {
		pgSize = 0
	}

	filters := []MesheryFilter{}

	query := mfp.DB.Where("visibility = ?", Published).Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(meshery_filters.name) like ?)", like)
	}

	var count int64
	err = query.Model(&MesheryFilter{}).Count(&count).Error

	if err != nil {
		return nil, err
	}

	if pgSize != 0 {
		Paginate(uint(pg), uint(pgSize))(query).Find(&filters)
	} else {
		query.Find(&filters)
	}

	response := FiltersAPIResponse{
		Page:       uint(pg),
		PageSize:   uint(pgSize),
		TotalCount: uint(count),
		Filters:    filters,
	}

	marshalledResponse, _ := json.Marshal(response)
	return marshalledResponse, nil
}

// CloneMesheryFilter clones meshery filter to private
func (mfp *MesheryFilterPersister) CloneMesheryFilter(filterID string, cloneFilterRequest *MesheryCloneFilterRequestBody) ([]byte, error) {
	var mesheryFilter MesheryFilter
	filterUUID, err := uuid.FromString(filterID)
	if err != nil {
		return nil, ErrInvalidUUID(err)
	}
	err = mfp.DB.First(&mesheryFilter, filterUUID).Error
	if err != nil || *mesheryFilter.ID == uuid.Nil {
		return nil, fmt.Errorf("unable to get filter: %w", err)
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	mesheryFilter.Visibility = Private
	mesheryFilter.ID = &id
	mesheryFilter.Name = cloneFilterRequest.Name

	return mfp.SaveMesheryFilter(&mesheryFilter)
}

// DeleteMesheryFilter takes in a profile id and delete it if it already exists
func (mfp *MesheryFilterPersister) DeleteMesheryFilter(id core.Uuid) ([]byte, error) {
	filter := MesheryFilter{ID: &id}
	res := mfp.DB.Delete(&filter)
	if res.Error != nil {
		return nil, ErrDBDelete(res.Error, id.String())
	}
	if res.RowsAffected == 0 {
		return nil, ErrResultNotFound(fmt.Errorf("filter %s not found", id.String()))
	}

	return marshalMesheryFilter(&filter), nil
}

// SaveMesheryFilter inserts a new filter, or updates an existing one with only
// the fields the save carries.
func (mfp *MesheryFilterPersister) SaveMesheryFilter(filter *MesheryFilter) ([]byte, error) {
	isNew := filter.ID == nil
	if isNew {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}

		filter.ID = &id
	}

	// Stamp before marshalling so the save/clone response carries the same
	// owner contract as the GET paths. Owner is gorm:"-", so this is response
	// shaping only - nothing extra is persisted.
	stampLocalProviderFilterOwner(filter)

	// An ID does not mean the row exists: clone and callers that choose their
	// own ID save new rows under it.
	exists := false
	if !isNew {
		var count int64
		if err := mfp.DB.Model(&MesheryFilter{}).Where("id = ?", *filter.ID).Count(&count).Error; err != nil {
			return nil, ErrDBRead(err)
		}
		exists = count > 0
	}

	var err error
	updated := exists
	if exists {
		err = mfp.updateSavedColumns(filter)
	} else {
		suppliedVisibility := filter.Visibility
		if filter.Visibility == "" {
			filter.Visibility = Private
		}
		res := mfp.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(filter)
		if res.Error != nil {
			return marshalMesheryFilters([]MesheryFilter{*filter}), ErrDBCreate(res.Error)
		}
		if res.RowsAffected == 0 {
			// A concurrent save inserted this ID first: update it like any
			// existing row, without the insert-only visibility default.
			filter.Visibility = suppliedVisibility
			updated = true
			err = mfp.updateSavedColumns(filter)
		}
	}

	// The saved struct holds only what the request carried; respond with the
	// stored row so the response shows the columns the update preserved.
	if updated && err == nil {
		var stored MesheryFilter
		if err = mfp.DB.First(&stored, "id = ?", *filter.ID).Error; err != nil {
			err = ErrDBRead(err)
		} else {
			stampLocalProviderFilterOwner(&stored)
			filter = &stored
		}
	}

	return marshalMesheryFilters([]MesheryFilter{*filter}), err
}

// updateSavedColumns updates an existing filter with the columns the save
// carries. POST /api/filter builds the record from the request, and a rename
// from the UI carries only the ID, the name and empty catalog data, so Updates
// skips zero-valued fields and those columns keep their stored values.
// created_at is never written by a save, and an empty filter file or catalog
// data is treated as not sent.
func (mfp *MesheryFilterPersister) updateSavedColumns(filter *MesheryFilter) error {
	omit := []string{"created_at"}
	if len(filter.FilterFile) == 0 {
		omit = append(omit, "filter_file")
	}
	if len(filter.CatalogData) == 0 {
		omit = append(omit, "catalog_data")
	}
	res := mfp.DB.Model(filter).Omit(omit...).Updates(filter)
	// updated_at always changes, so no affected row means the filter was
	// deleted after the existence check.
	if res.Error == nil && res.RowsAffected == 0 {
		res.Error = fmt.Errorf("filter %s was deleted while saving", filter.ID)
	}
	if res.Error != nil {
		return ErrDBPut(res.Error)
	}
	return nil
}

// SaveMesheryFilters batch inserts the given filters
func (mfp *MesheryFilterPersister) SaveMesheryFilters(filters []MesheryFilter) ([]byte, error) {
	finalFilters := []MesheryFilter{}
	for _, filter := range filters {
		if filter.Visibility == "" {
			filter.Visibility = Private
		}
		if filter.ID == nil {
			id, err := uuid.NewV4()
			if err != nil {
				return nil, ErrGenerateUUID(err)
			}

			filter.ID = &id
		}

		stampLocalProviderFilterOwner(&filter)

		finalFilters = append(finalFilters, filter)
	}

	return marshalMesheryFilters(finalFilters), mfp.DB.Create(finalFilters).Error
}

func (mfp *MesheryFilterPersister) GetMesheryFilter(id core.Uuid) ([]byte, error) {
	var mesheryFilter MesheryFilter

	err := mfp.DB.First(&mesheryFilter, id).Error
	stampLocalProviderFilterOwner(&mesheryFilter)
	return marshalMesheryFilter(&mesheryFilter), err
}

func (mfp *MesheryFilterPersister) GetMesheryFilterFile(id core.Uuid) ([]byte, error) {
	var mesheryFilter MesheryFilter

	err := mfp.DB.First(&mesheryFilter, id).Error
	return []byte(mesheryFilter.FilterFile), err
}

func marshalMesheryFilterPage(mfp *MesheryFilterPage) []byte {
	res, _ := json.Marshal(mfp)

	return res
}

func marshalMesheryFilter(mf *MesheryFilter) []byte {
	res, _ := json.Marshal(mf)

	return res
}

func marshalMesheryFilters(mps []MesheryFilter) []byte {
	res, _ := json.Marshal(mps)

	return res
}
