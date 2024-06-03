package helpers

import (
	"fmt"
	"strconv"
	"strings"

	gofrs "github.com/gofrs/uuid"
	"github.com/google/uuid"
	"github.com/layer5io/meshery/server/models"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/spf13/viper"

	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
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
}

type logRegistryHandler struct {
	NonImportModel   map[string]v1beta1.EntitySummary
	RegisterAttempts map[string]*EntityTypeCountWithErrors
}

var LogHandler logRegistryHandler

func init() {
	LogHandler = logRegistryHandler{
		NonImportModel:   make(map[string]v1beta1.EntitySummary),
		RegisterAttempts: make(map[string]*EntityTypeCountWithErrors),
	}
}

func HandleError(h v1beta1.Host, en entity.Entity, err error, isModelError bool, isRegistrantError bool) {
	if LogHandler.RegisterAttempts == nil {
		LogHandler.RegisterAttempts = make(map[string]*EntityTypeCountWithErrors)
	}

	if LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)] == nil {
		LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)] = &EntityTypeCountWithErrors{
			Model:        make(map[string]EntityErrorCount),
			Component:    make(map[string]EntityErrorCount),
			Relationship: make(map[uuid.UUID]EntityErrorCount),
			Policy:       make(map[uuid.UUID]EntityErrorCount),
			Registry:     make(map[string]EntityErrorCount),
		}
	}

	switch entity := en.(type) {
	case *v1beta1.ComponentDefinition:
		entityName := entity.DisplayName
		isAnnotation, _ := entity.Metadata["isAnnotation"].(bool)
		if entity.Component.Schema == "" && !isAnnotation && err == nil {
			err = meshmodel.ErrEmptySchema()
		}
		if err != nil {
			handleModelOrRegistrantError(h, entity.Model.Name, err, isModelError, isRegistrantError)

			if entityCount, ok := LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Component[entityName]; ok {
				entityCount.Attempt++
				LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Component[entityName] = entityCount
			} else {
				LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Component[entityName] = EntityErrorCount{Attempt: 1, Error: err}
			}

			if LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Component[entityName].Attempt == 1 {
				currentValue := LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)]
				currentValue.Components++
				LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)] = currentValue
			}
		}

	case *v1alpha2.RelationshipDefinition:
		if err != nil {
			handleModelOrRegistrantError(h, entity.Model.Name, err, isModelError, isRegistrantError)

			if entityCount, ok := LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Relationship[entity.ID]; ok {
				entityCount.Attempt++
				LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Relationship[entity.ID] = entityCount
			} else {
				LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Relationship[entity.ID] = EntityErrorCount{Attempt: 1, Error: err}
			}

			if LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Relationship[entity.ID].Attempt == 1 {
				currentValue := LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)]
				currentValue.Relationships++
				LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)] = currentValue
			}
		}
	}
}

func handleModelOrRegistrantError(h v1beta1.Host, modelName string, err error, isModelError, isRegistrantError bool) {
	if LogHandler.RegisterAttempts == nil {
		LogHandler.RegisterAttempts = make(map[string]*EntityTypeCountWithErrors)
	}

	if LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)] == nil {
		LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)] = &EntityTypeCountWithErrors{
			Model:    make(map[string]EntityErrorCount),
			Registry: make(map[string]EntityErrorCount),
		}
	}

	switch {
	case isModelError:
		if entityCount, ok := LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Model[modelName]; ok {
			entityCount.Attempt++
			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Model[modelName] = entityCount
		} else {
			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Model[modelName] = EntityErrorCount{Attempt: 1, Error: err}
		}

		if LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Model[modelName].Attempt == 1 {
			currentValue := LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)]
			currentValue.Models++
			LogHandler.NonImportModel[meshmodel.HostnameToPascalCase(h.Hostname)] = currentValue
		}
	case isRegistrantError:
		if entityCount, ok := LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Registry[meshmodel.HostnameToPascalCase(h.Hostname)]; ok {
			entityCount.Attempt++
			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Registry[meshmodel.HostnameToPascalCase(h.Hostname)] = entityCount
		} else {
			LogHandler.RegisterAttempts[meshmodel.HostnameToPascalCase(h.Hostname)].Registry[meshmodel.HostnameToPascalCase(h.Hostname)] = EntityErrorCount{Attempt: 1}
		}
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
func FailedEventCompute(hostname string, mesheryInstanceID gofrs.UUID, provider *models.Provider) (string, error) {
	failedMsg, err := FailedMsgCompute("", hostname)
	if err != nil {
		return "", err
	}
	if failedMsg != "" {
		filePath := viper.GetString("REGISTRY_LOG_FILE")
		errorEventBuilder := events.NewEvent().FromUser(mesheryInstanceID).FromSystem(mesheryInstanceID).WithCategory("kubernetes_components").WithAction("registration")
		errorEventBuilder.WithSeverity(events.Error).WithDescription(failedMsg)
		errorEvent := errorEventBuilder.Build()
		errorEventBuilder.WithMetadata(map[string]interface{}{
			"LongDescription":      fmt.Sprintf("The import process for a registrant %s encountered difficulties,due to which %s. Specific issues during the import process resulted in certain entities not being successfully registered in the table.", "Kubernetes", failedMsg),
			"SuggestedRemediation": fmt.Sprintf("Visit docs with the error code %s", "https://docs.meshery.io/reference/error-codes"),
			"DownloadLink":         filePath,
			"ViewLink":             filePath,
		})
		_ = (*provider).PersistEvent(errorEvent)

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
		if !isEmpty(nonEmptyEntity) {
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
func isEmpty(e EntityTypeCountWithErrors) bool {
	return len(e.Model) == 0 && len(e.Component) == 0 && len(e.Relationship) == 0 && len(e.Policy) == 0 && len(e.Registry) == 0
}
