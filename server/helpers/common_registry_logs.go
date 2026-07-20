package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/meshery/schemas/models/core"

	gofrs "github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	mutils "github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta3/component"

	"github.com/spf13/viper"

	"github.com/meshery/meshkit/models/events"
	_models "github.com/meshery/meshkit/models/meshmodel/core/v1beta1"
	entity "github.com/meshery/meshkit/models/meshmodel/entity"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
)

type EntityErrorCount struct {
	Attempt int
	Error   error
}

type EntityTypeCountWithErrors struct {
	Model        map[string]EntityErrorCount
	Component    map[string]EntityErrorCount
	Relationship map[core.Uuid]EntityErrorCount
	Policy       map[core.Uuid]EntityErrorCount
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

// recordAttempt increments the attempt count for key in m, recording err as
// the entry's Error only on the first attempt (matching the zero value a
// missing map key already returns, so no presence check is needed).
func recordAttempt[K comparable](m map[K]EntityErrorCount, key K, err error) {
	ec := m[key]
	ec.Attempt++
	if ec.Attempt == 1 {
		ec.Error = err
	}
	m[key] = ec
}

func HandleError(c connection.Connection, en entity.Entity, err error, isModelError bool, isRegistrantError bool) {
	if LogHandler.RegisterAttempts == nil {
		LogHandler.RegisterAttempts = make(map[string]*EntityTypeCountWithErrors)
	}

	hostKey := meshmodel.HostnameToPascalCase(c.Kind)
	if LogHandler.RegisterAttempts[hostKey] == nil {
		LogHandler.RegisterAttempts[hostKey] = &EntityTypeCountWithErrors{
			Model:        make(map[string]EntityErrorCount),
			Component:    make(map[string]EntityErrorCount),
			Relationship: make(map[core.Uuid]EntityErrorCount),
			Policy:       make(map[core.Uuid]EntityErrorCount),
			Registry:     make(map[string]EntityErrorCount),
		}
	}
	attempts := LogHandler.RegisterAttempts[hostKey]

	switch entity := en.(type) {
	case *component.ComponentDefinition:
		entityName := "[ " + entity.Model.Name + " " + entity.Model.Model.Version + " ]" + "( " + entity.DisplayName + " )"
		isAnnotation := entity.Metadata.IsAnnotation
		if entity.Component.Schema == "" && !isAnnotation && err == nil {
			err = meshmodel.ErrEmptySchema()
		}
		if err != nil {
			handleModelOrRegistrantError(c, entity.Model.Name, err, isModelError, isRegistrantError)

			attempts.mu.Lock()
			recordAttempt(attempts.Component, entityName, err)
			attempts.mu.Unlock()

			if attempts.Component[entityName].Attempt == 1 {
				currentValue := LogHandler.NonImportModel[hostKey]
				currentValue.Components++
				LogHandler.NonImportModel[hostKey] = currentValue
			}
		}

	case *relationship.RelationshipDefinition:
		if err != nil {
			handleModelOrRegistrantError(c, entity.Model.Name, err, isModelError, isRegistrantError)

			attempts.mu.Lock()
			recordAttempt(attempts.Relationship, entity.GetID(), err)
			attempts.mu.Unlock()

			if attempts.Relationship[entity.GetID()].Attempt == 1 {
				currentValue := LogHandler.NonImportModel[hostKey]
				currentValue.Relationships++
				LogHandler.NonImportModel[hostKey] = currentValue
			}
		}
	}
}

// handleModelOrRegistrantError is only ever called from HandleError, above,
// after HandleError has already ensured LogHandler.RegisterAttempts[hostKey]
// is initialized - so it can rely on that instead of repeating the check.
func handleModelOrRegistrantError(c connection.Connection, modelName string, err error, isModelError, isRegistrantError bool) {
	hostKey := meshmodel.HostnameToPascalCase(c.Kind)
	attempts := LogHandler.RegisterAttempts[hostKey]

	switch {
	case isModelError:
		attempts.mu.Lock()
		recordAttempt(attempts.Model, modelName, err)
		attempts.mu.Unlock()

		if attempts.Model[modelName].Attempt == 1 {
			currentValue := LogHandler.NonImportModel[hostKey]
			currentValue.Models++
			LogHandler.NonImportModel[hostKey] = currentValue
		}
	case isRegistrantError:
		attempts.mu.Lock()
		recordAttempt(attempts.Registry, hostKey, err)
		attempts.mu.Unlock()
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

func FailedEventCompute(hostname string, mesheryInstanceID core.Uuid, provider *models.Provider, userID string, ec *models.Broadcast) (string, error) {

	failedMsg, err := FailedMsgCompute("", hostname)
	if err != nil {
		return "", err
	}
	if failedMsg != "" {
		filePath := viper.GetString("REGISTRY_LOG_FILE")
		errorEventBuilder := events.NewEvent().FromOwner(mesheryInstanceID).FromSystem(mesheryInstanceID).WithCategory("registration").WithAction("get_summary")
		errorEventBuilder.WithSeverity(events.Error).WithDescription(failedMsg)
		errorEvent := errorEventBuilder.Build()
		errorEventBuilder.WithMetadata(map[string]interface{}{
			"Long_Description": fmt.Sprintf("One or more entities failed to register. The import process for registrant, %s, encountered the following issue: %s.", hostname, failedMsg),
			"DownloadLink":     filePath,
			"ViewLink":         filePath,
		})
		_ = (*provider).PersistSystemEvent(*errorEvent)
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
		fmt.Fprintf(&logMessage, "%s failed to register:\n  Components:\n", host)
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
func filterUUIDEmpty(m map[core.Uuid]EntityErrorCount) map[core.Uuid]EntityErrorCount {
	result := make(map[core.Uuid]EntityErrorCount)
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
