package resolver

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestResolver(t *testing.T) *Resolver {
	t.Helper()
	log, err := logger.New("test", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: int(logrus.ErrorLevel),
		Output:   io.Discard,
	})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	return &Resolver{Log: log}
}

func k8sContextForServer(id, server string) *models.K8sContext {
	return &models.K8sContext{
		ID:   id,
		Name: id,
		Cluster: map[string]interface{}{
			"name":    id,
			"cluster": map[string]interface{}{"server": server},
		},
		Auth: map[string]interface{}{
			"name": id + "-user",
			"user": map[string]interface{}{"token": "test-token"},
		},
		Server: server,
	}
}

func TestGetKubectlDescribe_NoK8sContextsInRequest(t *testing.T) {
	r := newTestResolver(t)

	details, err := r.getKubectlDescribe(context.Background(), "nginx", "pod", "default", "ctx-1")

	assert.Nil(t, details)
	assert.ErrorIs(t, err, ErrEmptyCurrentK8sContext)
}

func TestGetKubectlDescribe_UnknownK8sContextID(t *testing.T) {
	r := newTestResolver(t)
	ctx := context.WithValue(context.Background(), models.AllKubeClusterKey, []*models.K8sContext{
		nil,
		k8sContextForServer("ctx-1", "https://127.0.0.1:6443"),
	})

	details, err := r.getKubectlDescribe(ctx, "nginx", "pod", "default", "ctx-2")

	assert.Nil(t, details)
	assert.ErrorIs(t, err, ErrEmptyCurrentK8sContext)
}

func TestGetKubectlDescribe_UsesSelectedK8sContext(t *testing.T) {
	newCluster := func() (*httptest.Server, *[]string) {
		var mu sync.Mutex
		paths := []string{}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			mu.Lock()
			paths = append(paths, req.URL.Path)
			mu.Unlock()
			if req.URL.Path != "/api/v1/namespaces/default/pods/nginx" {
				http.NotFound(w, req)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"apiVersion":"v1","kind":"Pod","metadata":{"name":"nginx","namespace":"default"}}`)
		}))
		return srv, &paths
	}

	selected, selectedPaths := newCluster()
	defer selected.Close()
	other, otherPaths := newCluster()
	defer other.Close()

	r := newTestResolver(t)
	ctx := context.WithValue(context.Background(), models.AllKubeClusterKey, []*models.K8sContext{
		k8sContextForServer("ctx-other", other.URL),
		k8sContextForServer("ctx-selected", selected.URL),
	})

	details, err := r.getKubectlDescribe(ctx, "nginx", "pod", "default", "ctx-selected")

	require.NoError(t, err)
	require.NotNil(t, details)
	require.NotNil(t, details.Describe)
	assert.Contains(t, *details.Describe, "nginx")
	assert.Contains(t, *selectedPaths, "/api/v1/namespaces/default/pods/nginx")
	assert.Empty(t, *otherPaths)
}
