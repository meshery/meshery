package model

import "strings"

func modelDisplayPath(path string) string {
	if path == "" {
		return "."
	}
	// Normalize displayed paths to forward slashes (filesystem behavior is unchanged).
	return strings.ReplaceAll(path, "\\", "/")
}
