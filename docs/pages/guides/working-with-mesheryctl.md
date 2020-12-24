---
layout: default
title: Using mesheryctl
description: How to use mesheryctl
permalink: guides/mesheryctl 
type: Guides
---

`mesheryctl` is the command line interface to manage Meshery and interface with its functionality using a terminal. `mesheryctl` commands are categorized into three main areas:

- Lifecycle management of Meshery (control Meshery's lifecycle with commands like `system start`, `stop`, `status`, `reset`. )
- Lifecycle management of Service Meshes
- Performance management of Service Meshes and Workloads

<!-- Running `reset` will remove all active container instances, prune pulled images and remove any local volumes created by starting Meshery. -->

## Related Guides

- For an exhaustive list of commands and syntax, refer to the **[`mesheryctl` Command Reference]({{ site.baseurl }}/guides/mesheryctl-commands)**.
- To upgrade `mesheryctl`, refer to the **[Upgrade Guide]({{ site.baseurl }}/guides/upgrade)**.

## Installing `mesheryctl`

### Mac or Linux

Use your choice of homebrew or bash to install `mesheryctl`. You only need to use one.
### Homebrew

Install `mesheryctl` and run Meshery on Mac with Homebrew.

#### Installing with Homebrew

To install `mesheryctl`, execute the following commands:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew tap layer5io/tap
 brew install mesheryctl
 mesheryctl system start
 </div></div>
 </pre>

**Upgrading with Homebrew**

To upgrade `mesheryctl`, execute the following command:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew upgrade mesheryctl
 </div></div>
 </pre>

#### Bash

**Installing with Bash**

Install `mesheryctl` and run Meshery on Mac or Linux with this script:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 curl -L https://git.io/meshery | bash -
 </div></div>
 </pre>

**Upgrading with Bash**

Upgrade `mesheryctl` and run Meshery on Mac or Linux with this script:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 curl -L https://git.io/meshery | bash -
 </div></div>
 </pre>

## Windows

### Installing the `mesheryctl` binary

Download and unzip `mesheryctl` from the [Meshery releases](https://github.com/layer5io/meshery/releases/latest) page. Add `mesheryctl` to your PATH for ease of use. Then, execute:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 ./mesheryctl system start
 </div></div>
 </pre>

### Scoop

Use [Scoop](https://scoop.sh) to install Meshery on your Windows machine.

**Installing with Scoop**

Add the Meshery Scoop Bucket and install:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 scoop bucket add mesheryctl https://github.com/layer5io/scoop-bucket.git
 scoop install mesheryctl
 </div></div>
 </pre>

**Upgrading with Scoop**

To upgrade `mesheryctl`, execute the following command:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 scoop update mesheryctl
 </div></div>
 </pre>

## Advanced Installation

Users can control the specific container image and tag (version) of Meshery that they would like to run by editing their local *~/.meshery/meshery.yaml* (a docker compose file).
Aligned with the Meshery container image, instead of leaving the implicit :stable-latest tag behind image: layer5/meshery, users will instead identify a specific image tag like so:

```
bash
version: '3'
services:
  meshery:
    image: layer5/meshery:v0.5.0
    labels:
      - "com.centurylinklabs.watchtower.enable=true"
```

# Configuring Autocompletion for `mesheryctl`

If you would like to have `mesheryctl` commands automatically completed for use as you use `mesheryctl`, then use the following instructions to configure automatic completion within your environment.

## Autocompletion for Bash

### bash <= 3.2

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 source /dev/stdin <<< "$(mesheryctl system completion bash)"
 </div></div>
 </pre>

### bash >= 4.0

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 source <(mesheryctl system completion bash)
 </div></div>
 </pre>

### bash <= 3.2 on MacOS

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew install bash-completion # ensure you have bash-completion 1.3+
 mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
 </div></div>
 </pre>

### bash >= 4.0 on MacOS

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew install bash-completion@2
 mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
 </div></div>
 </pre>

## Autocompletion for zsh

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 source <(mesheryctl system completion zsh)
 </div></div>
 </pre><br>

If shell completion is not already enabled in your environment you will need to enable it.  You can execute the following once:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 ~/.zshrc > echo "autoload -U compinit; compinit" 
 </div></div>
 </pre>
_Note_ : You might need to restart your shell for this setup to take effect.

#### zsh on MacOS and Oh My zsh

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system completion zsh > "${fpath[1]}/_mesheryctl"
 </div></div>
 </pre>

### Autocompletion for fish

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system completion fish | source
 </div></div>
 </pre><br>

To load fish shell completions for each session, execute once:
 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system completion fish > ~/.config/fish/completions/mesheryctl.fish
 </div></div>
 </pre>

# Suggested Reading

- For an exhaustive list of commands and syntax, refer to the **[`mesheryctl` Command Reference]({{ site.baseurl }}/guides/mesheryctl-commands)**.
- To upgrade `mesheryctl`, refer to the **[Upgrade Guide]({{ site.baseurl }}/guides/upgrade)**.
