package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
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

// TestK8sFSMMiddleware_NilContextValues verifies that K8sFSMMiddleware defensively
// handles missing context values, type mismatches on AllKubeClusterKey, nil users,
// and nil providers without triggering a runtime panic.
func TestK8sFSMMiddleware_NilContextValues(t *testing.T) {
	h := &Handler{}

	t.Run("missing AllKubeClusterKey does not panic", func(t *testing.T) {
		ctx := context.Background()
		user := &models.User{ID: uuid.Must(uuid.NewV4())}
		// Must not panic on missing AllKubeClusterKey
		K8sFSMMiddleware(ctx, h, nil, user)
	})

	t.Run("invalid type for AllKubeClusterKey does not panic", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.AllKubeClusterKey, "invalid-type")
		user := &models.User{ID: uuid.Must(uuid.NewV4())}
		// Must not panic on interface conversion error
		K8sFSMMiddleware(ctx, h, nil, user)
	})

	t.Run("nil user does not panic", func(t *testing.T) {
		ctx := context.Background()
		// Must not panic on user.ID dereference
		K8sFSMMiddleware(ctx, h, nil, nil)
	})

	t.Run("nil Handler does not panic", func(t *testing.T) {
		ctx := context.Background()
		// Must return immediately without nil-dereferencing h
		K8sFSMMiddleware(ctx, nil, nil, nil)
	})

	t.Run("HTTP middleware wrapper executes next handler", func(t *testing.T) {
		called := false
		next := func(w http.ResponseWriter, req *http.Request, pref *models.Preference, user *models.User, prov models.Provider) {
			called = true
		}
		mw := h.K8sFSMMiddleware(next)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/system/sync", nil)

		mw(rec, req, nil, nil, nil)
		if !called {
			t.Fatalf("expected next handler to be invoked by K8sFSMMiddleware")
		}
	})
}

// TestPersistAndPublishEvent_DefensiveGuards verifies that persistAndPublishEvent
// safely ignores nil event or provider pointers and does not panic on nil loggers.
func TestPersistAndPublishEvent_DefensiveGuards(t *testing.T) {
	dummyUUID := core.Uuid(uuid.Must(uuid.NewV4()))

	t.Run("nil provider returns safely without panic", func(t *testing.T) {
		event := events.NewEvent().ActedUpon(uuid.UUID(dummyUUID)).WithSeverity(events.Informational).Build()
		// Must return early without attempting to call provider.PersistEvent
		persistAndPublishEvent(nil, nil, "dummy-token", nil, dummyUUID, event)
	})

	t.Run("nil event returns safely without panic", func(t *testing.T) {
		// Must return early without attempting to dereference nil event
		persistAndPublishEvent(nil, nil, "dummy-token", nil, dummyUUID, nil)
	})

	t.Run("nil event with non-nil provider returns safely without calling PersistEvent", func(t *testing.T) {
		prov := &recordingEventsProvider{}
		// Must return early without attempting to call provider.PersistEvent
		persistAndPublishEvent(nil, prov, "dummy-token", nil, dummyUUID, nil)
		if len(prov.persisted) != 0 {
			t.Fatalf("expected no events to be persisted when event is nil, got %d", len(prov.persisted))
		}
	})
}

// TestKubernetesMiddleware_NilGuards verifies that KubernetesMiddleware returns
// descriptive errors when required arguments (user, provider, token) are missing.
func TestKubernetesMiddleware_NilGuards(t *testing.T) {
	h := &Handler{}
	ctx := context.Background()

	t.Run("nil user returns descriptive error", func(t *testing.T) {
		resCtx, err := KubernetesMiddleware(ctx, h, nil, nil, nil)
		if err == nil {
			t.Fatalf("expected error when user is nil, got nil")
		}
		if resCtx != nil {
			t.Fatalf("expected nil context when error is returned, got %v", resCtx)
		}
	})

	t.Run("nil provider returns descriptive error", func(t *testing.T) {
		user := &models.User{ID: uuid.Must(uuid.NewV4())}
		resCtx, err := KubernetesMiddleware(ctx, h, nil, user, nil)
		if err == nil {
			t.Fatalf("expected error when provider is nil, got nil")
		}
		if resCtx != nil {
			t.Fatalf("expected nil context when error is returned, got %v", resCtx)
		}
	})

	t.Run("missing token returns ErrRetrieveUserToken", func(t *testing.T) {
		user := &models.User{ID: uuid.Must(uuid.NewV4())}
		prov := &recordingEventsProvider{}
		resCtx, err := KubernetesMiddleware(ctx, h, prov, user, nil)
		if err == nil {
			t.Fatalf("expected error when token is missing from context, got nil")
		}
		if resCtx != nil {
			t.Fatalf("expected nil context when error is returned, got %v", resCtx)
		}
		if code := meshkiterrors.GetCode(err); code != ErrRetrieveUserTokenCode {
			t.Fatalf("expected error code %s, got %s", ErrRetrieveUserTokenCode, code)
		}
	})
}

// TestRegistrationAuthMiddleware pins the auth boundary for component
// registration: selector-free posts on a PROVIDER-unset deployment must reach
// the handler (adapters register with no provider selector and local mode has
// no session to validate), while a remote selector or an enforced provider
// must go through ProviderAuth instead of reaching the handler.
func TestRegistrationAuthMiddleware(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()
	providers := map[string]models.Provider{local.Name(): local}

	tests := []struct {
		name             string
		providerOverride string
		setup            func(*http.Request)
		wantReached      bool
	}{
		{
			name:        "no selector, no enforced provider reaches handler",
			wantReached: true,
		},
		{
			name: "Local cookie reaches handler",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: "meshery-provider", Value: "Local"})
			},
			wantReached: true,
		},
		{
			name: "legacy None cookie reaches handler",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: "meshery-provider", Value: "None"})
			},
			wantReached: true,
		},
		{
			name: "remote selector requires auth and does not reach handler",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: "meshery-provider", Value: "Meshery"})
			},
			wantReached: false,
		},
		{
			name: "remote provider query param requires auth and does not reach handler",
			setup: func(r *http.Request) {
				r.URL.RawQuery = "provider=Meshery"
			},
			wantReached: false,
		},
		{
			name:             "enforced provider requires auth and does not reach handler",
			providerOverride: "Meshery",
			wantReached:      false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := newTestHandler(t, providers, tc.providerOverride)
			reached := false
			inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reached = true
				w.WriteHeader(http.StatusOK)
			})
			chain := h.ProviderMiddleware(h.RegistrationAuthMiddleware(inner))
			req := httptest.NewRequest(http.MethodPost, "/api/meshmodel/components/register", nil)
			if tc.setup != nil {
				tc.setup(req)
			}
			rec := httptest.NewRecorder()
			chain.ServeHTTP(rec, req)
			if reached != tc.wantReached {
				t.Fatalf("reached = %v, want %v (status %d, location %q)", reached, tc.wantReached, rec.Result().StatusCode, rec.Header().Get("Location"))
			}
			if tc.wantReached && rec.Result().StatusCode != http.StatusOK {
				t.Fatalf("expected 200 when handler is reached, got %d", rec.Result().StatusCode)
			}
		})
	}
}
