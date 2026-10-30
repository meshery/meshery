# Helm

> Install Meshery on Kubernetes using Helm. Deploy Meshery in Kubernetes in-cluster.

---

LLMS index: [llms.txt](/pr-preview/pr-22676/llms.txt)

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

Optionally, verify the health of your Meshery deployment with <a href='/pr-preview/pr-22676/reference/mesheryctl/system/check/'>mesheryctl system check</a>.

You're ready to use Meshery. Open Meshery UI in your browser.


<h2>Accessing Meshery UI</h2>


<p>Once Meshery is running, <code>mesheryctl system start</code> opens Meshery UI in your default browser. On Docker, Meshery UI listens at <a href="http://localhost:9081">http://localhost:9081</a> by default. On Kubernetes, <code>mesheryctl</code> opens the address of the <code>meshery</code> service in the <code>meshery</code> namespace, which the Helm chart creates as a <code>LoadBalancer</code> service on port 9081.</p>

<p>To open Meshery UI in your default browser at any time, run:</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">mesheryctl system dashboard</code>
</div>
</pre>

<p>If your cluster does not give the <code>meshery</code> service a reachable address, forward a local port to Meshery instead. By default, <code>mesheryctl</code> listens on local port 9081; use <code>-p</code> to choose another port.</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">mesheryctl system dashboard --port-forward</code>
</div>
</pre>

<p>You can also port-forward with <code>kubectl</code>:</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">kubectl port-forward svc/meshery 9081:9081 --namespace meshery</code>
</div>
</pre>

<h4>Verify Kubernetes Connection</h4>

Whichever way you installed Meshery, verify that Meshery is connected to your cluster:

<ol>
<li>In the Meshery UI, navigate to <strong>Lifecycle</strong> from the menu on the left.</li>
<li>Select <strong>Connections</strong>.</li>
<li>Ensure that your cluster appears in the list of connections and is marked as <code>Connected</code>.</li>
<li>Click the cluster name to run a connectivity test and confirm that Meshery can reach your cluster.</li>
</ol>

<details>
<summary>Customizing Your Meshery Provider Callback URL</summary>

<p>
  You can customize the callback URL that your <a href="/pr-preview/pr-22676/reference/extensibility/providers/">Meshery Provider</a> redirects to after sign-in. This helps when Meshery runs behind several layers of networking infrastructure.
</p>

<p>
  For production deployments, we recommend that you expose Meshery UI through a reverse proxy or a load balancer. A custom callback URL makes sure that sign-in completes even when requests pass through several routing layers.
</p>

<p>
  <b>Note</b>: For production deployments, use a Remote Provider so that you control which identity providers can sign users in. Learn more about this in the <a href="/pr-preview/pr-22676/reference/extensibility/providers/">Extensibility: Providers</a> guide.
</p>

<p>
  To set a custom callback URL, set the <code>MESHERY_SERVER_CALLBACK_URL</code> environment variable when you start Meshery:
</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start</code>
</div>
</pre>

<p>
  On Kubernetes, Meshery UI is then available at the <code>EXTERNAL-IP</code> of the <code>meshery</code> service. Run <code>kubectl get svc meshery -n meshery</code> to find it.
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
