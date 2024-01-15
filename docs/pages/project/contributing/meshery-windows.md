---
layout: page
title: Setting up Meshery Development Environment on Windows
permalink: project/contributing/meshery-windows
abstract: How to set up Meshery Development Environment on Windows
language: en
type: project
category: contributing
--- 

Meshery can be run on the Windows platform in three different modes.
1. Windows native mode
2. Under WSL or WSL2
3. Using Hyper-V to install a Linux VM to run Meshery.

Running Meshery on WSL(2) or Hyper-V is not much different from running Meshery on a native Linux platform. The same instructions can be followed as already documented.

Essentially, all three platforms require the same prerequisites but the Linux platforms (as there are more than one of them) come with some of the prerequisites that Windows does not.

In general, all three require that Docker and Kubernetes are present. If Docker Desktop is installed, then Docker actually uses Hyper-V to install a VM. For the requisite images, both Linux and Windows containers are required depending on what one wants to do but the Linux images are certainly required. This is where the exception that was mentioned earlier comes into play, that is if a Kubernetes elsewhere is used. This installation is specified in the config file in the .kube directory in the home directory.

## Prerequisite Dependencies

All three platforms require the following prerequisites:

- make
- gcc
- nodejs (for npm)
- go

## Using WSL2

Here is a link to setting up WSL2: [https://learn.microsoft.com/en-us/windows/wsl/install](https://learn.microsoft.com/en-us/windows/wsl/install)

## On Windows Native Mode

Installing Meshery on native Windows may be a little more involved than installing Meshery on a Linux platform. That is because Linux installations usually come with a lot of the prerequisites preinstalled. Moreover, one is given a hint on the Linux platform on what to do next including the command, especially on Ubuntu.

On Windows, make, gcc and some of the other prerequisites do not come pre-installed, These will need to be installed manually. Also, after installing any of the items, it may be necessary to set its path. This can be done in the control panel by searching for setting the system environment variables.

Please research the internet to install these components. However, here are a few suggestions:

- make: [https://gnuwin32.sourceforge.net/packages/make.htm](https://gnuwin32.sourceforge.net/packages/make.htm)
- gcc: [https://gcc.gnu.org/install/binaries.html](https://gcc.gnu.org/install/binaries.html)
- nodejs (for npm): [https://phoenixnap.com/kb/install-node-js-npm-on-windows](https://phoenixnap.com/kb/install-node-js-npm-on-windows)
- go: [https://go.dev/doc/install](https://go.dev/doc/install)


## Compatibility

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
       <li>  Go to <a href="https://github.com/docker/toolbox/releases">Toolbox Releases</a> and download the latest release <b>.exe</b> file </li>
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





## Using Hyper-V to install a linux VM

Here is a link to install linux VM using Hyper-V: [https://wiki.ubuntu.com/Hyper-V](https://wiki.ubuntu.com/Hyper-V)

### Git Bash  
Git Bash is a terminal emulator which provides git command line experience. This will make working with git easier. You can download it from here: [https://git-scm.com/downloads](https://git-scm.com/downloads)

### Get the code

- Fork and then clone the [Meshery repository](https://github.com/layer5io/meshery)
  ```bash
  $ git clone https://github.com/YOUR-USERNAME/meshery
  ```