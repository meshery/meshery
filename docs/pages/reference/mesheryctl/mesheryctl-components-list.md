---
layout: default
title: mesheryctl-components-list
permalink: reference/mesheryctl/components/list
redirect_from: reference/mesheryctl/components/list/
type: reference
display-title: "false"
language: en
command: components
subcommand: list
---

# mesheryctl components list

List registered components

## Synopsis

List all components registered in Meshery Server
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl components list [flags]

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

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
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
