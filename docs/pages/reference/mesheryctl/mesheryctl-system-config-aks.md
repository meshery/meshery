---
layout: default
title: mesheryctl-system-config-aks
permalink: reference/mesheryctl/system/config/aks
redirect_from: reference/mesheryctl/system/config/aks/
type: reference
display-title: "false"
language: en
command: system
subcommand: config
---

# mesheryctl system config aks

Configure Meshery to use AKS cluster

## Synopsis

Configure Meshery to connect to AKS cluster

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system config aks [flags]

</div>
</pre> 

## Examples

Configure Meshery to connect to AKS cluster using auth token
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system config aks --token auth.json

</div>
</pre> 

Configure Meshery to connect to AKS cluster (if session is logged in using login subcommand)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system config aks

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help           help for aks
  -t, --token string   Path to token for authenticating to Meshery API

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
