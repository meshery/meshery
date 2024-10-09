---
layout: default
title: mesheryctl-exp-connections
permalink: reference/mesheryctl/exp/connections
redirect_from: reference/mesheryctl/exp/connections/
type: reference
display-title: "false"
language: en
command: exp
subcommand: connections
---

# mesheryctl exp connections

Manage Meshery connections

## Synopsis

View and manage your Meshery connections
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connections [flags]

</div>
</pre> 

## Examples

List all the connections
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connections list

</div>
</pre> 

Delete a connection
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connections delete [connection_id]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for connections

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
