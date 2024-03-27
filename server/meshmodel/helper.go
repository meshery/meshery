package meshmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/gofrs/uuid"
	google "github.com/google/uuid"

	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var ArtifactHubComponentsHandler = meshmodel.ArtifactHub{} //The components generated in output directory will be handled by kubernetes
var ModelsPath = "../meshmodel"
var RelativeRelationshipsPath = "relationships"
var NonImportModel = make(map[string]v1alpha1.EntitySummary) // Count the quantity of error entity
var RegisterAttempts models.EntityCountWithErrors            //actual error and attempts we took to register the entity

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
	// Reset RegisterAttempts before processing components and relationships
	RegisterAttempts = models.EntityCountWithErrors{
		Model:        make(map[string]models.EntityErrorCount),
		Registry:     make(map[string]models.EntityErrorCount),
		Component:    make(map[string]models.EntityErrorCount),
		Relationship: make(map[google.UUID]models.EntityErrorCount),
		Policy:       make(map[google.UUID]models.EntityErrorCount),
	}

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
	var isModelError bool
	var isRegistrantError bool
	for {
		select {
		case comp := <-erh.componentChan:
			isModelError, isRegistrantError, err = erh.regManager.RegisterEntity(meshmodel.Host{
				Hostname: ArtifactHubComponentsHandler.String(),
			}, comp)

			handleError(meshmodel.Host{
				Hostname: ArtifactHubComponentsHandler.String(),
			}, comp, err, isModelError, isRegistrantError)

		case rel := <-erh.relationshipChan:
			isModelError, isRegistrantError, err = erh.regManager.RegisterEntity(meshmodel.Host{
				Hostname: ArtifactHubComponentsHandler.String(),
			}, rel)

			handleError(meshmodel.Host{
				Hostname: ArtifactHubComponentsHandler.String(),
			}, rel, err, isModelError, isRegistrantError)

		//Watching and logging errors from error channel
		case mhErr := <-erh.errorChan:
			if err != nil {
				fmt.Println(mhErr)
			}

		case <-ctx.Done():
			erh.registryLog(ctx)
			return
		}
		//Is this necessary because this what the current logic is doing as well this and mhErr but for schema less we need something else and if any error this just stops
		if err != nil {
			erh.errorChan <- err
		}
	}
}

// handle error for entity
func handleError(h meshmodel.Host, en meshmodel.Entity, err error, isModelError, isRegistrantError bool) {

	switch entity := en.(type) {
	case v1alpha1.ComponentDefinition:
		entityName := entity.DisplayName
		isAnnotation, _ := entity.Metadata["isAnnotation"].(bool)
		if entity.Schema == "" && !isAnnotation {
			//If the schema is empty or isAnnotation is false we give it an emptySchema error and handle it manually
			err = meshmodel.ErrEmptySchema()
		}
		if err != nil {
			// Check if the component registration attempt exists in RegisterAttempts.
			// If it does, increment the attempt count. If not, add a new entry with attempt count as 1 and store the error.
			handleModelOrRegistrantError(h, entity.Model.Name, err, isModelError, isRegistrantError)
			//
			if entityCount, ok := RegisterAttempts.Component[entityName]; ok {
				entityCount.Attempt++
				RegisterAttempts.Component[entityName] = entityCount
			} else {
				RegisterAttempts.Component[entityName] = models.EntityErrorCount{Attempt: 1, Error: err}
			}
			// If this is the first attempt for the component, increment the component count in NonImportModel.

			if RegisterAttempts.Component[entityName].Attempt == 1 {
				currentValue := NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)]
				currentValue.Components++
				NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)] = currentValue
			}
		}

	case v1alpha1.RelationshipDefinition:
		if err != nil {
			handleModelOrRegistrantError(h, entity.Model.Name, err, isModelError, isRegistrantError)
			entityID := entity.ID
			if entityCount, ok := RegisterAttempts.Relationship[entityID]; ok {
				entityCount.Attempt++
				RegisterAttempts.Relationship[entityID] = entityCount
			} else {
				RegisterAttempts.Relationship[entityID] = models.EntityErrorCount{Attempt: 1, Error: err}
			}

			if RegisterAttempts.Relationship[entityID].Attempt == 1 {
				currentValue := NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)]
				currentValue.Relationships++
				NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)] = currentValue
			}
		}

	}

}

// handle common error for model and registrant for both component and relationship
func handleModelOrRegistrantError(h meshmodel.Host, modelName string, err error, isModelError, isRegistrantError bool) {

	switch {
	case isModelError:
		if entityCount, ok := RegisterAttempts.Model[modelName]; ok {
			entityCount.Attempt++
			RegisterAttempts.Model[modelName] = entityCount
		} else {
			RegisterAttempts.Model[modelName] = models.EntityErrorCount{Attempt: 1, Error: err}
		}

		if RegisterAttempts.Model[modelName].Attempt == 1 {
			currentValue := NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)]
			currentValue.Models++
			NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)] = currentValue
		}
	case isRegistrantError:
		if entityCount, ok := RegisterAttempts.Registry[meshmodel.HostnameToPascalCase(h.Hostname)]; ok {
			entityCount.Attempt++
			RegisterAttempts.Registry[meshmodel.HostnameToPascalCase(h.Hostname)] = entityCount
		} else {
			RegisterAttempts.Registry[meshmodel.HostnameToPascalCase(h.Hostname)] = models.EntityErrorCount{Attempt: 1}
		}

	}

}

