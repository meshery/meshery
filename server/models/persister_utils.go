package models

import (
	"gorm.io/gorm"
)

// Paginate is the utility for paginating the results
func Paginate(page, pageSize uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}
