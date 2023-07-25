package resolver

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func getAdapterInformationByName(adapterName string) *models.Adapter {
	var adapter *models.Adapter

	for _, v := range models.ListAvailableAdapters {
		if adapterName == v.Name {
			adapter = &v
		}
	}

	return adapter
}

func (r *Resolver) changeAdapterStatus(_ context.Context, _ models.Provider, targetStatus model.Status, adapterName, targetPort string) (model.Status, error) {
	// not able to perform any operation when the name is not there
	if adapterName == "" && targetPort == "" {
		return model.StatusUnknown, helpers.ErrAdapterInsufficientInformation(fmt.Errorf("either of adapter name or target port is not provided, please provide a name of adapter or target-port to perform operation on"))
	}

	// in case of empty target, prefer the default ports
	if targetPort == "" {
		r.Log.Warn(fmt.Errorf("target port is not specified in the request body, searching for default ports"))
		selectedAdapter := getAdapterInformationByName(adapterName)
		if selectedAdapter == nil {
			return model.StatusUnknown, helpers.ErrAdapterInsufficientInformation(fmt.Errorf("adapter name is not available, not able to figure out target port"))
		}

		targetPort = selectedAdapter.Location
	}

	platform := utils.GetPlatform()
	if platform == "kubernetes" {
		r.Log.Info("Feature for kuberenetes disabled")
		return model.StatusDisabled, nil
	}

	deleteAdapter := true

	if targetStatus == model.StatusEnabled {
		r.Log.Info("Deploying Adapter")
		deleteAdapter = false
	} else {
		r.Log.Info("Undeploying Adapter")
	}

	r.Log.Debug(fmt.Printf("changing adapter status of %s on port %s to status %s \n", adapterName, targetPort, targetStatus))
	adapter := models.Adapter{Name: adapterName, Location: targetPort}
	go func(routineCtx context.Context, del bool) {
		var operation string
		if del {
			operation = "Undeploy"
			err = r.Config.AdapterTracker.UndeployAdapter(routineCtx, adapter)
		} else {
			operation = "Deploy"
			err = r.Config.AdapterTracker.DeployAdapter(routineCtx, adapter)
		}
		if err != nil {
			// r.Log.Error(errors.Errorf("Failed to "+operation+" adapter: %w", err))
			r.Log.Info("Failed to " + operation + " adapter")
			r.Log.Error(err)
		} else {
			r.Log.Info("Successfully " + operation + "ed adapter")
		}
	}(context.Background(), deleteAdapter)

	return model.StatusProcessing, nil
}