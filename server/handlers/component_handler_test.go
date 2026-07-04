package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/core"
)

// ---- spy provider for credential resolution tests ----

type ociCredentialSpyProvider struct {
	*models.DefaultLocalProvider
	getCredentialFn func(token string, credentialID core.Uuid) (*models.Credential, int, error)
}

func (m *ociCredentialSpyProvider) GetCredentialByID(token string, credentialID core.Uuid) (*models.Credential, int, error) {
	if m.getCredentialFn != nil {
		return m.getCredentialFn(token, credentialID)
	}
	return m.DefaultLocalProvider.GetCredentialByID(token, credentialID)
}

// ---- helpers ----

func withTokenCtx(t *testing.T) context.Context {
	t.Helper()
	return context.WithValue(context.Background(), models.TokenCtxKey, "test-token")
}

func jsonBody(t *testing.T, v interface{}) *bytes.Buffer {
	t.Helper()
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(v); err != nil {
		t.Fatal(err)
	}
	return &b
}

// writeOCIStore creates a minimal OCI layout at dir containing the given
// layer files, tagged with tag.
func writeOCIStore(t *testing.T, dir, tag string, layerFiles map[string]string) {
	t.Helper()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, content := range layerFiles {
		if err := tw.WriteHeader(&tar.Header{Name: name, Size: int64(len(content)), Mode: 0644, Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(tw, content); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	layerData := buf.Bytes()
	h := sha256.Sum256(layerData)
	layerDigest := fmt.Sprintf("sha256:%s", hex.EncodeToString(h[:]))

	// manifest
	m := struct {
		Layers []struct {
			Digest string `json:"digest"`
		} `json:"layers"`
	}{Layers: []struct {
		Digest string `json:"digest"`
	}{{Digest: layerDigest}}}
	md, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	hm := sha256.Sum256(md)
	manifestDigest := fmt.Sprintf("sha256:%s", hex.EncodeToString(hm[:]))

	// write blobs
	blobDir := filepath.Join(dir, "blobs", "sha256")
	if err := os.MkdirAll(blobDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(blobDir, strings.TrimPrefix(layerDigest, "sha256:")), layerData, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(blobDir, strings.TrimPrefix(manifestDigest, "sha256:")), md, 0644); err != nil {
		t.Fatal(err)
	}

	// index.json
	idx := struct {
		Manifests []struct {
			Digest      string            `json:"digest"`
			Annotations map[string]string `json:"annotations,omitempty"`
		} `json:"manifests"`
	}{
		Manifests: []struct {
			Digest      string            `json:"digest"`
			Annotations map[string]string `json:"annotations,omitempty"`
		}{{Digest: manifestDigest, Annotations: map[string]string{"org.opencontainers.image.ref.name": tag}}},
	}
	idxData, _ := json.Marshal(idx)
	if err := os.WriteFile(filepath.Join(dir, "index.json"), idxData, 0644); err != nil {
		t.Fatal(err)
	}
}

// ---- validateOCIDigest ----

func TestValidateOCIDigest(t *testing.T) {
	cases := []struct {
		name, digest, want string
	}{
		{"empty", "", "unsupported digest algorithm"},
		{"no prefix", "abc", "unsupported digest algorithm"},
		{"wrong prefix", "sha1:abc", "unsupported digest algorithm"},
		{"wrong length", "sha256:abc", "invalid digest length"},
		{"non-hex", fmt.Sprintf("sha256:%s", strings.Repeat("z", 64)), "invalid hex character"},
		{"valid", fmt.Sprintf("sha256:%s", strings.Repeat("a", 64)), ""},
		{"valid mixed", fmt.Sprintf("sha256:%s", "a"+strings.Repeat("b", 63)), ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateOCIDigest(c.digest)
			if c.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("expected %q in %v", c.want, err)
			}
		})
	}
}

// ---- validateOCIRegistryDestination ----

func TestValidateOCIRegistryDestination_RejectsEmpty(t *testing.T) {
	if err := validateOCIRegistryDestination("", "repo"); err == nil || !strings.Contains(err.Error(), "registry is required") {
		t.Fatal("expected registry required error")
	}
	if err := validateOCIRegistryDestination("reg.example.com", ""); err == nil || !strings.Contains(err.Error(), "repository is required") {
		t.Fatal("expected repository required error")
	}
}

func TestValidateOCIRegistryDestination_RejectsLocal(t *testing.T) {
	for _, h := range []string{"localhost", "127.0.0.1", "::1"} {
		t.Run(h, func(t *testing.T) {
			err := validateOCIRegistryDestination(h, "repo")
			if err == nil || !strings.Contains(err.Error(), "loopback") {
				t.Fatalf("expected loopback error, got %v", err)
			}
		})
	}
}

func TestValidateOCIRegistryDestination_RejectsPath(t *testing.T) {
	err := validateOCIRegistryDestination("reg.example.com/path", "repo")
	if err == nil || !strings.Contains(err.Error(), "path") {
		t.Fatal("expected path error, got:", err)
	}
}

func TestValidateOCIRegistryDestination_RejectsTraversal(t *testing.T) {
	err := validateOCIRegistryDestination("reg.example.com", "../evil")
	if err == nil || !strings.Contains(err.Error(), "path traversal") {
		t.Fatal("expected path traversal error, got:", err)
	}
}

func TestValidateOCIRegistryDestination_RejectsBadScheme(t *testing.T) {
	err := validateOCIRegistryDestination("ftp://reg.example.com", "repo")
	if err == nil || !strings.Contains(err.Error(), "unsupported registry scheme") {
		t.Fatal("expected scheme error, got:", err)
	}
}

func TestValidateOCIRegistryDestination_RejectsUnspecified(t *testing.T) {
	for _, h := range []string{"0.0.0.0", "::"} {
		t.Run(h, func(t *testing.T) {
			err := validateOCIRegistryDestination(h, "repo")
			if err == nil || !strings.Contains(err.Error(), "unspecified") {
				t.Fatalf("expected unspecified address error, got %v", err)
			}
		})
	}
}

func TestValidateOCIRegistryDestination_RejectsUnspecifiedWithPort(t *testing.T) {
	err := validateOCIRegistryDestination("0.0.0.0:5000", "repo")
	if err == nil || !strings.Contains(err.Error(), "unspecified") {
		t.Fatal("expected unspecified address error, got:", err)
	}
}

// ---- resolveCredentials ----

func TestResolveCredentials_Inline(t *testing.T) {
	p := &ociCredentialSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	p.Initialize()

	u, pw, err := resolveCredentials(
		httptest.NewRequest(http.MethodPost, "/", nil).WithContext(withTokenCtx(t)),
		nil, "alice", "secret", p,
	)
	if err != nil {
		t.Fatal(err)
	}
	if u != "alice" {
		t.Errorf("got username %q, want alice", u)
	}
	if pw != "secret" {
		t.Errorf("got password %q, want secret", pw)
	}
}

func TestResolveCredentials_EmptyIDFallsBack(t *testing.T) {
	p := &ociCredentialSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	p.Initialize()

	cid := ""
	u, pw, err := resolveCredentials(
		httptest.NewRequest(http.MethodPost, "/", nil).WithContext(withTokenCtx(t)),
		&cid, "bob", "pass", p,
	)
	if err != nil {
		t.Fatal(err)
	}
	if u != "bob" || pw != "pass" {
		t.Fatalf("got %q/%q, want bob/pass", u, pw)
	}
}

func TestResolveCredentials_InvalidUUID(t *testing.T) {
	p := &ociCredentialSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	p.Initialize()

	cid := "not-a-uuid"
	_, _, err := resolveCredentials(
		httptest.NewRequest(http.MethodPost, "/", nil).WithContext(withTokenCtx(t)),
		&cid, "", "", p,
	)
	if err == nil || !strings.Contains(err.Error(), "invalid credentialId") {
		t.Fatalf("expected invalid credentialId error, got %v", err)
	}
}

func TestResolveCredentials_LookupByID(t *testing.T) {
	p := &ociCredentialSpyProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getCredentialFn: func(_ string, _ core.Uuid) (*models.Credential, int, error) {
			return &models.Credential{
				Secret: map[string]interface{}{"username": "su", "password": "sp"},
			}, 200, nil
		},
	}
	p.Initialize()

	cid := uuid.Must(uuid.NewV4()).String()
	u, pw, err := resolveCredentials(
		httptest.NewRequest(http.MethodPost, "/", nil).WithContext(withTokenCtx(t)),
		&cid, "", "", p,
	)
	if err != nil {
		t.Fatal(err)
	}
	if u != "su" {
		t.Errorf("got username %q, want su", u)
	}
	if pw != "sp" {
		t.Errorf("got password %q, want sp", pw)
	}
}

