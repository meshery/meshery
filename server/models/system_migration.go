package models

import (
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/meshery/schemas/models/v1beta1/environment"
	"github.com/meshery/schemas/models/v1beta1/workspace"
	schemasOrganization "github.com/meshery/schemas/models/v1beta2/organization"
)

// SystemDatabaseModels is the single, authoritative list of every system table
// Meshery boots with. It is the source of truth used by every path that
// migrates the system database: server boot (cmd/main.go), the database reset
// handler (handlers/database_handlers.go), and the GraphQL hard reset
// (internal/graphql/resolver/meshsync.go).
//
// Both reset paths drop every table and then re-migrate. Historically each site
// hand-maintained its own subset of this list, and the two reset subsets drifted
// out of date: tables added to boot over time (connections, environments,
// environment_connection_mappings, workspaces, events, ...) were never
// recreated after a reset. That is why GET /api/integrations/connections
// returned HTTP 500 (meshery-server-1032) after a reset - ConnectionPersister
// .GetConnections LEFT JOINs environment_connection_mappings, which no longer
// existed. Routing all three sites through this one list makes that class of bug
// impossible: no path can drop every table and re-migrate a subset.
//
// The pointer/value receiver forms mirror exactly what boot has always used;
// keep any new system table in sync here rather than at a call site.
func SystemDatabaseModels() []interface{} {
	return []interface{}{
		&meshsyncmodel.KubernetesKeyValue{},
		&meshsyncmodel.KubernetesResource{},
		&meshsyncmodel.KubernetesResourceSpec{},
		&meshsyncmodel.KubernetesResourceStatus{},
		&meshsyncmodel.KubernetesResourceObjectMeta{},
		&PerformanceProfile{},
		&MesheryResult{},
		&MesheryPattern{},
		&MesheryFilter{},
		&PatternResource{},
		&MesheryApplication{},
		&UserPreference{},
		&UserCapabilities{},
		&PerformanceTestConfig{},
		&SmiResultWithID{},
		K8sContext{},
		schemasOrganization.Organization{},
		Key{},
		&Credential{},
		connections.Connection{},
		environment.Environment{},
		environment.EnvironmentConnectionMapping{},
		workspace.Workspace{},
		workspace.WorkspacesEnvironmentsMapping{},
		workspace.WorkspacesDesignsMapping{},
		workspace.WorkspacesTeamsMapping{},
		workspace.WorkspacesViewsMapping{},
		events.Event{},
		&SystemSetting{},
	}
}

// AutoMigrateSystemTables migrates every system table in SystemDatabaseModels
// onto db. Every path that (re)builds the system database must call this rather
// than migrating an inline list, so no path can ever again drop all tables and
// re-migrate only a subset.
func AutoMigrateSystemTables(db *database.Handler) error {
	return db.AutoMigrate(SystemDatabaseModels()...)
}
