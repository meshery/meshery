---
layout: default
title: mesheryctl-exp-connection-list
permalink: reference/mesheryctl/exp/connection/list
redirect_from: reference/mesheryctl/exp/connection/list/
type: reference
display-title: "false"
language: en
command: exp
subcommand: connection
---

# mesheryctl exp connection list

List all the connections

## Synopsis

List all available connections.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/exp/connection/list
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connection list [flags]

</div>
</pre> 

## Examples

List all the connections
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connection list

</div>
</pre> 

List all the connections with page number
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connection list --page [page-number]

</div>
</pre> 

Display total count of all available connections
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connection list --count

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count      Display the count of total available connections
  -h, --help       help for list
  -p, --page int   Page number (default 1)

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
