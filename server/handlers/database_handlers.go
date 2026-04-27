package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/utils"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"
)

const (
	maxDatabaseArchives      = 5
	maxDatabaseArchiveAge    = 30 * 24 * time.Hour
	databaseArchiveDirectory = ".archive"
)

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
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)

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
	order = models.SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
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
	if _, err := fmt.Fprint(w, string(val)); err != nil {
		h.log.Error(err)
	}
}

// Reset the system database to its initial state.
func (h *Handler) ResetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {

	mesherydbPath := path.Join(utils.GetHome(), ".meshery/config")
	err := os.Mkdir(path.Join(mesherydbPath, ".archive"), os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		http.Error(w, "Directory could not be created due to a non-existent path.", http.StatusInternalServerError)
		return
	}
	src := path.Join(mesherydbPath, "mesherydb.sql")
	currentTime := time.Now().Format("20060102150407")
	newFileName := ".archive/mesherydb" + currentTime + ".sql"
	dst := path.Join(mesherydbPath, newFileName)

	fin, err := os.Open(src)
	if err != nil {
		http.Error(w, "The database does not exist or you don't have enough permission to access it", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := fin.Close(); err != nil {
			h.log.Error(err)
		}
	}()

	fout, err := os.Create(dst)
	if err != nil {
		http.Error(w, "Destination file can not be created", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := fout.Close(); err != nil {
			h.log.Error(err)
		}
	}()

	_, err = io.Copy(fout, fin)
	if err != nil {
		http.Error(w, "Can not copy file from source to destination", http.StatusInternalServerError)
		return
	}

	dbHandler := provider.GetGenericPersister()
	if dbHandler == nil {
		http.Error(w, "Failed to obtain database handler", http.StatusInternalServerError)
		return
	} else {
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

		go h.pruneDatabaseBackups(path.Join(mesherydbPath, databaseArchiveDirectory), maxDatabaseArchives, maxDatabaseArchiveAge)
	}
}

type archiveFile struct {
	name    string
	fullPath string
	modTime time.Time
}

func (h *Handler) pruneDatabaseBackups(archiveDir string, maxBackups int, maxAge time.Duration) {
	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		h.log.Warnf("failed to read database archive directory %s: %v", archiveDir, err)
		return
	}

	cutoff := time.Now().Add(-maxAge)
	archives := make([]archiveFile, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, "mesherydb") || filepath.Ext(name) != ".sql" {
			continue
		}

		info, statErr := entry.Info()
		if statErr != nil {
			h.log.Warnf("failed to stat archive file %s: %v", path.Join(archiveDir, name), statErr)
			continue
		}

		filePath := path.Join(archiveDir, name)
		if info.ModTime().Before(cutoff) {
			if rmErr := os.Remove(filePath); rmErr != nil {
				h.log.Warnf("failed to remove old database archive %s: %v", filePath, rmErr)
				continue
			}
			h.log.Infof("removed old database archive file: %s", filePath)
			continue
		}

		archives = append(archives, archiveFile{name: name, fullPath: filePath, modTime: info.ModTime()})
	}

	if len(archives) <= maxBackups {
		return
	}

	sort.Slice(archives, func(i, j int) bool {
		return archives[i].modTime.After(archives[j].modTime)
	})

	for _, file := range archives[maxBackups:] {
		if err := os.Remove(file.fullPath); err != nil {
			h.log.Warnf("failed to prune database archive %s: %v", file.fullPath, err)
			continue
		}
		h.log.Infof("pruned database archive file: %s", file.fullPath)
	}
}
