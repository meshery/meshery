package models

import (
	"fmt"
	"os"
	"strings"
	"sync"

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

	// Create a map for O(1) lookup instead of O(n) loop
	validColsMap := make(map[string]bool, len(validColumns))
	for _, col := range validColumns {
		validColsMap[col] = true
	}

	if validColsMap[inputCol] {
		if typ == "desc" {
			return fmt.Sprintf("%s desc", inputCol)
		}
		return fmt.Sprintf("%s asc", inputCol)
	}

	return ""
}

var (
	dbHandler     database.Handler
	mx            sync.Mutex
	dbOnce        sync.Once
	isInitialized bool
)

func initDBInstance() {
	mx.Lock()
	defer mx.Unlock()

	if isInitialized {
		return
	}

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

	isInitialized = true
}

func GetNewDBInstance() *database.Handler {
	dbOnce.Do(initDBInstance)
	return &dbHandler
}

func GetDBInstance() *database.Handler {
	dbOnce.Do(initDBInstance)
	return &dbHandler
}
