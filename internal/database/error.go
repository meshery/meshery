package database

import "github.com/layer5io/meshkit/errors"

var (
	ErrNoneDatabaseCode    = "test"
	ErrDatabaseOpenCode    = "test"
	ErrDatabaseMigrateCode = "test"

	ErrNoneDatabase = errors.NewDefault(ErrNoneDatabaseCode, "No Database selected")
)

func ErrDatabaseOpen(err error) error {
	return errors.NewDefault(ErrDatabaseOpenCode, "Unable to open database", err.Error())
}

func ErrDatabaseMigrate(err error) error {
	return errors.NewDefault(ErrDatabaseMigrateCode, "Unable to migrate database schema", err.Error())
}
