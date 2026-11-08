package resolver

import (
	"context"
	"testing"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
)

// stubTestProvider implements models.Provider for testing context resolution.
type stubTestProvider struct {
	models.Provider
}

// TestGetProvider verifies defensive context extraction of models.Provider.
func TestGetProvider(t *testing.T) {
	t.Run("nil context returns ErrInvalidRequest", func(t *testing.T) {
		//nolint:staticcheck // intentional nil context to test the getProvider guard
		prov, err := getProvider(nil)
		if err == nil {
			t.Fatalf("expected error for nil context, got nil")
		}
		if prov != nil {
			t.Fatalf("expected nil provider, got %v", prov)
		}
		if code := meshkiterrors.GetCode(err); code != ErrResolverInvalidRequestCode {
			t.Fatalf("expected code %s, got %s", ErrResolverInvalidRequestCode, code)
		}
	})

	t.Run("context without provider key returns ErrInvalidRequest", func(t *testing.T) {
		prov, err := getProvider(context.Background())
		if err == nil {
			t.Fatalf("expected error for missing provider key, got nil")
		}
		if prov != nil {
			t.Fatalf("expected nil provider, got %v", prov)
		}
		if code := meshkiterrors.GetCode(err); code != ErrResolverInvalidRequestCode {
			t.Fatalf("expected code %s, got %s", ErrResolverInvalidRequestCode, code)
		}
	})

	t.Run("context with unexpected type returns ErrInvalidRequest", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.ProviderCtxKey, "invalid-type")
		prov, err := getProvider(ctx)
		if err == nil {
			t.Fatalf("expected error for non-provider type, got nil")
		}
		if prov != nil {
			t.Fatalf("expected nil provider, got %v", prov)
		}
	})

	t.Run("context with valid provider returns provider successfully", func(t *testing.T) {
		expected := &stubTestProvider{}
		ctx := context.WithValue(context.Background(), models.ProviderCtxKey, expected)
		prov, err := getProvider(ctx)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if prov != expected {
			t.Fatalf("expected provider %v, got %v", expected, prov)
		}
	})
}

// TestMutationResolvers_NilGuards verifies that mutation resolvers guard against nil inputs
// and missing provider contexts without panicking.
func TestMutationResolvers_NilGuards(t *testing.T) {
	r := &Resolver{}
	mut := r.Mutation()
	ctx := context.Background()

	t.Run("ChangeOperatorStatus nil input returns ErrInvalidRequest", func(t *testing.T) {
		status, err := mut.ChangeOperatorStatus(ctx, nil)
		if err == nil {
			t.Fatalf("expected error for nil input, got nil")
		}
		if status != model.StatusUnknown {
			t.Fatalf("expected StatusUnknown, got %v", status)
		}
	})

	t.Run("ChangeOperatorStatus missing provider returns ErrInvalidRequest", func(t *testing.T) {
		status, err := mut.ChangeOperatorStatus(ctx, &model.OperatorStatusInput{
			TargetStatus: model.StatusEnabled,
			ContextID:    "ctx-1",
		})
		if err == nil {
			t.Fatalf("expected error for missing provider, got nil")
		}
		if status != model.StatusUnknown {
			t.Fatalf("expected StatusUnknown, got %v", status)
		}
	})

	t.Run("ChangeAdapterStatus nil input returns ErrInvalidRequest", func(t *testing.T) {
		status, err := mut.ChangeAdapterStatus(ctx, nil)
		if err == nil {
			t.Fatalf("expected error for nil input, got nil")
		}
		if status != model.StatusUnknown {
			t.Fatalf("expected StatusUnknown, got %v", status)
		}
	})

	t.Run("ChangeAdapterStatus missing provider returns ErrInvalidRequest", func(t *testing.T) {
		status, err := mut.ChangeAdapterStatus(ctx, &model.AdapterStatusInput{
			TargetStatus: model.StatusEnabled,
			Adapter:      "meshery-istio",
			TargetPort:   "10000",
		})
		if err == nil {
			t.Fatalf("expected error for missing provider, got nil")
		}
		if status != model.StatusUnknown {
			t.Fatalf("expected StatusUnknown, got %v", status)
		}
	})
}

