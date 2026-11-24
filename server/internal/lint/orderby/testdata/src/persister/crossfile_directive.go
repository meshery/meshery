package persister

import "gorm.io/gorm"

// crossFileDirective carries the only //nolint:orderby in this file, on line 9.
// crossfile_violation.go has an unsanitized call on its own line 9. Suppression
// is keyed by file *and* line, so this directive must not reach it.
func crossFileDirective(db *gorm.DB, order string) *gorm.DB {
	return db.Order(order) //nolint:orderby // fixture: this file, this line only.
}
