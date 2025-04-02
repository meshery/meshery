---
layout: default
title: mesheryctl-exp-connections-list
permalink: reference/mesheryctl/exp/connections/list
redirect_from: reference/mesheryctl/exp/connections/list/
type: reference
display-title: "false"
language: en
command: exp
subcommand: connections
---

# mesheryctl exp connections list

List all the connections

## Synopsis

List all the connections
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connections list [flags]

</div>
</pre> 

## Examples

List all the connections
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connections list

</div>
</pre> 

List all the connections with page number
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp connections list --page 2

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count      Display the count of models
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
