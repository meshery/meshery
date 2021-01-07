package database

import (
	gormpkg "gorm.io/gorm"
)

type Options struct {
	Filename string `json:"filename"`
}

type Handler struct {
	*gormpkg.DB
	// Implement methods if necessary
}

func New(engine string, opts Options) (Handler, error) {
	switch engine {
	case GORM:
		return newGorm(opts)
	}

	return Handler{}, ErrNoneDatabase
}
