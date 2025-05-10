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
TMP_KUBECONFIG_PATH="meshery-integration-test-meshsync-kubeconfig.yaml"
LOCAL_SQLITE_PATH="meshery-integration-test-meshsync-mesherydb.sql"

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

setup_all() {
  setup_cluster
  setup_connection  
}

# TODO (maybe) for local set up.
# right now if there is a separate kind cluster with meshery-operator running
# operator in test cluster will fail to retrieve resource lock:
# error retrieving resource lock meshery/operator-981fc876-f149-4e96-a716-9cd9bfb0bd3f.meshery.io: Get "https://10.96.0.1:443/apis/coordination.k8s.io/v1/namespaces/meshery/leases/operator-981fc876-f149-4e96-a716-9cd9bfb0bd3f.meshery.io?timeout=5s": net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)
# and hence will not deploy meshsync and broker
# probably could be fixed by set up a custom network for test cluster.
# --
setup_cluster() {
  check_dependencies
  echo "🔧 Setting up..."
  echo ""

  # image must be build and present with tag "latest"
  # (docker-build make target could be used to build an image) 
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

  echo "Taging image with custom tag..."
  docker tag $DOCKER_IMAGE:latest $DOCKER_IMAGE:$DOCKER_IMAGE_TAG
  echo ""

  echo "Uploading image to kind..."
  kind load docker-image $DOCKER_IMAGE:$DOCKER_IMAGE_TAG --name $CLUSTER_NAME 
  echo ""

  echo "Deploying meshery server..."
  helm install meshery $PATH_TO_MESHERY_CHART \
    --namespace $MESHERY_K8S_NAMESPACE \
    --set image.tag=$DOCKER_IMAGE_TAG \
    --set image.pullPolicy=Never
  kubectl --namespace $MESHERY_K8S_NAMESPACE wait --for=condition=available deployment/meshery --timeout=60s
  echo ""

  echo "Outputing cluster resources..."
  kubectl --namespace $MESHERY_K8S_NAMESPACE get po
  echo ""
}

setup_connection() {
  echo "Preparing tmp kubeconfig with current contexts..."
  kubectl config view --minify --raw > $TMP_KUBECONFIG_PATH
  # Replace localhost API server address with in-cluster DNS.
  # '127.0.0.1:<port>' won't work inside pods; use 'kubernetes.default.svc' instead.
  sed -i.bak -E 's|^( *server: ).*|\1https://kubernetes.default.svc|' "$TMP_KUBECONFIG_PATH"
  echo ""

  echo "Submitting kubeconfig..." 
  JOB_NAME="integration-test-meshsync-curl-upload-kubeconfig-job"
  kubectl --namespace $MESHERY_K8S_NAMESPACE create configmap integration-test-meshsync-curl-upload-kubeconfig-script --from-file=$SCRIPT_DIR/curl-upload-kubeconfig.sh
  kubectl --namespace $MESHERY_K8S_NAMESPACE create configmap integration-test-meshsync-kubeconfig-file --from-file=kubeconfig.yaml=$TMP_KUBECONFIG_PATH
  # sometimes server not able to register k8s connection from first api call, 
  # sending request few times we ensure that connection is registered.
  # TODO (maybe) implement in more clever way then loop in a cycle.
  for i in {1..4}; do
    echo "🔁 Run #$i"

    kubectl --namespace "$MESHERY_K8S_NAMESPACE" apply -f "$SCRIPT_DIR/curl-upload-kubeconfig-job.yaml"
    kubectl --namespace "$MESHERY_K8S_NAMESPACE" wait --for=condition=complete --timeout=60s job/$JOB_NAME
    kubectl --namespace "$MESHERY_K8S_NAMESPACE" get job

    # Get the pod name for the job
    JOBS_POD_NAME=$(kubectl get pods --namespace "$MESHERY_K8S_NAMESPACE" --selector=job-name="$JOB_NAME" -o jsonpath='{.items[0].metadata.name}')
    # Output logs from the pod
    kubectl --namespace "$MESHERY_K8S_NAMESPACE" logs "$JOBS_POD_NAME"

    # delete job
    kubectl --namespace "$MESHERY_K8S_NAMESPACE" delete job "$JOB_NAME"
  done


  echo "Collecting meshsync events..."
  # TODO maybe there is a better way to be sure events are delivered?
  sleep 32

  echo "Printing server logs..." 
  # Get the pod name for the sdeployment
  SERVER_DEPLOYMENT_POD_NAME=$(kubectl get pods --namespace "$MESHERY_K8S_NAMESPACE" --selector=app.kubernetes.io/name="meshery" -o jsonpath='{.items[0].metadata.name}')
  # Output logs from the pod
  kubectl --namespace $MESHERY_K8S_NAMESPACE logs $SERVER_DEPLOYMENT_POD_NAME

  echo "Copying sqlite database file from pod..."
    NAMESPACE=$MESHERY_K8S_NAMESPACE \
    DEPLOYMENT="meshery" \
    REMOTE_PATH="/home/appuser/.meshery/config/mesherydb.sql" \
    LOCAL_PATH=$LOCAL_SQLITE_PATH \
  $SCRIPT_DIR/copy-file-from-deployment.sh
  ls -la | grep "$LOCAL_SQLITE_PATH"
  echo ""
}

cleanup_all() {
  echo "🧹 Cleaning up..."
  echo ""

  cleanup_connection
  cleanup_cluster
}

cleanup_connection () {
  if [ -f "$LOCAL_SQLITE_PATH" ]; then
    echo "Removing sqlite dataabse file from local..."
    rm "$LOCAL_SQLITE_PATH"
    echo ""
  fi

  if [ -f "$TMP_KUBECONFIG_PATH" ]; then
    echo "Removing tmp kubeconfig..."
    rm "$TMP_KUBECONFIG_PATH"
    echo ""
  fi
}

cleanup_cluster () {
  echo "Deleting KinD cluster..."
  kind delete cluster --name $CLUSTER_NAME
  echo ""

  if docker image inspect "$DOCKER_IMAGE:$DOCKER_IMAGE_TAG" > /dev/null 2>&1; then
    echo "Untagging custom image tag..."
    docker rmi $DOCKER_IMAGE:$DOCKER_IMAGE_TAG
  fi
}

print_help() {
  echo "Usage: $0 {check_dependencies|setup_all|setup_cluster|setup_connection|cleanup_all|cleanup_cluster|cleanup_connection|help}"
}

# Main dispatcher
case "$1" in
  check_dependencies)
    check_dependencies
    ;;
  setup_all)
    setup_all
    ;;
  setup_cluster)
    setup_cluster
    ;;
  setup_connection)
    setup_connection
    ;;
  cleanup_all)
    cleanup_all
    ;;
  cleanup_cluster)
    cleanup_cluster
    ;;
  cleanup_connection)
    cleanup_connection
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