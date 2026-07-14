package model

import "path/filepath"

func modelDisplayPath(path string) string {
	if path == "" {
		return "."
	}
	// Normalize displayed paths to forward slashes (filesystem behavior is unchanged).
	return filepath.ToSlash(path)
}
