package models

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/models/registration"
	"github.com/meshery/schemas/models/core"
	connectionv1beta1 "github.com/meshery/schemas/models/v1beta1/connection"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
	"github.com/sirupsen/logrus"
)

// resolveModelDir returns a shipped model's registration directory
// (models/<name>/<version>/v1.0.0) without pinning a specific version, so the
// test keeps working as the model registry is re-synced. meshery-core and
// kubernetes are both retained by the recommended sparse checkout, so this keeps
// the default `go test ./...` passing against a sparse clone.
func resolveModelDir(t *testing.T, name string) string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join("..", "..", "models", name, "*", "v1.0.0", "model.json"))
	if err != nil {
		t.Fatalf("glob model dir for %q: %v", name, err)
	}
	if len(matches) == 0 {
		t.Fatalf("no shipped model directory found under ../../models/%s/*/v1.0.0 - is it excluded by a sparse checkout?", name)
	}
	sort.Strings(matches)
	return filepath.Dir(matches[len(matches)-1])
}

// mesheryCoreModelDir resolves the meshery-core model, which carries the shipped
// connection definitions and registers under the `meshery` registrant.
func mesheryCoreModelDir(t *testing.T) string { return resolveModelDir(t, "meshery-core") }

// gitHubRegistrantModelDir resolves the kubernetes model, which registers under
// the `github` registrant - the registrant that makes the GitHub connection
// definition seedable.
func gitHubRegistrantModelDir(t *testing.T) string { return resolveModelDir(t, "kubernetes") }

// The Artifact Hub registrant models are committed fixtures, not shipped models:
// every Artifact Hub model lives under models/, which the recommended sparse
// checkout excludes, so depending on one would break `go test ./...` on a sparse
// clone. The registrant blobs are deliberately not uniform -
// artifactHubOwnedRegistrantModelDir carries a `user_id`, so it hashes to a
// different registrant row than artifactHubRegistrantModelDir. That is the
// two-spellings case a kind legitimately holds (meshery/meshery#20950), which
// lets the canonical-pick test exercise real duplication rather than a synthetic
// fixture.
const (
	artifactHubRegistrantModelDir      = "testdata/artifacthub-registrant-model"
	artifactHubOwnedRegistrantModelDir = "testdata/artifacthub-registrant-model-owned"
)

func newSeedTestLogger(t *testing.T) logger.Handler {
	t.Helper()
	log, err := logger.New("seed-connections-test", logger.Options{Format: logger.JsonLogFormat, LogLevel: int(logrus.ErrorLevel)})
	if err != nil {
		t.Fatalf("build logger: %v", err)
	}
	return log
}

// seedTestRegistry registers the given model directories into a fresh in-memory
// registry and returns it alongside its database handler.
//
// The canonical v1beta3 Connection is migrated alongside the registry manager's
// v1beta1 one, mirroring what cmd/main.go does on boot: the `connections` table
// carries both column sets, and a fixture holding only the legacy half would let
// a legacy-only write pass unnoticed.
func seedTestRegistry(t *testing.T, modelDirs ...string) (*database.Handler, *meshmodel.RegistryManager) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	regm, err := meshmodel.NewRegistryManager(&db)
	if err != nil {
		t.Fatalf("new registry manager: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}); err != nil {
		t.Fatalf("migrate canonical connection: %v", err)
	}

	regHelper := registration.NewRegistrationHelper(t.TempDir(), regm, NewRegistrationFailureLogHandler())
	for _, dir := range modelDirs {
		regHelper.Register(registration.NewDir(dir))
	}
	return &db, regm
}

// connectionsByKind reads the connections table through the canonical v1beta3
// model - the same model ConnectionPersister serves the API from - and groups
// the rows by kind. Asserting through this path is what makes a write that lands
// only in the legacy v1beta1 columns fail.
func connectionsByKind(t *testing.T, db *database.Handler) map[string][]connectionv1beta3.Connection {
	t.Helper()

	var conns []connectionv1beta3.Connection
	if err := db.Find(&conns).Error; err != nil {
		t.Fatalf("read connections: %v", err)
	}

	byKind := make(map[string][]connectionv1beta3.Connection, len(conns))
	for _, conn := range conns {
		byKind[conn.Kind] = append(byKind[conn.Kind], conn)
	}
	return byKind
}

