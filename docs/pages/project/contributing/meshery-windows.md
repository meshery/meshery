---
layout: page
title: Setting up Meshery Development Environment on Windows
permalink: project/contributing/meshery-windows
description: How to set up Meshery Development Environment on Windows
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

## Using Hyper-V to install a linux VM

Here is a link to install linux VM using Hyper-V: [https://wiki.ubuntu.com/Hyper-V](https://wiki.ubuntu.com/Hyper-V)

### Git Bash  
Git Bash is a terminal emulator which provides git command line experience. This will make working with git easier. You can download it from here: [https://git-scm.com/downloads](https://git-scm.com/downloads)

### Get the code

- Fork and then clone the [Meshery repository](https://github.com/layer5io/meshery)
  ```bash
  $ git clone https://github.com/YOUR-USERNAME/meshery
  ```