func TestResolveCredentials_ProviderError(t *testing.T) {
	p := &ociCredentialSpyProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getCredentialFn: func(_ string, _ core.Uuid) (*models.Credential, int, error) {
			return nil, 500, fmt.Errorf("provider err")
		},
	}
	p.Initialize()

	cid := uuid.Must(uuid.NewV4()).String()
	_, _, err := resolveCredentials(
		httptest.NewRequest(http.MethodPost, "/", nil).WithContext(withTokenCtx(t)),
		&cid, "", "", p,
	)
	if err == nil || !strings.Contains(err.Error(), "failed to retrieve credential") {
		t.Fatalf("expected credential retrieval error, got %v", err)
	}
}

func TestResolveCredentials_NilCredential(t *testing.T) {
	p := &ociCredentialSpyProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getCredentialFn: func(_ string, _ core.Uuid) (*models.Credential, int, error) {
			return nil, 200, nil
		},
	}
	p.Initialize()

	cid := uuid.Must(uuid.NewV4()).String()
	_, _, err := resolveCredentials(
		httptest.NewRequest(http.MethodPost, "/", nil).WithContext(withTokenCtx(t)),
		&cid, "", "", p,
	)
	if err == nil || !strings.Contains(err.Error(), "credential not found") {
		t.Fatalf("expected credential not found error, got %v", err)
	}
}

