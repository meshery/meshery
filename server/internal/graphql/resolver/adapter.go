package resolver

import (
	"context"
	"fmt"

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

func (r *Resolver) changeAdapterStatus(ctx context.Context, _ models.Provider, targetStatus model.Status, adapterName, targetPort string) (model.Status, error) {
	// not able to perform any operation when the name is not there
	if adapterName == "" && targetPort == "" {
		return model.StatusUnknown, ErrAdapterInsufficientInformation(fmt.Errorf("adapter name or targetport or both are missing"))
	}

	// in case of empty target, prefer the default ports
	if targetPort == "" {
		r.Log.Warn(fmt.Errorf("target port is not provided, looking for default ports"))
		selectedAdapter := getAdapterInformationByName(adapterName)
		if selectedAdapter == nil {
			return model.StatusUnknown, ErrAdapterInsufficientInformation(fmt.Errorf("adapter name is not available, not able to figure out target port"))
		}

		targetPort = selectedAdapter.Location
	}

	deleteAdapter := true

	if targetStatus == model.StatusEnabled {
		r.Log.Info("Deploying Adapter")
		deleteAdapter = false
	} else {
		r.Log.Info("Undeploying Adapter")
	}

	r.Log.Debug(fmt.Printf("changing adapter status for %s on port %s to %s \n", adapterName, targetPort, targetStatus))

	adapter := models.Adapter{Name: adapterName, Location: fmt.Sprintf("%s:%s", adapterName, targetPort)}
	go func(ctx context.Context, del bool) {
		var operation string
		if del {
			operation = "Undeploy"
			err = r.Config.AdapterTracker.UndeployAdapter(ctx, adapter)
		} else {
			operation = "Deploy"
			err = r.Config.AdapterTracker.DeployAdapter(ctx, adapter)
		}
		if err != nil {
			r.Log.Info("Failed to " + operation + " adapter")
			r.Log.Error(err)
		} else {
			r.Log.Info("Successfully " + operation + "ed adapter")
		}
	}(ctx, deleteAdapter)

	return model.StatusProcessing, nil
}
