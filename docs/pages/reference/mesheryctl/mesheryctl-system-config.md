---
layout: default
title: mesheryctl-system-config
permalink: /reference/mesheryctl/system/config/
redirect_from: /reference/mesheryctl/system/config/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system config

Configure Meshery

## Synopsis

Configure the Kubernetes cluster used by Meshery.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system config [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>

	// Set configuration according to k8s cluster
	mesheryctl system config [aks|eks|gke|minikube]

	// Path to token for authenticating to Meshery API (optional)
	mesheryctl system config --token "~/Downloads/auth.json"
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for config

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "/home/admin-pc/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

* [mesheryctl system](system/)	 - Meshery Lifecycle Management
* [mesheryctl system config aks](config/aks/)	 - Configure Meshery to use AKS cluster
* [mesheryctl system config eks](config/eks/)	 - Configure Meshery to use EKS cluster
* [mesheryctl system config gke](config/gke/)	 - Configure Meshery to use GKE cluster
* [mesheryctl system config minikube](config/minikube/)	 - Configure Meshery to use minikube cluster

