package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// GenerateConfigMinikube generates kube config file in ~/.meshery/kubeconfig.yaml for a Minikube cluster
func GenerateConfigMinikube(configPath string) error {
	script := fmt.Sprintf(`
	set -e
	set -o pipefail

	TARGET_FILE="%s"

	kubectl config view --minify --flatten > ${TARGET_FILE}
	`, configPath)

	generateCFG := exec.Command("bash", "-c", script)
	generateCFG.Stdout = os.Stdout
	generateCFG.Stderr = os.Stderr

	return generateCFG.Run()
}

// GenerateConfigGKE generates kube config file in ~/.meshery/kubeconfig.yaml for a GKE cluster
func GenerateConfigGKE(configPath, SAName, namespc string) error {
	script := fmt.Sprintf(`
	set -e
	set -o pipefail

	KUBECFG_FILE_NAME="%s"
	SERVICE_ACCOUNT_NAME=%s
	NAMESPACE="%s"
	TARGET_FOLDER=$(dirname ${KUBECFG_FILE_NAME})

	create_service_account() {
		echo -e "\\nCreating a service account in ${NAMESPACE} namespace: ${SERVICE_ACCOUNT_NAME}"
		kubectl create sa "${SERVICE_ACCOUNT_NAME}" --namespace "${NAMESPACE}"
		kubectl create clusterrolebinding "${SERVICE_ACCOUNT_NAME}" --clusterrole=cluster-admin --serviceaccount=default:"${SERVICE_ACCOUNT_NAME}" --namespace "${NAMESPACE}"
	}

	get_secret_name_from_service_account() {
		echo -e "\\nGetting secret of service account ${SERVICE_ACCOUNT_NAME} on ${NAMESPACE}"
		SECRET_NAME=$(kubectl get sa "${SERVICE_ACCOUNT_NAME}" --namespace="${NAMESPACE}" -o json | jq -r .secrets[].name)
		echo "Secret name: ${SECRET_NAME}"
	}

	extract_ca_crt_from_secret() {
		echo -e -n "\\nExtracting ca.crt from secret..."
		kubectl get secret --namespace "${NAMESPACE}" "${SECRET_NAME}" -o json | jq \
		-r '.data["ca.crt"]' | base64 --decode > "${TARGET_FOLDER}/ca.crt"
		printf "done"
	}

	get_user_token_from_secret() {
		echo -e -n "\\nGetting user token from secret..."
		USER_TOKEN=$(kubectl get secret --namespace "${NAMESPACE}" "${SECRET_NAME}" -o json | jq -r '.data["token"]' | base64 --decode)
		printf "done"
	}

	set_kube_config_values() {
		context=$(kubectl config current-context)
		echo -e "\\nSetting current context to: $context"

		CLUSTER_NAME=$(kubectl config get-contexts "$context" | awk '{print $3}' | tail -n 1)
		echo "Cluster name: ${CLUSTER_NAME}"

		ENDPOINT=$(kubectl config view \
		-o jsonpath="{.clusters[?(@.name == \"${CLUSTER_NAME}\")].cluster.server}")
		echo "Endpoint: ${ENDPOINT}"

		# Set up the config
		echo -e "\\nPreparing k8s-${SERVICE_ACCOUNT_NAME}-${NAMESPACE}-conf"
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

	create_service_account
	get_secret_name_from_service_account
	extract_ca_crt_from_secret
	get_user_token_from_secret
	set_kube_config_values

	echo -e "\\nAll done! Test with:"
	echo "KUBECONFIG=${KUBECFG_FILE_NAME} kubectl get pods"
	echo "you should not have any permissions by default - you have just created the authentication part"
	echo "You will need to create RBAC permissions"
	`, configPath, SAName, namespc)

	generateCFG := exec.Command("sh", "-c", script)
	generateCFG.Stdout = os.Stdout
	generateCFG.Stderr = os.Stderr

	return generateCFG.Run()
}

// GenerateConfigAKS generates kube config file in ~/.meshery/kubeconfig.yaml for a AKS cluster
func GenerateConfigAKS(configPath, resourceGroup, clusterName string) error {
	script := fmt.Sprintf(`
	set -e
	set -o pipefail

	TARGET_FILE="%s"
	Resource_Group="%s"
	Cluster_Name="%s"

	fetch_aks_script() {
		printf "\n"
		az aks get-credentials --resource-group "${Resource_Group}" --name "${Cluster_Name}" --file "${TARGET_FILE}"
	}

	fetch_aks_script
	"
	`, configPath, resourceGroup, clusterName)

	generateCFG := exec.Command("bash", "-c", script)
	generateCFG.Stdout = os.Stdout
	generateCFG.Stderr = os.Stderr

	return generateCFG.Run()
}

// GenerateConfigEKS generates kube config file in .meshery/kubeconfig.yaml for an EKS cluster
func GenerateConfigEKS(configPath, region, cluster string) error {
	script := fmt.Sprintf(`
	set -e
	set -o pipefail

	KUBECFG_FILE_NAME="%s"
	REGION_NAME=%s
	CLUSTER_NAME="%s"
	KUBECONFIG=${KUBECFG_FILE_NAME}

	create_update_kubeconfig() {
		echo -e "\\nGenerating kubeconfig for EKS cluster ${CLUSTER_NAME}..."
		aws eks --region ${REGION_NAME} update-kubeconfig --name ${CLUSTER_NAME} --kubeconfig ${KUBECONFIG} >/dev/null 2>&1 || \
		(echo -e "\\naws CLI is not available on the system.\nInstall aws CLI and run 'mesheryctl system config eks' command again" && \
		aws --version >/dev/null 2>&1)
	}

	create_update_kubeconfig

	echo -e "\\nEKS kubeconfig ready for use by Meshery: "${KUBECFG_FILE_NAME}
	`, configPath, region, cluster)

	generateCFG := exec.Command("sh", "-c", script)
	generateCFG.Stdout = os.Stdout
	generateCFG.Stderr = os.Stderr

	return generateCFG.Run()
}
