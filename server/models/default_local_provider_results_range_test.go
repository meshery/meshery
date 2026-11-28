package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GET /api/performance/results documents `from` and `to` as the dates that
// bound a result's test start time; the performance calendar asks for the
// month it shows. The results outside that range must not fill the page.
func TestFetchAllResultsKeepsOnlyTheRequestedDateRange(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&MesheryResult{}))
	log, err := logger.New("test", logger.Options{})
	require.NoError(t, err)
	provider := &DefaultLocalProvider{
		ResultPersister: &MesheryResultsPersister{DB: &database.Handler{DB: db}},
		Log:             log,
	}

	seed := func(name string, start time.Time) {
		t.Helper()
		// Stored the way WriteResult stores a run.
		require.NoError(t, db.Table("meshery_results").Save(convertMesheryResultToLocalRepresentation(&MesheryResult{
			ID:            uuid.Must(uuid.NewV4()),
			Name:          name,
			TestStartTime: &start,
		})).Error)
	}
	// Ten October runs were saved first, then one in November, and the last
	// day of the range counts.
	for i := 1; i <= 10; i++ {
		seed("october", time.Date(2026, time.October, i, 12, 0, 0, 0, time.UTC))
	}
	seed("november", time.Date(2026, time.November, 24, 12, 0, 0, 0, time.UTC))
	seed("last day", time.Date(2026, time.November, 30, 23, 0, 0, 0, time.UTC))
	seed("december", time.Date(2026, time.December, 1, 0, 30, 0, 0, time.UTC))

	body, err := provider.FetchAllResults("", "0", "10", "", "", "2026-11-01", "2026-11-30")
	require.NoError(t, err)

	var page MesheryResultPage
	require.NoError(t, json.Unmarshal(body, &page))
	names := []string{}
	for _, r := range page.Results {
		names = append(names, r.Name)
	}
	require.ElementsMatch(t, []string{"november", "last day"}, names)
	require.Equal(t, 2, page.TotalCount)

	for _, tc := range []struct {
		name string
		from string
		to   string
	}{
		{name: "from", from: "2026-13-40"},
		{name: "to", to: "2026-02-30"},
	} {
		t.Run("invalid "+tc.name, func(t *testing.T) {
			_, err := provider.FetchAllResults("", "0", "10", "", "", tc.from, tc.to)
			require.ErrorIs(t, err, ErrInvalidPerformanceResultDate)
			require.ErrorContains(t, err, tc.name+" must use YYYY-MM-DD")
		})
	}
}
