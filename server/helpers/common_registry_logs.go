package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	gofrs "github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"

	"github.com/spf13/viper"

	"github.com/layer5io/meshkit/models/events"
	_models "github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	entity "github.com/layer5io/meshkit/models/meshmodel/entity"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
)

type EntityErrorCount struct {
	Attempt int
	Error   error
}

type EntityTypeCountWithErrors struct {
	Model        map[string]EntityErrorCount
	Component    map[string]EntityErrorCount
	Relationship map[uuid.UUID]EntityErrorCount
	Policy       map[uuid.UUID]EntityErrorCount
	Registry     map[string]EntityErrorCount
	mu           sync.Mutex
}

type logRegistryHandler struct {
	NonImportModel   map[string]_models.EntitySummary
	RegisterAttempts map[string]*EntityTypeCountWithErrors
}

var LogHandler logRegistryHandler

func init() {
	LogHandler = logRegistryHandler{
		NonImportModel:   make(map[string]_models.EntitySummary),
		RegisterAttempts: make(map[string]*EntityTypeCountWithErrors),
	}
}

func HandleError(c connection.Connection, en entity.Entity, err error, isModelError bool, isRegistrantError bool) {
	if LogHandler.RegisterAttempts == nil {
		LogHandler.RegisterAttempts = make(map[string]*EntityTypeCountWithErrors)
	}

	if LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)] == nil {
		LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)] = &EntityTypeCountWithErrors{
			Model:        make(map[string]EntityErrorCount),
			Component:    make(map[string]EntityErrorCount),
			Relationship: make(map[uuid.UUID]EntityErrorCount),
			Policy:       make(map[uuid.UUID]EntityErrorCount),
			Registry:     make(map[string]EntityErrorCount),
		}
	}

	switch entity := en.(type) {
	case *component.ComponentDefinition:
		entityName := "[ " + entity.Model.Name + " " + entity.Model.Model.Version + " ]" + "( " + entity.DisplayName + " )"
		isAnnotation := entity.Metadata.IsAnnotation
		if entity.Component.Schema == "" && !isAnnotation && err == nil {
			err = meshmodel.ErrEmptySchema()
		}
		if err != nil {
			handleModelOrRegistrantError(c, entity.Model.Name, err, isModelError, isRegistrantError)

			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].mu.Lock()
			if entityCount, ok := LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Component[entityName]; ok {
				entityCount.Attempt++
				LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Component[entityName] = entityCount
			} else {
				LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Component[entityName] = EntityErrorCount{Attempt: 1, Error: err}
			}
			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].mu.Unlock()

			if LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Component[entityName].Attempt == 1 {
				currentValue := LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(c.Kind)]
				currentValue.Components++
				LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(c.Kind)] = currentValue
			}
		}

	case *relationship.RelationshipDefinition:
		if err != nil {
			handleModelOrRegistrantError(c, entity.Model.Name, err, isModelError, isRegistrantError)

			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].mu.Lock()
			if entityCount, ok := LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Relationship[entity.GetID()]; ok {
				entityCount.Attempt++
				LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Relationship[entity.Id] = entityCount
			} else {
				LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Relationship[entity.Id] = EntityErrorCount{Attempt: 1, Error: err}
			}
			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].mu.Unlock()

			if LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Relationship[entity.Id].Attempt == 1 {
				currentValue := LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(c.Kind)]
				currentValue.Relationships++
				LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(c.Kind)] = currentValue
			}
		}
	}
}

func handleModelOrRegistrantError(c connection.Connection, modelName string, err error, isModelError, isRegistrantError bool) {
	if LogHandler.RegisterAttempts == nil {
		LogHandler.RegisterAttempts = make(map[string]*EntityTypeCountWithErrors)
	}

	if LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)] == nil {
		LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)] = &EntityTypeCountWithErrors{
			Model:    make(map[string]EntityErrorCount),
			Registry: make(map[string]EntityErrorCount),
		}
	}

	switch {
	case isModelError:
		LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].mu.Lock()
		if entityCount, ok := LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Model[modelName]; ok {
			entityCount.Attempt++
			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Model[modelName] = entityCount
		} else {
			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Model[modelName] = EntityErrorCount{Attempt: 1, Error: err}
		}
		LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].mu.Unlock()

		if LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Model[modelName].Attempt == 1 {
			currentValue := LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(c.Kind)]
			currentValue.Models++
			LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(c.Kind)] = currentValue
		}
	case isRegistrantError:
		LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].mu.Lock()
		if entityCount, ok := LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Registry[meshmodel.HostnameToPascalCase(c.Kind)]; ok {
			entityCount.Attempt++
			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Registry[meshmodel.HostnameToPascalCase(c.Kind)] = entityCount
		} else {
			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].Registry[meshmodel.HostnameToPascalCase(c.Kind)] = EntityErrorCount{Attempt: 1}
		}
		LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(c.Kind)].mu.Unlock()
	}
}

