package models

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshkit/logger"
	meshkitUtils "github.com/layer5io/meshkit/utils"

	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/registration"
)

var ModelsPath = "../meshmodel"

// Define the maximum number of worker goroutines
const maxWorkers = 10

func getModelDirectoryPaths() ([]string, error) {
	modelsDirs, err := os.ReadDir(ModelsPath)
	if err != nil {
		return nil, err
	}

	var dirEntries []string
	var mutex sync.Mutex
	errors := make(chan error, len(modelsDirs))

	workerFunc := func(modelDir os.DirEntry) {
		if !modelDir.IsDir() {
			return
		}
		modelVersionsDirPath := filepath.Join(ModelsPath, modelDir.Name())
		modelVersionsDir, err := os.ReadDir(modelVersionsDirPath)
		if err != nil {
			errors <- err
			return
		}
		for _, version := range modelVersionsDir {
			if !version.IsDir() {
				continue
			}
			modelDefVersionsDirPath := filepath.Join(modelVersionsDirPath, version.Name())
			modelDefVersionsDir, err := os.ReadDir(modelDefVersionsDirPath)
			if err != nil {
				errors <- err
				return
			}
			for _, defVersion := range modelDefVersionsDir {
				if defVersion.IsDir() {
					defPath := filepath.Join(modelDefVersionsDirPath, defVersion.Name())
					mutex.Lock()
					dirEntries = append(dirEntries, defPath)
					mutex.Unlock()
				}
			}
		}
	}

	meshkitUtils.WorkerPool(modelsDirs, workerFunc, maxWorkers)

	close(errors)

	var firstError error
	for err := range errors {
		if firstError == nil && err != nil {
			firstError = err
		}
	}

	return dirEntries, firstError
}

func SeedComponents(log logger.Handler, hc *HandlerConfig, regm *meshmodel.RegistryManager) {
	regErrorStore := NewRegistrationFailureLogHandler()
	regHelper := registration.NewRegistrationHelper(utils.UI, regm, regErrorStore)

	modelDirPaths, err := getModelDirectoryPaths()
	if err != nil {
		log.Error(ErrSeedingComponents(err))
		return
	}

	workerFunc := func(dirPath string) {
		dir := registration.NewDir(dirPath)
		regHelper.Register(dir)
	}

	meshkitUtils.WorkerPool(modelDirPaths, workerFunc, maxWorkers)

	RegistryLog(log, hc, regm, regErrorStore)
}
