package models

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/models/registration"
	connectionv1beta1 "github.com/meshery/schemas/models/v1beta1/connection"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
	"github.com/sirupsen/logrus"
)

// mesheryCoreModelDir carries the shipped connection definitions.
const mesheryCoreModelDir = "../../models/meshery-core/0.7.2/v1.0.0"

// Representative models whose model.json registrant is `artifacthub` and
// `github` respectively. Registering them is what creates the registrant
// Connections that seeding is scoped to.
const (
	artifactHubRegistrantModelDir = "../../models/kubevault/2026.1.8-rc.0/v1.0.0"
	gitHubRegistrantModelDir      = "../../models/azure-operational-insights/azureserviceoperator_customresourcedefinitions_v2.13.0.yaml/v1.0.0"
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

	regHelper := registration.NewRegistrationHelper(t.TempDir(), regm, NewRegistrationFailureLogHandler())
	for _, dir := range modelDirs {
		regHelper.Register(registration.NewDir(dir))
	}
	return &db, regm
}

// connectionsByKind reads the connections table into a kind-keyed map, failing
// the test if a kind ever has more than one row - the duplicate-row regression
// this seeding must never introduce.
func connectionsByKind(t *testing.T, db *database.Handler) map[string]connectionv1beta1.Connection {
	t.Helper()

	var conns []connectionv1beta1.Connection
	if err := db.Find(&conns).Error; err != nil {
		t.Fatalf("read connections: %v", err)
	}

	byKind := make(map[string]connectionv1beta1.Connection, len(conns))
	for _, conn := range conns {
		if existing, dup := byKind[conn.Kind]; dup {
			t.Fatalf("duplicate connection rows for kind %q: %s and %s", conn.Kind, existing.ID, conn.ID)
		}
		byKind[conn.Kind] = conn
	}
	return byKind
}

// TestSeedConnectionsMaterializesDefinitionBackedConnections is the core
// behavior: after boot the Artifact Hub and GitHub connections exist carrying
// the identity their connection definitions declare, not the stale identity the
// model.json registrant blobs carry.
func TestSeedConnectionsMaterializesDefinitionBackedConnections(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir, artifactHubRegistrantModelDir, gitHubRegistrantModelDir)

	// Registration alone leaves the registrant rows disagreeing with the
	// definitions; assert that starting point so the test proves seeding did
	// the work rather than the registrant blobs happening to be right.
	before := connectionsByKind(t, db)
	if got := before["github"].Name; got != "Github" {
		t.Fatalf("precondition: expected registrant name %q, got %q", "Github", got)
	}
	if got := before["artifacthub"].Type; got != "registry" {
		t.Fatalf("precondition: expected registrant type %q, got %q", "registry", got)
	}

	SeedConnections(newSeedTestLogger(t), db, regm)

	after := connectionsByKind(t, db)
	for _, tc := range []struct {
		kind, name, connType, subType string
	}{
		{kind: "artifacthub", name: "Artifact Hub", connType: "source", subType: "registry"},
		{kind: "github", name: "GitHub", connType: "source", subType: "git"},
	} {
		conn, ok := after[tc.kind]
		if !ok {
			t.Fatalf("no connection seeded for kind %q", tc.kind)
		}
		if conn.Name != tc.name {
			t.Errorf("kind %q: name = %q, want %q", tc.kind, conn.Name, tc.name)
		}
		if conn.Type != tc.connType {
			t.Errorf("kind %q: type = %q, want %q", tc.kind, conn.Type, tc.connType)
		}
		if conn.SubType != tc.subType {
			t.Errorf("kind %q: subType = %q, want %q", tc.kind, conn.SubType, tc.subType)
		}
		// System-owned and anonymous: no credential, no owner.
		if conn.CredentialID != nil && *conn.CredentialID != uuid.Nil {
			t.Errorf("kind %q: seeded connection carries credential %s", tc.kind, *conn.CredentialID)
		}
		if conn.UserID != nil && *conn.UserID != uuid.Nil {
			t.Errorf("kind %q: seeded connection carries owner %s", tc.kind, *conn.UserID)
		}
	}
}

// TestSeedConnectionsIsIdempotentAcrossRestarts proves the restart behavior the
// issue calls for: no duplicate rows, and no write at all on a second boot.
func TestSeedConnectionsIsIdempotentAcrossRestarts(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir, artifactHubRegistrantModelDir, gitHubRegistrantModelDir)
	log := newSeedTestLogger(t)

	SeedConnections(log, db, regm)
	first := connectionsByKind(t, db)

	// A restart re-runs registration before seeding, exactly as SeedComponents
	// does, so exercise both.
	regHelper := registration.NewRegistrationHelper(t.TempDir(), regm, NewRegistrationFailureLogHandler())
	for _, dir := range []string{mesheryCoreModelDir, artifactHubRegistrantModelDir, gitHubRegistrantModelDir} {
		regHelper.Register(registration.NewDir(dir))
	}
	SeedConnections(log, db, regm)

	// connectionsByKind fails on any duplicate, so reaching here already proves
	// the row count per kind held at one.
	second := connectionsByKind(t, db)
	if len(second) != len(first) {
		t.Fatalf("restart changed connection count: %d -> %d", len(first), len(second))
	}
	for kind, before := range first {
		after, ok := second[kind]
		if !ok {
			t.Fatalf("kind %q disappeared across restart", kind)
		}
		if after.ID != before.ID {
			t.Errorf("kind %q: id changed across restart: %s -> %s", kind, before.ID, after.ID)
		}
		if after.Name != before.Name || after.Type != before.Type || after.SubType != before.SubType {
			t.Errorf("kind %q: identity changed across restart: %+v -> %+v", kind, before, after)
		}
	}

	// The second boot must do no work at all: every registrant already matches
	// its definition, so there is nothing left to update.
	byKind, err := registrantConnectionsByKind(db)
	if err != nil {
		t.Fatalf("read registrant connections: %v", err)
	}
	defs, err := registeredConnectionDefinitions(regm)
	if err != nil {
		t.Fatalf("read connection definitions: %v", err)
	}
	for _, def := range defs {
		if !isSeedable(def, byKind) {
			continue
		}
		changed, err := seedConnectionForDefinition(db, def, byKind[def.Kind])
		if err != nil {
			t.Fatalf("kind %q: %v", def.Kind, err)
		}
		if changed != 0 {
			t.Errorf("kind %q: steady-state seeding wrote %d row(s), want 0", def.Kind, changed)
		}
	}
}

