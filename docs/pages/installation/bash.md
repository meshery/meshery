---
layout: default
title: Bash
permalink: installation/linux/bash
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/bash.png
---

{% include installation_prerequisites.html %}

`mesheryctl` can be installed via **Bash** ( a command-line interpreter and scripting language that is used on Unix-like operating systems, including Linux and macOS ). `mesheryctl` is also available through **Homebrew**.

### Prerequisites

You need to have `Bash` ( which is commonly available on both Linux and macOS ) , to perform these actions.

### Install

To install `mesheryctl` using `Bash`, execute the following commands.

**Install** and **Upgrade**

Install `mesheryctl` command

 <pre class="codeblock-pre">
 <div class="codeblock">
 <div class="clipboardjs">
  $ curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -
 </div></div>
 </pre>

Install `mesheryctl` command and deploy Meshery on Docker

 <pre class="codeblock-pre">
 <div class="codeblock">
 <div class="clipboardjs">
  $ curl -L https://meshery.io/install | PLATFORM=docker bash -
 </div></div>
 </pre>

Install `mesheryctl` command and deploy Meshery on Kubernetes

 <pre class="codeblock-pre">
 <div class="codeblock">
 <div class="clipboardjs">
  $ curl -L https://meshery.io/install | PLATFORM=kubernetes bash -
 </div></div>
 </pre>

Install `mesheryctl` command and choose an [adapter]({{ site.baseurl }}/concepts/architecture/adapters) to be loaded.

 <pre class="codeblock-pre">
 <div class="codeblock">
 <div class="clipboardjs">
  $ curl -L https://meshery.io/install | ADAPTERS=consul PLATFORM=kubernetes bash -
 </div></div>
 </pre>

You are ready to run **Meshery** via **mesheryctl**. To do so, execute the following command.

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">mesheryctl system start</div></div>
 </pre>

```
Note -> To access Meshery Dashboard, execute the following command.
```

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">mesheryctl system dashboard</div></div>
 </pre>

Stuck at another error? [Tell us about it](https://slack.meshery.io/)