package models

import (
	"context"
	"crypto/x509"
	"encoding/json"
	stderrors "errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"syscall"
	"testing"
	"time"

	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/events"
	mkkube "github.com/meshery/meshkit/utils/kubernetes"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// TestErrUnreachableKubeAPICauses pins the nested causes: each cause carries
// its own remediations, one remediation is shared by every cause, the ids
// match the exported constants, and the flat lists older readers use are
// derived from the nesting.
func TestErrUnreachableKubeAPICauses(t *testing.T) {
	err := ErrUnreachableKubeAPI(stderrors.New("boom"), "https://cluster:6443")

	if code := meshkiterrors.GetCode(err); code != ErrUnreachableKubeAPICode {
		t.Fatalf("code = %q, want %q (the matched-cause wrapper must stay transparent)", code, ErrUnreachableKubeAPICode)
	}

	causes := meshkiterrors.GetCauses(err)
	ids := make([]string, len(causes))
	for i, c := range causes {
		ids[i] = c.ID
		if len(c.Remediations) < 2 {
			t.Errorf("cause %q has %d remediations, want its own plus the shared retry", c.ID, len(c.Remediations))
			continue
		}
		// NewV3 stores an empty policy as Manual, so compare what the cause says.
		if shared := c.Remediations[len(c.Remediations)-1]; shared.ID != remRetryKubeConnection.ID ||
			shared.Description != remRetryKubeConnection.Description || shared.ActionPolicy != meshkiterrors.Manual {
			t.Errorf("cause %q does not end with the shared retry remediation: %+v", c.ID, shared)
		}
	}
	wantIDs := []string{KubeAPICauseUnreachable, KubeAPICauseCredentialsRejected, KubeAPICauseCertificateUntrusted}
	if !reflect.DeepEqual(ids, wantIDs) {
		t.Fatalf("cause ids = %v, want %v", ids, wantIDs)
	}
	if verr := meshkiterrors.ValidateCauses(causes); verr != nil {
		t.Fatalf("causes break the error definition contract: %v", verr)
	}

	// The shared remediation appears once in the flattened list.
	wantRemediations := []string{
		"Verify network connectivity and Kubernetes API responsiveness between Meshery Server and your cluster.",
		remRetryKubeConnection.Description,
		"Ensure client credential is not expired and is properly formed.",
		"Remove the cluster credential and enable 'insecure-skip-tls-verify'.",
	}
	mk, _ := meshkiterrors.AsErrorV2(err)
	if !reflect.DeepEqual(mk.SuggestedRemediation, wantRemediations) {
		t.Fatalf("flattened remediations = %q, want %q", mk.SuggestedRemediation, wantRemediations)
	}
	if len(mk.ProbableCause) != 3 || mk.ProbableCause[0] != "The Kubernetes API server is not reachable." {
		t.Fatalf("flattened causes = %q", mk.ProbableCause)
	}
}

func TestClassifyKubeAPIError(t *testing.T) {
	refused := &url.Error{Op: "Get", URL: "https://cluster:6443", Err: &net.OpError{Op: "dial", Net: "tcp", Err: syscall.ECONNREFUSED}}
	untrusted := &url.Error{Op: "Get", URL: "https://cluster:6443", Err: x509.UnknownAuthorityError{}}
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"connection refused", refused, KubeAPICauseUnreachable},
		{"deadline", context.DeadlineExceeded, KubeAPICauseUnreachable},
		{"unknown certificate authority", untrusted, KubeAPICauseCertificateUntrusted},
		{"unauthorized", apierrors.NewUnauthorized("token expired"), KubeAPICauseCredentialsRejected},
		{"forbidden is not a credential failure", apierrors.NewForbidden(schema.GroupResource{Resource: "namespaces"}, "kube-system", stderrors.New("rbac")), ""},
		{"unknown", stderrors.New("something else"), ""},
	}
	for _, c := range cases {
		if got := classifyKubeAPIError(c.err); got != c.want {
			t.Errorf("%s: classifyKubeAPIError = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestAttachEventErrorNamesTheMatchedCause(t *testing.T) {
	refused := &net.OpError{Op: "dial", Net: "tcp", Err: syscall.ECONNREFUSED}
	err := ErrUnreachableKubeAPI(refused, "https://cluster:6443")

	event := AttachEventError(events.NewEvent().WithMetadata(map[string]interface{}{"connectionName": "prod"}), err, nil).Build()

	if event.Metadata["connectionName"] != "prod" {
		t.Fatalf("AttachEventError dropped other metadata: %v", event.Metadata)
	}
	ee, ok := event.Metadata[events.MetadataErrorKey].(events.EventError)
	if !ok {
		t.Fatalf("metadata.error is %T, want events.EventError", event.Metadata[events.MetadataErrorKey])
	}
	if ee.Code != ErrUnreachableKubeAPICode {
		t.Fatalf("code = %q", ee.Code)
	}
	if ee.ProbableCauses == nil || len(*ee.ProbableCauses) != 3 {
		t.Fatalf("probableCauses = %v, want the three nested causes", ee.ProbableCauses)
	}
	if ee.MatchedCauses == nil || len(*ee.MatchedCauses) != 1 || (*ee.MatchedCauses)[0].CauseID != KubeAPICauseUnreachable {
		t.Fatalf("matchedCauses = %v, want %q", ee.MatchedCauses, KubeAPICauseUnreachable)
	}
	if verr := event.Validate(); verr != nil {
		t.Fatalf("event breaks the schemas event contract: %v", verr)
	}

	// It serializes in the schemas shape, keyed by camelCase names.
	raw, mErr := json.Marshal(event.Metadata[events.MetadataErrorKey])
	if mErr != nil {
		t.Fatalf("marshal: %v", mErr)
	}
	var wire map[string]interface{}
	_ = json.Unmarshal(raw, &wire)
	for _, key := range []string{"code", "probableCauses", "matchedCauses", "probableCause", "suggestedRemediation"} {
		if _, ok := wire[key]; !ok {
			t.Errorf("metadata.error JSON lacks %q: %s", key, raw)
		}
	}
}

// An error that is not a MeshKit error keeps the legacy metadata.error form,
// so the event still reports the failure.
func TestAttachEventErrorKeepsLegacyFormForNonMeshkitErrors(t *testing.T) {
	plain := stderrors.New("plain failure")
	event := AttachEventError(events.NewEvent(), plain, nil).Build()
	// Build deep-copies the metadata, so the stored error is a copy of the
	// error, equal to it but not the same pointer.
	got, ok := event.Metadata[events.MetadataErrorKey].(error)
	if !ok || got.Error() != plain.Error() {
		t.Fatalf("metadata.error = %#v, want the error itself", event.Metadata[events.MetadataErrorKey])
	}
}

// A matched-cause error serializes like the MeshKit error it wraps, so rows
// that store the legacy form are unchanged.
func TestMatchedCauseErrorSerializesAsTheMeshkitError(t *testing.T) {
	err := ErrUnreachableKubeAPI(context.DeadlineExceeded, "https://cluster:6443")
	mk, _ := meshkiterrors.AsErrorV2(err)

	got, gErr := json.Marshal(err)
	want, wErr := json.Marshal(&mk)
	if gErr != nil || wErr != nil {
		t.Fatalf("marshal: %v %v", gErr, wErr)
	}
	if string(got) != string(want) {
		t.Fatalf("wrapper JSON = %s, want the MeshKit error's %s", got, want)
	}
}

// A Discovery that cannot reach the cluster reports AssignServerID's error,
// which is already an ErrUnreachableKubeAPI: the event the discover action
// builds from it carries the nested causes and the matched cause, not the
// legacy raw error.
func TestAssignServerIDFailureAttachesTheMatchedCause(t *testing.T) {
	// Nothing listens on port 1, so the API call is refused at once.
	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: "http://127.0.0.1:1", Timeout: 5 * time.Second})
	if err != nil {
		t.Fatalf("build clientset: %v", err)
	}
	kc := &K8sContext{Server: "http://127.0.0.1:1"}
	assignErr := kc.AssignServerID(context.Background(), &mkkube.Client{KubeClient: clientset})
	if assignErr == nil {
		t.Fatal("AssignServerID succeeded against a closed port")
	}

	event := AttachEventError(events.NewEvent().WithMetadata(map[string]interface{}{}), assignErr, nil).Build()
	ee, ok := event.Metadata[events.MetadataErrorKey].(events.EventError)
	if !ok {
		t.Fatalf("metadata.error is %T, want events.EventError", event.Metadata[events.MetadataErrorKey])
	}
	if ee.Code != ErrUnreachableKubeAPICode {
		t.Fatalf("code = %q, want %q", ee.Code, ErrUnreachableKubeAPICode)
	}
	if ee.MatchedCauses == nil || len(*ee.MatchedCauses) != 1 || (*ee.MatchedCauses)[0].CauseID != KubeAPICauseUnreachable {
		t.Fatalf("matchedCauses = %v, want %q", ee.MatchedCauses, KubeAPICauseUnreachable)
	}
}

