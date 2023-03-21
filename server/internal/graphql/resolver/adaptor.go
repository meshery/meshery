package resolver

import (
	"context"
	"path"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"

	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *Resolver) changeAdapterStatus(ctx context.Context, provider models.Provider, adapter models.Adapter, status model.Status, ctxID string) (model.Status, error) {
	delete := true

	if status == model.StatusEnabled {
		r.Log.Info("Installing Adapter")
		delete = false
	} else {
		r.Log.Info("Uninstalling Adapter in context ", ctxID)
	}

	var kubeclient *mesherykube.Client
	var k8scontext models.K8sContext
	var err error
	if ctxID != "" {
		allContexts, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
		if !ok || len(allContexts) == 0 {
			r.Log.Error(ErrNilClient)
			return model.StatusUnknown, ErrNilClient
		}
		for _, ctx := range allContexts {
			if ctx.ID == ctxID {
				k8scontext = ctx
				break
			}
		}
	} else {
		k8scontexts, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
		if !ok || len(k8scontexts) == 0 {
			return model.StatusUnknown, ErrMesheryClient(nil)
		}
		k8scontext = k8scontexts[0]
	}

	kubeclient, err = k8scontext.GenerateKubeHandler()
	if err != nil {
		return model.StatusUnknown, ErrMesheryClient(err)
	}

	if kubeclient.KubeClient == nil {
		r.Log.Error(ErrNilClient)
		return model.StatusUnknown, ErrNilClient
	}

	go func(del bool, kubeclient *mesherykube.Client) {
		var act mesherykube.HelmChartAction
		if del {
			r.Config.AdapterTracker.RemoveAdapter(ctx, adapter)
			act = mesherykube.UNINSTALL
		} else {
			r.Config.AdapterTracker.AddAdapter(ctx, adapter)
			act = mesherykube.INSTALL
		}
		err := kubeclient.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
			Namespace:       utils.MesheryNamespace,
			ReleaseName:     "meshery",
			CreateNamespace: true,
			ChartLocation: mesherykube.HelmChartLocation{
				Repository: utils.HelmChartURL,
				Chart:      utils.HelmChartName,
				Version:    "latest",
			},
			OverrideValues: map[string]interface{}{},
			Action:         act,
			// the helm chart will be downloaded to ~/.meshery/manifests if it doesn't exist
			DownloadLocation: path.Join(utils.MesheryFolder, utils.ManifestsFolder),
			DryRun:           true,
		})

		if err != nil {
			r.Log.Error(err)
			return
		}

		r.Log.Info("Adapter operation executed")

		if !del {
			endpoint, err := model.SubscribeToBroker(provider, kubeclient, r.brokerChannel, r.BrokerConn, connectionTrackerSingleton)
			r.Log.Debug("Endpoint: ", endpoint)
			if err != nil {
				r.Log.Error(err)
				return
			}
			connectionTrackerSingleton.Set(k8scontext.ID, endpoint)
			r.Log.Info("Connected to broker at:", endpoint)
			connectionTrackerSingleton.Log(r.Log)
		}

		r.Log.Info("Meshsync operation executed")
	}(delete, kubeclient)

	return model.StatusProcessing, nil
}
