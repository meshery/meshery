---
layout: page
title: Installation Guide
permalink: /installation
---

<a name="getting-started"></a>

# Quick Start

Getting Meshery up and running on a locally on Docker-enabled system is easy. Use the Meshery command line interface, *mesheryctl*, to start Meshery on any of its [supported platforms](/docs/installation/platforms).

## Using *mesheryctl*

*mesheryctl* is a command line interface to manage a Meshery deployment. *mesheryctl* allows you to control Meshery's lifecycle with commands like **start**, **stop**, **status**, **reset**. Running **reset** will remove all active container instances, prune pulled images and remove any local volumes created by starting Meshery.

### Mac or Linux

Use your choice of homebrew or bash to install *mesheryctl*. You only need to use one.

#### Homebrew

Install *mesheryctl* and run Meshery on Mac with Homebrew.

**Installing with Homebrew**

To install *mesheryctl*, execute the following commands:

<pre><code>
brew tap layer5io/tap
brew install mesheryctl
mesheryctl system start
</code></pre>

**Upgrading with Homebrew**

To upgrade *mesheryctl*, execute the following command:

```brew upgrade mesheryctl```

#### Bash

**Installing with Bash**

Install *mesheryctl* and run Meshery on Mac or Linux with this script:

`curl -L https://git.io/meshery | bash -`

**Upgrading with Bash**

Upgrade *mesheryctl* and run Meshery on Mac or Linux with this script:

`curl -L https://git.io/meshery | bash -`

### Windows

**Installing the *mesheryctl* binary**

Download and unzip *mesheryctl* from the [Meshery releases](https://github.com/layer5io/meshery/releases/latest) page. Add *mesheryctl* to your PATH for ease of use. Then, execute:

`./mesheryctl system start`

#### Scoop

Use [Scoop](https://scoop.sh) to install Meshery on your Windows machine.

**Installing with Scoop**

Add the Meshery Scoop Bucket and install:

<pre><code>
scoop bucket add mesheryctl https://github.com/layer5io/scoop-bucket.git
scoop install mesheryctl
</code></pre>

**Upgrading with Scoop**

To upgrade *mesheryctl*, execute the following command:

`scoop update mesheryct`

# Advanced Installation

Users can control the specific container image and tag (version) of Meshery that they would like to run by editing their local *~/.meshery/meshery.yaml* (a docker compose file).
Aligned with the Meshery container image, instead of leaving the implicit *:stable-latest* tag behind image: layer5/meshery, users will instead identify a specific image tag like so:

<pre><code>
version: '3'
services:
  meshery:
    image: layer5/meshery:v0.5.0
    labels:
      - "com.centurylinklabs.watchtower.enable=true"
</code></pre>

When Meshery is up and running, instructions to access Meshery will be printed on the screen and your default browser should be directed to the Meshery login screen.
