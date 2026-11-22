package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	mkkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/core"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func TestDiscoveryReleasesMachineAfterCanceledRequest(t *testing.T) {
	for _, stalledPath := range []string{"/api/v1/namespaces/kube-system", "/version"} {
		t.Run(stalledPath, func(t *testing.T) {
			started := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == stalledPath {
					close(started)
					<-r.Context().Done()
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"metadata": map[string]string{"uid": "30e099e0-ef31-477a-9f2d-754ee9b1d41a"}})
			}))
			defer server.Close()
			clientset, err := k8sclient.NewForConfig(&rest.Config{Host: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			log, err := logger.New("test", logger.Options{})
			if err != nil {
				t.Fatal(err)
			}
			machine, err := New(uuid.Must(uuid.NewV4()).String(), core.Uuid(uuid.Nil), log)
			if err != nil {
				t.Fatal(err)
			}
			machine.Context = &MachineCtx{K8sContext: models.K8sContext{Server: server.URL}, clientset: &mkkube.Client{KubeClient: clientset}, log: log}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan error, 1)
			go func() { _, err := machine.SendEvent(ctx, machines.Discovery, nil); done <- err }()
			select {
			case <-started:
			case <-time.After(5 * time.Second):
				t.Fatal("request did not start")
			}
			cancel()
			select {
			case err := <-done:
				if err != nil && err != context.Canceled {
					t.Logf("discovery returned: %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("discovery did not return")
			}
			// SendEvent takes the same mutex; an already canceled second call must return.
			again := make(chan struct{})
			go func() { _, _ = machine.SendEvent(ctx, machines.Discovery, nil); close(again) }()
			select {
			case <-again:
			case <-time.After(5 * time.Second):
				t.Fatal("state machine lock was not released")
			}
		})
	}
}
