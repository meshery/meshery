// Package models is a stand-in for github.com/meshery/meshery/server/models,
// carrying only the sanitizer the orderby analyzer trusts. The file name
// matters: sql-utils.go is the one file the rule exempts, and the call below
// pins that exemption.
package models

import "gorm.io/gorm"

func SanitizeOrderInput(order string, validColumns []string) string {
	for _, col := range validColumns {
		if col == order {
			return col + " asc"
		}
	}
	return ""
}

// exemptHere is not reported: sql-utils.go is where an ORDER BY fragment is
// legitimately assembled from unsanitized input.
func exemptHere(db *gorm.DB, order string) *gorm.DB {
	return db.Order(order)
}
