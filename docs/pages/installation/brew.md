---
layout: default
title: Brew
permalink: installation/linux/brew
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/homebrew.png
---

{% include installation_prerequisites.html %}

In addition to **Bash**, you can also use **Brew** ( a package manager for Mac and Linux, just like apt for Ubuntu ) to install `mesheryctl`.

### Prerequisites

You need to have `Brew` installed on your **Linux** or **Mac** system to perform these actions.

### Install

To install `mesheryctl` using **Brew**, execute the following commands.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">brew install mesheryctl</div></div>
</pre>

You're ready to run Meshery. To do so, execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">mesheryctl system start</div></div>
</pre>

### Upgrade

To upgrade `mesheryctl`, just execute the following command.

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">brew upgrade mesheryctl</div></div>
</pre>

Example output of a successful upgrade

```
➜  ~ brew upgrade mesheryctl
==> Upgrading 1 outdated package:
meshery/tap/mesheryctl 0.3.2 -> 0.3.4
==> Upgrading meshery/tap/mesheryctl
==> Downloading https://github.com/layer5io/meshery/releases/download/v0.3.4/mesheryctl_0.3.4_Darwin_x86_64.zip
==> Downloading from https://github-production-release-asset-2e65be.s3.amazonaws.com/157554479/17522b00-2af0-11ea-8aef-cbfe8
######################################################################## 100.0%
🍺  /usr/local/Cellar/mesheryctl/0.3.4: 5 files, 10.2MB, built in 4 seconds
Removing: /usr/local/Cellar/mesheryctl/0.3.2... (5 files, 10.2MB)
Removing: /Users/lee/Library/Caches/Homebrew/mesheryctl--0.3.2.zip... (3.9MB)
==> Checking for dependents of upgraded formulae...
==> No dependents found!
``` 

Stuck at another error? [Tell us about it](https://slack.meshery.io/)