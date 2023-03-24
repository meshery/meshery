package resolver

import (
	"context"
	"os"
	"os/exec"

	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *Resolver) changeAdapterStatus(ctx context.Context, provider models.Provider, adapter models.Adapter, status model.Status, adaptorPort string) (model.Status, error) {
	deleteAdaptor := true

	if status == model.StatusEnabled {
		r.Log.Info("Deploying Adapter")
		deleteAdaptor = false
	} else {
		r.Log.Info("Undeploying Adapter")
	}

	go func(del bool) {
		// Update AdapterTracker
		if del {
			r.Config.AdapterTracker.RemoveAdapter(ctx, adapter)
		} else {
			r.Config.AdapterTracker.AddAdapter(ctx, adapter)
		}

		// Clone the Meshery-Istio repository
		cmd := exec.Command("git", "clone", "https://github.com/layer5io/"+adapter.Name+".git")
		err := cmd.Run()
		if err != nil {
			r.Log.Error(err)
			return
		}

		// Navigate to the Meshery-Istio directory
		err = os.Chdir(adapter.Name)
		if err != nil {
			r.Log.Error(err)
		}

		// Run the Meshery-Istio
		cmd = exec.Command("make", "run")
		err = cmd.Run()
		if err != nil {
			r.Log.Error(err)
		}

		if status == model.StatusEnabled {
			r.Log.Info("Deployed Adapter")
			deleteAdaptor = false
		} else {
			r.Log.Info("Undeployed Adapter")
		}
	}(deleteAdaptor)

	return model.StatusProcessing, nil
}
