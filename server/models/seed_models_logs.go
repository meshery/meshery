package models

import (
	"fmt"

	gofrs "github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/spf13/viper"
	"strings"
)

var TAB = "    "

// representation of what could not be registered
type RegistrationFailureLog struct {
	// {'artifacthub':{'modelname': 'component' : {'jaegar': error}}}
	// for models, the structure will be like this:
	// {'artifacthub': {'': 'model': {'model': error}}}
	// think of it like every model belongs to a global model called ''
	failureData map[string](map[string](map[string](map[string]error)))

	// invalid definitions while parsing the given input (oci, tar, dir) into meshmodel entities
	invalidDefinitions map[string](error)
}
type EntityRegError struct {
	HostName   string
	ModelName  string
	EntityType entity.EntityType
	EntityName string
	Err        error
}

func NewRegistrationFailureLogHandler() *RegistrationFailureLog {
	return &RegistrationFailureLog{
		failureData:        make(map[string](map[string](map[string](map[string]error)))),
		invalidDefinitions: make(map[string]error),
	}
}

func (rfl *RegistrationFailureLog) AddInvalidDefinition(path string, err error) {
	rfl.invalidDefinitions[path] = err
}

func (rfl *RegistrationFailureLog) InsertEntityRegError(hostname string, modelName string, entityType entity.EntityType, name string, err error) {
	if rfl.failureData[hostname] == nil {
		rfl.failureData[hostname] = make(map[string]map[string]map[string]error)
	}
	if rfl.failureData[hostname][modelName] == nil {
		rfl.failureData[hostname][modelName] = make(map[string]map[string]error)
	}
	if rfl.failureData[hostname][modelName][string(entityType)] == nil {
		rfl.failureData[hostname][modelName][string(entityType)] = make(map[string]error)
	}
	rfl.failureData[hostname][modelName][string(entityType)][name] = err
}

func (rfl *RegistrationFailureLog) GetNonRegisteredEntitiesCount() map[string]v1beta1.EntitySummary {
	res := make(map[string]v1beta1.EntitySummary)
	for registrant, modelNamespacedData := range rfl.failureData {
		entitySummary := v1beta1.EntitySummary{}
		for modelName, entityTypeNamespacedData := range modelNamespacedData {
			if modelName == "" {
				modelsCount := len(entityTypeNamespacedData[string(entity.Model)])
				entitySummary.Models = int64(modelsCount)
				continue
			}
			entitySummary.Components = int64(len(entityTypeNamespacedData[string(entity.ComponentDefinition)]))
			entitySummary.Relationships = int64(len(entityTypeNamespacedData[string(entity.RelationshipDefinition)]))
			entitySummary.Policies = int64(len(entityTypeNamespacedData[string(entity.PolicyDefinition)]))
		}
		res[registrant] = entitySummary
	}
	return res
}

