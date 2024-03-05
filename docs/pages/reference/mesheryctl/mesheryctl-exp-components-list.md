---
layout: default
title: mesheryctl-exp-components-list
permalink: reference/mesheryctl/exp/components/list
redirect_from: reference/mesheryctl/exp/components/list/
type: reference
display-title: "false"
language: en
command: exp
subcommand: components
---

# mesheryctl exp components list

List registered components

## Synopsis

List all components registered in Meshery Server
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp components list [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
	// View list of components

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp components list

</div>
</pre> 

View list of components with specified page number (25 components per page)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp components list --page 2

</div>
</pre> 

View Total number of components
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp components list --count

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count      (optional) Get the number of components in total
  -h, --help       help for list
  -p, --page int   (optional) List next set of models with --page (default = 1) (default 1)

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
