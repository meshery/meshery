package meshmodel

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/server/models"
	reg "github.com/layer5io/meshery/server/models/meshmodel/registration"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/pkg/errors"
)

var ArtifactHubComponentsHandler = v1beta1.ArtifactHub{} //The components generated in output directory will be handled by kubernetes
var ModelsPath = "../meshmodel"
var RelativeRelationshipsPath = "relationships"

type EntityRegistrationHelper struct {
	handlerConfig    *models.HandlerConfig
	regManager       *meshmodel.RegistryManager
	componentChan    chan v1beta1.ComponentDefinition
	relationshipChan chan v1alpha2.RelationshipDefinition
	errorChan        chan error
	log              logger.Handler
}

func NewEntityRegistrationHelper(hc *models.HandlerConfig, rm *meshmodel.RegistryManager, log logger.Handler) *EntityRegistrationHelper {
	return &EntityRegistrationHelper{
		handlerConfig:    hc,
		regManager:       rm,
		componentChan:    make(chan v1beta1.ComponentDefinition),
		relationshipChan: make(chan v1alpha2.RelationshipDefinition),
		errorChan:        make(chan error),
		log:              log,
	}
}

// seed the local meshmodel components
func (erh *EntityRegistrationHelper) SeedComponents() {
	models, err := os.ReadDir(ModelsPath)
	if err != nil {
		 erh.log.Error(mutils.ErrReadDir(errors.Wrapf(err, "error while reading directory for generating components"), ModelsPath))
		return
	}

	for _, model := range models {
		modelVersionsDirPath := filepath.Join(ModelsPath, model.Name())
		// contains all versions for the model
		modelVersionsDir, err := os.ReadDir(modelVersionsDirPath)
		if err != nil {
			erh.log.Error(mutils.ErrReadDir(errors.Wrapf(err, "error while reading directory for registering components"), modelVersionsDirPath))
			continue
		}
		for _, version := range modelVersionsDir {
			modelDefVersionsDirPath := filepath.Join(modelVersionsDirPath, version.Name())

			modelDefVersionsDir, err := os.ReadDir(modelDefVersionsDirPath)
			if err != nil {
				erh.log.Error(mutils.ErrReadDir(errors.Wrapf(err, "error while reading directory for registering components"), modelVersionsDirPath))
			}
			for _, defVersion := range modelDefVersionsDir {
				defPath := filepath.Join(modelDefVersionsDirPath, defVersion.Name())
				fmt.Println("DEBUG: registering for models inside: ", defPath)
				// contains all def versions for a particular version of a model.
				dir := reg.NewDir(defPath)
				pkgUnit, err := dir.PkgUnit()
				if(err != nil){
					erh.log.Error(mutils.ErrReadDir(errors.Wrapf(err, "Given model directory is not a valid unit of packaging"), modelVersionsDirPath))
					continue
				}
				pipeline := reg.NewPipeline(erh.log, erh.handlerConfig, erh.regManager, pkgUnit)
				errorChan := make(chan error, 0)
				go func(echan chan error){
					for {
						err := <- echan
						// Using this temp
						fmt.Println(err)
						erh.log.Error(mutils.ErrCloseFile(errors.Wrapf(err, "Error in the registration pipeline")))
					}
				}(errorChan)
				pipeline.Start(errorChan)


			}


		}
	}

	erh.registryLog()
}
