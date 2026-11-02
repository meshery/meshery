package models

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/meshery/meshkit/database"
	"github.com/spf13/viper"
)

// BackupSystemDatabase takes a snapshot of the active SQLite system database
// into the .archive directory under USER_DATA_FOLDER using SQLite's VACUUM INTO.
// It creates a point-in-time, standalone backup that includes all committed WAL
// transactions, safely refuses to overwrite existing archives, and respects
// custom USER_DATA_FOLDER configuration.
func BackupSystemDatabase(dbHandler *database.Handler) (string, error) {
	if dbHandler == nil {
		return "", fmt.Errorf("nil database handler")
	}

	userDataFolder := viper.GetString("USER_DATA_FOLDER")
	if userDataFolder == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to determine user home directory: %w", err)
		}
		userDataFolder = filepath.Join(home, ".meshery", "config")
	}

	archiveDir := filepath.Join(userDataFolder, ".archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create database archive directory %q: %w", archiveDir, err)
	}

	currentTime := time.Now().UTC().Format("20060102150405.000000000")
	archiveFileName := fmt.Sprintf("mesherydb-%s.sql", currentTime)
	backupPath := filepath.Join(archiveDir, archiveFileName)

	if err := dbHandler.Exec("VACUUM INTO ?", backupPath).Error; err != nil {
		return "", fmt.Errorf("failed to vacuum system database into %q: %w", backupPath, err)
	}

	return backupPath, nil
}
