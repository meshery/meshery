package models

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode"

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
// the query may look like "updatedAt DESC" or "name ASC"
//
// The `order` query parameter is camelCase on the wire ("?order=updatedAt desc")
// while validColumns are the snake_case database columns the ORDER BY runs
// against, so the requested column is folded to snake_case before matching.
// Callers therefore pass DB columns only, and never need to list a key twice to
// accept both spellings. Legacy snake_case callers keep working because a
// snake_case identifier folds to itself.
//
// The returned column is always one of validColumns and is never echoed back
// from user input, which is what keeps the result safe to interpolate into an
// ORDER BY clause.
func SanitizeOrderInput(order string, validColumns []string) string {
	parsedOrderStr := strings.Split(order, " ")
	if len(parsedOrderStr) != 2 {
		return ""
	}

	inputCol := toSnakeCase(parsedOrderStr[0])
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

// toSnakeCase folds a camelCase wire identifier onto the snake_case database
// column it maps to ("createdAt" -> "created_at", "subType" -> "sub_type").
// An identifier that is already snake_case has no uppercase runes and so folds
// to itself.
func toSnakeCase(s string) string {
	var b strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		if !unicode.IsUpper(r) {
			b.WriteRune(r)
			continue
		}
		// Break only at a lower-to-upper boundary so an acronym run stays
		// intact: "userID" folds to "user_id" rather than "user_i_d".
		if i > 0 && !unicode.IsUpper(runes[i-1]) {
			b.WriteByte('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}

	return b.String()
}

var (
	dbHandler database.Handler
	mx        sync.Mutex
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
}

func GetNewDBInstance() *database.Handler {
	setNewDBInstance()
	return &dbHandler
}

func GetDBInstance() *database.Handler {
	return &dbHandler
}
