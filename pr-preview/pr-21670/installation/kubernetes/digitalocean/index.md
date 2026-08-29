# DigitalOcean

> Install Meshery on DigitalOcean. Deploy Meshery out-of-cluster with Docker on a Droplet, or in-cluster on DigitalOcean Kubernetes (DOKS).

Source: /pr-preview/pr-21670/installation/kubernetes/digitalocean/

<h1>Quick Start with DigitalOcean <img src="images/digitalocean.svg" style="width:35px;height:35px;" /></h1>

Deploy and manage your DigitalOcean infrastructure with Meshery. You can run Meshery on DigitalOcean in two ways: [out-of-cluster](#option-1-docker-on-a-droplet-out-of-cluster) using Docker on a Droplet, or [in-cluster](#option-2-digitalocean-kubernetes-in-cluster) on a DigitalOcean Kubernetes (DOKS) cluster. **_Note: It is advisable to install Meshery in your DOKS cluster._**

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install the Meshery command line client, <a href="/pr-preview/pr-21670/installation/mesheryctl/" class="meshery-light">mesheryctl</a>.</li>
<li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
<li>Install the <a href="https://docs.digitalocean.com/reference/doctl/how-to/install/">DigitalOcean CLI (doctl)</a>, authenticated for your account.</li>
<li>A <a href="https://www.digitalocean.com/">DigitalOcean</a> account with access to either a Droplet or an active DOKS cluster.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes](/pr-preview/pr-21670/installation/kubernetes/)

## Available Deployment Methods

- [Option 1: Docker on a Droplet (Out-of-Cluster)](#option-1-docker-on-a-droplet-out-of-cluster)
  - [Provision a Droplet](#provision-a-droplet)
  - [Install Meshery on Docker](#install-meshery-on-docker)
  - [Access Meshery UI](#access-meshery-ui)
- [Option 2: DigitalOcean Kubernetes (In-Cluster)](#option-2-digitalocean-kubernetes-in-cluster)
  - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)

# Option 1: Docker on a Droplet (Out-of-Cluster)

Run Meshery as a standalone management plane on a DigitalOcean Droplet using Docker. This out-of-cluster deployment is well suited for managing one or more remote clusters from a single, always-on host.

## Provision a Droplet

Create an Ubuntu Droplet from the [DigitalOcean Control Panel](https://docs.digitalocean.com/products/droplets/how-to/create/), or with `doctl`:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">doctl compute droplet create meshery --image ubuntu-22-04-x64 --size s-2vcpu-4gb --region nyc1 --ssh-keys [YOUR_SSH_KEY_FINGERPRINT]</code>
	</div>
</pre>


Meshery runs comfortably on a Droplet with at least 2 vCPUs and 4 GB of memory. Once the Droplet is ready, connect to it over SSH:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">ssh root@[DROPLET_IP]</code>
	</div>
</pre>


Install [Docker](https://docs.docker.com/engine/install/ubuntu/) and [Docker Compose](https://docs.docker.com/compose/install/) on the Droplet, followed by [mesheryctl](/pr-preview/pr-21670/installation/mesheryctl/).

## Install Meshery on Docker

On the Droplet, start Meshery on Docker:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system start -p docker</code>
	</div>
</pre>


To manage a remote cluster (for example, a DOKS cluster) from this out-of-cluster deployment, make the cluster's kubeconfig available to Meshery. See [Customizing Kubernetes Configuration Location](/pr-preview/pr-21670/installation/docker/#customizing-kubernetes-configuration-location).

## Access Meshery UI

By default, Meshery UI is served on port `9081`. To reach it from your browser, allow inbound traffic to that port using a [DigitalOcean Cloud Firewall](https://docs.digitalocean.com/products/networking/firewalls/how-to/configure-rules/):



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">doctl compute firewall create --name meshery-ui --inbound-rules &#34;protocol:tcp,ports:9081,address:[YOUR_IP]/32&#34; --droplet-ids [DROPLET_ID]</code>
	</div>
</pre>


Open your browser and navigate to `http://[DROPLET_IP]:9081`.

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Secure your Droplet</div>


Avoid exposing Meshery UI to the public internet. Restrict the firewall rule to your own IP address, or keep port `9081` closed and reach the UI through an SSH tunnel instead: `ssh -L 9081:localhost:9081 root@[DROPLET_IP]`.
</div>


# Option 2: DigitalOcean Kubernetes (In-Cluster)

Follow the steps below to install Meshery into your DigitalOcean Kubernetes (DOKS) cluster.

## Preflight: Cluster Connectivity

1. Authenticate `doctl` with your DigitalOcean account using a [personal access token](https://docs.digitalocean.com/reference/api/create-personal-access-token/).



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">doctl auth init</code>
	</div>
</pre>

2. Download your cluster's credentials and set it as the current `kubectl` context. Replace `[CLUSTER_NAME]` with the name or ID of your DOKS cluster.



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">doctl kubernetes cluster kubeconfig save [CLUSTER_NAME]</code>
	</div>
</pre>

3. Verify your kubeconfig's current context.



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl config current-context</code>
	</div>
</pre>


## Installation: Using `mesheryctl`

Ensure that the current platform is set to `kubernetes` in `~/.meshery/config.yaml`, then execute <a href='/pr-preview/pr-21670/reference/references/mesheryctl/system/start/'>mesheryctl system start</a> to start Meshery.



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system start</code>
	</div>
</pre>


If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/pr-preview/pr-21670/guides/mesheryctl/authenticate-with-meshery-via-cli/) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/pr-preview/pr-21670/installation/kubernetes/helm/) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='/pr-preview/pr-21670/reference/references/mesheryctl/system/check/'>mesheryctl system check</a>.

To expose Meshery UI outside the cluster, create a `LoadBalancer` service; DigitalOcean automatically provisions a [DigitalOcean Load Balancer](https://docs.digitalocean.com/products/kubernetes/how-to/add-load-balancers/) and assigns an external IP. Alternatively, use port-forwarding by following the [mesheryctl system dashboard](/pr-preview/pr-21670/reference/references/mesheryctl/system/dashboard/) guide.

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
