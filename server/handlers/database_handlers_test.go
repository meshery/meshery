package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	system "github.com/meshery/schemas/models/v1beta1/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newSystemDatabaseTestDB(t *testing.T) *database.Handler {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		if strings.Contains(err.Error(), "requires cgo") {
			t.Skip("skipping database test: sqlite requires cgo")
		}
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	if err := db.Exec("CREATE TABLE test_zebra (id int)").Error; err != nil {
		t.Fatalf("failed to create test_zebra table: %v", err)
	}
	if err := db.Exec("CREATE TABLE test_apple (id int)").Error; err != nil {
		t.Fatalf("failed to create test_apple table: %v", err)
	}
	if err := db.Exec("CREATE TABLE test_mango (id int)").Error; err != nil {
		t.Fatalf("failed to create test_mango table: %v", err)
	}

	return &database.Handler{DB: db, Mutex: &sync.Mutex{}}
}

func TestGetSystemDatabase_Ordering(t *testing.T) {
	dbHandler := newSystemDatabaseTestDB(t)
	h := &Handler{
		config:    &models.HandlerConfig{},
		log:       newTestLogger(t),
		dbHandler: dbHandler,
	}

	tests := []struct {
		name          string
		query         string
		expectedOrder []string
		assertOrder   bool
	}{
		{
			name:          "Separate REST ascending format",
			query:         "?order=name&sort=asc&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   true,
		},
		{
			name:          "Separate REST descending format",
			query:         "?order=name&sort=desc&search=test_",
			expectedOrder: []string{"test_zebra", "test_mango", "test_apple"},
			assertOrder:   true,
		},
		{
			name:          "Separate REST uppercase descending format",
			query:         "?order=name&sort=DESC&search=test_",
			expectedOrder: []string{"test_zebra", "test_mango", "test_apple"},
			assertOrder:   true,
		},
		{
			name:          "Separate REST mixed-case descending format",
			query:         "?order=name&sort=Desc&search=test_",
			expectedOrder: []string{"test_zebra", "test_mango", "test_apple"},
			assertOrder:   true,
		},
		{
			name:          "Separate REST descending format trims whitespace",
			query:         "?order=name&sort=%20DESC%20&search=test_",
			expectedOrder: []string{"test_zebra", "test_mango", "test_apple"},
			assertOrder:   true,
		},
		{
			name:          "Separate REST descending format trims column whitespace",
			query:         "?order=%20name%20&sort=desc&search=test_",
			expectedOrder: []string{"test_zebra", "test_mango", "test_apple"},
			assertOrder:   true,
		},
		{
			name:          "Bare column with omitted sort defaults to ascending",
			query:         "?order=name&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   true,
		},
		{
			name:          "UI compound ascending format",
			query:         "?order=name%20asc&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   true,
		},
		{
			name:          "UI compound descending format",
			query:         "?order=name%20desc&search=test_",
			expectedOrder: []string{"test_zebra", "test_mango", "test_apple"},
			assertOrder:   true,
		},
		{
			name:          "UI compound format trims surrounding whitespace",
			query:         "?order=%20name%20desc%20&search=test_",
			expectedOrder: []string{"test_zebra", "test_mango", "test_apple"},
			assertOrder:   true,
		},
		{
			name:          "Compound descending with separate sort=desc supplied",
			query:         "?order=name%20desc&sort=desc&search=test_",
			expectedOrder: []string{"test_zebra", "test_mango", "test_apple"},
			assertOrder:   true,
		},
		{
			name:          "Compound descending with conflicting separate sort=asc (compound precedence)",
			query:         "?order=name%20desc&sort=asc&search=test_",
			expectedOrder: []string{"test_zebra", "test_mango", "test_apple"},
			assertOrder:   true,
		},
		{
			name:          "Invalid direction defaults to ascending",
			query:         "?order=name&sort=invalid_dir&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   true,
		},
		{
			name:          "Invalid column falls back safely without error",
			query:         "?order=invalid_col&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   false,
		},
		{
			name:          "Non-existent created_at column falls back safely without error",
			query:         "?order=created_at&sort=desc&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   false,
		},
		{
			name:          "Non-existent updated_at column falls back safely without error",
			query:         "?order=updated_at&sort=desc&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   false,
		},
		{
			name:          "Uppercase direction handling in compound format",
			query:         "?order=name%20DESC&search=test_",
			expectedOrder: []string{"test_zebra", "test_mango", "test_apple"},
			assertOrder:   true,
		},
		{
			name:          "Empty order parameter falls back safely without error",
			query:         "?order=&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   false,
		},
		{
			name:          "Empty direction parameter defaults to ascending",
			query:         "?order=name&sort=&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   true,
		},
		{
			name:          "Multiple spaces in compound input falls back safely without error",
			query:         "?order=name%20%20desc&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   false,
		},
		{
			name:          "SQL-injection-style column input is rejected safely",
			query:         "?order=name;DROP%20TABLE%20test_apple;--&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   false,
		},
		{
			name:          "SQL-injection-style direction input is rejected or sanitized safely",
			query:         "?order=name%20desc;DROP%20TABLE%20test_apple;--&search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   false,
		},
		{
			name:          "Omitted order parameter returns tables safely",
			query:         "?search=test_",
			expectedOrder: []string{"test_apple", "test_mango", "test_zebra"},
			assertOrder:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/system/database"+tc.query, nil)
			rw := httptest.NewRecorder()

			h.GetSystemDatabase(rw, req, nil, nil, nil)

			require.Equal(t, http.StatusOK, rw.Code, rw.Body.String())

			var summary system.SystemDatabaseSummary
			err := json.Unmarshal(rw.Body.Bytes(), &summary)
			require.NoError(t, err, "response body must be valid JSON")

			var tableNames []string
			for _, tbl := range summary.Tables {
				tableNames = append(tableNames, tbl.Name)
			}

			if tc.assertOrder {
				assert.Equal(t, tc.expectedOrder, tableNames, "table names order mismatch for query %s", tc.query)
			} else {
				assert.ElementsMatch(t, tc.expectedOrder, tableNames, "tables must be present for query %s", tc.query)
			}
		})
	}
}
