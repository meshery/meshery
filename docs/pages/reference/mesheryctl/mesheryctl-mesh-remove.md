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

Remove a service mesh(linkerd)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh remove linkerd

</div>
</pre> 

Remove a service mesh(linkerd) under a specific namespace(linkerd-ns)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh remove linkerd --namespace linkerd-ns

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
		

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help               help for remove
  -n, --namespace string   Kubernetes namespace where the mesh is deployed (default "default")

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -t, --token string    Path to token for authenticating to Meshery API
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
