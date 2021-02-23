---
layout: default
title: Quick Start Guide
permalink: installation/quick-start
redirect_from: installation/
language: en
list: exclude
---

<a name="getting-started"></a>

# Meshery Quick Start Guide

Getting Meshery up and running locally on a Docker-enabled system is easy with Meshery's command line interface, <a href="/docs/guides/mesheryctl">mesheryctl</a>.

## Configure Your Environment

Meshery deploys as a set of Docker containers, which can be deployed to either a Docker host or Kubernetes cluster. See the complete list of its [supported platforms](/docs/installation/platforms). 

## Install Meshery

Use the Meshery command line interface, , to install and start Meshery. Begin with Meshery by installing its command line client: [mesheryctl](/docs/guides/mesheryctl). If you are on a MacOS or Linux system, you can download, install, and run to the management plane with the command shown in the figure.

```
$ curl -L https://git.io/meshery | bash - 
```

_Download, install, and run Meshery in a single command._

## Access Meshery

Visit Meshery's web-based user interface `http://<hostname>:9081`. When Meshery is installed on a Docker host, your default browser will be opened and navigated to `http://localhost:9081`.

## Select a Provider
Select from the list of [Providers](/docs/reference/extensibility#providers) in order to login to Meshery. Authenticate with your chosen your Provider.

<a href="/docs/assets/img/meshery-server-page.png">
  <img style="width:300px;" src="/docs/assets/img/meshery-server-page.png" />
</a>


## Configure Connection to Kubernetes
Meshery attempts to auto detect your kubeconfig if it is stored in the default path (`$HOME/.kube`) on your system. In most deployments, Meshery will automatically connect to your Kubernetes cluster. Ensure that Meshery is connected to your your Kubernetes cluster. 

Visit <i class="fas fa-cog"></i> Settings:

  <a href="/docs/assets/img/adapters/meshery-settings.png">
  <img style="width:600px;" src="/docs/assets/img/adapters/meshery-settings.png" />
  </a>

If your config has not been auto-detected, you may manually locate and upload your **kube config** file and select the **context name** (docker-desktop, kind-clsuter, minikube etc.)

## Verify Meshery's Deployment
Run connectivity tests and verify the health of your Meshery system. Verify Meshery's connection to Kubernetes by clicking on your configuration `context` name. You will be notified of your connection status. You can also verify a successful connection between Meshery and its adapters by clicking on any of the available [Adapter Ports](/docs/concepts/architecture#adapter-ports).

<a href="/docs/assets/img/adapters/meshery-ui.png">
<img style="width:450px;height=auto;" src="/docs/assets/img/adapters/meshery-ui.png" />
</a>


## Operate service meshes and their workloads
You may now proceed to install and work with any [service mesh](/docs/service-meshes) supported by Meshery.

<video class="videoTest" width="750" height="auto" autoplay muted loop>
  <source src="/docs/assets/img/adapters/meshery-ui-setup.mp4" type="video/mp4">
 Your browser does not support the video tag
</video>
