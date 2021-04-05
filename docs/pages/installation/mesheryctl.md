---
layout: default
title: Install mesheryctl
permalink: installation/mesheryctl
type: installation
display-title: "true"
language: en
list: exclude
# image: /assets/img/platforms/brew.png
---

Meshery's command line client is `mesheryctl`. To install `mesheryctl` on your system, you may choose from any of the following supported methods.

## Bash

**Install** and **Upgrade**

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 curl -L https://git.io/meshery | bash -
 </div></div>
 </pre>

## Homebrew

**Install**

To install `mesheryctl` using homebrew, execute the following commands.

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ brew tap layer5io/tap
 $ brew install mesheryctl
 </div></div>
</pre>

You need to have `brew` installed on your Mac or Linux system to perform these actions.

**Upgrade**

To upgrade `mesheryctl`, execute the following command.

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ brew upgrade mesheryctl
 </div></div>
 </pre>

Example output of a successful upgrade: 
```
➜  ~ brew upgrade mesheryctl
==> Upgrading 1 outdated package:
layer5io/tap/mesheryctl 0.3.2 -> 0.3.4
==> Upgrading layer5io/tap/mesheryctl
==> Downloading https://github.com/layer5io/meshery/releases/download/v0.3.4/mesheryctl_0.3.4_Darwin_x86_64.zip
==> Downloading from https://github-production-release-asset-2e65be.s3.amazonaws.com/157554479/17522b00-2af0-11ea-8aef-cbfe8
######################################################################## 100.0%
🍺  /usr/local/Cellar/mesheryctl/0.3.4: 5 files, 10.2MB, built in 4 seconds
Removing: /usr/local/Cellar/mesheryctl/0.3.2... (5 files, 10.2MB)
Removing: /Users/lee/Library/Caches/Homebrew/mesheryctl--0.3.2.zip... (3.9MB)
==> Checking for dependents of upgraded formulae...
==> No dependents found!
```

## Scoop

`mesheryctl` can be installed via Scoop (a package manager for Windows, just like apt for Ubuntu). To install `mesheryctl` using Scoop, execute the following commands.

**Install**
<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
scoop bucket add mesheryctl https://github.com/layer5io/scoop-bucket.git
scoop install mesheryctl

</div></div>
</pre>

You need to have `scoop` installed on your Windows system to perform these actions.

You're ready to run Meshery. To do so, execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
mesheryctl system start

</div></div>
</pre>

**Upgrade**

To upgrade `mesheryctl`, just execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
scoop update mesheryctl

</div></div>
</pre>

Continue deploying Meshery onto one of the [Supported Platforms]({{ site.baseurl }}/installation/platforms).

