---
layout : default
title : Mesheryctl system commands
permalink : guides/mesheryctl/system-commands
language : en
type : Guides
category : mesheryctl
---

Let's get familiar with mesheryctl system commands. The syntax of commands goes like this : `mesheryctl <Main_command> <Argument> <Flags>`

## Main_command : system
### start 
`mesheryctl system start` : it will initiate the Meshery & automatically open it in your default web browser.

<a href="{{ site.baseurl }}/assets/img/syscmd/start.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/start.png" /></a>

`mesheryctl system start --skip-browser` : If you don't want to open the URL on your default browser you can use this command. It will provide the URL, so that you can open the URL on any browser.

<a href="{{ site.baseurl }}/assets/img/syscmd/skipbrowser.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/skipbrowser.png" /></a>

`mesheryctl system start --skip-update` : It is used when you want to skip updating Meshery if an update is available.

<a href="{{ site.baseurl }}/assets/img/syscmd/context view.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/context view.png" /></a>

`mesheryctl system start --reset` : It reset your Meshery configuration file to its default configuration.

<a href="{{ site.baseurl }}/assets/img/syscmd/skipbrowser.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/skipbrowser.png" /></a>

`mesheryctl system start --platform string` : It allows you to specify the platform for running Meshery.

<a href="{{ site.baseurl }}/assets/img/syscmd/platform.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/platform.png" /></a>


### stop 
`mesheryctl system stop` : it stops Meshery resources & delete the associated namespaces.

<a href="{{ site.baseurl }}/assets/img/syscmd/stop.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/stop.png" /></a>

`mesheryctl system stop --reset` : It stops Meshery and resets the Meshery configuration file to its default configuration.

<a href="{{ site.baseurl }}/assets/img/syscmd/stop reset.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/stop reset.png" /></a>

`mesheryctl system stop --keep-namespace` : It stops Meshery without deleting the associated namespaces.

<a href="{{ site.baseurl }}/assets/img/syscmd/keep namespace.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/keep namespace.png" /></a>

`mesheryctl system stop --force` : Use this in emergency situations when `mesheryctl system stop` can't halt Meshery.

<a href="{{ site.baseurl }}/assets/img/syscmd/force stop.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/force stop.png" /></a>

### update
`mesheryctl system update` : This updates Meshery itself, not the mesheryctl. Ensure Meshery is running when using this.

<a href="{{ site.baseurl }}/assets/img/syscmd/system update.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/system update.png" /></a>

`mesheryctl system update --skip-reset` : Skips the check for a new manifest file.

<a href="{{ site.baseurl }}/assets/img/syscmd/update skip reset.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/update skip reset.png" /></a>

### reset
`mesheryctl system reset` : Resets Meshery to its default configuration.

<a href="{{ site.baseurl }}/assets/img/syscmd/reset.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/reset.png" /></a>

### restart 
`meshryctl system restart` : Stops Meshery and then starts it again. Opens the website in your default browser.

<a href="{{ site.baseurl }}/assets/img/syscmd/restart.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/restart.png" /></a>

### status 
`mesheryctl system status` : Displays the status of Meshery components.

`mesheryctl system status --verbose` : Provides additional data along with Meshery and its component status.

<a href="{{ site.baseurl }}/assets/img/syscmd/system status.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/system status.png" /></a>


### dashboard
`mesheryctl system dashboard` : Opens the Meshery dashboard in your default browser.

<a href="{{ site.baseurl }}/assets/img/syscmd/system dahboard.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/system dahboard.png" /></a>

`mesheryctl system dashboard --skip-browser` : Provides the link to the dashboard, allowing you to open it in any browser.

<a href="{{ site.baseurl }}/assets/img/syscmd/dashboard skip.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/dashboard skip.png" /></a>

`mesheryctl system dashboard --port-forward` : If the current port is busy, it opens the dashboard on another port.

<a href="{{ site.baseurl }}/assets/img/syscmd/portforward.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/portforward.png" /></a>


### login 
`mesheryctl system login` : Authenticates you with your selected provider.

<a href="{{ site.baseurl }}/assets/img/syscmd/system login.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/system login.png" /></a>

### check
`mesheryctl system check` : Performs checks for mesh deployment.

<a href="{{ site.baseurl }}/assets/img/syscmd/system check.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/system check.png" /></a>

`mesheryctl system check --preflight` : Runs pre-deployment checks.

`mesheryctl system check --adapter` : Runs checks for specific adapter

`mesheryctl system check --adapters` : Runs checks for adapters

`mesheryctl system check --components` : Runs checks for components

`mesheryctl system check --operator` : Runs checks for Meshery Operator

## Main_command : system channel
### channel
`mesheyctl system channel --set [stable|stable-version|edge|edge-version]` : Used to set the channel.

`mesheyctl system channel --switch [stable|stable-version|edge|edge-version]` : Used to switch between channels.

<a href="{{ site.baseurl }}/assets/img/syscmd/channel set.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/channel set.png" /></a>

`mesheyctl system channel view --all` : Displays all available channels.

`mesheyctl system channel view` : Displays the current channel.

<a href="{{ site.baseurl }}/assets/img/syscmd/channel view.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/channel view.png" /></a>


## Main_command : system context
### create 
`mesheryctl system context create 'context-name'` : Creates a new context with default parameters.

<a href="{{ site.baseurl }}/assets/img/syscmd/context create.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/context create.png" /></a>

`mesheryctl system context create --component stringArray` : specifies the component.

`mesheryctl system context create --platform string` : Specifies the platform.

`mesheryctl system context create --set` : Sets this  context as default context.

`mesheryctl system context create --url string` : Specifies the URL

<a href="{{ site.baseurl }}/assets/img/syscmd/context flag.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/context flag.png" /></a>


###  switch
`mesheryctl system context switch` :  Switches between different contexts.

<a href="{{ site.baseurl }}/assets/img/syscmd/context flag.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/context flag.png" /></a>

###  list
`mesheryctl system context list` : lists all your contexts.

<a href="{{ site.baseurl }}/assets/img/syscmd/context list.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/context list.png" /></a>

###  delete
`mesheryctl system context delete` : Delete context.

<a href="{{ site.baseurl }}/assets/img/syscmd/context delete.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/context flag.png" /></a>


###  view
`mesheryctl system context view` : Display all your contexts with additional information.

<a href="{{ site.baseurl }}/assets/img/syscmd/context view.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/context view.png" /></a>


## Main_command : system provider
### switch
`mesheryctl system provider switch` : Changes your provider

<a href="{{ site.baseurl }}/assets/img/syscmd/pro switch.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/pro switch.png" /></a>

### list
`mesheryctl system provider list` : Lists all available providers

<a href="{{ site.baseurl }}/assets/img/syscmd/pro list.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/pro list.png" /></a>

### set
`mesheryctl system provider set` : Set your provider

<a href="{{ site.baseurl }}/assets/img/syscmd/pro set.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/pro set.png" /></a>

### view
`mesheryctl system provider view` : Lists your current context and provider

<a href="{{ site.baseurl }}/assets/img/syscmd/pro view.png"><img alt="skip-browser" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/syscmd/pro view.png" /></a>

