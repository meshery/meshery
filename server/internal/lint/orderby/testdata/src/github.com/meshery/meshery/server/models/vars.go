package models

import "gorm.io/gorm"

const defaultOrderUpdatedAtDesc = "updated_at desc"

// notExemptHere proves the exemption is scoped to sql-utils.go rather than to
// the whole models package.
func notExemptHere(db *gorm.DB, order string) *gorm.DB {
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}
