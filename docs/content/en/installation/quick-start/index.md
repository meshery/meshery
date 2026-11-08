---
title: Quick Start Guide
description: Install mesheryctl and start Meshery on a Kubernetes cluster or a Docker host with a single command, then connect your clusters and start designing.
weight: -10
---

<a name="getting-started"></a>

Meshery runs as a set of containers on either a Kubernetes cluster or a Docker host. This guide installs Meshery's command line interface, `mesheryctl`, and uses it to start Meshery Server, then walks you through signing in and connecting your clusters.

{{% alert color="warning" title="Quick Start Assumptions" %}}
This guide downloads, installs, and runs Meshery with a single command, and assumes that your current kubeconfig context points at the cluster you want to use. For instructions specific to your environment, see all <a href='{{< ref "installation/_index.md" >}}'>supported platforms</a>.
{{% /alert %}}

## 1. Download, install, and run Meshery

On macOS or Linux, run the following command to download and install `mesheryctl` and then deploy Meshery Server to the Kubernetes cluster in your current kubeconfig context.

<!-- <pre class="codeblock-pre" style="padding: 0; font-size:0px;">
<div class="codeblock" style="display: block;">
  <div class="clipboardjs" style="visibility:hidden;padding: 0;">
    <span style="visibility:hidden">curl -L https://meshery.io/install | PLATFORM=kubernetes bash -</span>
  </div>
  <div class="window-buttons"></div>
  <div id="termynal0" style="width:fit-content;min-height:content-fit;" data-termynal="">
    <span data-ty="input">curl -L https://meshery.io/install | PLATFORM=kubernetes bash -</span>
  </div>
</div>
</pre>-->
<!-- <script src="../../assets/js/terminal.js" data-termynal-container="#termynal0"></script> -->

<pre class="codeblock-pre">
  <div class="codeblock">
  <div class="clipboardjs">curl -L https://meshery.io/install | PLATFORM=kubernetes bash -</div>
  </div>
</pre>
<br/>

To run Meshery on your local Docker host instead, replace `PLATFORM=kubernetes` with `PLATFORM=docker`. If you omit `PLATFORM`, the script prompts you to choose a platform.

On Windows, add the Meshery bucket to <a href="https://scoop.sh">Scoop</a>, install `mesheryctl`, and start Meshery:

<pre class="codeblock-pre">
  <div class="codeblock">
  <div class="clipboardjs">scoop bucket add mesheryctl https://github.com/meshery/scoop-bucket.git
scoop install mesheryctl
mesheryctl system start</div>
  </div>
</pre>

`mesheryctl system start` deploys to the platform set in your current `mesheryctl` context, which defaults to `kubernetes`. Add `--platform docker` to deploy to Docker instead.

{{% alert color="info" title="Meshery CLI" %}}
You can install <code>mesheryctl</code> in <a href='{{< ref "installation/mesheryctl/_index.md" >}}'>several ways</a>: with <a href='{{< ref "installation/mesheryctl/linux-mac/bash.md" >}}'>Bash</a>, <a href='{{< ref "installation/mesheryctl/linux-mac/brew.md" >}}'>Brew</a>, or <a href='{{< ref "installation/mesheryctl/windows/scoop.md" >}}'>Scoop</a>, or by <a href='https://github.com/meshery/meshery/releases/latest'>direct download</a> from the latest release.
{{% /alert %}}

## 2. Access Meshery

When Meshery is ready, `mesheryctl` opens Meshery UI in your default browser. On Docker, Meshery UI listens at `http://localhost:9081`. On Kubernetes, `mesheryctl` opens the address of the `meshery` service in the `meshery` namespace. If your browser cannot reach that address, run `mesheryctl system dashboard --port-forward` to forward Meshery UI to `http://localhost:9081`. To skip opening the browser, pass `--skip-browser` to `mesheryctl system start`.

