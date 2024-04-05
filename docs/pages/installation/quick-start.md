---
layout: default
title: Quick Start Guide
permalink: installation/quick-start
# redirect_from: installation/quick-start/
language: en
type: installation
abstract: Getting Meshery up and running locally on a Docker-enabled system or in Kubernetes is easy. Meshery deploys as a set of Docker containers, which can be deployed to either a Docker host or Kubernetes cluster.
---

<a name="getting-started"></a>

Getting Meshery up and running locally on a Docker-enabled system or in Kubernetes is easy. Meshery deploys as a set of Docker containers, which can be deployed to either a Docker host or Kubernetes cluster.

{% include alert.html type="info" title="All Supported Platforms" content="Download, install, and run Meshery in a single command. See all <a href='/installation'>supported platforms</a>." %}

## Install Meshery

To install and start Meshery, begin by [installing mesheryctl]({{ site.baseurl }}/installation/mesheryctl). If you are on macOS or Linux system, you can download, install, and run both `mesheryctl` and Meshery Server with the command shown in the figure.

 <pre class="codeblock-pre" style="padding: 0; font-size:0px;"><div class="codeblock" style="display: block;">
 <div class="clipboardjs" stytle="padding: 0">
 <span style="font-size: 0;">curl -L https://meshery.io/install | PLATFORM=kubernetes bash -</span>
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

 <script src="/assets/js/terminal.js" data-termynal-container="#termynal0"></script>

{% include alert.html type="info" title="Meshery CLI Package Managers" content="In addition to <a href='/installation/linux-mac/bash'>Bash</a>, you can also use <a href='/installation/linux-mac/brew'>Brew</a> or <a href='/installation/windows/scoop'>Scoop</a> to install <code>mesheryctl</code>. Alternatively, <code>mesheryctl</code> is also available <a href='https://github.com/meshery/meshery/releases/latest'>direct download</a>." %}

## Install using Docker Meshery Extension

You can find the [Docker Meshery Extension in Docker Hub](https://hub.docker.com/extensions/meshery/docker-extension-meshery) marketplace or use the Extensions browser in Docker Desktop to install the Docker Meshery Extension.

[![Docker Meshery Extension]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension.png)]({{site.baseurl}}/assets/img/platforms/docker-desktop-meshery-extension.png)

## Access Meshery

Your default browser will be opened and directed to Meshery's web-based user interface typically found at `http://localhost:9081`.

{% include alert.html type="info" title="Accessing Meshery UI" content="Meshery's web-based user interface is embedded in Meshery Server and is available as soon as Meshery starts. The location and port that Meshery UI is exposed varies depending upon your mode of deployment. See <a href='/tasks/accessing-meshery-ui'>accessing Meshery UI</a> for deployment-specific details." %}

### Select a Provider

Select from the list of [Providers]({{ site.baseurl }}/extensibility/providers) in order to login to Meshery. Authenticate with your chosen Provider.

<a href="/assets/img/meshery-server-page.png">
  <img class="center" style="width:min(100%,650px)" src="/assets/img/meshery-server-page.png" />
</a>

### Configure Connections to your Kubernetes Clusters

If you have deployed Meshery out-of-cluster, Meshery will automatically connect to any available Kubernetes clusters found in your kubeconfig (under `$HOME/.kube/config`). If you have deployed Meshery out-of-cluster, Meshery will automatically connect to the Kubernetes API Server availabe in the control plane. Ensure that Meshery is connected to one or more of your Kubernetes clusters.

Visit <i class="fas fa-cog"></i> Settings:

<a href="/assets/img/platforms/meshery-settings.png">
  <img class="center" style="width:min(100%,650px);" src="/assets/img/platforms/meshery-settings.png" />
</a>

If your config has not been autodetected, you can manually upload your kubeconfig file (or any number of kubeconfig files). By default, Meshery will attempt to connect to and deploy Meshery Operator to each reachable context contained in the imported kubeconfig files. See Managing Kubernetes Clusters for more information.

### Verify Deployment

Run connectivity tests and verify the health of your Meshery system. Verify Meshery's connection to your Kubernetes clusters by clicking on the connection chip. A quick connectivity test will run and inform you of Meshery's ability to reach and authenticate to your Kubernetes control plane(s). You will be notified of your connection status. You can also verify any other connection between Meshery and either its components (like [Meshery Adapters]({{ site.baseurl }}/concepts/architecture/adapters)) or other managed infrastructure by clicking on any of the connection chips. When clicked, a chip will perform an ad hoc connectivity test.

<a href="{{site.baseurl}}/assets/img/platforms/k8s-context-switcher.png" alt="Meshery Kubernetes Context Switcher">
  <img class="center" style="width:min(100%,350px);" src="{{site.baseurl}}/assets/img/platforms/k8s-context-switcher.png" />
</a>

### Design and operate Kubernetes clusters and their workloads

You may now proceed to managed any cloud native infrastructure supported by Meshery. See all integrations for a complete list of supported infrastructure.

<a href="{{site.baseurl}}/assets/img/platforms/meshery-designs.png">
  <img class="center" style="width:min(100%,650px);" src="{{site.baseurl}}/assets/img/platforms/meshery-designs.png" />
</a>

## Additional Guides

<div class="section">
    <ul>
        <li><a href="{{ site.baseurl }}/guides/troubleshooting/installation">Troubleshooting Meshery Installations</a></li>
        <li><a href="{{ site.baseurl }}/reference/error-codes">Meshery Error Code Reference</a></li>
        <li><a href="{{ site.baseurl }}/reference/mesheryctl/system/check">Mesheryctl system check</a></li> 
    </ul>
</div>
<script src="/assets/js/terminal.js" data-termynal-container="#termynal0|#termynal1|#termynal2"></script>

