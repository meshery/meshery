# Install on Kubernetes

> Install Meshery on Kubernetes. Deploy Meshery in Kubernetes in-cluster or outside of Kubernetes out-of-cluster.

---

LLMS index: [llms.txt](/pr-preview/pr-22666/llms.txt)

---

<h1>Quick Start with Kubernetes <img src="images/kubernetes.svg" style="width:35px;height:35px;" /></h1>

Manage your Kubernetes clusters with Meshery. Deploy Meshery in Kubernetes [in-cluster](#in-cluster-installation) or outside of Kubernetes [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to install Meshery in your Kubernetes clusters_**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="/pr-preview/pr-22666/installation/mesheryctl/" class="meshery-light">mesheryctl</a>.</li>
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
  - [Accessing Meshery UI](#accessing-meshery-ui)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Set up Ingress on Minikube with the NGINX Ingress Controller](#set-up-ingress-on-minikube-with-the-nginx-ingress-controller)
  - [Installing cert-manager with kubectl](#installing-cert-manager-with-kubectl)

# In-cluster Installation

Follow the steps below to install Meshery in your Kubernetes cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on Kubernetes.

### Preflight: Cluster Connectivity

Verify your kubeconfig's current context is set to the Kubernetes cluster you want to deploy Meshery to.



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl config current-context</code>
	</div>
</pre>


## Installation: Using `mesheryctl`

Once configured, execute the following command to start Meshery.

Before executing the below command, go to ~/.meshery/config.yaml and ensure that the current platform is set to Kubernetes.



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system start</code>
	</div>
</pre>


## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/pr-preview/pr-22666/installation/kubernetes/helm/) guide. This guide is common for all cloud providers including AWS, Azure, GCP, and Minikube cluster services.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='/pr-preview/pr-22666/reference/mesheryctl/system/check/'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

## Accessing Meshery UI

After installing Meshery in-cluster, you can access the Meshery UI using one of several methods depending on your environment and requirements.

### Option 1: LoadBalancer Service (Recommended for Production)

When Meshery is installed in-cluster, Meshery UI is served by a Kubernetes `Service` named `meshery` in the `meshery` namespace. This `Service` is created as type `LoadBalancer` by default, forwarding port `9081` (Meshery UI) to the Meshery Server container on port `8080`.

On a managed Kubernetes offering - such as GKE, EKS, AKS, or DigitalOcean Kubernetes - a `LoadBalancer` `Service` instructs the cloud provider to provision an external load balancer and assign it a routable `EXTERNAL-IP`. Retrieve the address with:



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


A `LoadBalancer` `Service` is only assigned an external address when the cluster runs a load balancer controller. Managed clouds provide one out of the box; bare-metal and local clusters (for example, Minikube or KinD) do not. On those clusters, install a load balancer implementation such as [MetalLB](https://metallb.universe.tf/), or use port-forwarding instead.
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



- **LoadBalancer**: Best for production deployments requiring external access
- **Port Forwarding**: Best for development, testing, and local access

Choose the method that best fits your use case. You can use both methods simultaneously if needed.
</div>


### Option 1: LoadBalancer Service (Recommended for Production)

When Meshery is installed in-cluster, Meshery UI is served by a Kubernetes `Service` named `meshery` in the `meshery` namespace. This `Service` is created as type `LoadBalancer` by default, forwarding port `9081` (Meshery UI) to the Meshery Server container on port `8080`.

On a managed Kubernetes offering - such as GKE, EKS, AKS, or DigitalOcean Kubernetes - a `LoadBalancer` `Service` instructs the cloud provider to provision an external load balancer and assign it a routable `EXTERNAL-IP`. Retrieve the address with:



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


A `LoadBalancer` `Service` is only assigned an external address when the cluster runs a load balancer controller. Managed clouds provide one out of the box; bare-metal and local clusters (for example, Minikube or KinD) do not. On those clusters, install a load balancer implementation such as [MetalLB](https://metallb.universe.tf/), or use port-forwarding instead.
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



- **LoadBalancer**: Best for production deployments requiring external access
- **Port Forwarding**: Best for development, testing, and local access

Choose the method that best fits your use case. You can use both methods simultaneously if needed.
</div>


# Out-of-cluster Installation

Install Meshery on Docker (out-of-cluster) and connect it to your Kubernetes cluster.

<!-- ## Installation: Upload Config File in Meshery Web UI

- Run the below command to generate the _"config_minikube.yaml"_ file for your cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">kubectl config view --minify --flatten > config_minikube.yaml</div></div>
 </pre>

- Upload the generated config file by navigating to _Settings > Environment > Out of Cluster Deployment_ in the Web UI and using the _"Upload kubeconfig"_ option. -->

## Set up Ingress on Minikube with the NGINX Ingress Controller

- Run the below command to enable the NGINX Ingress controller for your cluster:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">minikube addons enable ingress</code>
	</div>
</pre>


- To check if NGINX Ingress controller is running



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl get pods -n ingress-nginx</code>
	</div>
</pre>


## Installing cert-manager with kubectl

- Run the below command to install cert-manager for your cluster:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.3/cert-manager.yaml</code>
	</div>
</pre>


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

---

Section pages:

- [AKS](/pr-preview/pr-22666/installation/kubernetes/aks/): Manage your AKS clusters with Meshery. Deploy Meshery in AKS in-cluster or out-of-cluster.
- [DigitalOcean](/pr-preview/pr-22666/installation/kubernetes/digitalocean/): Install Meshery on DigitalOcean. Deploy Meshery out-of-cluster with Docker on a Droplet, or in-cluster on DigitalOcean Kubernetes (DOKS).
- [EKS](/pr-preview/pr-22666/installation/kubernetes/eks/): Install Meshery on Elastic Kubernetes Service. Deploy Meshery in EKS in-cluster or outside of EKS out-of-cluster.
- [GKE](/pr-preview/pr-22666/installation/kubernetes/gke/): Install Meshery on Google Kubernetes Engine. Deploy Meshery in GKE in-cluster or outside of GKE out-of-cluster.
- [Helm](/pr-preview/pr-22666/installation/kubernetes/helm/): Install Meshery on Kubernetes using Helm. Deploy Meshery in Kubernetes in-cluster. This guide applies to all cloud providers including AWS (EKS), Azure (AKS), GCP (GKE), and local clusters (Minikube, KinD).
- [KinD](/pr-preview/pr-22666/installation/kubernetes/kind/): Install Meshery on KinD. Deploy Meshery in KinD in-cluster or outside of KinD out-of-cluster.
- [KubeSphere](/pr-preview/pr-22666/installation/kubernetes/kubesphere/): Install Meshery on KubeSphere
- [Minikube](/pr-preview/pr-22666/installation/kubernetes/minikube/): Install Meshery on Minikube. Deploy Meshery in Minikube in-cluster or outside of Minikube out-of-cluster.
- [Port Forwarding](/pr-preview/pr-22666/installation/kubernetes/port-forwarding/): Access Meshery UI using port forwarding when LoadBalancer is not available.
