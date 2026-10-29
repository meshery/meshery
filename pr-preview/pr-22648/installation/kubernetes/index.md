# Install on Kubernetes

> Install Meshery on Kubernetes. Deploy Meshery in Kubernetes in-cluster or outside of Kubernetes out-of-cluster.

---

LLMS index: [llms.txt](/pr-preview/pr-22648/llms.txt)

---

<h1>Quick Start with Kubernetes <img src="images/kubernetes.svg" style="width:35px;height:35px;" /></h1>

Manage your Kubernetes clusters with Meshery. Deploy Meshery inside a Kubernetes cluster ([in-cluster](#in-cluster-installation)) or outside of it ([out-of-cluster](#out-of-cluster-installation)). **_Note: We recommend installing Meshery in your Kubernetes clusters._**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="/pr-preview/pr-22648/installation/mesheryctl/" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Access to an active Kubernetes cluster.</li>
  </ol>
</div>

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)
  - [Exposing Meshery Service (LoadBalancer)](#exposing-meshery-serviceloadbalancer)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Set up Ingress on Minikube with the NGINX Ingress Controller](#set-up-ingress-on-minikube-with-the-nginx-ingress-controller)
  - [Installing cert-manager with kubectl](#installing-cert-manager-with-kubectl)

# In-cluster Installation

Follow the steps below to install Meshery in your Kubernetes cluster.

## Preflight Checks

Before you deploy Meshery on Kubernetes, check the following.

### Preflight: Cluster Connectivity

Verify that your kubeconfig's current context points at the cluster you want to deploy Meshery to.



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl config current-context</code>
	</div>
</pre>


## Installation: Using `mesheryctl`

Start Meshery. `mesheryctl` deploys to the platform of your current `mesheryctl` context, which is `kubernetes` by default. To check it, run `mesheryctl system context view`. To deploy to Kubernetes regardless of the context, pass `-p kubernetes`:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system start -p kubernetes</code>
	</div>
</pre>


`mesheryctl` installs the Meshery Helm chart into the `meshery` namespace of your current cluster.

## Installation: Using Helm

To install Meshery with Helm 3, see the [Helm Installation](/pr-preview/pr-22648/installation/kubernetes/helm/) guide.

## Post-Installation Steps

Optionally, verify the health of your Meshery deployment with <a href='/pr-preview/pr-22648/reference/mesheryctl/system/check/'>mesheryctl system check</a>.

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
  You can customize the callback URL that your <a href="/pr-preview/pr-22648/reference/extensibility/providers/">Meshery Provider</a> redirects to after sign-in. This helps when Meshery runs behind several layers of networking infrastructure.
</p>

<p>
  For production deployments, we recommend that you expose Meshery UI through a reverse proxy or a load balancer. A custom callback URL makes sure that sign-in completes even when requests pass through several routing layers.
</p>

<p>
  <b>Note</b>: For production deployments, use a Remote Provider so that you control which identity providers can sign users in. Learn more about this in the <a href="/pr-preview/pr-22648/reference/extensibility/providers/">Extensibility: Providers</a> guide.
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


## Exposing Meshery Service (LoadBalancer) {#exposing-meshery-serviceloadbalancer}

When you install Meshery in-cluster, a Kubernetes `Service` named `meshery` in the `meshery` namespace serves Meshery UI. The Helm chart creates this `Service` as type `LoadBalancer` by default (`service.type`), forwarding port `9081` to the Meshery Server container on port `8080`.

On a managed Kubernetes offering, such as GKE, EKS, AKS, or DigitalOcean Kubernetes, a `LoadBalancer` `Service` tells the cloud provider to provision an external load balancer and assign it a routable `EXTERNAL-IP`. Retrieve the address with:



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


A `LoadBalancer` `Service` gets an external address only when the cluster runs a load balancer controller. Managed clouds provide one out of the box, but bare-metal and local clusters (for example, Minikube or kind) do not. On those clusters, install a load balancer implementation such as [MetalLB](https://metallb.universe.tf/), expose Meshery through a `NodePort` `Service` instead, or reach the UI with port-forwarding by following the [mesheryctl system dashboard](/pr-preview/pr-22648/reference/mesheryctl/system/dashboard/) guide.
</div>


# Out-of-cluster Installation

Install Meshery on Docker (out-of-cluster) and connect it to your Kubernetes cluster. Follow [Install on Docker](/pr-preview/pr-22648/installation/docker/), then sign in with `mesheryctl system login` and add your cluster with `mesheryctl connection create --type <aks|eks|gke|minikube>`, or upload its kubeconfig in Meshery UI under **Lifecycle** > **Connections**.

If you expose Meshery through an ingress on Minikube, the following steps set up an ingress controller and cert-manager.

<!-- ## Installation: Upload Config File in Meshery Web UI

- Run the below command to generate the _"config_minikube.yaml"_ file for your cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">kubectl config view --minify --flatten > config_minikube.yaml</div></div>
 </pre>

- Upload the generated config file by navigating to _Settings > Environment > Out of Cluster Deployment_ in the Web UI and using the _"Upload kubeconfig"_ option. -->

## Set up Ingress on Minikube with the NGINX Ingress Controller

- Enable the NGINX Ingress controller for your cluster:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">minikube addons enable ingress</code>
	</div>
</pre>


- Check that the NGINX Ingress controller is running:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl get pods -n ingress-nginx</code>
	</div>
</pre>


## Installing cert-manager with kubectl

- Install cert-manager in your cluster:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.3/cert-manager.yaml</code>
	</div>
</pre>


<div class="related-discussions">
  <h3>Recent Discussions with "meshery" Tag</h3><ul><li>
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
        </li><li>
          <a href="https://discuss.meshery.io/t/play-meshery-io-is-down/8045" target="_blank" rel="noopener noreferrer">
            play.meshery.io is DOWN
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>

---

Section pages:

- [AKS](/pr-preview/pr-22648/installation/kubernetes/aks/): Manage your AKS clusters with Meshery. Deploy Meshery in AKS in-cluster or out-of-cluster.
- [DigitalOcean](/pr-preview/pr-22648/installation/kubernetes/digitalocean/): Install Meshery on DigitalOcean. Deploy Meshery out-of-cluster with Docker on a Droplet, or in-cluster on DigitalOcean Kubernetes (DOKS).
- [EKS](/pr-preview/pr-22648/installation/kubernetes/eks/): Install Meshery on Elastic Kubernetes Service. Deploy Meshery in EKS in-cluster or outside of EKS out-of-cluster.
- [GKE](/pr-preview/pr-22648/installation/kubernetes/gke/): Install Meshery on Google Kubernetes Engine. Deploy Meshery in GKE in-cluster or outside of GKE out-of-cluster.
- [Helm](/pr-preview/pr-22648/installation/kubernetes/helm/): Install Meshery on Kubernetes using Helm. Deploy Meshery in Kubernetes in-cluster.
- [KinD](/pr-preview/pr-22648/installation/kubernetes/kind/): Install Meshery on KinD. Deploy Meshery in KinD in-cluster or outside of KinD out-of-cluster.
- [KubeSphere](/pr-preview/pr-22648/installation/kubernetes/kubesphere/): Install Meshery on KubeSphere
- [Minikube](/pr-preview/pr-22648/installation/kubernetes/minikube/): Install Meshery on Minikube. Deploy Meshery in Minikube in-cluster or outside of Minikube out-of-cluster.
