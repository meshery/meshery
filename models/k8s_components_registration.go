package models

import (
	"context"

	"github.com/layer5io/meshkit/logger"
	"github.com/spf13/viper"
)

type RegistrationStatus int

const (
	RegistrationComplete RegistrationStatus = iota
	NotRegistered
	Registering
)

type ComponentsRegistrationHelper struct {
	// map that holds the registration status of each of the contexts in this runtime of the server
	// it should be private
	ctxRegStatusMap map[string]RegistrationStatus
	log             logger.Handler
}

func NewComponentsRegistrationHelper(logger logger.Handler) *ComponentsRegistrationHelper {
	return &ComponentsRegistrationHelper{
		ctxRegStatusMap: make(map[string]RegistrationStatus),
		log:             logger,
	}
}

// update the map with the given list of contexts
func (cg *ComponentsRegistrationHelper) UpdateContexts(ctxs []*K8sContext) *ComponentsRegistrationHelper {
	for _, ctx := range ctxs {
		ctxId := ctx.ID
		if _, ok := cg.ctxRegStatusMap[ctxId]; !ok {
			cg.ctxRegStatusMap[ctxId] = NotRegistered
		}
	}
	return cg
}

type k8sRegistrationFunction func(ctxt context.Context, config []byte, ctxId string) error

// start registration of components for the contexts
func (cg *ComponentsRegistrationHelper) RegisterComponents(ctxs []*K8sContext, regFunc k8sRegistrationFunction) {
	for _, ctx := range ctxs {
		ctxId := ctx.ID
		// do not do anything about the contexts that are not present in the ctxRegStatusMap
		if status, ok := cg.ctxRegStatusMap[ctxId]; ok {
			if !viper.GetBool("SKIP_COMP_GEN") {
				// only start registering components for contexts whose status is NotRegistered
				if status == NotRegistered {
					// update the status
					cg.ctxRegStatusMap[ctxId] = Registering
					cg.log.Info("registration of k8s native components started for contextID: ", ctxId)

					go func() {
						// set the status to RegistrationComplete
						defer func() {
							cg.ctxRegStatusMap[ctxId] = RegistrationComplete

							cg.log.Info("registration of k8s native components completed for contextID: ", ctxId)
						}()

						// start registration
						cfg, err := ctx.GenerateKubeConfig()
						if err != nil {
							cg.log.Error(err)
							return
						}
						err = regFunc(context.Background(), cfg, ctxId)
						if err != nil {
							cg.log.Error(err)
							return
						}
					}()

				}
			}
		}
	}
}
