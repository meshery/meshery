---
layout: default
title: mesheryctl-mesh
permalink: /reference/mesheryctl/mesh/
redirect_from: /reference/mesheryctl/mesh/
type: reference
display-title: "false"
language: en
command: mesh
---

# mesheryctl mesh

Service Mesh Lifecycle Management

## Synopsis

Provisioning, configuration, and on-going operational management of service meshes

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl mesh [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>


	// Lifecycle management of service meshes
	mesheryctl mesh [subcommand] 
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for mesh

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

* [mesheryctl](/reference/mesheryctl/main)	 - Meshery Command Line tool
* [mesheryctl mesh deploy](deploy/)	 - Deploy a service mesh to the Kubernetes cluster
* [mesheryctl mesh remove](remove/)	 - remove a service mesh in the kubernetes cluster
* [mesheryctl mesh validate](validate/)	 - Validate conformance to service mesh standards