func TestResolveCredentials_MissingUsernameInStore(t *testing.T) {
	p := &ociCredentialSpyProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getCredentialFn: func(_ string, _ core.Uuid) (*models.Credential, int, error) {
			return &models.Credential{
				Secret: map[string]interface{}{"password": "sp"},
			}, 200, nil
		},
	}
	p.Initialize()

	cid := uuid.Must(uuid.NewV4()).String()
	_, _, err := resolveCredentials(
		httptest.NewRequest(http.MethodPost, "/", nil).WithContext(withTokenCtx(t)),
		&cid, "", "", p,
	)
	if err == nil || !strings.Contains(err.Error(), "no username") {
		t.Fatalf("expected missing username error, got %v", err)
	}
}

func TestResolveCredentials_MissingToken(t *testing.T) {
	p := &ociCredentialSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	p.Initialize()

	cid := uuid.Must(uuid.NewV4()).String()
	_, _, err := resolveCredentials(
		httptest.NewRequest(http.MethodPost, "/", nil),
		&cid, "", "", p,
	)
	if err == nil || !strings.Contains(err.Error(), "auth token") {
		t.Fatalf("expected auth token error, got %v", err)
	}
}

// ---- extractFromOCIStore ----

func TestExtractFromOCIStore_ValidArtifact(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	writeOCIStore(t, src, "latest", map[string]string{
		"model.json":  `{"name":"test"}`,
		"comp/c.json": `{}`,
		"rel/r.json":  `{}`,
	})

	if err := extractFromOCIStore(src, "latest", dst); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"model.json", "comp/c.json", "rel/r.json"} {
		if _, err := os.Stat(filepath.Join(dst, path)); err != nil {
			t.Errorf("%s not extracted: %v", path, err)
		}
	}
	data, _ := os.ReadFile(filepath.Join(dst, "model.json"))
	if !bytes.Contains(data, []byte("test")) {
		t.Error("model.json has wrong content")
	}
}

func TestExtractFromOCIStore_PathTraversal(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	writeOCIStore(t, src, "latest", map[string]string{"../../../etc/passwd": "x"})

	err := extractFromOCIStore(src, "latest", dst)
	if err == nil || !strings.Contains(err.Error(), "path traversal") {
		t.Fatal("expected path traversal error, got:", err)
	}
}

func TestExtractFromOCIStore_Symlink(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: "link", Typeflag: tar.TypeSymlink, Linkname: "/etc/passwd"}); err != nil {
		t.Fatal(err)
	}
	tw.Close()
	gz.Close()

	h := sha256.Sum256(buf.Bytes())
	ld := fmt.Sprintf("sha256:%s", hex.EncodeToString(h[:]))
	layerData := buf.Bytes()
	md, mDigest := func() ([]byte, string) {
		m := struct {
			Layers []struct {
				Digest string `json:"digest"`
			} `json:"layers"`
		}{Layers: []struct {
			Digest string `json:"digest"`
		}{{Digest: ld}}}
		d, _ := json.Marshal(m)
		hm := sha256.Sum256(d)
		return d, fmt.Sprintf("sha256:%s", hex.EncodeToString(hm[:]))
	}()

	os.MkdirAll(filepath.Join(src, "blobs", "sha256"), 0755)
	os.WriteFile(filepath.Join(src, "blobs", "sha256", strings.TrimPrefix(ld, "sha256:")), layerData, 0644)
	os.WriteFile(filepath.Join(src, "blobs", "sha256", strings.TrimPrefix(mDigest, "sha256:")), md, 0644)
	idx, _ := json.Marshal(struct {
		Manifests []struct {
			Digest      string            `json:"digest"`
			Annotations map[string]string `json:"annotations,omitempty"`
		} `json:"manifests"`
	}{
		Manifests: []struct {
			Digest      string            `json:"digest"`
			Annotations map[string]string `json:"annotations,omitempty"`
		}{{Digest: mDigest, Annotations: map[string]string{"org.opencontainers.image.ref.name": "latest"}}},
	})
	os.WriteFile(filepath.Join(src, "index.json"), idx, 0644)

	err := extractFromOCIStore(src, "latest", dst)
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatal("expected symlink error, got:", err)
	}
}