func FailedMsgCompute(failedMsg string, hostName string) (string, error) {
	nonImportModel := LogHandler.NonImportModel[hostName]

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

func FailedEventCompute(hostname string, mesheryInstanceID gofrs.UUID, provider *models.Provider, userID string, ec *models.Broadcast) (string, error) {

	failedMsg, err := FailedMsgCompute("", hostname)
	if err != nil {
		return "", err
	}
	if failedMsg != "" {
		filePath := viper.GetString("REGISTRY_LOG_FILE")
		errorEventBuilder := events.NewEvent().FromUser(mesheryInstanceID).FromSystem(mesheryInstanceID).WithCategory("registration").WithAction("get_summary")
		errorEventBuilder.WithSeverity(events.Error).WithDescription(failedMsg)
		errorEvent := errorEventBuilder.Build()
		errorEventBuilder.WithMetadata(map[string]interface{}{
			"LongDescription":      fmt.Sprintf("One or more entities failed to register. The import process for registrant, %s, encountered the following issue: %s.", hostname, failedMsg),
			"SuggestedRemediation": fmt.Sprintf("Open Meshery Error Reference with error code %s", "https://docs.meshery.io/reference/error-codes"),
			"DownloadLink":         filePath,
			"ViewLink":             filePath,
		})
		_ = (*provider).PersistEvent(errorEvent)
		if userID != "" {
			userUUID := gofrs.FromStringOrNil(userID)
			ec.Publish(userUUID, errorEvent)

		}
	}
	return failedMsg, nil

}

func WriteLogsToFiles() error {
	filePath := viper.GetString("REGISTRY_LOG_FILE")

	// Initialize the formatted log message
	var logMessage strings.Builder

	// Marshal RegisterAttempts without empty entries
	nonEmptyRegisterAttempts := make(map[string]*EntityTypeCountWithErrors)
	for host, attempts := range LogHandler.RegisterAttempts {
		nonEmptyEntity := EntityTypeCountWithErrors{
			Model:        filterEmpty(attempts.Model),
			Component:    filterEmpty(attempts.Component),
			Relationship: filterUUIDEmpty(attempts.Relationship),
			Policy:       filterUUIDEmpty(attempts.Policy),
			Registry:     filterEmpty(attempts.Registry),
		}
		if !isEmpty(&nonEmptyEntity) {
			nonEmptyRegisterAttempts[host] = &nonEmptyEntity
		}
	}

	// Iterate over non-empty register attempts and construct the log message
	for host, attempts := range nonEmptyRegisterAttempts {
		logMessage.WriteString(fmt.Sprintf("%s failed to register:\n  Components:\n", host))
		for entityType, entityCount := range attempts.Component {
			logMessage.WriteString("    " + entityType + " (Attempt " + strconv.Itoa(entityCount.Attempt) + "): " + entityCount.Error.Error() + "\n")
		}
	}

	// Write the formatted log message to the file
	err := mutils.WriteToFile(filePath, logMessage.String())
	if err != nil {
		return meshmodel.ErrWritingRegisteryAttempts(err)
	}
	return nil
}

// filterEmpty removes empty entries from a map
func filterEmpty(m map[string]EntityErrorCount) map[string]EntityErrorCount {
	result := make(map[string]EntityErrorCount)
	for k, v := range m {
		if v.Attempt > 0 || v.Error != nil {
			result[k] = v
		}
	}
	return result
}

// filterUUIDEmpty removes empty entries from a map with UUID keys
func filterUUIDEmpty(m map[uuid.UUID]EntityErrorCount) map[uuid.UUID]EntityErrorCount {
	result := make(map[uuid.UUID]EntityErrorCount)
	for k, v := range m {
		if v.Attempt > 0 || v.Error != nil {
			result[k] = v
		}
	}
	return result
}

// isEmpty checks if an EntityTypeCountWithErrors is empty
func isEmpty(e *EntityTypeCountWithErrors) bool {
	return len(e.Model) == 0 && len(e.Component) == 0 && len(e.Relationship) == 0 && len(e.Policy) == 0 && len(e.Registry) == 0
}
