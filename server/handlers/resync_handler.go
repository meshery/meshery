// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	mhelpers "github.com/meshery/meshery/server/machines/helpers"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/utils"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/spf13/viper"
)

// resyncRequest is the body of
// POST /api/system/kubernetes/contexts/{contextID}/resync.
// The legacy GraphQL resolver accepted these as string enums ("true"/"false")
// — the REST surface accepts proper booleans, which is the migration's
// canonical wire shape.
type resyncRequest struct {
	ReSync    bool `json:"reSync"`
	ClearDb   bool `json:"clearDb"`
	HardReset bool `json:"hardReset"`
}

type resyncResponse struct {
	Status string `json:"status"`
}

// ResyncClusterHandler triggers a meshsync resync (and optionally drops the
// database) for the supplied Kubernetes context. It replaces the
// `resyncCluster` GraphQL query, which was mis-typed as a Query but mutates
// the entire database — this endpoint is POST as it ought to be.
func (h *Handler) ResyncClusterHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	contextID := mux.Vars(req)["contextID"]
	if contextID == "" {
		writeJSONError(w, "missing contextID", http.StatusBadRequest)
		return
	}

	defer func() {
		_ = req.Body.Close()
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	var in resyncRequest
	if err := json.Unmarshal(body, &in); err != nil {
		h.log.Error(models.ErrUnmarshal(err, "resync request"))
		writeMeshkitError(w, models.ErrUnmarshal(err, "resync request"), http.StatusBadRequest)
		return
	}

	if in.ClearDb {
		if in.HardReset {
			if err := h.runHardReset(provider); err != nil {
				h.log.Error(err)
				writeJSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if err := h.runClearDbForContext(req, contextID, provider); err != nil {
				h.log.Error(err)
				writeJSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if in.ReSync {
		if in.HardReset {
			h.log.Info("Skipping resync after hard reset due to missing Kubernetes context")
			writeJSONMessage(w, resyncResponse{Status: "PROCESSING"}, http.StatusAccepted)
			return
		}
		if err := h.runResyncForContext(req, contextID); err != nil {
			h.log.Error(err)
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writeJSONMessage(w, resyncResponse{Status: "PROCESSING"}, http.StatusAccepted)
}

func (h *Handler) runHardReset(provider models.Provider) error {
	mesherydbPath := path.Join(utils.GetHome(), ".meshery/config")
	// 0750 keeps the archive directory readable only by the meshery user
	// and its primary group, rather than world-readable like the previous
	// os.ModePerm (0777) — there is no need for other users on the host
	// to traverse the meshery state directory.
	if err := os.Mkdir(path.Join(mesherydbPath, ".archive"), 0o750); err != nil && !os.IsExist(err) {
		return err
	}

	src := path.Join(mesherydbPath, "mesherydb.sql")
	dst := path.Join(mesherydbPath, ".archive/mesherydb.sql")

	if err := copyFile(src, dst); err != nil {
		return err
	}

	dbHandler := provider.GetGenericPersister()
	if dbHandler == nil {
		return fmt.Errorf("database handler is nil")
	}

	dbHandler.Lock()
	defer dbHandler.Unlock()

	h.log.Info("Dropping Meshery Database")
	tables, err := dbHandler.Migrator().GetTables()
	if err != nil {
		return err
	}
	for _, table := range tables {
		if table == "events" {
			continue
		}
		if err := dbHandler.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	h.log.Info("Migrating Meshery Database")
	if err := dbHandler.AutoMigrate(
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
	); err != nil {
		return err
	}

	krh, err := models.NewKeysRegistrationHelper(dbHandler, h.log)
	if err != nil {
		return err
	}
	rm, err := registry.NewRegistryManager(dbHandler)
	if err != nil {
		return err
	}
	go func() {
		models.SeedComponents(h.log, h.config, rm)
		krh.SeedKeys(viper.GetString("KEYS_PATH"))
	}()
	h.log.Info("Hard reset complete.")
	return nil
}

func (h *Handler) runClearDbForContext(req *http.Request, contextID string, provider models.Provider) error {
	k8sctxs, ok := req.Context().Value(models.AllKubeClusterKey).([]*models.K8sContext)
	if !ok || len(k8sctxs) == 0 {
		return fmt.Errorf("no kubernetes contexts available")
	}
	var sid string
	for _, k8ctx := range k8sctxs {
		if k8ctx == nil {
			continue
		}
		if k8ctx.ID == contextID && k8ctx.KubernetesServerID != nil {
			sid = k8ctx.KubernetesServerID.String()
			break
		}
	}
	if sid == "" {
		return fmt.Errorf("no matching cluster for context %s", contextID)
	}
	if provider.GetGenericPersister() == nil {
		return fmt.Errorf("database handler is nil")
	}
	db := provider.GetGenericPersister()

	// Mirrors the legacy resolver deletes — id-IN-subquery against the same
	// cluster id, repeated for each meshsync side table.
	subquery := db.Table("kubernetes_resources").Select("id").Where("cluster_id=?", sid)
	if err := db.Where("id IN (?)", subquery).Delete(&meshsyncmodel.KubernetesKeyValue{}).Error; err != nil {
		return err
	}
	subquery = db.Table("kubernetes_resources").Select("id").Where("cluster_id=?", sid)
	if err := db.Where("id IN (?)", subquery).Delete(&meshsyncmodel.KubernetesResourceSpec{}).Error; err != nil {
		return err
	}
	subquery = db.Table("kubernetes_resources").Select("id").Where("cluster_id=?", sid)
	if err := db.Where("id IN (?)", subquery).Delete(&meshsyncmodel.KubernetesResourceStatus{}).Error; err != nil {
		return err
	}
	subquery = db.Table("kubernetes_resources").Select("id").Where("cluster_id=?", sid)
	if err := db.Where("id IN (?)", subquery).Delete(&meshsyncmodel.KubernetesResourceObjectMeta{}).Error; err != nil {
		return err
	}
	if err := db.Where("cluster_id = ?", sid).Delete(&meshsyncmodel.KubernetesResource{}).Error; err != nil {
		return err
	}
	return nil
}

func (h *Handler) runResyncForContext(req *http.Request, contextID string) error {
	instanceTracker := h.ConnectionToStateMachineInstanceTracker
	if instanceTracker == nil {
		return fmt.Errorf("instance tracker is nil")
	}

	k8sCtxs, ok := req.Context().Value(models.AllKubeClusterKey).([]*models.K8sContext)
	if !ok || len(k8sCtxs) == 0 {
		return fmt.Errorf("no kubernetes contexts available")
	}

	var k8sCtx *models.K8sContext
	for _, v := range k8sCtxs {
		if v != nil && v.ID == contextID {
			k8sCtx = v
			break
		}
	}
	if k8sCtx == nil {
		return fmt.Errorf("no k8s context for contextID %s", contextID)
	}
	if k8sCtx.ConnectionID == "" {
		return fmt.Errorf("connectionID is empty for k8s context %s", contextID)
	}

	machine, ok := instanceTracker.Get(uuid.FromStringOrNil(k8sCtx.ConnectionID))
	if !ok || machine == nil {
		return fmt.Errorf("no state machine for connection %s", k8sCtx.ConnectionID)
	}

	if err := mhelpers.ResyncResources(req.Context(), machine); err != nil {
		return fmt.Errorf("resync resources for context %s: %w", contextID, err)
	}
	return nil
}

func copyFile(src, dst string) error {
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	fout, err := os.Create(dst)
	if err != nil {
		if cerr := fin.Close(); cerr != nil {
			return errors.Join(err, cerr)
		}
		return err
	}
	if _, err := io.Copy(fout, fin); err != nil {
		if cerr := errors.Join(fin.Close(), fout.Close()); cerr != nil {
			return errors.Join(err, cerr)
		}
		return err
	}
	if err := fin.Close(); err != nil {
		if cerr := fout.Close(); cerr != nil {
			return errors.Join(err, cerr)
		}
		return err
	}
	return fout.Close()
}
