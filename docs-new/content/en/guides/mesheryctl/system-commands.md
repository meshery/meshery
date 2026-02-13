---
layout: default
title: Mesheryctl system commands
categories: [mesheryctl]
description: Mesheryctl system commands for managing Meshery deployments.
---

Let's get familiar with mesheryctl system commands. The syntax of the mesheryctl commands goes as follws : `mesheryctl <Main_command> <Argument> <Flags>`

## Main_command : system
### start 
`mesheryctl system start` : This will initiate Meshery & automatically open it in your default web browser.

<a href="/guides/mesheryctl/images/start.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/start.png" /></a>

`mesheryctl system start --skip-browser` : It skips opening Meshery in your browser with the provided URL.

<a href="/guides/mesheryctl/images/skipbrowser.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/skipbrowser.png" /></a>

`mesheryctl system start --skip-update` : It is used when you want to skip updating Meshery if an update is available.

<a href="/guides/mesheryctl/images/system update.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/system update.png" /></a>

`mesheryctl system start --reset` : It resets your Meshery configuration file to its default configuration.

`mesheryctl system start --platform string` : It allows you specify a platform for deploying Meshery.

<a href="/guides/mesheryctl/images/platform.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/platform.png" /></a>


### stop 
`mesheryctl system stop` : It stops Meshery resources & delete its associated namespaces.

<a href="/guides/mesheryctl/images/stop.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/stop.png" /></a>

`mesheryctl system stop --reset` : It stops Meshery and resets the Meshery configuration file to its default configuration.

<a href="/guides/mesheryctl/images/stop reset.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/stop reset.png" /></a>

`mesheryctl system stop --keep-namespace` : It stops Meshery without deleting the associated namespaces.

<a href="/guides/mesheryctl/images/keep namespace.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/keep namespace.png" /></a>

`mesheryctl system stop --force` : Force stops Meshery instead of gentle way. This is only used in emergency situations when `mesheryctl system stop` can't halt Meshery.

<a href="/guides/mesheryctl/images/force stop.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/force stop.png" /></a>

### update
`mesheryctl system update` : This updates Meshery itself, not the mesheryctl. Ensure Meshery is running when using this.

<a href="/guides/mesheryctl/images/system update.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/system update.png" /></a>

`mesheryctl system update --skip-reset` : Skips the check for a new manifest file.

<a href="/guides/mesheryctl/images/update skip reset.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/update skip reset.png" /></a>

### reset
`mesheryctl system reset` : Resets Meshery to its default configuration.

<a href="/guides/mesheryctl/images/reset.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/reset.png" /></a>

### restart 
`meshryctl system restart` : Stops Meshery and then starts it again. Opens the website in your default browser.

<a href="/guides/mesheryctl/images/restart.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/restart.png" /></a>

### status 
`mesheryctl system status` : Displays the status of Meshery components.

`mesheryctl system status --verbose` : Provides additional data along with Meshery and its component status.

<a href="/guides/mesheryctl/images/system status.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/system status.png" /></a>


### dashboard
`mesheryctl system dashboard` : Opens the Meshery dashboard in your default browser.

<a href="/guides/mesheryctl/images/system dahboard.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/system dahboard.png" /></a>

`mesheryctl system dashboard --skip-browser` : Provides the link to the dashboard, allowing you to open it in any browser.

<a href="/guides/mesheryctl/images/dashboard skip.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/dashboard skip.png" /></a>

`mesheryctl system dashboard --port-forward` : If the current port is busy, it opens the dashboard on another port.

<a href="/guides/mesheryctl/images/portforward.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/portforward.png" /></a>


### login 
`mesheryctl system login` : Authenticates you with your selected provider.

<a href="/guides/mesheryctl/images/system login.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/system login.png" /></a>

### check
`mesheryctl system check` : Performs checks for both pre & post mesh deployment scenarios on Meshery.

<a href="/guides/mesheryctl/images/system check.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/system check.png" /></a>

`mesheryctl system check --preflight` : Runs pre-deployment checks.

`mesheryctl system check --adapter` : Runs checks for a specific Mesh adapter.

`mesheryctl system check --adapters` : Runs checks for Meshery adapters

`mesheryctl system check --components` : Runs checks for Meshery components

`mesheryctl system check --operator` : Runs checks for Meshery Operator

## Main_command : system channel
### channel
`mesheryctl system channel set [stable|stable-version|edge|edge-version]` : Used to set the channel.

`mesheryctl system channel switch [stable|stable-version|edge|edge-version]` : Used to switch between channels.

<a href="/guides/mesheryctl/images/channel set.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/channel set.png" /></a>

`mesheryctl system channel view --all` : Displays all available channels.

`mesheryctl system channel view` : Displays the current channel.

<a href="/guides/mesheryctl/images/channel view.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/channel view.png" /></a>


## Main_command : system context
### create 
`mesheryctl system context create 'context-name'` : Creates a new context with default parameters.

<a href="/guides/mesheryctl/images/context create.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/context create.png" /></a>

`mesheryctl system context create --component stringArray` : Specifies the component to be created in the context.

`mesheryctl system context create --platform string` : Specifies the platform.

`mesheryctl system context create --set` : Sets this  context as default context.

`mesheryctl system context create --url string` : Specifies the target URL.

<a href="/guides/mesheryctl/images/context flag.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/context flag.png" /></a>


###  switch
`mesheryctl system context switch` : Easily switch between different contexts.

###  list
`mesheryctl system context list` : Lists all your available Meshery contexts.

<a href="/guides/mesheryctl/images/context list.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/context list.png" /></a>

###  delete
`mesheryctl system context delete` : Delete context.

<a href="/guides/mesheryctl/images/context delete.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/context delete.png" /></a>


###  view
`mesheryctl system context view` : Display all your contexts with additional information.

<a href="/guides/mesheryctl/images/context view.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/context view.png" /></a>


## Main_command : system provider
### switch
`mesheryctl system provider switch` : Changes your provider

<a href="/guides/mesheryctl/images/pro switch.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/pro switch.png" /></a>

### list
`mesheryctl system provider list` : Lists all available providers

<a href="/guides/mesheryctl/images/pro list.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/pro list.png" /></a>

### set
`mesheryctl system provider set` : Set your provider

<a href="/guides/mesheryctl/images/pro set.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/pro set.png" /></a>

### view
`mesheryctl system provider view` : Lists your current context and provider

<a href="/guides/mesheryctl/images/pro view.png"><img alt="skip-browser" style="width:500px;height:auto;" src="/guides/mesheryctl/images/pro view.png" /></a>

