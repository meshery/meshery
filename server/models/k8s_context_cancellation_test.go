package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mkkube "github.com/meshery/meshkit/utils/kubernetes"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func TestK8sDiscoveryRequestsHonorCancellation(t *testing.T) {
	for _, tc := range []struct {
		name string
		path string
		call func(context.Context, *K8sContext, *mkkube.Client) error
	}{
		{"server ID", "/api/v1/namespaces/kube-system", func(ctx context.Context, kc *K8sContext, client *mkkube.Client) error {
			return kc.AssignServerID(ctx, client)
		}},
		{"version", "/version", func(ctx context.Context, kc *K8sContext, client *mkkube.Client) error {
			return kc.AssignVersion(ctx, client)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			started := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.path {
					http.NotFound(w, r)
					return
				}
				close(started)
				<-r.Context().Done()
			}))
			defer server.Close()
			clientset, err := k8sclient.NewForConfig(&rest.Config{Host: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- tc.call(ctx, &K8sContext{Server: server.URL}, &mkkube.Client{KubeClient: clientset}) }()
			select {
			case <-started:
			case <-time.After(5 * time.Second):
				t.Fatal("request did not start")
			}
			cancel()
			select {
			case err := <-done:
				if err == nil {
					t.Fatal("canceled request succeeded")
				}
			case <-time.After(5 * time.Second):
				t.Fatal("request did not stop after cancellation")
			}
		})
	}
}

func TestAssignVersionWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/version" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"gitVersion":"v1.32.1"}`))
	}))
	defer server.Close()
	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	kc := &K8sContext{Server: server.URL}
	if err := kc.AssignVersion(context.Background(), &mkkube.Client{KubeClient: clientset}); err != nil {
		t.Fatal(err)
	}
	if kc.Version != "v1.32.1" {
		t.Fatalf("version = %q, want v1.32.1", kc.Version)
	}
}

// TestK8sContextsFromKubeconfigWithContext_CancelsOnContextDone verifies that
// when the caller's context is canceled, K8sContextsFromKubeconfigWithContext
// bails out immediately without processing remaining contexts.
func TestK8sContextsFromKubeconfigWithContext_CancelsOnContextDone(t *testing.T) {
	kubeconfig := []byte(`apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://127.0.0.1:59901
  name: cluster-1
- cluster:
    server: https://127.0.0.1:59902
  name: cluster-2
contexts:
- context:
    cluster: cluster-1
    user: user-1
  name: ctx-1
- context:
    cluster: cluster-2
    user: user-2
  name: ctx-2
users:
- name: user-1
  user:
    token: tok-1
- name: user-2
  user:
    token: tok-2
`)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel context

	eventMetadata := map[string]interface{}{}
	contexts := K8sContextsFromKubeconfigWithContext(ctx, nil, "user-id", nil, kubeconfig, nil, eventMetadata, nil, true)
	if len(contexts) != 0 {
		t.Fatalf("expected 0 contexts on canceled context, got %d", len(contexts))
	}
}

// TestAssignServerID_TimeoutOnHangingServer verifies that AssignServerID
// returns an ErrUnreachableKubeAPI error when the cluster API server hangs.
func TestAssignServerID_TimeoutOnHangingServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Hang until client aborts
		<-r.Context().Done()
	}))
	defer server.Close()

	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: server.URL})
	if err != nil {
		t.Fatal(err)
	}

	kc := &K8sContext{Server: server.URL}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = kc.AssignServerID(ctx, &mkkube.Client{KubeClient: clientset})
	if err == nil {
		t.Fatal("expected AssignServerID to fail on timeout, got nil")
	}
}

// TestAssignVersion_TimeoutOnHangingServer verifies that AssignVersion
// returns an ErrUnreachableKubeAPI error when the cluster API server hangs.
func TestAssignVersion_TimeoutOnHangingServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Hang until client aborts
		<-r.Context().Done()
	}))
	defer server.Close()

	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: server.URL})
	if err != nil {
		t.Fatal(err)
	}

	kc := &K8sContext{Server: server.URL}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = kc.AssignVersion(ctx, &mkkube.Client{KubeClient: clientset})
	if err == nil {
		t.Fatal("expected AssignVersion to fail on timeout, got nil")
	}
}
