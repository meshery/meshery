---
name: migration-risk-analyzer
description: 'Assess the risk of a Meshery database schema change - a gorm model/struct change, a schemas bump that renames fields, a new system table, or hand-written SQL - for data loss, failed boots on upgraded installs, and downgrade safety. Triggers on: "analyze this migration", "migration risk", "is this schema change safe", "will this break existing installs", "/migration-risk-analyzer".'
metadata:
  version: 1.0.0-meshery
  upstream: https://github.com/Mathews-Tom/armory/tree/main/skills/migration-risk-analyzer
  upstream-version: 1.1.1
  license: MIT (see .agents/skills/ATTRIBUTIONS.md)
---

# Migration Risk Analyzer

Meshery has no versioned migration files. The server opens a local **SQLite**
database (`server/models/sql-utils.go`, `database.SQLITE`, file
`$USER_DATA_FOLDER/mesherydb.sql`) and reconciles it on every boot with gorm
`AutoMigrate` over `SystemDatabaseModels()` (`server/models/system_migration.go`).
The registry manager migrates its own tables the same way.

So "a migration" here is any change that alters what `AutoMigrate` sees: a
struct field added/removed/renamed/retyped, a changed `gorm:` tag, a new model,
or a `meshery/schemas` bump that renames a Go field (which renames the column -
gorm derives it from the field name, not the `db:` tag).

The risk is not lock contention. It is: **an existing user upgrades, the new
binary boots against their old `mesherydb.sql`, and either the boot fails, data
silently disappears from queries, or a later downgrade breaks.**

AGENTS.md lists database migrations under "Require Human Review" - always say so
in the report.

## Workflow

### 1. Find the schema delta

- `git diff "$(git merge-base HEAD origin/master)"...HEAD -- server/ go.mod`
- For each model in `SystemDatabaseModels()` (and registry models) touched -
  directly, or through a `github.com/meshery/schemas` / `meshkit` / `meshsync`
  version bump - list field-level changes. For a dependency bump, diff the
  model structs between the two module versions.
- Translate each change into what gorm's SQLite migrator will do (table below).

### 2. Classify each operation

| Change | gorm `AutoMigrate` on SQLite | Risk |
|---|---|---|
| New model / table | `CREATE TABLE` | Low - but must be added to `SystemDatabaseModels()` or the reset paths will not recreate it |
| New nullable field | `ALTER TABLE ADD COLUMN` | Low |
| New `NOT NULL` field without default | `ADD COLUMN ... NOT NULL` | **High** - SQLite rejects it on a non-empty table; the boot migration fails |
| Field removed from struct | Nothing - column stays | Low for upgrade; stale column keeps old data |
| Field **renamed** (incl. schemas rename) | Adds a new empty column; old column and its data stay orphaned | **High** - existing rows read back as zero values; raw `Select`/`Where`/`Order`/`Joins` strings using the old name fail at runtime |
| Type changed | SQLite migrator **rebuilds the table** (create temp, copy, drop, rename) | **High** - lossy casts, constraint failures on copy, long boot on large tables (`kubernetes_resources`, `events`) |
| Index / unique added | `CREATE [UNIQUE] INDEX` | Medium - unique fails if existing rows already collide |
| Primary key / embedded struct change | Table rebuild or mismatched keys | **High** |
| Hand-written `db.Exec` DDL/DML | Whatever it says | Review as raw SQL; must be idempotent - it runs on every boot |

### 3. Check the consumers

For every renamed or retyped column:

- `grep -rn '<old_column>' server/ mesheryctl/` - raw column strings in
  `Select`, `Where`, `Order`, `Joins`, `Scan`, `Pluck`, and test fixtures
- Any `SanitizeOrderInput` allow-list naming the column
- GraphQL resolvers and REST handlers that sort/filter by it
- Whether gorm errors on those queries are propagated or swallowed (a
  swallowed error turns the rename into silently empty results)

### 4. Upgrade, downgrade and reset paths

- **Upgrade:** does the new binary boot against a DB created by the previous
  release, with realistic data (designs, connections, events, MeshSync
  resources)? Is any data that existing users rely on now unreachable?
- **Backfill:** does a rename need a one-time copy from the old column? If so,
  it must be idempotent, run before readers, and tolerate the old column not
  existing (fresh installs).
- **Downgrade:** an older binary opening the newer DB - does it hit unknown
  `NOT NULL` columns on insert, or a rebuilt table it cannot read?
- **Reset:** `handlers/database_handlers.go` and the GraphQL hard reset drop
  every table and re-migrate via `AutoMigrateSystemTables` - confirm a new
  table is in that list.
- **Remote provider:** if the change also alters a contract the remote provider
  persists (meshery-cloud runs PostgreSQL), flag it - lock behavior and
  rollback planning then apply on that side.

### 5. Verify

Prefer evidence over reasoning. Suggested check, adapted to the change:

1. Build the previous release's server, boot it once so it creates
   `mesherydb.sql`, seed representative rows.
2. Boot the branch build against that same file. Confirm boot succeeds and
   the migrated rows read back correctly through the real handlers.
3. Inspect the schema: `sqlite3 mesherydb.sql '.schema <table>'`.
4. Add a Go test that `AutoMigrate`s the old struct, inserts rows, migrates the
   new struct, and asserts the data. For a schemas rename,
   `TestGetPerformanceProfilesReadsMigratedColumns` in
   `server/models/performance_profile_persister_test.go` is the regression
   test to model on.

## Output

```
## Migration Risk Analysis

Summary: <N> schema changes across <tables>. Overall risk: High | Medium | Low.
Requires human review (AGENTS.md): yes.

| # | Model.Field | Change | gorm/SQLite effect | Risk | Data preserved on upgrade |
|---|---|---|---|---|---|

Consumers needing update: <file:line list, or "none found" with the grep used>
Upgrade path: <result, observed or reasoned - say which>
Downgrade path: <result>
Reset path: <table in SystemDatabaseModels? yes/no>
Required before landing: <backfill, test, consumer fixes>
```

## Calibration

1. **Assume real data.** User DBs hold years of events and MeshSync resources;
   a table rebuild that is instant in a test can stall boot for minutes.
2. **A rename is a data-loss bug until proven otherwise.** gorm never moves
   data between columns.
3. **Boot failure is the worst outcome.** A failing `AutoMigrateSystemTables`
   takes down the whole server for that user.
4. **Say what you observed vs. inferred.** "Booted old DB on branch build, rows
   intact" and "should be fine because the column is nullable" are different
   claims.
