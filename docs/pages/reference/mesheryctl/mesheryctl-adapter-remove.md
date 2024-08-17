---
layout: default
title: mesheryctl-adapter-remove
permalink: reference/mesheryctl/adapter/remove
redirect_from: reference/mesheryctl/adapter/remove/
type: reference
display-title: "false"
language: en
command: adapter
subcommand: remove
---

# mesheryctl adapter remove

remove cloud and cloud native infrastructure

## Synopsis

remove cloud and cloud native infrastructure
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl adapter remove [flags]

</div>
</pre> 

## Examples

Remove Linkerd deployment
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl adapter remove linkerd

</div>
</pre> 

Remove a Linkerd control plane found under a specific namespace (linkerd-ns)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl adapter remove linkerd --namespace linkerd-ns

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
