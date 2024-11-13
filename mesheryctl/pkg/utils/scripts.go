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

// GenerateConfigGKE generates kubernetes config for GKE
func GenerateConfigGKE(configPath, SAName, namespace string) error {
	cfg := &GKEConfig{
		ConfigPath: configPath,
		SAName:     SAName,
		Namespace:  namespace,
	}

	if err := cfg.validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return cfg.executeScript()
}

func (c *GKEConfig) validate() error {
	if c.ConfigPath == "" || c.SAName == "" || c.Namespace == "" {
		return fmt.Errorf("configPath, SAName, and namespc are required")
	}
	return nil
}

func (c *GKEConfig) checkPrerequisites() error {
	var missingCommands []string
	requiredCommands := []string{"jq", "base64", "awk", "tail"}

	fmt.Println("\nRequired Commands")
	fmt.Println("--------------")

	for _, cmd := range requiredCommands {
		if _, err := exec.LookPath(cmd); err != nil {
			fmt.Printf("!! %s is not available\n", cmd)
			missingCommands = append(missingCommands, cmd)
		} else {
			fmt.Printf("✓ %s is available\n", cmd)
		}
	}

	if len(missingCommands) > 0 {
		return fmt.Errorf("missing required commands: %v", missingCommands)
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
	fmt.Println("\nKubernetes API")
	fmt.Println("--------------")

	// Check if client can be initialized
	cmd := exec.Command("kubectl", "cluster-info")
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		fmt.Println("!! cannot connect to Kubernetes API server")
		return fmt.Errorf("failed to connect to Kubernetes API server: %w", err)
	}
	fmt.Println("✓ can connect to Kubernetes API server")

	// Check API query
	cmd = exec.Command("kubectl", "get", "nodes")
	if err := cmd.Run(); err != nil {
		fmt.Println("!! cannot query Kubernetes API")
		return fmt.Errorf("failed to query Kubernetes API: %w", err)
	}
	fmt.Println("✓ can query Kubernetes API")

	fmt.Println("\n--------------")
	fmt.Println("--------------")

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
