---
layout: default
title: Using mesheryctl
description: How to use mesheryctl
permalink: guides/mesheryctl
type: Guides
---

## Using `mesheryctl`

`mesheryctl` is the command line interface to manage Meshery and interface with its functionality using a terminal.

`mesheryctl` allows you to control Meshery's lifecycle with commands like `start`, `stop`, `status`, `reset`. Running `reset` will remove all active container instances, prune pulled images and remove any local volumes created by starting Meshery. For more command cases, refer to our exhaustive list of [mesheryctl commands](/docs/guides/mesheryctl-commands).

To upgrade an existing version of mesheryctl installed on your local system, refer to the [mesheryctl upgradation guide](/docs/guides/upgrade)

## Installing `mesheryctl`

### Mac or Linux

Use your choice of homebrew or bash to install `mesheryctl`. You only need to use one.

#### Homebrew

Install `mesheryctl` and run Meshery on Mac with Homebrew.

**Installing with Homebrew**

To install `mesheryctl`, execute the following commands:

```bash
brew tap layer5io/tap
brew install mesheryctl
mesheryctl system start
```

**Upgrading with Homebrew**

To upgrade `mesheryctl`, execute the following command:

```bash
brew upgrade mesheryctl
```

#### Bash

**Installing with Bash**

Install `mesheryctl` and run Meshery on Mac or Linux with this script:

```bash
curl -L https://git.io/meshery | bash -
```

**Upgrading with Bash**

Upgrade `mesheryctl` and run Meshery on Mac or Linux with this script:

```bash
curl -L https://git.io/meshery | bash -
```

### Windows

#### Installing the `mesheryctl` binary

Download and unzip `mesheryctl` from the [Meshery releases](https://github.com/layer5io/meshery/releases/latest) page. Add `mesheryctl` to your PATH for ease of use. Then, execute:

```bash
./mesheryctl system start
```

#### Scoop

Use [Scoop](https://scoop.sh) to install Meshery on your Windows machine.

**Installing with Scoop**

Add the Meshery Scoop Bucket and install:

```bash
scoop bucket add mesheryctl https://github.com/layer5io/scoop-bucket.git
scoop install mesheryctl
```

**Upgrading with Scoop**

To upgrade `mesheryctl`, execute the following command:

```bash
scoop update mesheryctl
```

# Advanced Installation

Users can control the specific container image and tag (version) of Meshery that they would like to run by editing their local `~/.meshery/meshery.yaml` (a docker compose file).
Aligned with the Meshery container image, instead of leaving the implicit `:stable-latest` tag behind image: layer5/meshery, users will instead identify a specific image tag like so:

```bash
version: '3'
services:
  meshery:
    image: layer5/meshery:v0.5.0
    labels:
      - "com.centurylinklabs.watchtower.enable=true"
```
