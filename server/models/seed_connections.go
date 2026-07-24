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

// registeredConnectionDefinitions reads every connection definition currently in
// the registry.
func registeredConnectionDefinitions(regm *meshmodel.RegistryManager) ([]*connectionv1beta3.ConnectionDefinition, error) {
	entities, _, _, err := regm.GetEntities(&regv1beta1.ConnectionFilter{})
	if err != nil {
		return nil, err
	}

	defs := make([]*connectionv1beta3.ConnectionDefinition, 0, len(entities))
	for _, en := range entities {
		if def, ok := en.(*connectionv1beta3.ConnectionDefinition); ok && def != nil {
			defs = append(defs, def)
		}
	}
	return defs, nil
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
	// wrote it. The canonical row is the one with the lowest id: ids are content
	// hashes of data that does not change between boots, so the same row is
	// picked every time, independent of registration order or how many registry
	// entries happen to point at each row.
	canonical := registrants[0]

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
	// A seeded Connection is owned by the system and uses its kind's API
	// anonymously; it must never carry a credential or an owner.
	if carriesID(conn.CredentialID) || carriesID(legacy.CredentialID) {
		updates["credential_id"] = uuid.Nil
	}
	if carriesID(conn.Owner) || carriesID(legacy.UserID) {
		updates["owner"] = uuid.Nil
		updates["user_id"] = uuid.Nil
	}

	// `status` is deliberately absent. The definition's status is the state a
	// *new* Connection starts in; once the Connection exists its state belongs
	// to the connection state machine, and re-asserting it on every boot would
	// undo a user who connected or ignored it.

	return updates
}

// carriesID reports whether an optional id column actually holds an identity.
func carriesID(id *core.Uuid) bool {
	return id != nil && *id != uuid.Nil
}
