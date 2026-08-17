---
title: Mesheryctl system commands
categories: [mesheryctl]
description: Mesheryctl system commands for managing Meshery deployments.
---

Let's get familiar with mesheryctl system commands. The syntax of the mesheryctl commands goes as follows : `mesheryctl <Main_command> <Argument> <Flags>`

## Main_command : system
### start 
`mesheryctl system start` : This will initiate Meshery & automatically open it in your default web browser.

![start](images/start.png)

`mesheryctl system start --skip-browser` : It skips opening Meshery in your browser with the provided URL.

![skip browser](images/skipbrowser.png)

`mesheryctl system start --skip-update` : It is used when you want to skip updating Meshery if an update is available.

![system update](images/system update.png)

`mesheryctl system start --reset` : It resets your Meshery configuration file to its default configuration.

`mesheryctl system start --platform string` : It allows you specify a platform for deploying Meshery.

![platform](images/platform.png)


### stop 
`mesheryctl system stop` : It stops Meshery resources & delete its associated namespaces.

![stop](images/stop.png)

`mesheryctl system stop --reset` : It stops Meshery and resets the Meshery configuration file to its default configuration.

![stop reset](images/stop reset.png)

`mesheryctl system stop --keep-namespace` : It stops Meshery without deleting the associated namespaces.

![keep namespace](images/keep namespace.png)

`mesheryctl system stop --force` : Force stops Meshery instead of gentle way. This is only used in emergency situations when `mesheryctl system stop` can't halt Meshery.

![force stop](images/force stop.png)

### update
`mesheryctl system update` : This updates Meshery itself, not the mesheryctl. Ensure Meshery is running when using this.

![system update](images/system update.png)

`mesheryctl system update --skip-reset` : Skips the check for a new manifest file.

![update skip reset](images/update skip reset.png)

### reset
`mesheryctl system reset` : Resets Meshery to its default configuration.

![reset](images/reset.png)

### restart 
`meshryctl system restart` : Stops Meshery and then starts it again. Opens the website in your default browser.

![restart](images/restart.png)

### status 
`mesheryctl system status` : Displays the status of Meshery components.

`mesheryctl system status --verbose` : Provides additional data along with Meshery and its component status.

![system status](images/system status.png)


### dashboard
`mesheryctl system dashboard` : Opens the Meshery dashboard in your default browser.

![system dashboard](images/system dahboard.png)

`mesheryctl system dashboard --skip-browser` : Provides the link to the dashboard, allowing you to open it in any browser.

![dashboard skip](images/dashboard skip.png)

`mesheryctl system dashboard --port-forward` : If the current port is busy, it opens the dashboard on another port.

![portforward](images/portforward.png)


### login 
`mesheryctl system login` : Authenticates you with your selected provider.

![system login](images/system login.png)

### check
`mesheryctl system check` : Performs checks for both pre & post mesh deployment scenarios on Meshery.

![system check](images/system check.png)

`mesheryctl system check --preflight` : Runs pre-deployment checks.

`mesheryctl system check --adapter` : Runs checks for a specific Mesh adapter.

`mesheryctl system check --adapters` : Runs checks for Meshery adapters

`mesheryctl system check --components` : Runs checks for Meshery components

`mesheryctl system check --operator` : Runs checks for Meshery Operator

## Main_command : system channel
### channel
`mesheryctl system channel set [stable|stable-version|edge|edge-version]` : Used to set the channel.

`mesheryctl system channel switch [stable|stable-version|edge|edge-version]` : Used to switch between channels.

![channel set](images/channel set.png)

`mesheryctl system channel view --all` : Displays all available channels.

`mesheryctl system channel view` : Displays the current channel.

![channel view](images/channel view.png)


## Main_command : system context
### create 
`mesheryctl system context create 'context-name'` : Creates a new context with default parameters.

![context create](images/context create.png)

`mesheryctl system context create --component stringArray` : Specifies the component to be created in the context.

`mesheryctl system context create --platform string` : Specifies the platform.

`mesheryctl system context create --set` : Sets this  context as default context.

`mesheryctl system context create --url string` : Specifies the target URL.

![context flag](images/context flag.png)


###  switch
`mesheryctl system context switch` : Easily switch between different contexts.

###  list
`mesheryctl system context list` : Lists all your available Meshery contexts.

![context list](images/context list.png)

###  delete
`mesheryctl system context delete` : Delete context.

![context delete](images/context delete.png)


###  view
`mesheryctl system context view` : Display all your contexts with additional information.

![context view](images/context view.png)


## Main_command : system provider
### switch
`mesheryctl system provider switch` : Changes your provider

![provider switch](images/pro switch.png)

### list
`mesheryctl system provider list` : Lists all available providers

![provider list](images/pro list.png)

### set
`mesheryctl system provider set` : Set your provider

![provider set](images/pro set.png)

### view
`mesheryctl system provider view` : Lists your current context and provider

![provider view](images/pro view.png)