func failedMsgCompute(failedMsg string, hostName string, regLog *RegistrationFailureLog) (string, error) {
	nonImportModel := regLog.GetNonRegisteredEntitiesCount()[hostName]

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

func FailedEventCompute(hostname string, mesheryInstanceID gofrs.UUID, provider *Provider, userID string, ec *Broadcast, regErrorStore *RegistrationFailureLog) (string, error) {
	failedMsg, err := failedMsgCompute("", hostname, regErrorStore)
	if err != nil {
		return "", err
	}
	if failedMsg != "" {
		filePath := viper.GetString("REGISTRY_LOG_FILE")
		errorEventBuilder := events.NewEvent().FromUser(mesheryInstanceID).FromSystem(mesheryInstanceID).WithCategory("registration").WithAction("get_summary")
		errorEventBuilder.WithSeverity(events.Error).WithDescription(failedMsg)
		errorEvent := errorEventBuilder.Build()
		errorEventBuilder.WithMetadata(map[string]interface{}{
			"DownloadLink": filePath,
			"ViewLink":     filePath,
			"error":        ErrImportFailure(hostname, failedMsg),
		})
		_ = (*provider).PersistEvent(errorEvent)
		if userID != "" {
			userUUID := gofrs.FromStringOrNil(userID)
			ec.Publish(userUUID, errorEvent)

		}
	}
	return failedMsg, nil
}

func writeLogsToFiles(regLog *RegistrationFailureLog) error {
	filePath := viper.GetString("REGISTRY_LOG_FILE")
	// Initialize the formatted log message
	var logMessage strings.Builder

	logMessage.WriteString("Invalid Definitions: \n")
	for path, err := range regLog.invalidDefinitions {
		logMessage.WriteString(TAB + "[" + path + "] " + err.Error() + " \n")
	}
	logMessage.WriteString("---------------------------------------- \n")
	for host, modelNamespacedData := range regLog.failureData {
		logMessage.WriteString(fmt.Sprintf("%s failed to register:\n", host))
		for modelName, entityTypeNamespacedData := range modelNamespacedData {
			if modelName == "" {
				logMessage.WriteString(fmt.Sprintf("%sModels:\n", TAB))
				for name, modelData := range entityTypeNamespacedData[string(entity.Model)] {
					logMessage.WriteString(TAB + TAB + name + ": " + modelData.Error() + "\n")
				}
				continue
			}
			logMessage.WriteString(fmt.Sprintf("%sFor Model %s:\n", TAB, modelName))
			for entityType, entityNameNamespacedData := range entityTypeNamespacedData {
				logMessage.WriteString(fmt.Sprintf("%s%s%s:\n", TAB, TAB, entityType))
				for name, err := range entityNameNamespacedData {
					logMessage.WriteString(fmt.Sprintf("%s%s%s%s: %s\n", TAB, TAB, TAB, name, err.Error()))
				}

			}
		}
	}
	// Write the formatted log message to the file
	err := mutils.WriteToFile(filePath, logMessage.String())
	if err != nil {
		return meshmodel.ErrWritingRegisteryAttempts(err)
	}
	return nil
}

func RegistryLog(log logger.Handler, handlerConfig *HandlerConfig, regManager *meshmodel.RegistryManager, regErrorStore *RegistrationFailureLog) {
	provider := handlerConfig.Providers["None"]

	systemID := viper.GetString("INSTANCE_ID")

	sysID := gofrs.FromStringOrNil(systemID)
	hosts, _, err := regManager.GetRegistrants(&v1beta1.HostFilter{})
	if err != nil {
		log.Error(err)
	}

	for _, host := range hosts {
		eventBuilder := events.NewEvent().FromSystem(sysID).FromUser(sysID).WithCategory("entity").WithAction("get_summary")
		successMessage := fmt.Sprintf("For registrant %s successfully imported", host.Kind)
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
			"kind": host.Kind,
		})
		eventBuilder.WithSeverity(events.Informational).WithDescription(successMessage)
		successEvent := eventBuilder.Build()
		_ = provider.PersistEvent(successEvent)

		failLog, err := FailedEventCompute(host.Kind, sysID, &provider, "", handlerConfig.EventBroadcaster, regErrorStore)
		if err != nil {
			log.Error(err)
		}
		if failLog != "" {
			log.Error(meshmodel.ErrRegisteringEntity(failLog, host.Kind))
		}

	}
	err = writeLogsToFiles(regErrorStore)
	if err != nil {
		log.Error(err)
	}
}

func (rfl *RegistrationFailureLog) GetEntityRegErrors() []EntityRegError {
	var errors []EntityRegError
	for host, modelData := range rfl.failureData {
		for modelName, entityTypeData := range modelData {
			for entityTypeStr, entityNameData := range entityTypeData {
				entityType := entity.EntityType(entityTypeStr)
				for entityName, err := range entityNameData {
					errors = append(errors, EntityRegError{
						HostName:   host,
						ModelName:  modelName,
						EntityType: entityType,
						EntityName: entityName,
						Err:        err,
					})
				}
			}
		}
	}
	return errors
}
