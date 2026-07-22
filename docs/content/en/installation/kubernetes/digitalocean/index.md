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
<li>Install the <a href="https://docs.digitalocean.com/reference/doctl/how-to/install/">DigitalOcean CLI (doctl)</a> on your local machine.</li>
<li>A <a href="https://www.digitalocean.com/">DigitalOcean</a> account with access to either a Droplet or an active DOKS cluster.</li>
<li><strong>Option 2 only:</strong> Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> and the Meshery CLI, <a href="{{< ref "installation/mesheryctl/_index.md" >}}" class="meshery-light">mesheryctl</a>, on your local machine. If you use Helm, also install <a href="https://helm.sh/docs/intro/install/">Helm</a> v3.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes]({{< ref "installation/kubernetes/_index.md" >}})

## Available Deployment Methods

- [Option 1: Docker on a Droplet (Out-of-Cluster)](#option-1-docker-on-a-droplet-out-of-cluster)
  - [Authenticate `doctl`](#authenticate-doctl)
  - [Register an SSH key](#register-an-ssh-key)
  - [Provision a Droplet](#provision-a-droplet)
  - [Install Docker and mesheryctl on the Droplet](#install-docker-and-mesheryctl-on-the-droplet)
  - [Install Meshery on Docker](#install-meshery-on-docker)
  - [Access Meshery UI](#access-meshery-ui)
  - [Connect a remote Kubernetes cluster](#connect-a-remote-kubernetes-cluster)
- [Option 2: DigitalOcean Kubernetes (In-Cluster)](#option-2-digitalocean-kubernetes-in-cluster)
  - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)

# Option 1: Docker on a Droplet (Out-of-Cluster)

Run Meshery as a standalone management plane on a DigitalOcean Droplet using Docker. This out-of-cluster deployment is well suited for managing one or more remote clusters from a single, always-on host.

## Authenticate `doctl`

Before creating resources, authenticate `doctl` with a [personal access token](https://docs.digitalocean.com/reference/api/create-personal-access-token/) from your DigitalOcean account:

{{< code code="doctl auth init" >}}

## Register an SSH key

Droplet creation requires an SSH key that is already registered with DigitalOcean. List keys already on your account:

{{< code code="doctl compute ssh-key list" >}}

If the list is empty, generate a key locally (if you do not already have one), then import the **public** key into DigitalOcean:

{{< code code=`ssh-keygen -t ed25519 -C "meshery-droplet" -f ~/.ssh/id_ed25519
doctl compute ssh-key import meshery-key --public-key-file ~/.ssh/id_ed25519.pub
doctl compute ssh-key list` >}}

Copy the fingerprint (or key ID) from the list. You will pass it to `--ssh-keys` in the next step.

## Provision a Droplet

Create an Ubuntu Droplet from the [DigitalOcean Control Panel](https://docs.digitalocean.com/products/droplets/how-to/create/), or with `doctl`:

{{< code code="doctl compute droplet create meshery --image ubuntu-22-04-x64 --size s-2vcpu-4gb --region nyc1 --ssh-keys [YOUR_SSH_KEY_FINGERPRINT]" >}}

Meshery runs comfortably on a Droplet with at least 2 vCPUs and 4 GB of memory.

`doctl compute droplet create` returns immediately while the Droplet is still provisioning (`status: new`) and often **before** a public IPv4 address is assigned. Note the Droplet `ID` from the create output, then poll until the Droplet is `active` and has a public IP:

{{< code code="doctl compute droplet get [DROPLET_ID] --format ID,Name,PublicIPv4,Status" >}}

Repeat the get command until `Status` is `active` and `PublicIPv4` is populated. Then connect over SSH:

{{< code code="ssh root@[DROPLET_IP]" >}}

## Install Docker and mesheryctl on the Droplet

On the Droplet (Ubuntu), install Docker Engine and the Compose plugin using Docker's distribution-specific guides (keep these links so you can follow the correct steps for your image and version):

- [Install Docker Engine on Ubuntu](https://docs.docker.com/engine/install/ubuntu/)
- [Install Docker Compose](https://docs.docker.com/compose/install/)

This guide uses an **Ubuntu** Droplet, so install `mesheryctl` with the **Bash** installer (not Homebrew or Scoop). Install the CLI only, then start Meshery in the next section:

{{< code code="curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -" >}}

For other operating systems or package managers, see [Install mesheryctl]({{< ref "installation/mesheryctl/_index.md" >}}).

{{% alert title="Install CLI only on the Droplet" color="info" %}}
Using `DEPLOY_MESHERY=false` installs `mesheryctl` without auto-deploying Meshery. That keeps this guide's next step (`mesheryctl system start -p docker`) as the single, intentional start command. If you instead run `curl -L https://meshery.io/install | PLATFORM=docker bash -`, the script installs `mesheryctl` **and** starts Meshery for you—you can skip the start command below.
{{% /alert %}}

## Install Meshery on Docker

On the Droplet, start Meshery on Docker:

{{< code code="mesheryctl system start -p docker" >}}

## Access Meshery UI

By default, Meshery UI is served on port `9081` on the Droplet. You are typically SSHed into a headless host, so open the UI from **your local machine**, not from a browser on the Droplet.

**Recommended:** keep port `9081` closed to the public internet and use an SSH tunnel from your local machine:

{{< code code="ssh -L 9081:localhost:9081 root@[DROPLET_IP]" >}}

Then open [http://localhost:9081](http://localhost:9081) in your local browser.

**Optional:** allow inbound traffic with a [DigitalOcean Cloud Firewall](https://docs.digitalocean.com/products/networking/firewalls/how-to/configure-rules/) restricted to your IP, then open `http://[DROPLET_IP]:9081`:

{{< code code=`doctl compute firewall create --name meshery-ui --inbound-rules "protocol:tcp,ports:9081,address:[YOUR_IP]/32" --droplet-ids [DROPLET_ID]` >}}

{{% alert title="Secure your Droplet" color="warning" %}}
Avoid exposing Meshery UI to the public internet. Prefer the SSH tunnel, or restrict any firewall rule to your own IP address only.
{{% /alert %}}

## Connect a remote Kubernetes cluster

After a fresh out-of-cluster install, Meshery UI will not show any connected Kubernetes clusters. That is expected: the Droplet does not ship with a kubeconfig for DOKS or any other cluster.

To manage a remote cluster (for example, a DOKS cluster) from this deployment:

1. Obtain the cluster kubeconfig on a machine that can reach the cluster (for DOKS: `doctl kubernetes cluster kubeconfig save [CLUSTER_NAME]`).
2. Make that kubeconfig available to Meshery on the Droplet. See [Customizing Kubernetes Configuration Location]({{< ref "installation/docker/_index.md#customizing-kubernetes-configuration-location" >}}).
3. In the Meshery UI, open **Lifecycle → Connections** and confirm the cluster appears and is marked `Connected`.

# Option 2: DigitalOcean Kubernetes (In-Cluster)

Follow the steps below to install Meshery into your DigitalOcean Kubernetes (DOKS) cluster.

## Preflight: Cluster Connectivity

1. Authenticate `doctl` with your DigitalOcean account using a [personal access token](https://docs.digitalocean.com/reference/api/create-personal-access-token/). Your token must include permission to **manage Kubernetes clusters** (read and write). A token that can only list account resources may pass `doctl auth init` but fail later when downloading kubeconfig (`403`).

{{< code code="doctl auth init" >}}

2. Download your cluster's credentials and set it as the current `kubectl` context. Replace `[CLUSTER_NAME]` with the name or ID of your DOKS cluster.

{{< code code="doctl kubernetes cluster kubeconfig save [CLUSTER_NAME]" >}}

3. Verify your kubeconfig's current context and that the API server is reachable.

{{< code code=`kubectl config current-context
kubectl get nodes` >}}

## Installation: Using `mesheryctl`

Set the Meshery context platform to Kubernetes, then start Meshery. Prefer the CLI over editing `~/.meshery/config.yaml` by hand:

{{< code code=`mesheryctl system context create doks --platform kubernetes --set
mesheryctl system start` >}}

You can also pass the platform on start:

{{< code code="mesheryctl system start -p kubernetes" >}}

If you encounter Meshery **provider** authentication issues (sign-in to a Remote Provider), use `mesheryctl system login`. For more information, see [Authenticate with Meshery via CLI]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}}). This is separate from Kubernetes API authentication configured by your kubeconfig.

## Installation: Using Helm

Install [Helm](https://helm.sh/docs/intro/install/) v3 if it is not already available, then deploy Meshery. Enable the Meshery Operator so MeshSync and Broker are installed; without the operator, the UI may load while cluster lifecycle features remain unavailable.

{{< code code=`helm repo add meshery https://meshery.io/charts/
helm repo update
helm install meshery meshery/meshery --namespace meshery --create-namespace --set meshery-operator.enabled=true` >}}

For additional configuration options, see the [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='{{< ref "reference/references/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

Meshery is deployed with a Kubernetes `Service` of type `LoadBalancer` by default. On DOKS, DigitalOcean provisions a [Load Balancer](https://docs.digitalocean.com/products/kubernetes/how-to/add-load-balancers/) and assigns an external IP. Watch for the address:

{{< code code="kubectl get svc -n meshery meshery -w" >}}

Alternatively, use port-forwarding via [mesheryctl system dashboard]({{< ref "reference/references/mesheryctl/system/dashboard.md" >}}) or:

{{< code code="kubectl port-forward svc/meshery -n meshery 9081:9081" >}}

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< installation/accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}
