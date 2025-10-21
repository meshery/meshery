package helpers

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGetK8SClientSetWithInsecureSkipVerify(t *testing.T) {
    // Test kubeconfig with insecure-skip-tls-verify set to true
    kubeconfigContent := []byte(`
apiVersion: v1
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://127.0.0.1:6443
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
kind: Config
users:
- name: test-user
  user:
    token: test-token
`)
    
    // Test that the clientset can be created without errors
    clientset, err := getK8SClientSet(kubeconfigContent, "")
    
    // The clientset should be created successfully
    require.NoError(t, err)
    require.NotNil(t, clientset)
}

func TestGetK8SClientSetWithCADataAndInsecureSkipVerify(t *testing.T) {
    // Test that when both CA data and insecure-skip-tls-verify are present,
    // the insecure flag takes precedence and CA data is cleared
    kubeconfigContent := []byte(`
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJd01EWXdOakEzTVRReE0xb1hEVE13TURZd05EQTNNVFF4TTFvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTDNqCg==
    insecure-skip-tls-verify: true
    server: https://127.0.0.1:6443
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
kind: Config
users:
- name: test-user
  user:
    token: test-token
`)
    
    clientset, err := getK8SClientSet(kubeconfigContent, "")
    
    // Should succeed without "specifying a root certificates file with the insecure flag is not allowed" error
    require.NoError(t, err)
    require.NotNil(t, clientset)
}

func TestGetK8SClientSetWithContextName(t *testing.T) {
    // Test that context switching works correctly
    kubeconfigContent := []byte(`
apiVersion: v1
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://127.0.0.1:6443
  name: cluster-1
- cluster:
    insecure-skip-tls-verify: true
    server: https://192.168.1.1:6443
  name: cluster-2
contexts:
- context:
    cluster: cluster-1
    user: user-1
  name: context-1
- context:
    cluster: cluster-2
    user: user-2
  name: context-2
current-context: context-1
kind: Config
users:
- name: user-1
  user:
    token: token-1
- name: user-2
  user:
    token: token-2
`)
    
    // Test with specific context name
    clientset, err := getK8SClientSet(kubeconfigContent, "context-2")
    
    require.NoError(t, err)
    require.NotNil(t, clientset)
}

func TestGetK8SClientSetWithEmptyKubeconfig(t *testing.T) {
    // Test that in-cluster config is attempted when kubeconfig is empty
    _, err := getK8SClientSet([]byte{}, "")
    
    // We expect an error since we're not running inside a cluster
    assert.Error(t, err)
}
