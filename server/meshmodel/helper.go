package meshmodel

import (
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/server/models"
	reg "github.com/layer5io/meshery/server/models/meshmodel/registration"
	"github.com/layer5io/meshkit/logger"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
)

var ModelsPath = "../meshmodel"

func GetModelDirectoryPaths() ([]string, error) {
	dirEntries := make([]string, 0)
	modelsDirs, err := os.ReadDir(ModelsPath)
	if(err != nil){
		return dirEntries, err
	}
	for _, modelDir := range modelsDirs {
		if(!modelDir.IsDir()){
			continue
		}
		modelVersionsDirPath := filepath.Join(ModelsPath, modelDir.Name())
		modelVersionsDir, err := os.ReadDir(modelVersionsDirPath)
		if(err != nil){
			return dirEntries, err
		}
		for _, version := range modelVersionsDir {
			if(!version.IsDir()){
				continue
			}
			modelDefVersionsDirPath := filepath.Join(modelVersionsDirPath, version.Name())
			modelDefVersionsDir, err := os.ReadDir(modelDefVersionsDirPath)
			if(err != nil){
				return dirEntries, err
			}
			for _, defVersion := range modelDefVersionsDir {
				if(!defVersion.IsDir()){
					continue
				}
				defPath := filepath.Join(modelDefVersionsDirPath, defVersion.Name())
				dirEntries = append(dirEntries, defPath)
			}
		}
	}
	return dirEntries, nil
}

func SeedComponents(log logger.Handler, hc *models.HandlerConfig, regm *meshmodel.RegistryManager) {
	regHelper := reg.NewRegistrationHelper(log, hc, regm)
	modelDirPaths, err := GetModelDirectoryPaths()
	if(err != nil){
		ErrSeedingComponents(err)
	}
	for _, dirPath := range modelDirPaths {
		dir := 	reg.NewDir(dirPath)
		regHelper.Register(dir)
	}
	regHelper.RegistryLog()
}
