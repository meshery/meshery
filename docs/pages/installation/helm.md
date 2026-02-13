---
layout: default
title: Helm
permalink: installation/kubernetes/helm
type: installation
category: kubernetes
redirect_from:
- installation/platforms/helm
display-title: "false"
language: en
list: include
abstract: Install Meshery on Kubernetes using Helm. Deploy Meshery in Kubernetes in-cluster.
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

{% capture code_content %}helm repo add meshery https://meshery.io/charts/
helm install meshery meshery/meshery --namespace meshery --create-namespace
{% endcapture %}
{% include code.html code=code_content %}

Optionally, Meshery Server supports customizing the callback URL for your remote provider, like so:

{% capture code_content %}helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host --create-namespace{% endcapture %}
{% include code.html code=code_content %}

### Customizing Meshery's Installation with values.yaml

Meshery's Helm chart supports a number of configuration options. Please refer to the [Meshery Helm chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery#readme) and [Meshery Operator Helm Chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery-operator#readme) for more information.

#### Configuring Kubernetes Configuration Location

By default, Meshery looks for Kubernetes configuration in the `/home/appuser/.kube` directory within the container. You can customize this location by setting the `KUBECONFIG_FOLDER` environment variable:

{% capture code_content %}helm install meshery meshery/meshery --namespace meshery \
  --set env.KUBECONFIG_FOLDER=/custom/path/to/.kube \
  --create-namespace
{% endcapture %}
{% include code.html code=code_content %}

This is useful when providing a Meshery deployment with a predefined Kubernetes context or when using custom volume mounts for kubeconfig files.

## Upgrading Meshery with Helm

To upgrade an existing Meshery deployment:

{% capture code_content %}helm repo update
helm upgrade meshery meshery/meshery --namespace meshery
{% endcapture %}
{% include code.html code=code_content %}

For optimal upgrade performance with health check support, use the upgrade-specific values:

{% capture code_content %}helm upgrade meshery meshery/meshery --namespace meshery \
  -f https://raw.githubusercontent.com/meshery/meshery/master/install/kubernetes/helm/meshery/values-upgrade.yaml \
  --wait --timeout 10m
{% endcapture %}
{% include code.html code=code_content %}

The upgrade configuration includes:
- **Startup probes** to protect pods during initialization
- **Optimized probe timing** for capability reloading
- **Higher failure thresholds** to tolerate temporary unavailability during upgrades

See the [Health Check Configuration Guide](https://github.com/meshery/meshery/blob/master/install/kubernetes/helm/meshery/HEALTHCHECKS.md) for detailed information.

## Health Checks and Monitoring

Meshery implements Kubernetes-compliant health check endpoints that follow best practices from the Kubernetes API server:

- **Liveness probe** (`/healthz/live`) - Checks if Meshery is running and responsive
- **Readiness probe** (`/healthz/ready`) - Checks if Meshery is ready to accept traffic

### Monitoring Deployment Status

Monitor the status of your Meshery deployment:

{% capture code_content %}kubectl get pods --namespace meshery -w
{% endcapture %}
{% include code.html code=code_content %}

### Checking Health Status

Verify health status with detailed information using verbose mode:

{% capture code_content %}kubectl exec --namespace meshery deployment/meshery -- \
  curl -s "http://localhost:8080/healthz/ready?verbose=1"
{% endcapture %}
{% include code.html code=code_content %}

**Example output:**
```
[+]capabilities ok
[i]extension extension package found
healthz check passed
```

**Legend:**
- `[+]` - Health check passed
- `[-]` - Health check failed (causes pod to be marked unhealthy)
- `[i]` - Informational status (does not affect health)

### Health Check Configuration

The Helm chart includes pre-configured health checks with sensible defaults:

- **Liveness probe**: Initial delay of 80 seconds to allow for server startup and provider initialization
- **Readiness probe**: Initial delay of 10 seconds with frequent checks for faster readiness detection
- **Startup probe**: Optional (disabled by default) for handling slow-starting containers

To customize probe settings, modify your `values.yaml`:

{% capture code_content %}probe:
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
    timeoutSeconds: 3
{% endcapture %}
{% include code.html code=code_content %}

For comprehensive guidance on configuring health checks for different scenarios (installation, upgrades, troubleshooting), see the [Health Check Configuration Guide](https://github.com/meshery/meshery/blob/master/install/kubernetes/helm/meshery/HEALTHCHECKS.md).

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{% include_cached installation/accessing-meshery-ui.md display-title="true" %}

{% include related-discussions.html tag="meshery" %}
