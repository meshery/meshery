package models

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// TestRestrictedAccessJSONTagCanonicalSpelling pins the wire-side spelling
// of RestrictedAccess.IsMesheryUIRestricted to the canonical post-flip
// `isMesheryUIRestricted` (uppercase UI). Meshery Cloud emits this exact
// key on the /{version}/capabilities response; the legacy lowercase-Ui tag
// silently dropped that value and left the playground/restricted-UI gate
// stuck at false, which presented as a cloud.meshery.io ⇄ localhost:9081
// redirect loop on enforced-provider hosts. Guarding the tag here means a
// future tag refactor can't reintroduce that class of bug without flipping
// this assertion as well.
func TestRestrictedAccessJSONTagCanonicalSpelling(t *testing.T) {
	const payload = `{"isMesheryUIRestricted": true, "allowedComponents": {}}`

	var ra RestrictedAccess
	if err := json.Unmarshal([]byte(payload), &ra); err != nil {
		t.Fatalf("unmarshal canonical payload: %v", err)
	}
	if !ra.IsMesheryUIRestricted {
		t.Fatalf("IsMesheryUIRestricted = false; want true — JSON tag must be `isMesheryUIRestricted` (uppercase UI) to match the cloud emitter")
	}

	out, err := json.Marshal(ra)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(out, &raw); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}
	if _, ok := raw["isMesheryUIRestricted"]; !ok {
		t.Fatalf("server output missing canonical key `isMesheryUIRestricted`; got keys: %v", keys(raw))
	}
	if _, ok := raw["isMesheryUiRestricted"]; ok {
		t.Fatalf("server output emitted legacy key `isMesheryUiRestricted`; tag must be canonical only")
	}
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// TestResolveProviderKeyWithProbe_ResolvesCanonicalRemoteName protects the
// enforced-provider boot path: remotes register under their URL host first,
// then publish their canonical providerName from /capabilities. A configured
// PROVIDER=Meshery must still resolve to that remote before the chooser can
// render.
func TestResolveProviderKeyWithProbe_ResolvesCanonicalRemoteName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"providerName":"Meshery","providerType":"remote","capabilities":[]}`))
	}))
	defer server.Close()

	remote := newTestRemoteProvider(t, server.URL)
	remote.Initialize()

	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}

	providers := map[string]Provider{parsed.Host: remote}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	key, ok := ResolveProviderKeyWithProbe(ctx, "Meshery", providers)
	if !ok {
		t.Fatal("ResolveProviderKeyWithProbe did not resolve canonical remote provider name")
	}
	if key != parsed.Host {
		t.Fatalf("resolved key = %q, want %q", key, parsed.Host)
	}
	// Succeeds because VerifyAvailability updated the provider's stored
	// ProviderProperties as a side effect of the probe above, so the plain
	// (probe-less) resolver can now map "Meshery" to the remote's host key.
	if resolvedKey, ok := ResolveProviderKey("Meshery", providers); !ok || resolvedKey != parsed.Host {
		t.Fatalf("ResolveProviderKey after probe = (%q, %v), want (%q, true)", resolvedKey, ok, parsed.Host)
	}
}

func TestRestrictToEnforcedProvider(t *testing.T) {
	local := &DefaultLocalProvider{}
	local.Initialize()
	remote := &RemoteProvider{RemoteProviderURL: "https://cloud.meshery.io"}
	remote.Initialize()
	remote.SetProviderProperties(ProviderProperties{
		ProviderName: "Meshery",
		ProviderType: RemoteProviderType,
		ProviderURL:  remote.RemoteProviderURL,
	})

	t.Run("empty key is a no-op", func(t *testing.T) {
		provs := map[string]Provider{LocalProviderName: local, "Meshery": remote}
		RestrictToEnforcedProvider(provs, "")
		if len(provs) != 2 {
			t.Fatalf("len = %d, want 2", len(provs))
		}
	})

	t.Run("unknown key is a no-op", func(t *testing.T) {
		provs := map[string]Provider{LocalProviderName: local, "Meshery": remote}
		RestrictToEnforcedProvider(provs, "DoesNotExist")
		if len(provs) != 2 {
			t.Fatalf("len = %d, want 2", len(provs))
		}
	})

	t.Run("remote pin drops Local and other remotes", func(t *testing.T) {
		other := &RemoteProvider{RemoteProviderURL: "https://cloud.layer5.io"}
		other.Initialize()
		provs := map[string]Provider{
			LocalProviderName: local,
			"Meshery":         remote,
			"Layer5":          other,
		}
		RestrictToEnforcedProvider(provs, "Meshery")
		if len(provs) != 1 {
			t.Fatalf("len = %d, want 1; keys %v", len(provs), providerMapKeys(provs))
		}
		if _, ok := provs["Meshery"]; !ok {
			t.Fatal("expected Meshery to remain")
		}
		if _, ok := provs[LocalProviderName]; ok {
			t.Fatal("Local must not remain when a remote is pinned")
		}
	})
}

func providerMapKeys(m map[string]Provider) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
