package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/utils"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	system "github.com/meshery/schemas/models/v1beta1/system"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"
)

func (h *Handler) GetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	var tables []system.SystemDatabaseTable
	var recordCount int
	var totalTables int64
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)

	tableQueryBase := h.dbHandler.DB.Table("sqlite_schema AS s").
		Where("type = ?", "table")

	if search != "" {
		tableQueryBase = tableQueryBase.Where("name LIKE ?", "%"+search+"%")
	}

	tableQueryBase.Count(&totalTables)

	tableFinder := tableQueryBase.
		Select("s.*, CAST(COALESCE(st.stat, '0') AS INTEGER) as count").
		Joins("LEFT JOIN sqlite_stat1 st ON s.name = st.tbl")

	if limit != 0 {
		tableFinder = tableFinder.Limit(limit)
	}
	if offset != 0 {
		tableFinder = tableFinder.Offset(offset)
	}
	order = models.SanitizeOrderInput(order, []string{"created_at", "updated_at", "name", "count"})
	if order != "" {
		if sort == "desc" {
			tableFinder = tableFinder.Order(clause.OrderByColumn{Column: clause.Column{Name: order}, Desc: true})
		} else {
			tableFinder = tableFinder.Order(order)
		}
	}

	tableFinder.Find(&tables)

	for i := range tables {
		recordCount += int(tables[i].Count)
	}

	databaseSummary := &system.SystemDatabaseSummary{
		Page:        page,
		PageSize:    limit,
		TotalTables: int(totalTables),
		RecordCount: recordCount,
		Tables:      tables,
	}

	w.Header().Set("Content-Type", "application/json")

	val, err := json.Marshal(databaseSummary)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := fmt.Fprint(w, string(val)); err != nil {
		h.log.Error(err)
	}
}

// Reset the system database to its initial state.
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
		writeJSONMessage(w, system.SystemMessageResponse{Message: "Database reset successful"}, http.StatusOK)
	}
}
