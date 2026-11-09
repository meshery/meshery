package resolver

import (
	"context"
	"fmt"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
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

// changeAdapterStatus deploys or undeploys an adapter given its name and port.
// It verifies adapter parameters, tracker availability, and guards logger invocations.
func (r *Resolver) changeAdapterStatus(ctx context.Context, _ models.Provider, targetStatus model.Status, adapterName, targetPort string) (model.Status, error) {
	// not able to perform any operation when the name is not there
	if adapterName == "" && targetPort == "" {
		return model.StatusUnknown, ErrAdapterInsufficientInformation(fmt.Errorf("adapter name or targetport or both are missing"))
	}

	// in case of empty target, prefer the default ports
	if targetPort == "" {
		if r.Log != nil {
			r.Log.Warn(fmt.Errorf("target port is not provided, looking for default ports"))
		}
		selectedAdapter := getAdapterInformationByName(adapterName)
		if selectedAdapter == nil {
			return model.StatusUnknown, ErrAdapterInsufficientInformation(fmt.Errorf("adapter name is not available, not able to figure out target port"))
		}

		targetPort = selectedAdapter.Location
	}

	deleteAdapter := true

	if targetStatus == model.StatusEnabled {
		if r.Log != nil {
			r.Log.Info("Deploying Adapter")
		}
		deleteAdapter = false
	} else {
		if r.Log != nil {
			r.Log.Info("Undeploying Adapter")
		}
	}

	if r.Log != nil {
		r.Log.Debug(fmt.Sprintf("changing adapter status for %s on port %s to %s \n", adapterName, targetPort, targetStatus))
	}

	if r.Config == nil || r.Config.AdapterTracker == nil {
		return model.StatusUnknown, ErrAdapterInsufficientInformation(fmt.Errorf("adapter tracker is not available"))
	}

	var err error
	adapter := models.Adapter{Name: adapterName, Location: fmt.Sprintf("%s:%s", adapterName, targetPort)}
	var operation string
	if deleteAdapter {
		operation = "Undeploy"
		err = r.Config.AdapterTracker.UndeployAdapter(ctx, adapter)
	} else {
		operation = "Deploy"
		err = r.Config.AdapterTracker.DeployAdapter(ctx, adapter)
	}
	if err != nil {
		if r.Log != nil {
			r.Log.Info("Failed to " + operation + " adapter")
			r.Log.Error(err)
		}
		return model.StatusUnknown, err
	}

	if r.Log != nil {
		r.Log.Info(operation + "ed adapter")
	}

	return model.StatusProcessing, nil
}
