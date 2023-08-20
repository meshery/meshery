package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
)

// MesheryFilterPersister is the persister for persisting
// wasm filters on the database
type MesheryFilterPersister struct {
	DB *database.Handler
}

// MesheryFilterPage represents a page of filters
type MesheryFilterPage struct {
	Page       uint64           `json:"page"`
	PageSize   uint64           `json:"page_size"`
	TotalCount int              `json:"total_count"`
	Filters    []*MesheryFilter `json:"filters"`
}

// GetMesheryFilters returns all of the 'private' filters. Though private has no meaning here since there is only
// one local user. We make this distinction to be consistent with the remote provider
func (mfp *MesheryFilterPersister) GetMesheryFilters(search, order string, page, pageSize uint64, visibility string) ([]byte, error) {
	order = sanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = "updated_at desc"
	}

	if visibility == "" {
		visibility = "private"
	}

	count := int64(0)
	filters := []*MesheryFilter{}

	query := mfp.DB.Where("visibility = ?", visibility).Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(meshery_filters.name) like ?)", like)
	}

	query.Table("meshery_filters").Count(&count)

	Paginate(uint(page), uint(pageSize))(query).Find(&filters)

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
	order = sanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = "updated_at desc"
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
	filterUUID, _ := uuid.FromString(filterID)
	err := mfp.DB.First(&mesheryFilter, filterUUID).Error
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
func (mfp *MesheryFilterPersister) DeleteMesheryFilter(id uuid.UUID) ([]byte, error) {
	filter := MesheryFilter{ID: &id}
	mfp.DB.Delete(&filter)

	return marshalMesheryFilter(&filter), nil
}

func (mfp *MesheryFilterPersister) SaveMesheryFilter(filter *MesheryFilter) ([]byte, error) {
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

	return marshalMesheryFilters([]MesheryFilter{*filter}), mfp.DB.Save(filter).Error
}

// SaveMesheryFilters batch inserts the given filters
func (mfp *MesheryFilterPersister) SaveMesheryFilters(filters []MesheryFilter) ([]byte, error) {
	finalFilters := []MesheryFilter{}
	nilUserID := ""
	for _, filter := range filters {
		if filter.Visibility == "" {
			filter.Visibility = Private
		}
		filter.UserID = &nilUserID
		if filter.ID == nil {
			id, err := uuid.NewV4()
			if err != nil {
				return nil, ErrGenerateUUID(err)
			}

			filter.ID = &id
		}

		finalFilters = append(finalFilters, filter)
	}

	return marshalMesheryFilters(finalFilters), mfp.DB.Create(finalFilters).Error
}

func (mfp *MesheryFilterPersister) GetMesheryFilter(id uuid.UUID) ([]byte, error) {
	var mesheryFilter MesheryFilter

	err := mfp.DB.First(&mesheryFilter, id).Error
	return marshalMesheryFilter(&mesheryFilter), err
}

func (mfp *MesheryFilterPersister) GetMesheryFilterFile(id uuid.UUID) ([]byte, error) {
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
