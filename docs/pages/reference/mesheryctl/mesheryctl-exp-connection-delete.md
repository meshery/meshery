---
layout: default
title: mesheryctl-exp-connection-delete
permalink: reference/mesheryctl/exp/connection/delete
redirect_from: reference/mesheryctl/exp/connection/delete/
type: reference
display-title: "false"
language: en
command: exp
subcommand: connection
---

# mesheryctl exp connection delete

Delete a connection

## Synopsis

Delete a connection providing the connection ID.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/exp/connection/delete
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connection delete [flags]

</div>
</pre> 

## Examples

Delete a specific connection
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connection delete [connection_id]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for delete

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
