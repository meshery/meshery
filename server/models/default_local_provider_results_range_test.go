package models

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFetchAllResultsUsesOneSnapshot(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "results.db")), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })
	var journalMode string
	require.NoError(t, db.Raw("PRAGMA journal_mode=WAL").Scan(&journalMode).Error)
	require.Equal(t, "wal", journalMode)
	require.NoError(t, db.AutoMigrate(&MesheryResult{}))
	persister := &MesheryResultsPersister{DB: &database.Handler{DB: db}}

	lateResult, err := json.Marshal(&MesheryResult{
		Result: map[string]interface{}{"StartTime": "2026-11-30T23:59:00Z"},
	})
	require.NoError(t, err)
	saved := false
	require.NoError(t, db.Callback().Query().After("gorm:query").Register("save-late-result", func(_ *gorm.DB) {
		if !saved {
			saved = true
			require.NoError(t, persister.WriteResult(uuid.Must(uuid.NewV4()), lateResult))
		}
	}))

	fetch := func() MesheryResultPage {
		t.Helper()
		body, err := persister.GetAllResults(0, 10, "2026-11-30", "2026-11-30", nil)
		require.NoError(t, err)
		var page MesheryResultPage
		require.NoError(t, json.Unmarshal(body, &page))
		return page
	}

	first := fetch()
	require.True(t, saved)
	require.Zero(t, first.TotalCount)
	require.Empty(t, first.Results)
	second := fetch()
	require.Equal(t, 1, second.TotalCount)
	require.Len(t, second.Results, 1)
}

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

	seedLegacy := func(name, runnerStart string, savedAt time.Time) {
		t.Helper()
		runnerResults := map[string]interface{}{}
		if runnerStart != "" {
			runnerResults["StartTime"] = runnerStart
		}
		require.NoError(t, db.Table("meshery_results").Save(convertMesheryResultToLocalRepresentation(&MesheryResult{
			ID:            uuid.Must(uuid.NewV4()),
			Name:          name,
			Result:        runnerResults,
			TestStartTime: &savedAt,
		})).Error)
	}
	for i := 1; i <= 10; i++ {
		seedLegacy("october", "", time.Date(2026, time.October, i, 12, 0, 0, 0, time.UTC))
	}
	seedLegacy("crossed midnight", "2026-11-30T23:59:00Z", time.Date(2026, time.December, 1, 0, 1, 0, 0, time.UTC))
	seedLegacy("legacy offset", "", time.Date(2026, time.December, 1, 0, 15, 0, 0, time.FixedZone("", 5*3600+30*60)))

	id := uuid.Must(uuid.NewV4())
	result, err := json.Marshal(&MesheryResult{
		Name:   "first day",
		Result: map[string]interface{}{"StartTime": "2026-11-01T00:15:00.123456789+05:30"},
	})
	require.NoError(t, err)
	require.NoError(t, provider.ResultPersister.WriteResult(id, result))

	fetch := func(pageNumber, size, from, to string) MesheryResultPage {
		t.Helper()
		body, err := provider.FetchAllResults("", pageNumber, size, "", "", from, to)
		require.NoError(t, err)
		var page MesheryResultPage
		require.NoError(t, json.Unmarshal(body, &page))
		require.NotNil(t, page.Results)
		return page
	}

	page := fetch("0", "10", "2026-11-01", "2026-11-30")
	require.Equal(t, 2, page.TotalCount)
	starts := map[string]string{}
	for _, result := range page.Results {
		require.NotNil(t, result.TestStartTime)
		starts[result.Name] = result.TestStartTime.Format(time.RFC3339Nano)
	}
	require.Equal(t, map[string]string{
		"first day":        "2026-11-01T00:15:00.123456789+05:30",
		"crossed midnight": "2026-11-30T23:59:00Z",
	}, starts)

	first := fetch("0", "1", "2026-11-01", "2026-11-30")
	second := fetch("1", "1", "2026-11-01", "2026-11-30")
	require.Equal(t, 2, first.TotalCount)
	require.Equal(t, 2, second.TotalCount)
	require.Len(t, first.Results, 1)
	require.Len(t, second.Results, 1)
	require.ElementsMatch(t, []string{"first day", "crossed midnight"}, []string{first.Results[0].Name, second.Results[0].Name})

	december := fetch("0", "10", "2026-12-01", "")
	require.Equal(t, 1, december.TotalCount)
	require.Len(t, december.Results, 1)
	require.Equal(t, "legacy offset", december.Results[0].Name)
	require.NotNil(t, december.Results[0].TestStartTime)
	require.Equal(t, "2026-12-01T00:15:00+05:30", december.Results[0].TestStartTime.Format(time.RFC3339Nano))
	require.Equal(t, 12, fetch("0", "20", "", "2026-11-30").TotalCount)
	require.Equal(t, 13, fetch("0", "20", "", "").TotalCount)

	empty := fetch("0", "10", "2027-01-01", "2027-01-31")
	require.Zero(t, empty.TotalCount)
	require.Empty(t, empty.Results)
	pastEnd := fetch("2", "1", "2026-11-01", "2026-11-30")
	require.Equal(t, 2, pastEnd.TotalCount)
	require.Empty(t, pastEnd.Results)

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
