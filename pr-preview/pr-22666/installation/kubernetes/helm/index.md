# Helm

> Install Meshery on Kubernetes using Helm. Deploy Meshery in Kubernetes in-cluster. This guide applies to all cloud providers including AWS (EKS), Azure (AKS), GCP (GKE), and local clusters (Minikube, KinD).

---

LLMS index: [llms.txt](/pr-preview/pr-22666/llms.txt)

---

# Install Meshery on Kubernetes Using Helm

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install <a href="https://helm.sh/docs/intro/install/" class="meshery-light">Helm</a> and <a href="https://kubernetes.io/docs/tasks/tools/" class="meshery-light">kubectl</a> on your local machine.</li>
<li>Access to the cluster where you want to deploy Meshery.</li>
<li>A kubeconfig whose current context points at that cluster.</li>
</ol>
</div>

## Install Meshery on Your Kubernetes Cluster Using Helm



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm repo add meshery https://meshery.io/charts/
helm install meshery meshery/meshery --namespace meshery --create-namespace</code>
	</div>
</pre>


Optionally, customize the callback URL that your Remote Provider redirects to after sign-in:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host --create-namespace</code>
	</div>
</pre>


### Customizing Meshery's Installation with values.yaml

Meshery's Helm chart supports many configuration options. See the [Meshery Helm chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery#readme) and [Meshery Operator Helm Chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery-operator#readme) for more information.

#### Configuring Kubernetes Configuration Location

By default, Meshery Server looks for Kubernetes configuration in the `/home/appuser/.kube` directory inside the container. To change this location, set the `KUBECONFIG_FOLDER` environment variable:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm install meshery meshery/meshery --namespace meshery --set env.KUBECONFIG_FOLDER=/custom/path/to/.kube --create-namespace</code>
	</div>
</pre>


A custom location helps when you give a Meshery deployment a predefined Kubernetes context or mount kubeconfig files from a custom volume.

## Upgrading Meshery with Helm

To upgrade an existing Meshery deployment:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm repo update
helm upgrade meshery meshery/meshery --namespace meshery</code>
	</div>
</pre>


Alternatively, run this command instead to use the upgrade-specific values file, which tunes the health checks for a rolling upgrade:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">helm upgrade meshery meshery/meshery --namespace meshery -f https://raw.githubusercontent.com/meshery/meshery/master/install/kubernetes/helm/meshery/values-upgrade.yaml --wait --timeout 10m</code>
	</div>
</pre>


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



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl get pods --namespace meshery -w</code>
	</div>
</pre>


### Checking Health Status

To see each health check's result, query the readiness endpoint in verbose mode. The Meshery Server image does not include `curl`, so forward a local port to the deployment first:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl port-forward --namespace meshery deployment/meshery 8080:8080</code>
	</div>
</pre>


Then, in another terminal:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">curl -s &#39;http://localhost:8080/healthz/ready?verbose=1&#39;</code>
	</div>
</pre>


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


For guidance on configuring health checks for installation, upgrades, and troubleshooting, see the [Health Check Configuration Guide](https://github.com/meshery/meshery/blob/master/install/kubernetes/helm/meshery/HEALTHCHECKS.md).

## Post-Installation Steps

Optionally, verify the health of your Meshery deployment with <a href='/pr-preview/pr-22666/reference/mesheryctl/system/check/'>mesheryctl system check</a>.

You're ready to use Meshery. Open Meshery UI in your browser.


<h2>Accessing Meshery UI</h2>


<p>After successfully deploying Meshery, you can access Meshery's web-based user interface. Your default browser will automatically open and navigate to Meshery UI (default location is <a href="http://localhost:9081">http://localhost:9081</a>).</p>

<p>You can use the following command to open Meshery UI in your default browser:</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">mesheryctl system dashboard</code>
</div>
</pre>

<p>If you have installed Meshery on Kubernetes or a remote host, you can access Meshery UI by exposing it as a Kubernetes service or by port forwarding to Meshery UI.</p>

<p>For detailed instructions on accessing Meshery UI, see the <a href="/pr-preview/pr-22666/installation/kubernetes/#accessing-meshery-ui">Accessing Meshery UI</a> section in the Kubernetes installation guide.</p>

<p><strong>Quick options:</strong></p>
<ul>
  <li><strong>Port forwarding (recommended for local access):</strong> <code>mesheryctl system dashboard --port-forward</code></li>
  <li><strong>LoadBalancer service (recommended for production):</strong> See the <a href="/pr-preview/pr-22666/installation/kubernetes/#accessing-meshery-ui">LoadBalancer section</a> for details</li>
</ul>

<p>For comprehensive port forwarding instructions including troubleshooting, see the <a href="/pr-preview/pr-22666/installation/kubernetes/port-forwarding/">Port Forwarding Guide</a>.</p>

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
  Meshery Server supports customizing your <a href="/pr-preview/pr-22666/reference/extensibility/providers/">Meshery Provider</a> authentication flow callback URL. This is helpful when deploying Meshery behind multiple layers of networking infrastructure.
</p>

<p>
  For production deployments, it is recommended to access the Meshery UI by setting up a reverse proxy or using a LoadBalancer. By specifying a custom redirect endpoint, you can ensure that authentication flows complete successfully, even when multiple routing layers are involved.
</p>

<p>
  <b>Note</b>: For production deployments, it is important to select the <code>Remote Provider</code> in order to control which identity providers are authorized. Learn more about this in the <a href="/pr-preview/pr-22666/reference/extensibility/providers/">Extensibility: Providers</a> guide.
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
          <a href="https://discuss.meshery.io/t/new-contributor-intro-vedit-agrawal-bughunterx2101-looking-to-pair-with-a-meshmate-for-meshery-ui-ai-adapter/8172" target="_blank" rel="noopener noreferrer">
            New Contributor Intro: Vedit Agrawal (@BugHunterX2101) | Looking to pair with a MeshMate for Meshery UI &amp; AI Adapter
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-seeking-an-experienced-meshery-contributor-as-a-guide/8166" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate: Seeking an Experienced Meshery Contributor as a Guide
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-to-reactivate-my-slack-account/8168" target="_blank" rel="noopener noreferrer">
            Request to Reactivate My Slack Account
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-frontend-backend-cloud-native-development/8161" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate — Frontend, Backend &amp; Cloud-Native Development
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-for-a-meshmate-active-meshery-contributor/8158" target="_blank" rel="noopener noreferrer">
            Request for a MeshMate — Active Meshery Contributor
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/bug-new-catalog-publish-request-email-renders-with-oversized-social-icons-and-squashed-content-on-mobile/8138" target="_blank" rel="noopener noreferrer">
            [Bug] &#34;New Catalog publish request&#34; email renders with oversized social icons and squashed content on mobile
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/how-to-present-an-issue-pr-in-a-development-meeting-feedback-wanted/8107" target="_blank" rel="noopener noreferrer">
            How to Present an Issue/PR in a Development Meeting: feedback wanted
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/newbie-here-excited-to-join/8052" target="_blank" rel="noopener noreferrer">
            Newbie here -- Excited to join
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
