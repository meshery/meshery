---
layout: default
title: Upgrade `mesheryctl` and Meshery
description: How to Upgrade mesheryctl
permalink: guides/upgrade
type: Guides
---

Various components of Meshery will need to be upgraded as new releases become available. Meshery is comprised of a number of components including a server, adapters, UI, and CLI.

# Upgrade Meshery Server, Adapters, and UI
All three of these components are released as part of the same set of artifacts. In order to upgrade Meshery server, UI and adapters, you may execute the following command:

### Meshery Docker Deployments
In order to upgrade Meshery Server, Adapters, and UI, execute the following command:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system upgrade
 </div></div>
 </pre>
### Meshery Kubernetes Deployments

Use `kubectl apply` or `helm` to upgrade the Meshery application manifests in your Kubernetes cluster.

## Upgrading `mesheryctl`

The Meshery command line client is available in different package managers. Use the instructions relevant to your environment.

### Upgrading `mesheryctl` using Homebrew

<p>To upgrade `mesheryctl`, execute the following command:</p>

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew upgrade mesheryctl
 </div></div>
 </pre>

### Upgrading `mesheryctl` using Bash

Upgrade `mesheryctl` and run Meshery on Mac or Linux with this script:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 curl -L https://git.io/meshery | bash -
 </div></div>
 </pre>

### Upgrading `mesheryctl` using Scoop

To upgrade `mesheryctl`, execute the following command:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 scoop update mesheryctl
 </div></div>
 </pre>
