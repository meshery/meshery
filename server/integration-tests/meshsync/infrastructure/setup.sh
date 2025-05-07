#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLUSTER_NAME="meshery-integration-test-meshsync-cluster"
HELM_LOCAL_REPO_NAME="meshery-integration-test-meshsync"
DOCKER_IMAGE="layer5/meshery"
DOCKER_IMAGE_TAG="integration-test"

MESHERY_K8S_NAMESPACE="meshery"
CUSTOM_K8S_NAMESPACE="calm-koala"
PATH_TO_CRDS_YAML="install/kubernetes/helm/meshery-operator/crds/crds.yaml"
PATH_TO_MESHERY_OPERATOR_CHART="install/kubernetes/helm/meshery-operator"
PATH_TO_MESHERY_CHART="install/kubernetes/helm/meshery"

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
  echo ""

  # image must be build and present with latest tag 
  echo "Checking meshery docker image is present..."
  docker image inspect --format '{{.RepoTags}}' $DOCKER_IMAGE:latest
  echo ""

  echo "Creating KinD cluster..."
  kind create cluster --name $CLUSTER_NAME
  echo ""

  echo "Creating meshery namespace..."
  kubectl create namespace $MESHERY_K8S_NAMESPACE
  echo ""

  echo "Applying meshery resources..."
  kubectl apply -f $PATH_TO_CRDS_YAML
  helm install meshery-operator $PATH_TO_MESHERY_OPERATOR_CHART --namespace $MESHERY_K8S_NAMESPACE --dependency-update
  echo ""

  echo "Creating $CUSTOM_K8S_NAMESPACE namespace..."
  kubectl create namespace $CUSTOM_K8S_NAMESPACE
  echo ""

  echo "Applying $CUSTOM_K8S_NAMESPACE resources..."
  kubectl --namespace $CUSTOM_K8S_NAMESPACE apply -f $SCRIPT_DIR/test-deployment.yaml
  echo ""

  echo "Tag image with custom tag..."
  docker tag $DOCKER_IMAGE:latest $DOCKER_IMAGE:$DOCKER_IMAGE_TAG
  echo ""

  echo "Upload image to kind..."
  kind load docker-image $DOCKER_IMAGE:$DOCKER_IMAGE_TAG --name $CLUSTER_NAME 
  echo ""

  echo "Deploying meshery server..."
  helm install meshery $PATH_TO_MESHERY_CHART \
    --namespace $MESHERY_K8S_NAMESPACE \
    --set image.tag=$DOCKER_IMAGE_TAG \
    --set image.pullPolicy=Never
  echo ""

  sleep 16
  echo "Outputing cluster resources..."
  kubectl --namespace meshery get po

}

cleanup() {
  echo "🧹 Cleaning up..."
  echo ""

  echo "Deleting KinD cluster..."
  kind delete cluster --name $CLUSTER_NAME
  echo ""

  if docker image inspect "$DOCKER_IMAGE:$DOCKER_IMAGE_TAG" > /dev/null 2>&1; then
    echo "Untagging custom image tag..."
    docker rmi $DOCKER_IMAGE:$DOCKER_IMAGE_TAG
  fi
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