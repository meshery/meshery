---
layout: default
title: mesheryctl-system-config
permalink: reference/mesheryctl/system/config
redirect_from: reference/mesheryctl/system/config/
type: reference
display-title: "false"
language: en
command: system
subcommand: config
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

Set configuration according to k8s cluster
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system config [aks|eks|gke|minikube]

</div>
</pre> 

Path to token for authenticating to Meshery API (optional, can be done alternatively using "login")
<pre class='codeblock-pre'>
<div class='codeblock'>
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
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
