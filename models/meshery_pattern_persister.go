package models

import (
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
)

// MesheryPatternPersister is the persister for persisting
// performance profiles on the database
type MesheryPatternPersister struct {
	DB *database.Handler
}

// MesheryPatternPage represents a page of performance profiles
type MesheryPatternPage struct {
	Page       uint64            `json:"page"`
	PageSize   uint64            `json:"page_size"`
	TotalCount int               `json:"total_count"`
	Patterns   []*MesheryPattern `json:"patterns"`
}

// GetMesheryPatterns returns all of the performance profiles
func (mpp *MesheryPatternPersister) GetMesheryPatterns(search, order string, page, pageSize uint64) ([]byte, error) {
	order = sanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = "updated_at desc"
	}

	count := int64(0)
	patterns := []*MesheryPattern{}

	query := mpp.DB.Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(meshery_patterns.name) like ?)", like)
	}

	query.Table("meshery_patterns").Count(&count)

	Paginate(uint(page), uint(pageSize))(query).Find(&patterns)

	mesheryPatternPage := &MesheryPatternPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Patterns:   patterns,
	}

	return marshalMesheryPatternPage(mesheryPatternPage), nil
}

// DeleteMesheryPattern takes in a profile id and delete it if it already exists
func (mpp *MesheryPatternPersister) DeleteMesheryPattern(id uuid.UUID) ([]byte, error) {
	pattern := MesheryPattern{ID: &id}
	mpp.DB.Delete(&pattern)

	return marshalMesheryPattern(&pattern), nil
}

// DeleteMesheryPatterns takes in a meshery-patterns and delete those if exist
func (mpp *MesheryPatternPersister) DeleteMesheryPatterns(patterns MesheryPatternDeleteRequestBody) ([]byte, error) {
	var deletedMaptterns []MesheryPattern
	for _, pObj := range patterns.Patterns {
		id := uuid.FromStringOrNil(pObj.ID)
		pattern := MesheryPattern{ID: &id}
		mpp.DB.Delete(&pattern)
		deletedMaptterns = append(deletedMaptterns, pattern)
	}

	return marshalMesheryPatterns(deletedMaptterns), nil
}

func (mpp *MesheryPatternPersister) SaveMesheryPattern(pattern *MesheryPattern) ([]byte, error) {
	if pattern.ID == nil {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}

		pattern.ID = &id
	}

	return marshalMesheryPatterns([]MesheryPattern{*pattern}), mpp.DB.Save(pattern).Error
}

// SaveMesheryPatterns batch inserts the given patterns
func (mpp *MesheryPatternPersister) SaveMesheryPatterns(patterns []MesheryPattern) ([]byte, error) {
	finalPatterns := []MesheryPattern{}
	for _, pattern := range patterns {
		if pattern.ID == nil {
			id, err := uuid.NewV4()
			if err != nil {
				return nil, ErrGenerateUUID(err)
			}

			pattern.ID = &id
		}

		finalPatterns = append(finalPatterns, pattern)
	}

	return marshalMesheryPatterns(finalPatterns), mpp.DB.Create(finalPatterns).Error
}

func (mpp *MesheryPatternPersister) GetMesheryPattern(id uuid.UUID) ([]byte, error) {
	var mesheryPattern MesheryPattern

	err := mpp.DB.First(&mesheryPattern, id).Error
	return marshalMesheryPattern(&mesheryPattern), err
}

func marshalMesheryPatternPage(mpp *MesheryPatternPage) []byte {
	res, _ := json.Marshal(mpp)

	return res
}

func marshalMesheryPattern(mp *MesheryPattern) []byte {
	res, _ := json.Marshal(mp)

	return res
}

func marshalMesheryPatterns(mps []MesheryPattern) []byte {
	res, _ := json.Marshal(mps)

	return res
}
