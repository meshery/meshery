---
layout: default
title: mesheryctl-exp-relationship
permalink: reference/mesheryctl/exp/relationship
redirect_from: reference/mesheryctl/exp/relationship/
type: reference
display-title: "false"
language: en
command: exp
subcommand: relationship
---

# mesheryctl exp relationship

View list of relationships and details of relationship

## Synopsis

Meshery uses relationships to define how interconnected components interact. View list of relationships and detailed information of a specific relationship
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationship [flags]

</div>
</pre> 

## Examples

To view list of relationships
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationships list

</div>
</pre> 

To view a specific relationship
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationships view [model-name]

</div>
</pre> 

//To search a specific relationship
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationships search --[flag] [query-text]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for relationship

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
