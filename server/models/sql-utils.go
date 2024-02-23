package models

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
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

func setNewDBInstance() {
	mx.Lock()
	defer mx.Unlock()

	// Initialize Logger instance
	log, err := logger.New("meshery", logger.Options{
		Format: logger.SyslogLogFormat,
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	dbHandler, err = database.New(database.Options{
		Filename: fmt.Sprintf("file:%s/mesherydb.sql?cache=private&mode=rwc&_busy_timeout=10000&_journal_mode=WAL", viper.GetString("USER_DATA_FOLDER")),
		Engine:   database.SQLITE,
		Logger:   log,
	})
	if err != nil {
		logrus.Fatal(err)
	}
}

func GetNewDBInstance() *database.Handler {
	setNewDBInstance()
	return &dbHandler
}

func GetDBInstance() *database.Handler {
	return &dbHandler
}
