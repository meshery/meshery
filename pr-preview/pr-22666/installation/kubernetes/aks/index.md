# AKS

> Manage your AKS clusters with Meshery. Deploy Meshery in AKS in-cluster or out-of-cluster.

---

LLMS index: [llms.txt](/pr-preview/pr-22666/llms.txt)

---

<h1>Quick Start with AKS <img src="images/aks.svg" style="width:35px;height:35px;" /></h1>

Manage your AKS clusters with Meshery. Deploy Meshery in AKS [in-cluster](#in-cluster-installation) or outside of AKS [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to install Meshery in your AKS clusters_**

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install the Meshery command line client, <a href="/pr-preview/pr-22666/installation/mesheryctl/" class="meshery-light">mesheryctl</a>.</li>
<li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
<li>Install <a href="https://learn.microsoft.com/en-us/cli/azure/install-azure-cli">Azure CLI</a>, configured for your environment.</li>
<li>Access to an active AKS cluster in one of your resource groups.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes](/pr-preview/pr-22666/installation/kubernetes/)

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
    - [Preflight: Plan your access to Meshery UI](#preflight-plan-your-access-to-meshery-ui)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)
  - [Accessing Meshery UI](#accessing-meshery-ui)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Install Meshery on Docker](#install-meshery-on-docker)

# In-cluster Installation

Follow the steps below to install Meshery in your AKS cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on AKS.

### Preflight: Cluster Connectivity

Verify your connection to an Azure Kubernetes Service cluster using Azure CLI.

1. Login to Azure account using [az login](https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli).
2. After a successful login, identify the subscription associated with your AKS cluster:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">az account set --subscription [SUBSCRIPTION_ID]</code>
	</div>
</pre>


3. After setting the subscription, set the cluster context:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">az aks get-credentials --resource-group [RESOURCE_GROUP] --name [AKS_CLUSTER_NAME]</code>
	</div>
</pre>


4. Verify your kubeconfig's current context:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl config current-context</code>
	</div>
</pre>


### Preflight: Plan your access to Meshery UI

1. If you are using port-forwarding, please refer to the [Port Forwarding Guide](/pr-preview/pr-22666/installation/kubernetes/port-forwarding/) for detailed instructions. This guide is common for all cloud providers including AWS, Azure, GCP, and Minikube cluster services.
2. If you are using a LoadBalancer, please refer to the [LoadBalancer](/pr-preview/pr-22666/installation/kubernetes/#accessing-meshery-ui) guide for detailed instructions.
3. Customize your Meshery Provider Callback URL. Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start</code>
	</div>
</pre>


## Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your AKS cluster. Configure Meshery to connect to your AKS cluster by executing:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system config aks</code>
	</div>
</pre>


Once configured, execute the following command to start Meshery.



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system start</code>
	</div>
</pre>


If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/pr-preview/pr-22666/guides/mesheryctl/authenticate-with-meshery-via-cli/) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/pr-preview/pr-22666/installation/kubernetes/helm/) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='/pr-preview/pr-22666/reference/mesheryctl/system/check/'>mesheryctl system check</a>.

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


## Accessing Meshery UI

After installing Meshery in-cluster, you can access the Meshery UI using one of several methods depending on your environment and requirements.

### Option 1: LoadBalancer Service (Recommended for Production)

When Meshery is installed in-cluster, Meshery UI is served by a Kubernetes `Service` named `meshery` in the `meshery` namespace. This `Service` is created as type `LoadBalancer` by default, forwarding port `9081` (Meshery UI) to the Meshery Server container on port `8080`.

On AKS, a `LoadBalancer` `Service` instructs Azure to provision an external load balancer and assign it a routable `EXTERNAL-IP`. Retrieve the address with:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl get service meshery --namespace meshery</code>
	</div>
</pre>


Once the `EXTERNAL-IP` column shows an address instead of `<pending>`, open Meshery UI in your browser at `http://[EXTERNAL-IP]:9081`.

If the `Service` was previously set to another type (for example, `ClusterIP`), switch it back to `LoadBalancer`:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl patch service meshery --namespace meshery --type merge -p &#39;{&#34;spec&#34;:{&#34;type&#34;:&#34;LoadBalancer&#34;}}&#39;</code>
	</div>
</pre>


<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">EXTERNAL-IP not assigned?</div>


A `LoadBalancer` `Service` is only assigned an external address when the cluster runs a load balancer controller. Azure provides one out of the box for AKS; however, if you encounter issues, ensure your AKS cluster has the cloud provider integration enabled. Alternatively, you can use port-forwarding instead.
</div>


### Option 2: Port Forwarding (Universal)

Port forwarding is a universal method to access Meshery UI that works across all Kubernetes environments, including local clusters and managed clouds. This is particularly useful when:

- A LoadBalancer is not available or practical
- You're working in a local development environment
- You need temporary access without exposing Meshery externally
- You're behind a firewall or in a restricted network

For detailed instructions on using port forwarding, see the [Port Forwarding Guide](/pr-preview/pr-22666/installation/kubernetes/port-forwarding/).

**Quick example using mesheryctl:**



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system dashboard --port-forward</code>
	</div>
</pre>


**Quick example using kubectl:**



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl port-forward svc/meshery 9081:9081 --namespace meshery</code>
	</div>
</pre>


Then open `http://localhost:9081` in your browser.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Port Forwarding vs LoadBalancer</div>



- **LoadBalancer**: Best for production deployments requiring external access on AKS
- **Port Forwarding**: Best for development, testing, and local access

Choose the method that best fits your use case. You can use both methods simultaneously if needed.
</div>


# Out-of-cluster Installation

Install Meshery on Docker (out-of-cluster) and connect it to your AKS cluster.

## Install Meshery on Docker



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system start -p docker</code>
	</div>
</pre>


Configure Meshery to connect to your cluster by executing:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system config aks</code>
	</div>
</pre>


Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081).

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
