package models

import (
	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
)

// OrganizationPersister is the persister for persisting
// organizations on the database
type OrganizationPersister struct {
	BasePersister
}

// NewOrganizationPersister creates a new organization persister
func NewOrganizationPersister(db *database.Handler) *OrganizationPersister {
	return &OrganizationPersister{
		BasePersister: BasePersister{DB: db},
	}
}

// GetOrganizations returns all organizations with optimized query building
func (op *OrganizationPersister) GetOrganizations(search, order string, page, pageSize uint64, updatedAfter string) ([]byte, error) {
	organizations := []*Organization{}

	qb := op.NewQueryBuilder(&Organization{}).
		WithSearch(search, []string{"organizations.name"}).
		WithOrder(order, []string{"created_at", "updated_at", "name"}).
		WithUpdatedAfter(updatedAfter).
		Paginate(page, pageSize)

	count := qb.Count("organizations")
	if err := qb.Find(&organizations); err != nil {
		return nil, ErrDBRead(err)
	}

	organizationsPage := &OrganizationsPage{
		Page:          page,
		PageSize:      pageSize,
		TotalCount:    int(count),
		Organizations: organizations,
	}

	return MarshalJSON(organizationsPage), nil
}

func (op *OrganizationPersister) SaveOrganization(organization *Organization) ([]byte, error) {
	if err := op.SaveEntity(organization, organization.ID); err != nil {
		return nil, ErrDBCreate(err)
	}
	return MarshalJSON([]Organization{*organization}), nil
}

func (op *OrganizationPersister) GetOrganzation(id uuid.UUID) ([]byte, error) {
	var organization Organization
	if err := op.FindByID(&organization, id); err != nil {
		return nil, ErrDBRead(err)
	}
	return MarshalJSON(&organization), nil
}

func (op *OrganizationPersister) GetOrganizationsCount() (int64, error) {
	var count int64
	if err := op.DB.Model(&Organization{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
