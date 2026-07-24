package models

import (
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/v1beta1/environment"
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
