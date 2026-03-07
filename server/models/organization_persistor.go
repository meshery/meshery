package models

import (
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/v1beta1/organization"
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
	organizations := []*organization.Organization{}

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

func (op *OrganizationPersister) SaveOrganization(org *organization.Organization) ([]byte, error) {
	if org.Id == uuid.Nil {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}

		org.Id = id
	}

	return marshalOrganizations([]organization.Organization{*org}), op.DB.Save(org).Error
}

func (op *OrganizationPersister) GetOrganzation(id uuid.UUID) ([]byte, error) {
	var organization organization.Organization
	err := op.DB.First(&organization, id).Error
	return marshalOrganization(&organization), err
}

func (op *OrganizationPersister) GetOrganizationsCount() (int64, error) {
	var count int64
	if err := op.DB.Model(&organization.Organization{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func marshalOrganization(org *organization.Organization) []byte {
	res, _ := json.Marshal(org)

	return res
}

func marshalOrganizations(orgs []organization.Organization) []byte {
	res, _ := json.Marshal(orgs)

	return res
}

func marshalOrganizationsPage(orgsPage *OrganizationsPage) []byte {
	res, _ := json.Marshal(orgsPage)

	return res
}
