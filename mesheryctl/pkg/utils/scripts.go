package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	// tokenSecretPollInterval is how often to check whether the controller has
	// populated the service-account token Secret.
	tokenSecretPollInterval = 1 * time.Second
	// tokenSecretWaitTimeout bounds how long we wait for token/CA data.
	// The previous shell script used a fixed 10s sleep; polling with a longer
	// timeout is more reliable across slow API servers without OS-specific tools.
	tokenSecretWaitTimeout = 30 * time.Second
)

// GKEConfig holds configuration for GKE kubeconfig generation.
// Despite the name, this generates a ServiceAccount-based kubeconfig against
// the ambient Kubernetes context (typically obtained via
// `gcloud container clusters get-credentials` for GKE).
type GKEConfig struct {
	ConfigPath string
	SAName     string
	Namespace  string
}

// GenerateConfigGKE generates a kubeconfig that authenticates as a newly
// created cluster-admin ServiceAccount and writes it to configPath.
//
// parameters: configPath, SAName = Service Account Name, namespace
//
// Implementation is pure Go (client-go + stdlib) so it works on Linux, macOS,
// and Windows without jq/base64/awk/tail/sh.
func GenerateConfigGKE(configPath, SAName, namespace string) error {
	cfg := &GKEConfig{
		ConfigPath: configPath,
		SAName:     SAName,
		Namespace:  namespace,
	}

	if err := cfg.validate(); err != nil {
		return ErrInvalidArgument(err)
	}

	return cfg.generate()
}

func (c *GKEConfig) validate() error {
	if c.ConfigPath == "" || c.SAName == "" || c.Namespace == "" {
		return fmt.Errorf("configPath, SAName, and namespace are required")
	}
	return nil
}

func (c *GKEConfig) generate() error {
	ctx := context.Background()

	restConfig, rawConfig, err := loadAmbientKubeConfig()
	if err != nil {
		return ErrKubernetesConnectivity(fmt.Errorf("failed to load kubeconfig: %w", err))
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return ErrKubernetesConnectivity(fmt.Errorf("failed to create Kubernetes client: %w", err))
	}

	if err := c.checkConnectivity(ctx, clientset); err != nil {
		return err
	}

	clusterName, endpoint, err := clusterEndpointFromConfig(rawConfig)
	if err != nil {
		return ErrKubernetesQuery(err)
	}

	if err := os.MkdirAll(filepath.Dir(c.ConfigPath), 0o755); err != nil {
		return ErrCreateFile(c.ConfigPath, fmt.Errorf("failed to create config directory: %w", err))
	}

	logInfo("Service Account Creation")
	if err := c.createServiceAccount(ctx, clientset); err != nil {
		return ErrKubernetesQuery(err)
	}
	if err := c.createClusterRoleBinding(ctx, clientset); err != nil {
		return ErrKubernetesQuery(err)
	}
	if err := c.createTokenSecret(ctx, clientset); err != nil {
		return ErrKubernetesQuery(err)
	}

	logInfo("Waiting for service account token Secret...")
	token, caCRT, err := c.waitForTokenSecret(ctx, clientset)
	if err != nil {
		return ErrKubernetesQuery(err)
	}

	if err := c.writeKubeconfig(clusterName, endpoint, token, caCRT); err != nil {
		return err
	}

	logInfo(fmt.Sprintf("Configuration generated at: %s", c.ConfigPath))
	logInfo(fmt.Sprintf("Test access with: KUBECONFIG=%s kubectl get pods", c.ConfigPath))
	return nil
}

func loadAmbientKubeConfig() (*rest.Config, clientcmdapi.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, clientcmdapi.Config{}, err
	}

	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, clientcmdapi.Config{}, err
	}

	return restConfig, rawConfig, nil
}

func (c *GKEConfig) checkConnectivity(ctx context.Context, clientset kubernetes.Interface) error {
	if _, err := clientset.Discovery().ServerVersion(); err != nil {
		return ErrKubernetesConnectivity(fmt.Errorf("failed to connect to Kubernetes API server: %w", err))
	}

	if _, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{Limit: 1}); err != nil {
		return ErrKubernetesQuery(fmt.Errorf("failed to query Kubernetes API: %w", err))
	}

	return nil
}

func clusterEndpointFromConfig(rawConfig clientcmdapi.Config) (clusterName, endpoint string, err error) {
	if rawConfig.CurrentContext == "" {
		return "", "", fmt.Errorf("no current context set in kubeconfig")
	}

	ctx, ok := rawConfig.Contexts[rawConfig.CurrentContext]
	if !ok || ctx == nil {
		return "", "", fmt.Errorf("current context %q not found in kubeconfig", rawConfig.CurrentContext)
	}

	clusterName = ctx.Cluster
	cluster, ok := rawConfig.Clusters[clusterName]
	if !ok || cluster == nil {
		return "", "", fmt.Errorf("cluster %q not found in kubeconfig", clusterName)
	}

	if cluster.Server == "" {
		return "", "", fmt.Errorf("cluster %q has empty server endpoint", clusterName)
	}

	return clusterName, cluster.Server, nil
}

