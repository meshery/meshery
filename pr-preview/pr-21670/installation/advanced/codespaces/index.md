# Codespaces

> Build and contribute to Meshery using GitHub Codespaces

Source: /pr-preview/pr-21670/installation/advanced/codespaces/

Use Minikube in GitHub Codespaces to setup your development environment for Meshery.

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Install the Meshery command line client, <a href="/pr-preview/pr-21670/installation/mesheryctl/" class="meshery-light">mesheryctl</a>.</li>
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

Follow the steps below to install Meshery in your Minikube cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on Minikube.

### Preflight: Cluster Connectivity


You can develop and run Meshery in GitHub Codespaces using your choice of tool:

- A command shell, via an SSH connection initiated using GitHub CLI.
- One of the JetBrains IDEs, via the JetBrains Gateway.
- The Visual Studio Code desktop application.
- A browser-based version of Visual Studio Code.

<div class="alert alert-dark" role="alert"><div class="h4 alert-heading" role="heading">Choice of Codespaces Tool</div>

For the best experience, run Codespaces in your locally <a href='https://docs.github.com/en/codespaces/developing-in-codespaces/developing-in-a-codespace'>installed IDE</a>. Alternatively, you can <br /><a href='https://github.com/codespaces/new?hide_repo_select=true&ref=master&repo=157554479&machine=premiumLinux'><img alt='Open in GitHub Codespaces' src='https://github.com/codespaces/badge.svg' /></a></div>


Start Minikube, if it is not already running, using the following command:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">minikube start --cpus 4 --memory 4096</code>
	</div>
</pre>


Please allocate CPUs based on the machine you selected in GitHub Codespaces. To check the status of your Minikube cluster:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">minikube status</code>
	</div>
</pre>


Verify your kubeconfig's current context.



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl cluster-info</code>
	</div>
</pre>


## Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your Minikube cluster. Configure Meshery to connect to your Minikube cluster by executing:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system config minikube</code>
	</div>
</pre>


Once configured, execute the following command to start Meshery.



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system start</code>
	</div>
</pre>


If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/pr-preview/pr-21670/guides/mesheryctl/authenticate-with-meshery-via-cli/) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/pr-preview/pr-21670/installation/kubernetes/helm/) guide.

## Installation: Manual Steps

You may also manually generate and load the kubeconfig file for Meshery to use:

**The following configuration yaml will be used by Meshery. Copy and paste the following in your config file** :



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">apiVersion: v1
clusters:

- cluster:
  certificate-authority-data: &lt; cert shortcutted &gt;
  server: https://192.168.99.100:8443
  name: minikube
  contexts:
- context:
  cluster: minikube
  user: minikube
  name: minikube
  current-context: minikube
  kind: Config
  preferences: {}
  users:
- name: minikube
  user:
  client-certificate-data: &lt; cert shortcutted &gt;
  client-key-data: &lt; key shortcutted &gt;</code>
	</div>
</pre>


_Note_: Make sure _current-context_ is set to _minikube_.

<br />
**To allow Meshery to auto detect your config file, Run** :



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl config view --minify --flatten &gt; config_minikube.yaml</code>
	</div>
</pre>


<br />
Meshery should now be connected with your managed Kubernetes instance. Take a look at the [Meshery guides](/pr-preview/pr-21670/guides/) for advanced usage tips.

# Post-Installation Steps

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


For further information to access Meshery UI or port forwarding in GitHub Codespaces, read the [docs](https://docs.github.com/en/codespaces/developing-in-a-codespace/forwarding-ports-in-your-codespace?tool=vscode)

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
