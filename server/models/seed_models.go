package models

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"time"

	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/models/registration"
	meshkitUtils "github.com/meshery/meshkit/utils"
)

var ModelsPath = "../../models"

const PoliciesPath = "../../models/meshery-core/0.7.2/v1.0.0/policies"

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
		// NOTE
		// Temporarily:  remove this once the connection and credentials of k8s is written to repective version is implemented in the generator, and the namespace bug (where the component isNamespace is incorrectly marked as true) is resolved.
		// if modelName == "kubernetes" {
		// 	sortedVersionDirs[0] = "../../models/kubernetes/v1.32.0-alpha.3"
		// }
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

// RunSeedStage runs one seeding stage, converting a panic inside it into a
// logged, structured error.
//
// Nine of the ten stages run in a bare goroutine - the boot goroutine in
// cmd/main.go and the two paths that reseed after a reset
// (handlers.ResetSystemDatabase and the GraphQL resyncCluster hard reset) -
// where Go gives a panic no second chance. The tenth, boot's "content" stage,
// runs synchronously on main's own goroutine, where a panic unwinds main and
// exits just as surely. Either shape takes down the whole process, HTTP
// listener included, which is not a proportionate response to a seeding fault
// - a Meshery Server with an incomplete registry is still useful, and the
// operator can read the error and act on it, whereas a server that exits at
// boot leaves them with a crash loop and no UI to read it in
// (meshery/meshery#21584).
//
// Each stage is wrapped separately so one faulting stage does not skip the
// others.
//
// The cover is bounded by what recover can reach: only panics on the wrapped
// stage's own goroutine. A stage that itself spawns a goroutine - SeedKeys
// does, in keys_helper.go - leaves that goroutine outside this recover, and a
// panic there still terminates the process. Recovering it belongs where it is
// spawned, not here.
func RunSeedStage(log logger.Handler, stage string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(ErrSeedingStagePanic(stage, r, debug.Stack()))
		}
	}()
	fn()
}

// SeedComponents registers the latest versions of models
func SeedComponents(log logger.Handler, hc *HandlerConfig, regm *meshmodel.RegistryManager, db *database.Handler) {
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

	// Registration has now put every connection definition in the registry and
	// created the registrant Connections they describe. Seeding is folded in
	// here, rather than left to each caller, so that every path that rebuilds
	// the registry - boot, database reset, hard reset - seeds too.
	SeedConnections(log, db, regm)
}
