package models

import (
	"fmt"
	"strings"
)

// sanitizeOrderInput takes in the "order by" string and a validColums
// string slice and returns a sanitized the string
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
