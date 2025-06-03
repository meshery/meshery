---
layout: default
title: mesheryctl-exp-connection
permalink: reference/mesheryctl/exp/connection
redirect_from: reference/mesheryctl/exp/connection/
type: reference
display-title: "false"
language: en
command: exp
subcommand: connection
---

# mesheryctl exp connection

Manage Meshery connection

## Synopsis

View and manage your Meshery connection.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/exp/connection
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connection [flags]

</div>
</pre> 

## Examples

Display total count of all available connections
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connection --count

</div>
</pre> 

List all the connection
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connection list

</div>
</pre> 

Delete a connection
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connection delete [connection_id]

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
      --config string   path to config file (default "/home/n2/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

* [mesheryctl exp connection delete](/reference/mesheryctl/exp/connection/delete)
* [mesheryctl exp connection list](/reference/mesheryctl/exp/connection/list)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
