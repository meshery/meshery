package models

import (
	"strings"

	"errors"
	"gorm.io/gorm"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
)

type PatternResourcePersister struct {
	DB *database.Handler
}

type PatternResourcePage struct {
	Page       uint64             `json:"page,omitempty"`
	PageSize   uint64             `json:"pageSize,omitempty"`
	TotalCount int                `json:"totalCount,omitempty"`
	Resources  []*PatternResource `json:"resources,omitempty"`
}

func (prp *PatternResourcePersister) SavePatternResource(pr *PatternResource) (*PatternResource, error) {
	if prp == nil || prp.DB == nil || prp.DB.DB == nil {
		return nil, ErrDBConnection
	}
	if pr == nil {
		return nil, ErrNilConfigData
	}
	if pr.ID == nil {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}

		pr.ID = &id
	}

	return pr, prp.DB.Create(pr).Error
}

func (prp *PatternResourcePersister) DeletePatternResource(id core.Uuid) error {
	if prp == nil || prp.DB == nil || prp.DB.DB == nil {
		return ErrDBConnection
	}
	result := prp.DB.Model(&PatternResource{}).Where("id = ? AND deleted = false", id).Update("deleted", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrResultNotFound(gorm.ErrRecordNotFound)
	}
	return nil
}

func (prp *PatternResourcePersister) GetPatternResource(id core.Uuid) (*PatternResource, error) {
	if prp == nil || prp.DB == nil || prp.DB.DB == nil {
		return nil, ErrDBConnection
	}
	var pr PatternResource
	err := prp.DB.Where("deleted = false").First(&pr, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrResultNotFound(err)
	}
	if err != nil {
		return nil, ErrDBRead(err)
	}
	return &pr, nil
}

func (prp *PatternResourcePersister) GetPatternResourceByAttributes(name, namespace, typ, oamType string) (*PatternResource, error) {
	if prp == nil || prp.DB == nil || prp.DB.DB == nil {
		return nil, ErrDBConnection
	}
	var pr PatternResource
	err := prp.DB.Where("name = ? AND namespace = ? AND type = ? AND oam_type = ? AND deleted = false", name, namespace, typ, oamType).First(&pr).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrResultNotFound(err)
	}
	if err != nil {
		return nil, ErrDBRead(err)
	}
	return &pr, nil
}

func (prp *PatternResourcePersister) GetPatternResources(search, order, name, namespace, typ, oamType string, page, pageSize uint64) (*PatternResourcePage, error) {
	if prp == nil || prp.DB == nil || prp.DB.DB == nil {
		return nil, ErrDBConnection
	}
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	count := int64(0)
	resources := []*PatternResource{}

	query := prp.DB.Order(order).Where("deleted = false")

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(pattern_resources.name) like ?)", like)
	}
	if name != "" {
		query = query.Where("name = ?", name)
	}
	if namespace != "" {
		query = query.Where("namespace = ?", namespace)
	}
	if typ != "" {
		query = query.Where("type = ?", typ)
	}
	if oamType != "" {
		query = query.Where("oam_type = ?", oamType)
	}

	if err := query.Model(&PatternResource{}).Count(&count).Error; err != nil {
		return nil, ErrDBRead(err)
	}

	if err := Paginate(uint(page), uint(pageSize))(query).Find(&resources).Error; err != nil {
		return nil, ErrDBRead(err)
	}

	patternResourcePage := &PatternResourcePage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Resources:  resources,
	}

	return patternResourcePage, nil
}

func (prp *PatternResourcePersister) Exists(name, namespace, typ, oamType string) bool {
	if prp == nil || prp.DB == nil || prp.DB.DB == nil {
		return false
	}
	var result struct {
		Found bool
	}

	prp.DB.
		Raw(`
	SELECT
		EXISTS(SELECT 1 FROM pattern_resources WHERE name = ? AND namespace = ? AND type = ? AND oam_type = ? AND deleted = false) AS "found"`,
			name,
			namespace,
			typ,
			oamType,
		).
		Scan(&result)

	return result.Found
}
