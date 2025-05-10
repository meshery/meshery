  # echo "Port forward port to broker..."
  # kubectl --namespace $MESHERY_K8S_NAMESPACE port-forward svc/meshery-broker 32043:4222  &
  # PORT_FORWARD_PID=$!
  # echo "Port forward PID=$PORT_FORWARD_PID"

  # kubectl config view --minify --raw > kubeconfig-current.yaml

  # # port forward
  # kubectl --namespace meshery port-forward svc/meshery-broker 32043:4222