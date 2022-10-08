package resolver

import (
	"context"
	"io"
	"os"
	"path"
	"time"

	operatorClient "github.com/layer5io/meshery-operator/pkg/client"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/utils"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

// Global singleton instance of k8s connection tracker to map Each K8sContext to a unique Broker URL
var connectionTrackerSingleton = model.NewK8sConnctionTracker()
var (
	MeshSyncSubscriptionError = model.Error{
		Description: "Failed to get MeshSync data",
		Code:        ErrResolverMeshsyncSubscriptionCode,
	}
	MeshSyncMesheryClientMissingError = model.Error{
		Code:        ErrResolverMeshsyncSubscriptionCode,
		Description: "Cannot find Meshery Client",
	}
)

func (r *Resolver) listenToMeshSyncEvents(ctx context.Context, provider models.Provider, k8scontextIDs []string) (<-chan *model.OperatorControllerStatusPerK8sContext, error) {
	channel := make(chan *model.OperatorControllerStatusPerK8sContext)
	if r.brokerChannel == nil {
		r.brokerChannel = make(chan *broker.Message)
	}

	k8sctxs, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
	if !ok || len(k8sctxs) == 0 {
		return nil, ErrNilClient
	}
	var k8sContexts []models.K8sContext
	if len(k8scontextIDs) == 1 && k8scontextIDs[0] == "all" {
		k8sContexts = k8sctxs
	} else if len(k8scontextIDs) != 0 {
		var k8sContextIDsMap = make(map[string]bool)
		for _, k8sContext := range k8scontextIDs {
			k8sContextIDsMap[k8sContext] = true
		}
		for _, k8Context := range k8sctxs {
			if k8sContextIDsMap[k8Context.ID] {
				k8sContexts = append(k8sContexts, k8Context)
			}
		}
	}
	for _, k8sctx := range k8sContexts {
		prevStatus := r.getMeshSyncStatus(k8sctx)
		go func(k8sctx models.K8sContext, ch chan *model.OperatorControllerStatusPerK8sContext) {
			r.Log.Info("Initializing MeshSync subscription")
			if r.MeshSyncChannelPerK8sContext[k8sctx.ID] == nil {
				r.MeshSyncChannelPerK8sContext[k8sctx.ID] = make(chan struct{})
			}
			go model.PersistClusterNames(ctx, r.Log, provider.GetGenericPersister(), r.MeshSyncChannelPerK8sContext[k8sctx.ID])
			go model.ListernToEvents(r.Log, provider.GetGenericPersister(), r.brokerChannel, r.MeshSyncChannelPerK8sContext[k8sctx.ID], r.Broadcast)
			for {
				status := r.getMeshSyncStatus(k8sctx)
				statusWithContext := model.OperatorControllerStatusPerK8sContext{
					ContextID:                k8sctx.ID,
					OperatorControllerStatus: &status,
				}
				ch <- &statusWithContext
				if status != prevStatus {
					prevStatus = status
				}
				time.Sleep(10 * time.Second)
			}
		}(k8sctx, channel)
	}

	return channel, nil
}

func (r *Resolver) getMeshSyncStatus(k8sctx models.K8sContext) model.OperatorControllerStatus {
	var status model.OperatorControllerStatus
	kubeclient, err := k8sctx.GenerateKubeHandler()
	if err != nil {
		r.Log.Error(ErrNilClient)
		return model.OperatorControllerStatus{
			Name:    "",
			Version: "",
			Status:  model.StatusDisabled,
			Error:   &MeshSyncMesheryClientMissingError,
		}
	}

	mesheryclient, err := operatorClient.New(&kubeclient.RestConfig)
	if err != nil {
		return model.OperatorControllerStatus{
			Name:    "",
			Version: "",
			Status:  model.StatusDisabled,
			Error:   &MeshSyncMesheryClientMissingError,
		}
	}

	status, err = model.GetMeshSyncInfo(mesheryclient, kubeclient)

	if err != nil {
		return model.OperatorControllerStatus{
			Name:    "",
			Version: "",
			Status:  model.StatusDisabled,
			Error:   &MeshSyncSubscriptionError,
		}
	}

	return status
}

func (r *Resolver) resyncCluster(ctx context.Context, provider models.Provider, actions *model.ReSyncActions, k8scontextID string) (model.Status, error) {
	if actions.ClearDb == "true" {
		// copies the contents .meshery/config/mesherydb.sql to .meshery/config/.archive/mesherydb.sql
		// then drops all the DB table and then migrate/create tables, missing foreign keys, constraints, columns and indexes.
		if actions.HardReset == "true" {
			mesherydbPath := path.Join(utils.GetHome(), ".meshery/config")
			err := os.Mkdir(path.Join(mesherydbPath, ".archive"), os.ModePerm)
			if err != nil && os.IsNotExist(err) {
				return "", err
			}

			src := path.Join(mesherydbPath, "mesherydb.sql")
			dst := path.Join(mesherydbPath, ".archive/mesherydb.sql")

			fin, err := os.Open(src)
			if err != nil {
				return "", err
			}
			defer fin.Close()

			fout, err := os.Create(dst)
			if err != nil {
				return "", err
			}
			defer fout.Close()

			_, err = io.Copy(fout, fin)
			if err != nil {
				return "", err
			}

			dbHandler := provider.GetGenericPersister()
			if dbHandler == nil {
				return "", ErrEmptyHandler
			}

			dbHandler.Lock()
			defer dbHandler.Unlock()

			r.Log.Info("Dropping Meshery Database")
			err = dbHandler.Migrator().DropTable(
				&meshsyncmodel.KeyValue{},
				&meshsyncmodel.Object{},
				&meshsyncmodel.ResourceSpec{},
				&meshsyncmodel.ResourceStatus{},
				&meshsyncmodel.ResourceObjectMeta{},
				&models.PerformanceProfile{},
				&models.MesheryResult{},
				&models.MesheryPattern{},
				&models.MesheryFilter{},
				&models.PatternResource{},
				&models.MesheryApplication{},
				&models.UserPreference{},
				&models.PerformanceTestConfig{},
				&models.SmiResultWithID{},
				models.K8sContext{},
			)
			if err != nil {
				r.Log.Error(err)
				return "", err
			}

			r.Log.Info("Migrating Meshery Database")
			err = dbHandler.AutoMigrate(
				&meshsyncmodel.KeyValue{},
				&meshsyncmodel.Object{},
				&meshsyncmodel.ResourceSpec{},
				&meshsyncmodel.ResourceStatus{},
				&meshsyncmodel.ResourceObjectMeta{},
				&models.PerformanceProfile{},
				&models.MesheryResult{},
				&models.MesheryPattern{},
				&models.MesheryFilter{},
				&models.PatternResource{},
				&models.MesheryApplication{},
				&models.UserPreference{},
				&models.PerformanceTestConfig{},
				&models.SmiResultWithID{},
				models.K8sContext{},
			)
			if err != nil {
				r.Log.Error(err)
				return "", err
			}

			r.Log.Info("Hard reset successfully completed")
		} else { //Delete meshsync objects coming from a particular cluster
			k8sctxs, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
			if !ok || len(k8sctxs) == 0 {
				r.Log.Error(ErrEmptyCurrentK8sContext)
				return "", ErrEmptyCurrentK8sContext
			}
			var sid string
			for _, k8ctx := range k8sctxs {
				if k8ctx.ID == k8scontextID && k8ctx.KubernetesServerID != nil {
					sid = k8ctx.KubernetesServerID.String()
					break
				}
			}
			if provider.GetGenericPersister() == nil {
				return "", ErrEmptyHandler
			}
			err := provider.GetGenericPersister().Where("cluster_id = ?", sid).Delete(&meshsyncmodel.Object{}).Error
			if err != nil {
				return "", ErrEmptyHandler
			}
		}
	}

	if actions.ReSync == "true" {
		if r.BrokerConn.Info() != broker.NotConnected {
			err := r.BrokerConn.Publish(model.RequestSubject, &broker.Message{
				Request: &broker.RequestObject{
					Entity: broker.ReSyncDiscoveryEntity,
				},
			})
			if err != nil {
				return "", ErrPublishBroker(err)
			}
		}
	}
	return model.StatusProcessing, nil
}

func (r *Resolver) connectToBroker(ctx context.Context, provider models.Provider, ctxID string) error {
	status, err := r.getOperatorStatus(ctx, provider, ctxID)
	if err != nil {
		return err
	}
	var currContext *models.K8sContext
	var newContextFound bool
	if ctxID == "" {
		currContexts, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
		if !ok || len(currContexts) == 0 {
			r.Log.Error(ErrNilClient)
			return ErrNilClient
		}
		currContext = &currContexts[0]
	} else {
		allContexts, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
		if !ok || len(allContexts) == 0 {
			r.Log.Error(ErrNilClient)
			return ErrNilClient
		}
		for _, ctx := range allContexts {
			if ctx.ID == ctxID {
				currContext = &ctx
				break
			}
		}
	}
	if currContext == nil {
		r.Log.Error(ErrNilClient)
		return ErrNilClient
	}
	if connectionTrackerSingleton.Get(currContext.ID) == "" {
		newContextFound = true
	}
	kubeclient, err := currContext.GenerateKubeHandler()
	if err != nil {
		r.Log.Error(ErrNilClient)
		return ErrNilClient
	}
	if (r.BrokerConn.IsEmpty() || newContextFound) && status != nil && status.Status == model.StatusEnabled {
		endpoint, err := model.SubscribeToBroker(provider, kubeclient, r.brokerChannel, r.BrokerConn, connectionTrackerSingleton)
		if err != nil {
			r.Log.Error(ErrAddonSubscription(err))
			return err
		}
		r.Log.Info("Connected to broker at:", endpoint)
		connectionTrackerSingleton.Set(currContext.ID, endpoint)
		connectionTrackerSingleton.Log(r.Log)
		return nil
	}

	if r.BrokerConn.Info() == broker.NotConnected {
		return ErrBrokerNotConnected
	}

	return nil
}

func (r *Resolver) deployMeshsync(ctx context.Context, provider models.Provider, ctxID string) (model.Status, error) {
	//err := model.RunMeshSync(r.Config.KubeClient, false)
	return model.StatusProcessing, nil
}

func (r *Resolver) connectToNats(ctx context.Context, provider models.Provider, k8scontextID string) (model.Status, error) {
	err := r.connectToBroker(ctx, provider, k8scontextID)
	if err != nil {
		return model.StatusDisabled, err
	}
	return model.StatusConnected, nil
}
