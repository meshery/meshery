package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/utils"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	systemv1beta1 "github.com/meshery/schemas/models/v1beta1/system"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"
)

func (h *Handler) GetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	var tables []*models.SqliteSchema
	var recordCount int
	var totalTables int64
	page, offset, limit, search, sortColumn, sortDesc, err := getSystemDatabaseQueryParams(r)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	tableFinder := h.dbHandler.DB.Table("sqlite_schema").
		Where("type = ?", "table")

	if search != "" {
		tableFinder = tableFinder.Where("name LIKE ?", "%"+search+"%")
	}

	tableFinder.Count(&totalTables)

	if limit != 0 {
		tableFinder = tableFinder.Limit(limit)
	}
	if offset != 0 {
		tableFinder = tableFinder.Offset(offset)
	}
	if sortColumn != "" {
		tableFinder = tableFinder.Order(clause.OrderByColumn{Column: clause.Column{Name: sortColumn}, Desc: sortDesc})
	}

	tableFinder.Find(&tables)

	responseTables := make([]systemv1beta1.SystemDatabaseTable, 0, len(tables))
	for _, table := range tables {
		h.dbHandler.DB.Table(table.Name).Count(&table.Count)
		recordCount += int(table.Count)
		name := table.Name
		tableType := table.Type
		responseTables = append(responseTables, systemv1beta1.SystemDatabaseTable{
			Count: table.Count,
			Name:  &name,
			Type:  &tableType,
		})
	}

	databaseSummary := &systemv1beta1.SystemDatabaseSummary{
		Page:        page,
		PageSize:    limit,
		TotalTables: int(totalTables),
		RecordCount: recordCount,
		Tables:      responseTables,
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

func getSystemDatabaseQueryParams(req *http.Request) (page, offset, limit int, search, sortColumn string, sortDesc bool, err error) {
	urlValues := req.URL.Query()

	page, _ = strconv.Atoi(urlValues.Get("page"))
	if page < 0 {
		page = 0
	}

	pageSize := urlValues.Get("pageSize")
	if pageSize == "" {
		pageSize = urlValues.Get("pagesize")
	}
	if pageSize == "" {
		limit = defaultPageSize
	} else {
		limit, err = strconv.Atoi(pageSize)
		if err != nil || limit < 1 {
			return 0, 0, 0, "", "", false, fmt.Errorf("pageSize must be an integer greater than 0")
		}
	}

	sortColumns := map[string]string{
		"name": "name",
	}
	sort := urlValues.Get("sort")
	if sort != "" {
		var ok bool
		sortColumn, ok = sortColumns[sort]
		if !ok {
			return 0, 0, 0, "", "", false, fmt.Errorf("sort must be name")
		}
	}

	order := urlValues.Get("order")
	switch order {
	case "", "asc":
	case "desc":
		sortDesc = true
	default:
		return 0, 0, 0, "", "", false, fmt.Errorf("order must be asc or desc")
	}

	search = urlValues.Get("search")
	offset = page * limit
	return
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
		if h.ConnectionToStateMachineInstanceTracker != nil {
			h.ConnectionToStateMachineInstanceTracker.Clear()
		}

		krh, err := models.NewKeysRegistrationHelper(dbHandler, h.log)
		if err != nil {
			writeMeshkitError(w, ErrMigrateDatabaseTables(err), http.StatusInternalServerError)
			return
		}
		go func() {
			models.SeedComponents(h.log, h.config, h.registryManager)
			krh.SeedKeys(viper.GetString("KEYS_PATH"))
		}()
		writeJSONMessage(w, systemv1beta1.SystemMessageResponse{Message: "Database reset"}, http.StatusOK)
	}
}
