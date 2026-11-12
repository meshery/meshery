package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/models/registration"
	"github.com/meshery/schemas/models/v1beta1/environment"
)

// artifactHubRegistrantModelDir is a minimal in-repo model published under the
// `artifacthub` registrant. The shipped Artifact Hub models it would otherwise
// come from live under models/, which the recommended sparse checkout excludes,
// so committing a fixture keeps the Artifact Hub seeding path covered without
// pulling in the model registry.
const artifactHubRegistrantModelDir = "testdata/artifacthub-registrant-model"

// seedModelDirs resolves the model directories that boot-time registration reads
// for this fixture. All are available in a sparse clone (see the "Cloning the
// Repository" contributing guide), so the default `go test ./...` keeps passing
// there:
//
//   - meshery-core carries the Connection *definitions* (Artifact Hub, GitHub,
//     Kubernetes, ...) and registers under the `meshery` registrant.
//   - kubernetes registers under the `github` registrant, which is what makes the
//     GitHub connection definition seedable.
//   - the committed artifacthub-registrant-model registers under the
//     `artifacthub` registrant, which makes the Artifact Hub definition seedable.
//
// Together they create the `artifacthub` and `github` registrant Connections
// that SeedConnections is scoped to, so this fixture is the same input the
// server gets on boot rather than a synthetic one.
func seedModelDirs(t *testing.T) []string {
	t.Helper()
	return []string{
		resolveModelDir(t, "meshery-core"),
		resolveModelDir(t, "kubernetes"),
		artifactHubRegistrantModelDir,
	}
}

// resolveModelDir returns a shipped model's registration directory
// (models/<name>/<version>/v1.0.0) without pinning a specific version, so the
// test keeps working as the model registry is re-synced. It selects the highest
// version directory that contains a model.json.
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

// apiConnection is the connection as it reaches the wire (and therefore the
// Connections page), decoded from the handler's own response rather than read
// back out of the database.
type apiConnection struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Type    string `json:"type"`
	SubType string `json:"subType"`
	Status  string `json:"status"`
	Owner   string `json:"owner"`
}

type apiConnectionPage struct {
	Connections []apiConnection `json:"connections"`
	TotalCount  int             `json:"totalCount"`
}

// newSeedConnectionsAPIFixture stands up the connections API over a database
// that has been through boot-time model registration, exactly as the server
// does: registry manager first (it migrates the legacy v1beta1 `connections`
// columns), then the canonical v1beta3 Connection on top.
func newSeedConnectionsAPIFixture(t *testing.T) (*Handler, models.Provider, *database.Handler, *meshmodel.RegistryManager) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	regm, err := meshmodel.NewRegistryManager(&db)
	if err != nil {
		t.Fatalf("new registry manager: %v", err)
	}
	if err := db.AutoMigrate(
		connections.Connection{},
		environment.Environment{},
		environment.EnvironmentConnectionMapping{},
	); err != nil {
		t.Fatalf("migrate connection tables: %v", err)
	}

	regHelper := registration.NewRegistrationHelper(t.TempDir(), regm, models.NewRegistrationFailureLogHandler())
	for _, dir := range seedModelDirs(t) {
		regHelper.Register(registration.NewDir(dir))
	}

	h := &Handler{config: &models.HandlerConfig{}, log: newTestLogger(t)}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister: &models.ConnectionPersister{DB: &db},
	}
	return h, provider, &db, regm
}

// getConnectionsAPI issues a real request through the same handler the
// `/api/integrations/connections` route is bound to and returns the decoded
// response together with its raw body, so a test can both assert on it and log
// what a client actually received.
func getConnectionsAPI(t *testing.T, h *Handler, provider models.Provider, query string) (apiConnectionPage, string) {
	t.Helper()

	url := "/api/integrations/connections?pageSize=100"
	if query != "" {
		url += "&" + query
	}
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	h.GetConnections(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s: status = %d, body: %s", url, rec.Code, rec.Body.String())
	}

	var page apiConnectionPage
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode response: %v; body: %s", err, rec.Body.String())
	}
	sort.Slice(page.Connections, func(i, j int) bool { return page.Connections[i].Kind < page.Connections[j].Kind })
	return page, rec.Body.String()
}

func summarize(page apiConnectionPage) string {
	out, _ := json.MarshalIndent(page.Connections, "", "  ")
	return string(out)
}

func connectionOfKind(page apiConnectionPage, kind string) (apiConnection, bool) {
	for _, conn := range page.Connections {
		if conn.Kind == kind {
			return conn, true
		}
	}
	return apiConnection{}, false
}

