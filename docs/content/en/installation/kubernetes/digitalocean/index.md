---
title: DigitalOcean
categories: [kubernetes]
aliases:
- /installation/platforms/digitalocean
- /installation/kubernetes/doks
display_title: false
image: installation/kubernetes/digitalocean/images/digitalocean.svg
description: Install Meshery on DigitalOcean. Deploy Meshery out-of-cluster with Docker on a Droplet, or in-cluster on DigitalOcean Kubernetes (DOKS).
---

<h1>Quick Start with DigitalOcean <img src="images/digitalocean.svg" style="width:35px;height:35px;" /></h1>

Deploy and manage your DigitalOcean infrastructure with Meshery. You can run Meshery on DigitalOcean in two ways: [out-of-cluster](#option-1-docker-on-a-droplet-out-of-cluster) using Docker on a Droplet, or [in-cluster](#option-2-digitalocean-kubernetes-in-cluster) on a DigitalOcean Kubernetes (DOKS) cluster. **_Note: It is advisable to install Meshery in your DOKS cluster._**

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install the Meshery command line client, <a href="{{< ref "installation/mesheryctl/_index.md" >}}" class="meshery-light">mesheryctl</a>.</li>
<li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
<li>Install the <a href="https://docs.digitalocean.com/reference/doctl/how-to/install/">DigitalOcean CLI (doctl)</a>, authenticated for your account.</li>
<li>A <a href="https://www.digitalocean.com/">DigitalOcean</a> account with access to either a Droplet or an active DOKS cluster.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes]({{< ref "installation/kubernetes/_index.md" >}})

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

{{< code code="doctl compute droplet create meshery --image ubuntu-22-04-x64 --size s-2vcpu-4gb --region nyc1 --ssh-keys [YOUR_SSH_KEY_FINGERPRINT]" >}}

Meshery runs comfortably on a Droplet with at least 2 vCPUs and 4 GB of memory. Once the Droplet is ready, connect to it over SSH:

{{< code code="ssh root@[DROPLET_IP]" >}}

Install [Docker](https://docs.docker.com/engine/install/ubuntu/) and [Docker Compose](https://docs.docker.com/compose/install/) on the Droplet, followed by [mesheryctl]({{< ref "installation/mesheryctl/_index.md" >}}).

## Install Meshery on Docker

On the Droplet, start Meshery on Docker:

{{< code code="mesheryctl system start -p docker" >}}

To manage a remote cluster (for example, a DOKS cluster) from this out-of-cluster deployment, make the cluster's kubeconfig available to Meshery. See [Customizing Kubernetes Configuration Location]({{< ref "installation/docker/_index.md#customizing-kubernetes-configuration-location" >}}).

## Access Meshery UI

By default, Meshery UI is served on port `9081`. To reach it from your browser, allow inbound traffic to that port using a [DigitalOcean Cloud Firewall](https://docs.digitalocean.com/products/networking/firewalls/how-to/configure-rules/):

{{< code code=`doctl compute firewall create --name meshery-ui --inbound-rules "protocol:tcp,ports:9081,address:[YOUR_IP]/32" --droplet-ids [DROPLET_ID]` >}}

Open your browser and navigate to `http://[DROPLET_IP]:9081`.

{{% alert title="Secure your Droplet" color="warning" %}}
Avoid exposing Meshery UI to the public internet. Restrict the firewall rule to your own IP address, or keep port `9081` closed and reach the UI through an SSH tunnel instead: `ssh -L 9081:localhost:9081 root@[DROPLET_IP]`.
{{% /alert %}}

# Option 2: DigitalOcean Kubernetes (In-Cluster)

Follow the steps below to install Meshery into your DigitalOcean Kubernetes (DOKS) cluster.

## Preflight: Cluster Connectivity

1. Authenticate `doctl` with your DigitalOcean account using a [personal access token](https://docs.digitalocean.com/reference/api/create-personal-access-token/).

{{< code code="doctl auth init" >}}
2. Download your cluster's credentials and set it as the current `kubectl` context. Replace `[CLUSTER_NAME]` with the name or ID of your DOKS cluster.

{{< code code="doctl kubernetes cluster kubeconfig save [CLUSTER_NAME]" >}}
3. Verify your kubeconfig's current context.

{{< code code="kubectl config current-context" >}}

## Installation: Using `mesheryctl`

Ensure that the current platform is set to `kubernetes` in `~/.meshery/config.yaml`, then execute <a href='{{< ref "reference/references/mesheryctl/system/start.md" >}}'>mesheryctl system start</a> to start Meshery.

{{< code code="mesheryctl system start" >}}

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}}) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) guide. This guide is common for all cloud providers including AWS, Azure, GCP, and Minikube cluster services.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='{{< ref "reference/references/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

