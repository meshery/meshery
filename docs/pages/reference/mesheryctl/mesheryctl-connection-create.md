---
layout: default
title: mesheryctl-connection-create
permalink: reference/mesheryctl/connection/create
redirect_from: reference/mesheryctl/connection/create/
type: reference
display-title: "false"
language: en
command: connection
subcommand: create
---

# mesheryctl connection create

Create a new connection

## Synopsis

Create a new connection to a Kubernetes cluster or other supported platform
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection create [flags]

</div>
</pre> 

## Examples

Create a new Kubernetes connection using a specific type
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection create --type aks

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection create --type eks

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection create --type gke

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection create --type minikube

</div>
</pre> 

Create a connection with a token
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection create --type gke --token auth.json

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help           help for create
      --token string   Path to token for authenticating to Meshery API
  -t, --type string    Type of connection to create (aks|eks|gke|minikube)

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

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
