package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/helpers"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/viper"
)

// TestK8sEventMetadataHasError guards the decision that raises a Kubernetes
// connection receipt event to Error severity. A receipt whose per-context
// metadata records any failure must be reported as Error so it persists in the
// notification center and stays retrievable under the Error severity filter
// (issue #20725); a receipt describing only successful connections must remain
// Informational.
func TestK8sEventMetadataHasError(t *testing.T) {
	tests := []struct {
		name          string
		eventMetadata map[string]interface{}
		want          bool
	}{
		{
			name:          "nil metadata reports no error",
			eventMetadata: nil,
			want:          false,
		},
		{
			name:          "empty metadata reports no error",
			eventMetadata: map[string]interface{}{},
			want:          false,
		},
		{
			name: "present-but-nil error value reports no error",
			eventMetadata: map[string]interface{}{
				"prod": map[string]interface{}{
					"error": nil,
				},
			},
			want: false,
		},
		{
			name: "only successful contexts report no error",
			eventMetadata: map[string]interface{}{
				"prod": map[string]interface{}{
					"description": "Connection registered with kubernetes context \"prod\".",
				},
				"staging": map[string]interface{}{
					"description": "Connection already exists with Kubernetes context \"staging\".",
				},
			},
			want: false,
		},
		{
			name: "a single failed context reports an error",
			eventMetadata: map[string]interface{}{
				"unreachable": map[string]interface{}{
					"description": "Unable to establish connection with context \"unreachable\".",
					"error":       errors.New("api server unreachable"),
				},
			},
			want: true,
		},
		{
			name: "a failure mixed with successes reports an error",
			eventMetadata: map[string]interface{}{
				"prod": map[string]interface{}{
					"description": "Connection registered with kubernetes context \"prod\".",
				},
				"unreachable": map[string]interface{}{
					"error": errors.New("api server unreachable"),
				},
			},
			want: true,
		},
		{
			name: "in-cluster nil context error reports an error",
			eventMetadata: map[string]interface{}{
				"in-cluster": map[string]interface{}{
					"description": "No contexts detected in the in-cluster kubeconfig.",
					"error":       ErrNilClusterContext(nil),
				},
			},
			want: true,
		},
		{
			name: "non-map entries are ignored",
			eventMetadata: map[string]interface{}{
				"weird": "not a metadata map",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := k8sEventMetadataHasError(tt.eventMetadata); got != tt.want {
				t.Errorf("k8sEventMetadataHasError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestErrNilClusterContext verifies that ErrNilClusterContext handles a nil error gracefully
// without panicking on err.Error(), and includes appropriate error code.
func TestErrNilClusterContext(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "nil error",
			err:  nil,
		},
		{
			name: "non-nil error",
			err:  errors.New("in-cluster context not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := ErrNilClusterContext(tc.err)
			if res == nil {
				t.Fatal("expected non-nil error, got nil")
			}
			if code := meshkiterrors.GetCode(res); code != ErrNilClusterContextCode {
				t.Fatalf("expected error code %s, got %s", ErrNilClusterContextCode, code)
			}
			if tc.err != nil && !strings.Contains(res.Error(), tc.err.Error()) {
				t.Fatalf("expected error to contain %q, got %q", tc.err.Error(), res.Error())
			}
		})
	}
}

func TestDiscoverK8SContextFromKubeConfigTypedNilInstanceID(t *testing.T) {
	orig := viper.Get("INSTANCE_ID")
	t.Cleanup(func() { viper.Set("INSTANCE_ID", orig) })

	var typedNil *core.Uuid
	viper.Set("INSTANCE_ID", typedNil)

	_, err := (&Handler{}).DiscoverK8SContextFromKubeConfig("user-1", "tok", nil)
	if !errors.Is(err, models.ErrMesheryInstanceID) {
		t.Fatalf("err = %v, want models.ErrMesheryInstanceID", err)
	}
}

func TestK8SConfigHandler_Post_NilUser_ReturnsUnauthorized(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := &models.DefaultLocalProvider{}

	req := httptest.NewRequest(http.MethodPost, "/api/system/kubernetes", nil)
	rec := httptest.NewRecorder()

	h.K8SConfigHandler(rec, req, nil, nil, p)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if body["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", body["error"])
	}
}

func TestK8SConfigHandler_Post_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/system/kubernetes", nil)
	rec := httptest.NewRecorder()

	h.K8SConfigHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetContextsFromK8SConfig_NilUser_ReturnsUnauthorized(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := &models.DefaultLocalProvider{}

	req := httptest.NewRequest(http.MethodPost, "/api/system/kubernetes/contexts", nil)
	rec := httptest.NewRecorder()

	h.GetContextsFromK8SConfig(rec, req, nil, nil, p)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if body["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", body["error"])
	}
}

func TestGetContextsFromK8SConfig_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/system/kubernetes/contexts", nil)
	rec := httptest.NewRecorder()

	h.GetContextsFromK8SConfig(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestKubernetesPingHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/ping", nil)
	rec := httptest.NewRecorder()

	h.KubernetesPingHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestKubernetesPingHandler_MissingToken_ReturnsUnauthorized(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := &models.DefaultLocalProvider{}

	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/ping", nil)
	rec := httptest.NewRecorder()

	h.KubernetesPingHandler(rec, req, nil, nil, p)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for missing token, got %d", rec.Code)
	}
}

func TestK8sRegistrationHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := &models.DefaultLocalProvider{}

	req := httptest.NewRequest(http.MethodPost, "/api/system/kubernetes/register", nil)
	rec := httptest.NewRecorder()

	h.K8sRegistrationHandler(rec, req, nil, nil, p)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if body["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", body["error"])
	}
}

