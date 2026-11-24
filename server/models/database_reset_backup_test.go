package models

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshkit/database"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestBackupEntity struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func TestBackupSystemDatabase(t *testing.T) {
	tempDir := t.TempDir()
	originalUserData := viper.GetString("USER_DATA_FOLDER")
	viper.Set("USER_DATA_FOLDER", tempDir)
	defer viper.Set("USER_DATA_FOLDER", originalUserData)

	dbFile := filepath.Join(tempDir, "mesherydb.sql")
	db, err := database.New(database.Options{
		Filename: fmt.Sprintf("file:%s?cache=private&mode=rwc&_busy_timeout=10000&_journal_mode=WAL", filepath.ToSlash(dbFile)),
		Engine:   database.SQLITE,
	})
	require.NoError(t, err, "failed to initialize SQLite DB")
	t.Cleanup(func() {
		_ = db.DBClose()
	})

	err = db.AutoMigrate(&TestBackupEntity{})
	require.NoError(t, err)

	entity := TestBackupEntity{ID: "design-1", Name: "My Pattern"}
	err = db.Create(&entity).Error
	require.NoError(t, err)

	// 1. Successful backup
	backupPath1, err := BackupSystemDatabase(&db)
	require.NoError(t, err)
	assert.FileExists(t, backupPath1)
	assert.Contains(t, backupPath1, filepath.Join(tempDir, ".archive"))

	// Verify the backup contains the committed data
	backupDB, err := database.New(database.Options{
		Filename: backupPath1,
		Engine:   database.SQLITE,
	})
	require.NoError(t, err)
	defer func() {
		_ = backupDB.DBClose()
	}()

	var readEntity TestBackupEntity
	err = backupDB.First(&readEntity, "id = ?", "design-1").Error
	require.NoError(t, err)
	assert.Equal(t, "My Pattern", readEntity.Name)

	// 2. Distinct timestamp for second backup
	backupPath2, err := BackupSystemDatabase(&db)
	require.NoError(t, err)
	assert.FileExists(t, backupPath2)
	assert.NotEqual(t, backupPath1, backupPath2, "subsequent backups must not overwrite each other")

	// 3. Nil handler returns error
	_, err = BackupSystemDatabase(nil)
	assert.Error(t, err)

	// 4. Refuse existing destination
	err = os.WriteFile(filepath.Join(tempDir, ".archive", "existing.sql"), []byte("data"), 0644)
	require.NoError(t, err)
	// Direct VACUUM INTO existing path fails
	err = db.Exec("VACUUM INTO ?", filepath.Join(tempDir, ".archive", "existing.sql")).Error
	assert.Error(t, err, "VACUUM INTO must refuse to overwrite existing file")
}
