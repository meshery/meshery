---
layout: default
title: mesheryctl-connection-list
permalink: reference/mesheryctl/connection/list
redirect_from: reference/mesheryctl/connection/list/
type: reference
display-title: "false"
language: en
command: connection
subcommand: list
---

# mesheryctl connection list

List all the connections

## Synopsis

List all available connections.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/connection/list
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection list [flags]

</div>
</pre> 

## Examples

List all the connections
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection list

</div>
</pre> 

List all the connections with page number
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection list --page [page-number]

</div>
</pre> 

List all the connections matching a specific kind and status
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection list --kind [kind] --status [status]

</div>
</pre> 

List all the connections matching a set of kinds and statuses
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection list --kind [kind] --kind [kind] --status [status] --status [status]

</div>
</pre> 

Display total count of all available connections
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection list --count

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count            Display the count of total available connections
  -h, --help             help for list
  -k, --kind strings     Filter connections by kind
  -p, --page int         Page number (default 1)
      --pagesize int     Number of connections per page (default 10)
  -s, --status strings   Filter connections by status

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
