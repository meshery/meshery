---
layout: default
title: mesheryctl-system-check
permalink: reference/mesheryctl/system/check
redirect_from: reference/mesheryctl/system/check/
type: reference
display-title: "false"
language: en
command: system
subcommand: check
---

# mesheryctl system check

Meshery environment check

## Synopsis

Verify environment pre/post-deployment of Meshery.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system check [flags]

</div>
</pre> 

## Examples

Run system checks for both pre and post mesh deployment scenarios on Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system check

</div>
</pre> 

Run Pre-mesh deployment checks (Docker and Kubernetes)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system check --preflight

</div>
</pre> 

Run checks on specific mesh adapter
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system check --adapter meshery-istio:10000

</div>
</pre> 

or
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system check --adapter meshery-istio

</div>
</pre> 

Run checks for all the mesh adapters
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system check --adapters

</div>
</pre> 

Verify the health of Meshery Operator's deployment with MeshSync and Broker
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system check --operator

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --adapter string   Check status of specified meshery adapter
      --adapters         Check status of meshery adapters
      --components       Check status of Meshery components
  -h, --help             help for check
      --operator         Verify the health of Meshery Operator's deployment with MeshSync and Broker
      --pre              Verify environment readiness to deploy Meshery
      --preflight        Verify environment readiness to deploy Meshery

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

## Screenshots

Usage of mesheryctl system check
![check-usage](/assets/img/mesheryctl/check.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
