---
layout: default
title: Mesheryctl system commands
description: Get to know what are the meshrey commands present for system model
permalink: guides/mesheryctl/mesheryctl-system
type: Guides
language: en
---


These commands can be executed to do tasks related to system configuration, control & maintainance of meshery itself.


### Mesheryctl system channel :

* `mesheryctl system channel`   -  Switch between release channels.
* `mesheryctl system channel set` - set release channel and version
* `mesheryctl system channel switch` - switch release channel and version
* `mesheryctl system channel view` - view release channel and version
* `mesheryctl system channel view --all` - Show release channel for all contexts

### Mesheryctl system check :

* `mesheryctl system check` -  Meshery environment check.
* `mesheryctl system check --preflight` - Verify environment readiness to deploy Meshery
* `mesheryctl system check --adapters` - Check status of meshery adapters
* `mesheryctl system check --adapter string` - Check status of specified meshery adapter
* `mesheryctl system check --components` - Check status of Meshery components
* `mesheryctl system check --pre` - Verify environment readiness to deploy Meshery
* `mesheryctl system check --operator` - Verify the health of Meshery Operator's deployment with MeshSync and Broker
* `mesheryctl system check --report` - Runs diagnostic checks and bundles up to open an issue if present

### Mesheryctl system config :

* `mesheryctl system config`    -  Configure meshery.
* `mesheryctl system config --token` - Path to token for authenticating to Meshery API (optional, can be done alternatively using "login")
* `mesheryctl system config aks` - Configure Meshery to use AKS cluster
* `mesheryctl system config eks` - Configure Meshery to use EKS cluster
* `mesheryctl system config gke` - Configure Meshery to use GKE cluster
* `mesheryctl system config minikube` - Configure Meshery to use minikube cluster

### Mesheryctl system context :

* `mesheryctl system context`   -  Configure your meshery deployments
* `mesheryctl system context create` - Create a new context
* `mesheryctl system context delete` - Delete context
* `mesheryctl system context list` - List contexts
* `mesheryctl system context switch` - Switch context
* `mesheryctl system context view` - View current context

### Mesheryctl system dashboard :

* `mesheryctl system dashboard` -  Open Meshery-UI in browser.
* `mesheryctl system dashboard --port-forward` - (optional) Use port forwarding to access Meshery UI
* `mesheryctl system dashboard --skip-browser` - (optional) skip opening of MesheryUI in browser.

### Mesheryctl system login :

* `mesheryctl system login`     -  Authenticate to a Meshery server.
* `mesheryctl system login --provider string` - Login Meshery with specified provider

### Mesheryctl system Logout :

* `mesheryctl system logout`    -  Remove authetication for Meshery server.

### Mesheryctl system logs :

* `mesheryctl system logs`      -  Print logs.
* `mesheryctl system logs --follow` - (Optional) Follow the stream of the Meshery's logs. Defaults to false.

### Mesheryctl system provider :

* `mesheryctl system provider`  -  Switch between providers.
* `mesheryctl system provider set` - Set provider
* `mesheryctl system provider list` - List available providers
* `mesheryctl system provider reset` - Reset provider to default
* `mesheryctl system provider switch` - switch provider and redeploy
* `mesheryctl system provider view` - view provider

### Mesheryctl system reset :

* `mesheryctl system reset`     -  Reset Meshery's configuration.

### Mesheryctl system restart :

* `mesheryctl system restart`   -  Stop, then start Meshery.
* `mesheryctl system restart --skip-update` - (optional) skip checking for new Meshery's container images.

### Mesheryctl system start :

* `mesheryctl system start`     -  Start Meshery.
* `mesheryctl system start --skip-update` - (optional) skip checking for new Meshery's container images.
* `mesheryctl system start --reset` - (optional) reset Meshery's configuration file to default settings.
* `mesheryctl system start --yes` - Silently create Meshery's configuration file with default settings
* `mesheryctl system start --platform string` - platform to deploy Meshery to.
* `mesheryctl system start --skip-browser` - (optional) skip opening of MesheryUI in browser.

### Mesheryctl system status :

* `mesheryctl system status`    -  Check Meshery status.

### Mesheryctl system stop :

* `mesheryctl system stop`      -  Stop Meshery.
* `mesheryctl system stop --force` - (optional) uninstall Meshery resources forcefully
* `mesheryctl system stop --reset` - (optional) reset Meshery's configuration file to default settings.
* `mesheryctl system stop --keep-namespace` - (optional) keep the Meshery namespace during uninstallation

### Mesheryctl system token :

* `mesheryctl system token`     -  Manage Meshery user tokens.
* `mesheryctl system token create` - Create a token in your meshconfig
* `mesheryctl system token delete` - Delete a token in your meshconfig
* `mesheryctl system token list` - List tokens
* `mesheryctl system token set` - Set token for context
* `mesheryctl system token view` - View token

### Mesheryctl system update :

* `mesheryctl system update`    -  Pull new Meshery images/manifest files.
* `mesheryctl system update --skip-reset` - (optional) skip checking for new Meshery manifest files.
<br><br>

To know more about mesheryctl commands type `mesheryctl system [command] -h` or `mesheryctl system context [command] -h` in your terminal (after installing mesheryctl).