package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
)

// PerformanceProfilePersister is the persister for persisting
// performance profiles on the database
type PerformanceProfilePersister struct {
	DB *database.Handler
}

// nilUUIDString keeps reads backward-compatible with local databases whose
// performance_profiles rows predate the owner column. The schema-native model
// types Owner as a non-null UUID, so NULL owner values need a scan-safe
// fallback.
const nilUUIDString = "00000000-0000-0000-0000-000000000000"

// The column lists below MUST name the columns AutoMigrate derives from the
// schemas PerformanceProfile - notably `owner`, from its Owner field. They named
// the pre-rename `user_id` after schemas renamed UserID to Owner; that column no
// longer exists, so every read failed with "no such column: user_id" and, since
// the gorm errors were discarded, reported an empty profile list rather than the
// failure. Re-check these lists whenever the schemas model changes.

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
		// COALESCE(owner, nilUUIDString) preserves old local rows that were
		// saved before the schema-native PerformanceProfile introduced owner.
		Select(`
			id, name, COALESCE(owner, ?) as owner, load_generators,
			endpoints, service_mesh, concurrent_request, qps,
			duration, request_headers, request_cookies,
			request_body, content_type, metadata, created_at,
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

	if err := query.Table("performance_profiles").Count(&count).Error; err != nil {
		return nil, err
	}

	if err := Paginate(uint(page), uint(pageSize))(query).Find(&profiles).Error; err != nil {
		return nil, err
	}

	performanceProfilePage := &PerformanceProfilePage{
		Page:       int(page),
		PageSize:   int(pageSize),
		TotalCount: int(count),
		Profiles:   profiles,
	}

	return json.Marshal(performanceProfilePage)
}

// DeletePerformanceProfile takes in a profile id and deletes it. When no profile
// with that id is stored it returns ErrResultNotFound, so a delete of a profile
// that never existed is distinguishable from a successful delete.
func (ppp *PerformanceProfilePersister) DeletePerformanceProfile(id core.Uuid) ([]byte, error) {
	profile := PerformanceProfile{ID: id}
	result := ppp.DB.Delete(&profile)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrResultNotFound(fmt.Errorf("performance profile %s not found", id.String()))
	}

	return json.Marshal(&profile)
}

// SavePerformanceProfile updates the profile with this ID, or inserts it when
// none is stored. An edit from the performance form carries every form field
// but no createdAt, so an update writes every column except created_at.
func (ppp *PerformanceProfilePersister) SavePerformanceProfile(_ core.Uuid, profile *PerformanceProfile) error {
	res := ppp.DB.Model(profile).Select("*").Omit("created_at").Updates(profile)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return nil
	}
	return ppp.DB.Create(profile).Error
}

func (ppp *PerformanceProfilePersister) GetPerformanceProfile(id core.Uuid) (*PerformanceProfile, error) {
	var performanceProfile PerformanceProfile

	err := ppp.DB.
		Table("performance_profiles").
		// Keep single-profile reads compatible with old rows whose owner is
		// absent/NULL after AutoMigrate adds the column.
		Select(`
			id, name, COALESCE(owner, ?) as owner, schedule, load_generators,
			endpoints, service_mesh, concurrent_request, qps, duration,
			request_headers, request_cookies, request_body, content_type,
			metadata, last_run, total_results, created_at, updated_at`,
			nilUUIDString,
		).
		Where("id = ?", id).
		First(&performanceProfile).Error
	return &performanceProfile, err
}
