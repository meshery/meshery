// Package gorm is a stand-in for gorm.io/gorm, reproducing only the shape the
// orderby analyzer keys on: a concrete *DB whose Order method takes `any`.
// analysistest loads its fixtures in GOPATH mode with the module proxy off, so
// the real dependency is not reachable from here.
package gorm

type DB struct{}

func (db *DB) Order(value any) *DB { return db }

func (db *DB) Where(query any, args ...any) *DB { return db }

func (db *DB) Find(dest any) *DB { return db }
