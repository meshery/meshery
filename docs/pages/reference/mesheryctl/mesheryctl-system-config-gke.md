---
layout: default
title: mesheryctl-system-config-gke
permalink: /reference/mesheryctl/system/config/gke/
redirect_from: /reference/mesheryctl/system/config/gke/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system config gke

Configure Meshery to use GKE cluster

## Synopsis

Configure Meshery to connect to GKE cluster

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system config gke [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>


	Configure Meshery to connect to GKE cluster using auth token
	mesheryctl system config gke --token auth.json
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help           help for gke
  -t, --token string   Path to token for authenticating to Meshery API

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

* [mesheryctl system config](config/)	 - Configure Meshery

