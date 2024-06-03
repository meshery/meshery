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

	systemID := viper.GetString("INSTANCE_ID")

	sysID := uuid.FromStringOrNil(systemID)
	hosts, _, err := erh.regManager.GetRegistrants(&v1beta1.HostFilter{})
	if err != nil {
		log.Error(err)
	}

	for _, host := range hosts {

		eventBuilder := events.NewEvent().FromSystem(sysID).FromUser(sysID).WithCategory("entity").WithAction("get_summary")
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

		failLog, err := helpers.FailedEventCompute(host.Hostname, sysID, &provider)
		if err != nil {
			log.Error(err)
		}
		if failLog != "" {
			log.Error(meshmodel.ErrRegisteringEntity(failLog, host.Hostname))
		}

	}
	err = helpers.WriteLogsToFiles()
	if err != nil {
		log.Error(err)
	}
}
