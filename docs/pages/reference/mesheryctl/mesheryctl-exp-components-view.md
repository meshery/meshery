---
layout: default
title: mesheryctl-exp-components-view
permalink: reference/mesheryctl/exp/components/view
redirect_from: reference/mesheryctl/exp/components/view/
type: reference
display-title: "false"
language: en
command: exp
subcommand: components
---

# mesheryctl exp components view

view registered components

## Synopsis

view a component registered in Meshery Server
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp components view [flags]

</div>
</pre> 

## Examples

View details of a specific component
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp components view [component-name]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                   help for view
  -o, --output-format string   (optional) format to display in [json|yaml] (default "yaml")

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
