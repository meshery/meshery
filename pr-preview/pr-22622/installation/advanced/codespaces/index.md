# Codespaces

> Build and contribute to Meshery using GitHub Codespaces

---

LLMS index: [llms.txt](/pr-preview/pr-22622/llms.txt)

---

Use Minikube in GitHub Codespaces to set up a development environment for Meshery.

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Install the Meshery command line client, <a href="/pr-preview/pr-22622/installation/mesheryctl/" class="meshery-light">mesheryctl</a>.</li>
  </ol>
</div>

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Installation: Manual Steps](#installation-manual-steps)
- [Post-Installation Steps](#post-installation-steps)

# In-cluster Installation

Follow these steps to install Meshery in your Minikube cluster.

## Preflight Checks

Before you deploy Meshery on Minikube, check the following.

### Preflight: Cluster Connectivity


You can develop and run Meshery in GitHub Codespaces using your choice of tool:

- A command shell, via an SSH connection initiated using GitHub CLI.
- One of the JetBrains IDEs, via the JetBrains Gateway.
- The Visual Studio Code desktop application.
- A browser-based version of Visual Studio Code.

<div class="alert alert-dark" role="alert"><div class="h4 alert-heading" role="heading">Choice of Codespaces Tool</div>

For the best experience, run Codespaces in your locally <a href='https://docs.github.com/en/codespaces/developing-in-codespaces/developing-in-a-codespace'>installed IDE</a>. Alternatively, you can <br /><a href='https://github.com/codespaces/new?hide_repo_select=true&ref=master&repo=157554479&machine=premiumLinux'><img alt='Open in GitHub Codespaces' src='https://github.com/codespaces/badge.svg' /></a></div>


Start Minikube if it is not already running:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">minikube start --cpus 4 --memory 4096</code>
	</div>
</pre>


Set `--cpus` and `--memory` to fit the machine type you chose in GitHub Codespaces. To check the status of your Minikube cluster:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">minikube status</code>
	</div>
</pre>


Verify that `kubectl` can reach the cluster:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl cluster-info</code>
	</div>
</pre>


## Installation: Using `mesheryctl`

Start Meshery in your Minikube cluster:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system start -p kubernetes</code>
	</div>
</pre>


Meshery Server running in the cluster connects to it automatically. If the connection does not appear, sign in with `mesheryctl system login` and create it with `mesheryctl connection create --type minikube`, which uploads your kubeconfig to Meshery Server. (`mesheryctl system config minikube` is deprecated.) For more on signing in, see [Authenticate with Meshery via CLI](/pr-preview/pr-22622/guides/mesheryctl/authenticate-with-meshery-via-cli/).

## Installation: Using Helm

To install Meshery with Helm 3, see the [Helm Installation](/pr-preview/pr-22622/installation/kubernetes/helm/) guide.

## Installation: Manual Steps

You can also generate a self-contained kubeconfig for your Minikube cluster and upload it to Meshery yourself. Export the `minikube` context, with its certificates embedded:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl config view --minify --flatten --context minikube &gt; config_minikube.yaml</code>
	</div>
</pre>


The file looks like this, with your own server address and certificate data:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: &lt;certificate data&gt;
    server: https://192.168.49.2:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
preferences: {}
users:
- name: minikube
  user:
    client-certificate-data: &lt;certificate data&gt;
    client-key-data: &lt;key data&gt;</code>
	</div>
</pre>


_Note_: Make sure that `current-context` is `minikube`.

Then upload `config_minikube.yaml` in Meshery UI: go to **Lifecycle** > **Connections**, click **Create Connection**, choose **Kubernetes**, and upload the file in the **Import Kubeconfig** step. See [Registering a Connection](/pr-preview/pr-22622/guides/infrastructure-management/registering-a-connection/). Take a look at the [Meshery guides](/pr-preview/pr-22622/guides/) for advanced usage tips.

# Post-Installation Steps

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


To reach Meshery UI from your browser, forward its port from the codespace. See [Forwarding ports in your codespace](https://docs.github.com/en/codespaces/developing-in-a-codespace/forwarding-ports-in-your-codespace?tool=vscode).

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
