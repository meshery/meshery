package models

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SanitizeOrderInput takes in the "order by" query, a validColums
// string slice and returns a sanitized query
//
// it will allow to run order by query only on the columns that are present
// in the validColumns string slice, if any other column is requested in the
// query then it will be IGNORED and an empty query would be returned instead
//
// SanitizeOrderInput also expects the query to be no longer than two words, that is
// the query may look like "updated_at DESC" or "name ASC"
func SanitizeOrderInput(order string, validColumns []string) string {
	parsedOrderStr := strings.Split(order, " ")
	if len(parsedOrderStr) != 2 {
		return ""
	}

	inputCol := parsedOrderStr[0]
	typ := strings.ToLower(parsedOrderStr[1])
	for _, col := range validColumns {
		if col == inputCol {
			if typ == "desc" {
				return fmt.Sprintf("%s desc", col)
			}

			return fmt.Sprintf("%s asc", col)
		}
	}

	return ""
}

var (
	dbHandler database.Handler
	mx        sync.Mutex
)

const (
	// maxOpenConns sets the maximum number of open connections to the database.
	// Setting this to 1 is crucial for SQLite to prevent "database is locked" errors.
	maxOpenConns = 1
	// maxIdleConns sets the maximum number of connections in the idle connection pool.
	// ensures that the one connection is retained in the pool after use
	maxIdleConns = 1
	// connMaxLifetime sets the maximum amount of time a connection may be reused.
	// A value of 0 means that connections are not closed due to a connection's age.
	connMaxLifetime time.Duration = 0
)

func setNewDBInstance() {
	mx.Lock()
	defer mx.Unlock()

	// Initialize Logger instance
	logLevel := viper.GetInt("LOG_LEVEL")
	log, err := logger.New("meshery", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: logLevel,
	})
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	dbHandler, err = database.New(database.Options{
		Filename: fmt.Sprintf("file:%s/mesherydb.sql?cache=private&mode=rwc&_busy_timeout=10000&_journal_mode=WAL", viper.GetString("USER_DATA_FOLDER")),
		Engine:   database.SQLITE,
		Logger:   log,
	})
	if err != nil {
		err = ErrInitializeDBHandler(err)
		log.Error(err)
	}
	sqlDB, err := dbHandler.DB.DB()
	if err != nil {
		log.Error(err)
		return
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
}

func GetNewDBInstance() *database.Handler {
	setNewDBInstance()
	return &dbHandler
}

func GetDBInstance() *database.Handler {
	return &dbHandler
}
