# Helm

> Install Meshery on Kubernetes using Helm. Deploy Meshery in Kubernetes in-cluster.

Source: /pr-preview/pr-21670/installation/kubernetes/helm/

# Install Meshery on Kubernetes Using Helm

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li><a href="https://helm.sh/docs/intro/install/" class="meshery-light">Helm</a> should be installed on your local machine.</li>
<li>You should have access to the cluster/platform where you want to deploy Meshery.</li>
<li>Ensure that the kubeconfig file has the correct current context/cluster configuration.</li>
</ol>
</div>

## Install Meshery on Your Kubernetes Cluster Using Helm



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm repo add meshery https://meshery.io/charts/
helm install meshery meshery/meshery --namespace meshery --create-namespace</code>
	</div>
</pre>


Optionally, Meshery Server supports customizing the callback URL for your remote provider, like so:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host --create-namespace</code>
	</div>
</pre>


### Customizing Meshery's Installation with values.yaml

Meshery's Helm chart supports a number of configuration options. Please refer to the [Meshery Helm chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery#readme) and [Meshery Operator Helm Chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery-operator#readme) for more information.

#### Configuring Kubernetes Configuration Location

By default, Meshery looks for Kubernetes configuration in the `/home/appuser/.kube` directory within the container. You can customize this location by setting the `KUBECONFIG_FOLDER` environment variable:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm install meshery meshery/meshery --namespace meshery --set env.KUBECONFIG_FOLDER=/custom/path/to/.kube --create-namespace</code>
	</div>
</pre>


This is useful when providing a Meshery deployment with a predefined Kubernetes context or when using custom volume mounts for kubeconfig files.

## Upgrading Meshery with Helm

To upgrade an existing Meshery deployment:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm repo update
helm upgrade meshery meshery/meshery --namespace meshery</code>
	</div>
</pre>


For optimal upgrade performance with health check support, use the upgrade-specific values:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm upgrade meshery meshery/meshery --namespace meshery -f https://raw.githubusercontent.com/meshery/meshery/master/install/kubernetes/helm/meshery/values-upgrade.yaml --wait --timeout 10m</code>
	</div>
</pre>


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



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl get pods --namespace meshery -w</code>
	</div>
</pre>


### Checking Health Status

Verify health status with detailed information using verbose mode:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl exec --namespace meshery deployment/meshery -- curl -s &#39;http://localhost:8080/healthz/ready?verbose=1&#39;</code>
	</div>
</pre>


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



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">probe:
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
    timeoutSeconds: 3</code>
	</div>
</pre>


For comprehensive guidance on configuring health checks for different scenarios (installation, upgrades, troubleshooting), see the [Health Check Configuration Guide](https://github.com/meshery/meshery/blob/master/install/kubernetes/helm/meshery/HEALTHCHECKS.md).

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='/pr-preview/pr-21670/reference/references/mesheryctl/system/check/'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.


<h2>Accessing Meshery UI</h2>


<p>After successfully deploying Meshery, you can access Meshery's web-based user interface. Your default browser will automatically open and navigate to Meshery UI (default location is <a href="http://localhost:9081">http://localhost:9081</a>).</p>

<p>You can use the following command to open Meshery UI in your default browser:</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">mesheryctl system dashboard</code>
</div>
</pre>

<p>If you have installed Meshery on Kubernetes or a remote host, you can access Meshery UI by exposing it as a Kubernetes service or by port forwarding to Meshery UI.</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">mesheryctl system dashboard --port-forward</code>
</div>
</pre>

<p>Depending on how you have networking configured in Kubernetes, you can use kubectl to port forward to the Meshery UI.</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">kubectl port-forward svc/meshery 9081:9081 --namespace meshery</code>
</div>
</pre>

<h4>Verify Kubernetes Connection</h4>

After installing Meshery, regardless of the installation type, it is important to verify that your kubeconfig file has been uploaded correctly via the UI. 

<ol>
<li>In the Meshery UI, navigate to <strong>Lifecycle</strong> from the menu on the left.</li>
<li>Click on Connections.</li>
<li>Ensure that your cluster appears in the list of connections and is marked as <code>Connected</code>.</li>
<li>Click on the cluster name to perform a ping test and confirm that Meshery can communicate with your cluster.</li>
</ol>

<details>
<summary>Customizing Your Meshery Provider Callback URL</summary>

<p>
  Meshery Server supports customizing your <a href="/pr-preview/pr-21670/reference/extensibility/providers/">Meshery Provider</a> authentication flow callback URL. This is helpful when deploying Meshery behind multiple layers of networking infrastructure.
</p>

<p>
  For production deployments, it is recommended to access the Meshery UI by setting up a reverse proxy or using a LoadBalancer. By specifying a custom redirect endpoint, you can ensure that authentication flows complete successfully, even when multiple routing layers are involved.
</p>

<p>
  <b>Note</b>: For production deployments, it is important to select the <code>Remote Provider</code> in order to control which identity providers are authorized. Learn more about this in the <a href="/pr-preview/pr-21670/reference/extensibility/providers/">Extensibility: Providers</a> guide.
</p>

<p>
  Define a custom callback URL by setting up the <code>MESHERY_SERVER_CALLBACK_URL</code> environment variable before installing Meshery.
</p>

<p>
  To customize the authentication flow callback URL, use the following command:
</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start</code>
</div>
</pre>

<p>
  Meshery should now be running in your Kubernetes cluster and the Meshery UI should be accessible at the <code>EXTERNAL IP</code> of the <code>meshery</code> service.
</p>

</details>


<div class="related-discussions">
  <h3>Recent Discussions with "meshery" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/design-meshery-mcp-server-architecture-and-registration-interface/7954" target="_blank" rel="noopener noreferrer">
            Design: Meshery MCP Server architecture and registration interface
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/the-meshery-mcp-server-foundation-is-up-lets-agree-on-what-to-build-next/7952" target="_blank" rel="noopener noreferrer">
            The Meshery MCP Server foundation is up, let&#39;s agree on what to build next
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/new-intro-topic/7975" target="_blank" rel="noopener noreferrer">
            New intro topic
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-mcp-server-poc-an-ai-agent-managing-kubernetes-through-mcp/7974" target="_blank" rel="noopener noreferrer">
            Meshery MCP Server POC: an AI agent managing Kubernetes through MCP
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/approach-for-context-window-aware-retrieval-in-the-ai-adapter-issue-20994/7963" target="_blank" rel="noopener noreferrer">
            Approach for context-window-aware retrieval in the AI Adapter (Issue #20994)
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/there-is-some-error-in-running-the-localhost-of-layer5-in-my-server/7818" target="_blank" rel="noopener noreferrer">
            There is some error in running the localhost of layer5 in my server
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/rfc-aligning-the-foundation-for-the-meshery-mcp-server/7913" target="_blank" rel="noopener noreferrer">
            RFC: Aligning the Foundation for the Meshery MCP Server
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-development-meeting-august-12th-2026/7948" target="_blank" rel="noopener noreferrer">
            Meshery Development Meeting | August 12th, 2026
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
