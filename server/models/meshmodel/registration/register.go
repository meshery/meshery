package registration

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/meshmodel/core"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/spf13/viper"
)

type PackagingUnit struct {
	model v1beta1.Model
	components []v1beta1.ComponentDefinition
	relationships []v1alpha2.RelationshipDefinition
	policies []v1beta1.PolicyDefinition
}

type RegistrationHelper struct {
	handlerConfig    *models.HandlerConfig
	regManager       *meshmodel.RegistryManager
	log              logger.Handler
}

func NewRegistrationHelper(log logger.Handler, hc *models.HandlerConfig, regm *meshmodel.RegistryManager) RegistrationHelper {
	return RegistrationHelper{handlerConfig: hc, log: log, regManager: regm}
}

// Register will accept a dir, oci or tar. Register assumes that whatever is given in its argument is a single unit of
// packaging (model) and interpret it like that. inputs that have models inside models are not considered valid. This is
// enforce the fact that registration should happen through the unit of packaing (model)
// Given that, it does not care about the folder structure as long as we have one file defining the `model` and others
// `components` etc.
// For bulk registration, use BulkRegister
// Using any because of lack of type strength in go
func (rh *RegistrationHelper) Register(entity RegisterableEntity) error {
	// get the packaging units
	pu, err := entity.PkgUnit()
	if(err != nil){
		// given input is not a valid model, or it should not be packaged like this
		return err
	}
	return rh.register(pu)
}

// This function will take meshery's unit of packaging, a model in any format (tar.gz, oci, or a directory) and register
// it
func (rh *RegistrationHelper)register(pkg PackagingUnit) error {
	// 1. Register the model
	model := pkg.model
	_, _, err := rh.regManager.RegisterEntity(
		v1beta1.Host{Hostname: model.Registrant.Hostname,},
		&model,
		)
	if err != nil {
		err = core.ErrRegisterEntity(err, string(model.Type()), model.DisplayName)
		return err
	}

	// 2. Register components
	for _, comp := range pkg.components {
		// TODO: After we make sure that the models are accurately generated (by artifacthub specifically), we have to use pu.model as the model for
		// components
		isRegistrantError, isModelError, err := rh.regManager.RegisterEntity(
		v1beta1.Host{Hostname: comp.Model.Registrant.Hostname,},
		&comp,
		)
	if err != nil {
		err = core.ErrRegisterEntity(err, string(comp.Type()), comp.DisplayName)
		rh.log.Error(err)
	}
	helpers.HandleError(v1beta1.Host{
				Hostname: comp.Model.Registrant.Hostname,
			}, &comp, err, isModelError, isRegistrantError)
	}

	// 3. Register relationships
	for _, rel := range pkg.relationships {
		rel.Model = model
		isRegistrantError, isModelError, err := rh.regManager.RegisterEntity(v1beta1.Host{
			Hostname: rel.Model.Registrant.Hostname,
		}, &rel)
		if err != nil {
			err = core.ErrRegisterEntity(err, string(rel.Type()), rel.Kind)
			rh.log.Error(err)
			helpers.HandleError(v1beta1.Host{
			Hostname: rel.Model.Registrant.Hostname,
			}, &rel, err, isModelError, isRegistrantError)
		}
	}

	return nil
}

func (erh *RegistrationHelper) RegistryLog() {
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

		failLog, err := helpers.FailedEventCompute(host.Hostname, sysID, &provider, "", erh.handlerConfig.EventBroadcaster)
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























