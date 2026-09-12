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
	system "github.com/meshery/schemas/models/v1beta1/system"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"
)

func (h *Handler) GetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	var tables []system.SystemDatabaseTable
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

	for i := range tables {
		h.dbHandler.DB.Table(tables[i].Name).Count(&tables[i].Count)
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

// ResetSystemDatabase resets the system database to its initial state.
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
		// Held for the whole workflow, including the seeding goroutine below.
		// dbHandler's own lock is released when this function returns, which
		// leaves the goroutine unprotected: a second reset could acquire it and
		// drop tables mid-seed. Released by whichever goroutine finishes the
		// work, not on handler return - hence the seedingStarted flag, which
		// makes sure every early return below still releases.
		if !models.TryAcquireResetLock() {
			writeMeshkitError(w, ErrResetInProgress(), http.StatusConflict)
			return
		}
		seedingStarted := false
		defer func() {
			if !seedingStarted {
				models.ReleaseResetLock()
			}
		}()

		dbHandler.Lock()
		defer dbHandler.Unlock()

		tables, err := dbHandler.Migrator().GetTables()
		if err != nil {
			writeMeshkitError(w, ErrAccessDatabaseTables(err), http.StatusInternalServerError)
			return
		}

		for _, table := range tables {
			// The GraphQL hard reset (resolver/meshsync.go) skips this table
			// deliberately; this path did not, so the same reset produced
			// different results depending on the entry point. Nothing re-seeds
			// events, so dropping it here is unrecoverable data loss.
			// Re-migrating it below is idempotent.
			if table == "events" {
				continue
			}
			err = dbHandler.Migrator().DropTable(table)
			if err != nil {
				writeMeshkitError(w, ErrDropDatabaseTable(err), http.StatusInternalServerError)
				return
			}
		}

		// Re-migrate the full system-table set after dropping every table.
		// Sharing models.SystemDatabaseModels with boot and the GraphQL hard
		// reset is what keeps the reset from re-creating only a stale subset -
		// previously environments/environment_connection_mappings were never
		// recreated, so GetConnections (which LEFT JOINs
		// environment_connection_mappings) returned 500 until a restart.
		err = models.AutoMigrateSystemTables(dbHandler)

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

		// Seeded synchronously, before the success response: organizations is
		// dropped and re-migrated above but was only ever seeded at boot, and
		// the UI skips its keys query without an org - resolving to no
		// permissions until a restart.
		if lp, ok := provider.(*models.DefaultLocalProvider); ok {
			if err := lp.SeedDefaultOrganization(); err != nil {
				writeMeshkitError(w, err, http.StatusInternalServerError)
				return
			}
		}

		seedingStarted = true
		go func() {
			defer models.ReleaseResetLock()
			models.RunSeedStage(h.log, "user keys", func() {
				krh.SeedKeys(viper.GetString("KEYS_PATH"))
			})
			models.RunSeedStage(h.log, "content", func() {
				if lp, ok := provider.(*models.DefaultLocalProvider); ok {
					lp.SeedContent(h.log)
				}
			})
			models.RunSeedStage(h.log, "models", func() {
				models.SeedComponents(h.log, h.config, h.registryManager, dbHandler)
			})
		}()
		writeJSONMessage(w, system.SystemMessageResponse{Message: "Database reset successful"}, http.StatusOK)
	}
}
