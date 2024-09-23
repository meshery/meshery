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

{% include alert.html type="warning" title="Quick Start Assumptions" content="This quick start guide enables you to download, install, and run Meshery in a single command. See all <a href='/installation'>supported platforms</a> for more specific (and less presumptious) instructions." %}

## 1. Download, install, and run Meshery

If you are on macOS or Linux system, you can download, install, and run both `mesheryctl` and Meshery Server by executing the following command.

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
<!-- <script src="/assets/js/terminal.js" data-termynal-container="#termynal0"></script> -->

<pre class="codeblock-pre">
  <div class="codeblock">
  <div class="clipboardjs">$ curl -L https://meshery.io/install | PLATFORM=kubernetes bash -</div>
  </div>
</pre>

The script above downloads Meshery's command line interface, <code>mesheryctl</code>.

{% include alert.html type="info" title="Meshery CLI" content="Meshery's command line interface, <code>mesheryctl</code>, can be installed in <a href='/installation/mesheryctl'>various ways</a>. In addition to <a href='/installation/linux-mac/bash'>Bash</a>, you can also use <a href='/installation/linux-mac/brew'>Brew</a> or <a href='/installation/windows/scoop'>Scoop</a> to install <code>mesheryctl</code>. Alternatively, <code>mesheryctl</code> is also available via <a href='https://github.com/meshery/meshery/releases/latest'>direct download</a>." %}

## 2. Access Meshery

Your default browser will be opened and directed to Meshery's web-based user interface typically found at `http://localhost:9081`.

{% include alert.html type="light" title="Accessing Meshery Server with Meshery UI" content="Meshery's web-based user interface is embedded in Meshery Server and is available as soon as Meshery starts. The location and port that Meshery UI is exposed varies depending upon your mode of deployment. See <a href='/installation/accessing-meshery-ui'>accessing Meshery UI</a> for deployment-specific details." %}

{% include alert.html type="light" title="Accessing Meshery Server with Meshery CLI" content="Meshery's command line interface is a client of Meshery Server's REST API (just as Meshery UI is). Choose to use <code>mesheryctl</code> as an alternative client as it suits your needs." %}

### 3. Select a Provider

Select from the list of [Providers]({{ site.baseurl }}/extensibility/providers) in order to login to Meshery. Authenticate with your chosen Provider.

<a href="/assets/img/meshery-server-page.png">
  <img class="center" style="width:min(100%,650px)" src="/assets/img/meshery-server-page.png" />
</a>

## 4. Configure Connections to your Kubernetes Clusters

**Out-of-Cluster Deployments**
If you have deployed Meshery out-of-cluster, Meshery Server will automatically attempt to connect to any available Kubernetes clusters found in your kubeconfig (under `$HOME/.kube/config`) and in kubeconfigs uploaded through Meshery UI. Meshery Server deploys [Meshery Operator](/concepts/architecture/operator), [MeshSync](/concepts/architecture/meshsync), and Broker into the `meshery` namespace (by default).

**In-Cluster Deployments**
If you have deployed Meshery in-cluster, Meshery Server will automatically connect to the Kubernetes API Server availabe in the control plane.

Visit <i class="fas fa-cog"></i> Settings:

<a href="/assets/img/platforms/meshery-settings.png">
  <img class="center" style="width:min(100%,650px);" src="/assets/img/platforms/meshery-settings.png" />
</a>

If your config has not been autodetected, you can manually upload your kubeconfig file (or any number of kubeconfig files). By default, Meshery will attempt to connect to and deploy Meshery Operator to each reachable context contained in the imported kubeconfig files. See Managing Kubernetes Clusters for more information.

## 5. Verify Deployment

Run connectivity tests and verify the health of your Meshery system. Verify Meshery's connection to your Kubernetes clusters by clicking on the connection chip. A quick connectivity test will run and inform you of Meshery's ability to reach and authenticate to your Kubernetes control plane(s). You will be notified of your connection status. You can also verify any other connection between Meshery and either its components (like [Meshery Adapters]({{ site.baseurl }}/concepts/architecture/adapters)) or other managed infrastructure by clicking on any of the connection chips. When clicked, a chip will perform an ad hoc connectivity test.

<a href="{{site.baseurl}}/assets/img/platforms/k8s-context-switcher.png" alt="Meshery Kubernetes Context Switcher">
  <img class="center" style="width:min(100%,350px);" src="{{site.baseurl}}/assets/img/platforms/k8s-context-switcher.png" />
</a>

## 5. Design and operate Kubernetes clusters and their workloads

You may now proceed to managed any cloud native infrastructure supported by Meshery. See all integrations for a complete list of supported infrastructure.

<a href="{{site.baseurl}}/assets/img/platforms/meshery-designs.png">
  <img class="center" style="width:min(100%,650px);" src="{{site.baseurl}}/assets/img/platforms/meshery-designs.png" />
</a>

# Additional Guides

<div class="section">
    <ul>
        <li><a href="{{ site.baseurl }}/guides/troubleshooting/installation">Troubleshooting Meshery Installations</a></li>
        <li><a href="{{ site.baseurl }}/reference/error-codes">Meshery Error Code Reference</a></li>
        <li><a href="{{ site.baseurl }}/reference/mesheryctl/system/check">Mesheryctl system check</a></li> 
    </ul>
</div>