// ErrUnreachableKubeAPI is exported, so a nil client error must not panic.
func TestErrUnreachableKubeAPIAcceptsANilError(t *testing.T) {
	err := ErrUnreachableKubeAPI(nil, "https://cluster:6443")
	if mk, ok := meshkiterrors.AsErrorV2(err); !ok || mk.Code != ErrUnreachableKubeAPICode {
		t.Fatalf("got %v, want an ErrUnreachableKubeAPI", err)
	}
}

// TestAssignServerID_HonorsContextCancellation verifies that AssignServerID immediately
// aborts when invoked with a pre-canceled context.
func TestAssignServerID_HonorsContextCancellation(t *testing.T) {
	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: "http://127.0.0.1:1", Timeout: 10 * time.Second})
	if err != nil {
		t.Fatalf("build clientset: %v", err)
	}
	kc := &K8sContext{Server: "http://127.0.0.1:1"}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel context

	start := time.Now()
	err = kc.AssignServerID(ctx, &mkkube.Client{KubeClient: clientset})
	if err == nil {
		t.Fatal("expected error with canceled context, got nil")
	}
	if time.Since(start) > 2*time.Second {
		t.Fatalf("call took %v, expected immediate abort on canceled context", time.Since(start))
	}
}

// TestAssignVersion_HonorsContextCancellation verifies that AssignVersion immediately
// aborts when invoked with a pre-canceled context.
func TestAssignVersion_HonorsContextCancellation(t *testing.T) {
	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: "http://127.0.0.1:1", Timeout: 10 * time.Second})
	if err != nil {
		t.Fatalf("build clientset: %v", err)
	}
	kc := &K8sContext{Server: "http://127.0.0.1:1"}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel context

	start := time.Now()
	err = kc.AssignVersion(ctx, &mkkube.Client{KubeClient: clientset})
	if err == nil {
		t.Fatal("expected error with canceled context, got nil")
	}
	if time.Since(start) > 2*time.Second {
		t.Fatalf("call took %v, expected immediate abort on canceled context", time.Since(start))
	}
}

