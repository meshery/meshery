---
layout: default
title: mesheryctl-mesh-remove
permalink: reference/mesheryctl/mesh/remove
redirect_from: reference/mesheryctl/mesh/remove/
type: reference
display-title: "false"
language: en
command: mesh
subcommand: remove
---

# mesheryctl mesh remove

remove a service mesh in the kubernetes cluster

## Synopsis

remove service mesh in the connected kubernetes cluster

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh remove [flags]

</div>
</pre> 

## Examples

Remove a service mesh
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh remove [mesh adapter name]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
		

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --adapter string     Adapter to use for installation (default "meshery-istio:10000")
  -h, --help               help for remove
  -n, --namespace string   Kubernetes namespace to be used for deploying the validation tests and sample workload (default "default")
  -t, --token string       Path to token for authenticating to Meshery API

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