// soleConnection returns the single Connection of a kind, failing the test when
// the kind is missing or holds more than one row.
func soleConnection(t *testing.T, byKind map[string][]connectionv1beta3.Connection, kind string) connectionv1beta3.Connection {
	t.Helper()

	conns := byKind[kind]
	switch len(conns) {
	case 1:
		return conns[0]
	case 0:
		t.Fatalf("no connection row for kind %q", kind)
	default:
		t.Fatalf("expected exactly one connection row for kind %q, got %d", kind, len(conns))
	}
	return connectionv1beta3.Connection{}
}

// legacyConnection reads the duplicated v1beta1 columns (`type`, `user_id`) of a
// row, so a test can assert that both halves of the split schema were written.
func legacyConnection(t *testing.T, db *database.Handler, id core.Uuid) connectionv1beta1.Connection {
	t.Helper()

	var conn connectionv1beta1.Connection
	if err := db.Table("connections").Where("id = ?", id).First(&conn).Error; err != nil {
		t.Fatalf("read legacy columns of %s: %v", id, err)
	}
	return conn
}

// registerRegistrant creates a registrant Connection of the given kind owning one
// registry entry of the given entity type, the way registration itself does.
func registerRegistrant(t *testing.T, db *database.Handler, kind string, entityType entity.EntityType) core.Uuid {
	t.Helper()

	registrant := connectionv1beta1.Connection{Kind: kind, Type: "registry"}
	registrantID, err := registrant.Create(db)
	if err != nil {
		t.Fatalf("create %s registrant: %v", kind, err)
	}
	linkRegistryEntry(t, db, registrantID, entityType)
	return registrantID
}

// insertRegistrant creates a registrant Connection under an explicitly chosen id,
// so a test can place a row at a specific point in the id ordering.
func insertRegistrant(t *testing.T, db *database.Handler, id core.Uuid, kind string, entityType entity.EntityType) {
	t.Helper()

	registrant := connectionv1beta1.Connection{ID: id, Kind: kind, Type: "registry"}
	if err := db.Create(&registrant).Error; err != nil {
		t.Fatalf("create %s registrant %s: %v", kind, id, err)
	}
	linkRegistryEntry(t, db, id, entityType)
}

func linkRegistryEntry(t *testing.T, db *database.Handler, registrantID core.Uuid, entityType entity.EntityType) {
	t.Helper()

	if err := db.Create(&meshmodel.Registry{
		ID:           uuid.Must(uuid.NewV4()),
		RegistrantID: registrantID,
		Entity:       uuid.Must(uuid.NewV4()),
		Type:         entityType,
	}).Error; err != nil {
		t.Fatalf("create registry entry for %s: %v", registrantID, err)
	}
}