func TestK8sRegistrationHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/system/kubernetes/register", nil)
	rec := httptest.NewRecorder()

	h.K8sRegistrationHandler(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

type mockK8sConfigProvider struct {
	*models.DefaultLocalProvider
	persistedEvents []events.Event
}

func (m *mockK8sConfigProvider) PersistEvent(event events.Event, _ string) (models.PersistedEvent, error) {
	m.persistedEvents = append(m.persistedEvents, event)
	return models.PersistedEvent{Event: &event, Disposition: models.EventCreated}, nil
}

func TestGetContextsFromK8SConfig_NilLogger_UploadSingleContextSucceeds(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	h.log = nil
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}
	provider := &mockK8sConfigProvider{}

	kubeconfig := []byte(`apiVersion: v1
kind: Config
current-context: ctx-0
clusters:
- cluster:
    server: https://127.0.0.1:59900
  name: cluster-0
contexts:
- context:
    cluster: cluster-0
    user: user-0
  name: ctx-0
users:
- name: user-0
  user:
    token: token-0
`)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, err := w.CreateFormFile("k8sfile", "kubeconfig")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(kubeconfig); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/system/kubernetes/contexts", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()

	h.GetContextsFromK8SConfig(rec, req, nil, authUser, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}
	var contexts []models.K8sContext
	if err := json.Unmarshal(rec.Body.Bytes(), &contexts); err != nil {
		t.Fatalf("unmarshaling contexts response: %v", err)
	}
	if len(contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(contexts))
	}
	if contexts[0].Name != "ctx-0" {
		t.Fatalf("expected context name 'ctx-0', got %q", contexts[0].Name)
	}
}

func TestKubeconfigDiscoveryPayloadStaysAnArrayUnlessAsked(t *testing.T) {
	discovery := models.KubeconfigDiscovery{
		Contexts: []*models.K8sContext{{Name: "prod", Server: "https://cluster.example"}},
		SkippedContexts: []models.SkippedKubeContext{{
			Name:   "eks",
			Reason: "Authentication method is not supported (exec plugin).",
		}},
	}

	arrayBody, err := json.Marshal(kubeconfigDiscoveryPayload(false, discovery))
	if err != nil {
		t.Fatal(err)
	}
	var asArray []map[string]interface{}
	if err := json.Unmarshal(arrayBody, &asArray); err != nil {
		t.Fatalf("default payload must stay a JSON array: %v (%s)", err, arrayBody)
	}
	if len(asArray) != 1 || asArray[0]["name"] != "prod" {
		t.Fatalf("array = %#v, want the prod context", asArray)
	}

	emptyBody, err := json.Marshal(kubeconfigDiscoveryPayload(false, models.KubeconfigDiscovery{}))
	if err != nil {
		t.Fatal(err)
	}
	if string(emptyBody) != "[]" {
		t.Fatalf("empty default payload = %s, want []", emptyBody)
	}

	detailBody, err := json.Marshal(kubeconfigDiscoveryPayload(true, discovery))
	if err != nil {
		t.Fatal(err)
	}
	var asObject map[string]json.RawMessage
	if err := json.Unmarshal(detailBody, &asObject); err != nil {
		t.Fatalf("detail payload must be an object: %v (%s)", err, detailBody)
	}
	var skipped []map[string]interface{}
	if err := json.Unmarshal(asObject["skippedContexts"], &skipped); err != nil {
		t.Fatalf("skippedContexts: %v", err)
	}
	if len(skipped) != 1 || skipped[0]["name"] != "eks" {
		t.Fatalf("skipped = %#v, want eks", skipped)
	}
}

