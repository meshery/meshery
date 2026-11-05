package models

import (
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
)

// K8sContextFromConnection reads the cluster and auth blocks straight out of the
// stored credential, whichever shape that credential was persisted in. The ~36.9k
// Kubernetes credentials in production hold both blocks at the top of the secret
// map; the credential form instead writes a double-nested wrapper whose inner
// payload names the cluster with clusterName/clusterServerURL. The shape
// catalogue lives in credential_secret.go.
func TestK8sContextFromConnectionCredentialShapes(t *testing.T) {
	storedAuth := map[string]interface{}{"clusterToken": "tok", "clusterUserName": "u1"}
	storedCluster := map[string]interface{}{"server": "https://k8s.example", "name": "cluster-a"}

	// Exactly the auth block ui/components/schemas/credentials/kubernetes.tsx
	// renders and MesheryCredentialComponent persists.
	formAuth := map[string]interface{}{
		"clusterUserName":                 "u1",
		"clusterToken":                    "tok",
		"clusterClientCertificateData":    "cert",
		"clusterClientKeyData":            "key",
		"clusterCertificateAuthorityData": "ca",
	}

	tests := []struct {
		name        string
		secret      map[string]interface{}
		wantAuth    map[string]interface{}
		wantCluster map[string]interface{}
	}{
		{
			name:        "stored kubernetes shape",
			secret:      map[string]interface{}{"auth": storedAuth, "cluster": storedCluster},
			wantAuth:    storedAuth,
			wantCluster: storedCluster,
		},
		{
			name: "credential form double-nested shape",
			secret: map[string]interface{}{
				"credentialName": "kube-cred",
				"secret": map[string]interface{}{
					"clusterName":      "cluster-a",
					"clusterServerURL": "https://k8s.example",
					"auth":             formAuth,
				},
			},
			wantAuth: formAuth,
			// Cluster is nil here, and that is the tolerant read working rather
			// than failing: the wrapper is unwrapped correctly and the auth block
			// arrives intact, but this payload describes the cluster with
			// clusterName/clusterServerURL instead of the kubeconfig-style
			// `cluster` block K8sContext reads. Pinned deliberately - deciding
			// which side moves is tracked in
			// https://github.com/meshery/meshery/issues/21336. Update this
			// expectation only alongside that decision.
			wantCluster: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
			if err != nil {
				t.Fatalf("open database: %v", err)
			}
			if err := db.AutoMigrate(Credential{}); err != nil {
				t.Fatalf("migrate tables: %v", err)
			}
			provider := &DefaultLocalProvider{GenericPersister: &db}

			credential, err := provider.SaveUserCredential("", &Credential{
				Name:   "kube-cred",
				Type:   "kubernetes",
				Secret: tt.secret,
			})
			if err != nil {
				t.Fatalf("save credential: %v", err)
			}

			credentialID := credential.ID
			ctx, err := K8sContextFromConnection(provider, "", &connections.Connection{
				ID:           uuid.Must(uuid.NewV4()),
				Name:         "cluster-a",
				Kind:         "kubernetes",
				CredentialID: &credentialID,
				Metadata:     core.Map{"name": "cluster-a"},
			})
			if err != nil {
				t.Fatalf("K8sContextFromConnection: %v", err)
			}

			// ctx.Auth/ctx.Cluster are sql.Map, so compare the contents rather
			// than the named map type the assignment converts to.
			if !reflect.DeepEqual(map[string]interface{}(ctx.Auth), tt.wantAuth) {
				t.Errorf("Auth = %#v, want %#v", ctx.Auth, tt.wantAuth)
			}
			if !reflect.DeepEqual(map[string]interface{}(ctx.Cluster), tt.wantCluster) {
				t.Errorf("Cluster = %#v, want %#v", ctx.Cluster, tt.wantCluster)
			}
		})
	}
}
