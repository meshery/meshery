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

Deploy a service mesh to the Kubernetes cluster

## Synopsis

Deploy a service mesh to the connected Kubernetes cluster

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh deploy [flags]

</div>
</pre> 

## Examples

Deploy a service mesh from an interactive on the default namespace
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh deploy

</div>
</pre> 

Deploy Linkerd mesh on a specific namespace
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh deploy --adapter meshery-linkerd --namespace linkerd-ns

</div>
</pre> 

Deploy Linkerd mesh and wait for it to be deployed
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh deploy --adapter meshery-linkerd --watch

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

Go back to [command reference index](/reference/mesheryctl/) 