// TestConnectionsAPISurfacesSeededSystemConnections is the user-visible end of
// meshery/meshery#20932: after boot the Connections page must already list
// Artifact Hub and GitHub, and filtering it by the `source` connection type must
// return both.
//
// The assertions run against the handler's own HTTP response - the exact bytes
// the Connections page consumes - rather than against the database, so a write
// that lands in the database but never reaches the wire fails here.
func TestConnectionsAPISurfacesSeededSystemConnections(t *testing.T) {
	h, provider, db, regm := newSeedConnectionsAPIFixture(t)

	// Before seeding: registration has already created a connection row for
	// both kinds, but with the stale identity the model.json `registrant` blobs
	// carry - so filtering by `source` finds neither, which is the reported
	// defect.
	before, _ := getConnectionsAPI(t, h, provider, "")
	t.Logf("BEFORE seeding - GET /api/integrations/connections\n%s", summarize(before))

	beforeSource, _ := getConnectionsAPI(t, h, provider, "filter=type+source")
	t.Logf("BEFORE seeding - GET /api/integrations/connections?filter=type+source\n%s", summarize(beforeSource))
	if len(beforeSource.Connections) != 0 {
		t.Fatalf("precondition: expected no `source` connections before seeding, got %d", len(beforeSource.Connections))
	}
	for _, kind := range []string{"artifacthub", "github"} {
		conn, ok := connectionOfKind(before, kind)
		if !ok {
			t.Fatalf("precondition: expected registration to create a %q connection row", kind)
		}
		if conn.Type == "source" {
			t.Fatalf("precondition: %q was already typed source before seeding", kind)
		}
	}

	models.SeedConnections(newTestLogger(t), db, regm)

	after, rawAfter := getConnectionsAPI(t, h, provider, "")
	t.Logf("AFTER seeding - GET /api/integrations/connections\n%s", summarize(after))

	sourceOnly, rawSource := getConnectionsAPI(t, h, provider, "filter=type+source")
	t.Logf("AFTER seeding - GET /api/integrations/connections?filter=type+source\n%s", summarize(sourceOnly))
	t.Logf("AFTER seeding - raw response body (unfiltered): %s", rawAfter)
	t.Logf("AFTER seeding - raw response body (type=source): %s", rawSource)

	want := map[string]apiConnection{
		"artifacthub": {Kind: "artifacthub", Name: "Artifact Hub", Type: "source", SubType: "registry"},
		"github":      {Kind: "github", Name: "GitHub", Type: "source", SubType: "git"},
	}

	if len(sourceOnly.Connections) != len(want) {
		t.Fatalf("filtering by type=source returned %d connection(s), want %d: %s",
			len(sourceOnly.Connections), len(want), summarize(sourceOnly))
	}
	for _, conn := range sourceOnly.Connections {
		expected, ok := want[conn.Kind]
		if !ok {
			t.Errorf("unexpected kind %q in the source-filtered connections", conn.Kind)
			continue
		}
		if conn.Name != expected.Name || conn.Type != expected.Type || conn.SubType != expected.SubType {
			t.Errorf("kind %q: got name=%q type=%q subType=%q, want name=%q type=%q subType=%q",
				conn.Kind, conn.Name, conn.Type, conn.SubType, expected.Name, expected.Type, expected.SubType)
		}
		// System-owned: no owner reaches the wire.
		if conn.Owner != "" && conn.Owner != uuid.Nil.String() {
			t.Errorf("kind %q: seeded connection reached the wire with owner %q", conn.Kind, conn.Owner)
		}
	}

	// The sub-type filter is the other half of the same raw-SQL filter, and is
	// what tells the two seeded source connections apart on the page.
	for kind, expected := range want {
		filtered, _ := getConnectionsAPI(t, h, provider, "filter=sub_type+"+expected.SubType)
		conn, ok := connectionOfKind(filtered, kind)
		if !ok {
			t.Errorf("filtering by sub_type=%s did not return the %q connection: %s",
				expected.SubType, kind, summarize(filtered))
			continue
		}
		if conn.Name != expected.Name {
			t.Errorf("kind %q: name = %q, want %q", kind, conn.Name, expected.Name)
		}
	}
}

// TestConnectionsAPISeededConnectionsSurviveRestart pins the restart outcome at
// the API surface: a second boot re-registers every model and seeds again, and
// the Connections page must still show exactly one Artifact Hub and one GitHub
// connection - not a duplicate pair.
func TestConnectionsAPISeededConnectionsSurviveRestart(t *testing.T) {
	h, provider, db, regm := newSeedConnectionsAPIFixture(t)
	log := newTestLogger(t)

	models.SeedConnections(log, db, regm)
	first, _ := getConnectionsAPI(t, h, provider, "filter=type+source")

	// A restart: registration runs again, then seeding, exactly as
	// SeedComponents sequences them.
	regHelper := registration.NewRegistrationHelper(t.TempDir(), regm, models.NewRegistrationFailureLogHandler())
	for _, dir := range seedModelDirs(t) {
		regHelper.Register(registration.NewDir(dir))
	}
	models.SeedConnections(log, db, regm)

	second, _ := getConnectionsAPI(t, h, provider, "filter=type+source")
	t.Logf("AFTER restart - GET /api/integrations/connections?filter=type+source\n%s", summarize(second))

	if len(second.Connections) != len(first.Connections) {
		t.Fatalf("restart changed the source-filtered connection count: %d -> %d\n%s",
			len(first.Connections), len(second.Connections), summarize(second))
	}
	for i := range first.Connections {
		if second.Connections[i] != first.Connections[i] {
			t.Errorf("restart changed the %q connection: %+v -> %+v",
				first.Connections[i].Kind, first.Connections[i], second.Connections[i])
		}
	}
}
