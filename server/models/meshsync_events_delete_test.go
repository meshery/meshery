package models

import (
	"slices"
	"sync"
	"testing"

	meshkitBroker "github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/logger"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
)

// newMeshsyncEventsTestHandler returns a MeshsyncDataHandler over an in-memory
// database holding the MeshSync tables.
func newMeshsyncEventsTestHandler(t *testing.T) *MeshsyncDataHandler {
	t.Helper()
	db := newMeshSyncTestDB(t)
	db.Mutex = &sync.Mutex{}
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	return &MeshsyncDataHandler{dbHandler: *db, log: log}
}

// configMapObject is the ConfigMap MeshSync reports for name in namespace
// "default". Every field that lives in a child row (resourceVersion, the
// label, the annotation, spec and status) carries version, so a test can tell
// which incarnation of the object a row came from.
func configMapObject(name, version string) meshsyncmodel.KubernetesResource {
	return meshsyncmodel.KubernetesResource{
		APIVersion: "v1",
		Kind:       "ConfigMap",
		ClusterID:  "cluster-a",
		Data:       "data-" + version,
		KubernetesResourceMeta: &meshsyncmodel.KubernetesResourceObjectMeta{
			Name:            name,
			Namespace:       "default",
			ResourceVersion: version,
			Labels: []*meshsyncmodel.KubernetesKeyValue{
				{Kind: meshsyncmodel.KindLabel, Key: "app", Value: version},
			},
			Annotations: []*meshsyncmodel.KubernetesKeyValue{
				{Kind: meshsyncmodel.KindAnnotation, Key: "owner", Value: version},
			},
		},
		Spec:   &meshsyncmodel.KubernetesResourceSpec{Attribute: `{"version":"` + version + `"}`},
		Status: &meshsyncmodel.KubernetesResourceStatus{Attribute: `{"version":"` + version + `"}`},
	}
}

func configMapID(t *testing.T, name string) string {
	t.Helper()
	obj := configMapObject(name, "")
	meshsyncmodel.SetID(&obj)
	if obj.ID == "" {
		t.Fatalf("MeshSync derived no id for ConfigMap %s", name)
	}
	return obj.ID
}

func sendMeshsyncEvent(t *testing.T, mh *MeshsyncDataHandler, eventType meshkitBroker.EventType, obj meshsyncmodel.KubernetesResource) {
	t.Helper()
	if err := mh.meshsyncEventsAccumulator(&meshkitBroker.Message{EventType: eventType, Object: obj}); err != nil {
		t.Fatalf("%s event for %s: %v", eventType, obj.KubernetesResourceMeta.Name, err)
	}
}

// TestMeshsyncDeleteEventRemovesChildRows guards the DELETED event MeshSync
// publishes when an object is removed from the cluster: the resource and every
// row sharing its id (object meta, spec, status, labels and annotations) must
// be deleted, and other resources left alone. Before the fix only the
// kubernetes_resources row was deleted.
func TestMeshsyncDeleteEventRemovesChildRows(t *testing.T) {
	mh := newMeshsyncEventsTestHandler(t)
	deletedID, keptID := configMapID(t, "web"), configMapID(t, "api")

	sendMeshsyncEvent(t, mh, meshkitBroker.Add, configMapObject("web", "v1"))
	sendMeshsyncEvent(t, mh, meshkitBroker.Add, configMapObject("api", "v1"))
	sendMeshsyncEvent(t, mh, meshkitBroker.Delete, configMapObject("web", "v1"))

	assertMeshSyncRows(t, &mh.dbHandler, deletedID, 0)
	kept := countMeshSyncRows(t, &mh.dbHandler, keptID)
	want := map[string]int64{"resource": 1, "object meta": 1, "spec": 1, "status": 1, "key value": 2}
	for table, n := range want {
		if kept[table] != n {
			t.Errorf("%s rows of the resource that was not deleted: got %d, want %d", table, kept[table], n)
		}
	}
}