// TestSeedConnectionsMaterializesDefinitionBackedConnections is the core
// behavior: after boot the Artifact Hub and GitHub connections exist carrying
// the identity their connection definitions declare, not the stale identity the
// model.json registrant blobs carry.
func TestSeedConnectionsMaterializesDefinitionBackedConnections(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t), artifactHubRegistrantModelDir, gitHubRegistrantModelDir(t))

	// Registration alone leaves the registrant rows disagreeing with the
	// definitions; assert that starting point so the test proves seeding did
	// the work rather than the registrant blobs happening to be right. The
	// registrant is written through the v1beta1 model, so the canonical
	// `connection_type` column is not populated at all.
	before := connectionsByKind(t, db)
	if got := soleConnection(t, before, "github").Name; got != "Github" {
		t.Fatalf("precondition: expected registrant name %q, got %q", "Github", got)
	}
	if got := soleConnection(t, before, "artifacthub").ConnectionType; got != "" {
		t.Fatalf("precondition: expected an unset canonical connection type, got %q", got)
	}
	if got := legacyConnection(t, db, soleConnection(t, before, "artifacthub").ID).Type; got != "registry" {
		t.Fatalf("precondition: expected legacy registrant type %q, got %q", "registry", got)
	}

	SeedConnections(newSeedTestLogger(t), db, regm)

	after := connectionsByKind(t, db)
	for _, tc := range []struct {
		kind, name, connType, subType string
	}{
		{kind: "artifacthub", name: "Artifact Hub", connType: "source", subType: "registry"},
		{kind: "github", name: "GitHub", connType: "source", subType: "git"},
	} {
		conn := soleConnection(t, after, tc.kind)
		if conn.Name != tc.name {
			t.Errorf("kind %q: name = %q, want %q", tc.kind, conn.Name, tc.name)
		}
		// The canonical column is what the API and UI render.
		if conn.ConnectionType != tc.connType {
			t.Errorf("kind %q: connection_type = %q, want %q", tc.kind, conn.ConnectionType, tc.connType)
		}
		if conn.SubType != tc.subType {
			t.Errorf("kind %q: subType = %q, want %q", tc.kind, conn.SubType, tc.subType)
		}
		// The legacy column is what ConnectionPersister.GetConnections filters
		// on; both halves have to agree or filtering and rendering disagree.
		legacy := legacyConnection(t, db, conn.ID)
		if legacy.Type != tc.connType {
			t.Errorf("kind %q: legacy type = %q, want %q", tc.kind, legacy.Type, tc.connType)
		}
		// System-owned and anonymous: no credential, no owner, on either half.
		if carriesID(conn.CredentialID) || carriesID(legacy.CredentialID) {
			t.Errorf("kind %q: seeded connection carries a credential", tc.kind)
		}
		if carriesID(conn.Owner) || carriesID(legacy.UserID) {
			t.Errorf("kind %q: seeded connection carries an owner", tc.kind)
		}
	}
}

// TestSeedConnectionsIsIdempotentAcrossRestarts proves the restart behavior the
// issue calls for: no new rows, and no write at all on a second boot.
func TestSeedConnectionsIsIdempotentAcrossRestarts(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t), artifactHubRegistrantModelDir, gitHubRegistrantModelDir(t))
	log := newSeedTestLogger(t)

	SeedConnections(log, db, regm)
	first := connectionsByKind(t, db)

	// A restart re-runs registration before seeding, exactly as SeedComponents
	// does, so exercise both.
	regHelper := registration.NewRegistrationHelper(t.TempDir(), regm, NewRegistrationFailureLogHandler())
	for _, dir := range []string{mesheryCoreModelDir(t), artifactHubRegistrantModelDir, gitHubRegistrantModelDir(t)} {
		regHelper.Register(registration.NewDir(dir))
	}
	SeedConnections(log, db, regm)

	second := connectionsByKind(t, db)
	if len(second) != len(first) {
		t.Fatalf("restart changed the number of connection kinds: %d -> %d", len(first), len(second))
	}
	for kind, before := range first {
		after, ok := second[kind]
		if !ok {
			t.Fatalf("kind %q disappeared across restart", kind)
		}
		if len(after) != len(before) {
			t.Errorf("kind %q: row count changed across restart: %d -> %d", kind, len(before), len(after))
			continue
		}
		for i := range before {
			if after[i].ID != before[i].ID {
				t.Errorf("kind %q: id changed across restart: %s -> %s", kind, before[i].ID, after[i].ID)
			}
			if after[i].Name != before[i].Name || after[i].ConnectionType != before[i].ConnectionType || after[i].SubType != before[i].SubType {
				t.Errorf("kind %q: identity changed across restart: %+v -> %+v", kind, before[i], after[i])
			}
		}
	}

	// The second boot must do no work at all: every canonical registrant already
	// matches its definition, on both halves of the split schema, so there is
	// nothing left to update.
	assertSteadyStateWritesNothing(t, db, regm)
}

// runSeedingPass performs one seeding pass exactly as SeedConnections does and
// reports how many rows it wrote, which SeedConnections itself only surfaces as
// a log line.
func runSeedingPass(t *testing.T, db *database.Handler, regm *meshmodel.RegistryManager) int {
	t.Helper()

	byKind, err := registrantConnectionsByKind(db)
	if err != nil {
		t.Fatalf("read registrant connections: %v", err)
	}
	legacyByID, err := legacyRegistrantColumnsByID(db)
	if err != nil {
		t.Fatalf("read legacy registrant columns: %v", err)
	}
	defs, err := registeredConnectionDefinitions(regm)
	if err != nil {
		t.Fatalf("read connection definitions: %v", err)
	}

	written := 0
	for _, def := range defs {
		if !isSeedable(def, byKind) {
			continue
		}
		changed, err := seedConnectionForDefinition(db, def, byKind[def.Kind], legacyByID)
		if err != nil {
			t.Fatalf("kind %q: %v", def.Kind, err)
		}
		written += changed
	}
	return written
}

