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
<li>Install the <a href="https://docs.digitalocean.com/reference/doctl/how-to/install/">DigitalOcean CLI (doctl)</a> and authenticate it for your account with <code>doctl auth init</code>.</li>
<li>A <a href="https://www.digitalocean.com/">DigitalOcean</a> account with access to either a Droplet or an active DOKS cluster.</li>
<li><strong>Option 1 only:</strong> An SSH key registered with your DigitalOcean account. See DigitalOcean's guide to <a href="https://docs.digitalocean.com/products/droplets/how-to/add-ssh-keys/">adding SSH keys</a>.</li>
<li><strong>Option 2 only:</strong> Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> and the Meshery CLI, <a href="{{< ref "installation/mesheryctl/_index.md" >}}" class="meshery-light">mesheryctl</a>, on your local machine. If you use Helm, also install <a href="https://helm.sh/docs/intro/install/">Helm</a> v3.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes]({{< ref "installation/kubernetes/_index.md" >}})

## Available Deployment Methods

- [Option 1: Docker on a Droplet (Out-of-Cluster)](#option-1-docker-on-a-droplet-out-of-cluster)
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

Run Meshery on a DigitalOcean Droplet with Docker so you can manage remote clusters from one host.

## Provision a Droplet

On your local machine, authenticate `doctl` with a [personal access token](https://docs.digitalocean.com/reference/api/create-personal-access-token/):

{{< code code="doctl auth init" >}}

Create an Ubuntu Droplet (2 vCPU / 4 GB or larger) from the [Control Panel](https://docs.digitalocean.com/products/droplets/how-to/create/) or with `doctl`. Replace `[SSH_KEY_FINGERPRINT]` with your key fingerprint:

{{< code code="doctl compute droplet create meshery --image ubuntu-22-04-x64 --size s-2vcpu-4gb --region nyc1 --ssh-keys [SSH_KEY_FINGERPRINT]" >}}

From the create output, copy **ID**. Public IPv4 may still be empty while **Status** is `new`. Poll until **Status** is `active` and **Public IPv4** is set:

{{< code code="doctl compute droplet get [ID] --format ID,Name,PublicIPv4,Status" >}}

**Example output:**

| ID | Name | Public IPv4 | Status |
| --- | --- | --- | --- |
| `XXXXXXXXX` | meshery | `XXX.XXX.XXX.XXX` | active |

Use **ID** for later `doctl` commands (get, firewall). Use **Public IPv4** for SSH, browser, and tunnel:

{{< code code="ssh root@[Public IPv4]" >}}

## Install Docker and mesheryctl on the Droplet

On the Droplet, install [Docker Engine](https://docs.docker.com/engine/install/ubuntu/) and [Docker Compose](https://docs.docker.com/compose/install/), then install `mesheryctl` only (Ubuntu Bash path):

{{< code code="curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -" >}}

`DEPLOY_MESHERY=false` installs the CLI without starting Meshery. To install and start in one step instead, use `curl -L https://meshery.io/install | PLATFORM=docker bash -` and skip the next section. Other install methods: [Install mesheryctl]({{< ref "installation/mesheryctl/_index.md" >}}).

## Install Meshery on Docker

{{< code code="mesheryctl system start -p docker" >}}

## Access Meshery UI

After Meshery is running, open `http://[Public IPv4]:9081` in your browser.

{{% alert color="warning" title="Secure your Droplet" %}}
By default, a DigitalOcean Droplet does not restrict inbound traffic at the network edge. Once Meshery is listening, port `9081` may be reachable from the public internet. Use a [DigitalOcean Cloud Firewall](https://docs.digitalocean.com/products/networking/firewalls/how-to/configure-rules/) to control access. Create the firewall from your local machine with `doctl`:

<pre class="codeblock-pre" style="margin-bottom: 1rem;"><div class="codeblock"><code class="clipboardjs" style="white-space: pre-wrap; overflow-wrap: anywhere;">doctl compute firewall create --name meshery-ui --inbound-rules "protocol:tcp,ports:22,address:0.0.0.0/0 protocol:tcp,ports:9081,address:[YOUR_IP]/32" --outbound-rules "protocol:tcp,ports:0,address:0.0.0.0/0 protocol:udp,ports:0,address:0.0.0.0/0 protocol:icmp,address:0.0.0.0/0" --droplet-ids [ID]</code></div></pre>

<ul>
<li><code>--inbound-rules</code>: allow SSH on port <code>22</code>; limit Meshery UI (<code>9081</code>) to <code>[YOUR_IP]</code>.</li>
<li><code>--outbound-rules</code>: permit DNS and HTTPS egress so Meshery can communicate with remote Kubernetes APIs.</li>
<li><code>--droplet-ids [ID]</code>: attach the firewall to the Droplet <strong>ID</strong> from provisioning.</li>
</ul>

Alternatively, leave port `9081` closed on the public interface and access the UI over an SSH tunnel: `ssh -L 9081:localhost:9081 root@[Public IPv4]`, then open [http://localhost:9081](http://localhost:9081).
{{% /alert %}}

## Connect a remote Kubernetes cluster

An empty cluster list in the UI is expected until you attach a kubeconfig.

1. For DOKS: `doctl kubernetes cluster kubeconfig save [CLUSTER_NAME]`.
2. Make that kubeconfig available to Meshery on the Droplet ([Customizing Kubernetes Configuration Location]({{< ref "installation/docker/_index.md#customizing-kubernetes-configuration-location" >}})). If the file uses `exec`/`doctl` from your laptop, use a portable token-based kubeconfig on the Droplet instead.
3. In the UI: **Lifecycle → Connections** → confirm the cluster is `Connected`.

# Option 2: DigitalOcean Kubernetes (In-Cluster)

Install Meshery into a DigitalOcean Kubernetes (DOKS) cluster.

## Preflight: Cluster Connectivity

**1.** Authenticate `doctl` with a [personal access token](https://docs.digitalocean.com/reference/api/create-personal-access-token/) that can manage Kubernetes clusters:

{{< code code="doctl auth init" >}}

**2.** Download cluster credentials (replace `[CLUSTER_NAME]`):

{{< code code="doctl kubernetes cluster kubeconfig save [CLUSTER_NAME]" >}}

**3.** Confirm connectivity:

{{< code code=`kubectl config current-context
kubectl get nodes` >}}

## Installation: Using `mesheryctl`

Set the platform to Kubernetes, then start Meshery:

{{< code code=`mesheryctl system context create doks --platform kubernetes --set
mesheryctl system start` >}}

Or in one step: `mesheryctl system start -p kubernetes`.

If provider login fails, see [Authenticate with Meshery via CLI]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}}). That login is separate from your kubeconfig.

## Installation: Using Helm

Deploy with the Operator enabled (required for MeshSync and Broker):

{{< code code=`helm repo add meshery https://meshery.io/charts/
helm repo update
helm install meshery meshery/meshery --namespace meshery --create-namespace --set meshery-operator.enabled=true` >}}

More options: [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}).

## Post-Installation Steps

Optional health check: <a href='{{< ref "reference/references/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

Meshery uses a `LoadBalancer` Service by default. On DOKS, wait for the external IP:

{{< code code="kubectl get svc -n meshery meshery -w" >}}

Or port-forward with [mesheryctl system dashboard]({{< ref "reference/references/mesheryctl/system/dashboard.md" >}}) or:

{{< code code="kubectl port-forward svc/meshery -n meshery 9081:9081" >}}

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< installation/accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}
