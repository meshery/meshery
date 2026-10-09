package models

import (
	"sort"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	"github.com/meshery/schemas/models/core"
	connectionv1beta1 "github.com/meshery/schemas/models/v1beta1/connection"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
)

// registrantScope narrows a `connections` query to the rows Meshery created as
// registrants of registered Models - the hosts that Models, and the Components
// and Relationships beneath them, are registered under.
//
// Scoping to registrants is what keeps seeding from touching a Connection a
// user created: a user's Connection is never a registrant.
//
// Narrowing to `model` entities is what keeps the scope Meshery-internal. A
// registrant row on its own is not a Meshery-only signal - POST
// /api/registry/connections registers a caller-supplied connection definition
// under a caller-supplied registrant kind, which creates a registrant row owning
// a single `connection` entity. Requiring registered Models makes the rule
// uninducible from request input while keeping it registry-derived rather than
// a list of kinds.
const registrantScope = "EXISTS (SELECT 1 FROM registries WHERE registries.registrant_id = connections.id AND registries.type = ?)"

// SeedConnections materializes the system-owned Connections that Meshery holds
// out of the box, so a user finds them already present instead of having to
// create them by hand.
//
// It is driven entirely by the registry: the set of Connections seeded is
// derived from the registered connection definitions
// (models/<model>/<version>/connections/*.json), never from a list of kinds
// enumerated here. Shipping a definition that satisfies the two rules in
// isSeedable is all that is required to have its Connection seeded.
//
// Must run after model registration has completed - it reads both the
// connection definitions and the registrant Connections that registration
// creates. SeedComponents calls it at the end for exactly that reason.
func SeedConnections(log logger.Handler, db *database.Handler, regm *meshmodel.RegistryManager) {
	if db == nil || regm == nil {
		missing := make([]string, 0, 2)
		if db == nil {
			missing = append(missing, "database handler")
		}
		if regm == nil {
			missing = append(missing, "registry manager")
		}
		log.Warnf("Skipping seeding of system-owned connections: %s unavailable", strings.Join(missing, " and "))
		return
	}

	defs, err := registeredConnectionDefinitions(regm)
	if err != nil {
		log.Error(ErrSeedingConnections(err))
		return
	}

	byKind, err := registrantConnectionsByKind(db)
	if err != nil {
		log.Error(ErrSeedingConnections(err))
		return
	}

	legacyByID, err := legacyRegistrantColumnsByID(db)
	if err != nil {
		log.Error(ErrSeedingConnections(err))
		return
	}

	seeded := 0
	for _, def := range defs {
		if !isSeedable(def, byKind) {
			continue
		}
		changed, err := seedConnectionForDefinition(db, def, byKind[def.Kind], legacyByID)
		if err != nil {
			// One kind failing to seed must not stop the rest: a partially
			// seeded registry is still better than none, and the next boot
			// retries whatever did not land.
			log.Error(ErrSeedingConnectionKind(err, def.Kind))
			continue
		}
		seeded += changed
	}

	// Steady state is zero: every seeded Connection already matches its
	// definition, so a restart writes nothing.
	if seeded > 0 {
		log.Infof("Seeded %d system-owned connection(s) from registered connection definitions", seeded)
	}
}

// registeredConnectionDefinitions reads the connection definitions currently in
// the registry, reduced to one per kind and ordered by kind.
//
// The registry accumulates copies: ConnectionDefinition.Create mints a fresh
// random id and inserts unconditionally, and every boot re-registers every model
// directory, so connection_definition_dbs grows by a row per definition per
// boot. Seeding is a per-kind operation, so the copies are reduced here rather
// than left to re-issue the same UPDATE once per copy and report copies instead
// of rows.
//
// The surviving copy is the most recently registered one - a later models
// release supersedes an accumulated older row - with the definition's own
// identity as the tie-break, so the choice is content-derived and cannot
// oscillate between boots even when accumulated copies of a kind disagree.
func registeredConnectionDefinitions(regm *meshmodel.RegistryManager) ([]*connectionv1beta3.ConnectionDefinition, error) {
	entities, _, _, err := regm.GetEntities(&regv1beta1.ConnectionFilter{})
	if err != nil {
		return nil, err
	}

	byKind := make(map[string]*connectionv1beta3.ConnectionDefinition, len(entities))
	for _, en := range entities {
		def, ok := en.(*connectionv1beta3.ConnectionDefinition)
		if !ok || def == nil || def.Kind == "" {
			continue
		}
		if incumbent, seen := byKind[def.Kind]; !seen || supersedes(def, incumbent) {
			byKind[def.Kind] = def
		}
	}

	defs := make([]*connectionv1beta3.ConnectionDefinition, 0, len(byKind))
	for _, def := range byKind {
		defs = append(defs, def)
	}
	sort.Slice(defs, func(i, j int) bool { return defs[i].Kind < defs[j].Kind })
	return defs, nil
}

