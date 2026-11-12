---
name: migration-risk-analyzer
description: Risk assessment and migration safety analysis for GORM AutoMigrate schema changes on SQLite.
---

# Migration Risk Analyzer

`migration-risk-analyzer` evaluates schema changes and GORM struct modifications in Meshery Server. Meshery reconciles its local SQLite database with GORM `AutoMigrate` on every boot. Because SQLite handles schema alterations differently from client-server relational databases, struct changes can silently orphan data, lock databases, or break server restarts.

## Key SQLite & GORM AutoMigrate Hazards

### 1. Column Renames from Schemas Bumps
* **Hazard**: GORM `AutoMigrate` on SQLite cannot detect that a column was renamed. When a struct field or its `gorm:"column:..."` tag changes, GORM simply adds a *new* column.
* **Impact**: The old column remains in SQLite, but application queries read the new (empty) column. User data is orphaned silently without migration errors.
* **Verification**: Check if the change renames any field or column. If so, a data backfill or dual-read fallback must be implemented in the model persister (see reference test: `TestGetPerformanceProfilesReadsMigratedColumns` in `server/models/performance_profile_persister_test.go`).

### 2. `NOT NULL` Constraints on Existing Tables
* **Hazard**: Adding `gorm:"not null"` to a struct field on an existing table will fail during SQLite migration if existing rows contain `NULL` or if SQLite attempts table recreation without a default value.
* **Remediation**: Always specify a sensible `gorm:"default:..."` or ensure existing tables are safely backfilled prior to adding strict nullability constraints.

### 3. Type Alterations and Table Rebuilds
* **Hazard**: SQLite does not support `ALTER COLUMN TYPE`. GORM emulates this by creating a temporary table (`_new_<table>`), copying data, dropping the old table, and renaming the temporary table.
* **Impact**: Can lead to deadlocks, dropped indexes, or foreign key constraint violations if foreign keys are enabled.

### 4. Registration in `SystemDatabaseModels()`
* **Hazard**: Meshery maintains an authoritative central list of system tables in `server/models/system_migration.go:SystemDatabaseModels()`.
* **Impact**: If a new model or join table is introduced but omitted from `SystemDatabaseModels()`, database reset routines (`DELETE /api/system/database/reset` or GraphQL `resyncCluster`) will drop or skip the table, leading to 500 errors on subsequent operations.
* **Verification**: Run `go test -run TestSystemDatabaseModelsCoversResetJoinTables ./server/models/...` to ensure full coverage.

### 5. SQLite WAL Mode Hazards
* **Hazard**: Meshery runs SQLite in WAL mode (`_journal_mode=WAL`). Reset backups that perform a naive file copy of `mesherydb.sql` miss active transactions held in `mesherydb.sql-wal`.
* **Remediation**: Database resets and backups must snapshot data safely (e.g. via `VACUUM INTO` or atomic flush) to avoid corrupted or partial recovery states.

---

## Schema Review Checklist

Before approving any PR touching `server/models/*`:

- [ ] **Field Renames**: Does the PR change any Go struct field or GORM column tag? If so, verify that existing data in SQLite is copied or read with backwards compatibility.
- [ ] **Table List Coverage**: Is any newly introduced table or join model included in `SystemDatabaseModels()`?
- [ ] **Defaults for New Columns**: Do newly added non-nullable columns have explicit `gorm:"default:..."` tags?
- [ ] **Reset & Rollback Safety**: Is the model covered in `SystemDatabaseModels()` so database resets and hard resets do not drop or corrupt tables?
- [ ] **WAL Mode Handling**: Are SQLite resets and backups capturing the database and active WAL journal safely (e.g. via `VACUUM INTO` or atomic copy)?
- [ ] **Regression Tests**: Does the PR include tests in `server/models/*_test.go` confirming fresh-install migrations, existing database upgrades, and reset behaviors succeed?