// TestMeshsyncRecreatedObjectKeepsOnlyItsOwnMetadata follows an object that is
// deleted and then created again with the same kind, namespace and name, and
// so the same MeshSync id. The new object must be stored with its own labels,
// annotations, object meta, spec and status. Before the fix the rows the
// deleted object left behind took precedence: creating the object did not
// overwrite them, so it showed the old spec and status and both objects'
// labels.
func TestMeshsyncRecreatedObjectKeepsOnlyItsOwnMetadata(t *testing.T) {
	mh := newMeshsyncEventsTestHandler(t)
	id := configMapID(t, "web")

	sendMeshsyncEvent(t, mh, meshkitBroker.Add, configMapObject("web", "v1"))
	sendMeshsyncEvent(t, mh, meshkitBroker.Delete, configMapObject("web", "v1"))
	sendMeshsyncEvent(t, mh, meshkitBroker.Add, configMapObject("web", "v2"))

	var keyValues []meshsyncmodel.KubernetesKeyValue
	if err := mh.dbHandler.Where("id = ?", id).Find(&keyValues).Error; err != nil {
		t.Fatalf("reading labels and annotations: %v", err)
	}
	got := make([]string, 0, len(keyValues))
	for _, kv := range keyValues {
		got = append(got, kv.Kind+" "+kv.Key+"="+kv.Value)
	}
	slices.Sort(got)
	if want := []string{"annotation owner=v2", "label app=v2"}; !slices.Equal(got, want) {
		t.Errorf("labels and annotations = %v, want %v", got, want)
	}

	var resource meshsyncmodel.KubernetesResource
	if err := mh.dbHandler.Preload("KubernetesResourceMeta").Preload("Spec").Preload("Status").
		Where("id = ?", id).First(&resource).Error; err != nil {
		t.Fatalf("reading the recreated object: %v", err)
	}
	if resource.KubernetesResourceMeta == nil || resource.KubernetesResourceMeta.ResourceVersion != "v2" {
		t.Errorf("object meta = %+v, want resourceVersion v2", resource.KubernetesResourceMeta)
	}
	if resource.Spec == nil || resource.Spec.Attribute != `{"version":"v2"}` {
		t.Errorf("spec = %+v, want the recreated object's spec", resource.Spec)
	}
	if resource.Status == nil || resource.Status.Attribute != `{"version":"v2"}` {
		t.Errorf("status = %+v, want the recreated object's status", resource.Status)
	}
}

// TestMeshsyncDeleteEventForUnknownObject checks that a DELETED event for an
// object that was never stored changes nothing and is not an error, as before.
func TestMeshsyncDeleteEventForUnknownObject(t *testing.T) {
	mh := newMeshsyncEventsTestHandler(t)
	keptID := configMapID(t, "api")
	sendMeshsyncEvent(t, mh, meshkitBroker.Add, configMapObject("api", "v1"))

	sendMeshsyncEvent(t, mh, meshkitBroker.Delete, configMapObject("web", "v1"))

	if n := countMeshSyncRows(t, &mh.dbHandler, keptID)["resource"]; n != 1 {
		t.Errorf("resource rows of the stored object: got %d, want 1", n)
	}
}

// TestMeshsyncDeleteEventWithoutIdentity checks that a DELETED event whose
// object carries neither an id nor the metadata to derive one is reported as
// an error rather than deleting anything.
func TestMeshsyncDeleteEventWithoutIdentity(t *testing.T) {
	mh := newMeshsyncEventsTestHandler(t)
	keptID := configMapID(t, "api")
	sendMeshsyncEvent(t, mh, meshkitBroker.Add, configMapObject("api", "v1"))

	err := mh.meshsyncEventsAccumulator(&meshkitBroker.Message{
		EventType: meshkitBroker.Delete,
		Object:    meshsyncmodel.KubernetesResource{APIVersion: "v1", Kind: "ConfigMap"},
	})
	if err == nil {
		t.Fatal("expected an error for a DELETED event without an object identity")
	}
	if n := countMeshSyncRows(t, &mh.dbHandler, keptID)["resource"]; n != 1 {
		t.Errorf("resource rows of the stored object: got %d, want 1", n)
	}
}
