package helpers

import (
	"strings"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

func uploadKubeconfig(userBlock string) []byte {
	return []byte(`apiVersion: v1
kind: Config
current-context: c
clusters:
- name: c
  cluster:
    server: https://127.0.0.1:6443
` + "    CLUSTEREXTRA\n" + `contexts:
- name: c
  context:
    cluster: c
    user: u
users:
- name: u
  user:
` + userBlock)
}

func TestRejectKubeConfigFileReferences(t *testing.T) {
	tests := []struct {
		name    string
		config  []byte
		wantErr bool
	}{
		{"embedded data only", []byte(strings.ReplaceAll(string(uploadKubeconfig("    client-certificate-data: YQ==\n    client-key-data: Yg==\n")), "    CLUSTEREXTRA\n", "    certificate-authority-data: Yw==\n")), false},
		{"token only", []byte(strings.ReplaceAll(string(uploadKubeconfig("    token: abc\n")), "    CLUSTEREXTRA\n", "")), false},
		{"client-key path that does not exist", []byte(strings.ReplaceAll(string(uploadKubeconfig("    client-key: /no/such/dir/client.key\n")), "    CLUSTEREXTRA\n", "")), true},
		{"client-key path that exists", []byte(strings.ReplaceAll(string(uploadKubeconfig("    client-key: /etc/hostname\n")), "    CLUSTEREXTRA\n", "")), true},
		{"client-certificate path", []byte(strings.ReplaceAll(string(uploadKubeconfig("    client-certificate: /etc/hostname\n")), "    CLUSTEREXTRA\n", "")), true},
		{"tokenFile path", []byte(strings.ReplaceAll(string(uploadKubeconfig("    tokenFile: /var/run/secrets/token\n")), "    CLUSTEREXTRA\n", "")), true},
		{"certificate-authority path", []byte(strings.ReplaceAll(string(uploadKubeconfig("    token: abc\n")), "    CLUSTEREXTRA\n", "    certificate-authority: /etc/ssl/ca.crt\n")), true},
		{"relative path", []byte(strings.ReplaceAll(string(uploadKubeconfig("    client-key: client.key\n")), "    CLUSTEREXTRA\n", "")), true},
		{"not yaml", []byte("\t: : ["), false},
		// The generic decode fails on the unknown "junk" key but the typed decode
		// the consumers use does not, so the check must not rely on one decode.
		{"unparsable unknown key beside a tokenFile", []byte("apiVersion: v1\nkind: Config\njunk: {[a]: b}\nclusters: [{name: c, cluster: {server: \"https://attacker.example:6443\"}}]\ncontexts: [{name: c, context: {cluster: c, user: u}}]\nusers:\n- name: u\n  user: {tokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token}\n"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RejectKubeConfigFileReferences(tt.config)
			if tt.wantErr {
				if meshkiterrors.GetCode(err) != ErrKubeConfigFileReferenceCode {
					t.Fatalf("err = %v, want ErrKubeConfigFileReferenceCode", err)
				}
				if strings.Contains(err.Error(), "/") {
					t.Fatalf("error %q must not carry a path", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
