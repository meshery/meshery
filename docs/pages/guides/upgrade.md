---
layout: default
title: Upgrade*mesheryctl*and Meshery
description: How to Upgrade mesheryctl
permalink: guides/upgrade
type: Guides
---
Various components of Meshery will need to be upgraded as new releases become available. Meshery is comprised of a number of components including a server, adapters, UI, and CLI.

## Upgrade Meshery Server, Adapters, and UI
All three of these components are released as part of the same set of artifacts. In order to upgrade Meshery server, UI and adapters, you may execute the following command:

```*mesheryctl*system upgrade
```

## Upgrade Meshery Client *mesheryctl*

### Upgrading*mesheryctl*using Homebrew

To upgrade *mesheryctl*, execute the following command:

```bash
brew upgrade*mesheryctl*```

Example output of a successful upgrade:

```bash
➜  ~ brew upgrade*mesheryctl*==> Upgrading 1 outdated package:
layer5io/tap*mesheryctl*0.3.2 -> 0.3.4
==> Upgrading layer5io/tap*mesheryctl*==> Downloading https://github.com/layer5io/meshery/releases/download/v0.3.4/mesheryctl_0.3.4_Darwin_x86_64.zip
==> Downloading from https://github-production-release-asset-2e65be.s3.amazonaws.com/157554479/17522b00-2af0-11ea-8aef-cbfe8
######################################################################## 100.0%
🍺  /usr/local/Cellar*mesheryctl*0.3.4: 5 files, 10.2MB, built in 4 seconds
Removing: /usr/local/Cellar*mesheryctl*0.3.2... (5 files, 10.2MB)
Removing: /Users/lee/Library/Caches/Homebrew*mesheryctl*-0.3.2.zip... (3.9MB)
==> Checking for dependents of upgraded formulae...
==> No dependents found!
```


### Upgrading*mesheryctl*using Bash

Upgrade *mesheryctl* and run Meshery on Mac or Linux with this script:

```bash
curl -L https://git.io/meshery | bash -
```

### Upgrading*mesheryctl*using Scoop

To upgrade *mesheryctl*, execute the following command:

```bash
scoop update *mesheryctl*
```
