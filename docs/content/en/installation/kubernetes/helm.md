---
title: Helm
categories: [kubernetes]
aliases:
- /installation/platforms/helm
display_title: false
description: Install Meshery on Kubernetes using Helm. Deploy Meshery in Kubernetes in-cluster.
---
# Install Meshery on Kubernetes Using Helm

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install <a href="https://helm.sh/docs/intro/install/" class="meshery-light">Helm</a> on your local machine.</li>
<li>Access to the cluster where you want to deploy Meshery.</li>
<li>A kubeconfig whose current context points at that cluster.</li>
</ol>
</div>

## Install Meshery on Your Kubernetes Cluster Using Helm

{{< code code=`helm repo add meshery https://meshery.io/charts/
helm install meshery meshery/meshery --namespace meshery --create-namespace` >}}

Optionally, customize the callback URL that your Remote Provider redirects to after sign-in:

{{< code code="helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host --create-namespace" >}}

### Customizing Meshery's Installation with values.yaml

Meshery's Helm chart supports many configuration options. See the [Meshery Helm chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery#readme) and [Meshery Operator Helm Chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery-operator#readme) for more information.

#### Configuring Kubernetes Configuration Location

By default, Meshery Server looks for Kubernetes configuration in the `/home/appuser/.kube` directory inside the container. To change this location, set the `KUBECONFIG_FOLDER` environment variable:

{{< code code="helm install meshery meshery/meshery --namespace meshery --set env.KUBECONFIG_FOLDER=/custom/path/to/.kube --create-namespace" >}}

A custom location helps when you give a Meshery deployment a predefined Kubernetes context or mount kubeconfig files from a custom volume.

## Upgrading Meshery with Helm

To upgrade an existing Meshery deployment:

{{< code code=`helm repo update
helm upgrade meshery meshery/meshery --namespace meshery` >}}

For upgrades, use the upgrade-specific values file, which tunes the health checks for a rolling upgrade:

{{< code code=`helm upgrade meshery meshery/meshery --namespace meshery -f https://raw.githubusercontent.com/meshery/meshery/master/install/kubernetes/helm/meshery/values-upgrade.yaml --wait --timeout 10m` >}}

The upgrade configuration includes:
- **Startup probes** to protect pods during initialization
- **Optimized probe timing** for capability reloading
- **Higher failure thresholds** to tolerate temporary unavailability during upgrades

For details, see the [Health Check Configuration Guide](https://github.com/meshery/meshery/blob/master/install/kubernetes/helm/meshery/HEALTHCHECKS.md).

## Health Checks and Monitoring

Meshery Server exposes health check endpoints modeled on the Kubernetes API server's:

- **Liveness probe** (`/healthz/live`): checks that Meshery is running and responsive.
- **Readiness probe** (`/healthz/ready`): checks that Meshery is ready to accept traffic.

### Monitoring Deployment Status

Watch the status of your Meshery pods:

{{< code code="kubectl get pods --namespace meshery -w" >}}

### Checking Health Status

To see each health check's result, query the readiness endpoint in verbose mode. The Meshery Server image does not include `curl`, so forward a local port to the deployment first:

{{< code code=`kubectl port-forward --namespace meshery deployment/meshery 8080:8080` >}}

Then, in another terminal:

{{< code code=`curl -s 'http://localhost:8080/healthz/ready?verbose=1'` >}}

**Example output:**
```
[+]capabilities ok
[i]extension extension package found
healthz check passed
```

**Legend:**
- `[+]`: the check passed.
- `[-]`: the check failed, and Kubernetes marks the pod unhealthy.
- `[i]`: informational only; it does not affect health.

### Health Check Configuration

The Helm chart configures health checks with these defaults:

- **Liveness probe**: an initial delay of 80 seconds, to allow for server startup and provider initialization.
- **Readiness probe**: an initial delay of 10 seconds, with frequent checks to detect readiness quickly.
- **Startup probe**: optional and disabled by default, for slow-starting containers.

To customize the probes, set them in your `values.yaml`:

{{< code code=`probe:
  livenessProbe:
    enabled: true
    initialDelaySeconds: 80
    periodSeconds: 12
    failureThreshold: 4
    timeoutSeconds: 5

  readinessProbe:
    enabled: true
    initialDelaySeconds: 10
    periodSeconds: 4
    failureThreshold: 4
    timeoutSeconds: 3` >}}

For guidance on configuring health checks for installation, upgrades, and troubleshooting, see the [Health Check Configuration Guide](https://github.com/meshery/meshery/blob/master/install/kubernetes/helm/meshery/HEALTHCHECKS.md).

## Post-Installation Steps

Optionally, verify the health of your Meshery deployment with <a href='{{< ref "reference/references/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

You're ready to use Meshery. Open Meshery UI in your browser.

{{< installation/accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}
