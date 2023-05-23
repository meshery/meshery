---
layout: default
title: mesheryctl-system-start
permalink: reference/mesheryctl/system/start
redirect_from: reference/mesheryctl/system/start/
type: reference
display-title: "false"
language: en
command: system
subcommand: start
---

# mesheryctl system start

Start Meshery

## Synopsis

Start Meshery and each of its cloud native components.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start [flags]

</div>
</pre> 

## Examples

Start meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start

</div>
</pre> 

To create a new context for in-cluster Kubernetes deployments and set the new context as your current-context
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context create k8s -p kubernetes -s

</div>
</pre> 

(optional) skip checking for new updates available in Meshery.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start --skip-update

</div>
</pre> 

Reset Meshery's configuration file to default settings.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start --reset

</div>
</pre> 

Silently create Meshery's configuration file with default settings
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start --yes

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help              help for start
  -p, --platform string   platform to deploy Meshery to.
      --reset             (optional) reset Meshery's configuration file to default settings.
      --skip-browser      (optional) skip opening of MesheryUI in browser.
      --skip-update       (optional) skip checking for new Meshery's container images.

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
