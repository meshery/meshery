package models

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshkit/logger"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/registration"
	meshkitUtils "github.com/layer5io/meshkit/utils"
)

var ModelsPath = "../meshmodel"

const PoliciesPath = "../meshmodel/kubernetes/v1.25.2/v1.0.0/policies"

// versionInfo holds information about a version directory
type versionInfo struct {
	dirName string
	modTime time.Time
	dirPath string
}

// GetModelDirectoryPaths retrieves model definition directories based on the following criteria:
// 1. Find the latest version of each model that contains a non-empty 'components' directory.
//   - If the latest version has a non-empty 'components', use its path.
//   - If not, search previous versions in descending order to find the nearest version with non-empty 'components'.
//
// 2. For all versions (including the one used for 'components'), check if the 'relationships' directory is non-empty and include their paths if so.
// The returned directories are sorted with the latest version first.
func GetModelDirectoryPaths(modelPath string) ([]string, error) {
	dirEntries := []string{}

	// Read all model directories (e.g., accurate, kubernetes)
	modelsDirs, err := os.ReadDir(modelPath)
	if err != nil {
		return dirEntries, meshkitUtils.ErrReadDir(err, fmt.Sprintf("failed to read models directory '%s'", modelPath))
	}

	for _, modelDir := range modelsDirs {
		if !modelDir.IsDir() {
			continue
		}

		modelName := modelDir.Name()
		modelVersionsDirPath := filepath.Join(modelPath, modelName)

		// Get all version directories sorted in descending order (latest first)
		sortedVersionDirs, err := meshkitUtils.GetAllVersionDirsSortedDesc(modelVersionsDirPath)
		if err != nil {
			continue
		}
		if len(sortedVersionDirs) == 0 {
			continue
		}
		modelDefDirPath, err := getLatestModelDefDir(sortedVersionDirs[0])
		if err != nil {
			continue
		}
		dirEntries = append(dirEntries, modelDefDirPath)
	}

	return dirEntries, nil
}

// getLatestModelDefDir returns the path to the latest model definition directory based on modification time
func getLatestModelDefDir(latestVersionDirPath string) (string, error) {
	entries, err := os.ReadDir(latestVersionDirPath)
	if err != nil {
		return "", meshkitUtils.ErrReadDir(err, fmt.Sprintf("failed to read model definition directory '%s'", latestVersionDirPath))
	}

	if len(entries) == 0 {
		return "", nil
	}

	modelDefs := []versionInfo{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		defDirPath := filepath.Join(latestVersionDirPath, entry.Name())
		fi, err := os.Stat(defDirPath)
		if err != nil {
			continue
		}

		modelDefs = append(modelDefs, versionInfo{
			dirName: entry.Name(),
			modTime: fi.ModTime(),
			dirPath: defDirPath,
		})
	}

	if len(modelDefs) == 0 {
		return "", nil
	}

	sort.Slice(modelDefs, func(i, j int) bool {
		return modelDefs[i].modTime.After(modelDefs[j].modTime)
	})

	return modelDefs[0].dirPath, nil
}

// SeedComponents registers the latest versions of models
func SeedComponents(log logger.Handler, hc *HandlerConfig, regm *meshmodel.RegistryManager) {
	regErrorStore := NewRegistrationFailureLogHandler()
	regHelper := registration.NewRegistrationHelper(utils.UI, regm, regErrorStore)

	modelDirPaths, err := GetModelDirectoryPaths(ModelsPath)
	if err != nil {
		log.Error(ErrSeedingComponents(err))
	}

	for _, dirPath := range modelDirPaths {
		dir := registration.NewDir(dirPath)
		regHelper.Register(dir)
	}

	RegistryLog(log, hc, regm, regErrorStore)
}
