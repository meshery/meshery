---
title: "Quick Start Guide"
description: "Getting Meshery up and running locally on a Docker-enabled system or in Kubernetes"
weight: 1
aliases:
  - /installation/platforms/quick-start
---

<a name="getting-started"></a>

Getting Meshery up and running locally on a Docker-enabled system or in Kubernetes is easy. Meshery deploys as a set of Docker containers, which can be deployed to either a Docker host or Kubernetes cluster.

> [!WARNING]
> **Quick Start Assumptions**
> This quick start guide enables you to download, install, and run Meshery in a single command. See all [supported platforms](/installation) for more specific (and less presumptuous) instructions.

## 1. Download, install, and run Meshery

If you are on a macOS or Linux system, you can download, install, and run both `mesheryctl` and Meshery Server by executing the following command.

{{< code >}}
$ curl -L https://meshery.io/install | PLATFORM=kubernetes bash -
{{< /code >}}

> [!NOTE]
> **Meshery CLI**
> Meshery's command line interface, `mesheryctl`, can be installed in [various ways](/installation/mesheryctl). In addition to [Bash](/installation/linux-mac/bash), you can also use [Brew](/installation/linux-mac/brew) or [Scoop](/installation/windows/scoop) to install `mesheryctl`. Alternatively, `mesheryctl` is also available via [direct download](https://github.com/meshery/meshery/releases/latest).

## 2. Access Meshery

Your default browser will be opened and directed to Meshery's web-based user interface typically found at `http://localhost:9081`.

> [!NOTE]
> **Accessing Meshery Server with Meshery UI**
> Meshery's web-based user interface is embedded in Meshery Server and is available as soon as Meshery starts. The location and port that Meshery UI is exposed varies depending upon your mode of deployment. See [accessing Meshery UI](/installation/accessing-meshery-ui) for deployment-specific details.

> [!NOTE]
> **Accessing Meshery Server with Meshery CLI**
> Meshery's command line interface is a client of Meshery Server's REST API (just as Meshery UI is). Choose to use `mesheryctl` as an alternative client as it suits your needs.

## 3. Select a Provider

Select from the list of [Providers](/extensibility/providers) in order to log in to Meshery. Authenticate with your chosen Provider.

<a href="/images/meshery-server-page.png">
  <img class="center" style="width:min(100%,650px)" src="/images/meshery-server-page.png" />
</a>

## 4. Configure Connections to your Kubernetes Clusters

**Out-of-Cluster Deployments**
If you have deployed Meshery out-of-cluster, Meshery Server will automatically attempt to connect to any available Kubernetes clusters found in your kubeconfig (under `$HOME/.kube/config`) and in kubeconfigs uploaded through Meshery UI. Meshery Server deploys [Meshery Operator](/concepts/architecture/operator), [MeshSync](/concepts/architecture/meshsync), and Broker into the `meshery` namespace (by default).

**In-Cluster Deployments**
If you have deployed Meshery in-cluster, Meshery Server will automatically connect to the Kubernetes API Server available in the control plane.

Visit <i class="fas fa-cog"></i> Settings:

<a href="/images/platforms/meshery-settings.png">
  <img class="center" style="width:min(100%,650px);" src="/images/platforms/meshery-settings.png" />
</a>

If your config has not been autodetected, you can manually upload your kubeconfig file (or any number of kubeconfig files). By default, Meshery will attempt to connect to and deploy Meshery Operator to each reachable context contained in the imported kubeconfig files. See [Managing Kubernetes Clusters](/installation/kubernetes) for more information.

## 5. Verify Deployment

Run connectivity tests and verify the health of your Meshery system. Verify Meshery's connection to your Kubernetes clusters by clicking on the connection chip. A quick connectivity test will run and inform you of Meshery's ability to reach and authenticate to your Kubernetes control plane(s). You will be notified of your connection status. You can also verify any other connection between Meshery and either its components (like [Meshery Adapters](/concepts/architecture/adapters)) or other managed infrastructure by clicking on any of the connection chips. When clicked, a chip will perform an ad hoc connectivity test.

<a href="/images/platforms/k8s-context-switcher.png" alt="Meshery Kubernetes Context Switcher">
  <img class="center" style="width:min(100%,350px);" src="/images/platforms/k8s-context-switcher.png" />
</a>

## 6. Design and operate Kubernetes clusters and their workloads

You may now proceed to manage any cloud native infrastructure supported by Meshery. See all [integrations](/extensibility/integrations) for a complete list of supported infrastructure.

<a href="/images/platforms/meshery-designs.png">
  <img class="center" style="width:min(100%,650px);" src="/images/platforms/meshery-designs.png" />
</a>

## Additional Guides

- [Troubleshooting Meshery Installations](/guides/troubleshooting/installation)
- [Meshery Error Code Reference](/reference/error-codes)
- [mesheryctl system check](/reference/mesheryctl/system/check)

{{< related-discussions tag="meshery" >}}
