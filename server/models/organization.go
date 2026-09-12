package models

import (
	"encoding/json"

	"github.com/meshery/schemas/models/v1beta2/organization"
)

// OrganizationsPage is the Meshery-local response envelope for
// `GET /api/identity/orgs`. Canonical wire form is camelCase per the
// identifier-naming migration; legacy snake_case keys (`total_count`,
// `page_size`) are emitted alongside for the deprecation window so
// external consumers still reading the old spellings keep working.
// Unmarshal likewise accepts either spelling so any inbound serialization
// (round-trip via clients) tolerates both.
//
// Once every known consumer has migrated off the snake_case keys, drop
// MarshalJSON/UnmarshalJSON and keep only the camelCase struct tags.
//
// TODO: Move to schemas
type OrganizationsPage struct {
	Organizations []*organization.Organization `json:"organizations"`
	TotalCount    int                          `json:"totalCount"`
	Page          uint64                       `json:"page"`
	PageSize      uint64                       `json:"pageSize"`
}

// MarshalJSON emits both camelCase (canonical) and snake_case (legacy)
// keys for TotalCount and PageSize so external consumers reading either
// spelling continue to work during the deprecation window.
func (p OrganizationsPage) MarshalJSON() ([]byte, error) {
	type alias OrganizationsPage
	return json.Marshal(struct {
		alias
		TotalCountLegacy int    `json:"total_count"`
		PageSizeLegacy   uint64 `json:"page_size"`
	}{
		alias:            alias(p),
		TotalCountLegacy: p.TotalCount,
		PageSizeLegacy:   p.PageSize,
	})
}

// UnmarshalJSON accepts either the canonical camelCase keys or the
// legacy snake_case keys. Canonical keys win when both are present.
func (p *OrganizationsPage) UnmarshalJSON(data []byte) error {
	var raw struct {
		Organizations    []*organization.Organization `json:"organizations"`
		TotalCount       *int                         `json:"totalCount"`
		TotalCountLegacy *int                         `json:"total_count"`
		Page             uint64                       `json:"page"`
		PageSize         *uint64                      `json:"pageSize"`
		PageSizeLegacy   *uint64                      `json:"page_size"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Organizations = raw.Organizations
	p.Page = raw.Page
	// Match stdlib json.Unmarshal semantics: fields absent from the input
	// reset to their zero value on the destination (important when callers
	// reuse the same OrganizationsPage across decodes).
	p.TotalCount = 0
	switch {
	case raw.TotalCount != nil:
		p.TotalCount = *raw.TotalCount
	case raw.TotalCountLegacy != nil:
		p.TotalCount = *raw.TotalCountLegacy
	}
	p.PageSize = 0
	switch {
	case raw.PageSize != nil:
		p.PageSize = *raw.PageSize
	case raw.PageSizeLegacy != nil:
		p.PageSize = *raw.PageSizeLegacy
	}
	return nil
}
