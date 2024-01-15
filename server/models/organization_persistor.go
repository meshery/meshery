package models

import (
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
)

// MesheryApplicationPersister is the persister for persisting
// applications on the database
type OrganizationPersister struct {
	DB *database.Handler
}

// GetMesheryApplications returns all of the applications
func (op *OrganizationPersister) GetOrganizations(search, order string, page, pageSize uint64, updatedAfter string) ([]byte, error) {
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = "updated_at desc"
	}

	count := int64(0)
	organizations := []*Organization{}

	query := op.DB.Where("updated_at > ?", updatedAfter).Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(organizations.name) like ?)", like)
	}

	query.Table("organizations").Count(&count)

	Paginate(uint(page), uint(pageSize))(query).Find(&organizations)

	organizationsPage := &OrganizationsPage{
		Page:          page,
		PageSize:      pageSize,
		TotalCount:    int(count),
		Organizations: organizations,
	}

	return marshalOrganizationsPage(organizationsPage), nil
}

func (op *OrganizationPersister) SaveOrganization(organization *Organization) ([]byte, error) {
	if organization.ID == nil {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}

		organization.ID = &id
	}

	return marshalOrganizations([]Organization{*organization}), op.DB.Save(organization).Error
}

func (op *OrganizationPersister) GetOrganzation(id uuid.UUID) ([]byte, error) {
	var organization Organization
	err := op.DB.First(&organization, id).Error
	return marshalOrganization(&organization), err
}

func marshalOrganization(org *Organization) []byte {
	res, _ := json.Marshal(org)

	return res
}

func marshalOrganizations(orgs []Organization) []byte {
	res, _ := json.Marshal(orgs)

	return res
}

func marshalOrganizationsPage(orgsPage *OrganizationsPage) []byte {
	res, _ := json.Marshal(orgsPage)

	return res
}