// assertSteadyStateWritesNothing re-runs the seeding decision and fails if it
// would write anything.
func assertSteadyStateWritesNothing(t *testing.T, db *database.Handler, regm *meshmodel.RegistryManager) {
	t.Helper()

	if written := runSeedingPass(t, db, regm); written != 0 {
		t.Errorf("steady-state seeding wrote %d row(s), want 0", written)
	}
}

// TestSeedConnectionsStampsOneCanonicalRegistrantPerKind covers the case the
// single-model-per-kind fixtures cannot reach: the shipped `registrant` blobs
// are not uniform, so a kind legitimately holds more than one registrant row
// (meshery/meshery#20950). Seeding must stamp exactly one - deterministically,
// the lowest id when none is stamped yet - and leave the sibling exactly as
// registration wrote it, rather than turning both rows into indistinguishable
// copies of the definition.
func TestSeedConnectionsStampsOneCanonicalRegistrantPerKind(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t), artifactHubRegistrantModelDir, artifactHubOwnedRegistrantModelDir)

	before := connectionsByKind(t, db)
	if got := len(before["artifacthub"]); got != 2 {
		t.Fatalf("precondition: expected the two registrant blob variants to produce 2 artifacthub rows, got %d", got)
	}

	SeedConnections(newSeedTestLogger(t), db, regm)

	after := connectionsByKind(t, db)
	if got := len(after["artifacthub"]); got != 2 {
		t.Fatalf("seeding changed the artifacthub row count: 2 -> %d", got)
	}

	registrants, err := registrantConnectionsByKind(db)
	if err != nil {
		t.Fatalf("read registrant connections: %v", err)
	}
	canonicalID := registrants["artifacthub"][0].ID
	for _, other := range registrants["artifacthub"][1:] {
		if other.ID.String() < canonicalID.String() {
			t.Fatalf("canonical pick is not the lowest id: picked %s, saw %s", canonicalID, other.ID)
		}
	}

	stamped := 0
	for _, conn := range after["artifacthub"] {
		legacy := legacyConnection(t, db, conn.ID)
		if conn.ID == canonicalID {
			stamped++
			if conn.ConnectionType != "source" || conn.SubType != "registry" {
				t.Errorf("canonical registrant %s was not stamped: type=%q subType=%q", conn.ID, conn.ConnectionType, conn.SubType)
			}
			if legacy.Type != "source" {
				t.Errorf("canonical registrant %s: legacy type = %q, want %q", conn.ID, legacy.Type, "source")
			}
			continue
		}
		// The sibling is left exactly as registration wrote it.
		if conn.ConnectionType != "" || conn.SubType != "" {
			t.Errorf("sibling registrant %s was rewritten: type=%q subType=%q", conn.ID, conn.ConnectionType, conn.SubType)
		}
		if legacy.Type != "registry" {
			t.Errorf("sibling registrant %s: legacy type = %q, want %q", conn.ID, legacy.Type, "registry")
		}
	}
	if stamped != 1 {
		t.Errorf("stamped %d artifacthub registrants, want exactly 1", stamped)
	}

	// Picking one row of several must still be a fixed point across restarts.
	assertSteadyStateWritesNothing(t, db, regm)
}

