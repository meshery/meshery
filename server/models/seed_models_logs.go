package models

import (
	"fmt"
	"strings"

	gofrs "github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/models/meshmodel/core/v1beta1"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	mutils "github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/viper"
)

var TAB = "    "

// RegistrationFailureLog represents what could not be registered.
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

// FailedEventCompute reports the registration failures recorded for hostname
// and, when there are any, persists an error event describing them.
//
// persister is the sink for that event. It is a SystemEventPersister rather
// than a Provider because the failure event is raised outside any user request:
// the caller may be boot-time seeding, which has no provider to route through
// once PROVIDER enforcement has emptied the registration map.
//
// A failure to persist that event is returned, not swallowed, alongside the
// failure summary: both return values are independent, so the caller must
// report the error and still act on a non-empty summary.
func FailedEventCompute(hostname string, mesheryInstanceID core.Uuid, persister SystemEventPersister, userID string, ec *Broadcast, regErrorStore *RegistrationFailureLog) (string, error) {
	failedMsg, err := failedMsgCompute("", hostname, regErrorStore)
	if err != nil {
		return "", err
	}
	var persistErr error
	if failedMsg != "" {
		filePath := viper.GetString("REGISTRY_LOG_FILE")
		errorEventBuilder := events.NewEvent().FromOwner(mesheryInstanceID).FromSystem(mesheryInstanceID).WithCategory("registration").WithAction("get_summary")
		errorEventBuilder.WithSeverity(events.Error).WithDescription(failedMsg)
		errorEvent := errorEventBuilder.Build()
		errorEventBuilder.WithMetadata(map[string]interface{}{
			"DownloadLink": filePath,
			"ViewLink":     filePath,
			"error":        ErrImportFailure(hostname, failedMsg),
		})
		persistErr = persister.PersistSystemEvent(*errorEvent)
		if userID != "" {
			userUUID := gofrs.FromStringOrNil(userID)
			ec.Publish(userUUID, errorEvent)

		}
	}
	return failedMsg, persistErr
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
		fmt.Fprintf(&logMessage, "%s failed to register:\n", host)
		for modelName, entityTypeNamespacedData := range modelNamespacedData {
			if modelName == "" {
				fmt.Fprintf(&logMessage, "%sModels:\n", TAB)
				for name, modelData := range entityTypeNamespacedData[string(entity.Model)] {
					logMessage.WriteString(TAB + TAB + name + ": " + modelData.Error() + "\n")
				}
				continue
			}
			fmt.Fprintf(&logMessage, "%sFor Model %s:\n", TAB, modelName)
			for entityType, entityNameNamespacedData := range entityTypeNamespacedData {
				fmt.Fprintf(&logMessage, "%s%s%s:\n", TAB, TAB, entityType)
				for name, err := range entityNameNamespacedData {
					fmt.Fprintf(&logMessage, "%s%s%s%s: %s\n", TAB, TAB, TAB, name, err.Error())
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

// aggregateSummariesByKind merges the per-connection entity summaries returned
// by GetRegistrants into a single summary per registrant Kind, preserving the
// order in which each Kind is first encountered.
//
// A registrant is identified here by Kind (for example "artifacthub", "github",
// "meshery") rather than by connection, because the startup summary reports per
// Kind and registration failures are already tracked per Kind upstream - see
// meshkit's registration handler, which keys failures with
// InsertEntityRegError(model.Registrant.Kind, ...). When a single Kind is backed
// by more than one connection, GetRegistrants returns a row per connection;
// folding their counts together yields exactly one summary line per Kind.
func aggregateSummariesByKind(hosts []v1beta1.MeshModelHostsWithEntitySummary) (orderedKinds []string, summaries map[string]v1beta1.EntitySummary) {
	orderedKinds = make([]string, 0, len(hosts))
	summaries = make(map[string]v1beta1.EntitySummary)
	for _, host := range hosts {
		// A registrant with no Kind cannot be meaningfully named in the summary,
		// so skip it rather than emit a "For registrant  imported" line.
		if host.Kind == "" {
			continue
		}
		current, exists := summaries[host.Kind]
		if !exists {
			orderedKinds = append(orderedKinds, host.Kind)
		}
		current.Models += host.Summary.Models
		current.Components += host.Summary.Components
		current.Relationships += host.Summary.Relationships
		current.Policies += host.Summary.Policies
		summaries[host.Kind] = current
	}
	return orderedKinds, summaries
}

// formatRegistrantSummary renders the one-line import summary for a registrant
// Kind, e.g. "For registrant artifacthub imported 150 models, 1326 components.".
// GetRegistrants only returns registrants with at least one model, so the
// "0 entities" fallback is defensive: it keeps the line grammatical if that
// contract ever changes, instead of leaving a dangling "imported.".
func formatRegistrantSummary(kind string, summary v1beta1.EntitySummary) string {
	var parts []string
	if summary.Models != 0 {
		parts = append(parts, fmt.Sprintf("%d models", summary.Models))
	}
	if summary.Components != 0 {
		parts = append(parts, fmt.Sprintf("%d components", summary.Components))
	}
	if summary.Relationships != 0 {
		parts = append(parts, fmt.Sprintf("%d relationships", summary.Relationships))
	}
	if summary.Policies != 0 {
		parts = append(parts, fmt.Sprintf("%d policies", summary.Policies))
	}
	if len(parts) == 0 {
		parts = append(parts, "0 entities")
	}
	return fmt.Sprintf("For registrant %s imported %s.", kind, strings.Join(parts, ", "))
}

func RegistryLog(log logger.Handler, handlerConfig *HandlerConfig, regManager *meshmodel.RegistryManager, regErrorStore *RegistrationFailureLog) {
	// Seeding runs before any user request, so there is no provider to route
	// these events through - and, since PROVIDER enforcement, no guarantee
	// that any particular provider is even registered. This used to read
	// handlerConfig.Providers[LocalProviderName], which RestrictToEnforcedProvider
	// deletes on a pinned deployment; the zero Provider it returned then
	// panicked the seeding goroutine and, with it, the whole server
	// (meshery/meshery#21584).
	persister := handlerConfig.SystemEventPersister
	if persister == nil {
		// Reachable only if a caller built a HandlerConfig without wiring the
		// persister. Report it rather than dropping the events silently, then
		// keep going: the startup summaries below are still worth logging, and
		// a missing sink must not cost the operator the boot log too.
		log.Error(ErrNoSystemEventSink())
		persister = discardSystemEvents{}
	}

	systemID := viper.GetString("INSTANCE_ID")

	sysID := gofrs.FromStringOrNil(systemID)
	hosts, _, err := regManager.GetRegistrants(&v1beta1.HostFilter{})
	if err != nil {
		log.Error(err)
	}

	// GetRegistrants returns one row per registrant connection, so a Kind backed
	// by more than one connection would emit a duplicate summary line. Aggregate
	// by Kind so each registrant is reported exactly once in the startup logs.
	orderedKinds, kindSummaries := aggregateSummariesByKind(hosts)
	for _, kind := range orderedKinds {
		summary := kindSummaries[kind]
		eventBuilder := events.NewEvent().FromSystem(sysID).FromOwner(sysID).WithCategory("entity").WithAction("get_summary")
		successMessage := formatRegistrantSummary(kind, summary)

		log.Info(successMessage)
		eventBuilder.WithMetadata(map[string]interface{}{
			"kind":    kind,
			"doclink": "https://docs.meshery.io/concepts/logical#logical-concepts",
		})
		eventBuilder.WithSeverity(events.Informational).WithDescription(successMessage)
		successEvent := eventBuilder.Build()
		if err := persister.PersistSystemEvent(*successEvent); err != nil {
			log.Error(err)
		}

		failLog, err := FailedEventCompute(kind, sysID, persister, "", handlerConfig.EventBroadcaster, regErrorStore)
		if err != nil {
			log.Error(err)
		}
		if failLog != "" {
			log.Error(meshmodel.ErrRegisteringEntity(failLog, kind))
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

// discardSystemEvents is the fallback sink used when a HandlerConfig reaches
// RegistryLog without a SystemEventPersister. It drops events, which is only
// ever correct because the omission has already been reported as an error - it
// is not a silent skip.
type discardSystemEvents struct{}

func (discardSystemEvents) PersistSystemEvent(events.Event) error { return nil }
