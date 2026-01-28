package models

import (
	"github.com/meshery/schemas/models/v1beta1/organization"
)

// TODO: Move to schemas
type OrganizationsPage struct {
	Organizations []*organization.Organization `json:"organizations"`
	TotalCount    int                          `json:"total_count"`
	Page          uint64                       `json:"page"`
	PageSize      uint64                       `json:"page_size"`
}
