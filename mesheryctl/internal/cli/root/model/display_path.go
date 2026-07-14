package model

import (
	"path/filepath"
)

func modelDisplayPath(path string) string {
	if path == "" {
		return "."
	}
	if filepath.Separator == '\\' {
		return filepath.ToSlash(path)
	}
	return path
}
