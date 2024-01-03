package meshmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ArtifactHubComponentsHandler = meshmodel.ArtifactHub{} //The components generated in output directory will be handled by kubernetes
var ModelsPath = "../meshmodel"
var RelativeRelationshipsPath = "relationships"

type EntityRegistrationHelper struct {
	handlerConfig    *models.HandlerConfig
	regManager       *meshmodel.RegistryManager
	componentChan    chan v1alpha1.ComponentDefinition
	relationshipChan chan v1alpha1.RelationshipDefinition
	errorChan        chan error
	log              logger.Handler
}

func NewEntityRegistrationHelper(hc *models.HandlerConfig, rm *meshmodel.RegistryManager, log logger.Handler) *EntityRegistrationHelper {
	return &EntityRegistrationHelper{
		handlerConfig:    hc,
		regManager:       rm,
		componentChan:    make(chan v1alpha1.ComponentDefinition),
		relationshipChan: make(chan v1alpha1.RelationshipDefinition),
		errorChan:        make(chan error),
		log:              log,
	}
}

// seed the local meshmodel components
func (erh *EntityRegistrationHelper) SeedComponents() {
	// Watch channels and register components and relationships with the registry manager
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go erh.watchComponents(ctx)

	models, err := os.ReadDir(ModelsPath)
	if err != nil {
		erh.errorChan <- errors.Wrapf(err, "error while reading directory for generating components")
		return
	}

	relationships := make([]string, 0)

	// change to queue approach to register comps, relationships and policies
	// Read component and relationship definitions from files and send them to respective channels
	for _, model := range models {
		entitiesPath := filepath.Join(ModelsPath, model.Name())
		entities, err := os.ReadDir(entitiesPath)
		if err != nil {
			erh.errorChan <- errors.Wrapf(err, "error while reading directory for generating components")
			continue
		}

		for _, entity := range entities {
			entityPath := filepath.Join(entitiesPath, entity.Name())
			if entity.IsDir() {
				switch entity.Name() {
				case "relationships":
					relationships = append(relationships, entityPath)
				case "policies":
				default:
					erh.generateComponents(entityPath) // register components first
				}
			}
		}
	}

	for _, relationship := range relationships {
		erh.generateRelationships(relationship)
	}
}

// reads component definitions from files and sends them to the component channel
func (erh *EntityRegistrationHelper) generateComponents(pathToComponents string) {
	path, err := filepath.Abs(pathToComponents)
	if err != nil {
		erh.errorChan <- errors.Wrapf(err, "error while getting absolute path for generating components")
		return
	}

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if !info.IsDir() {
			// Read the component definition from file
			var comp v1alpha1.ComponentDefinition
			byt, err := os.ReadFile(path)
			if err != nil {
				erh.errorChan <- errors.Wrapf(err, fmt.Sprintf("unable to read file at %s", path))
				return nil
			}
			err = json.Unmarshal(byt, &comp)
			if err != nil {
				erh.errorChan <- errors.Wrapf(err, fmt.Sprintf("unmarshal json failed for %s", path))
				return nil
			}
			// Only register components that have been marked as published
			if comp.Metadata != nil && comp.Metadata["published"] == true {
				// Generate SVGs for the component and save them on the file system
				utils.WriteSVGsOnFileSystem(&comp)
				erh.componentChan <- comp
			}
		}
		return nil
	})
	if err != nil {
		erh.errorChan <- errors.Wrapf(err, "error while generating components")
	}
}

// reads relationship definitions from files and sends them to the relationship channel
func (erh *EntityRegistrationHelper) generateRelationships(pathToComponents string) {
	path, err := filepath.Abs(pathToComponents)
	if err != nil {
		erh.errorChan <- errors.Wrapf(err, "error while getting absolute path for generating relationships")
		return
	}

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if !info.IsDir() {
			var rel v1alpha1.RelationshipDefinition
			byt, err := os.ReadFile(path)
			if err != nil {
				erh.errorChan <- errors.Wrapf(err, fmt.Sprintf("unable to read file at %s", path))
				return nil
			}
			err = json.Unmarshal(byt, &rel)
			if err != nil {
				erh.errorChan <- errors.Wrapf(err, fmt.Sprintf("unmarshal json failed for %s", path))
				return nil
			}
			erh.relationshipChan <- rel
		}
		return nil
	})
	if err != nil {
		erh.errorChan <- errors.Wrapf(err, "error while generating relationships")
	}
}

// watches the component and relationship channels for incoming definitions and registers them with the registry manager
// If an error occurs, it logs the error
func (erh *EntityRegistrationHelper) watchComponents(ctx context.Context) {
	var err error
	meshmodel.Mutex.Lock()
	for {
		select {
		case comp := <-erh.componentChan:
			err = erh.regManager.RegisterEntity(meshmodel.Host{
				Hostname: ArtifactHubComponentsHandler.String(),
			}, comp)
		case rel := <-erh.relationshipChan:
			err = erh.regManager.RegisterEntity(meshmodel.Host{
				Hostname: ArtifactHubComponentsHandler.String(),
			}, rel)

		//Watching and logging errors from error channel
		case mhErr := <-erh.errorChan:
			if err != nil {
				erh.log.Error(mhErr)
			}

		case <-ctx.Done():
			registryLog(erh.regManager)
			meshmodel.Mutex.Unlock()
			return
		}

		if err != nil {
			erh.errorChan <- err
		}
	}
}
func writeToFile(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
func registryLog(regManager *meshmodel.RegistryManager) {
	logLevel := viper.GetInt("LOG_LEVEL")
	if viper.GetBool("DEBUG") {
		logLevel = int(logrus.DebugLevel)
	}
	// Initialize Logger instance
	log, err := logger.New("meshery", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: logLevel,
	})

	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	hosts, _, err := regManager.GetRegistrants(&v1alpha1.HostFilter{})
	if err != nil {
		log.Error(err)
	}
	for _, host := range hosts {
		summary := host.Summary
		addNonZero := func(a, b int64) int64 {
			if b != 0 {
				return a + b
			}
			return a
		}
		totalModels := addNonZero(0, summary.Models)
		totalComponents := addNonZero(0, summary.Components)
		totalRelationships := addNonZero(0, summary.Relationships)
		totalPolicies := addNonZero(0, summary.Policies)

		log.Info(fmt.Sprintf("For registrant %s successfully imported %d models %d components %d relationships %d policy",
			host.Hostname, totalModels, totalComponents, totalRelationships, totalPolicies))

		failedMsg, _ := meshmodel.FailedMsgCompute("", host.Hostname)
		if failedMsg != "" {
			log.Error(meshmodel.ErrRegisteringEntity(failedMsg, host.Hostname))
		}
	}

	filePath := "register_attempts.json"
	jsonData, err := json.MarshalIndent(meshmodel.RegisterAttempts, "", "  ")
	if err != nil {
		log.Error(meshmodel.ErrMarshalingRegisteryAttempts(err))
		return
	}

	err = writeToFile(filePath, jsonData)
	if err != nil {
		log.Error(meshmodel.ErrWritingRegisteryAttempts(err))
		return
	}
}
