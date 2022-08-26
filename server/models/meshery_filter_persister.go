package models

import (
	"encoding/json"
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

// GetMesheryFilters returns all of the filters
func (mfp *MesheryFilterPersister) GetMesheryFilters(search, order string, page, pageSize uint64) ([]byte, error) {
	order = sanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = "updated_at desc"
	}

	count := int64(0)
	filters := []*MesheryFilter{}

	query := mfp.DB.Order(order)

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

// DeleteMesheryFilter takes in a profile id and delete it if it already exists
func (mfp *MesheryFilterPersister) DeleteMesheryFilter(id uuid.UUID) ([]byte, error) {
	filter := MesheryFilter{ID: &id}
	mfp.DB.Delete(&filter)

	return marshalMesheryFilter(&filter), nil
}

func (mfp *MesheryFilterPersister) SaveMesheryFilter(filter *MesheryFilter) ([]byte, error) {
	filter.Visibility = "private"
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
	for _, filter := range filters {
		filter.Visibility = "private"
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
