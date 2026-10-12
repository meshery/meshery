package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sort"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshsync/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMeshSyncClusterIDsFromQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    []string
		wantErr bool
	}{
		{
			name:  "given repeated clusterId when parsed then every value is returned",
			query: "clusterId=c1&clusterId=c2",
			want:  []string{"c1", "c2"},
		},
		{
			name:  "given legacy JSON clusterIds when parsed then every value is returned",
			query: "clusterIds=" + url.QueryEscape(`["c1","c2"]`),
			want:  []string{"c1", "c2"},
		},
		{
			name:  "given both forms when parsed then values are merged",
			query: "clusterId=c1&clusterIds=" + url.QueryEscape(`["c3"]`),
			want:  []string{"c1", "c3"},
		},
		{
			name:  "given an empty clusterId value when parsed then it is skipped",
			query: "clusterId=&clusterId=c2",
			want:  []string{"c2"},
		},
		{
			name:  "given no cluster parameter when parsed then the result is empty, not nil",
			query: "",
			want:  []string{},
		},
		{
			name:    "given malformed legacy clusterIds when parsed then an error is returned",
			query:   "clusterIds=not-json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := url.ParseQuery(tt.query)
			if err != nil {
				t.Fatalf("parsing test query %q: %v", tt.query, err)
			}

			got, err := meshSyncClusterIDsFromQuery(values)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %#v, got %#v", tt.want, got)
			}
		})
	}
}

// newMeshSyncTestProvider returns a local provider backed by an in-memory
// database holding the MeshSync tables, seeded with the given resources.
func newMeshSyncTestProvider(t *testing.T, resources ...*model.KubernetesResource) models.Provider {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}
	if err := db.AutoMigrate(
		&model.KubernetesKeyValue{},
		&model.KubernetesResource{},
		&model.KubernetesResourceSpec{},
		&model.KubernetesResourceStatus{},
		&model.KubernetesResourceObjectMeta{},
	); err != nil {
		t.Fatalf("failed to migrate meshsync tables: %v", err)
	}
	for _, resource := range resources {
		if err := db.Create(resource).Error; err != nil {
			t.Fatalf("failed to seed resource %s: %v", resource.ID, err)
		}
	}

	provider := &models.DefaultLocalProvider{}
	provider.Initialize()
	provider.GenericPersister = &database.Handler{DB: db}
	return provider
}

func newMeshSyncTestResource(id, clusterID string, keyValues ...*model.KubernetesKeyValue) *model.KubernetesResource {
	meta := &model.KubernetesResourceObjectMeta{
		ID:        id,
		Name:      id,
		Namespace: "default",
		UID:       id,
		ClusterID: clusterID,
	}
	for _, kv := range keyValues {
		kv.ID = id
		kv.UniqueID = id + "/" + kv.Kind + "/" + kv.Key
		if kv.Kind == model.KindAnnotation {
			meta.Annotations = append(meta.Annotations, kv)
		} else {
			meta.Labels = append(meta.Labels, kv)
		}
	}
	return &model.KubernetesResource{
		ID:                     id,
		APIVersion:             "v1",
		Kind:                   "Pod",
		Model:                  "kubernetes",
		ClusterID:              clusterID,
		KubernetesResourceMeta: meta,
	}
}

