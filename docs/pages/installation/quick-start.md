---
layout: default
title: Quick Start Guide
permalink: installation/quick-start
redirect_from: installation/quick-start/
language: en
---
<!-- <link href="/assets/css/terminal.css" rel="stylesheet"> -->
<a name="getting-started"></a>

Getting Meshery up and running locally on a Docker-enabled system or in Kubernetes is easy. Meshery deploys as a set of Docker containers, which can be deployed to either a Docker host or Kubernetes cluster. See the complete list of its [supported platforms]({{ site.baseurl }}/installation/platforms).

## Install Meshery

Use the Meshery command line interface, [mesheryctl]({{ site.baseurl }}/guides/mesheryctl), to install and start Meshery. Begin with Meshery by installing its command line client: [mesheryctl]({{ site.baseurl }}/guides/mesheryctl). If you are on a MacOS or Linux system, you can download, install, and run to the management plane with the command shown in the figure.

### Install on Kubernetes

 <pre class="codeblock-pre" style="padding: 0; font-size:0px;"><div class="codeblock" style="display: block;">
 <div class="clipboardjs" style="padding: 0">
    <span style="font-size:0;">curl -L https://meshery.io/install | PLATFORM=kubernetes bash -</span> 
 </div>
 <div class="window-buttons"></div>
  <div id="termynal0" style="width:100%; height:150px; max-width:100%;" data-termynal="">
            <span data-ty="input">curl -L https://meshery.io/install | PLATFORM=kubernetes bash -</span>
            <span data-ty="progress"></span>
            <span data-ty="">Successfully installed Meshery</span>
            <span data-ty="input">mesheryctl system dashboard</span>
  </div>
  </div>
 </pre>
 <br>

{% include alert.html type="dark" title="All Supported Platforms" content="Don't find an answer to your question here? Ask on the <a href='https://discuss.layer5.io'>Discussion Forum</a>." %}


### Install using Docker Extension for Meshery
You can visit the [Docker Hub](https://hub.docker.com/extensions/meshery/docker-extension-meshery) marketplace to directly install Meshery extension in your Docker Desktop.


_Download, install, and run Meshery in a single command. See all [installation methods]({{ site.baseurl }}/installation)._

## Access Meshery

Visit Meshery's web-based user interface `http://<hostname>:9081`. When Meshery is installed on a Docker host, your default browser will be opened and navigated to `http://localhost:9081`.

## Select a Provider

Select from the list of [Providers]({{site.baseurl}}/extensibility#providers) in order to login to Meshery. Authenticate with your chosen Provider.

<a href="/assets/img/meshery-server-page.png">
  <img  style="width:min(100%,650px)" src="/assets/img/meshery-server-page.png" />
</a>

## Configure Connection to Kubernetes

Meshery attempts to auto detect your kubeconfig if it is stored in the default path (`$HOME/.kube`) on your system. In most deployments, Meshery will automatically connect to your Kubernetes cluster. Ensure that Meshery is connected to your Kubernetes cluster.

Visit <i class="fas fa-cog"></i> Settings:

  <a href="/assets/img/adapters/meshery-settings.png">
  <img style="width:min(100%,650px);" src="/assets/img/adapters/meshery-settings.png" />
  </a>

If your config has not been auto-detected, you may manually locate and upload your **kube config** file and select the **context name** (docker-desktop, kind-clsuter, minikube etc.)

## Verify Meshery's Deployment

Run connectivity tests and verify the health of your Meshery system. Verify Meshery's connection to Kubernetes by clicking on your configuration `context` name. You will be notified of your connection status. You can also verify a successful connection between Meshery and its adapters by clicking on any of the available [Adapter Ports]({{ site.baseurl }}/concepts/architecture#adapter-ports).

<a href="/assets/img/adapters/meshery-ui.png">
<img style="width:min(100%,650px);height=auto;" src="/assets/img/adapters/meshery-ui.png" />
</a>

## Operate service meshes and their workloads

You may now proceed to install and work with any [service mesh]({{ site.baseurl }}/service-meshes) supported by Meshery.

<video class="videoTest" style="width:min(100%,650px)" height="auto" autoplay muted loop>
  <source src="/assets/img/adapters/meshery-ui-setup.mp4" type="video/mp4">
 Your browser does not support the video tag
</video>

## Additional Guides

<div class="section">
    <ul>
        <li><a href="{{ site.baseurl }}/guides/troubleshooting/installation">Troubleshooting Meshery Installations</a></li>
        <li><a href="{{ site.baseurl }}/reference/error-codes">Meshery Error Code Reference</a></li>
        <li><a href="{{ site.baseurl }}/reference/mesheryctl/system/check">Mesheryctl system check</a></li> 
    </ul>
</div>
<script src="/assets/js/terminal.js" data-termynal-container="#termynal0|#termynal1|#termynal2"></script>
