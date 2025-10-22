package models

import (
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"gorm.io/gorm"
)

// BasePersister provides common database operations for all persisters
type BasePersister struct {
	DB *database.Handler
}

// PaginationResult represents a paginated result
type PaginationResult struct {
	Page       uint64      `json:"page"`
	PageSize   uint64      `json:"page_size"`
	TotalCount int         `json:"total_count"`
	Data       interface{} `json:"data"`
}

// QueryBuilder helps build common database queries
type QueryBuilder struct {
	query *gorm.DB
}

// NewQueryBuilder creates a new query builder
func (bp *BasePersister) NewQueryBuilder(model interface{}) *QueryBuilder {
	return &QueryBuilder{
		query: bp.DB.Model(model),
	}
}

// WithSearch adds search functionality to the query
func (qb *QueryBuilder) WithSearch(search string, searchFields []string) *QueryBuilder {
	if search != "" && len(searchFields) > 0 {
		like := "%" + strings.ToLower(search) + "%"
		conditions := make([]string, len(searchFields))
		args := make([]interface{}, len(searchFields))

		for i, field := range searchFields {
			conditions[i] = "lower(" + field + ") like ?"
			args[i] = like
		}

		qb.query = qb.query.Where("("+strings.Join(conditions, " OR ")+")", args...)
	}
	return qb
}

// WithOrder adds ordering to the query with sanitization
func (qb *QueryBuilder) WithOrder(order string, validColumns []string) *QueryBuilder {
	sanitizedOrder := SanitizeOrderInput(order, validColumns)
	if sanitizedOrder == "" {
		sanitizedOrder = "updated_at desc"
	}
	qb.query = qb.query.Order(sanitizedOrder)
	return qb
}

// WithUpdatedAfter filters by updated_at timestamp
func (qb *QueryBuilder) WithUpdatedAfter(updatedAfter string) *QueryBuilder {
	if updatedAfter != "" {
		qb.query = qb.query.Where("updated_at > ?", updatedAfter)
	}
	return qb
}

// WithFilters adds multiple filter conditions
func (qb *QueryBuilder) WithFilters(filters map[string]interface{}) *QueryBuilder {
	for field, value := range filters {
		switch v := value.(type) {
		case []string:
			if len(v) > 0 {
				qb.query = qb.query.Where(field+" IN (?)", v)
			}
		case string:
			if v != "" {
				qb.query = qb.query.Where(field+" = ?", v)
			}
		default:
			qb.query = qb.query.Where(field+" = ?", v)
		}
	}
	return qb
}

// Paginate applies pagination to the query
func (qb *QueryBuilder) Paginate(page, pageSize uint64) *QueryBuilder {
	if pageSize > 0 {
		qb.query = Paginate(uint(page), uint(pageSize))(qb.query)
	}
	return qb
}

// Count returns the total count of records
func (qb *QueryBuilder) Count(tableName string) int64 {
	var count int64
	qb.query.Table(tableName).Count(&count)
	return count
}

// Find executes the query and populates the result
func (qb *QueryBuilder) Find(dest interface{}) error {
	return qb.query.Find(dest).Error
}

// GetQuery returns the underlying GORM query
func (qb *QueryBuilder) GetQuery() *gorm.DB {
	return qb.query
}

// GenerateUUID generates a new UUID for entities
func (bp *BasePersister) GenerateUUID() (uuid.UUID, error) {
	return uuid.NewV4()
}

// SaveEntity saves an entity with automatic UUID generation
func (bp *BasePersister) SaveEntity(entity interface{}, idField *uuid.UUID) error {
	if idField != nil && *idField == uuid.Nil {
		id, err := bp.GenerateUUID()
		if err != nil {
			return ErrGenerateUUID(err)
		}
		*idField = id
	}

	return bp.DB.Save(entity).Error
}

// DeleteEntity performs soft delete by setting deleted flag
func (bp *BasePersister) DeleteEntity(model interface{}, id uuid.UUID) error {
	return bp.DB.Model(model).Where("id = ?", id).Update("deleted", true).Error
}

// HardDeleteEntity performs hard delete
func (bp *BasePersister) HardDeleteEntity(entity interface{}) error {
	return bp.DB.Delete(entity).Error
}

// FindByID finds an entity by ID
func (bp *BasePersister) FindByID(dest interface{}, id uuid.UUID) error {
	return bp.DB.First(dest, id).Error
}
