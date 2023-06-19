---
layout: default
title: Windows
permalink: installation/platforms/windows
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/wsl2.png
---

{% include installation_prerequisites.html %}

# Overview

To set up and run Meshery on Windows

1. <a href="#step1">Configure Windows and enable Docker </a>
2. <a href="#step4">Install a Kubernetes cluster on top </a>
3. <a href="#step5">Run Meshery</a>

### Compatibility

The following minimum Windows build versions are required:

<table id="compatibility-table">
  <tr>
    <th id="model">Name</th>
    <th id="model">Version</th> 
  </tr>
  <tr>
    <td><a href="#wsl1">WSL1</a></td>
    <td><b>x64</b> - Windows 7 </td>
  </tr>
  <tr>
    <td><a href="#wsl2">WSL2</a></td>
    <td><b>x64</b> - Version 1903, Build 18362; <b>ARM 64</b> - Version 2004, Build 19041</td>
  </tr>
  <tr>
    <td><a href="https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-18945">Custom Kernel</a></td>
    <td>Build 18945</td>
  </tr>
  <tr>
    <td><a href="https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-19013">Kernel with K8s required modules</a></td>
    <td>Build 19013</td>
  </tr>
</table>

Note
<br />Run the following command on Powershell to check your Windows build and version:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">[System.Environment]::OSVersion.Version</div></div>
 </pre>

### Steps

Perform the following steps in order:

### 1. <a name="step1" href="https://docs.microsoft.com/en-us/windows/wsl/install-win10"><b>Install Windows Subsystem for Linux (WSL)</b></a>

Open Powershell in administrator mode and run:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart
Restart-Computer</div></div>
 </pre>

##### Choosing your WSL version:

<h6><b><a href="https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-18917" name="wsl2">WSL2</a></b> (Recommended)</h6>
Set the default version to *WSL2*, which will be inherited by any distro you wish to use.

Enable VM (Virtual Machine) feature:
Open PowerShell in administrator mode and run:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart</div></div>
 </pre>

Set WSL2 as the default version:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">wsl --set-default-version 2</div></div>
 </pre>

<h6><b><a href="https://docs.microsoft.com/en-us/windows/wsl/install-win10" name="wsl1"> WSL1 </a></b></h6>

<b>Warning:</b>
It is recommended to update to <a href="#wsl2">WSL2</a> as WSL1 doesn't support the Docker Desktop application for Windows. Instead, it only supports the deprecated version, [Docker Toolbox](https://docs.docker.com/toolbox/toolbox_install_windows/).

If you still wish to continue, follow the instructions for <button class="toggle-button" onclick="HideToggleFunction()"><b>WSL1</b></button>

<div id="hiddendiv">
<p>
1. The default version of WSL is set to WSL1 by default. You can move forward to <a href="https://docs.microsoft.com/en-us/windows/wsl/install-win10#install-your-linux-distribution-of-choice">install the distro</a> of your choice. <br /><br />

2. <b><a href="https://docs.docker.com/toolbox/toolbox_install_windows/">Docker Toolbox</a></b> <br />

<b>Warning</b>: Docker Toolbox is a deprecated version. It is recommended to update your system and install the Docker Desktop application with WSL2. <br/><br />

Docker Toolbox uses Linux-specific kernel features, and can’t run natively on Windows. Instead, it creates and uses a small Linux VM on your machine along with <a href="https://docs.docker.com/machine/overview/"><b>docker-machine</b></a>, and uses VirtualBox to run Docker. <br />

<ul>
<li> Go to <a href="https://github.com/docker/toolbox/releases">Toolbox Releases</a> and download the latest release <b>.exe</b> file </li>
<li> Follow these <a href="https://docs.docker.com/toolbox/toolbox_install_windows/#step-2-install-docker-toolbox">instructions</a> to successfully set up the Docker Toolbox application. </li>
</ul>

</p>
</div>

### 2. <b>[Install a new distro](https://docs.microsoft.com/en-us/windows/wsl/install-win10#install-your-linux-distribution-of-choice)</b>

In this tutorial, [Ubuntu 18.04](https://www.microsoft.com/en-us/p/ubuntu-1804-lts/9n9tngvndl3q?activetab=pivot:overviewtab) will be the distro used. Feel free to use any distro of your choice.

<strong>Note:</strong> If you choose to run Meshery without installing a distro, skip [step 2](#2-install-a-new-distro).

### 3. <b>Enable Docker</b>

The Docker Desktop application for Windows includes a comprehensive set of tools, including Docker Engine, Docker CLI client, Docker Compose, Notary, Kubernetes, and a Credential Helper.

<table id="compatibility-table">
  <tr>
    <th id="model">Windows 10 Version</th>
    <th id="model">Docker Desktop</th> 
  </tr>
  <tr>
    <td>Pro/Education/Enterprise</td>
    <td><a href="https://docs.docker.com/docker-for-windows/install/">Docker Desktop for Windows Pro</a></td>
  </tr>
  <tr>
    <td>Home</td>
    <td><a href="https://docs.docker.com/docker-for-windows/install-windows-home/">Docker Desktop for Windows Home</a></td>
  </tr>
</table>

### 4. <a name="step4"> <b>Install a Kubernetes cluster</b></a>

Once Docker is installed, the next step will be to install a Kubernetes cluster.
Under "Settings" in the Docker Desktop application, enable _Kubernetes_.
In this how-to, [K3d](https://github.com/rancher/k3d) will be used as it relies only on Docker.

  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
k3d cluster create
export KUBECONFIG="$(k3d kubeconfig get 'k3s-default')"</div></div></div></div>
  </pre>

If using Scoop, run the following in the PowerShell to install a Kubernetes cluster :

  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">scoop install k3d
k3d cluster create
$env:KUBECONFIG = "$(k3d.exe kubeconfig get 'k3s-default')"</div></div>
  </pre>

### 5. <a name="step5"><b>Set up Meshery</b></a>

Follow the [installation steps]({{ site.baseurl }}/installation#windows) to install the mesheryctl CLI. Then, execute:

  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">./mesheryctl system start</div></div>
  </pre>

Meshery server supports customizing authentication flow callback URL, which can be configured in the following way

  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host ./mesheryctl system start</div></div>
  </pre>

Type **yes** when prompted to choose to configure a file. To get started, choose Docker as your platform to deploy Meshery.

### 6. <a name="step6"> <b>Aftermath errors</b></a>

Meshery is unable to detect the Kubernetes connection running on your local system, even after manually uploading the `.kube config` file :

1. To display a list of the contexts in Kubernetes, run:

  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">kubectl config get-contexts</div></div>
  </pre>

2. Change your current context to docker-desktop:

  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">kubectl config use-context
docker-desktop</div></div>
  </pre>

3. Restart Meshery.

Stuck at another error? [Tell us about it](http://slack.layer5.io/)
