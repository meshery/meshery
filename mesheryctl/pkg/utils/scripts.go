package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// GKEConfig holds configuration for GKE script generation

type GKEConfig struct {
	ConfigPath string
	SAName     string
	Namespace  string
}

/*
GenerateConfigGKE generates kubernetes config for GKE
parameters :- configPath ,SAName=Service Account Name, namespace
*/
func GenerateConfigGKE(configPath, SAName, namespace string) error {
	cfg := &GKEConfig{
		ConfigPath: configPath,
		SAName:     SAName,
		Namespace:  namespace,
	}

	if err := cfg.validate(); err != nil {
		return ErrInvalidArgument(err)
	}

	return cfg.executeScript()
}

func (c *GKEConfig) validate() error {
	if c.ConfigPath == "" || c.SAName == "" || c.Namespace == "" {
		return ErrInvalidArgument(fmt.Errorf("configPath, SAName, and namespace are required"))
	}
	return nil
}

func (c *GKEConfig) checkPrerequisites() error {
	var missingCommands []string
	requiredCommands := []string{"jq", "base64", "awk", "tail"}

	for _, cmd := range requiredCommands {
		if _, err := exec.LookPath(cmd); err != nil {
			missingCommands = append(missingCommands, cmd)
		}
	}

	if len(missingCommands) > 0 {
		return ErrMissingCommands(fmt.Errorf("missing required commands: %v", missingCommands))
	}
	return nil
}

func (c *GKEConfig) executeScript() error {
	if err := c.checkConnectivity(); err != nil {
		return err
	}

	if err := c.checkPrerequisites(); err != nil {
		return err
	}

	script := c.generateScript()
	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *GKEConfig) checkConnectivity() error {

	// Check if client can be initialized
	cmd := exec.Command("kubectl", "cluster-info")
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return ErrKubernetesConnectivity(fmt.Errorf("failed to connect to Kubernetes API server: %w", err))
	}

	// Check API query
	cmd = exec.Command("kubectl", "get", "nodes")
	if err := cmd.Run(); err != nil {
		return ErrKubernetesQuery(fmt.Errorf("failed to query Kubernetes API: %w", err))
	}

	return nil
}

func (c *GKEConfig) generateScript() string {
	return fmt.Sprintf(`
set -e
set -o pipefail

KUBECFG_FILE_NAME="%s"
SERVICE_ACCOUNT_NAME=%s
NAMESPACE="%s"
TARGET_FOLDER=$(dirname ${KUBECFG_FILE_NAME})

%s
%s
%s
%s

main
`, c.ConfigPath, c.SAName, c.Namespace,
		serviceAccountTemplate(),
		secretOpsTemplate(),
		kubeConfigTemplate(),
		mainOpsTemplate())
}

func serviceAccountTemplate() string {
	return `create_service_account() {
    echo "\nService Account Creation"
    echo "--------------"
    
    if kubectl create sa "${SERVICE_ACCOUNT_NAME}" --namespace "${NAMESPACE}"; then
        echo "✓ Service account created in ${NAMESPACE} namespace"
    else
        echo "!! Failed to create service account"
        exit 1
    fi
    
    if kubectl create clusterrolebinding "${SERVICE_ACCOUNT_NAME}" --clusterrole=cluster-admin --serviceaccount=default:"${SERVICE_ACCOUNT_NAME}" --namespace "${NAMESPACE}"; then
        echo "✓ Cluster role binding created"
    else
        echo "!! Failed to create cluster role binding"
        exit 1
    fi
    
    echo "\nToken Secret Creation"
    echo "--------------"
    cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: "${SERVICE_ACCOUNT_NAME}-token"
  namespace: "${NAMESPACE}"
  annotations:
    kubernetes.io/service-account.name: "${SERVICE_ACCOUNT_NAME}"
type: kubernetes.io/service-account-token
EOF
    echo "✓ Token secret created"
    sleep 10
}`
}

func secretOpsTemplate() string {
	return `
get_secret_name_from_service_account() {
    echo "\nSecret Operations"
    echo "--------------"
    SECRET_NAME="${SERVICE_ACCOUNT_NAME}-token"
    echo "✓ Found secret: ${SECRET_NAME}"
}

extract_ca_crt_from_secret() {
    if kubectl get secret --namespace "${NAMESPACE}" "${SECRET_NAME}" -o jsonpath='{.data.ca\.crt}' | base64 --decode > "${TARGET_FOLDER}/ca.crt"; then
        echo "✓ Extracted ca.crt from secret"
    else
        echo "!! Failed to extract ca.crt"
        exit 1
    fi
}

get_user_token_from_secret() {
    USER_TOKEN=$(kubectl get secret --namespace "${NAMESPACE}" "${SECRET_NAME}" -o json | jq -r '.data["token"]' | base64 --decode)
    if [ -n "${USER_TOKEN}" ]; then
        echo "✓ Retrieved user token"
    else
        echo "!! Failed to get user token"
        exit 1
    fi
}`
}

func kubeConfigTemplate() string {
	return `
set_kube_config_values() {
    echo "\nKubeconfig Setup"
    echo "--------------"
    
    context=$(kubectl config current-context)
    CLUSTER_NAME=$(kubectl config get-contexts "$context" | awk '{print $3}' | tail -n 1)
    ENDPOINT=$(kubectl config view -o jsonpath="{.clusters[?(@.name == \"${CLUSTER_NAME}\")].cluster.server}")

    if [ -z "${ENDPOINT}" ]; then
        echo "!! Failed to get cluster endpoint"
        exit 1
    fi
    echo "✓ Retrieved cluster information"

    echo "\nConfiguring Authentication"
    echo "--------------"
    
    kubectl config set-cluster "${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}" \
    --server="${ENDPOINT}" \
    --certificate-authority="${TARGET_FOLDER}/ca.crt" \
    --embed-certs=true && echo "✓ Cluster configuration set"

    kubectl config set-credentials \
    "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}" \
    --token="${USER_TOKEN}" && echo "✓ Credentials configured"

    kubectl config set-context \
    "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}" \
    --cluster="${CLUSTER_NAME}" \
    --user="${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --namespace="${NAMESPACE}" && echo "✓ Context created"

    kubectl config use-context "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}" && echo "✓ Context activated"
}`
}

func mainOpsTemplate() string {
	return `
main() {
    echo "\nPrerequisites"
    echo "--------------"
    mkdir -p "${TARGET_FOLDER}" && echo "✓ Target directory created"
    
    create_service_account
    get_secret_name_from_service_account
    extract_ca_crt_from_secret
    get_user_token_from_secret
    set_kube_config_values

    echo "\nSetup Complete"
    echo "--------------"
    echo "✓ Configuration generated at: ${KUBECFG_FILE_NAME}"
    echo "✓ Test access with: KUBECONFIG=${KUBECFG_FILE_NAME} kubectl get pods"
    echo "!! Note: RBAC permissions need to be configured separately"
}`
}
