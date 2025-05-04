#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLUSTER_NAME="meshery-integration-test-meshsync-cluster"
HELM_LOCAL_REPO_NAME="meshery-integration-test-meshsync"

MESHERY_K8S_NAMESPACE="meshery"
CUSTOM_K8S_NAMESPACE="calm-koala"
PATH_TO_CRDS_YAML="install/kubernetes/helm/meshery-operator/crds/crds.yaml"
PATH_TO_MESHERY_OPERATOR_CHART="install/kubernetes/helm/meshery-operator"

check_dependencies() {
  # Check for docker
  if ! command -v docker &> /dev/null; then
  echo "❌ docker is not installed. Please install docker first."
  exit 1
  fi
  echo "✅ docker is installed;"
  docker --version
  echo ""

  # Check for kind
  if ! command -v kind &> /dev/null; then
  echo "❌ kind is not installed. Please install KinD first."
  exit 1
  fi
  echo "✅ kind is installed;"
  kind --version
  echo ""
  
  # Check for kubectl
  if ! command -v kubectl &> /dev/null; then
  echo "❌ kubectl is not installed. Please install kubectl first."
  exit 1
  fi
  echo "✅ kubectl is installed;"
  kubectl version || echo "if you see [The connection to the server was refused] it is likely that kubectl context is pointing to a deleted cluster which is not an issue."
  echo ""
  
  # Check for helm
  if ! command -v helm &> /dev/null; then
    echo "❌ Helm is not installed. Please install Helm first."
    exit 1
  fi
  echo "✅ Helm is installed."
  helm version
  echo ""
}

# TODO
# 1) for local set up, 
# right now if there is a separate kind cluster with meshery-operator running
# operator in test cluster will fail to retrieve resource lock:
# error retrieving resource lock meshery/operator-981fc876-f149-4e96-a716-9cd9bfb0bd3f.meshery.io: Get "https://10.96.0.1:443/apis/coordination.k8s.io/v1/namespaces/meshery/leases/operator-981fc876-f149-4e96-a716-9cd9bfb0bd3f.meshery.io?timeout=5s": net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)
# and hence will not deploy meshsync and broker
# could be fixed by set up a custom network for test cluster
# --
setup() {
  check_dependencies
  echo "🔧 Setting up..."

  # we do not need this, no standalone nats, use broker in cluster
  # echo "Running docker compose..."
  # docker compose -f $SCRIPT_DIR/docker-compose.yaml up -d || exit 1

  echo "Creating KinD cluster..."
  kind create cluster --name $CLUSTER_NAME
  # kind create cluster --config $SCRIPT_DIR/kind-config.yaml

  echo "Creating meshery namespace..."
  kubectl create namespace $MESHERY_K8S_NAMESPACE

  # this is not working for now, it is applying only operator, not subcharts (meshsync, broker)
  echo "Applying meshery resources..."
  kubectl apply -f $PATH_TO_CRDS_YAML
  helm install meshery-operator $PATH_TO_MESHERY_OPERATOR_CHART --namespace $MESHERY_K8S_NAMESPACE --dependency-update

  echo "Creating $CUSTOM_K8S_NAMESPACE namespace..."
  kubectl create namespace $CUSTOM_K8S_NAMESPACE

  echo "Applying $CUSTOM_K8S_NAMESPACE resources..."
  kubectl --namespace $CUSTOM_K8S_NAMESPACE apply -f $SCRIPT_DIR/test-deployment.yaml

  echo "Outputing cluster resources..."
  kubectl --namespace default get deployment
  kubectl --namespace default get rs
  kubectl --namespace default get po
  kubectl --namespace default get service
  kubectl --namespace default get configmap
  kubectl --namespace $CUSTOM_K8S_NAMESPACE get deployment
  kubectl --namespace $CUSTOM_K8S_NAMESPACE get rs
  kubectl --namespace $CUSTOM_K8S_NAMESPACE get po
  kubectl --namespace $CUSTOM_K8S_NAMESPACE get service
  kubectl --namespace $CUSTOM_K8S_NAMESPACE get configmap
}

cleanup() {
  echo "🧹 Cleaning up..."

  # we do not need this
  # echo "Stopping docker compose..."
  # docker compose -f $SCRIPT_DIR/docker-compose.yaml down

  echo "Deleting KinD cluster..."
  kind delete cluster --name $CLUSTER_NAME
}

print_help() {
  echo "Usage: $0 {check_dependencies|setup|cleanup|help}"
}

# Main dispatcher
case "$1" in
  check_dependencies)
    check_dependencies
    ;;
  setup)
    setup
    ;;
  cleanup)
    cleanup
    ;;
  help)
    print_help
    ;;
  *)
    echo "❌ Unknown command: $1"
    print_help
    exit 1
    ;;
esac