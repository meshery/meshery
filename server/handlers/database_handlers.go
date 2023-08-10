package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/utils"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"gorm.io/gorm/clause"
)

const defaultPageSize = 10

// swagger:route GET /api/system/database GetSystemDatabase idGetSystemDatabase
// Handle GET request for getting summary about the system database.
//
// # Tables can be further filtered through query parameter
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 10. To return all results: ```pagesize=all```
//
// ```?search={tablename}``` If search is non empty then a greedy search is performed
// responses:
//
//	200: systemDatabaseResponseWrapper
func (h *Handler) GetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	var tables []*models.SqliteSchema
	var recordCount int
	var totalTables int64

	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit <= 0 {
			limit = defaultPageSize
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")

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
	if order != "" {
		if sort == "desc" {
			tableFinder = tableFinder.Order(clause.OrderByColumn{Column: clause.Column{Name: order}, Desc: true})
		} else {
			tableFinder = tableFinder.Order(order)
		}
	}

	tableFinder.Find(&tables)

	for _, table := range tables {
		h.dbHandler.DB.Table(table.Name).Count(&table.Count)
		recordCount += int(table.Count)
	}

	databaseSummary := &models.DatabaseSummary{
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
	fmt.Fprint(w, string(val))
}

func (h *Handler) ResetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {

	mesherydbPath := path.Join(utils.GetHome(), ".meshery/config")
	err := os.Mkdir(path.Join(mesherydbPath, ".archive"), os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	src := path.Join(mesherydbPath, "mesherydb.sql")
	dst := path.Join(mesherydbPath, ".archive/mesherydb.sql")

	fin, err := os.Open(src)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer fin.Close()

	fout, err := os.Create(dst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	dbHandler := provider.GetGenericPersister()
	if dbHandler == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		dbHandler.Lock()
		defer dbHandler.Unlock()

		tables, err := dbHandler.Migrator().GetTables()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		for _, table := range tables {
			err = dbHandler.Migrator().DropTable(table)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

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
			&models.K8sContext{},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "Database reset successful")

	}
}