func TestGetMeshSyncResourcesScopesByCluster(t *testing.T) {
	handler := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newMeshSyncTestProvider(t,
		newMeshSyncTestResource("r1", "c1"),
		newMeshSyncTestResource("r2", "c2"),
		newMeshSyncTestResource("r3", "c3"),
	)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedIDs    []string
	}{
		{
			name:           "given repeated clusterId when listing then only those clusters are returned",
			query:          "clusterId=c1&clusterId=c2",
			expectedStatus: http.StatusOK,
			expectedIDs:    []string{"r1", "r2"},
		},
		{
			name:           "given legacy JSON clusterIds when listing then those clusters are still returned",
			query:          "clusterIds=" + url.QueryEscape(`["c3"]`),
			expectedStatus: http.StatusOK,
			expectedIDs:    []string{"r3"},
		},
		{
			name:           "given both forms when listing then the union is returned",
			query:          "clusterId=c1&clusterIds=" + url.QueryEscape(`["c3"]`),
			expectedStatus: http.StatusOK,
			expectedIDs:    []string{"r1", "r3"},
		},
		{
			name:           "given no cluster parameter when listing then no resources are returned",
			query:          "",
			expectedStatus: http.StatusOK,
			expectedIDs:    []string{},
		},
		{
			name:           "given malformed legacy clusterIds when listing then status 400",
			query:          "clusterIds=not-json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources?"+tt.query, nil)
			rw := httptest.NewRecorder()

			handler.GetMeshSyncResources(rw, req, nil, nil, provider)

			if rw.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d: %s", tt.expectedStatus, rw.Code, rw.Body.String())
			}
			if tt.expectedStatus != http.StatusOK {
				return
			}

			// meshsync's create hook rewrites the resource id into a derived key,
			// so identify the seeded resources by metadata.name instead.
			var response struct {
				TotalCount int64 `json:"totalCount"`
				Resources  []struct {
					Metadata struct {
						Name string `json:"name"`
					} `json:"metadata"`
				} `json:"resources"`
			}
			if err := json.NewDecoder(rw.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			gotIDs := make([]string, 0, len(response.Resources))
			for _, resource := range response.Resources {
				gotIDs = append(gotIDs, resource.Metadata.Name)
			}
			sort.Strings(gotIDs)

			if !reflect.DeepEqual(gotIDs, tt.expectedIDs) {
				t.Fatalf("expected resources %v, got %v", tt.expectedIDs, gotIDs)
			}
			if response.TotalCount != int64(len(tt.expectedIDs)) {
				t.Fatalf("expected totalCount %d, got %d", len(tt.expectedIDs), response.TotalCount)
			}
		})
	}
}

// The summary's labels are the distinct key/value pairs in scope, and that is
// all the query selects. Pin the wire shape to exactly {key, value}: the
// meshsync KubernetesKeyValue it used to be encoded as sent always-empty
// id/unique_id/kind fields and dropped `value` for empty-valued labels such as
// node-role.kubernetes.io/control-plane, which every control-plane node carries.
func TestGetMeshSyncResourcesSummaryLabelsAreKeyValuePairs(t *testing.T) {
	handler := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newMeshSyncTestProvider(t,
		newMeshSyncTestResource("r1", "c1",
			&model.KubernetesKeyValue{Kind: model.KindLabel, Key: "app", Value: "nginx"},
			&model.KubernetesKeyValue{Kind: model.KindLabel, Key: "node-role.kubernetes.io/control-plane", Value: ""},
			&model.KubernetesKeyValue{Kind: model.KindAnnotation, Key: "note", Value: "not-a-label"},
		),
		newMeshSyncTestResource("r2", "c2",
			&model.KubernetesKeyValue{Kind: model.KindLabel, Key: "other-cluster", Value: "yes"},
		),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/summary?clusterId=c1", nil)
	rw := httptest.NewRecorder()

	handler.GetMeshSyncResourcesSummary(rw, req, nil, nil, provider)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rw.Code, rw.Body.String())
	}

	var response struct {
		Labels []map[string]any `json:"labels"`
	}
	if err := json.NewDecoder(rw.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	got := map[string]any{}
	for _, label := range response.Labels {
		keys := make([]string, 0, len(label))
		for k := range label {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		if !reflect.DeepEqual(keys, []string{"key", "value"}) {
			t.Fatalf("expected each label to carry exactly key and value, got %v", label)
		}
		got[label["key"].(string)] = label["value"]
	}

	want := map[string]any{
		"app":                                   "nginx",
		"node-role.kubernetes.io/control-plane": "",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected labels %v, got %v", want, got)
	}
}
