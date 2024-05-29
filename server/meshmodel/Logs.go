package meshmodel

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"

	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/spf13/viper"
)

func (erh *EntityRegistrationHelper) registryLog() {
	log := erh.log
	provider := erh.handlerConfig.Providers["None"]
	filePath := viper.GetString("REGISTRY_LOG_FILE")

	systemID := viper.GetString("INSTANCE_ID")

	sysID := uuid.FromStringOrNil(systemID)
	hosts, _, err := erh.regManager.GetRegistrants(&v1beta1.HostFilter{})
	if err != nil {
		log.Error(err)
	}
	eventBuilder := events.NewEvent().FromSystem(sysID).FromUser(sysID).WithCategory("entity").WithAction("get_summary")
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
		_ = provider.PersistEvent(successEvent)

		failedMsg, _ := helpers.FailedMsgCompute("", host.Hostname)
		if failedMsg != "" {
			log.Error(meshmodel.ErrRegisteringEntity(failedMsg, host.Hostname))
			errorEventBuilder := events.NewEvent().FromUser(sysID).FromSystem(sysID).WithCategory("entity").WithAction("get_summary")
			errorEventBuilder.WithSeverity(events.Error).WithDescription(failedMsg)
			errorEvent := errorEventBuilder.Build()
			errorEventBuilder.WithMetadata(map[string]interface{}{
				"LongDescription":      fmt.Sprintf("The import process for a registrant %s encountered difficulties,due to which %s. Specific issues during the import process resulted in certain entities not being successfully registered in the table.", host.Hostname, failedMsg),
				"SuggestedRemediation": fmt.Sprintf("Visit docs with the error code %s", "https://docs.meshery.io/reference/error-codes"),
				"DownloadLink":         filePath,
				"ViewLink":             filePath,
			})
			_ = provider.PersistEvent(errorEvent)
		}

	}
	err = helpers.WriteLogsToFiles()
	if err != nil {
		log.Error(err)
	}
}