// An uploaded kubeconfig must never make Meshery Server run a command named in
// it. models.DiscoverKubeconfigContexts proves the skip for each source; this
// pins the HTTP entry points to the upload source, so switching one of them to
// KubeconfigFromServer fails here even though every model test stays green.
func TestUploadedExecKubeconfigNeverRunsItsCommand(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("needs a POSIX shell to stand in for an exec credential plugin")
	}
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	// The command is the bare name "sh", not an absolute path: the import
	// flattens the kubeconfig and inlines any value that is a path on disk,
	// which would break the plugin for a reason unrelated to the skip.
	upload := func(t *testing.T, marker string, fields map[string]string) *http.Request {
		t.Helper()
		kubeconfig := fmt.Sprintf(`apiVersion: v1
kind: Config
current-context: eks
clusters:
- cluster:
    server: https://127.0.0.1:1
  name: eks
contexts:
- context:
    cluster: eks
    user: eks-user
  name: eks
users:
- name: eks-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: sh
      args: ["-c", "touch '%s'; exit 1"]
`, marker)
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		part, err := w.CreateFormFile("k8sfile", "kubeconfig")
		if err != nil {
			t.Fatalf("create form file: %v", err)
		}
		if _, err := part.Write([]byte(kubeconfig)); err != nil {
			t.Fatalf("write form file: %v", err)
		}
		for k, v := range fields {
			if err := w.WriteField(k, v); err != nil {
				t.Fatalf("write field %s: %v", k, err)
			}
		}
		if err := w.Close(); err != nil {
			t.Fatalf("close writer: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/system/kubernetes", &b)
		req.Header.Set("Content-Type", w.FormDataContentType())
		return req
	}
	assertNotRun := func(t *testing.T, marker string) {
		t.Helper()
		if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("exec command from an uploaded kubeconfig ran on the server (marker stat err = %v)", err)
		}
	}

	t.Run("discovery skips it and reports why", func(t *testing.T) {
		h := newTestHandler(t, map[string]models.Provider{}, "")
		marker := filepath.Join(t.TempDir(), "ran")
		rec := httptest.NewRecorder()
		h.GetContextsFromK8SConfig(rec, upload(t, marker, map[string]string{DiscoveryDetailFormKey: "true"}), nil, authUser, &mockK8sConfigProvider{})
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}
		assertNotRun(t, marker)
		var body struct {
			Contexts        []models.K8sContext         `json:"contexts"`
			SkippedContexts []models.SkippedKubeContext `json:"skippedContexts"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshaling discovery response: %v", err)
		}
		if len(body.Contexts) != 0 || len(body.SkippedContexts) != 1 || body.SkippedContexts[0].Name != "eks" {
			t.Fatalf("contexts=%d skipped=%+v, want only eks skipped", len(body.Contexts), body.SkippedContexts)
		}
	})

	t.Run("import does not run it or register the context", func(t *testing.T) {
		h := newTestHandler(t, map[string]models.Provider{}, "")
		marker := filepath.Join(t.TempDir(), "ran")
		rec := httptest.NewRecorder()
		req := upload(t, marker, nil)
		req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
		h.K8SConfigHandler(rec, req, nil, authUser, &mockK8sConfigProvider{})
		if rec.Code != http.StatusOK {
			t.Fatalf("expected the import to reach discovery and return 200, got %d: %s", rec.Code, rec.Body.String())
		}
		assertNotRun(t, marker)
		var body SaveK8sContextResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshaling import response (status %d): %v: %s", rec.Code, err, rec.Body.String())
		}
		if n := len(body.RegisteredContexts) + len(body.ConnectedContexts) + len(body.IgnoredContexts) + len(body.ErroredContexts); n != 0 {
			t.Fatalf("import returned %d contexts, want the exec context left out", n)
		}
	})
}

// uploadKubeconfigRequest builds the multipart request both kubeconfig upload
// endpoints read, carrying kubeconfig under the k8sfile form field.
func uploadKubeconfigRequest(t *testing.T, target string, kubeconfig []byte) *http.Request {
	t.Helper()
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, err := w.CreateFormFile("k8sfile", "kubeconfig")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(kubeconfig); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, target, &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

// TestUploadedKubeconfigNeverReadsServerSideFiles pins issue #24219: an uploaded
// kubeconfig that points client-key (or another credential field) at a path is
// refused outright. Before the fix FlattenMinifyKubeConfig read the file on the
// Meshery Server host and inlined it as client-key-data, and a missing path was
// silently kept, so responses told the uploader which server paths exist.
func TestUploadedKubeconfigNeverReadsServerSideFiles(t *testing.T) {
	secretFile := filepath.Join(t.TempDir(), "server-secret.key")
	const secretContent = "TOP-SECRET-SERVER-FILE-CONTENT"
	if err := os.WriteFile(secretFile, []byte(secretContent), 0o600); err != nil {
		t.Fatalf("write secret file: %v", err)
	}
	missingFile := filepath.Join(t.TempDir(), "missing", "client.key")

	kubeconfigWith := func(keyPath string) []byte {
		return []byte(`apiVersion: v1
kind: Config
current-context: ctx-0
clusters:
- cluster:
    server: https://127.0.0.1:59900
  name: cluster-0
contexts:
- context:
    cluster: cluster-0
    user: user-0
  name: ctx-0
users:
- name: user-0
  user:
    client-certificate-data: YQ==
    client-key: ` + keyPath + `
`)
	}

	// Same credential field, but carried by a document whose generic YAML decode
	// fails on an unknown key while the typed decode the server uses succeeds.
	kubeconfigWithJunk := func(keyPath string) []byte {
		return append([]byte("junk: {[a]: b}\n"), kubeconfigWith(keyPath)...)
	}

	endpoints := map[string]func(h *Handler, rec *httptest.ResponseRecorder, req *http.Request, u *models.User, p models.Provider){
		"/api/system/kubernetes/contexts": func(h *Handler, rec *httptest.ResponseRecorder, req *http.Request, u *models.User, p models.Provider) {
			h.GetContextsFromK8SConfig(rec, req, nil, u, p)
		},
		"/api/system/kubernetes/register": func(h *Handler, rec *httptest.ResponseRecorder, req *http.Request, u *models.User, p models.Provider) {
			h.K8sRegistrationHandler(rec, req, nil, u, p)
		},
		"/api/system/kubernetes": func(h *Handler, rec *httptest.ResponseRecorder, req *http.Request, u *models.User, p models.Provider) {
			h.K8SConfigHandler(rec, req, nil, u, p)
		},
	}

	for target, call := range endpoints {
		for name, keyPath := range map[string]string{"existing file": secretFile, "missing file": missingFile} {
			for shape, build := range map[string]func(string) []byte{"plain": kubeconfigWith, "unknown junk key": kubeconfigWithJunk} {
				t.Run(target+" "+name+" "+shape, func(t *testing.T) {
					h := newTestHandler(t, map[string]models.Provider{}, "")
					authUser := &models.User{ID: uuid.Must(uuid.NewV4())}
					provider := &mockK8sConfigProvider{}
					req := uploadKubeconfigRequest(t, target, build(keyPath))
					req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "token"))
					rec := httptest.NewRecorder()

					call(h, rec, req, authUser, provider)

					if rec.Code != http.StatusBadRequest {
						t.Fatalf("status = %d, want 400; body: %s", rec.Code, rec.Body.String())
					}
					body := rec.Body.String()
					for _, leak := range []string{secretContent, "VE9QLVNFQ1JFVC1TRVJWRVItRklMRS1DT05URU5U", keyPath, filepath.Base(keyPath)} {
						if strings.Contains(body, leak) {
							t.Fatalf("response leaks %q: %s", leak, body)
						}
					}
					if decodeErrorBody(t, rec.Body.Bytes()).Code != helpers.ErrKubeConfigFileReferenceCode {
						t.Fatalf("unexpected error code in %s", body)
					}
				})
			}
		}
	}
}

// TestImportReceiptIgnoresFailingDeselectedContext pins issue #24220: a context
// that fails discovery but that the user deselected must leave no error in the
// import receipt, which then stays Informational like the wizard's success.
func TestImportReceiptIgnoresFailingDeselectedContext(t *testing.T) {
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}
	// Nothing listens on port 1, so discovery cannot reach the API server and
	// records an error for the context.
	kubeconfig := []byte(`apiVersion: v1
kind: Config
current-context: broken
clusters:
- cluster:
    server: https://127.0.0.1:1
  name: broken
contexts:
- context:
    cluster: broken
    user: broken-user
  name: broken
users:
- name: broken-user
  user:
    token: token-0
`)
	receiptOf := func(t *testing.T, form map[string]string) events.Event {
		t.Helper()
		h := newTestHandler(t, map[string]models.Provider{}, "")
		provider := &mockK8sConfigProvider{}
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		part, err := w.CreateFormFile("k8sfile", "kubeconfig")
		if err != nil {
			t.Fatalf("create form file: %v", err)
		}
		if _, err := part.Write(kubeconfig); err != nil {
			t.Fatalf("write form file: %v", err)
		}
		for k, v := range form {
			if err := w.WriteField(k, v); err != nil {
				t.Fatalf("write field %s: %v", k, err)
			}
		}
		if err := w.Close(); err != nil {
			t.Fatalf("close writer: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/system/kubernetes", &b)
		req.Header.Set("Content-Type", w.FormDataContentType())
		req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
		rec := httptest.NewRecorder()
		h.K8SConfigHandler(rec, req, nil, authUser, provider)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}
		if len(provider.persistedEvents) != 1 {
			t.Fatalf("persisted %d events, want the single import receipt", len(provider.persistedEvents))
		}
		return provider.persistedEvents[0]
	}

	t.Run("deselected failing context leaves an informational receipt", func(t *testing.T) {
		receipt := receiptOf(t, map[string]string{SelectedContextsFormKey: "[]"})
		if receipt.Severity != events.Informational {
			t.Fatalf("receipt severity = %v, want Informational (description %q)", receipt.Severity, receipt.Description)
		}
		if k8sEventMetadataHasError(receipt.Metadata) {
			t.Fatalf("receipt metadata carries an error for a context the user deselected: %#v", receipt.Metadata)
		}
		if _, ok := receipt.Metadata["broken"]; ok {
			t.Fatalf("receipt metadata describes the deselected context: %#v", receipt.Metadata)
		}
	})
}

func TestRetainSelectedContextMetadata(t *testing.T) {
	picked := &models.K8sContext{ID: "id-picked", Name: "picked"}
	dropped := &models.K8sContext{ID: "id-dropped", Name: "dropped"}
	failing := map[string]interface{}{"error": errors.New("unreachable")}
	newMetadata := func() map[string]interface{} {
		return map[string]interface{}{"picked": failing, "dropped": failing, "skipped": failing}
	}

	t.Run("a selection keeps only the selected contexts", func(t *testing.T) {
		metadata := newMetadata()
		retainSelectedContextMetadata(metadata, []*models.K8sContext{picked, dropped}, map[string]struct{}{"id-picked": {}})
		if len(metadata) != 1 || metadata["picked"] == nil {
			t.Fatalf("metadata = %#v, want only the selected context's entry", metadata)
		}
	})

	t.Run("no selection removes nothing, so skipped and failed contexts still show", func(t *testing.T) {
		metadata := newMetadata()
		retainSelectedContextMetadata(metadata, []*models.K8sContext{picked, dropped}, nil)
		if len(metadata) != 3 || !k8sEventMetadataHasError(metadata) {
			t.Fatalf("metadata = %#v, want all three entries kept", metadata)
		}
	})
}

type mockReimportProvider struct {
	models.Provider
	sync.Mutex
	persistedEvents []events.Event
	connectionID    core.Uuid
}

func (m *mockReimportProvider) SaveK8sContext(token string, k8sContext models.K8sContext, metadata map[string]any) (connections.Connection, error) {
	return connections.Connection{ID: m.connectionID, Status: connections.DISCOVERED}, nil
}

func (m *mockReimportProvider) GetConnectionByID(token string, connectionID core.Uuid) (*connections.Connection, int, error) {
	return &connections.Connection{ID: m.connectionID, Status: connections.DISCOVERED}, 0, nil
}

func (m *mockReimportProvider) UpdateConnectionStatusByID(token string, connectionID core.Uuid, connectionStatus connections.ConnectionStatus) (*connections.Connection, int, error) {
	return &connections.Connection{ID: m.connectionID, Status: connectionStatus}, 0, nil
}

func (m *mockReimportProvider) UpdateConnectionStatusIf(token string, connectionID core.Uuid, fromStatuses []connections.ConnectionStatus, newStatus connections.ConnectionStatus) (bool, *connections.Connection, error) {
	return true, &connections.Connection{ID: m.connectionID, Status: newStatus}, nil
}

func (m *mockReimportProvider) PersistEvent(event events.Event, _ string) (models.PersistedEvent, error) {
	m.Lock()
	defer m.Unlock()
	m.persistedEvents = append(m.persistedEvents, event)
	return models.PersistedEvent{Event: &event, Disposition: models.EventCreated}, nil
}

func TestK8SConfigHandler_ReimportRefreshesCachedStateMachineContext(t *testing.T) {
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}
	connUUID := uuid.Must(uuid.NewV4())
	connID := core.Uuid(connUUID)

	kubeconfig := []byte(`apiVersion: v1
kind: Config
current-context: reimported
clusters:
- cluster:
    server: https://127.0.0.1:6443
  name: reimported
contexts:
- context:
    cluster: reimported
    user: reimported-user
  name: reimported
users:
- name: reimported-user
  user:
    token: token-reimported
`)

	t.Run("refreshes credentials on cache hit with existing context", func(t *testing.T) {
		h := newTestHandler(t, map[string]models.Provider{}, "")
		h.config.K8scontextChannel = &models.K8scontextChan{}
		tracker := &machines.ConnectionToStateMachineInstanceTracker{
			ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
		}
		h.ConnectionToStateMachineInstanceTracker = tracker

		provider := &mockReimportProvider{
			connectionID: connID,
		}

		inst, err := kubernetes.New(connID.String(), core.Uuid(uuid.Nil), h.log)
		if err != nil {
			t.Fatalf("failed to create state machine: %v", err)
		}
		inst.Provider = provider
		storedCtx := &kubernetes.MachineCtx{
			K8sContext: models.K8sContext{
				Name:   "old-name",
				Server: "https://old-server:6443",
			},
		}
		inst.Context = storedCtx
		tracker.Add(connID, inst)

		req := uploadKubeconfigRequest(t, "/api/system/kubernetes", kubeconfig)
		req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
		rec := httptest.NewRecorder()

		h.K8SConfigHandler(rec, req, nil, authUser, provider)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}

		ctxObj := inst.GetContext()
		updatedCtx, ok := ctxObj.(*kubernetes.MachineCtx)
		if !ok || updatedCtx == nil {
			t.Fatalf("expected *kubernetes.MachineCtx, got %#v", ctxObj)
		}
		if updatedCtx.K8sContext.Server != "https://127.0.0.1:6443" {
			t.Fatalf("expected server https://127.0.0.1:6443, got %s", updatedCtx.K8sContext.Server)
		}
	})

	t.Run("initializes context on cache hit when initial startup left Context nil", func(t *testing.T) {
		h := newTestHandler(t, map[string]models.Provider{}, "")
		h.config.K8scontextChannel = &models.K8scontextChan{}
		tracker := &machines.ConnectionToStateMachineInstanceTracker{
			ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
		}
		h.ConnectionToStateMachineInstanceTracker = tracker

		provider := &mockReimportProvider{
			connectionID: connID,
		}

		inst, err := kubernetes.New(connID.String(), core.Uuid(uuid.Nil), h.log)
		if err != nil {
			t.Fatalf("failed to create state machine: %v", err)
		}
		inst.Provider = provider
		inst.Context = nil
		tracker.Add(connID, inst)

		req := uploadKubeconfigRequest(t, "/api/system/kubernetes", kubeconfig)
		req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
		rec := httptest.NewRecorder()

		h.K8SConfigHandler(rec, req, nil, authUser, provider)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}

		ctxObj := inst.GetContext()
		initializedCtx, ok := ctxObj.(*kubernetes.MachineCtx)
		if !ok || initializedCtx == nil {
			t.Fatalf("expected initialized *kubernetes.MachineCtx, got %#v", ctxObj)
		}
		if initializedCtx.K8sContext.Server != "https://127.0.0.1:6443" {
			t.Fatalf("expected server https://127.0.0.1:6443, got %s", initializedCtx.K8sContext.Server)
		}
	})
}
