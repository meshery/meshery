---
layout: default
title: Using mesheryctl
description: How to use mesheryctl
permalink: guides/mesheryctl 
type: Guides
---

# Using `mesheryctl`

`mesheryctl` is the command line interface to manage Meshery and interface with its functionality using a terminal. `mesheryctl` commands are categorized into three main areas:

- Lifecycle management of Meshery (control Meshery's lifecycle with commands like `system start`, `stop`, `status`, `reset`. )
- Lifecycle management of Service Meshes
- Performance management of Service Meshes and Workloads

<!-- Running `reset` will remove all active container instances, prune pulled images and remove any local volumes created by starting Meshery. -->
## Related Guides

- For an exhaustive list of commands and syntax, refer to the **[`mesheryctl` Command Reference](/docs/guides/mesheryctl-commands)**.
- To upgrade `mesheryctl`, refer to the **[Upgrade Guide](/docs/guides/upgrade)**.

# Installing `mesheryctl`

## Mac or Linux

Use your choice of homebrew or bash to install *mesheryctl*. You only need to use one.

### Homebrew

Install *mesheryctl* and run Meshery on Mac with Homebrew.

**Installing with Homebrew**

To install *mesheryctl*, execute the following commands:

```bash
brew install layer5io mesheryctl system start
```

**Upgrading with Homebrew**

To upgrade *mesheryctl*, execute the following command:

```bash
brew upgrade mesheryctl
```

#### Bash

**Installing with Bash**

Install *mesheryctl* and run Meshery on Mac or Linux with this script:

```bash
curl -L https://git.io/meshery | bash -
```

**Upgrading with Bash**

Upgrade *mesheryctl* and run Meshery on Mac or Linux with this script:

```bash
curl -L https://git.io/meshery | bash -
```

## Windows

### Installing the `mesheryctl` binary

Download and unzip *mesheryctl* from the [Meshery releases](https://github.com/layer5io/meshery/releases/latest) page. Add *mesheryctl* to your PATH for ease of use. Then, execute:

```bash
.*mesheryctl*system start
```

### Scoop

Use [Scoop](https://scoop.sh) to install Meshery on your Windows machine.

**Installing with Scoop**

Add the Meshery Scoop Bucket and install:

```bash
scoop bucket add*mesheryctl*https://github.com/layer5io/scoop-bucket.git
scoop install*mesheryctl*
```

**Upgrading with Scoop**

To upgrade *mesheryctl*, execute the following command:

```bash
scoop update mesheryctl
```

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

# Configuring Auto Completion for `mesheryctl`
If you would like to have `mesheryctl` commands automatically completed for use as you use `mesheryctl`, then use the following instructions to configure automatic completion within your environment.

## Autocompletion for Bash

### bash <= 3.2

```bash
source /dev/stdin <<< "$(mesheryctl system completion bash)"
```

### bash >= 4.0

```bash
source <(mesheryctl system completion bash)
```

### bash <= 3.2 on osx

```bash
brew install bash-completion # ensure you have bash-completion 1.3+
mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
```

### bash >= 4.0 on osx

```bash
brew install bash-completion@2
mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
```

## Autocompletion for zsh

```bash
source <(mesheryctl system completion zsh)
```

If shell completion is not already enabled in your environment you will need to enable it.  You can execute the following once:

```bash
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

Might need to start a new shell for this setup to take effect.

### zsh on osx / oh-my-zsh

```bash
mesheryctl system completion zsh > "${fpath[1]}/_mesheryctl"
```

## Autocompletion for fish

```bash
mesheryctl system completion fish | source
```

To load fish shell completions for each session, execute once:

```bash
mesheryctl system completion fish > ~/.config/fish/completions/mesheryctl.fish
```
