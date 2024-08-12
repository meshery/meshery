package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshkit/logger"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/registration"
)

var ModelsPath = "../meshmodel"

func getModelDirectoryPaths() ([]string, error) {
	dirEntries := make([]string, 0)
	modelsDirs, err := os.ReadDir(ModelsPath)
	if err != nil {
		return dirEntries, err
	}
	for _, modelDir := range modelsDirs {
		if modelDir.Name() != "kubernetes" {
			continue
		}

		if !modelDir.IsDir() {
			continue
		}
		modelVersionsDirPath := filepath.Join(ModelsPath, modelDir.Name())
		modelVersionsDir, err := os.ReadDir(modelVersionsDirPath)
		if err != nil {
			return dirEntries, err
		}
		for _, version := range modelVersionsDir {
			if !version.IsDir() {
				continue
			}
			modelDefVersionsDirPath := filepath.Join(modelVersionsDirPath, version.Name())
			modelDefVersionsDir, err := os.ReadDir(modelDefVersionsDirPath)
			if err != nil {
				return dirEntries, err
			}
			for _, defVersion := range modelDefVersionsDir {
				if !defVersion.IsDir() {
					continue
				}
				defPath := filepath.Join(modelDefVersionsDirPath, defVersion.Name())
				dirEntries = append(dirEntries, defPath)
			}
		}
	}
	return dirEntries, nil
}

func SeedComponents(log logger.Handler, hc *HandlerConfig, regm *meshmodel.RegistryManager) {
	regErrorStore := NewRegistrationFailureLogHandler()
	regHelper := registration.NewRegistrationHelper(utils.UI, regm, regErrorStore)
	modelDirPaths, err := getModelDirectoryPaths()
	if err != nil {
		log.Error(ErrSeedingComponents(err))
	}
	fmt.Println("MODEL DIR PATHS : ", modelDirPaths)
	for _, dirPath := range modelDirPaths {
		dir := registration.NewDir(dirPath)
		regHelper.Register(dir)
	}
	RegistryLog(log, hc, regm, regErrorStore)
}
