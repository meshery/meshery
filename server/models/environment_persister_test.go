package models

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/v1beta3/environment"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newEnvironmentTestDB(t *testing.T) *database.Handler {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	handler := &database.Handler{
		DB: db,
	}

	if err := handler.AutoMigrate(
		&environment.Environment{},
		&environment.EnvironmentConnectionMapping{},
		&connections.Connection{},
	); err != nil {
		t.Fatalf("failed to auto-migrate environment models: %v", err)
	}

	return handler
}

// Regression test for #21011.
// Empty page/pageSize parameters previously caused strconv.ParseUint("")
// to fail before default pagination values were applied.
func TestGetEnvironmentConnections_UsesDefaultPaginationWhenParamsAreEmpty(t *testing.T) {
	dbHandler := newEnvironmentTestDB(t)

	ep := &EnvironmentPersister{
		DB: dbHandler,
	}

	envID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate environment ID: %v", err)
	}

	result, err := ep.GetEnvironmentConnections(
		envID,
		"", // search
		"", // order
		"", // page
		"", // pageSize
		"", // filter
	)
	if err != nil {
		t.Fatalf("expected empty pagination params to use defaults, got error: %v", err)
	}

	var page connections.ConnectionPage

	if err := json.Unmarshal(result, &page); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if page.Page != 0 {
		t.Fatalf("expected page to default to 0, got %d", page.Page)
	}

}
