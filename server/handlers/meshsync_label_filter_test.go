package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"sort"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshsync/pkg/model"
	"github.com/stretchr/testify/require"
)

// newLabelFilterFixture stores two Pods of cluster c1 that share app=web and
// differ in tier, the way MeshSync stores them: each label is a
// kubernetes_key_values row keyed by the object's metadata ID.
func newLabelFilterFixture(t *testing.T) (*Handler, models.Provider) {
	t.Helper()
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "meshsync.db")})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.KubernetesResource{}, &model.KubernetesResourceObjectMeta{},
		&model.KubernetesKeyValue{}, &model.KubernetesResourceSpec{}, &model.KubernetesResourceStatus{}))

	for _, pod := range []struct{ id, name, tier string }{
		{"pod-frontend", "web-frontend", "frontend"},
		{"pod-backend", "web-backend", "backend"},
	} {
		resource := model.KubernetesResource{
			ID: pod.id, APIVersion: "v1", Kind: "Pod", Model: "kubernetes", ClusterID: "c1",
			KubernetesResourceMeta: &model.KubernetesResourceObjectMeta{
				ID: pod.id, Name: pod.name, Namespace: "default", ClusterID: "c1",
				Labels: []*model.KubernetesKeyValue{
					{ID: pod.id, Kind: model.KindLabel, Key: "app", Value: "web"},
					{ID: pod.id, Kind: model.KindLabel, Key: "tier", Value: pod.tier},
				},
			},
		}
		require.NoError(t, db.Create(&resource).Error)
	}
	return &Handler{log: newTestLogger(t)}, &models.DefaultLocalProvider{GenericPersister: &db}
}

// listByLabels sends GET /api/system/meshsync/resources for cluster c1 with the
// given labels and extra query values.
func listByLabels(t *testing.T, h *Handler, provider models.Provider, labels []string, extra url.Values) (*httptest.ResponseRecorder, models.MeshSyncResourcesAPIResponse) {
	t.Helper()
	q := url.Values{"clusterIds": {`["c1"]`}, "label": labels}
	for key, values := range extra {
		q[key] = values
	}
	req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	h.GetMeshSyncResources(rec, req, nil, nil, provider)

	var page models.MeshSyncResourcesAPIResponse
	if rec.Code == http.StatusOK {
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &page))
	}
	return rec, page
}

func resourceNames(page models.MeshSyncResourcesAPIResponse) []string {
	names := []string{}
	for _, resource := range page.Resources {
		if resource.KubernetesResourceMeta != nil {
			names = append(names, resource.KubernetesResourceMeta.Name)
		}
	}
	sort.Strings(names)
	return names
}

// Several labels select the resources that carry all of them, each once.
func TestGetMeshSyncResourcesMatchesAllLabels(t *testing.T) {
	h, provider := newLabelFilterFixture(t)

	rec, page := listByLabels(t, h, provider, []string{"app=web", "tier=frontend"}, nil)

	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.Equal(t, []string{"web-frontend"}, resourceNames(page))
	require.EqualValues(t, 1, page.TotalCount)
}

// One label still selects every resource that carries it.
func TestGetMeshSyncResourcesMatchesOneLabel(t *testing.T) {
	h, provider := newLabelFilterFixture(t)

	rec, page := listByLabels(t, h, provider, []string{"app=web"}, nil)

	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.Equal(t, []string{"web-backend", "web-frontend"}, resourceNames(page))
	require.EqualValues(t, 2, page.TotalCount)
}

// Sorting by kind, an allowed sort column, works together with a label filter.
func TestGetMeshSyncResourcesLabelFilterSortedByKind(t *testing.T) {
	h, provider := newLabelFilterFixture(t)

	rec, page := listByLabels(t, h, provider, []string{"app=web"}, url.Values{"order": {"kind asc"}})

	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.EqualValues(t, 2, page.TotalCount)
}
