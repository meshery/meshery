package models

import (
	"encoding/json"
	"strings"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
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
	// Sort-input whitelist dual-accepts the canonical camelCase key
	// (`lastRun`) and the legacy snake_case key (`last_run`) for the
	// Phase 2.K cascade deprecation window. Once all UI callers emit
	// `lastRun` exclusively, drop the snake_case entry.
	order = SanitizeOrderInput(order, []string{"updated_at", "created_at", "name", "last_run", "lastRun"})
	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}
	// Translate the canonical camelCase sort key to the SQL column
	// alias (`last_run`) below, where the SELECT aliases the
	// aggregated subquery back to snake_case. The wire contract is
	// camelCase; the DB column layout stays snake_case per the
	// identifier-naming migration (wire vs. DB are distinct layers).
	if strings.HasPrefix(order, "lastRun") {
		order = "last_run" + strings.TrimPrefix(order, "lastRun")
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
func (ppp *PerformanceProfilePersister) DeletePerformanceProfile(id core.Uuid) ([]byte, error) {
	profile := PerformanceProfile{ID: &id}
	ppp.DB.Delete(profile)

	return marshalPerformanceProfile(&profile), nil
}

func (ppp *PerformanceProfilePersister) SavePerformanceProfile(_ core.Uuid, profile *PerformanceProfile) error {
	return ppp.DB.Save(profile).Error
}

func (ppp *PerformanceProfilePersister) GetPerformanceProfile(id core.Uuid) (*PerformanceProfile, error) {
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