// supersedes reports whether an accumulated copy of a connection definition
// should replace the copy already chosen for its kind.
func supersedes(candidate, incumbent *connectionv1beta3.ConnectionDefinition) bool {
	if !candidate.CreatedAt.Equal(incumbent.CreatedAt) {
		return candidate.CreatedAt.After(incumbent.CreatedAt)
	}
	return definitionIdentity(candidate) < definitionIdentity(incumbent)
}

// definitionIdentity is the identity a definition stamps onto its Connection.
func definitionIdentity(def *connectionv1beta3.ConnectionDefinition) string {
	return def.Name + "\x00" + def.ConnectionType + "\x00" + def.SubType
}

// registrantConnectionsByKind returns the registrant Connections (see
// registrantScope) grouped by kind and ordered deterministically by id.
//
// The read is through the canonical v1beta3 Connection - the same model
// ConnectionPersister serves the API from - so the columns compared against a
// definition are the columns that actually reach the wire.
func registrantConnectionsByKind(db *database.Handler) (map[string][]connectionv1beta3.Connection, error) {
	var conns []connectionv1beta3.Connection
	err := db.Table("connections").
		Where(registrantScope, string(entity.Model)).
		Find(&conns).Error
	if err != nil {
		return nil, err
	}

	byKind := make(map[string][]connectionv1beta3.Connection, len(conns))
	for _, conn := range conns {
		byKind[conn.Kind] = append(byKind[conn.Kind], conn)
	}
	for kind := range byKind {
		ofKind := byKind[kind]
		sort.Slice(ofKind, func(i, j int) bool {
			return ofKind[i].ID.String() < ofKind[j].ID.String()
		})
	}
	return byKind, nil
}

// legacyRegistrantColumnsByID reads the same registrant rows through the v1beta1
// Connection, keyed by id.
//
// The `connections` table carries two column sets for the same fields: meshkit's
// registry manager AutoMigrates the v1beta1 Connection (`type`, `user_id`) while
// the server AutoMigrates the canonical v1beta3 Connection (`connection_type`,
// `owner`). Reads that serve the API scan v1beta3, but
// ConnectionPersister.GetConnections filters with raw SQL on `type`/`sub_type`,
// so both halves matter. Seeding reads and writes them as a pair so the filter
// and the payload cannot disagree.
func legacyRegistrantColumnsByID(db *database.Handler) (map[core.Uuid]connectionv1beta1.Connection, error) {
	var conns []connectionv1beta1.Connection
	err := db.Table("connections").
		Where(registrantScope, string(entity.Model)).
		Find(&conns).Error
	if err != nil {
		return nil, err
	}

	byID := make(map[core.Uuid]connectionv1beta1.Connection, len(conns))
	for _, conn := range conns {
		byID[conn.ID] = conn
	}
	return byID, nil
}

// isSeedable reports whether a registered connection definition describes a
// Connection that Meshery seeds for itself at boot.
//
// Two rules, both read off the registry rather than off a list of kinds:
//
//  1. Meshery already sources content through the kind - it holds a registrant
//     Connection of that kind that owns registered Models. This is what
//     distinguishes Artifact Hub and GitHub (Meshery generates models from both)
//     from Kubernetes, Grafana and Prometheus, which describe resources a user
//     brings and which would be meaningless as empty, endpoint-less rows.
//
//  2. The kind works anonymously - its credentialSchema marks nothing as
//     required. A seeded Connection is owned by the system and carries no
//     credential, so a kind that cannot be used without one must not be seeded.
func isSeedable(def *connectionv1beta3.ConnectionDefinition, registrantsByKind map[string][]connectionv1beta3.Connection) bool {
	if def == nil || def.Kind == "" {
		return false
	}
	if len(registrantsByKind[def.Kind]) == 0 {
		return false
	}
	return !requiresCredential(def)
}

// requiresCredential reports whether the definition's credentialSchema declares
// any required property. A definition with no credentialSchema, or one whose
// every property is optional, describes a kind usable anonymously.
func requiresCredential(def *connectionv1beta3.ConnectionDefinition) bool {
	if def == nil || def.CredentialSchema == nil {
		return false
	}
	// `required` survives the JSON round-trip through the registry as
	// []interface{}; accept the typed form too in case it is set in Go.
	switch required := def.CredentialSchema["required"].(type) {
	case []interface{}:
		return len(required) > 0
	case []string:
		return len(required) > 0
	}
	return false
}

