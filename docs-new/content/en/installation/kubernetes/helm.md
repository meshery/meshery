---
title: "Helm"
description: "Install Meshery on Kubernetes using Helm. Deploy Meshery in-cluster."
weight: 10
aliases:
  - /installation/platforms/helm
display_title: "false"
---

# Install Meshery on Kubernetes Using Helm

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li><a href="https://helm.sh/docs/intro/install/" class="meshery-light">Helm</a> should be installed on your local machine.</li>
<li>You should have access to the cluster/platform where you want to deploy Meshery.</li>
<li>Ensure that the kubeconfig file has the correct current context/cluster configuration.</li>
</ol>
</div>

## Install Meshery on Your Kubernetes Cluster Using Helm

{{< code >}}
helm repo add meshery https://meshery.io/charts/
helm install meshery meshery/meshery --namespace meshery --create-namespace
{{< /code >}}

Optionally, Meshery Server supports customizing the callback URL for your remote provider:

{{< code >}}
helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host --create-namespace
{{< /code >}}

### Customizing Meshery's Installation with values.yaml

Meshery's Helm chart supports a number of configuration options. Please refer to the [Meshery Helm chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery#readme) and [Meshery Operator Helm Chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery-operator#readme) for more information.

#### Configuring Kubernetes Configuration Location

By default, Meshery looks for Kubernetes configuration in the `/home/appuser/.kube` directory within the container. You can customize this location:

{{< code >}}
helm install meshery meshery/meshery --namespace meshery \
  --set env.KUBECONFIG_FOLDER=/custom/path/to/.kube \
  --create-namespace
{{< /code >}}

## Upgrading Meshery with Helm

To upgrade an existing Meshery deployment:

{{< code >}}
helm repo update
helm upgrade meshery meshery/meshery --namespace meshery
{{< /code >}}

For optimal upgrade performance with health check support:

{{< code >}}
helm upgrade meshery meshery/meshery --namespace meshery \
  -f https://raw.githubusercontent.com/meshery/meshery/master/install/kubernetes/helm/meshery/values-upgrade.yaml \
  --wait --timeout 10m
{{< /code >}}

## Health Checks and Monitoring

Meshery implements Kubernetes-compliant health check endpoints:

- **Liveness probe** (`/healthz/live`) - Checks if Meshery is running and responsive
- **Readiness probe** (`/healthz/ready`) - Checks if Meshery is ready to accept traffic

### Monitoring Deployment Status

{{< code >}}
kubectl get pods --namespace meshery -w
{{< /code >}}

### Checking Health Status

{{< code >}}
kubectl exec --namespace meshery deployment/meshery -- \
  curl -s "http://localhost:8080/healthz/ready?verbose=1"
{{< /code >}}

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using [mesheryctl system check](/reference/mesheryctl/system/check).

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}
