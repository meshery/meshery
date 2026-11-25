package models

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/viper"
)

// setKubeconfigFolder points the server's own kubeconfig lookup at folder for
// the test and restores whatever was set before.
func setKubeconfigFolder(t *testing.T, folder string) {
	t.Helper()
	prev := viper.Get("KUBECONFIG_FOLDER")
	viper.Set("KUBECONFIG_FOLDER", folder)
	t.Cleanup(func() { viper.Set("KUBECONFIG_FOLDER", prev) })
}

func kubeCredential(user map[string]interface{}, cluster map[string]interface{}) (auth, clusterBlock map[string]interface{}) {
	return map[string]interface{}{"name": "u", "user": user},
		map[string]interface{}{"name": "c", "cluster": cluster}
}

func TestUnsafeKubeCredentialFields(t *testing.T) {
	execPlugin := map[string]interface{}{"apiVersion": "client.authentication.k8s.io/v1", "command": "sh", "args": []interface{}{"-c", "id"}}
	tests := []struct {
		name    string
		user    map[string]interface{}
		cluster map[string]interface{}
		want    []string
	}{
		{"embedded data only", map[string]interface{}{"client-certificate-data": "YQ==", "client-key-data": "Yg==", "token": "t"}, map[string]interface{}{"server": "https://k8s", "certificate-authority-data": "Yw=="}, nil},
		{"exec plugin", map[string]interface{}{"exec": execPlugin}, map[string]interface{}{"server": "https://k8s"}, []string{"exec"}},
		{"auth-provider", map[string]interface{}{"auth-provider": map[string]interface{}{"name": "oidc"}}, map[string]interface{}{"server": "https://k8s"}, []string{"auth-provider"}},
		{"token file", map[string]interface{}{"tokenFile": "/var/run/secrets/kubernetes.io/serviceaccount/token"}, map[string]interface{}{"server": "https://attacker"}, []string{"tokenFile"}},
		{"key and cert paths", map[string]interface{}{"client-key": "/k", "client-certificate": "/c"}, map[string]interface{}{"server": "https://k8s"}, []string{"client-certificate", "client-key"}},
		{"ca path", map[string]interface{}{"token": "t"}, map[string]interface{}{"server": "https://k8s", "certificate-authority": "/ca"}, []string{"certificate-authority"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, cluster := kubeCredential(tt.user, tt.cluster)
			got := UnsafeKubeCredentialFields(auth, cluster)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("UnsafeKubeCredentialFields = %v, want %v", got, tt.want)
			}
		})
	}
	if got := UnsafeKubeCredentialFields(nil, nil); got != nil {
		t.Fatalf("empty credential = %v, want nil", got)
	}
}

