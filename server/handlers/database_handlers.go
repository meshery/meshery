package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/utils"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/spf13/viper"
)

func (h *Handler) GetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	var recordCount int
	page, offset, limit, search, order, sortOrder, _ := getPaginationParams(r)

	dbHandler := provider.GetGenericPersister()
	if dbHandler == nil {
		http.Error(w, "Failed to obtain database handler", http.StatusInternalServerError)
		return
	}

	allTables, err := dbHandler.Migrator().GetTables()
	if err != nil {
		http.Error(w, "Cannot access database tables", http.StatusInternalServerError)
		return
	}

	if search != "" {
		var filtered []string
		for _, t := range allTables {
			if strings.Contains(strings.ToLower(t), strings.ToLower(search)) {
				filtered = append(filtered, t)
			}
		}
		allTables = filtered
	}

	if order == "name" {
		if sortOrder == "desc" {
			slices.SortFunc(allTables, func(a, b string) int {
				return strings.Compare(b, a)
			})
		} else {
			slices.Sort(allTables)
		}
	}

	totalTables := len(allTables)

	start := offset
	end := offset + limit
	if limit == 0 || end > totalTables {
		end = totalTables
	}
	if start > totalTables {
		start = totalTables
	}
	pagedTables := allTables[start:end]

	var tables []*models.SqliteSchema
	for _, tableName := range pagedTables {
		var count int64
		
		if h.dbHandler != nil && h.dbHandler.DB != nil {
			h.dbHandler.DB.Table(tableName).Count(&count)
		}
		
		recordCount += int(count)

		tables = append(tables, &models.SqliteSchema{
			Name:  tableName,
			Type:  "table",
			Count: count,
		})
	}

	databaseSummary := &models.DatabaseSummary{
		Page:        page,
		PageSize:    limit,
		TotalTables: totalTables,
		RecordCount: recordCount,
		Tables:      tables,
	}

	w.Header().Set("Content-Type", "application/json")

	val, err := json.Marshal(databaseSummary)
	if err != nil {
		fmt.Println(err)
		h.log.Error(err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	if _, err := fmt.Fprint(w, string(val)); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) ResetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {

	mesherydbPath := path.Join(utils.GetHome(), ".meshery/config")
	err := os.Mkdir(path.Join(mesherydbPath, ".archive"), os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		writeMeshkitError(w, ErrCreateDatabaseArchiveDir(err), http.StatusInternalServerError)
		return
	}
	src := path.Join(mesherydbPath, "mesherydb.sql")
	currentTime := time.Now().Format("20060102150407")
	newFileName := ".archive/mesherydb" + currentTime + ".sql"
	dst := path.Join(mesherydbPath, newFileName)

	fin, err := os.Open(src)
	if err != nil {
		writeMeshkitError(w, ErrOpenDatabaseFile(err), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := fin.Close(); err != nil {
			h.log.Error(err)
		}
	}()

	fout, err := os.Create(dst)
	if err != nil {
		writeMeshkitError(w, ErrCreateDatabaseArchiveFile(err), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := fout.Close(); err != nil {
			h.log.Error(err)
		}
	}()

	_, err = io.Copy(fout, fin)
	if err != nil {
		writeMeshkitError(w, ErrCopyDatabaseFile(err), http.StatusInternalServerError)
		return
	}

	dbHandler := provider.GetGenericPersister()
	if dbHandler == nil {
		writeMeshkitError(w, ErrObtainDatabaseHandler(), http.StatusInternalServerError)
		return
	} else {
		dbHandler.Lock()
		defer dbHandler.Unlock()

		tables, err := dbHandler.Migrator().GetTables()
		if err != nil {
			writeMeshkitError(w, ErrAccessDatabaseTables(err), http.StatusInternalServerError)
			return
		}

		for _, table := range tables {
			err = dbHandler.Migrator().DropTable(table)
			if err != nil {
				writeMeshkitError(w, ErrDropDatabaseTable(err), http.StatusInternalServerError)
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
			writeMeshkitError(w, ErrMigrateDatabaseTables(err), http.StatusInternalServerError)
			return
		}

		rm, err := registry.NewRegistryManager(dbHandler)
		if err != nil {
			writeMeshkitError(w, ErrMigrateDatabaseTables(err), http.StatusInternalServerError)
			return
		}
		h.registryManager = rm

		krh, err := models.NewKeysRegistrationHelper(dbHandler, h.log)
		if err != nil {
			writeMeshkitError(w, ErrMigrateDatabaseTables(err), http.StatusInternalServerError)
			return
		}
		go func() {
			models.SeedComponents(h.log, h.config, h.registryManager)
			krh.SeedKeys(viper.GetString("KEYS_PATH"))
		}()
		writeJSONMessage(w, map[string]string{"message": "Database reset successful"}, http.StatusOK)
	}
}