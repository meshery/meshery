package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshsync/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestResolveProviderName covers the precedence rules for picking the
// provider name from a request, plus the enforced-default fallback that
// breaks the /user/login ⇄ /provider redirect loop on enforced-provider
// deployments where the cookie isn't propagating back (SameSite/popup/CDN).
func TestResolveProviderName(t *testing.T) {
	const cookieName = "meshery-provider"

	mkReq := func(setup func(*http.Request)) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
		if setup != nil {
			setup(req)
		}
		return req
	}

	tests := []struct {
		name             string
		setup            func(*http.Request)
		enforcedProvider string
		want             string
	}{
		{
			name:             "no signals, no enforced default → empty",
			setup:            nil,
			enforcedProvider: "",
			want:             "",
		},
		{
			name:             "no signals, enforced default returned",
			setup:            nil,
			enforcedProvider: "Enforced",
			want:             "Enforced",
		},
		{
			name: "enforced default wins over cookie",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "Meshery"})
			},
			enforcedProvider: "Enforced",
			want:             "Enforced",
		},
		{
			name: "enforced default wins over header when cookie empty",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: ""})
				r.Header.Set(cookieName, "Meshery")
			},
			enforcedProvider: "Enforced",
			want:             "Enforced",
		},
		{
			name: "enforced default wins over header and query",
			setup: func(r *http.Request) {
				r.Header.Set(cookieName, "Meshery")
				r.URL.RawQuery = "provider=None"
			},
			enforcedProvider: "Enforced",
			want:             "Enforced",
		},
		{
			name: "enforced default wins over query (None alias does not select Local)",
			setup: func(r *http.Request) {
				r.URL.RawQuery = "provider=None"
			},
			enforcedProvider: "Enforced",
			want:             "Enforced",
		},
		{
			name: "cookie with whitespace value is honored verbatim",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "  Meshery  "})
			},
			enforcedProvider: "",
			want:             "  Meshery  ",
		},
		// Backward-compat guard for the "None" -> "Local" rename. Stale
		// browser cookies, hand-written scripts, and older mesheryctl all
		// keep sending "None" — the request edge must accept it.
		{
			name: "legacy None cookie normalizes to Local",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "None"})
			},
			enforcedProvider: "",
			want:             "Local",
		},
		{
			name:             "legacy None enforced default normalizes to Local",
			setup:            nil,
			enforcedProvider: "None",
			want:             "Local",
		},
		{
			name: "alias matching is case-insensitive",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "none"})
			},
			enforcedProvider: "",
			want:             "Local",
		},
		{
			name: "lowercase canonical name normalizes to canonical Local",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "local"})
			},
			enforcedProvider: "",
			want:             "Local",
		},
		{
			name:             "uppercase canonical name normalizes to canonical Local",
			setup:            nil,
			enforcedProvider: "LOCAL",
			want:             "Local",
		},
		{
			name: "non-local provider name is returned verbatim (remote names not normalized)",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "Meshery"})
			},
			enforcedProvider: "",
			want:             "Meshery",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, _ := resolveProviderName(mkReq(tc.setup), cookieName, tc.enforcedProvider)
			if got != tc.want {
				t.Fatalf("resolveProviderName: want %q, got %q", tc.want, got)
			}
		})
	}
}

// TestResolveProviderName_ReturnsLegacyAliasFlag locks in the second return
// value used by ProviderMiddleware to drive the one-shot deprecation warning.
// The boolean must be true for any casing of the legacy alias and false for
// canonical names (including the canonical local provider name).
func TestResolveProviderName_ReturnsLegacyAliasFlag(t *testing.T) {
	const cookieName = "meshery-provider"
	mkReq := func(cookieValue string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
		if cookieValue != "" {
			req.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})
		}
		return req
	}

	cases := []struct {
		name     string
		cookie   string
		wantFlag bool
	}{
		{"canonical Local does not flag legacy", "Local", false},
		{"lowercase local does not flag legacy", "local", false},
		{"legacy None flags legacy", "None", true},
		{"lowercase none flags legacy", "none", true},
		{"uppercase NONE flags legacy", "NONE", true},
		{"remote Meshery does not flag legacy", "Meshery", false},
		{"empty does not flag legacy", "", false},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			_, got := resolveProviderName(mkReq(c.cookie), cookieName, "")
			if got != c.wantFlag {
				t.Errorf("usedLegacyAlias = %v, want %v", got, c.wantFlag)
			}
		})
	}
}

func TestAuthMiddleWare(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	t.Log("Need to run AuthMiddleware() skipping")
	//_, err := AuthMiddleware(nil, nil)
	//if err != nil {
	//	t.Errorf("AuthMiddleWare() failed with error: %s", err)
	//}
}

