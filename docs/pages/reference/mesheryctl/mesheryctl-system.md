---
layout: default
title: mesheryctl-system
permalink: /reference/mesheryctl/system/
redirect_from: /reference/mesheryctl/system/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system

Meshery Lifecycle Management

## Synopsis

Manage the state and configuration of Meshery server, components, and client.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system [flags]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --context string   (optional) temporarily change the current context.
  -h, --help             help for system
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/admin-pc/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

* [mesheryctl](reference/mesheryctl/main)	 - Meshery Command Line tool
* [mesheryctl system channel](channel/)	 - Switch between release channels
* [mesheryctl system check](check/)	 - Meshery environment check
* [mesheryctl system completion](completion/)	 - Output shell completion code
* [mesheryctl system config](config/)	 - Configure Meshery
* [mesheryctl system context](context/)	 - Configure your Meshery deployment(s)
* [mesheryctl system dashboard](dashboard/)	 - Open Meshery UI in browser.
* [mesheryctl system login](login/)	 - Authenticate to a Meshery Server
* [mesheryctl system logout](logout/)	 - Remove authentication for Meshery Server
* [mesheryctl system logs](logs/)	 - Print logs
* [mesheryctl system reset](reset/)	 - Reset Meshery's configuration
* [mesheryctl system restart](restart/)	 - Stop, then start Meshery
* [mesheryctl system start](start/)	 - Start Meshery
* [mesheryctl system status](status/)	 - Check Meshery status
* [mesheryctl system stop](stop/)	 - Stop Meshery
* [mesheryctl system token](token/)	 - Manage Meshery user tokens
* [mesheryctl system update](update/)	 - Pull new Meshery images/manifest files.