func (c *GKEConfig) createServiceAccount(ctx context.Context, clientset kubernetes.Interface) error {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.SAName,
			Namespace: c.Namespace,
		},
	}

	_, err := clientset.CoreV1().ServiceAccounts(c.Namespace).Create(ctx, sa, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			logInfo(fmt.Sprintf("Service account %q already exists in namespace %q", c.SAName, c.Namespace))
			return nil
		}
		return fmt.Errorf("failed to create service account %q: %w", c.SAName, err)
	}

	logInfo(fmt.Sprintf("Service account created in %s namespace", c.Namespace))
	return nil
}

func (c *GKEConfig) createClusterRoleBinding(ctx context.Context, clientset kubernetes.Interface) error {
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.SAName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      c.SAName,
				Namespace: c.Namespace,
			},
		},
	}

	_, err := clientset.RbacV1().ClusterRoleBindings().Create(ctx, crb, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			logInfo(fmt.Sprintf("Cluster role binding %q already exists", c.SAName))
			return nil
		}
		return fmt.Errorf("failed to create cluster role binding %q: %w", c.SAName, err)
	}

	logInfo("Cluster role binding created")
	return nil
}

func (c *GKEConfig) tokenSecretName() string {
	return c.SAName + "-token"
}

func (c *GKEConfig) createTokenSecret(ctx context.Context, clientset kubernetes.Interface) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.tokenSecretName(),
			Namespace: c.Namespace,
			Annotations: map[string]string{
				corev1.ServiceAccountNameKey: c.SAName,
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}

	_, err := clientset.CoreV1().Secrets(c.Namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			logInfo(fmt.Sprintf("Token secret %q already exists", c.tokenSecretName()))
			return nil
		}
		return fmt.Errorf("failed to create token secret %q: %w", c.tokenSecretName(), err)
	}

	logInfo("Token secret created")
	return nil
}

func (c *GKEConfig) waitForTokenSecret(ctx context.Context, clientset kubernetes.Interface) (token string, caCRT []byte, err error) {
	deadline := time.Now().Add(tokenSecretWaitTimeout)
	var lastErr error

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return "", nil, ctx.Err()
		default:
		}

		secret, getErr := clientset.CoreV1().Secrets(c.Namespace).Get(ctx, c.tokenSecretName(), metav1.GetOptions{})
		if getErr != nil {
			lastErr = getErr
			time.Sleep(tokenSecretPollInterval)
			continue
		}

		tokenBytes, okToken := secret.Data[corev1.ServiceAccountTokenKey]
		caBytes, okCA := secret.Data[corev1.ServiceAccountRootCAKey]
		if okToken && len(tokenBytes) > 0 && okCA && len(caBytes) > 0 {
			return string(tokenBytes), caBytes, nil
		}

		lastErr = fmt.Errorf("token or ca.crt not yet populated on secret %q", c.tokenSecretName())
		time.Sleep(tokenSecretPollInterval)
	}

	if lastErr != nil {
		return "", nil, fmt.Errorf("timed out waiting for service account token secret %q: %w", c.tokenSecretName(), lastErr)
	}
	return "", nil, fmt.Errorf("timed out waiting for service account token secret %q", c.tokenSecretName())
}

func (c *GKEConfig) writeKubeconfig(clusterName, endpoint, token string, caCRT []byte) error {
	userName := fmt.Sprintf("%s-%s-%s", c.SAName, c.Namespace, clusterName)
	contextName := userName

	cfg := clientcmdapi.NewConfig()
	cfg.Clusters[clusterName] = &clientcmdapi.Cluster{
		Server:                   endpoint,
		CertificateAuthorityData: caCRT,
	}
	cfg.AuthInfos[userName] = &clientcmdapi.AuthInfo{
		Token: token,
	}
	cfg.Contexts[contextName] = &clientcmdapi.Context{
		Cluster:   clusterName,
		AuthInfo:  userName,
		Namespace: c.Namespace,
	}
	cfg.CurrentContext = contextName

	if err := clientcmd.WriteToFile(*cfg, c.ConfigPath); err != nil {
		return ErrCreateFile(c.ConfigPath, fmt.Errorf("failed to write kubeconfig: %w", err))
	}

	// Optional: also write ca.crt next to the kubeconfig for parity with the
	// legacy script (which wrote ${TARGET_FOLDER}/ca.crt). Embed already
	// includes CA data; writing the file keeps debugging parity.
	caPath := filepath.Join(filepath.Dir(c.ConfigPath), "ca.crt")
	if err := os.WriteFile(caPath, caCRT, 0o600); err != nil {
		// Non-fatal: kubeconfig already embeds the CA.
		logInfo(fmt.Sprintf("Warning: could not write %s: %v", caPath, err))
	}

	logInfo("Kubeconfig written")
	return nil
}

func logInfo(msg string) {
	if Log != nil {
		Log.Info(msg)
		return
	}
	fmt.Println(msg)
}
