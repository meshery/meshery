package models

import (
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
)

// PerformanceProfilePersister is the persister for persisting
// performance profiles on the database
type PerformanceProfilePersister struct {
	DB *database.Handler
}

// PerformanceProfilePage represents a page of performance profiles
type PerformanceProfilePage struct {
	Page       uint64                `json:"page"`
	PageSize   uint64                `json:"page_size"`
	TotalCount int                   `json:"total_count"`
	Profiles   []*PerformanceProfile `json:"profiles"`
}

// GetPerformanceProfiles returns all of the performance profiles
func (ppp *PerformanceProfilePersister) GetPerformanceProfiles(_, search, order string, page, pageSize uint64) ([]byte, error) {
	order = sanitizeOrderInput(order, []string{"updated_at", "created_at", "name", "last_run"})
	if order == "" {
		order = "updated_at desc"
	}

	count := int64(0)
	profiles := []*PerformanceProfile{}

	query := ppp.DB.
		Select(`
		id, name, load_generators,
		endpoints, qps, service_mesh,
		duration, request_headers, request_cookies,
		request_body, content_type, created_at,
		updated_at, (?) as last_run, (?) as total_results`,
			ppp.DB.Table("meshery_results").Select("DATETIME(MAX(meshery_results.test_start_time))").Where("performance_profile = performance_profiles.id"),
			ppp.DB.Table("meshery_results").Select("COUNT(meshery_results.name)").Where("performance_profile = performance_profiles.id"),
		).
		Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(performance_profiles.name) like ?)", like)
	}

	query.Table("performance_profiles").Count(&count)

	Paginate(uint(page), uint(pageSize))(query).Find(&profiles)

	performanceProfilePage := &PerformanceProfilePage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Profiles:   profiles,
	}

	return marshalPerformanceProfilePage(performanceProfilePage), nil
}

// DeletePerformanceProfile takes in a profile id and delete it if it already exists
func (ppp *PerformanceProfilePersister) DeletePerformanceProfile(id uuid.UUID) ([]byte, error) {
	profile := PerformanceProfile{ID: &id}
	ppp.DB.Delete(profile)

	return marshalPerformanceProfile(&profile), nil
}

func (ppp *PerformanceProfilePersister) SavePerformanceProfile(_ uuid.UUID, profile *PerformanceProfile) error {
	return ppp.DB.Save(profile).Error
}

func (ppp *PerformanceProfilePersister) GetPerformanceProfile(id uuid.UUID) (*PerformanceProfile, error) {
	var performanceProfile PerformanceProfile

	err := ppp.DB.First(&performanceProfile, id).Error
	return &performanceProfile, err
}

func marshalPerformanceProfilePage(ppp *PerformanceProfilePage) []byte {
	res, _ := json.Marshal(ppp)

	return res
}

func marshalPerformanceProfile(pp *PerformanceProfile) []byte {
	res, _ := json.Marshal(pp)

	return res
}
