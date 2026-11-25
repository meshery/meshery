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

// fetchAllPatternPages collects every page of a design name search. The server
// paginates the substring search (10 designs per page by default) and orders by
// updated_at, so a single request can omit the exact design when many stored
// names contain the search term. fetchPage is called with increasing page
// numbers starting at 0 until every design counted by totalCount has been
// gathered or a page comes back empty.
func fetchAllPatternPages(fetchPage func(page uint) (*models.PatternsAPIResponse, error)) ([]models.MesheryPattern, error) {
	var all []models.MesheryPattern
	for page := uint(0); ; page++ {
		resp, err := fetchPage(page)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			break
		}
		all = append(all, resp.Patterns...)
		if len(resp.Patterns) == 0 || uint(len(all)) >= resp.TotalCount {
			break
		}
	}
	return all, nil
}