{{% alert color="light" title="Accessing Meshery Server with Meshery UI" %}}
Meshery Server serves Meshery UI, so the UI is available as soon as Meshery starts. The address and port depend on how you deployed Meshery. See <a href='{{< ref "installation/accessing-meshery-ui.md" >}}'>accessing Meshery UI</a> for deployment-specific details.
{{% /alert %}}

{{% alert color="light" title="Accessing Meshery Server with Meshery CLI" %}}
Like Meshery UI, <code>mesheryctl</code> is a client of Meshery Server's REST API. Use whichever client suits the task at hand.
{{% /alert %}}

## 3. Select a Provider

Choose one of the listed [Providers]({{< ref "reference/extensibility/providers/index.md" >}}) and sign in with it.

<a href="./images/meshery-server-page.png">
  <img class="center" style="width:min(100%,650px)" src="./images/meshery-server-page.png" />
</a>

## 4. Configure Connections to your Kubernetes Clusters

**Out-of-Cluster Deployments**
If you deployed Meshery out-of-cluster, Meshery Server automatically tries to connect to the Kubernetes clusters in your kubeconfig (`$HOME/.kube/config`) and in any kubeconfig files you upload through Meshery UI. By default, [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) runs embedded inside Meshery Server and nothing is installed into your cluster. Only when a connection uses operator mode does Meshery Server deploy [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}), MeshSync, and Broker into the `meshery` namespace. To switch a connection between the two, see [MeshSync deployment modes]({{< ref "guides/infrastructure-management/kubernetes-connection-lifecycle.md#meshsync-deployment-modes" >}}).

**In-Cluster Deployments**
If you deployed Meshery in-cluster, Meshery Server automatically connects to the API server of the cluster it runs in.

Visit <i class="fas fa-cog"></i> Settings:

<a href="./images/meshery-settings.png">
  <img class="center" style="width:min(100%,650px);" src="./images/meshery-settings.png" />
</a>

If Meshery did not detect your kubeconfig, upload one or more kubeconfig files yourself. Meshery tries to connect to each reachable context in the uploaded files, using embedded mode unless you set that context's connection to operator mode. For more information, see [Kubernetes connection lifecycle]({{< ref "guides/infrastructure-management/kubernetes-connection-lifecycle.md" >}}).

## 5. Verify Deployment

Run connectivity tests to verify the health of your Meshery deployment. To check Meshery's connection to a Kubernetes cluster, click its connection chip. Meshery runs a quick connectivity test and tells you whether it can reach and authenticate to that cluster's control plane. Every connection chip works the same way, whether it represents one of Meshery's own components, like [Meshery Adapters]({{< ref "concepts/architecture/adapters.md" >}}), or other managed infrastructure. From the command line, run `mesheryctl system check` for pre- and post-deployment health checks.

If Meshery Operator does not reach a healthy status, its status card and the connection's Diagnostics show the reason. Deploying the Operator requires outbound access to the Meshery chart repository. Meshery retries the chart lookup on each deploy, so if the repository was briefly unreachable, redeploy the Operator. See [Meshery Operator and MeshSync troubleshooting]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}}).

<a href="./images/k8s-context-switcher.png" alt="Meshery Kubernetes Context Switcher">
  <img class="center" style="width:min(100%,350px);" src="./images/k8s-context-switcher.png" />
</a>

## 6. Design and operate Kubernetes clusters and their workloads

You can now design and manage any infrastructure that Meshery supports. For the complete list, see all [integrations]({{< ref "extensions/models/_index.md" >}}).

<a href="./images/meshery-designs.png">
  <img class="center" style="width:min(100%,650px);" src="./images/meshery-designs.png" />
</a>

## Explore Tutorials

🧑‍🔬 Work through these tutorials to learn how to manage infrastructure collaboratively with Meshery.

{{< tutorials-list >}}

## Additional Resources

- [Troubleshooting Meshery Installations]({{< ref "guides/troubleshooting/installation.md" >}})
- [Meshery Error Code Reference]({{< ref "reference/error-codes.md" >}})
- [mesheryctl system check]({{< ref "reference/mesheryctl/system/check.md" >}})
