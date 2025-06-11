package resolver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/gofrs/uuid"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"github.com/meshery/meshery/server/handlers"
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/utils"
	"github.com/spf13/viper"
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
				return "", model.ErrEmptyHandler
			}

			dbHandler.Lock()
			defer dbHandler.Unlock()

			r.Log.Info("Dropping Meshery Database")
			tables, err := dbHandler.Migrator().GetTables()
			if err != nil {
				r.Log.Error(ErrGormDatabase(err))
				return "", err
			}

			for _, table := range tables {
				if table == "events" {
					continue
				}
				if err := dbHandler.Migrator().DropTable(table); err != nil {
					r.Log.Error(ErrGormDatabase(err))
					return "", err
				}
			}

			r.Log.Info("Migrating Meshery Database")
			err = dbHandler.AutoMigrate(
				&meshsyncmodel.KubernetesKeyValue{},
				&meshsyncmodel.KubernetesResource{},
				&meshsyncmodel.KubernetesResourceSpec{},
				&meshsyncmodel.KubernetesResourceStatus{},
				&meshsyncmodel.KubernetesResourceObjectMeta{},
				&models.PerformanceProfile{},
				&models.MesheryResult{},
				&models.MesheryPattern{},
				&models.MesheryFilter{},
				&models.PatternResource{},
				&models.MesheryApplication{},
				&models.UserPreference{},
				&models.PerformanceTestConfig{},
				&models.SmiResultWithID{},
				&models.K8sContext{},
			)
			if err != nil {
				r.Log.Error(err)
				return "", err
			}

			krh, err := models.NewKeysRegistrationHelper(dbHandler, r.Log)
			if err != nil {
				return "", err
			}

			rm, err := registry.NewRegistryManager(dbHandler)
			if err != nil {
				return "", err
			}

			go func() {
				models.SeedComponents(r.Log, r.Config, rm)
				krh.SeedKeys(viper.GetString("KEYS_PATH"))
			}()
			r.Log.Info("Hard reset complete.")
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
				return "", model.ErrEmptyHandler
			}

			err := provider.GetGenericPersister().Where("id IN (?)", provider.GetGenericPersister().Table("kubernetes_resources").Select("id").Where("cluster_id=?", sid)).Delete(&meshsyncmodel.KubernetesKeyValue{}).Error
			if err != nil {
				return "", model.ErrEmptyHandler
			}

			err = provider.GetGenericPersister().Where("id IN (?)", provider.GetGenericPersister().Table("kubernetes_resources").Select("id").Where("cluster_id=?", sid)).Delete(&meshsyncmodel.KubernetesResourceSpec{}).Error
			if err != nil {
				return "", model.ErrEmptyHandler
			}

			err = provider.GetGenericPersister().Where("id IN (?)", provider.GetGenericPersister().Table("kubernetes_resources").Select("id").Where("cluster_id=?", sid)).Delete(&meshsyncmodel.KubernetesResourceStatus{}).Error
			if err != nil {
				return "", model.ErrEmptyHandler
			}

			err = provider.GetGenericPersister().Where("id IN (?)", provider.GetGenericPersister().Table("kubernetes_resources").Select("id").Where("cluster_id=?", sid)).Delete(&meshsyncmodel.KubernetesResourceObjectMeta{}).Error
			if err != nil {
				return "", model.ErrEmptyHandler
			}

			err = provider.GetGenericPersister().Where("cluster_id = ?", sid).Delete(&meshsyncmodel.KubernetesResource{}).Error
			if err != nil {
				return "", model.ErrEmptyHandler
			}
		}
	}

	if actions.ReSync == "true" {
		handler, ok := ctx.Value(models.HandlerKey).(*handlers.Handler)
		if !ok {
			return "", ErrResyncCluster(
				fmt.Errorf("not able to take handler from context"),
			)
		}
		instanceTracker := handler.ConnectionToStateMachineInstanceTracker
		if instanceTracker == nil {
			return "", ErrResyncCluster(
				fmt.Errorf("instance tracker is nil in handler instance"),
			)
		}
		machine, ok := instanceTracker.Get(uuid.FromStringOrNil(k8scontextID))
		if !ok || machine == nil {
			return "", ErrResyncCluster(
				fmt.Errorf(
					"instance tracker does not contain machine for connection %s",
					k8scontextID,
				),
			)
		}
		// TODO what to put as event builder?
		mashineCtx, err := kubernetes.GetMachineCtx(machine.Context, events.NewEvent())
		if err != nil {
			return "", ErrResyncCluster(
				fmt.Errorf(
					"can not receive machine context for connection %s",
					k8scontextID,
				),
			)
		}
		mesheryCtrlsHelper := mashineCtx.MesheryCtrlsHelper
		if mesheryCtrlsHelper == nil {
			return "", ErrResyncCluster(
				fmt.Errorf(
					"machine context does not contain reference to MesheryCtrlsHelper for connection %s",
					k8scontextID,
				),
			)
		}

		if err := mesheryCtrlsHelper.ResyncMeshsync(ctx); err != nil {
			return "", ErrResyncCluster(
				errors.Join(
					fmt.Errorf("error calling ResyncMeshsync for connection %s", k8scontextID),
					err,
				),
			)
		}
	}
	return model.StatusProcessing, nil
}
