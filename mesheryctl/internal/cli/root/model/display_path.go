package model

import (
	"path/filepath"
	"strings"
)

func modelDisplayPath(path string) string {
	if path == "" {
		return "."
	}
	normalizedPath := strings.ReplaceAll(path, `\`, "/")
	return strings.ReplaceAll(filepath.Clean(normalizedPath), `\`, "/")
}