To expose Meshery UI outside the cluster, create a `LoadBalancer` service; DigitalOcean automatically provisions a [DigitalOcean Load Balancer](https://docs.digitalocean.com/products/kubernetes/how-to/add-load-balancers/) and assigns an external IP. Alternatively, use port-forwarding.

For more information on accessing Meshery UI, including detailed port-forwarding instructions, see the [Accessing Meshery UI]({{< ref "installation/kubernetes/_index.md#accessing-meshery-ui" >}}) section and the [Port Forwarding Guide]({{< ref "installation/kubernetes/port-forwarding.md" >}}).

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

## Accessing Meshery UI

After installing Meshery in-cluster, you can access the Meshery UI using one of several methods depending on your environment and requirements.

### Option 1: LoadBalancer Service (Recommended for Production)

When Meshery is installed in-cluster, Meshery UI is served by a Kubernetes `Service` named `meshery` in the `meshery` namespace. This `Service` is created as type `LoadBalancer` by default, forwarding port `9081` (Meshery UI) to the Meshery Server container on port `8080`.

On DOKS, a `LoadBalancer` `Service` instructs DigitalOcean to provision an external load balancer and assign it a routable `EXTERNAL-IP`. Retrieve the address with:

{{< code code="kubectl get service meshery --namespace meshery" >}}

Once the `EXTERNAL-IP` column shows an address instead of `<pending>`, open Meshery UI in your browser at `http://[EXTERNAL-IP]:9081`.

If the `Service` was previously set to another type (for example, `ClusterIP`), switch it back to `LoadBalancer`:

{{< code code=`kubectl patch service meshery --namespace meshery --type merge -p '{"spec":{"type":"LoadBalancer"}}'` >}}

{{% alert title="EXTERNAL-IP not assigned?" color="warning" %}}
A `LoadBalancer` `Service` is only assigned an external address when the cluster runs a load balancer controller. DigitalOcean provides one out of the box for DOKS; however, if you encounter issues, ensure your DOKS cluster has the load balancer integration enabled. Alternatively, you can use port-forwarding instead.
{{% /alert %}}

### Option 2: Port Forwarding (Universal)

Port forwarding is a universal method to access Meshery UI that works across all Kubernetes environments, including local clusters and managed clouds. This is particularly useful when:

- A LoadBalancer is not available or practical
- You're working in a local development environment
- You need temporary access without exposing Meshery externally
- You're behind a firewall or in a restricted network

For detailed instructions on using port forwarding, see the [Port Forwarding Guide]({{< ref "installation/kubernetes/port-forwarding.md" >}}).

**Quick example using mesheryctl:**

{{< code code="mesheryctl system dashboard --port-forward" >}}

**Quick example using kubectl:**

{{< code code="kubectl port-forward svc/meshery 9081:9081 --namespace meshery" >}}

Then open `http://localhost:9081` in your browser.

{{% alert title="Port Forwarding vs LoadBalancer" color="info" %}}
- **LoadBalancer**: Best for production deployments requiring external access on DOKS
- **Port Forwarding**: Best for development, testing, and local access

Choose the method that best fits your use case. You can use both methods simultaneously if needed.
{{% /alert %}}

{{< installation/accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}
