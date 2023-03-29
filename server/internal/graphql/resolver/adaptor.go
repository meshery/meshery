package resolver

import (
	"context"

	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *Resolver) changeAdaptorStatus(ctx context.Context, provider models.Provider, targetStatus model.Status, adaptorPort string) (model.Status, error) {
	deleteAdaptor := true

	if targetStatus == model.StatusEnabled {
		r.Log.Info("Deploying Adapter")
		deleteAdaptor = false
	} else {
		r.Log.Info("Undeploying Adapter")
	}

	adapter := models.Adapter{Name: string(models.Istio), Location: adaptorPort}
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
			r.Log.Info(err.Error())
		} else {
			r.Log.Info("Successfully " + operation + "ed adapter")
		}
	}(context.Background(), deleteAdaptor)

	return model.StatusProcessing, nil
}