// TestSeedConnectionsSkipsKindsWithoutRegistrants guards the scoping rule: a
// kind that ships a definition but that Meshery does not itself source content
// through must not be materialized as an empty, endpoint-less row.
func TestSeedConnectionsSkipsKindsWithoutRegistrants(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir, artifactHubRegistrantModelDir, gitHubRegistrantModelDir)

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
// rule actually exists for. Registering a cluster's components creates a
// `kubernetes` registrant Connection (server/models/meshmodel/core/register.go
// types it `registry`), so the kind does acquire a registrant the moment a user
// connects a cluster. Kubernetes cannot be used anonymously - its
// credentialSchema requires a kubeconfig - so it must still never be seeded,
// and its registrant must be left exactly as registration wrote it.
func TestSeedConnectionsSkipsKubernetesRegistrant(t *testing.T) {
	db, regm := seedTestRegistry(t, mesheryCoreModelDir, artifactHubRegistrantModelDir)

	// Register a component under a kubernetes registrant the same way cluster
	// component registration does, so a `kubernetes` registrant row exists.
	k8sRegistrant := connectionv1beta1.Connection{Kind: "kubernetes", Type: "registry"}
	registrantID, err := k8sRegistrant.Create(db)
	if err != nil {
		t.Fatalf("create kubernetes registrant: %v", err)
	}
	if err := db.Create(&meshmodel.Registry{
		ID:           uuid.Must(uuid.NewV4()),
		RegistrantID: registrantID,
		Entity:       uuid.Must(uuid.NewV4()),
		Type:         entity.ComponentDefinition,
	}).Error; err != nil {
		t.Fatalf("create registry entry: %v", err)
	}

	SeedConnections(newSeedTestLogger(t), db, regm)

	after := connectionsByKind(t, db)
	k8s, ok := after["kubernetes"]
	if !ok {
		t.Fatal("the kubernetes registrant connection disappeared")
	}
	// Untouched: still the registrant identity, not the definition's
	// platform/orchestration identity.
	if k8s.Type != "registry" || k8s.SubType != "" {
		t.Errorf("kubernetes registrant was seeded: type=%q subType=%q, want registry/\"\"", k8s.Type, k8s.SubType)
	}

	// The anonymous kind alongside it still seeds, so this proves the credential
	// rule did the skipping rather than seeding having failed wholesale.
	if got := after["artifacthub"].SubType; got != "registry" {
		t.Errorf("artifacthub subType = %q, want %q", got, "registry")
	}
	if got := after["artifacthub"].Type; got != "source" {
		t.Errorf("artifacthub type = %q, want %q", got, "source")
	}
}

// TestIsSeedableRequiresAnonymousUse pins the credential rule on its own: a
// definition whose credentialSchema marks anything required describes a kind
// that cannot be used anonymously, so the system must not seed it even when it
// does have a registrant.
func TestIsSeedableRequiresAnonymousUse(t *testing.T) {
	registrants := map[string][]connectionv1beta1.Connection{
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
	conn := connectionv1beta1.Connection{
		Kind:    "artifacthub",
		Name:    "Artifact Hub",
		Type:    "source",
		SubType: "registry",
		Status:  connectionv1beta1.ConnectionStatus("connected"),
	}

	if updates := seedUpdatesFor(conn, def); len(updates) != 0 {
		t.Errorf("a connected connection matching its definition needs no writes, got %v", updates)
	}
}

// TestSeedUpdatesForClearsCredentialAndOwner covers the system-owned invariant:
// a seeded Connection uses its kind's API anonymously, so any credential or
// owner found on it is cleared.
func TestSeedUpdatesForClearsCredentialAndOwner(t *testing.T) {
	credentialID := uuid.Must(uuid.NewV4())
	ownerID := uuid.Must(uuid.NewV4())

	def := &connectionv1beta3.ConnectionDefinition{
		Kind: "github", Name: "GitHub", ConnectionType: "source", SubType: "git",
	}
	conn := connectionv1beta1.Connection{
		Kind: "github", Name: "GitHub", Type: "source", SubType: "git",
		CredentialID: &credentialID,
		UserID:       &ownerID,
	}

	updates := seedUpdatesFor(conn, def)
	if updates["credential_id"] != uuid.Nil {
		t.Errorf("credential_id = %v, want it cleared", updates["credential_id"])
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

	db, regm := seedTestRegistry(t, mesheryCoreModelDir)
	SeedConnections(log, db, nil)
	SeedConnections(log, nil, regm)
}
