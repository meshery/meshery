package models

import (
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
)

type PatternResourcePersister struct {
	DB *database.Handler
}

type PatternResourcePage struct {
	Page       uint64             `json:"page,omitempty"`
	PageSize   uint64             `json:"page_size,omitempty"`
	TotalCount int                `json:"total_count,omitempty"`
	Resources  []*PatternResource `json:"resources,omitempty"`
}

func (prp *PatternResourcePersister) SavePatternResource(pr *PatternResource) (*PatternResource, error) {
	if pr.ID == nil {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}

		pr.ID = &id
	}

	return pr, prp.DB.Create(pr).Error
}

func (prp *PatternResourcePersister) DeletePatternResource(id uuid.UUID) error {
	return prp.DB.Model(&PatternResource{}).Where("id = ?", id).Update("deleted", true).Error
}

func (prp *PatternResourcePersister) GetPatternResource(id uuid.UUID) (*PatternResource, error) {
	var pr PatternResource

	err := prp.DB.First(&pr, id).Error
	return &pr, err
}

func (prp *PatternResourcePersister) GetPatternResourceByAttributes(name, namespace, typ, oamType string) (*PatternResource, error) {
	var pr PatternResource

	err := Paginate(0, 1)(prp.DB.Model(&PatternResource{}).Where(
		"name = ? AND namespace = ? AND type = ? AND oam_type = ? AND deleted = false",
		name,
		namespace,
		typ,
		oamType,
	)).Scan(&pr).Error

	return &pr, err
}

func (prp *PatternResourcePersister) GetPatternResources(search, order, name, namespace, typ, oamType string, page, pageSize uint64) (*PatternResourcePage, error) {
	order = sanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = "updated_at desc"
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

	query.Model(&PatternResource{}).Count(&count)

	Paginate(uint(page), uint(pageSize))(query).Find(&resources)

	patternResourcePage := &PatternResourcePage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Resources:  resources,
	}

	return patternResourcePage, nil
}

func (prp *PatternResourcePersister) Exists(name, namespace, typ, oamType string) bool {
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
