# GKE

> Install Meshery on Google Kubernetes Engine. Deploy Meshery in GKE in-cluster or outside of GKE out-of-cluster.

---

LLMS index: [llms.txt](/pr-preview/pr-22622/llms.txt)

---

<h1>Quick Start with GKE <img src="images/gke.png" style="width:35px;height:35px;" /></h1>

Manage your GKE clusters with Meshery. Deploy Meshery inside your GKE cluster ([in-cluster](#in-cluster-installation)), or run it outside of GKE ([out-of-cluster](/pr-preview/pr-22622/installation/kubernetes/#out-of-cluster-installation)) and connect it to your cluster. **_Note: We recommend installing Meshery in your GKE clusters._**

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Install the Meshery command line client, <a href="/pr-preview/pr-22622/installation/mesheryctl/" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Install <a href="https://cloud.google.com/sdk/docs/install">gcloud CLI</a>, configured for your environment.</li>
    <li>Access to an active GKE cluster in your Google Cloud project.</li>
  </ol>
</div>

Also see: [Install Meshery on Kubernetes](/pr-preview/pr-22622/installation/kubernetes/)

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
    - [Preflight: Plan your access to Meshery UI](#preflight-plan-your-access-to-meshery-ui)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)

# In-cluster Installation

Follow the steps below to install Meshery in your GKE cluster.

## Preflight Checks

Before you deploy Meshery on GKE, check the following.

### Preflight: Cluster Connectivity

Use the gcloud CLI to connect `kubectl` to your Google Kubernetes Engine (GKE) cluster. GKE kubeconfigs authenticate through `gke-gcloud-auth-plugin`, so [install the plugin](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl#install_plugin) first.

1. Sign in to your Google Cloud account with [gcloud auth login](https://cloud.google.com/sdk/gcloud/reference/auth/login).
2. Set the project that contains your cluster:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">gcloud config set project [PROJECT_ID]</code>
	</div>
</pre>

3. Download the cluster's credentials and make it your current context:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">gcloud container clusters get-credentials [CLUSTER_NAME] --location [CLUSTER_LOCATION]</code>
	</div>
</pre>

4. Verify your kubeconfig's current context:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl config current-context</code>
	</div>
</pre>


### Preflight: Plan your access to Meshery UI

1. To reach Meshery UI through port-forwarding, see the [port-forwarding](/pr-preview/pr-22622/reference/mesheryctl/system/dashboard/) guide.
2. To reach it through a load balancer, see [Exposing Meshery Service (LoadBalancer)](/pr-preview/pr-22622/installation/kubernetes/#exposing-meshery-serviceloadbalancer).
3. If users reach Meshery through a custom host name, set the callback URL that your Provider redirects to after sign-in when you start Meshery:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start -p kubernetes</code>
	</div>
</pre>


With a load balancer, Meshery UI is then available at `http://[EXTERNAL-IP]:9081`, where `EXTERNAL-IP` is the external address of the `meshery` service.

## Installation: Using `mesheryctl`

With your GKE cluster as the current context, start Meshery on Kubernetes:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system start -p kubernetes</code>
	</div>
</pre>


Meshery Server runs in the cluster and connects to it automatically, so you do not need to upload a kubeconfig. To connect a Meshery Server that runs elsewhere to this cluster, sign in with `mesheryctl system login` and run `mesheryctl connection create --type gke`. (`mesheryctl system config gke` is deprecated.)

If `mesheryctl` reports an authentication error, sign in with `mesheryctl system login`. See [Authenticate with Meshery via CLI](/pr-preview/pr-22622/guides/mesheryctl/authenticate-with-meshery-via-cli/).

## Installation: Using Helm

To install Meshery with Helm 3, see the [Helm Installation](/pr-preview/pr-22622/installation/kubernetes/helm/) guide.

## Post-Installation Steps

Optionally, verify the health of your Meshery deployment with <a href='/pr-preview/pr-22622/reference/mesheryctl/system/check/'>mesheryctl system check</a>.

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
  You can customize the callback URL that your <a href="/pr-preview/pr-22622/reference/extensibility/providers/">Meshery Provider</a> redirects to after sign-in. This helps when Meshery runs behind several layers of networking infrastructure.
</p>

<p>
  For production deployments, we recommend that you expose Meshery UI through a reverse proxy or a load balancer. A custom callback URL makes sure that sign-in completes even when requests pass through several routing layers.
</p>

<p>
  <b>Note</b>: For production deployments, use a Remote Provider so that you control which identity providers can sign users in. Learn more about this in the <a href="/pr-preview/pr-22622/reference/extensibility/providers/">Extensibility: Providers</a> guide.
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
        </li><li>
          <a href="https://discuss.meshery.io/t/play-meshery-io-is-down/8045" target="_blank" rel="noopener noreferrer">
            play.meshery.io is DOWN
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/your-path-to-becoming-a-certified-meshery-contributor-starts-here/7652" target="_blank" rel="noopener noreferrer">
            Your Path to Becoming a Certified Meshery Contributor Starts Here
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/design-meshery-mcp-server-architecture-and-registration-interface/7954" target="_blank" rel="noopener noreferrer">
            Design: Meshery MCP Server architecture and registration interface
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
