---
layout: page
title: Installation Guide
permalink: installation
---
<a name="getting-started"></a>

# Quick Start 
Getting Meshery up and running on a locally on Docker-enabled system is easy. Use the Meshery command line interface, `mesheryctl`, to start Meshery on any of its [supported platforms](platforms).

## Using `mesheryctl`
`mesheryctl` is a command line interface to manage a Meshery deployment. `mesheryctl` allows you to control Meshery's lifecycle with commands like `start`, `stop`, `status`, `cleanup`. Running `cleanup` will remove all active container instanaces, prune pulled images and remove any local volumes crated by starting Meshery.

### Mac or Linux
Use your choice of homebrew or bash to install `mesheryctl`. You only need to use one.

**Homebrew**
* Install `mesheryctl` and run Meshery on Mac with homebrew:

```
brew tap layer5io/tap
brew install mesheryctl
mesheryctl start
```

**Bash**
* Install `mesheryctl` and run Meshery on Mac or Linux with this script:

```
curl -L https://git.io/meshery | bash -
```

### Windows
* Download and unzip `mesheryctl` from the [Meshery releases](https://github.com/layer5io/meshery/releases/latest) page. Add `mesheryctl` to your PATH for ease of use. Then, execute:

```
./mesheryctl start
```

Upon starting Meshery successfully, instructions to access Meshery will be printed on the sceen.