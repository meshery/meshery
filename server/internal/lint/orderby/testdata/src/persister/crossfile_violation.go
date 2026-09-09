package persister

import "gorm.io/gorm"

// crossFileViolation's Order call is on line 9, the same line the directive in
// crossfile_directive.go occupies. Keying suppression on the line number alone
// would silence this diagnostic - the padding here keeps the lines aligned.
func crossFileViolation(db *gorm.DB, order string) *gorm.DB {
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}
