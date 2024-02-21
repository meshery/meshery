---
layout: default
title: mesheryctl-exp-components
permalink: reference/mesheryctl/exp/components
redirect_from: reference/mesheryctl/exp/components/
type: reference
display-title: "false"
language: en
command: exp
subcommand: components
---

# mesheryctl exp components

View list of components and detail of components

## Synopsis

View list of components and detailed information of a specific component
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp components [flags]

</div>
</pre> 

## Examples

To view list of components
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp components list

</div>
</pre> 

To view a specific component
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp components view [component-name]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for components

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
