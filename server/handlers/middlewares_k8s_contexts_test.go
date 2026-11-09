package handlers

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/sirupsen/logrus"
)

// loadContextsProvider stubs LoadAllK8sContext; every other Provider method
// is left nil and must not be reached by KubernetesMiddleware in these tests.
type loadContextsProvider struct {
	models.Provider
	contexts []*models.K8sContext
	err      error
}

func (p *loadContextsProvider) LoadAllK8sContext(string) ([]*models.K8sContext, error) {
	return p.contexts, p.err
}

func runKubernetesMiddleware(t *testing.T, prov models.Provider) (context.Context, string, error) {
	t.Helper()

	var out bytes.Buffer
	log, err := logger.New("test", logger.Options{Output: &out, ErrorOutput: &out, LogLevel: int(logrus.DebugLevel)})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	h := &Handler{log: log}
	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "test-token")
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	resCtx, err := KubernetesMiddleware(ctx, h, prov, user, []string{"all"})
	return resCtx, out.String(), err
}

// A failed context read must be logged with its cause, distinctly from an
// empty result, while the request still proceeds via kubeconfig discovery.
func TestKubernetesMiddleware_LoadContextsError(t *testing.T) {
	prov := &loadContextsProvider{err: errors.New("database is unavailable")}

	resCtx, logs, err := runKubernetesMiddleware(t, prov)

	if err != nil {
		t.Fatalf("expected the request to proceed, got error %v", err)
	}
	if resCtx == nil {
		t.Fatal("expected a context to be returned")
	}
	if !strings.Contains(logs, ErrGetK8sContextsCode) {
		t.Fatalf("expected error code %s in logs, got %s", ErrGetK8sContextsCode, logs)
	}
	if !strings.Contains(logs, "database is unavailable") {
		t.Fatalf("expected the read error's cause in logs, got %s", logs)
	}
	if strings.Contains(logs, ErrFailToGetK8SContextCode) {
		t.Fatalf("a read error must not be logged as an empty result, got %s", logs)
	}
}

func TestKubernetesMiddleware_NoContexts(t *testing.T) {
	resCtx, logs, err := runKubernetesMiddleware(t, &loadContextsProvider{})

	if err != nil {
		t.Fatalf("expected the request to proceed, got error %v", err)
	}
	if resCtx == nil {
		t.Fatal("expected a context to be returned")
	}
	if !strings.Contains(logs, ErrFailToGetK8SContextCode) {
		t.Fatalf("expected error code %s in logs, got %s", ErrFailToGetK8SContextCode, logs)
	}
	if strings.Contains(logs, ErrGetK8sContextsCode) {
		t.Fatalf("an empty result must not be logged as a read error, got %s", logs)
	}
}