// To log the event to Terminal and send an event to Ui with id 00000000-0000-0000-0000-000000000000
func (erh *EntityRegistrationHelper) registryLog(ctx context.Context) {
	log := erh.log
	fileRoute := path.Join(viper.GetString("SERVER_CONTENT_FOLDER"), "entities")
	errDir := os.MkdirAll(fileRoute, 0755)
	if errDir != nil {
		log.Error(meshmodel.ErrCreatingUserDataDirectory(viper.GetString("SERVER_CONTENT_FOLDER")))
		os.Exit(1)
	}
	filePath := fileRoute + "/entities.json"

	systemID := viper.GetString("INSTANCE_ID")

	sysID := uuid.FromStringOrNil(systemID)
	ownerID := "00000000-0000-0000-0000-000000000000"
	ownerUserID := uuid.FromStringOrNil(ownerID)

	hosts, _, err := erh.regManager.GetRegistrants(&v1alpha1.HostFilter{})
	if err != nil {
		log.Error(err)
	}
	eventBuilder := events.NewEvent().FromSystem(sysID).FromUser(ownerUserID).WithCategory("entity").WithAction("get_summary")
	for _, host := range hosts {

		successMessage := fmt.Sprintf("For registrant %s successfully imported", host.Hostname)
		appendIfNonZero := func(value int64, label string) {
			if value != 0 {
				successMessage += fmt.Sprintf(" %d %s", value, label)
			}
		}
		appendIfNonZero(host.Summary.Models, "models")
		appendIfNonZero(host.Summary.Components, "components")
		appendIfNonZero(host.Summary.Relationships, "relationships")
		appendIfNonZero(host.Summary.Policies, "policies")

		log.Info(successMessage)
		eventBuilder.WithMetadata(map[string]interface{}{
			"Hostname": host.Hostname,
		})
		eventBuilder.WithSeverity(events.Informational).WithDescription(successMessage)
		successEvent := eventBuilder.Build()
		provider := erh.handlerConfig.Providers["Meshery"]
		_ = provider.PersistEvent(successEvent)

		failedMsg, _ := FailedMsgCompute("", host.Hostname)
		if failedMsg != "" {
			log.Error(meshmodel.ErrRegisteringEntity(failedMsg, host.Hostname))
			errorEventBuilder := events.NewEvent().FromUser(ownerUserID).FromSystem(sysID).WithCategory("entity").WithAction("get_summary")
			errorEventBuilder.WithSeverity(events.Error).WithDescription(failedMsg)
			errorEvent := errorEventBuilder.Build()
			errorEventBuilder.WithMetadata(map[string]interface{}{
				"LongDescription":      fmt.Sprintf("The import process for a registrant %s encountered difficulties,due to which %s. Specific issues during the import process resulted in certain entities not being successfully registered in the table.", host.Hostname, failedMsg),
				"SuggestedRemediation": fmt.Sprintf("Visit docs with the error code %s", "https://docs.meshery.io/reference/error-codes"),
				"DownloadLink":         filePath,
				"ViewLink":             filePath,
			})
			provider := erh.handlerConfig.Providers["Meshery"]
			_ = provider.PersistEvent(errorEvent)
		}

	}

	jsonData, err := json.Marshal(RegisterAttempts)
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

// To compute failed Msg for each host
func FailedMsgCompute(failedMsg string, hostName string) (string, error) {
	nonImportModel, exists := NonImportModel[hostName]
	if !exists {

		return "", meshmodel.ErrUnknownHostInMap()
	}

	if nonImportModel.Models > 0 || nonImportModel.Components > 0 || nonImportModel.Relationships > 0 || nonImportModel.Policies > 0 {
		failedMsg = "Failed to import"
		appendIfNonZero := func(msg string, count int64, entityName string) string {
			if count > 0 {
				return fmt.Sprintf("%s %d %s", msg, count, entityName)
			}
			return msg
		}

		failedMsg = appendIfNonZero(failedMsg, nonImportModel.Models, "models")
		failedMsg = appendIfNonZero(failedMsg, nonImportModel.Components, "components")
		failedMsg = appendIfNonZero(failedMsg, nonImportModel.Relationships, "relationships")
		failedMsg = appendIfNonZero(failedMsg, nonImportModel.Policies, "policies")
	}
	return failedMsg, nil
}

// write the json file to SERVER_CONTENT_FOLDER
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