// TestSeedConnectionsCanonicalPickSurvivesRegistrantGrowth covers the stability
// the lowest-id rule alone does not give. Registrant ids are content hashes of
// the hand-authored model.json `registrant` blob, so a later models release can
// introduce a spelling that hashes BELOW the row already stamped. The stamp must
// not migrate to it: doing so would leave the previously stamped row carrying
// the definition identity forever, producing the duplicate-identity outcome the
// single-canonical rule exists to prevent, with no self-healing on later boots.
func TestSeedConnectionsCanonicalPickSurvivesRegistrantGrowth(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t), artifactHubRegistrantModelDir)
	log := newSeedTestLogger(t)

	SeedConnections(log, db, regm)
	stampedID := soleConnection(t, connectionsByKind(t, db), "artifacthub").ID

	// A registrant of the same kind whose id sorts below every content hash.
	lowestID := uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))
	insertRegistrant(t, db, lowestID, "artifacthub", entity.Model)
	if lowestID.String() >= stampedID.String() {
		t.Fatalf("precondition: new registrant %s must sort below the stamped one %s", lowestID, stampedID)
	}

	SeedConnections(log, db, regm)

	after := connectionsByKind(t, db)
	if got := len(after["artifacthub"]); got != 2 {
		t.Fatalf("expected the grown registrant set to still be 2 rows, got %d", got)
	}

	carriers := []core.Uuid{}
	for _, conn := range after["artifacthub"] {
		if conn.Name == "Artifact Hub" && conn.ConnectionType == "source" && conn.SubType == "registry" {
			carriers = append(carriers, conn.ID)
		}
	}
	if len(carriers) != 1 {
		t.Fatalf("expected exactly one row carrying the definition identity, got %d: %v", len(carriers), carriers)
	}
	if carriers[0] != stampedID {
		t.Errorf("the stamp moved on registrant growth: %s -> %s", stampedID, carriers[0])
	}

	assertSteadyStateWritesNothing(t, db, regm)
}

// TestSeedConnectionsWritesOncePerKindWithDuplicateDefinitions guards against
// accumulated connection definitions multiplying the work. Registration inserts
// a definition row unconditionally under a freshly minted random id, so every
// boot adds another copy of every definition; without reducing them by kind the
// same canonical row is rewritten once per copy in a single pass and the seeded
// count reports copies rather than rows.
func TestSeedConnectionsWritesOncePerKindWithDuplicateDefinitions(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t), artifactHubRegistrantModelDir, gitHubRegistrantModelDir(t))

	// Re-registering the definitions is what a restart does, and what leaves a
	// second copy of every definition behind.
	regHelper := registration.NewRegistrationHelper(t.TempDir(), regm, NewRegistrationFailureLogHandler())
	regHelper.Register(registration.NewDir(mesheryCoreModelDir(t)))

	var copies int64
	if err := db.Table("connection_definition_dbs").Count(&copies).Error; err != nil {
		t.Fatalf("count connection definitions: %v", err)
	}
	defs, err := registeredConnectionDefinitions(regm)
	if err != nil {
		t.Fatalf("read connection definitions: %v", err)
	}
	if copies <= int64(len(defs)) {
		t.Fatalf("precondition: expected accumulated definition copies, got %d rows for %d kinds", copies, len(defs))
	}

	seenKinds := map[string]bool{}
	for _, def := range defs {
		if seenKinds[def.Kind] {
			t.Errorf("kind %q returned more than once after de-duplication", def.Kind)
		}
		seenKinds[def.Kind] = true
	}

	if written := runSeedingPass(t, db, regm); written != 2 {
		t.Errorf("seeding wrote %d row(s) for 2 seedable kinds, want 2", written)
	}
	assertSteadyStateWritesNothing(t, db, regm)
}

// TestSeedConnectionsPreservesUserAttachedCredential mirrors the status
// carve-out for credentials. Both seedable kinds are seedable precisely because
// their credentialSchema marks nothing required - Artifact Hub's API key and
// GitHub's token are optional and raise the rate limit - so attaching one is a
// legitimate user action that a restart must not silently revoke.
func TestSeedConnectionsPreservesUserAttachedCredential(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t), artifactHubRegistrantModelDir, gitHubRegistrantModelDir(t))
	log := newSeedTestLogger(t)

	SeedConnections(log, db, regm)

	credentialID := uuid.Must(uuid.NewV4())
	gitHubID := soleConnection(t, connectionsByKind(t, db), "github").ID
	if err := db.Model(&connectionv1beta3.Connection{}).
		Where("id = ?", gitHubID).
		Update("credential_id", credentialID).Error; err != nil {
		t.Fatalf("attach credential: %v", err)
	}

	SeedConnections(log, db, regm)

	gitHub := soleConnection(t, connectionsByKind(t, db), "github")
	if gitHub.CredentialID == nil || *gitHub.CredentialID != credentialID {
		t.Errorf("credential_id = %v, want it preserved as %s", gitHub.CredentialID, credentialID)
	}
	assertSteadyStateWritesNothing(t, db, regm)
}