// TestWriteMeshkitError_ErrTransientProvider is a focused test for the
// transient-provider response shape emitted by SessionInjectorMiddleware when
// Meshery Cloud is temporarily unreachable. A full middleware-level test would
// require stubbing the 100+ method models.Provider interface, which the
// package doesn't yet have. Instead, we test the migration end-to-end at the
// response-helper level: the same two lines the middleware now executes.
//
// If a provider mocking pattern lands in this package later, promote this to
// a middleware-dispatch test by wiring a stub provider whose GetUserDetails
// returns "Could not reach remote provider".
func TestWriteMeshkitError_ErrTransientProvider(t *testing.T) {
	rec := httptest.NewRecorder()

	writeMeshkitError(
		rec,
		ErrTransientProvider(errors.New("Could not reach remote provider: dial tcp: i/o timeout")),
		http.StatusServiceUnavailable,
	)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type application/json; charset=utf-8, got %q", ct)
	}
	if nosniff := resp.Header.Get("X-Content-Type-Options"); nosniff != "nosniff" {
		t.Errorf("expected X-Content-Type-Options: nosniff, got %q", nosniff)
	}

	var decoded struct {
		Error                string   `json:"error"`
		Code                 string   `json:"code"`
		Severity             string   `json:"severity"`
		ProbableCause        []string `json:"probableCause"`
		SuggestedRemediation []string `json:"suggestedRemediation"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("body did not parse as JSON: %v", err)
	}
	if decoded.Code != ErrTransientProviderCode {
		t.Errorf("expected code %q, got %q", ErrTransientProviderCode, decoded.Code)
	}
	if decoded.Error == "" {
		t.Errorf("expected non-empty error message")
	}
	if len(decoded.SuggestedRemediation) == 0 {
		t.Errorf("expected suggestedRemediation to be populated for transient-provider errors")
	}
}

// TestIsTransientProviderError exercises the detector used by
// SessionInjectorMiddleware to decide between a 401 auth failure and a 503
// transient-provider failure. Regression guard in case the error-message
// substrings the detector relies on drift.
func TestIsTransientProviderError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil is not transient", nil, false},
		{"generic auth failure is not transient", fmt.Errorf("token expired"), false},
		{"Could not reach remote provider is transient", fmt.Errorf("Could not reach remote provider: dial tcp: i/o timeout"), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isTransientProviderError(tc.err); got != tc.want {
				t.Errorf("isTransientProviderError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// blockingDiscoverAction stands in for the kubernetes DiscoverAction. It waits
// until the test releases it, so the request context can be cancelled first,
// then reports the context error it observed.
type blockingDiscoverAction struct {
	release <-chan struct{}
	ctxErr  chan error
}

func (a *blockingDiscoverAction) ExecuteOnEntry(ctx context.Context, _ interface{}, _ interface{}) (machines.EventType, *events.Event, error) {
	<-a.release
	a.ctxErr <- ctx.Err()
	return machines.NoOp, nil, nil
}

func (a *blockingDiscoverAction) Execute(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *blockingDiscoverAction) ExecuteOnExit(context.Context, interface{}, interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

// The discovery K8sFSMMiddleware starts runs in a goroutine that outlives the
// request, so it must not inherit the request's cancellation.
func TestK8sFSMMiddleware_DetachesDiscoveryFromRequestContext(t *testing.T) {
	release := make(chan struct{})
	action := &blockingDiscoverAction{release: release, ctxErr: make(chan error, 1)}

	systemID := uuid.Must(uuid.NewV4())
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	connectionID := uuid.Must(uuid.NewV4())

	inst := &machines.StateMachine{
		ID:           connectionID,
		Name:         "kubernetes",
		InitialState: machines.InitialState,
		CurrentState: machines.InitialState,
		States: machines.States{
			machines.InitialState: {Events: machines.Events{machines.Discovery: machines.DISCOVERED}},
			machines.DISCOVERED:   {Action: action},
		},
		Context: &kubernetes.MachineCtx{MesheryCtrlsHelper: &models.MesheryControllersHelper{}},
		Log:     newTestLogger(t),
	}

	h := &Handler{
		config:                                  &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:                                     newTestLogger(t),
		SystemID:                                &systemID,
		ConnectionToStateMachineInstanceTracker: trackerWith(connectionID, inst),
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}
	if err := db.AutoMigrate(&model.KubernetesResource{}, &model.KubernetesResourceObjectMeta{}); err != nil {
		t.Fatalf("failed to migrate kubernetes resource tables: %v", err)
	}
	provider := &models.DefaultLocalProvider{}
	provider.Initialize()
	provider.GenericPersister = &database.Handler{DB: db}

	reqCtx, cancel := context.WithCancel(context.Background())
	reqCtx = context.WithValue(reqCtx, models.TokenCtxKey, "test-token")
	reqCtx = context.WithValue(reqCtx, models.UserCtxKey, user)
	reqCtx = context.WithValue(reqCtx, models.SystemIDKey, &systemID)
	reqCtx = context.WithValue(reqCtx, models.AllKubeClusterKey, []*models.K8sContext{{ConnectionID: connectionID.String()}})

	K8sFSMMiddleware(reqCtx, h, provider, user)
	cancel()
	close(release)

	select {
	case err := <-action.ctxErr:
		if err != nil {
			t.Fatalf("discover action saw a cancelled context: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for the discover action to run")
	}
}
