package registration

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/spf13/viper"
)

// packaingUnit is the representation of the atomic unit that can be registered into meshery server
type packagingUnit struct {
	model v1beta1.Model
	components []v1beta1.ComponentDefinition
	relationships []v1alpha2.RelationshipDefinition
	_ []v1beta1.PolicyDefinition
}

type RegistrationHelper struct {
	handlerConfig    *models.HandlerConfig
	regManager       *meshmodel.RegistryManager
	log              logger.Handler
}

func NewRegistrationHelper(log logger.Handler, hc *models.HandlerConfig, regm *meshmodel.RegistryManager) RegistrationHelper {
	return RegistrationHelper{handlerConfig: hc, log: log, regManager: regm}
}

/*
	Register will accept a RegisterableEntity (dir, tar or oci for now).
	Errors are written to the log file.
*/
func (rh *RegistrationHelper) Register(entity RegisterableEntity) {
	// get the packaging units
	pu, err := entity.PkgUnit()
	if(err != nil){
		// given input is not a valid model, or could not walk the directory
		return
	}
	// fmt.Printf("Packaging Unit: Model name: %s, comps: %d, rels: %d\n", pu.model.Name, len(pu.components), len(pu.relationships))
	rh.register(pu)
}


/*
	register will return an error if it is not able to register the `model`.
	If there are errors when registering other entities, they are handled properly but does not stop the registration process.
*/
func (rh *RegistrationHelper)register(pkg packagingUnit) {
	// 1. Register the model
	model := pkg.model

	// Dont register anything else if registrant is not there
	if(model.Registrant.Hostname == ""){
		err := ErrMissingRegistrant(model.Name)
		RegLog.InsertEntityRegFailure(model.Registrant.Hostname, "",entity.Model, model.Name, err)
		return
	}
	_, _, err := rh.regManager.RegisterEntity(
		v1beta1.Host{Hostname: model.Registrant.Hostname,},
		&model,
		)

	// If model cannot be registered, don't register anything else
	if err != nil {
		err = ErrRegisterEntity(err, string(model.Type()), model.DisplayName)
		RegLog.InsertEntityRegFailure(model.Registrant.Hostname, "",entity.Model, model.Name, err)
		return
	}

	hostname := model.Registrant.Hostname
	modelName := model.Name
	// 2. Register components
	for _, comp := range pkg.components {
		comp.Model = model
		_, _, err := rh.regManager.RegisterEntity(
		v1beta1.Host{Hostname: hostname,},
		&comp,
		)
	if err != nil {
		err = ErrRegisterEntity(err, string(comp.Type()), comp.DisplayName)
		RegLog.InsertEntityRegFailure(hostname, modelName ,entity.ComponentDefinition, comp.DisplayName, err)
		rh.log.Error(err)
	}
	}

	// 3. Register relationships
	for _, rel := range pkg.relationships {
		rel.Model = model
		_, _, err := rh.regManager.RegisterEntity(v1beta1.Host{
			Hostname: hostname,
		}, &rel)
		if err != nil {
			err = ErrRegisterEntity(err, string(rel.Type()), rel.Kind)
			rh.log.Error(err)
			RegLog.InsertEntityRegFailure(hostname, modelName ,entity.RelationshipDefinition, rel.ID.String(), err)
		}
	}
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

		failLog, err := FailedEventCompute(host.Hostname, sysID, &provider, "", erh.handlerConfig.EventBroadcaster)
		if err != nil {
			log.Error(err)
		}
		if failLog != "" {
			log.Error(meshmodel.ErrRegisteringEntity(failLog, host.Hostname))
		}

	}
	err = writeLogsToFiles()
	if err != nil {
		log.Error(err)
	}
}