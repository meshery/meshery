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

// nilUUIDString keeps reads backward-compatible with local databases created
// before performance_profiles had a user_id column/value. The schema-native
// model uses a non-null UUID, so NULL user_id values need a scan-safe fallback.
const nilUUIDString = "00000000-0000-0000-0000-000000000000"

// GetPerformanceProfiles returns all of the performance profiles
func (ppp *PerformanceProfilePersister) GetPerformanceProfiles(_, search, order string, page, pageSize uint64) ([]byte, error) {
	// The wire contract is camelCase (`?order=lastRun desc`) while the SELECT
	// below aliases the aggregated subquery to the snake_case `last_run`
	// column. SanitizeOrderInput folds the wire name onto the DB column, so
	// only the DB spelling is listed here.
	order = SanitizeOrderInput(order, []string{"updated_at", "created_at", "name", "last_run"})
	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	count := int64(0)
	profiles := []PerformanceProfile{}

	query := ppp.DB.
		// COALESCE(user_id, nilUUIDString) preserves old local rows that were
		// saved before the schema-native PerformanceProfile introduced userId.
		Select(`
			id, name, COALESCE(user_id, ?) as user_id, load_generators,
			endpoints, qps, service_mesh,
			duration, request_headers, request_cookies,
			request_body, content_type, created_at,
			updated_at, (?) as last_run, (?) as total_results`,
			nilUUIDString,
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
		Page:       int(page),
		PageSize:   int(pageSize),
		TotalCount: int(count),
		Profiles:   profiles,
	}

	return marshalPerformanceProfilePage(performanceProfilePage), nil
}

// DeletePerformanceProfile takes in a profile id and delete it if it already exists
func (ppp *PerformanceProfilePersister) DeletePerformanceProfile(id core.Uuid) ([]byte, error) {
	profile := PerformanceProfile{ID: id}
	ppp.DB.Delete(&profile)

	return marshalPerformanceProfile(&profile), nil
}

func (ppp *PerformanceProfilePersister) SavePerformanceProfile(_ core.Uuid, profile *PerformanceProfile) error {
	return ppp.DB.Save(profile).Error
}

func (ppp *PerformanceProfilePersister) GetPerformanceProfile(id core.Uuid) (*PerformanceProfile, error) {
	var performanceProfile PerformanceProfile

	err := ppp.DB.
		Table("performance_profiles").
		// Keep single-profile reads compatible with old rows whose user_id is
		// absent/NULL after AutoMigrate adds the column.
		Select(`
			id, name, COALESCE(user_id, ?) as user_id, schedule, load_generators,
			endpoints, service_mesh, concurrent_request, qps, duration,
			request_headers, request_cookies, request_body, content_type,
			metadata, last_run, total_results, created_at, updated_at`,
			nilUUIDString,
		).
		Where("id = ?", id).
		First(&performanceProfile).Error
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
