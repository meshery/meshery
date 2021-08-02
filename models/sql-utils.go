package models

import (
	"fmt"
	"strings"
)

// sanitizeOrderInput takes in the "order by" query, a validColums
// string slice and returns a sanitized query
//
// it will allow to run order by query only on the columns that are present
// in the validColumns string slice, if any other column is requested in the
// query then it will be IGNORED and an empty query would be returned instead
//
// sanitizeOrderInput also expects the query to be no longer than two words, that is
// the query may look like "updated_at DESC" or "name ASC"
func sanitizeOrderInput(order string, validColumns []string) string {
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
