package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// TODO@Aisuko this function is so hard to read and maintain, please refactor it
func GenerateConfigGKE(configPath, SAName, namespc string) error {
	if configPath == "" || SAName == "" || namespc == "" {
		return fmt.Errorf("configPath, SAName, and namespc are required")
	}

	script := fmt.Sprintf(`
set -e
set -o pipefail

KUBECFG_FILE_NAME="%s"
SERVICE_ACCOUNT_NAME=%s
NAMESPACE="%s"
TARGET_FOLDER=$(dirname ${KUBECFG_FILE_NAME})

create_service_account() {
    echo -e "\nCreating a service account in ${NAMESPACE} namespace: ${SERVICE_ACCOUNT_NAME}"
    kubectl create sa "${SERVICE_ACCOUNT_NAME}" --namespace "${NAMESPACE}"
    kubectl create clusterrolebinding "${SERVICE_ACCOUNT_NAME}" --clusterrole=cluster-admin --serviceaccount=default:"${SERVICE_ACCOUNT_NAME}" --namespace "${NAMESPACE}"
    
    # Create token secret for service account
    echo -e "\nCreating token secret for service account..."
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
    
    # Add delay to allow secret creation
    echo "Waiting for service account secret to be created..."
    sleep 10
}

get_secret_name_from_service_account() {
    echo -e "\nGetting secret of service account ${SERVICE_ACCOUNT_NAME} on ${NAMESPACE}"
    # Get the secret name
    SECRET_NAME="${SERVICE_ACCOUNT_NAME}-token"
    echo "Secret name: ${SECRET_NAME}"
}

extract_ca_crt_from_secret() {
    echo -e -n "\nExtracting ca.crt from secret..."
    kubectl get secret --namespace "${NAMESPACE}" "${SECRET_NAME}" -o jsonpath='{.data.ca\.crt}' | base64 --decode > "${TARGET_FOLDER}/ca.crt"
    printf "done"
}

get_user_token_from_secret() {
    echo -e -n "\nGetting user token from secret..."
    USER_TOKEN=$(kubectl get secret --namespace "${NAMESPACE}" "${SECRET_NAME}" -o json | jq -r '.data["token"]' | base64 --decode)
    if [ -z "${USER_TOKEN}" ]; then
        echo "Error: Failed to get user token"
        exit 1
    fi
    printf "done"
}

set_kube_config_values() {
    context=$(kubectl config current-context)
    echo -e "\nSetting current context to: $context"

    CLUSTER_NAME=$(kubectl config get-contexts "$context" | awk '{print $3}' | tail -n 1)
    echo "Cluster name: ${CLUSTER_NAME}"

    ENDPOINT=$(kubectl config view -o jsonpath="{.clusters[?(@.name == \"${CLUSTER_NAME}\")].cluster.server}")
    echo "Endpoint: ${ENDPOINT}"

    if [ -z "${ENDPOINT}" ]; then
        echo "Error: Failed to get cluster endpoint"
        exit 1
    fi

    # Set up the config
    echo -e "\nPreparing k8s-${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-conf"
    echo -n "Setting a cluster entry in kubeconfig..."
    kubectl config set-cluster "${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}" \
    --server="${ENDPOINT}" \
    --certificate-authority="${TARGET_FOLDER}/ca.crt" \
    --embed-certs=true

    echo -n "Setting token credentials entry in kubeconfig..."
    kubectl config set-credentials \
    "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}" \
    --token="${USER_TOKEN}"

    echo -n "Setting a context entry in kubeconfig..."
    kubectl config set-context \
    "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}" \
    --cluster="${CLUSTER_NAME}" \
    --user="${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --namespace="${NAMESPACE}"

    echo -n "Setting the current-context in the kubeconfig file..."
    kubectl config use-context "${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-${CLUSTER_NAME}" \
    --kubeconfig="${KUBECFG_FILE_NAME}"
}

main() {
    mkdir -p "${TARGET_FOLDER}"
    create_service_account
    get_secret_name_from_service_account
    extract_ca_crt_from_secret
    get_user_token_from_secret
    set_kube_config_values

    echo -e "\nAll done! Test with:"
    echo "KUBECONFIG=${KUBECFG_FILE_NAME} kubectl get pods"
    echo "you should not have any permissions by default - you have just created the authentication part"
    echo "You will need to create RBAC permissions"
}

main
`, configPath, SAName, namespc)

	generateCFG := exec.Command("sh", "-c", script)
	generateCFG.Stdout = os.Stdout
	generateCFG.Stderr = os.Stderr

	return generateCFG.Run()
}
