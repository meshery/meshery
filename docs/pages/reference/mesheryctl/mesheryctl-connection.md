---
layout: default
title: mesheryctl-connection
permalink: reference/mesheryctl/connection
redirect_from: reference/mesheryctl/connection/
type: reference
display-title: "false"
language: en
command: connection
subcommand: nil
---

# mesheryctl connection

Manage Meshery connections

## Synopsis

View and manage your Meshery connection.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/connection
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection [flags]

</div>
</pre> 

## Examples

Display total count of all available connections
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection --count

</div>
</pre> 

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

List all the connection
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection list

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection list --count

</div>
</pre> 

Delete a connection
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection delete [connection_id]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count   Display the count of total available connections
  -h, --help    help for connection

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

* [mesheryctl connection create](/reference/mesheryctl/connection/create)
* [mesheryctl connection delete](/reference/mesheryctl/connection/delete)
* [mesheryctl connection list](/reference/mesheryctl/connection/list)
* [mesheryctl connection view](/reference/mesheryctl/connection/view)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
