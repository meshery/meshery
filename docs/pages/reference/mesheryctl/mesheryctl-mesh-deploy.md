---
layout: default
title: mesheryctl-mesh-deploy
permalink: reference/mesheryctl/mesh/deploy
redirect_from: reference/mesheryctl/mesh/deploy/
type: reference
display-title: "false"
language: en
command: mesh
subcommand: deploy
---

# mesheryctl mesh deploy

Deploy infrastructure to the Kubernetes cluster

## Synopsis

Deploy infrastructure to the connected Kubernetes cluster
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh deploy [flags]

</div>
</pre> 

## Examples

Deploy a infrastructure from an interactive on the default namespace
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh deploy

</div>
</pre> 

Deploy infrastructure
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh deploy linkerd

</div>
</pre> 

Deploy Linkerd mesh on a specific namespace
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh deploy linkerd --namespace linkerd-ns

</div>
</pre> 

Deploy Linkerd mesh and wait for it to be deployed
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh deploy linkerd --watch

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
		

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help               help for deploy
  -n, --namespace string   Kubernetes namespace to be used for deploying the validation tests and sample workload (default "default")
  -t, --token string       Path to token for authenticating to Meshery API
  -w, --watch              Watch for events and verify operation (in beta testing)

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## Screenshots

Usage of mesheryctl mesh deploy
![mesh-deploy-usage](/assets/img/mesheryctl/deploy-mesh.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
