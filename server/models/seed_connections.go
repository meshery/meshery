package models

import (
	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	connectionv1beta1 "github.com/meshery/schemas/models/v1beta1/connection"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
)

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

	seeded := 0
	for _, def := range defs {
		if !isSeedable(def, byKind) {
			continue
		}
		changed, err := seedConnectionForDefinition(db, def, byKind[def.Kind])
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

// registrantConnectionsByKind returns the Connections that Meshery itself
// created as registrants - the hosts that models, components and relationships
// are registered under - grouped by kind.
//
// Restricting the seed to these rows is what keeps it from touching a
// Connection a user created: a user's Connection is never a registrant, so it
// is never returned here and never rewritten.
func registrantConnectionsByKind(db *database.Handler) (map[string][]connectionv1beta1.Connection, error) {
	var conns []connectionv1beta1.Connection
	err := db.Table("connections").
		Where("EXISTS (SELECT 1 FROM registries WHERE registries.registrant_id = connections.id)").
		Find(&conns).Error
	if err != nil {
		return nil, err
	}

	byKind := make(map[string][]connectionv1beta1.Connection, len(conns))
	for _, conn := range conns {
		byKind[conn.Kind] = append(byKind[conn.Kind], conn)
	}
	return byKind, nil
}

// isSeedable reports whether a registered connection definition describes a
// Connection that Meshery seeds for itself at boot.
//
// Two rules, both read off the registry rather than off a list of kinds:
//
//  1. Meshery already sources content through the kind - it holds a registrant
//     Connection of that kind. This is what distinguishes Artifact Hub and
//     GitHub (Meshery generates models from both) from Kubernetes, Grafana and
//     Prometheus, which describe resources a user brings and which would be
//     meaningless as empty, endpoint-less rows.
//
//  2. The kind works anonymously - its credentialSchema marks nothing as
//     required. A seeded Connection is owned by the system and carries no
//     credential, so a kind that cannot be used without one must not be seeded.
func isSeedable(def *connectionv1beta3.ConnectionDefinition, registrantsByKind map[string][]connectionv1beta1.Connection) bool {
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

// seedConnectionForDefinition brings the registrant Connections of a kind in
// line with the kind's connection definition, and reports how many rows it had
// to write.
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
// carrying the stale one - it would leave two Connections per kind. Updating
// keeps exactly one, and keeps every registries.registrant_id foreign key
// pointing at it.
func seedConnectionForDefinition(db *database.Handler, def *connectionv1beta3.ConnectionDefinition, registrants []connectionv1beta1.Connection) (int, error) {
	changed := 0
	for _, conn := range registrants {
		updates := seedUpdatesFor(conn, def)
		if len(updates) == 0 {
			continue
		}
		if err := db.Model(&connectionv1beta1.Connection{}).
			Where("id = ?", conn.ID).
			Updates(updates).Error; err != nil {
			return changed, err
		}
		changed++
	}
	return changed, nil
}

// seedUpdatesFor returns the columns that have to change for a registrant
// Connection to match its definition, and an empty map when it already does -
// which is what makes a restart write nothing.
func seedUpdatesFor(conn connectionv1beta1.Connection, def *connectionv1beta3.ConnectionDefinition) map[string]interface{} {
	updates := map[string]interface{}{}

	if conn.Name != def.Name {
		updates["name"] = def.Name
	}
	if conn.Type != def.ConnectionType {
		updates["type"] = def.ConnectionType
	}
	if conn.SubType != def.SubType {
		updates["sub_type"] = def.SubType
	}
	// A seeded Connection is owned by the system and uses its kind's API
	// anonymously; it must never carry a credential or an owner.
	if conn.CredentialID != nil && *conn.CredentialID != uuid.Nil {
		updates["credential_id"] = uuid.Nil
	}
	if conn.UserID != nil && *conn.UserID != uuid.Nil {
		updates["user_id"] = uuid.Nil
	}

	// `status` is deliberately absent. The definition's status is the state a
	// *new* Connection starts in; once the Connection exists its state belongs
	// to the connection state machine, and re-asserting it on every boot would
	// undo a user who connected or ignored it.

	return updates
}
