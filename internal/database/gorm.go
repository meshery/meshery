package database

import (
	sqlite "gorm.io/driver/sqlite"
	gormpkg "gorm.io/gorm"

	"github.com/layer5io/meshsync/pkg/model"
)

const (
	GORM = "gorm"
)

func newGorm(opts Options) (Handler, error) {
	db, err := gormpkg.Open(sqlite.Open(opts.Filename), &gormpkg.Config{})
	if err != nil {
		return Handler{}, ErrDatabaseOpen(err)
	}

	// Migrate the schema
	object := &model.Object{}
	err = db.AutoMigrate(object)
	if err != nil {
		return Handler{}, ErrDatabaseMigrate(err)
	}

	return Handler{
		db,
	}, nil
}
