package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResetSystemDatabase_RejectedWhenLockHeldWritesNoBackup(t *testing.T) {
	tempDir := t.TempDir()
	originalUserData := viper.GetString("USER_DATA_FOLDER")
	viper.Set("USER_DATA_FOLDER", tempDir)
	defer viper.Set("USER_DATA_FOLDER", originalUserData)

	dbFile := filepath.Join(tempDir, "mesherydb.sql")
	db, err := database.New(database.Options{
		Filename: models.SystemDatabaseDSN(dbFile),
		Engine:   database.SQLITE,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.DBClose()
	})

	provider := &models.DefaultLocalProvider{
		GenericPersister: &db,
	}

	h := &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}

	// Acquire reset lock beforehand to simulate an in-progress reset
	acquired := models.TryAcquireResetLock()
	require.True(t, acquired, "must acquire reset lock")
	defer models.ReleaseResetLock()

	req := httptest.NewRequest(http.MethodDelete, "/api/system/database/reset", nil)
	rw := httptest.NewRecorder()

	h.ResetSystemDatabase(rw, req, nil, nil, provider)

	resp := rw.Result()
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	// Ensure no backup files were created under .archive
	archiveDir := filepath.Join(tempDir, ".archive")
	entries, err := os.ReadDir(archiveDir)
	if err == nil {
		assert.Empty(t, entries, "rejected reset must not create or write any backup files")
	}
}

func TestResetSystemDatabase_NilPersister(t *testing.T) {
	h := &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}

	provider := &models.DefaultLocalProvider{
		GenericPersister: nil,
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/system/database/reset", nil)
	rw := httptest.NewRecorder()

	h.ResetSystemDatabase(rw, req, nil, nil, provider)

	resp := rw.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