// TestQueryResolvers_MissingProviderGuards verifies that all 14 query resolvers
// defensively reject contexts lacking ProviderCtxKey without crashing.
func TestQueryResolvers_MissingProviderGuards(t *testing.T) {
	r := &Resolver{}
	q := r.Query()
	ctx := context.Background()

	t.Run("ResyncCluster", func(t *testing.T) {
		status, err := q.ResyncCluster(ctx, nil, "ctx-1")
		if err == nil {
			t.Fatalf("expected error for missing provider, got nil")
		}
		if status != model.StatusUnknown {
			t.Fatalf("expected StatusUnknown, got %v", status)
		}
	})

	t.Run("GetAvailableAddons", func(t *testing.T) {
		res, err := q.GetAvailableAddons(ctx, nil)
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("GetControlPlanes", func(t *testing.T) {
		res, err := q.GetControlPlanes(ctx, nil)
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("GetDataPlanes", func(t *testing.T) {
		res, err := q.GetDataPlanes(ctx, nil)
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("GetAvailableNamespaces", func(t *testing.T) {
		res, err := q.GetAvailableNamespaces(ctx, []string{"cluster-1"})
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("GetPerfResult", func(t *testing.T) {
		res, err := q.GetPerfResult(ctx, "res-1")
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("FetchResults", func(t *testing.T) {
		res, err := q.FetchResults(ctx, model.PageFilter{}, "profile-1")
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("GetPerformanceProfiles", func(t *testing.T) {
		res, err := q.GetPerformanceProfiles(ctx, model.PageFilter{})
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("FetchAllResults", func(t *testing.T) {
		res, err := q.FetchAllResults(ctx, model.PageFilter{})
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("FetchPatterns", func(t *testing.T) {
		res, err := q.FetchPatterns(ctx, model.PageFilter{})
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("FetchPatternCatalogContent", func(t *testing.T) {
		res, err := q.FetchPatternCatalogContent(ctx, nil)
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("FetchFilterCatalogContent", func(t *testing.T) {
		res, err := q.FetchFilterCatalogContent(ctx, nil)
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("GetMeshModelSummary", func(t *testing.T) {
		res, err := q.GetMeshModelSummary(ctx, model.MeshModelSummarySelector{})
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})

	t.Run("FetchTelemetryComponents", func(t *testing.T) {
		res, err := q.FetchTelemetryComponents(ctx, []string{"ctx-1"})
		if err == nil || res != nil {
			t.Fatalf("expected error and nil result, got res=%v, err=%v", res, err)
		}
	})
}

// TestResyncCluster_EdgeCases verifies edge-case parameter validation in resyncCluster.
func TestResyncCluster_EdgeCases(t *testing.T) {
	r := &Resolver{}
	prov := &stubTestProvider{}
	ctx := context.Background()

	t.Run("nil actions returns ErrInvalidRequest", func(t *testing.T) {
		status, err := r.resyncCluster(ctx, prov, nil, "ctx-1")
		if err == nil {
			t.Fatalf("expected error for nil actions, got nil")
		}
		if status != model.StatusUnknown {
			t.Fatalf("expected StatusUnknown, got %v", status)
		}
	})

	t.Run("ClearDb missing clusters in context returns ErrEmptyCurrentK8sContext", func(t *testing.T) {
		actions := &model.ReSyncActions{
			ClearDb: "true",
		}
		status, err := r.resyncCluster(ctx, prov, actions, "unmatched-id")
		if err == nil {
			t.Fatalf("expected error when context has no clusters, got nil")
		}
		if status != "" {
			t.Fatalf("expected empty status, got %v", status)
		}
		if code := meshkiterrors.GetCode(err); code != ErrEmptyCurrentK8sContextCode {
			t.Fatalf("expected code %s, got %s", ErrEmptyCurrentK8sContextCode, code)
		}
	})

	t.Run("ClearDb with unmatched context ID returns ErrEmptyCurrentK8sContext", func(t *testing.T) {
		actions := &model.ReSyncActions{
			ClearDb: "true",
		}
		clusterCtx := context.WithValue(ctx, models.AllKubeClusterKey, []*models.K8sContext{
			{ID: "connected-ctx-1"},
		})
		status, err := r.resyncCluster(clusterCtx, prov, actions, "different-ctx-id")
		if err == nil {
			t.Fatalf("expected error for unmatched context ID, got nil")
		}
		if status != "" {
			t.Fatalf("expected empty status, got %v", status)
		}
		if code := meshkiterrors.GetCode(err); code != ErrEmptyCurrentK8sContextCode {
			t.Fatalf("expected code %s, got %s", ErrEmptyCurrentK8sContextCode, code)
		}
	})
}

// TestChangeAdapterStatus_NilGuards verifies nil guards on adapter change.
func TestChangeAdapterStatus_NilGuards(t *testing.T) {
	r := &Resolver{}
	prov := &stubTestProvider{}
	ctx := context.Background()

	t.Run("empty name and port returns ErrAdapterInsufficientInformation", func(t *testing.T) {
		status, err := r.changeAdapterStatus(ctx, prov, model.StatusEnabled, "", "")
		if err == nil {
			t.Fatalf("expected error for empty name/port, got nil")
		}
		if status != model.StatusUnknown {
			t.Fatalf("expected StatusUnknown, got %v", status)
		}
		if code := meshkiterrors.GetCode(err); code != ErrAdapterInsufficientInformationCode {
			t.Fatalf("expected code %s, got %s", ErrAdapterInsufficientInformationCode, code)
		}
	})

	t.Run("missing AdapterTracker in Config returns ErrAdapterInsufficientInformation", func(t *testing.T) {
		status, err := r.changeAdapterStatus(ctx, prov, model.StatusEnabled, "meshery-istio", "10000")
		if err == nil {
			t.Fatalf("expected error for uninitialized AdapterTracker, got nil")
		}
		if status != model.StatusUnknown {
			t.Fatalf("expected StatusUnknown, got %v", status)
		}
	})
}

// TestChangeOperatorStatus_NilGuards verifies nil and missing context guards in changeOperatorStatus.
func TestChangeOperatorStatus_NilGuards(t *testing.T) {
	r := &Resolver{}
	prov := &stubTestProvider{}
	ctx := context.Background()

	t.Run("missing KubeClustersKey with empty ctxID returns ErrMesheryClientNil", func(t *testing.T) {
		status, err := r.changeOperatorStatus(ctx, prov, model.StatusEnabled, "")
		if err == nil {
			t.Fatalf("expected error for missing cluster context, got nil")
		}
		if status != model.StatusUnknown {
			t.Fatalf("expected StatusUnknown, got %v", status)
		}
		if err != model.ErrMesheryClientNil {
			t.Fatalf("expected ErrMesheryClientNil, got %v", err)
		}
	})

	t.Run("missing AllKubeClusterKey with non-empty ctxID returns ErrNilClient", func(t *testing.T) {
		status, err := r.changeOperatorStatus(ctx, prov, model.StatusEnabled, "ctx-123")
		if err == nil {
			t.Fatalf("expected error for missing AllKubeClusterKey, got nil")
		}
		if status != model.StatusUnknown {
			t.Fatalf("expected StatusUnknown, got %v", status)
		}
		if code := meshkiterrors.GetCode(err); code != ErrResolverNilClientCode {
			t.Fatalf("expected code %s, got %s", ErrResolverNilClientCode, code)
		}
	})
}
