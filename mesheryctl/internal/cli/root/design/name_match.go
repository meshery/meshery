package design

import (
	"strings"

	"github.com/meshery/meshery/server/models"
)

// exactNameMatches keeps the designs whose name is exactly name. The server's
// search parameter is a substring match, so its results can include designs
// that only contain name. Case-sensitive matches win; otherwise fall back to
// case-insensitive ones.
func exactNameMatches(patterns []models.MesheryPattern, name string) []models.MesheryPattern {
	var exact, fold []models.MesheryPattern
	for _, p := range patterns {
		switch {
		case p.Name == name:
			exact = append(exact, p)
		case strings.EqualFold(p.Name, name):
			fold = append(fold, p)
		}
	}
	if len(exact) > 0 {
		return exact
	}
	return fold
}
