---
layout: default
title: mesheryctl-system-config-gke
permalink: reference/mesheryctl/system/config/gke
redirect_from: reference/mesheryctl/system/config/gke/
type: reference
display-title: "false"
language: en
command: system
subcommand: config
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

Configure Meshery to connect to GKE cluster using auth token
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system config gke --token auth.json

</div>
</pre> 

Configure Meshery to connect to GKE cluster (if session is logged in using login subcommand)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system config gke

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
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
