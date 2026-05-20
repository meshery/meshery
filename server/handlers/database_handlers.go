package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/utils"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/spf13/viper"
)

func (h *Handler) GetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	var tables []*models.TableSchema
	var recordCount int
	page, offset, limit, search, _, _, _ := getPaginationParams(r)

	tableNames, err := h.dbHandler.DB.Migrator().GetTables()
	if err != nil {
		http.Error(w, "Failed to fetch tables", http.StatusInternalServerError)
		return
	}

	for _, name := range tableNames {
		if search != "" && !strings.Contains(name, search) {
			continue
		}

		var count int64
		h.dbHandler.DB.Table(name).Count(&count)

		tables = append(tables, &models.TableSchema{
			Name:  name,
			Count: count,
		})
		recordCount += int(count)
	}

	totalTables := len(tables)

	if offset > totalTables {
		offset = totalTables
	}
	end := offset + limit
	if limit == 0 || end > totalTables {
		end = totalTables
	}

	paginatedTables := tables[offset:end]

	databaseSummary := &models.DatabaseSummary{
		Page:        page,
		PageSize:    limit,
		TotalTables: totalTables,
		RecordCount: recordCount,
		Tables:      paginatedTables,
	}

	w.Header().Set("Content-Type", "application/json")

	val, err := json.Marshal(databaseSummary)
	if err != nil {
		h.log.Error(err)
	}
	_, err = fmt.Fprint(w, string(val))
	if err != nil {
		return
	}
}

// Reset the system database to its initial state.
func (h *Handler) ResetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	isSQLite := viper.GetString("DATABASE_URL") == ""

	if isSQLite {
		mesherydbPath := path.Join(utils.GetHome(), ".meshery/config")
		err := os.MkdirAll(path.Join(mesherydbPath, ".archive"), os.ModePerm) // Заменил Mkdir на MkdirAll для надежности
		if err != nil {
			http.Error(w, "Directory could not be created.", http.StatusInternalServerError)
			return
		}

		src := path.Join(mesherydbPath, "mesherydb.sql")
		currentTime := time.Now().Format("20060102150407")
		dst := path.Join(mesherydbPath, ".archive/mesherydb"+currentTime+".sql")

		fin, err := os.Open(src)
		if err != nil {
			http.Error(w, "Database file does not exist", http.StatusInternalServerError)
			return
		}

		fout, err := os.Create(dst)
		if err != nil {
			fin.Close()
			http.Error(w, "Destination file can not be created", http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(fout, fin)
		fin.Close()
		fout.Close()
		if err != nil {
			http.Error(w, "Can not copy file", http.StatusInternalServerError)
			return
		}
	}

	dbHandler := provider.GetGenericPersister()
	if dbHandler == nil {
		writeMeshkitError(w, ErrObtainDatabaseHandler(), http.StatusInternalServerError)
		return
	}

	dbHandler.Lock()
	defer dbHandler.Unlock()

	tables, err := dbHandler.Migrator().GetTables()
	if err != nil {
		http.Error(w, "Can not access database tables", http.StatusInternalServerError)
		return
	}

	for _, table := range tables {
		err = dbHandler.Migrator().DropTable(table)
		if err != nil {
			http.Error(w, "Cannot drop table from database", http.StatusInternalServerError)
			return
		}
	}

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
		http.Error(w, "Can not migrate tables to database", http.StatusInternalServerError)
		return
	}

	rm, err := registry.NewRegistryManager(dbHandler)
	if err != nil {
		http.Error(w, "Can not migrate tables to database", http.StatusInternalServerError)
		return
	}
	h.registryManager = rm

	krh, err := models.NewKeysRegistrationHelper(dbHandler, h.log)
	if err != nil {
		http.Error(w, "Can not migrate tables to database", http.StatusInternalServerError)
		return
	}
	go func() {
		models.SeedComponents(h.log, h.config, h.registryManager)
		krh.SeedKeys(viper.GetString("KEYS_PATH"))
	}()
	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(w, "Database reset successful"); err != nil {
		h.log.Error(err)
	}
}