// TestAssignServerIDAndVersion_StalledServerHonorsDeadline asserts that both AssignServerID
// and AssignVersion respect context deadlines when connecting to an unresponsive server.
func TestAssignServerIDAndVersion_StalledServerHonorsDeadline(t *testing.T) {
	// A server that waits on the request context simulates a cluster that became unresponsive
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.CloseClientConnections()
	defer srv.Close()

	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: srv.URL, Timeout: 30 * time.Second})
	if err != nil {
		t.Fatalf("build clientset: %v", err)
	}
	kc := &K8sContext{Server: srv.URL}

	t.Run("AssignServerID", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		start := time.Now()
		err := kc.AssignServerID(ctx, &mkkube.Client{KubeClient: clientset})
		duration := time.Since(start)
		if err == nil {
			t.Fatal("expected timeout error on stalled server")
		}
		if duration > 1*time.Second {
			t.Fatalf("AssignServerID hung for %v instead of stopping at deadline", duration)
		}
	})

	t.Run("AssignVersion", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		start := time.Now()
		err := kc.AssignVersion(ctx, &mkkube.Client{KubeClient: clientset})
		duration := time.Since(start)
		if err == nil {
			t.Fatal("expected timeout error on stalled server")
		}
		if duration > 1*time.Second {
			t.Fatalf("AssignVersion hung for %v instead of stopping at deadline", duration)
		}
	})
}

// TestAssignVersion_Success verifies that AssignVersion correctly parses a valid version
// response from the /version endpoint into kc.Version.
func TestAssignVersion_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/version" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(version.Info{GitVersion: "v1.28.2"})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: srv.URL})
	if err != nil {
		t.Fatalf("build clientset: %v", err)
	}
	kc := &K8sContext{Server: srv.URL}

	if err := kc.AssignVersion(context.Background(), &mkkube.Client{KubeClient: clientset}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kc.Version != "v1.28.2" {
		t.Fatalf("got version %q, want %q", kc.Version, "v1.28.2")
	}
}

// TestAssignVersion_NilRESTClient verifies that AssignVersion returns a clear error
// rather than panicking when the discovery REST client is nil.
func TestAssignVersion_NilRESTClient(t *testing.T) {
	kc := &K8sContext{Server: "http://127.0.0.1:1"}
	// handler with nil DiscoveryClient REST client returns error without panicking
	err := kc.AssignVersion(context.Background(), &mkkube.Client{KubeClient: &k8sclient.Clientset{}})
	if err == nil {
		t.Fatal("expected error with nil discovery REST client, got nil")
	}
}

// TestAssignVersion_MalformedJSON verifies that AssignVersion wraps JSON unmarshaling
// errors in ErrUnreachableKubeAPI.
func TestAssignVersion_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/version" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("not-valid-json"))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: srv.URL})
	if err != nil {
		t.Fatalf("build clientset: %v", err)
	}
	kc := &K8sContext{Server: srv.URL}

	err = kc.AssignVersion(context.Background(), &mkkube.Client{KubeClient: clientset})
	if err == nil {
		t.Fatal("expected error on malformed JSON, got nil")
	}
	if code := meshkiterrors.GetCode(err); code != ErrUnreachableKubeAPICode {
		t.Fatalf("expected ErrUnreachableKubeAPICode (%s), got: %v (code: %s)", ErrUnreachableKubeAPICode, err, code)
	}
}
