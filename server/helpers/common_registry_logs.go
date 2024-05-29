package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/spf13/viper"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	entity "github.com/layer5io/meshkit/models/meshmodel/entity"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
)

type EntityErrorCount struct {
	Attempt int
	Error   error
}

type EntityCountWithErrors struct {
	Model        map[string]EntityErrorCount
	Registry     map[string]EntityErrorCount
	Component    map[string]EntityErrorCount
	Relationship map[uuid.UUID]EntityErrorCount
	Policy       map[uuid.UUID]EntityErrorCount
}

type logRegistryHandler struct {
	NonImportModel   map[string]v1beta1.EntitySummary
	RegisterAttempts *EntityCountWithErrors
}

var LogHandler logRegistryHandler

// init initializes the logRegistryHandler with default values.
func init() {
	// Initialize NonImportModel as an empty map
	nonImportModel := make(map[string]v1beta1.EntitySummary)

	// Initialize RegisterAttempts with empty maps for each entity type
	registerAttempts := &EntityCountWithErrors{
		Model:        make(map[string]EntityErrorCount),
		Registry:     make(map[string]EntityErrorCount),
		Component:    make(map[string]EntityErrorCount),
		Relationship: make(map[uuid.UUID]EntityErrorCount),
		Policy:       make(map[uuid.UUID]EntityErrorCount),
	}

	// Assign initialized values to the global LogHandler
	LogHandler = logRegistryHandler{
		NonImportModel:   nonImportModel,
		RegisterAttempts: registerAttempts,
	}
}
func HandleError(h v1beta1.Host, en entity.Entity, err error, isModelError bool, isRegistrantError bool) {

	switch entity := en.(type) {
	case *v1beta1.ComponentDefinition:
		entityName := entity.DisplayName
		isAnnotation, _ := entity.Metadata["isAnnotation"].(bool)
		if entity.Component.Schema == "" && !isAnnotation && err == nil {
			//If the schema is empty or isAnnotation is false we give it an emptySchema error and handle it manually
			err = meshmodel.ErrEmptySchema()
		}
		if err != nil {
			// Check if the component registration attempt exists in RegisterAttempts.
			// If it does, increment the attempt count. If not, add a new entry with attempt count as 1 and store the error.
			handleModelOrRegistrantError(h, entity.Model.Name, err, isModelError, isRegistrantError)
			//
			if entityCount, ok := LogHandler.RegisterAttempts.Component[entityName]; ok {
				entityCount.Attempt++
				LogHandler.RegisterAttempts.Component[entityName] = entityCount
			} else {
				LogHandler.RegisterAttempts.Component[entityName] = EntityErrorCount{Attempt: 1, Error: err}
			}
			// If this is the first attempt for the component, increment the component count in NonImportModel.

			if LogHandler.RegisterAttempts.Component[entityName].Attempt == 1 {
				currentValue := LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)]
				currentValue.Components++
				LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)] = currentValue
			}
		}

	case *v1alpha2.RelationshipDefinition:
		if err != nil {
			handleModelOrRegistrantError(h, entity.Model.Name, err, isModelError, isRegistrantError)
			entityID := entity.ID
			if entityCount, ok := LogHandler.RegisterAttempts.Relationship[entityID]; ok {
				entityCount.Attempt++
				LogHandler.RegisterAttempts.Relationship[entityID] = entityCount
			} else {
				LogHandler.RegisterAttempts.Relationship[entityID] = EntityErrorCount{Attempt: 1, Error: err}
			}

			if LogHandler.RegisterAttempts.Relationship[entityID].Attempt == 1 {
				currentValue := LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)]
				currentValue.Relationships++
				LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)] = currentValue
			}
		}

	}

}

// handle common error for model and registrant for both component and relationship
func handleModelOrRegistrantError(h v1beta1.Host, modelName string, err error, isModelError, isRegistrantError bool) {

	switch {
	case isModelError:
		if entityCount, ok := LogHandler.RegisterAttempts.Model[modelName]; ok {
			entityCount.Attempt++
			LogHandler.RegisterAttempts.Model[modelName] = entityCount
		} else {
			LogHandler.RegisterAttempts.Model[modelName] = EntityErrorCount{Attempt: 1, Error: err}
		}

		if LogHandler.RegisterAttempts.Model[modelName].Attempt == 1 {
			currentValue := LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)]
			currentValue.Models++
			LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)] = currentValue
		}
	case isRegistrantError:
		if entityCount, ok := LogHandler.RegisterAttempts.Registry[meshmodel.HostnameToPascalCase(h.Hostname)]; ok {
			entityCount.Attempt++
			LogHandler.RegisterAttempts.Registry[meshmodel.HostnameToPascalCase(h.Hostname)] = entityCount
		} else {
			LogHandler.RegisterAttempts.Registry[meshmodel.HostnameToPascalCase(h.Hostname)] = EntityErrorCount{Attempt: 1}
		}

	}

}

// To compute failed Msg for each host
func FailedMsgCompute(failedMsg string, hostName string) (string, error) {
	nonImportModel, exists := LogHandler.NonImportModel[hostName]
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
func WriteLogsToFiles() error {
	filePath := viper.GetString("REGISTRY_LOG_FILE")
	jsonData, err := json.Marshal(LogHandler.RegisterAttempts)
	if err != nil {

		return (meshmodel.ErrMarshalingRegisteryAttempts(err))
	}
	jsonDataString := string(jsonData)

	err = mutils.WriteToFile(filePath, jsonDataString)
	if err != nil {
		return (meshmodel.ErrWritingRegisteryAttempts(err))

	}
	return nil
}
