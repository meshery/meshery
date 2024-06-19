---
layout: default
title: mesheryctl-exp-relationship-list
permalink: reference/mesheryctl/exp/relationship/list
redirect_from: reference/mesheryctl/exp/relationship/list/
type: reference
display-title: "false"
language: en
command: exp
subcommand: relationship
---

# mesheryctl exp relationship list

List registered relationships

## Synopsis

List all relationships registered in Meshery Server
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationship list [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
	View list of relationship

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    mesheryctl exp relationship list

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    View list of relationship with specified page number (25 relationships per page)

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    mesheryctl exp relationship list --page 2

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help       help for list
  -p, --page int   (optional) List next set of relationships with --page (default = 1) (default 1)

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
