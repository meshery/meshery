---
title: Setting up Meshery Development Environment on Windows
description: How to set up Meshery Development Environment on Windows
category: [contributing]
--- 

You can use one of the following three approaches to setup a Meshery development environment in Windows:
1. Using WSL or WSL2 (Recommended)
2. Natively (directly) on Windows
3. Using Hyper-V to create a Linux VM

Setting up a Meshery development environment on WSL(2) or Hyper-V is not much different from setting up on a native Linux platform. The same instructions can be followed as already documented.

All three methods require the same prerequisites mentioned below in the [Prerequisites](#prerequisites) section but the Linux-based methods/platforms come with some of the prerequisites pre-installed and have no issues with Linux style commands making it easier to setup the environment.

Additionally, all three require methods will require access to a Kubernetes cluster which can be a kind, minikube, Kubernetes in Docker or a cluster elsewhere with network access.

## Prerequisites

All three methods require the following prerequisites:
- make
- gcc
- nodejs (for npm)
- go
- Git or Git Bash
- Visual Studio Code for development

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

## Using WSL2

Here is a link to setting up WSL2: [https://learn.microsoft.com/en-us/windows/wsl/install](https://learn.microsoft.com/en-us/windows/wsl/install)

Once WSL is setup, login to the WSL distribution and use the steps for Linux to setup the development environment. 

### Steps
Perform the following steps in order:

#### 1. Enable and configure Windows Subsystem for Linux (WSL)

Open Powershell in administrator mode and run:
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart

Restart-Computer</div></div>
</pre>

For more information visit [How to install Linux on Windows with WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10).

#### 2. Configuring WSL

Enable VM (Virtual Machine) feature. To do so open _PowerShell_ in administrator mode and run:
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart</div></div>
</pre>

Ensure WSL2 is the default version:
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">wsl --set-default-version 2</div></div>
</pre>

#### 3. Install a new distro

Refer to the instructions [here](https://docs.microsoft.com/en-us/windows/wsl/install-win10#install-your-linux-distribution-of-choice) to install a new distro such as Ubuntu by default.
Feel free to use any distro of your choice.

#### 4. Enable Docker

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

## Natively on Windows

This approach can be the most challenging to setup among all and is recommended for users that are very conversant with setting up development environments in Windows.
make, gcc and the other prerequisites do not come pre-installed on Windows and hence these need to be installed manually. You will also need to set the **Path** in _Environment Variables_ for some of them like _make_ and _gcc_. 

### Installing prerequsites:
Here are some links and recommendations to install the prerequisites:

#### - make 

Use the [setup program](https://gnuwin32.sourceforge.net/downlinks/make.php) to download the installation package and run it to install make and other dependencies. 

If the download link fails, visit the home page at [https://gnuwin32.sourceforge.net/packages/make.htm](https://gnuwin32.sourceforge.net/packages/make.htm)

#### - gcc

Use the [mingw-w64](https://www.mingw-w64.org/) project for gcc. The pre-built binaries can be downloaded from [https://github.com/niXman/mingw-builds-binaries/releases](https://github.com/niXman/mingw-builds-binaries/releases). For a typical 64-bit Windows, it is of the **x86_64-XX.X.X-release-posix-seh-ucrt-XX_XXX-XXXX.7z** format. Posix provides cross-platform compatibility and is known to work well. After extracting and copying (typically to the C:\ drive at root), add the path to the **Path** variable in _Environment Variables_. 

For example, here's a [direct download link](https://github.com/niXman/mingw-builds-binaries/releases/download/15.2.0-rt_v13-rev0/x86_64-15.2.0-release-posix-seh-ucrt-rt_v13-rev0.7z) to the latest release as of this writing. After installing, ensure the path is in the **Path** variable in _Environment Variables_. It is usually added automatically.
  
The _Cygwin_ project is known to fail with errors.

For more information, visit [https://gcc.gnu.org/install/binaries.html](https://gcc.gnu.org/install/binaries.html)

#### - nodejs (for npm)

You will need **Node 20**. Download it from [https://nodejs.org/en/download](https://nodejs.org/en/download) and install it. The path is usually automatically added.

#### - go

Visit [https://go.dev/doc/install](https://go.dev/doc/install) to download and install Go. Please use the version specified in the project's `go.mod` file. As of this writing, Meshery uses **1.25.5**. The path is usually automatically added.

#### - Git Bash

Git Bash is a terminal emulator which provides git and Linux-like command line experience. You can download it from: https://git-scm.com/downloads.

### Verify prerequisites

Launch Git Bash, PowerShell or VS Code Terminal and run the following commands to ensure prerequisites are met:
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">make --version</div></div>
</pre>

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">gcc --version</div></div>
</pre>

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">node --version</div></div>
</pre>

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">go version</div></div>
</pre>

If any of these result in command not found error, check if the PATH in environment variables has been set properly. To launch the environment variables window Here's an example:
<a href="/project/contributing/images/win-environment-variables.png">
  <img style= "max-width: 450px;" src="/project/contributing/images/win-environment-variables.png" />
</a>

You are now ready to [contribute](/project/contributing/) to Meshery.


## Using Hyper-V to install a linux VM

Here is a link to install linux VM using Hyper-V: [https://wiki.ubuntu.com/Hyper-V](https://wiki.ubuntu.com/Hyper-V)


## Get the code

- Fork and then clone the [Meshery repository](https://github.com/meshery/meshery)
  ```bash
  $ git clone https://github.com/YOUR-USERNAME/meshery
  ```
