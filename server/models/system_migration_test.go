package models

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/v1beta3/environment"
)

// indirectType returns the struct type behind a model entry, whether it was
// registered as a value or a pointer.
func indirectType(m interface{}) reflect.Type {
	t := reflect.TypeOf(m)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// TestSystemDatabaseModelsCoversResetJoinTables pins the invariant that broke
// the database-reset path: the shared migration list MUST include the tables
// ConnectionPersister.GetConnections LEFT JOINs. A reset that re-migrated a
// subset missing environment_connection_mappings left GET
// /api/integrations/connections returning HTTP 500 (meshery-server-1032) until a
// restart. Since all three migration sites now source this one list, asserting
// membership here guards every site at once.
func TestSystemDatabaseModelsCoversResetJoinTables(t *testing.T) {
	present := map[reflect.Type]bool{}
	for _, m := range SystemDatabaseModels() {
		present[indirectType(m)] = true
	}

	required := []reflect.Type{
		reflect.TypeOf(connections.Connection{}),
		reflect.TypeOf(environment.Environment{}),
		reflect.TypeOf(environment.EnvironmentConnectionMapping{}),
	}
	for _, rt := range required {
		if !present[rt] {
			t.Errorf("SystemDatabaseModels() is missing %s; a reset re-migrating this list would drop the table and 500 GetConnections", rt)
		}
	}
}

// TestAutoMigrateSystemTablesSurvivesConnectionsQuery reproduces the reset-path
// regression at the layer the bug lived: after AutoMigrateSystemTables runs on a
// fresh database (exactly what the reset handler now does after dropping every
// table), GetConnections - which LEFT JOINs environment_connection_mappings -
// must succeed rather than error on a missing table. Before the migration list
// was unified this query returned "no such table: environment_connection_mappings".
func TestAutoMigrateSystemTablesSurvivesConnectionsQuery(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	if err := AutoMigrateSystemTables(&db); err != nil {
		t.Fatalf("AutoMigrateSystemTables: %v", err)
	}

	for _, table := range []string{"connections", "environments", "environment_connection_mappings"} {
		if !db.Migrator().HasTable(table) {
			t.Fatalf("AutoMigrateSystemTables did not create the %q table", table)
		}
	}

	if !db.Migrator().HasColumn("environments", "purpose") {
		t.Fatalf("AutoMigrateSystemTables did not create the \"purpose\" column on environments table")
	}

	// The environment_connection_mappings join only runs when the page holds at
	// least one connection, so seed a row - otherwise an empty page short-
	// circuits the join and the missing-table regression would slip through.
	conn := connections.Connection{ID: uuid.Must(uuid.NewV4()), Name: "regression-fixture", Kind: "test"}
	if err := db.Table("connections").Create(&conn).Error; err != nil {
		t.Fatalf("seed a connection row: %v", err)
	}

	cp := &ConnectionPersister{DB: &db}
	page, err := cp.GetConnections("", "", 0, 10, "", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("GetConnections after reset migration errored (the 500 regression): %v", err)
	}
	if page.TotalCount != 1 {
		t.Fatalf("expected the seeded connection on the page, got TotalCount=%d", page.TotalCount)
	}
}

// TestAutoMigrateSystemTablesIndexesEventQueries pins the indexes behind the
// Notification Center queries. events.Event comes from meshkit without these
// index tags, so AutoMigrateSystemTables creates them itself; dropping that
// step would leave GetEvents and getCountBySeverity without them. Migrating
// twice mirrors a boot followed by a reset, which must not fail on indexes
// that already exist.
func TestAutoMigrateSystemTablesIndexesEventQueries(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	for i := 0; i < 2; i++ {
		if err := AutoMigrateSystemTables(&db); err != nil {
			t.Fatalf("AutoMigrateSystemTables run %d: %v", i+1, err)
		}
	}

	for _, index := range []string{"idx_events_owner_created_at", "idx_events_owner_status_severity"} {
		if !db.Migrator().HasIndex("events", index) {
			t.Errorf("AutoMigrateSystemTables did not create the %q index", index)
		}
	}

	var plan []struct{ Detail string }
	err = db.Raw("EXPLAIN QUERY PLAN SELECT severity, count(severity) as count FROM events WHERE status = ? AND (owner = ? OR owner = ?) GROUP BY severity",
		"unread", "owner-a", "owner-b").Scan(&plan).Error
	if err != nil {
		t.Fatalf("explain the severity count query: %v", err)
	}
	for _, step := range plan {
		if strings.Contains(step.Detail, "idx_events_owner_status_severity") {
			return
		}
	}
	t.Fatalf("the severity count query does not use idx_events_owner_status_severity; plan: %+v", plan)
}

// TestAutoMigrateSystemTablesMigratesPurposeOnExistingTable reproduces the migration
// scenario where an existing database created before the purpose column was introduced
// (e.g. v1beta1 Environment schema on stable-v1.0.76) undergoes migration. AutoMigrateSystemTables
// must add the purpose column, and pre-existing rows with empty purpose must be read
// without error and evaluate to non-administrative / user.
func TestAutoMigrateSystemTablesMigratesPurposeOnExistingTable(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	// Create legacy environments table without the purpose column (matching stable-v1.0.76 schema)
	legacyDDL := "CREATE TABLE `environments` (`id` text,`name` text,`description` text,`organization_id` text,`owner` text,`created_at` datetime,`metadata` text,`updated_at` datetime,`deleted_at` datetime,PRIMARY KEY (`id`));"
	if err := db.Exec(legacyDDL).Error; err != nil {
		t.Fatalf("create legacy environments table: %v", err)
	}

	// Seed pre-existing environment with no purpose column
	envID := uuid.Must(uuid.NewV4()).String()
	insertSQL := "INSERT INTO `environments` (`id`, `name`, `description`, `organization_id`) VALUES (?, ?, ?, ?);"
	if err := db.Exec(insertSQL, envID, "legacy-env", "legacy description", uuid.Must(uuid.NewV4()).String()).Error; err != nil {
		t.Fatalf("seed legacy environment: %v", err)
	}

	if db.Migrator().HasColumn("environments", "purpose") {
		t.Fatal("expected purpose column to not exist prior to migration")
	}

	// Run migration
	if err := AutoMigrateSystemTables(&db); err != nil {
		t.Fatalf("AutoMigrateSystemTables on legacy database failed: %v", err)
	}

	if !db.Migrator().HasColumn("environments", "purpose") {
		t.Fatal("AutoMigrateSystemTables did not add purpose column to existing environments table")
	}

	// Read back the pre-existing row
	var migratedEnv environment.Environment
	if err := db.Table("environments").Where("id = ?", envID).First(&migratedEnv).Error; err != nil {
		t.Fatalf("read migrated environment: %v", err)
	}

	if migratedEnv.Purpose != "" && migratedEnv.Purpose != environment.EnvironmentPurposeUser {
		t.Fatalf("expected legacy row to read empty or default user Purpose, got %q", migratedEnv.Purpose)
	}
	if migratedEnv.IsAdministrative() {
		t.Fatal("legacy row must not evaluate as administrative")
	}
	if got := migratedEnv.EffectivePurpose(); got != environment.EnvironmentPurposeUser {
		t.Fatalf("expected EffectivePurpose() to be %q, got %q", environment.EnvironmentPurposeUser, got)
	}
}