// seedConnectionForDefinition brings the canonical registrant Connection of a
// kind in line with the kind's connection definition, and reports how many rows
// it had to write.
//
// The registrant rows are created by model registration from the `registrant`
// blob hand-authored into every model.json, which predates connection
// definitions and disagrees with them (it types both Artifact Hub and GitHub as
// `registry`/`""` and names the latter "Github"). The definition is the
// authoritative identity, so it wins.
//
// Writing in place rather than inserting is deliberate. Connection IDs are
// content-addressed (an md5 over the whole connection), so inserting a row
// carrying the definition's identity would not collide with the registrant row
// carrying the stale one - it would add a Connection rather than correct one.
// Updating leaves the row count exactly as registration left it, and keeps every
// registries.registrant_id foreign key pointing at it.
func seedConnectionForDefinition(
	db *database.Handler,
	def *connectionv1beta3.ConnectionDefinition,
	registrants []connectionv1beta3.Connection,
	legacyByID map[core.Uuid]connectionv1beta1.Connection,
) (int, error) {
	if len(registrants) == 0 {
		return 0, nil
	}

	// A kind can already hold more than one registrant row. The `registrant`
	// blob is hand-authored per model.json and the shipped blobs are not
	// uniform - some carry `user_id`, some do not - so their content-addressed
	// ids differ and registration creates one row per variant. Collapsing those
	// rows would have to repoint registries.registrant_id, which this boot-time
	// step must not do; the duplication is pre-existing and tracked as
	// meshery/meshery#20950.
	//
	// Seeding therefore stamps the definition's identity onto exactly one
	// canonical registrant and leaves every sibling exactly as registration
	// wrote it.
	canonical := canonicalRegistrantFor(registrants, def)

	updates := seedUpdatesFor(canonical, legacyByID[canonical.ID], def)
	if len(updates) == 0 {
		return 0, nil
	}
	if err := db.Model(&connectionv1beta3.Connection{}).
		Where("id = ?", canonical.ID).
		Updates(updates).Error; err != nil {
		return 0, err
	}
	return 1, nil
}

// canonicalRegistrantFor picks the one registrant row of a kind that carries the
// kind's definition identity.
//
// The pick is sticky: a registrant that already carries the identity wins, and
// only when none does is the lowest id taken (registrants arrive ordered by id).
// Stickiness is what keeps the pick stable as the registrant set GROWS - a later
// models release can introduce a new `registrant` spelling whose content hash
// sorts below the stamped row, and a pure lowest-id rule would move the identity
// onto it while the previously stamped row kept the definition's name, type and
// subType forever, leaving two identical-looking Connections for the kind with
// no self-healing. Siblings are never un-stamped or otherwise mutated.
func canonicalRegistrantFor(registrants []connectionv1beta3.Connection, def *connectionv1beta3.ConnectionDefinition) connectionv1beta3.Connection {
	for _, conn := range registrants {
		if carriesDefinitionIdentity(conn, def) {
			return conn
		}
	}
	return registrants[0]
}

// carriesDefinitionIdentity reports whether a Connection already presents the
// identity its kind's definition declares.
func carriesDefinitionIdentity(conn connectionv1beta3.Connection, def *connectionv1beta3.ConnectionDefinition) bool {
	return conn.Name == def.Name &&
		conn.ConnectionType == def.ConnectionType &&
		conn.SubType == def.SubType
}

// seedUpdatesFor returns the columns that have to change for a registrant
// Connection to match its definition, and an empty map when it already does -
// which is what makes a restart write nothing.
//
// The duplicated v1beta1/v1beta3 columns are written as a pair and treated as
// stale when either half disagrees, so the raw-SQL filter and the canonical
// payload cannot drift apart.
func seedUpdatesFor(conn connectionv1beta3.Connection, legacy connectionv1beta1.Connection, def *connectionv1beta3.ConnectionDefinition) map[string]interface{} {
	updates := map[string]interface{}{}

	if conn.Name != def.Name {
		updates["name"] = def.Name
	}
	if conn.ConnectionType != def.ConnectionType || legacy.Type != def.ConnectionType {
		updates["connection_type"] = def.ConnectionType
		updates["type"] = def.ConnectionType
	}
	if conn.SubType != def.SubType {
		updates["sub_type"] = def.SubType
	}
	// A seeded Connection is owned by the system. No path creates a registrant
	// with an owner and none sets one afterwards, so asserting it can only
	// repair a row, never contend with a user's action.
	if carriesID(conn.Owner) || carriesID(legacy.UserID) {
		updates["owner"] = uuid.Nil
		updates["user_id"] = uuid.Nil
	}

	// `status` and `credential_id` are deliberately absent. Both describe the
	// state a *new* Connection starts in, and both belong to the user once the
	// Connection exists: re-asserting `status` on every boot would undo a user
	// who connected or ignored it, and clearing `credential_id` would silently
	// revoke the optional Artifact Hub API key or GitHub token they attached
	// through the wizard - a seeded kind is seedable precisely because its
	// credential is optional, so attaching one is a legitimate user action.

	return updates
}

// carriesID reports whether an optional id column actually holds an identity.
func carriesID(id *core.Uuid) bool {
	return id != nil && *id != uuid.Nil
}
