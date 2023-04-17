package resolver

import (
	"context"

	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *Resolver) changeAdapterStatus(_ context.Context, _ models.Provider, targetStatus model.Status, adapterName, targetPort string) (model.Status, error) {
	platform, _ := utils.GetPlatform()
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

	r.Log.Debug(adapterName, targetPort)
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
