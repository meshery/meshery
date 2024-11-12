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
func GenerateConfigGKE(configPath, SAName, namespc string) error {
	cfg := &GKEConfig{
		ConfigPath: configPath,
		SAName:     SAName,
		Namespace:  namespc,
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
	requiredCommands := []string{"jq", "base64", "awk", "tail"}
	for _, cmd := range requiredCommands {
		if _, err := exec.LookPath(cmd); err != nil {
			return fmt.Errorf("required command %s is not installed", cmd)
		}
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
	cmd := exec.Command("kubectl", "cluster-info")
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to connect to Kubernetes API server: %w", err)
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
    echo "Creating a service account in ${NAMESPACE} namespace: ${SERVICE_ACCOUNT_NAME}"
    if ! kubectl create sa "${SERVICE_ACCOUNT_NAME}" --namespace "${NAMESPACE}"; then
        echo "[ERROR] Failed to create service account"
        exit 1
    fi
    if ! kubectl create clusterrolebinding "${SERVICE_ACCOUNT_NAME}" --clusterrole=cluster-admin --serviceaccount=default:"${SERVICE_ACCOUNT_NAME}" --namespace "${NAMESPACE}"; then
        echo "[ERROR] Failed to create cluster role binding"
        exit 1
    fi
    
    echo "Creating token secret for service account..."
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
    
    sleep 10
}`
}

func secretOpsTemplate() string {
	return `
get_secret_name_from_service_account() {
    echo "Getting secret of service account ${SERVICE_ACCOUNT_NAME} on ${NAMESPACE}"
    SECRET_NAME="${SERVICE_ACCOUNT_NAME}-token"
    echo "Secret name: ${SECRET_NAME}"
}

extract_ca_crt_from_secret() {
    echo -n "Extracting ca.crt from secret..."
    kubectl get secret --namespace "${NAMESPACE}" "${SECRET_NAME}" -o jsonpath='{.data.ca\.crt}' | base64 --decode > "${TARGET_FOLDER}/ca.crt"
    echo "done"
}

get_user_token_from_secret() {
    echo -n "Getting user token from secret..."
    USER_TOKEN=$(kubectl get secret --namespace "${NAMESPACE}" "${SECRET_NAME}" -o json | jq -r '.data["token"]' | base64 --decode)
    if [ -z "${USER_TOKEN}" ]; then
        echo "Error: Failed to get user token"
        exit 1
    fi
    echo "done"
}`
}

func kubeConfigTemplate() string {
	return `
set_kube_config_values() {
    context=$(kubectl config current-context)
    echo "Setting current context to: $context"

    CLUSTER_NAME=$(kubectl config get-contexts "$context" | awk '{print $3}' | tail -n 1)
    echo "Cluster name: ${CLUSTER_NAME}"

    ENDPOINT=$(kubectl config view -o jsonpath="{.clusters[?(@.name == \"${CLUSTER_NAME}\")].cluster.server}")
    echo "Endpoint: ${ENDPOINT}"

    if [ -z "${ENDPOINT}" ]; then
        echo "Error: Failed to get cluster endpoint"
        exit 1
    fi

    echo "Preparing k8s-${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-conf"
    kubectl config set-cluster "${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}" \
    --server="${ENDPOINT}" \
    --certificate-authority="${TARGET_FOLDER}/ca.crt" \
    --embed-certs=true

    kubectl config set-credentials \
    "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}" \
    --token="${USER_TOKEN}"

    kubectl config set-context \
    "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}" \
    --cluster="${CLUSTER_NAME}" \
    --user="${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --namespace="${NAMESPACE}"

    kubectl config use-context "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}"
}`
}

func mainOpsTemplate() string {
	return `
main() {
    mkdir -p "${TARGET_FOLDER}"
    create_service_account
    get_secret_name_from_service_account
    extract_ca_crt_from_secret
    get_user_token_from_secret
    set_kube_config_values

    echo "All done! Test with:"
    echo "KUBECONFIG=${KUBECFG_FILE_NAME} kubectl get pods"
    echo "You should not have any permissions by default - you have just created the authentication part"
    echo "You will need to create RBAC permissions"
}`
}
