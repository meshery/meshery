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
// stored credential. The ~36.9k Kubernetes credentials in production hold them at
// the top of the secret map, while the credential form writes the same payload
// one level down under a `secret` key - both have to resolve. The shape
// catalogue lives in credential_secret.go.
func TestK8sContextFromConnectionCredentialShapes(t *testing.T) {
	auth := map[string]interface{}{"clusterToken": "tok", "clusterUserName": "u1"}
	cluster := map[string]interface{}{"server": "https://k8s.example", "name": "cluster-a"}

	tests := []struct {
		name   string
		secret map[string]interface{}
	}{
		{
			name:   "stored kubernetes shape",
			secret: map[string]interface{}{"auth": auth, "cluster": cluster},
		},
		{
			name: "legacy double-nested shape",
			secret: map[string]interface{}{
				"credentialName": "kube-cred",
				"secret":         map[string]interface{}{"auth": auth, "cluster": cluster},
			},
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

			// ctx.Auth/ctx.Cluster are core.Map, so compare the contents rather
			// than the named map type the assignment converts to.
			if !reflect.DeepEqual(map[string]interface{}(ctx.Auth), auth) {
				t.Errorf("Auth = %#v, want %#v", ctx.Auth, auth)
			}
			if !reflect.DeepEqual(map[string]interface{}(ctx.Cluster), cluster) {
				t.Errorf("Cluster = %#v, want %#v", ctx.Cluster, cluster)
			}
		})
	}
}