// TestSeedConnectionsSkipsKindsWithoutRegistrants guards the scoping rule: a
// kind that ships a definition but that Meshery does not itself source content
// through must not be materialized as an empty, endpoint-less row.
func TestSeedConnectionsSkipsKindsWithoutRegistrants(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t), artifactHubRegistrantModelDir, gitHubRegistrantModelDir(t))

	SeedConnections(newSeedTestLogger(t), db, regm)

	// Every shipped definition kind, so this fails if a future definition
	// silently starts being seeded.
	defs, err := registeredConnectionDefinitions(regm)
	if err != nil {
		t.Fatalf("read connection definitions: %v", err)
	}
	if len(defs) == 0 {
		t.Fatal("expected the shipped connection definitions to be registered")
	}

	seededKinds := connectionsByKind(t, db)
	for _, kind := range []string{"kubernetes", "grafana", "prometheus"} {
		if _, ok := seededKinds[kind]; ok {
			t.Errorf("kind %q has a definition but no registrant; it must not be seeded", kind)
		}
	}
}

// TestSeedConnectionsSkipsKubernetesRegistrant covers the case the credential
// rule actually exists for. A user who imports Models registered under a
// `kubernetes` registrant satisfies the registrant rule, so only the credential
// rule stands between Kubernetes and being seeded. Kubernetes cannot be used
// anonymously - its credentialSchema requires a kubeconfig - so it must still
// never be seeded, and its registrant must be left exactly as registration
// wrote it.
func TestSeedConnectionsSkipsKubernetesRegistrant(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t), artifactHubRegistrantModelDir)

	registerRegistrant(t, db, "kubernetes", entity.Model)

	SeedConnections(newSeedTestLogger(t), db, regm)

	after := connectionsByKind(t, db)
	k8s := soleConnection(t, after, "kubernetes")
	// Untouched: still the registrant identity, not the definition's
	// platform/orchestration identity.
	if k8s.ConnectionType != "" || k8s.SubType != "" {
		t.Errorf("kubernetes registrant was seeded: connection_type=%q subType=%q", k8s.ConnectionType, k8s.SubType)
	}
	if got := legacyConnection(t, db, k8s.ID).Type; got != "registry" {
		t.Errorf("kubernetes registrant was seeded: legacy type = %q, want %q", got, "registry")
	}

	// The anonymous kind alongside it still seeds, so this proves the credential
	// rule did the skipping rather than seeding having failed wholesale.
	artifactHub := soleConnection(t, after, "artifacthub")
	if artifactHub.ConnectionType != "source" || artifactHub.SubType != "registry" {
		t.Errorf("artifacthub connection_type/subType = %q/%q, want source/registry", artifactHub.ConnectionType, artifactHub.SubType)
	}
}

// TestSeedConnectionsIgnoresConnectionDefinitionOnlyRegistrant pins the rule
// that keeps the registrant signal Meshery-internal. POST
// /api/registry/connections registers a caller-supplied connection definition
// under a caller-supplied registrant kind, which leaves a registrant row owning
// nothing but that `connection` entity. Grafana ships a definition and requires
// no credential, so the registrant rule is the only thing standing between a
// request body and a system-owned, endpoint-less Grafana connection.
func TestSeedConnectionsIgnoresConnectionDefinitionOnlyRegistrant(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t), artifactHubRegistrantModelDir)

	registerRegistrant(t, db, "grafana", entity.ConnectionDefinition)

	registrants, err := registrantConnectionsByKind(db)
	if err != nil {
		t.Fatalf("read registrant connections: %v", err)
	}
	if len(registrants["grafana"]) != 0 {
		t.Errorf("a registrant owning only a connection definition must not count as a registrant, got %d row(s)", len(registrants["grafana"]))
	}

	SeedConnections(newSeedTestLogger(t), db, regm)

	grafana := soleConnection(t, connectionsByKind(t, db), "grafana")
	if grafana.ConnectionType != "" || grafana.SubType != "" {
		t.Errorf("grafana was seeded from an externally created registrant: connection_type=%q subType=%q", grafana.ConnectionType, grafana.SubType)
	}
	if got := legacyConnection(t, db, grafana.ID).Type; got != "registry" {
		t.Errorf("grafana was seeded from an externally created registrant: legacy type = %q, want %q", got, "registry")
	}
}