func TestExtractFromOCIStore_InvalidLayerDigest(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	md, mDigest := func() ([]byte, string) {
		m := struct {
			Layers []struct {
				Digest string `json:"digest"`
			} `json:"layers"`
		}{Layers: []struct {
			Digest string `json:"digest"`
		}{{Digest: "sha256:bad"}}}
		d, _ := json.Marshal(m)
		hm := sha256.Sum256(d)
		return d, fmt.Sprintf("sha256:%s", hex.EncodeToString(hm[:]))
	}()

	os.MkdirAll(filepath.Join(src, "blobs", "sha256"), 0755)
	os.WriteFile(filepath.Join(src, "blobs", "sha256", strings.TrimPrefix(mDigest, "sha256:")), md, 0644)
	idx, _ := json.Marshal(struct {
		Manifests []struct {
			Digest string `json:"digest"`
		} `json:"manifests"`
	}{Manifests: []struct {
		Digest string `json:"digest"`
	}{{Digest: mDigest}}})
	os.WriteFile(filepath.Join(src, "index.json"), idx, 0644)

	err := extractFromOCIStore(src, "latest", dst)
	if err == nil || !strings.Contains(err.Error(), "invalid layer digest") {
		t.Fatal("expected invalid layer digest error, got:", err)
	}
}

func TestExtractFromOCIStore_NoManifests(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	os.WriteFile(filepath.Join(src, "index.json"), []byte(`{"manifests":[]}`), 0644)

	err := extractFromOCIStore(src, "latest", dst)
	if err == nil || !strings.Contains(err.Error(), "no manifests") {
		t.Fatal("expected no manifests error, got:", err)
	}
}

func TestExtractFromOCIStore_MissingIndex(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	err := extractFromOCIStore(src, "latest", dst)
	if err == nil || !strings.Contains(err.Error(), "index.json") {
		t.Fatal("expected index.json error, got:", err)
	}
}

// ---- PushModel handler-level tests ----
//
// These tests exercise early-exit validation paths that do not require
// registryManager or live OCI network calls.

func TestPushModel_InvalidJSON(t *testing.T) {
	h := newTestHandler(t, nil, "")
	prov := &ociCredentialSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	prov.Initialize()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not json")).WithContext(withTokenCtx(t))
	rec := httptest.NewRecorder()

	h.PushModel(rec, req, nil, nil, prov)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestPushModel_MissingModelId(t *testing.T) {
	h := newTestHandler(t, nil, "")
	prov := &ociCredentialSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	prov.Initialize()

	body := jsonBody(t, map[string]string{
		"registry":   "reg.example.com",
		"repository": "my-org/my-model",
	})
	req := httptest.NewRequest(http.MethodPost, "/", body).WithContext(withTokenCtx(t))
	rec := httptest.NewRecorder()

	h.PushModel(rec, req, nil, nil, prov)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestPullModel_InvalidJSON(t *testing.T) {
	h := newTestHandler(t, nil, "")
	prov := &ociCredentialSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	prov.Initialize()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not json")).WithContext(withTokenCtx(t))
	rec := httptest.NewRecorder()

	h.PullModel(rec, req, nil, nil, prov)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

// ---- extractFromOCIStore ----

func TestExtractFromOCIStore_NoLayers(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	// Manifest with an empty layers array — exercises len(manifest.Layers) == 0.
	md := []byte(`{"layers":[]}`)
	hm := sha256.Sum256(md)
	mDigest := fmt.Sprintf("sha256:%s", hex.EncodeToString(hm[:]))

	os.MkdirAll(filepath.Join(src, "blobs", "sha256"), 0755)
	os.WriteFile(filepath.Join(src, "blobs", "sha256", strings.TrimPrefix(mDigest, "sha256:")), md, 0644)
	idx, _ := json.Marshal(struct {
		Manifests []struct {
			Digest string `json:"digest"`
		} `json:"manifests"`
	}{Manifests: []struct {
		Digest string `json:"digest"`
	}{{Digest: mDigest}}})
	os.WriteFile(filepath.Join(src, "index.json"), idx, 0644)

	err := extractFromOCIStore(src, "latest", dst)
	if err == nil || !strings.Contains(err.Error(), "no layers") {
		t.Fatal("expected no layers error, got:", err)
	}
}
