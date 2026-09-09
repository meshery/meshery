// Package clause is a stand-in for gorm.io/gorm/clause, carrying only the
// ORDER BY values gorm renders through its quoting clause builder.
package clause

type Column struct {
	Name string
	Raw  bool
}

type OrderByColumn struct {
	Column Column
	Desc   bool
}