// TestIsSeedableRequiresAnonymousUse pins the credential rule on its own: a
// definition whose credentialSchema marks anything required describes a kind
// that cannot be used anonymously, so the system must not seed it even when it
// does have a registrant.
func TestIsSeedableRequiresAnonymousUse(t *testing.T) {
	registrants := map[string][]connectionv1beta3.Connection{
		"anonymous":  {{Kind: "anonymous"}},
		"authed":     {{Kind: "authed"}},
		"typedreq":   {{Kind: "typedreq"}},
		"noregistry": nil,
	}

	tests := []struct {
		name string
		def  *connectionv1beta3.ConnectionDefinition
		want bool
	}{
		{
			name: "no credential schema at all",
			def:  &connectionv1beta3.ConnectionDefinition{Kind: "anonymous"},
			want: true,
		},
		{
			name: "credential schema with only optional properties",
			def: &connectionv1beta3.ConnectionDefinition{
				Kind:             "anonymous",
				CredentialSchema: map[string]any{"properties": map[string]any{"token": map[string]any{}}},
			},
			want: true,
		},
		{
			name: "empty required list is still anonymous",
			def: &connectionv1beta3.ConnectionDefinition{
				Kind:             "anonymous",
				CredentialSchema: map[string]any{"required": []interface{}{}},
			},
			want: true,
		},
		{
			name: "required credential from JSON round-trip",
			def: &connectionv1beta3.ConnectionDefinition{
				Kind:             "authed",
				CredentialSchema: map[string]any{"required": []interface{}{"kubeconfig"}},
			},
			want: false,
		},
		{
			name: "required credential set in Go",
			def: &connectionv1beta3.ConnectionDefinition{
				Kind:             "typedreq",
				CredentialSchema: map[string]any{"required": []string{"token"}},
			},
			want: false,
		},
		{
			name: "anonymous but Meshery holds no registrant of the kind",
			def:  &connectionv1beta3.ConnectionDefinition{Kind: "noregistry"},
			want: false,
		},
		{
			name: "definition without a kind",
			def:  &connectionv1beta3.ConnectionDefinition{},
			want: false,
		},
		{
			name: "nil definition",
			def:  nil,
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isSeedable(tc.def, registrants); got != tc.want {
				t.Errorf("isSeedable = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestSeedUpdatesForLeavesStatusAlone pins the one field seeding must never
// re-assert: a Connection's status belongs to the connection state machine
// once the Connection exists, so a user who connected or ignored it must not
// have that undone on the next boot.
func TestSeedUpdatesForLeavesStatusAlone(t *testing.T) {
	def := &connectionv1beta3.ConnectionDefinition{
		Kind:           "artifacthub",
		Name:           "Artifact Hub",
		ConnectionType: "source",
		SubType:        "registry",
		Status:         connectionv1beta3.ConnectionStatusRegistered,
	}
	conn := connectionv1beta3.Connection{
		Kind:           "artifacthub",
		Name:           "Artifact Hub",
		ConnectionType: "source",
		SubType:        "registry",
		Status:         connectionv1beta3.ConnectionStatus("connected"),
	}
	legacy := connectionv1beta1.Connection{Kind: "artifacthub", Type: "source"}

	if updates := seedUpdatesFor(conn, legacy, def); len(updates) != 0 {
		t.Errorf("a connected connection matching its definition needs no writes, got %v", updates)
	}
}

// TestSeedUpdatesForWritesBothSchemaHalves covers the split-column invariant: the
// legacy v1beta1 columns and the canonical v1beta3 columns describe the same
// field, so seeding writes them together and treats either half disagreeing as
// stale. Without this, `type` (what GetConnections filters on) and
// `connection_type` (what the API renders) drift apart.
func TestSeedUpdatesForWritesBothSchemaHalves(t *testing.T) {
	def := &connectionv1beta3.ConnectionDefinition{
		Kind: "artifacthub", Name: "Artifact Hub", ConnectionType: "source", SubType: "registry",
	}

	t.Run("canonical column already correct but legacy one stale", func(t *testing.T) {
		conn := connectionv1beta3.Connection{
			Kind: "artifacthub", Name: "Artifact Hub", ConnectionType: "source", SubType: "registry",
		}
		legacy := connectionv1beta1.Connection{Kind: "artifacthub", Type: "registry"}

		updates := seedUpdatesFor(conn, legacy, def)
		if updates["type"] != "source" || updates["connection_type"] != "source" {
			t.Errorf("both halves must be rewritten, got %v", updates)
		}
	})

	t.Run("legacy column already correct but canonical one unset", func(t *testing.T) {
		conn := connectionv1beta3.Connection{
			Kind: "artifacthub", Name: "Artifact Hub", SubType: "registry",
		}
		legacy := connectionv1beta1.Connection{Kind: "artifacthub", Type: "source"}

		updates := seedUpdatesFor(conn, legacy, def)
		if updates["type"] != "source" || updates["connection_type"] != "source" {
			t.Errorf("both halves must be rewritten, got %v", updates)
		}
	})
}

// TestSeedUpdatesForClearsOwnerButKeepsCredential covers the two halves of the
// system-owned invariant, which are not symmetric. Ownership is asserted: no
// path gives a registrant an owner, so writing it can only repair a row. A
// credential is not: both seedable kinds accept an optional one through the
// wizard, so it belongs to the user exactly as `status` does.
func TestSeedUpdatesForClearsOwnerButKeepsCredential(t *testing.T) {
	credentialID := uuid.Must(uuid.NewV4())
	ownerID := uuid.Must(uuid.NewV4())

	def := &connectionv1beta3.ConnectionDefinition{
		Kind: "github", Name: "GitHub", ConnectionType: "source", SubType: "git",
	}

	t.Run("carried on the canonical columns", func(t *testing.T) {
		conn := connectionv1beta3.Connection{
			Kind: "github", Name: "GitHub", ConnectionType: "source", SubType: "git",
			CredentialID: &credentialID,
			Owner:        &ownerID,
		}
		legacy := connectionv1beta1.Connection{Kind: "github", Type: "source"}

		assertOwnershipCleared(t, seedUpdatesFor(conn, legacy, def))
	})

	t.Run("carried on the legacy columns", func(t *testing.T) {
		conn := connectionv1beta3.Connection{
			Kind: "github", Name: "GitHub", ConnectionType: "source", SubType: "git",
		}
		legacy := connectionv1beta1.Connection{
			Kind: "github", Type: "source",
			CredentialID: &credentialID,
			UserID:       &ownerID,
		}

		assertOwnershipCleared(t, seedUpdatesFor(conn, legacy, def))
	})

	t.Run("a credential alone needs no write at all", func(t *testing.T) {
		conn := connectionv1beta3.Connection{
			Kind: "github", Name: "GitHub", ConnectionType: "source", SubType: "git",
			CredentialID: &credentialID,
		}
		legacy := connectionv1beta1.Connection{
			Kind: "github", Type: "source", CredentialID: &credentialID,
		}

		if updates := seedUpdatesFor(conn, legacy, def); len(updates) != 0 {
			t.Errorf("a credentialed connection matching its definition needs no writes, got %v", updates)
		}
	})
}

func assertOwnershipCleared(t *testing.T, updates map[string]interface{}) {
	t.Helper()

	if _, wrote := updates["credential_id"]; wrote {
		t.Errorf("credential_id must never be written, got %v", updates["credential_id"])
	}
	if updates["owner"] != uuid.Nil {
		t.Errorf("owner = %v, want it cleared", updates["owner"])
	}
	if updates["user_id"] != uuid.Nil {
		t.Errorf("user_id = %v, want it cleared", updates["user_id"])
	}
}

// TestSeedConnectionsToleratesMissingDependencies keeps boot resilient: seeding
// is a best-effort step and must never panic the server when it is called
// before the registry or database exist.
func TestSeedConnectionsToleratesMissingDependencies(t *testing.T) {
	log := newSeedTestLogger(t)
	SeedConnections(log, nil, nil)

	db, regm := seedTestRegistry(t, mesheryCoreModelDir(t))
	SeedConnections(log, db, nil)
	SeedConnections(log, nil, regm)
}