func saveKubeCredentialAndConnection(t *testing.T, secret map[string]interface{}) (*DefaultLocalProvider, *connections.Connection) {
	t.Helper()
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(Credential{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}
	provider := &DefaultLocalProvider{GenericPersister: &db}
	credential, err := provider.SaveUserCredential("", &Credential{Name: "kube-cred", Type: "kubernetes", Secret: secret})
	if err != nil {
		t.Fatalf("save credential: %v", err)
	}
	id := credential.ID
	return provider, &connections.Connection{ID: uuid.Must(uuid.NewV4()), Name: "cluster-a", Kind: "kubernetes", CredentialID: &id, Metadata: core.Map{"name": "cluster-a"}}
}

// TestK8sContextFromConnectionRefusesUnsafeCredential pins issue #24218: a saved
// connection's credential secret can be edited to carry an exec plugin or a file
// path, and ping, the state machine and graphql then build a client from it.
// Every one of them goes through K8sContextFromConnection first, so it must
// refuse, and it must do so without echoing the command or the path.
func TestK8sContextFromConnectionRefusesUnsafeCredential(t *testing.T) {
	setKubeconfigFolder(t, t.TempDir()) // no kubeconfig of its own inside

	for name, user := range map[string]map[string]interface{}{
		"exec":       {"exec": map[string]interface{}{"apiVersion": "client.authentication.k8s.io/v1", "command": "touch-marker-cmd", "args": []interface{}{"/tmp/marker-path"}}},
		"tokenFile":  {"tokenFile": "/tmp/marker-path"},
		"client-key": {"client-key": "/tmp/marker-path"},
	} {
		t.Run(name, func(t *testing.T) {
			auth, cluster := kubeCredential(user, map[string]interface{}{"server": "https://127.0.0.1:1"})
			provider, conn := saveKubeCredentialAndConnection(t, map[string]interface{}{"auth": auth, "cluster": cluster})

			_, err := K8sContextFromConnection(provider, "", conn)
			if meshkiterrors.GetCode(err) != ErrUnsafeKubeCredentialCode {
				t.Fatalf("err = %v, want ErrUnsafeKubeCredential", err)
			}
			for _, leak := range []string{"marker-path", "touch-marker-cmd"} {
				if strings.Contains(err.Error(), leak) || strings.Contains(meshkiterrors.GetLDescription(err), leak) {
					t.Fatalf("error leaks %q: %v", leak, err)
				}
			}
		})
	}

	t.Run("embedded data still works", func(t *testing.T) {
		auth, cluster := kubeCredential(map[string]interface{}{"token": "t"}, map[string]interface{}{"server": "https://127.0.0.1:1"})
		provider, conn := saveKubeCredentialAndConnection(t, map[string]interface{}{"auth": auth, "cluster": cluster})
		if _, err := K8sContextFromConnection(provider, "", conn); err != nil {
			t.Fatalf("K8sContextFromConnection: %v", err)
		}
	})
}

// The server's own kubeconfig legitimately uses exec plugins (EKS, GKE), and
// contexts imported from it are stored as connections and rebuilt through
// K8sContextFromConnection. Those must keep working, but only when the stored
// exec is identical to what the operator's kubeconfig holds.
func TestK8sContextFromConnectionAllowsTheServersOwnExec(t *testing.T) {
	folder := t.TempDir()
	if err := os.WriteFile(filepath.Join(folder, "config"), []byte(`apiVersion: v1
kind: Config
clusters:
- name: c
  cluster:
    server: https://127.0.0.1:1
users:
- name: u
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1
      command: aws
      args: [eks, get-token]
contexts:
- name: c
  context: {cluster: c, user: u}
current-context: c
`), 0o600); err != nil {
		t.Fatal(err)
	}
	setKubeconfigFolder(t, folder)

	ownExec := map[string]interface{}{"apiVersion": "client.authentication.k8s.io/v1", "command": "aws", "args": []interface{}{"eks", "get-token"}}
	cluster := map[string]interface{}{"server": "https://127.0.0.1:1"}

	t.Run("identical exec is the server's own", func(t *testing.T) {
		auth, cl := kubeCredential(map[string]interface{}{"exec": ownExec}, cluster)
		provider, conn := saveKubeCredentialAndConnection(t, map[string]interface{}{"auth": auth, "cluster": cl})
		if _, err := K8sContextFromConnection(provider, "", conn); err != nil {
			t.Fatalf("K8sContextFromConnection: %v", err)
		}
	})

	t.Run("the server's exec pointed at another host is refused", func(t *testing.T) {
		auth, cl := kubeCredential(map[string]interface{}{"exec": ownExec}, map[string]interface{}{"server": "https://attacker.example:6443"})
		provider, conn := saveKubeCredentialAndConnection(t, map[string]interface{}{"auth": auth, "cluster": cl})
		if _, err := K8sContextFromConnection(provider, "", conn); meshkiterrors.GetCode(err) != ErrUnsafeKubeCredentialCode {
			t.Fatalf("err = %v, want ErrUnsafeKubeCredential", err)
		}
	})

	t.Run("the server's exec under a user the server does not pair with the cluster is refused", func(t *testing.T) {
		auth := map[string]interface{}{"name": "other-user", "user": map[string]interface{}{"exec": ownExec}}
		cl := map[string]interface{}{"name": "c", "cluster": cluster}
		provider, conn := saveKubeCredentialAndConnection(t, map[string]interface{}{"auth": auth, "cluster": cl})
		if _, err := K8sContextFromConnection(provider, "", conn); meshkiterrors.GetCode(err) != ErrUnsafeKubeCredentialCode {
			t.Fatalf("err = %v, want ErrUnsafeKubeCredentialCode", err)
		}
	})

	t.Run("an edited exec is refused even under the same user name", func(t *testing.T) {
		edited := map[string]interface{}{"apiVersion": "client.authentication.k8s.io/v1", "command": "sh", "args": []interface{}{"-c", "id"}}
		auth, cl := kubeCredential(map[string]interface{}{"exec": edited}, cluster)
		provider, conn := saveKubeCredentialAndConnection(t, map[string]interface{}{"auth": auth, "cluster": cl})
		if _, err := K8sContextFromConnection(provider, "", conn); meshkiterrors.GetCode(err) != ErrUnsafeKubeCredentialCode {
			t.Fatalf("err = %v, want ErrUnsafeKubeCredential", err)
		}
	})
}
