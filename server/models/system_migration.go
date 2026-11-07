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
	if err := db.AutoMigrate(SystemDatabaseModels()...); err != nil {
		return err
	}
	return migrateEventCorrelationIndex(db)
}

// openEventCorrelationIndex is the partial unique index that keeps at most one
// open event (no resolved_at) per owner and correlation_key, which the
// createEvent correlation contract requires. EventsPersister.PersistEvent
// relies on it to make opening a condition atomic. Its predicate covers exactly
// the rows correlatesEvent correlates: rows without a correlation_key, and rows
// that identify no user (no owner or the zero owner), are recorded
// uncorrelated, so they are outside the index and never conflict.
const openEventCorrelationIndex = `CREATE UNIQUE INDEX IF NOT EXISTS idx_events_open_correlation
	ON events (owner, correlation_key)
	WHERE correlation_key IS NOT NULL AND correlation_key <> ''
		AND owner IS NOT NULL AND owner <> '00000000-0000-0000-0000-000000000000'
		AND resolved_at IS NULL`

// migrateEventCorrelationIndex creates openEventCorrelationIndex. GORM's tags
// cannot express a partial index, so it is created here, after AutoMigrate has
// added the correlation columns. The statement is idempotent and runs on every
// boot and reset.
func migrateEventCorrelationIndex(db *database.Handler) error {
	if err := db.Exec(openEventCorrelationIndex).Error; err != nil {
		return ErrMigrateEventCorrelationIndex(err)
	}
	return nil
}
