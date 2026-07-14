package model

import (
	"strings"
)

func modelDisplayPath(path string) string {
	if path == "" {
		return "."
	}
	return strings.ReplaceAll(path, `\`, "/")
}
