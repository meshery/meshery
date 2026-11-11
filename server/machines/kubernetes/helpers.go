package kubernetes

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/meshery/schemas/models/core"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/helpers"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/utils/kubernetes"
	"github.com/spf13/viper"
)

// adapterTracker is initialized lazily on first access so that the ADAPTER_URLS
// env-var read happens AFTER main() has called viper.AutomaticEnv(). A package-
// level eager init would run at import time, before viper sees env vars, and
// silently produce an empty tracker for env-only configurations.
var (
	adapterTrackerOnce sync.Once
	adapterTrackerInst *helpers.AdaptersTracker
)

func getAdapterTracker() *helpers.AdaptersTracker {
	adapterTrackerOnce.Do(func() {
		adapterURLs := utils.SplitAndTrim(viper.GetString("ADAPTER_URLS"), ", \t\n\r")
		adapterTrackerInst = helpers.NewAdaptersTracker(adapterURLs)
	})
	return adapterTrackerInst
}

func GenerateClientSetAction(k8sContext *models.K8sContext, eventBuilder *events.EventBuilder, log logger.Handler) (*kubernetes.Client, error) {
	eventBuilder.ActedUpon(uuid.FromStringOrNil(k8sContext.ConnectionID))

	eventMetadata := map[string]interface{}{}

	handler, err := models.GenerateK8sClientSet(k8sContext, eventBuilder, eventMetadata, log)
	if handler == nil {
		return nil, err
	}

	return handler, nil
}

func GetMachineCtx(machinectx interface{}, eb *events.EventBuilder) (*MachineCtx, error) {
	machineCtx, ok := machinectx.(*MachineCtx)
	if !ok {
		err := machines.ErrAssertMachineCtx(fmt.Errorf("asserting of context %v failed", machinectx))
		if eb != nil {
			eb.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			})
		}
		return nil, err
	}
	return machineCtx, nil
}

func AssignClientSetToContext(machinectx *MachineCtx, eventBuilder *events.EventBuilder) error {
	if machinectx.clientset != nil {
		return nil
	}

	k8sContext := machinectx.K8sContext
	eventBuilder.ActedUpon(uuid.FromStringOrNil(k8sContext.ConnectionID))

	handler, err := GenerateClientSetAction(&k8sContext, eventBuilder, machinectx.log)
	if err != nil {
		// perofmr event publishinh and err handling
		return err
	}
	machinectx.clientset = handler
	return nil
}

// RefreshK8sContext replaces the stored context when the newly loaded one
// carries different connection material (server or credentials) and rebuilds
// the clientset from it, so a credential update is picked up by the next
// Discovery instead of probing with the clientset built when the machine was
// created. It reports whether anything changed. When the clientset cannot be
// rebuilt the old context is kept and the error is returned, so a bad update
// never leaves the machine without a clientset.
//
// Building the clientset is local work (it assembles a kubeconfig and a
// client; the cluster is only contacted by the Discovery that follows), so
// this is safe to call while holding the machine lock. log is the caller's
// logger: the stored machine context's own log may be nil on machines built
// outside the normal init path, and a failed build logs through it.
func (m *MachineCtx) RefreshK8sContext(newCtx models.K8sContext, eventBuilder *events.EventBuilder, log logger.Handler) (bool, error) {
	if k8sContextCredentialsEqual(m.K8sContext, newCtx) {
		return false, nil
	}
	handler, err := GenerateClientSetAction(&newCtx, eventBuilder, log)
	if err != nil {
		return false, err
	}
	// The replacement runs under the setup lock: ConnectAction's setup
	// goroutine reads K8sContext under that lock (and snapshots it per
	// iteration), so a refresh racing a setup cannot hand the goroutine a
	// half-replaced context.
	m.setupLock.Lock()
	m.K8sContext = newCtx
	m.clientset = handler
	m.setupLock.Unlock()
	return true, nil
}

// k8sContextCredentialsEqual reports whether two contexts carry the same
// connection material: the API server and the auth/cluster maps. Other fields
// (name, ids, server id, version) do not affect how the clientset talks to
// the cluster, so a Discovery enrichment that sets them must not trigger a
// rebuild.
func k8sContextCredentialsEqual(a, b models.K8sContext) bool {
	return a.Server == b.Server &&
		sqlMapEqual(a.Auth, b.Auth) &&
		sqlMapEqual(a.Cluster, b.Cluster)
}

// sqlMapEqual compares two sql.Map values, treating a nil map and an empty
// map as equal: a context loaded from JSON has nil maps where one built in
// code has empty ones, and that difference must not trigger a rebuild.
func sqlMapEqual(a, b sql.Map) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		other, ok := b[k]
		if !ok || !reflect.DeepEqual(v, other) {
			return false
		}
	}
	return true
}

func AssignControllerHandlers(machinectx *MachineCtx, systemID *core.Uuid, provider models.Provider) {
	machinectx.MesheryCtrlsHelper = models.NewMesheryControllersHelper(
		machinectx.log,
		models.NewOperatorDeploymentConfig(getAdapterTracker()),
		models.GetDBInstance(),
		machinectx.EventBroadcaster,
		provider,
		systemID,
	)
}
