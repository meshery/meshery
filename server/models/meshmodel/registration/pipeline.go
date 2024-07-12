package registration

import (
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/meshmodel/core"
	"github.com/layer5io/meshkit/logger"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
)

type RegistrationPipeline struct {
	handlerConfig    *models.HandlerConfig
	regManager       *meshmodel.RegistryManager
	log              logger.Handler

	pu PackagingUnit
	breakpipeline bool
}

// 1. show initialiation summary, like how many components are there to be registered.
// 2. Register model
// 3. Register components (use queue, to get the count of how many have been done)
// 4. Write SVG for file system for all components
// 5. Register relationships (use queue, to get the count of how many have been done)
func NewPipeline(log logger.Handler, hc *models.HandlerConfig, regm *meshmodel.RegistryManager, pu PackagingUnit) *RegistrationPipeline {
	return &RegistrationPipeline{
		handlerConfig: hc,
		regManager: regm,
		log: log,
		pu: pu,
		breakpipeline: false,
	}
}

func (rp *RegistrationPipeline) Start(errorChan chan error) {
	 rp.initSummary(errorChan).registerModel(errorChan).registerComponents(errorChan).writeComponentSvgs(errorChan).registerRelationships(errorChan)
}

func (rp *RegistrationPipeline) initSummary(errorChan chan error) *RegistrationPipeline{
	return rp
}
func (rp *RegistrationPipeline) registerModel(errorChan chan error) *RegistrationPipeline{
	if(rp.breakpipeline){
		return rp
	}
	model := rp.pu.model
	isRegistrantError, isModelError, err := rp.regManager.RegisterEntity(
		v1beta1.Host{Hostname: model.Registrant.Hostname,},
		&model,
		)
	if err != nil {
		err = core.ErrRegisterEntity(err, string(model.Type()), model.DisplayName)
		errorChan <- err
	}
	// have to see if this works, cuz it was used for comps
	helpers.HandleError(v1beta1.Host{
				Hostname: model.Registrant.Hostname,
			}, &model, err, isModelError, isRegistrantError)
	if(err != nil){
		rp.breakpipeline = true
	}

	return rp
}


func (rp *RegistrationPipeline) registerComponents(errorChan chan error) *RegistrationPipeline{
	if(rp.breakpipeline){
		return rp
	}

	for _, comp := range rp.pu.components {
		// TODO: After we make sure that the models are accurately generated (by artifacthub specifically), we have to use pu.model as the model for
		// components
		isRegistrantError, isModelError, err := rp.regManager.RegisterEntity(
		v1beta1.Host{Hostname: comp.Model.Registrant.Hostname,},
		&comp,
		)
	if err != nil {
		err = core.ErrRegisterEntity(err, string(comp.Type()), comp.DisplayName)
		errorChan <- err
	}
	// have to see if this works, cuz it was used for comps
	helpers.HandleError(v1beta1.Host{
				Hostname: comp.Model.Registrant.Hostname,
			}, &comp, err, isModelError, isRegistrantError)
	}
	return rp
}
func (rp *RegistrationPipeline) writeComponentSvgs(errorChan chan error) *RegistrationPipeline{
	if(rp.breakpipeline){
		return rp
	}
	for _, comp := range rp.pu.components {
		utils.WriteSVGsOnFileSystem(&comp)
	}
	return rp
}
func (rp *RegistrationPipeline) registerRelationships(errorChan chan error) *RegistrationPipeline{
	if(rp.breakpipeline){
		return rp
	}
	model := rp.pu.model
	for _, rel := range rp.pu.relationships {
		rel.Model = model

		isRegistrantError, isModelError, err := rp.regManager.RegisterEntity(v1beta1.Host{
			Hostname: rel.Model.Registrant.Hostname,
		}, &rel)
		helpers.HandleError(v1beta1.Host{
			Hostname: rel.Model.Registrant.Hostname,
		}, &rel, err, isModelError, isRegistrantError)
		if err != nil {
			err = core.ErrRegisterEntity(err, string(rel.Type()), rel.Kind)
			errorChan <- err
		}
	}

	return rp
}