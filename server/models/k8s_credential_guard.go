package models

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// A Kubernetes credential is a kubeconfig user plus cluster ({auth, cluster}),
// and building a client from it is not passive: an `exec` plugin is run as a
// command on the Meshery Server host, and certificate, key and token file
// references are read from its disk. Meshery therefore only builds a client
// from fields that are plain data, unless the credential is the server's own
// kubeconfig, which the operator controls (see issue #24218).

// Names of the unsafe kubeconfig fields, as reported by UnsafeKubeCredentialFields.
const (
	kubeFieldExec                 = "exec"
	kubeFieldAuthProvider         = "auth-provider"
	kubeFieldClientCertificate    = "client-certificate"
	kubeFieldClientKey            = "client-key"
	kubeFieldTokenFile            = "tokenFile"
	kubeFieldCertificateAuthority = "certificate-authority"
)

// UnsafeKubeCredentialFields returns the names of the kubeconfig fields in a
// credential's auth and cluster that make building a client run a command or
// read a file on the Meshery Server host. It inspects what client-go will
// actually see (the kubeconfig GenerateKubeConfig renders, parsed by clientcmd)
// rather than the raw map, so no spelling of the YAML hides a field from it.
// The result is sorted and carries no values, so it is safe to return to a
// caller. It returns nil for a credential that has neither auth nor cluster.
func UnsafeKubeCredentialFields(auth, cluster map[string]interface{}) []string {
	if len(auth) == 0 && len(cluster) == 0 {
		return nil
	}
	cfg, err := loadCredentialKubeconfig(auth, cluster)
	if err != nil {
		// A credential that does not render to a kubeconfig cannot be used to
		// build a client either.
		return nil
	}
	return unsafeFields(cfg, nil)
}

func loadCredentialKubeconfig(auth, cluster map[string]interface{}) (*clientcmdapi.Config, error) {
	data, err := K8sContext{Name: "credential", Auth: auth, Cluster: cluster}.GenerateKubeConfig()
	if err != nil {
		return nil, err
	}
	return clientcmd.Load(data)
}

// unsafeFields lists the unsafe fields set in cfg. When trusted is non-nil, a
// field whose value is identical to the one the trusted kubeconfig holds for the
// same user (or cluster) name is not reported.
func unsafeFields(cfg, trusted *clientcmdapi.Config) []string {
	found := map[string]struct{}{}
	for name, u := range cfg.AuthInfos {
		var t *clientcmdapi.AuthInfo
		if trusted != nil {
			t = trusted.AuthInfos[name]
		}
		if u.Exec != nil && (t == nil || !reflect.DeepEqual(u.Exec, t.Exec)) {
			found[kubeFieldExec] = struct{}{}
		}
		if u.AuthProvider != nil && (t == nil || !reflect.DeepEqual(u.AuthProvider, t.AuthProvider)) {
			found[kubeFieldAuthProvider] = struct{}{}
		}
		if u.ClientCertificate != "" && (t == nil || u.ClientCertificate != t.ClientCertificate) {
			found[kubeFieldClientCertificate] = struct{}{}
		}
		if u.ClientKey != "" && (t == nil || u.ClientKey != t.ClientKey) {
			found[kubeFieldClientKey] = struct{}{}
		}
		if u.TokenFile != "" && (t == nil || u.TokenFile != t.TokenFile) {
			found[kubeFieldTokenFile] = struct{}{}
		}
	}
	for name, c := range cfg.Clusters {
		var t *clientcmdapi.Cluster
		if trusted != nil {
			t = trusted.Clusters[name]
		}
		if c.CertificateAuthority != "" && (t == nil || c.CertificateAuthority != t.CertificateAuthority) {
			found[kubeFieldCertificateAuthority] = struct{}{}
		}
	}
	var out []string
	for f := range found {
		out = append(out, f)
	}
	sort.Strings(out)
	return out
}

// serverKubeconfig loads the kubeconfig the Meshery Server itself runs with
// (KUBECONFIG_FOLDER/config), or nil when there is none (for example in-cluster).
func serverKubeconfig() *clientcmdapi.Config {
	folder := viper.GetString("KUBECONFIG_FOLDER")
	if folder == "" {
		return nil
	}
	path := filepath.Join(folder, "config")
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	cfg, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return nil
	}
	return cfg
}

// UntrustedKubeCredentialFields returns the unsafe fields of a credential's auth
// and cluster that the server's own kubeconfig does not vouch for. A user and
// cluster pair is the server's own only when a context of the server's
// kubeconfig names both, the cluster points at the same endpoint (server, proxy,
// TLS server name, skip-verify), and every unsafe field is identical to the
// server's. That keeps contexts imported from the server's kubeconfig (an EKS or
// GKE exec plugin, say) working, while a credential a user created or edited
// cannot introduce a command or a path of their choosing, nor reuse the
// operator's exec plugin to send a token to another host.
func UntrustedKubeCredentialFields(auth, cluster map[string]interface{}) []string {
	if len(auth) == 0 && len(cluster) == 0 {
		return nil
	}
	cfg, err := loadCredentialKubeconfig(auth, cluster)
	if err != nil {
		return nil
	}
	fields := unsafeFields(cfg, nil)
	if len(fields) == 0 {
		return nil
	}
	srv := serverKubeconfig()
	if srv == nil || !pairsMatchServerContext(cfg, srv) {
		return fields
	}
	return unsafeFields(cfg, srv)
}

// pairsMatchServerContext reports whether every user of cfg is paired, in a
// context of srv, with a cluster of srv that has the same name and endpoint as
// the one cfg carries.
func pairsMatchServerContext(cfg, srv *clientcmdapi.Config) bool {
	for uname := range cfg.AuthInfos {
		for cname, c := range cfg.Clusters {
			sc := srv.Clusters[cname]
			if sc == nil || c.Server != sc.Server || c.ProxyURL != sc.ProxyURL ||
				c.TLSServerName != sc.TLSServerName || c.InsecureSkipTLSVerify != sc.InsecureSkipTLSVerify {
				return false
			}
			paired := false
			for _, sctx := range srv.Contexts {
				if sctx.AuthInfo == uname && sctx.Cluster == cname {
					paired = true
					break
				}
			}
			if !paired {
				return false
			}
		}
	}
	return len(cfg.AuthInfos) > 0
}
